package api_test

import (
	"bytes"
	"context"
	"ecommerce/internal/cache"
	api_merchantawards_cache "ecommerce/internal/cache/api/merchant_awards"
	merchantawards_cache "ecommerce/internal/cache/merchant_awards"
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

type MerchantAwardApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	echo        *echo.Echo
	client      pb.MerchantAwardServiceClient
	conn        *grpc.ClientConn
	merchantID  int32
}

func (s *MerchantAwardApiTestSuite) SetupSuite() {
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
	log, _ := logger.NewLogger("test-merchant-award-api", lp)
	obs, _ := observability.NewObservability("test-merchant-award-api", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-merchant-award-api")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	
	awardCacheSrv := merchantawards_cache.NewMerchantAward(cacheStore)
	awardCacheApi := api_merchantawards_cache.NewMerchantAward(cacheStore)

	awardService := service.NewMerchantAwardService(service.MerchantAwardServiceDeps{
		MerchantAwardRepository: repos.MerchantAward,
		Logger:                  log,
		Observability:           obs,
		Cache:                   awardCacheSrv,
	})

	// Start gRPC Server
	awardGapi := gapi.NewMerchantAwardHandleGrpc(awardService)
	server := grpc.NewServer()
	pb.RegisterMerchantAwardServiceServer(server, awardGapi)
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
	s.client = pb.NewMerchantAwardServiceClient(conn)

	// Echo Setup
	s.echo = echo.New()
	mapping := response_api.NewMerchantAwardResponseMapper()
	mappingMerchant := response_api.NewMerchantResponseMapper()
	apiHandler := errors.NewApiHandler(obs, log)

	api.NewHandlerMerchantAward(s.echo, s.client, log, mapping, mappingMerchant, apiHandler, awardCacheApi)

	// Prerequisite: Create Merchant
	ctx := context.Background()
	_ = pool.QueryRow(ctx, "INSERT INTO users (firstname, lastname, email, password) VALUES ($1, $2, $3, $4) RETURNING user_id",
		"Award", "Api", "award-api@example.com", "password").Scan(new(int32))
	err = pool.QueryRow(ctx, "INSERT INTO merchants (user_id, name) VALUES (1, $1) RETURNING merchant_id", "Award Merchant").Scan(&s.merchantID)
	s.Require().NoError(err)
}

func (s *MerchantAwardApiTestSuite) TearDownSuite() {
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

func (s *MerchantAwardApiTestSuite) TestMerchantAwardApiLifecycle() {
	// 1. Create
	createBody := map[string]interface{}{
		"merchant_id":     s.merchantID,
		"title":           "Best API Award",
		"description":     "Award for API quality",
		"issued_by":       "Organization",
		"issue_date":      "2024-01-01",
		"expiry_date":     "2025-01-01",
		"certificate_url": "https://cert.com/1",
	}
	bodyBytes, _ := json.Marshal(createBody)

	req := httptest.NewRequest(http.MethodPost, "/api/merchant-certification/create", bytes.NewBuffer(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var createRes map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &createRes)
	data := createRes["data"].(map[string]interface{})
	awardID := int(data["id"].(float64))

	// 2. Find By ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant-certification/%d", awardID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Update
	updateBody := make(map[string]interface{})
	for k, v := range createBody {
		updateBody[k] = v
	}
	updateBody["merchant_certification_id"] = awardID
	updateBody["title"] = "Updated Award"
	bodyBytes, _ = json.Marshal(updateBody)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-certification/update/%d", awardID), bytes.NewBuffer(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestMerchantAwardApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantAwardApiTestSuite))
}
