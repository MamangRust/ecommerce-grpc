package api_test

import (
	"bytes"
	"ecommerce/internal/cache"
	api_product_cache "ecommerce/internal/cache/api/product"
	product_cache "ecommerce/internal/cache/product"
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
	"time"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	echo        *echo.Echo
	client      pb.ProductServiceClient
	conn        *grpc.ClientConn

	// IDs for reference
	merchantID int32
	categoryID int32
}

func (s *ProductApiTestSuite) SetupSuite() {
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
	log, _ := logger.NewLogger("test-product-api", lp)
	obs, _ := observability.NewObservability("test-product-api", log)

	cacheMetrics, _ := observability.NewCacheMetrics("test-product-api")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)

	// Service layer cache
	prodCacheSrv := product_cache.NewProductMencache(cacheStore)
	// API layer cache
	prodCacheApi := api_product_cache.NewProductMencache(cacheStore)

	productService := service.NewProductService(service.ProductServiceDeps{
		CategoryRepository: repos.Category,
		MerchantRepository: repos.Merchant,
		ProductRepository:  repos.Product,
		Logger:             log,
		Observability:      obs,
		Cache:              prodCacheSrv,
	})

	// Start gRPC Server
	productGapi := gapi.NewProductHandleGrpc(productService)
	server := grpc.NewServer()
	pb.RegisterProductServiceServer(server, productGapi)
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
	s.client = pb.NewProductServiceClient(conn)

	// Echo Setup
	s.echo = echo.New()
	mapping := response_api.NewProductResponseMapper()
	imgUpload := upload_image.NewImageUpload(log)
	apiHandler := errors.NewApiHandler(obs, log)

	api.NewHandlerProduct(s.echo, s.client, log, mapping, imgUpload, apiHandler, prodCacheApi)

	// Setup prerequisites directly via DB
	ctx := context.Background()
	var userID int32
	err = pool.QueryRow(ctx, "INSERT INTO users (firstname, lastname, email, password) VALUES ($1, $2, $3, $4) RETURNING user_id",
		"Api", "User", "api@example.com", "password").Scan(&userID)
	s.Require().NoError(err)

	err = pool.QueryRow(ctx, "INSERT INTO merchants (user_id, name, description) VALUES ($1, $2, $3) RETURNING merchant_id",
		userID, "Api Merchant", "Api Description").Scan(&s.merchantID)
	s.Require().NoError(err)

	err = pool.QueryRow(ctx, "INSERT INTO categories (name, description, slug_category) VALUES ($1, $2, $3) RETURNING category_id",
		"Api Category", "Api Category Description", "api-category").Scan(&s.categoryID)
	s.Require().NoError(err)
}

func (s *ProductApiTestSuite) TearDownSuite() {
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

func (s *ProductApiTestSuite) TestProductApiLifecycle() {
	timestamp := time.Now().UnixNano()
	slug := fmt.Sprintf("api-product-%d", timestamp)

	// 1. Create (Multipart Form)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("merchant_id", fmt.Sprintf("%d", s.merchantID))
	_ = writer.WriteField("category_id", fmt.Sprintf("%d", s.categoryID))
	_ = writer.WriteField("name", "Api Product")
	_ = writer.WriteField("description", "Via Echo API")
	_ = writer.WriteField("price", "500")
	_ = writer.WriteField("count_in_stock", "50")
	_ = writer.WriteField("brand", "ApiBrand")
	_ = writer.WriteField("weight", "200")
	_ = writer.WriteField("rating", "5")
	_ = writer.WriteField("slug_product", slug)

	part, _ := writer.CreateFormFile("image_product", "api_test.jpg")
	_, _ = part.Write([]byte("api test image content"))
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/product/create", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusCreated, rec.Code)
	var createRes map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &createRes)
	s.NoError(err)
	s.Equal("success", createRes["status"])

	data := createRes["data"].(map[string]interface{})
	productID := int(data["id"].(float64))

	// 2. Find By ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/product/%d", productID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Find All
	req = httptest.NewRequest(http.MethodGet, "/api/product?page=1&page_size=10", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Update
	body = new(bytes.Buffer)
	writer = multipart.NewWriter(body)
	_ = writer.WriteField("merchant_id", fmt.Sprintf("%d", s.merchantID))
	_ = writer.WriteField("category_id", fmt.Sprintf("%d", s.categoryID))
	_ = writer.WriteField("name", "Api Product Updated")
	_ = writer.WriteField("description", "Updated via Echo API")
	_ = writer.WriteField("price", "550")
	_ = writer.WriteField("count_in_stock", "45")
	_ = writer.WriteField("brand", "ApiBrand v2")
	_ = writer.WriteField("weight", "210")
	_ = writer.WriteField("rating", "4")
	_ = writer.WriteField("slug_product", slug+"-updated")
	_ = writer.Close()

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/product/update/%d", productID), body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/product/trashed/%d", productID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/product/restore/%d", productID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. Delete Permanent
	// Trash again
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/product/trashed/%d", productID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/product/permanent/%d", productID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestProductApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductApiTestSuite))
}
