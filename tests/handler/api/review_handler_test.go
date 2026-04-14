package api_test

import (
	"bytes"
	"context"
	"ecommerce/internal/cache"
	api_review_cache "ecommerce/internal/cache/api/review"
	review_cache "ecommerce/internal/cache/review"
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

type ReviewApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	echo        *echo.Echo
	client      pb.ReviewServiceClient
	conn        *grpc.ClientConn
	userID      int
	productID   int
}

func (s *ReviewApiTestSuite) SetupSuite() {
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
	log, _ := logger.NewLogger("test-review-api", lp)
	obs, _ := observability.NewObservability("test-review-api", log)

	cacheMetrics, _ := observability.NewCacheMetrics("test-review-api")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)

	revCacheSrv := review_cache.NewReviewMencache(cacheStore)
	revCacheApi := api_review_cache.NewReviewMencache(cacheStore)

	reviewService := service.NewReviewService(service.ReviewServiceDeps{
		ReviewRepository:  repos.Review,
		ProductRepository: repos.Product,
		UserRepository:    repos.User,
		Logger:            log,
		Observability:     obs,
		Cache:             revCacheSrv,
	})

	// Start gRPC Server
	reviewGapi := gapi.NewReviewHandleGrpc(reviewService)
	server := grpc.NewServer()
	pb.RegisterReviewServiceServer(server, reviewGapi)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewReviewServiceClient(conn)

	// Echo Setup
	s.echo = echo.New()
	mapping := response_api.NewReviewResponseMapper()
	apiHandler := errors.NewApiHandler(obs, log)

	api.NewHandlerReview(s.echo, s.client, log, mapping, apiHandler, revCacheApi)

	// Create dependencies
	ctx := context.Background()
	user, _ := repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "ApiUser",
		LastName:  "Reviewer",
		Email:     "api@example.com",
		Password:  "password123",
	})
	s.userID = int(user.UserID)

	merchant, _ := repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: s.userID,
		Name:   "Api Merchant",
	})

	catSlug := "api-cat"
	category, _ := repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "Api Category",
		SlugCategory: &catSlug,
	})

	prodSlug := "api-prod"
	product, _ := repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:  int(merchant.MerchantID),
		CategoryID:  int(category.CategoryID),
		Name:        "Api Product",
		SlugProduct: &prodSlug,
		Price:       100,
	})
	s.productID = int(product.ProductID)
}

func (s *ReviewApiTestSuite) TearDownSuite() {
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

func (s *ReviewApiTestSuite) TestReviewApiLifecycle() {
	// 1. Create Review
	createReq := requests.CreateReviewRequest{
		UserID:    s.userID,
		ProductID: s.productID,
		Rating:    5,
		Comment:   "Excellent thru API!",
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/review/create", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var createRes map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &createRes)
	
	data := createRes["data"].(map[string]interface{})
	reviewID := int(data["id"].(float64))

	// 2. Find All
	req = httptest.NewRequest(http.MethodGet, "/api/review", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Update Review
	updateReq := requests.UpdateReviewRequest{
		ReviewID: &reviewID,
		Name:     "ApiUser",
		Rating:   4,
		Comment:  "Updated via API",
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/review/update/%d", reviewID), bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/review/trashed/%d", reviewID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/review/restore/%d", reviewID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Delete Permanent
	// Trash again
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/review/trashed/%d", reviewID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/review/permanent/%d", reviewID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestReviewApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewApiTestSuite))
}
