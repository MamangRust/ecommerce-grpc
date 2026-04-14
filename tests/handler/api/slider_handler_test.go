package api_test

import (
	"context"
	"ecommerce/internal/cache"
	api_cache "ecommerce/internal/cache/api/slider"
	service_cache "ecommerce/internal/cache/slider"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/handler/api"
	"ecommerce/internal/handler/gapi"
	mapper "ecommerce/internal/mapper"
	"ecommerce/internal/pb"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"ecommerce/pkg/upload_image"
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
	"google.golang.org/grpc/test/bufconn"
)

type SliderApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	repos       *repository.Repositories
	echo        *echo.Echo
	gRpcConn    *grpc.ClientConn
	listener    *bufconn.Listener
}

func (s *SliderApiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	// DB Setup
	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	// Redis Setup
	opt, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opt)

	// Repositories
	queries := db.New(pool)
	s.repos = repository.NewRepositories(queries)

	// Logging & Observability
	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	l, err := logger.NewLogger("test-slider-api", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-slider-api", l)
	s.Require().NoError(err)

	// Cache Setup
	cacheMetrics, err := observability.NewCacheMetrics("test-slider-api")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.redisClient, l, cacheMetrics)
	slServiceCache := service_cache.NewSliderMencache(cacheStore)
	slApiCache := api_cache.NewSliderMencache(cacheStore)

	// Service
	slService := service.NewSliderService(service.SliderServiceDeps{
		SliderRepository: s.repos.Slider,
		Logger:           l,
		Observability:    obs,
		Cache:            slServiceCache,
	})

	// gRPC Server Setup
	s.listener = bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	hGrpc := gapi.NewSliderHandleGrpc(slService)
	pb.RegisterSliderServiceServer(server, hGrpc)

	go func() {
		if err := server.Serve(s.listener); err != nil {
			panic(err)
		}
	}()

	// gRPC Client Setup for API Handler
	conn, err := grpc.DialContext(context.Background(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return s.listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err)
	s.gRpcConn = conn
	slClient := pb.NewSliderServiceClient(conn)

	// API Handler Setup
	s.echo = echo.New()
	slMapper := mapper.NewSliderResponseMapper()
	apiHandler := errors.NewApiHandler(obs, l)
	uploader := upload_image.NewImageUpload(l)
	
	api.NewHandlerSlider(s.echo, slClient, l, slMapper, uploader, apiHandler, slApiCache)
}

func (s *SliderApiTestSuite) TearDownSuite() {
	if s.gRpcConn != nil {
		s.gRpcConn.Close()
	}
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *SliderApiTestSuite) TestSliderApiLifecycle() {
	// 1. Create via Repository
	ctx := context.Background()
	createReq := &requests.CreateSliderRequest{
		Nama:     "Api Test Slider",
		FilePath: "path/to/api_image.jpg",
	}

	slider, err := s.repos.Slider.CreateSlider(ctx, createReq)
	s.NoError(err)
	s.NotNil(slider)
	sliderID := int(slider.SliderID)

	// 2. Find All
	req := httptest.NewRequest(http.MethodGet, "/api/slider?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	var res map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	s.NoError(err)
	s.Equal("success", res["status"])

	// 3. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/slider/trashed/%d", sliderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/slider/restore/%d", sliderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Delete Permanent
	// Trash again
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/slider/trashed/%d", sliderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/slider/permanent/%d", sliderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestSliderApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SliderApiTestSuite))
}
