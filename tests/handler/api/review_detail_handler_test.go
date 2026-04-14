package api_test

import (
	"bytes"
	"context"
	"ecommerce/internal/cache"
	api_reviewdetail_cache "ecommerce/internal/cache/api/review_detail"
	reviewdetail_cache "ecommerce/internal/cache/review_detail"
	"ecommerce/internal/domain/requests"
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
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReviewDetailApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	echo        *echo.Echo
	client      pb.ReviewDetailServiceClient
	conn        *grpc.ClientConn
	reviewID    int
}

func (s *ReviewDetailApiTestSuite) SetupSuite() {
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
	log, _ := logger.NewLogger("test-review-detail-api", lp)
	obs, _ := observability.NewObservability("test-review-detail-api", log)

	cacheMetrics, _ := observability.NewCacheMetrics("test-review-detail-api")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)

	revDetCacheSrv := reviewdetail_cache.NewReviewDetailMencache(cacheStore)
	revDetCacheApi := api_reviewdetail_cache.NewReviewDetailMencache(cacheStore)

	reviewDetailService := service.NewReviewDetailService(service.ReviewDetailServiceDeps{
		ReviewDetailRepository: repos.ReviewDetail,
		Logger:                 log,
		Observability:          obs,
		Cache:                  revDetCacheSrv,
	})

	// Start gRPC Server
	reviewDetailGapi := gapi.NewReviewDetailHandleGrpc(reviewDetailService)
	server := grpc.NewServer()
	pb.RegisterReviewDetailServiceServer(server, reviewDetailGapi)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewReviewDetailServiceClient(conn)

	// Echo Setup
	s.echo = echo.New()
	mapping := response_api.NewReviewDetailResponseMapper()
	mappingReview := response_api.NewReviewResponseMapper()
	imgUpload := upload_image.NewImageUpload(log)
	apiHandler := errors.NewApiHandler(obs, log)

	api.NewHandlerReviewDetail(s.echo, s.client, log, mapping, mappingReview, imgUpload, apiHandler, revDetCacheApi)

	// Create dependencies
	ctx := context.Background()
	user, _ := repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "ApiDetailUser",
		Email:     "apidetail@example.com",
		Password:  "password123",
	})
	merchant, _ := repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: int(user.UserID),
		Name:   "ApiDetail Merchant",
	})
	catSlug := "api-detail-cat"
	category, _ := repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "ApiDetail Category",
		SlugCategory: &catSlug,
	})
	prodSlug := "api-detail-prod"
	product, _ := repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:  int(merchant.MerchantID),
		CategoryID:  int(category.CategoryID),
		Name:        "ApiDetail Product",
		SlugProduct: &prodSlug,
		Price:       100,
	})
	review, _ := repos.Review.CreateReview(ctx, &requests.CreateReviewRequest{
		UserID:    int(user.UserID),
		ProductID: int(product.ProductID),
		Rating:    5,
		Comment:   "ApiDetail Review",
	})
	s.reviewID = int(review.ReviewID)
}

func (s *ReviewDetailApiTestSuite) TearDownSuite() {
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

func (s *ReviewDetailApiTestSuite) TestReviewDetailApiLifecycle() {
	// 1. Create (Multipart Form)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("review_id", fmt.Sprintf("%d", s.reviewID))
	_ = writer.WriteField("type", "photo")
	_ = writer.WriteField("caption", "Via Echo Multipart")
	
	part, _ := writer.CreateFormFile("url", "echo-test.jpg")
	_, _ = part.Write([]byte("echo test image content"))
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/review-detail/create", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var createRes map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &createRes)
	
	data := createRes["data"].(map[string]interface{})
	detailID := int(data["id"].(float64))

	// 2. Find By ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/review-detail/%d", detailID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Update (Multipart)
	body = new(bytes.Buffer)
	writer = multipart.NewWriter(body)
	_ = writer.WriteField("review_id", fmt.Sprintf("%d", s.reviewID)) // review_id is required by parseReviewDetailForm
	_ = writer.WriteField("type", "photo")
	_ = writer.WriteField("caption", "Updated Via Echo Multipart")
	
	part, _ = writer.CreateFormFile("url", "echo-updated.jpg")
	_, _ = part.Write([]byte("echo updated content"))
	_ = writer.Close()

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/review-detail/update/%d", detailID), body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	
	// Get new URL path (old one should be deleted by service)
	var updateRes map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &updateRes)
	dataUpd := updateRes["data"].(map[string]interface{})
	newUrlPath := dataUpd["url"].(string)

	// 4. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/review-detail/trashed/%d", detailID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/review-detail/restore/%d", detailID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Delete Permanent
	// Trash again
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/review-detail/trashed/%d", detailID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/review-detail/permanent/%d", detailID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// Cleanup: newUrlPath should be gone if it was actually created on disk.
	_, err2 := os.Stat(newUrlPath)
	s.True(os.IsNotExist(err2), "New image should be deleted permanently")
}

func TestReviewDetailApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewDetailApiTestSuite))
}
