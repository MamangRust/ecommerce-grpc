package api_test

import (
	"bytes"
	"ecommerce/internal/cache"
	api_banner_cache "ecommerce/internal/cache/api/banner"
	banner_cache "ecommerce/internal/cache/banner"
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

type BannerApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	echo        *echo.Echo
	client      pb.BannerServiceClient
	conn        *grpc.ClientConn
}

func (s *BannerApiTestSuite) SetupSuite() {
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
	log, _ := logger.NewLogger("test-banner-api", lp)
	obs, _ := observability.NewObservability("test-banner-api", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-banner-api")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	
	bannerCacheSrv := banner_cache.NewBannerMencache(cacheStore)
	bannerCacheApi := api_banner_cache.NewBannerMencache(cacheStore)

	bannerService := service.NewBannerService(service.BannerServiceDeps{
		BannerRepository: repos.Banner,
		Logger:           log,
		Observability:      obs,
		Cache:              bannerCacheSrv,
	})

	// Start gRPC Server
	bannerGapi := gapi.NewBannerHandleGrpc(bannerService)
	server := grpc.NewServer()
	pb.RegisterBannerServiceServer(server, bannerGapi)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	// gRPC Client for the API Handler
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewBannerServiceClient(conn)

	// Echo Setup
	s.echo = echo.New()
	mapping := response_api.NewBannerResponseMapper()
	apiHandler := errors.NewApiHandler(obs, log)

	api.NewHandleBanner(s.echo, s.client, log, mapping, apiHandler, bannerCacheApi)
}

func (s *BannerApiTestSuite) TearDownSuite() {
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

func (s *BannerApiTestSuite) TestBannerApiLifecycle() {
	// 1. Create
	createReq := map[string]interface{}{
		"name":       "Winter Sale API",
		"start_date": "2026-12-01",
		"end_date":   "2026-12-31",
		"start_time": "00:00",
		"end_time":   "23:59",
		"is_active":  true,
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/banner/create", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var createRes map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &createRes)
	
	data := createRes["data"].(map[string]interface{})
	bannerID := int(data["id"].(float64))

	// 2. Find By ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/banner/%d", bannerID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Update
	updateReq := map[string]interface{}{
		"banner_id":  bannerID,
		"name":       "Winter Sale API Updated",
		"start_date": "2026-12-01",
		"end_date":   "2026-12-31",
		"start_time": "00:00",
		"end_time":   "23:59",
		"is_active":  true,
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/banner/update/%d", bannerID), bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/banner/trashed/%d", bannerID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/banner/restore/%d", bannerID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Delete Permanent
	// Trash again
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/banner/trashed/%d", bannerID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/banner/permanent/%d", bannerID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestBannerApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(BannerApiTestSuite))
}
