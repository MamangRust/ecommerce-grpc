package api_test

import (
	"bytes"
	"context"
	"ecommerce/internal/cache"
	api_merchantpolicies_cache "ecommerce/internal/cache/api/merchant_policies"
	merchantpolicies_cache "ecommerce/internal/cache/merchant_policies"
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

type MerchantPolicyApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	echo        *echo.Echo
	client      pb.MerchantPoliciesServiceClient
	conn        *grpc.ClientConn
	merchantID  int32
}

func (s *MerchantPolicyApiTestSuite) SetupSuite() {
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
	log, _ := logger.NewLogger("test-merchant-policy-api", lp)
	obs, _ := observability.NewObservability("test-merchant-policy-api", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-merchant-policy-api")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	
	policyCacheSrv := merchantpolicies_cache.NewMerchantPoliciesMencache(cacheStore)
	policyCacheApi := api_merchantpolicies_cache.NewMerchantPoliciesMencache(cacheStore)

	policyService := service.NewMerchantPoliciesService(service.MerchantPoliciesServiceDeps{
		MerchantPoliciesRepository: repos.MerchantPolicies,
		Logger:                     log,
		Observability:              obs,
		Cache:                      policyCacheSrv,
	})

	// Start gRPC Server
	policyGapi := gapi.NewMerchantPolicyHandleGrpc(policyService)
	server := grpc.NewServer()
	pb.RegisterMerchantPoliciesServiceServer(server, policyGapi)
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
	s.client = pb.NewMerchantPoliciesServiceClient(conn)

	// Echo Setup
	s.echo = echo.New()
	mapping := response_api.NewMerchantPolicyResponseMapper()
	mappingMerchant := response_api.NewMerchantResponseMapper()
	apiHandler := errors.NewApiHandler(obs, log)

	api.NewHandlerMerchantPolicies(s.echo, s.client, log, mapping, mappingMerchant, apiHandler, policyCacheApi)

	// Prerequisite: Create Merchant
	ctx := context.Background()
	_ = pool.QueryRow(ctx, "INSERT INTO users (firstname, lastname, email, password) VALUES ($1, $2, $3, $4) RETURNING user_id",
		"Policy", "Api", "policy-api@example.com", "password").Scan(new(int32))
	err = pool.QueryRow(ctx, "INSERT INTO merchants (user_id, name) VALUES (1, $1) RETURNING merchant_id", "Policy Merchant").Scan(&s.merchantID)
	s.Require().NoError(err)
}

func (s *MerchantPolicyApiTestSuite) TearDownSuite() {
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

func (s *MerchantPolicyApiTestSuite) TestMerchantPolicyApiLifecycle() {
	// 1. Create
	createBody := map[string]interface{}{
		"merchant_id": s.merchantID,
		"policy_type": "Refund",
		"title":       "30-day refund",
		"description": "Refund within 30 days",
	}
	bodyBytes, _ := json.Marshal(createBody)

	req := httptest.NewRequest(http.MethodPost, "/api/merchant-policy/create", bytes.NewBuffer(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var createRes map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &createRes)
	data := createRes["data"].(map[string]interface{})
	policyID := int(data["id"].(float64))

	// 2. Find By ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant-policy/%d", policyID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Update
	updateBody := make(map[string]interface{})
	for k, v := range createBody {
		updateBody[k] = v
	}
	updateBody["merchant_policy_id"] = policyID
	updateBody["title"] = "Updated Policy"
	bodyBytes, _ = json.Marshal(updateBody)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-policy/update/%d", policyID), bytes.NewBuffer(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-policy/trashed/%d", policyID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestMerchantPolicyApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantPolicyApiTestSuite))
}
