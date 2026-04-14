package api_test

import (
	"bytes"
	"context"
	"ecommerce/internal/cache"
	api_merchantbusiness_cache "ecommerce/internal/cache/api/merchant_business"
	merchantbusiness_cache "ecommerce/internal/cache/merchant_business"
	"ecommerce/internal/handler/api"
	"ecommerce/internal/handler/gapi"
	response_api "ecommerce/internal/mapper"
	"ecommerce/internal/pb"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"ecommerce/tests"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MerchantBusinessApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	echo        *echo.Echo
	client      pb.MerchantBusinessServiceClient
	conn        *grpc.ClientConn
	merchantID  int32
}

func (s *MerchantBusinessApiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)

	queries := db.New(pool)
	repos := repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test-merchant-business-api", lp)
	obs, _ := observability.NewObservability("test-merchant-business-api", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-merchant-business-api")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	
	businessCacheSrv := merchantbusiness_cache.NewMerchantBusinessMencache(cacheStore)
	businessCacheApi := api_merchantbusiness_cache.NewMerchantBusinessMencache(cacheStore)

	businessService := service.NewMerchantBusinessService(service.MerchantBusinessServiceDeps{
		MerchantBusinessRepository: repos.MerchantBusiness,
		Logger:                     log,
		Observability:              obs,
		Cache:                      businessCacheSrv,
	})

	// Start gRPC Server
	businessGapi := gapi.NewMerchantBusinessHandleGrpc(businessService)
	server := grpc.NewServer()
	pb.RegisterMerchantBusinessServiceServer(server, businessGapi)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	// gRPC Client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewMerchantBusinessServiceClient(conn)

	// Echo Setup
	s.echo = echo.New()
	mapping := response_api.NewMerchantBusinessResponseMapper()
	mappingMerchant := response_api.NewMerchantResponseMapper()
	apiHandler := errors.NewApiHandler(obs, log)

	api.NewHandlerMerchantBusiness(s.echo, s.client, log, mapping, mappingMerchant, apiHandler, businessCacheApi)

	// Prerequisite: Create Merchant
	ctx := context.Background()
	_ = pool.QueryRow(ctx, "INSERT INTO users (firstname, lastname, email, password) VALUES ($1, $2, $3, $4) RETURNING user_id",
		"Business", "Api", "business-api@example.com", "password").Scan(new(int32))
	err = pool.QueryRow(ctx, "INSERT INTO merchants (user_id, name) VALUES (1, $1) RETURNING merchant_id", "Business Merchant").Scan(&s.merchantID)
	s.Require().NoError(err)
}

func (s *MerchantBusinessApiTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *MerchantBusinessApiTestSuite) TestMerchantBusinessApiLifecycle() {
	// 1. Create
	createBody := map[string]interface{}{
		"merchant_id":         s.merchantID,
		"business_type":       "Corporation",
		"tax_id":              "12-345-678",
		"established_year":    2010,
		"number_of_employees": 100,
		"website_url":         "https://corp.com",
	}
	bodyBytes, _ := json.Marshal(createBody)

	req := httptest.NewRequest(http.MethodPost, "/api/merchant-business/create", bytes.NewBuffer(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var createRes map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &createRes)
	data := createRes["data"].(map[string]interface{})
	infoID := int(data["id"].(float64))

	// 2. Find By ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant-business/%d", infoID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Update
	updateBody := make(map[string]interface{})
	for k, v := range createBody {
		updateBody[k] = v
	}
	updateBody["merchant_business_info_id"] = infoID
	updateBody["business_type"] = "Partnership"
	bodyBytes, _ = json.Marshal(updateBody)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-business/update/%d", infoID), bytes.NewBuffer(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestMerchantBusinessApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantBusinessApiTestSuite))
}
