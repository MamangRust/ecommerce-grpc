package api_test

import (
	"bytes"
	"context"
	"ecommerce/internal/cache"
	api_merchant_cache "ecommerce/internal/cache/api/merchant"
	merchant_cache "ecommerce/internal/cache/merchant"
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

type MerchantApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	echo        *echo.Echo
	client      pb.MerchantServiceClient
	conn        *grpc.ClientConn
	userID      int32
}

func (s *MerchantApiTestSuite) SetupSuite() {
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
	log, _ := logger.NewLogger("test-merchant-api", lp)
	obs, _ := observability.NewObservability("test-merchant-api", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-merchant-api")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	
	merchCacheSrv := merchant_cache.NewMerchantMencache(cacheStore)
	merchCacheApi := api_merchant_cache.NewMerchantMencache(cacheStore)

	merchantService := service.NewMerchantService(service.MerchantServiceDeps{
		MerchantRepository: repos.Merchant,
		Logger:             log,
		Observability:      obs,
		Cache:              merchCacheSrv,
	})

	// Start gRPC Server
	merchantGapi := gapi.NewMerchantHandleGrpc(merchantService)
	server := grpc.NewServer()
	pb.RegisterMerchantServiceServer(server, merchantGapi)
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
	s.client = pb.NewMerchantServiceClient(conn)

	// Echo Setup
	s.echo = echo.New()
	mapping := response_api.NewMerchantResponseMapper()
	apiHandler := errors.NewApiHandler(obs, log)

	api.NewHandlerMerchant(s.echo, s.client, log, mapping, apiHandler, merchCacheApi)

	// Prerequisite: Create User
	ctx := context.Background()
	err = pool.QueryRow(ctx, "INSERT INTO users (firstname, lastname, email, password) VALUES ($1, $2, $3, $4) RETURNING user_id",
		"Merchant", "Api", "merchant-api@example.com", "password").Scan(&s.userID)
	s.Require().NoError(err)
}

func (s *MerchantApiTestSuite) TearDownSuite() {
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

func (s *MerchantApiTestSuite) TestMerchantApiLifecycle() {
	// 1. Create (POST JSON)
	createBody := map[string]interface{}{
		"user_id":       s.userID,
		"name":          "API Merchant",
		"description":   "API desc",
		"address":       "API Addr",
		"contact_email": "api@merchant.com",
		"contact_phone": "08123456789",
		"status":       "active",
	}
	bodyBytes, _ := json.Marshal(createBody)

	req := httptest.NewRequest(http.MethodPost, "/api/merchant/create", bytes.NewBuffer(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var createRes map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &createRes)
	
	data := createRes["data"].(map[string]interface{})
	merchantID := int(data["id"].(float64))

	// 2. Find By ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant/%d", merchantID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Update
	updateBody := make(map[string]interface{})
	for k, v := range createBody {
		updateBody[k] = v
	}
	updateBody["merchant_id"] = merchantID
	updateBody["name"] = "API Merchant Updated"
	bodyBytes, _ = json.Marshal(updateBody)

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant/update/%d", merchantID), bytes.NewBuffer(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant/trashed/%d", merchantID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant/restore/%d", merchantID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Delete Permanent
	_ = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant/trashed/%d", merchantID), nil)
	// (Actually call it)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant/trashed/%d", merchantID), nil))

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/merchant/permanent/%d", merchantID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestMerchantApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantApiTestSuite))
}
