package api_test

import (
	"bytes"
	"ecommerce/internal/cache"
	api_category_cache "ecommerce/internal/cache/api/category"
	category_cache "ecommerce/internal/cache/category"
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
	"ecommerce/pkg/upload_image"
	"ecommerce/tests"
	"encoding/json"
	"fmt"
	"mime/multipart"
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

type CategoryApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	echo        *echo.Echo
	client      pb.CategoryServiceClient
	conn        *grpc.ClientConn
}

func (s *CategoryApiTestSuite) SetupSuite() {
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
	log, _ := logger.NewLogger("test-api", lp)
	obs, _ := observability.NewObservability("test-api", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-api")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	
	// Service layer cache
	catCacheSrv := category_cache.NewCategoryMencache(cacheStore)
	// API layer cache
	catCacheApi := api_category_cache.NewCategoryMencache(cacheStore)

	categoryService := service.NewCategoryService(service.CategoryServiceDeps{
		CategoryRepository: repos.Category,
		Logger:             log,
		Observability:      obs,
		Cache:              catCacheSrv,
	})

	// Start gRPC Server
	categoryGapi := gapi.NewCategoryHandleGrpc(categoryService)
	server := grpc.NewServer()
	pb.RegisterCategoryServiceServer(server, categoryGapi)
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
	s.client = pb.NewCategoryServiceClient(conn)

	// Echo Setup
	s.echo = echo.New()
	mapping := response_api.NewCategoryResponseMapper()
	imgUpload := upload_image.NewImageUpload(log)
	apiHandler := errors.NewApiHandler(obs, log)

	api.NewHandlerCategory(s.echo, s.client, log, mapping, imgUpload, apiHandler, catCacheApi)
}

func (s *CategoryApiTestSuite) TearDownSuite() {
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

func (s *CategoryApiTestSuite) TestCategoryApiLifecycle() {
	// 1. Create (Multipart Form)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", "Electronics API")
	_ = writer.WriteField("description", "Via Echo")
	_ = writer.WriteField("slug_category", "electronics-api")
	
	part, _ := writer.CreateFormFile("image_category", "test.jpg")
	_, _ = part.Write([]byte("test image content"))
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/category/create", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var createRes map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &createRes)
	s.NoError(err)
	s.Equal("success", createRes["status"])
	
	data := createRes["data"].(map[string]interface{})
	categoryID := int(data["id"].(float64))

	// 2. Find By ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/category/%d", categoryID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Update (Multipart Form)
	body = new(bytes.Buffer)
	writer = multipart.NewWriter(body)
	_ = writer.WriteField("name", "Electronics API Updated")
	_ = writer.WriteField("description", "Updated via Echo")
	_ = writer.WriteField("slug_category", "electronics-api-updated")
	
	part, _ = writer.CreateFormFile("image_category", "updated.jpg")
	_, _ = part.Write([]byte("updated image content"))
	_ = writer.Close()

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/category/update/%d", categoryID), body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/category/trashed/%d", categoryID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/category/restore/%d", categoryID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Delete Permanent
	// Trash again
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/category/trashed/%d", categoryID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/category/permanent/%d", categoryID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestCategoryApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryApiTestSuite))
}
