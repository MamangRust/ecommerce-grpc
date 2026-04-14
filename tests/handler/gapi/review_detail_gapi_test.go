package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	reviewdetail_cache "ecommerce/internal/cache/review_detail"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/handler/gapi"
	"ecommerce/internal/pb"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"ecommerce/tests"
	"net"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReviewDetailGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.ReviewDetailServiceClient
	conn        *grpc.ClientConn
	reviewID    int
}

func (s *ReviewDetailGapiTestSuite) SetupSuite() {
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
	log, _ := logger.NewLogger("test-review-detail-gapi", lp)
	obs, _ := observability.NewObservability("test-review-detail-gapi", log)

	cacheMetrics, _ := observability.NewCacheMetrics("test-review-detail-gapi")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	revDetCache := reviewdetail_cache.NewReviewDetailMencache(cacheStore)

	reviewDetailService := service.NewReviewDetailService(service.ReviewDetailServiceDeps{
		ReviewDetailRepository: repos.ReviewDetail,
		Logger:                 log,
		Observability:          obs,
		Cache:                  revDetCache,
	})

	// Start gRPC Server
	reviewDetailHandler := gapi.NewReviewDetailHandleGrpc(reviewDetailService)
	server := grpc.NewServer()
	pb.RegisterReviewDetailServiceServer(server, reviewDetailHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	// Create Client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewReviewDetailServiceClient(conn)

	// Create deps
	ctx := context.Background()
	user, _ := repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "GapiDetailUser",
		Email:     "gapidetail@example.com",
		Password:  "password123",
	})
	merchant, _ := repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: int(user.UserID),
		Name:   "GapiDetail Merchant",
	})
	catSlug := "gapi-detail-cat"
	category, _ := repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "GapiDetail Category",
		SlugCategory: &catSlug,
	})
	prodSlug := "gapi-detail-prod"
	product, _ := repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:  int(merchant.MerchantID),
		CategoryID:  int(category.CategoryID),
		Name:        "GapiDetail Product",
		SlugProduct: &prodSlug,
		Price:       100,
	})
	review, _ := repos.Review.CreateReview(ctx, &requests.CreateReviewRequest{
		UserID:    int(user.UserID),
		ProductID: int(product.ProductID),
		Rating:    5,
		Comment:   "GapiDetail Review",
	})
	s.reviewID = int(review.ReviewID)
}

func (s *ReviewDetailGapiTestSuite) TearDownSuite() {
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

func (s *ReviewDetailGapiTestSuite) TestReviewDetailLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pb.CreateReviewDetailRequest{
		ReviewId: int32(s.reviewID),
		Type:     "photo",
		Url:      "uploads/products/gapi-test.jpg", // Valid path now that I removed 'url' tag
		Caption:  "Gapi Test Caption",
	}
	res, err := s.client.Create(ctx, createReq)
	s.NoError(err)
	s.Equal(createReq.Caption, res.Data.Caption)
	detailID := res.Data.Id

	// 2. Find By ID
	found, err := s.client.FindById(ctx, &pb.FindByIdReviewDetailRequest{Id: detailID})
	s.NoError(err)
	s.Equal(detailID, found.Data.Id)

	// 3. Update
	// Create dummy file for old image
	oldPath := "uploads/products/gapi-test.jpg"
	_ = os.MkdirAll("uploads/products", 0755)
	_ = os.WriteFile(oldPath, []byte("test"), 0644)

	updateReq := &pb.UpdateReviewDetailRequest{
		ReviewDetailId: detailID,
		Type:           "photo",
		Url:            "uploads/products/gapi-updated.jpg",
		Caption:        "Gapi Updated Caption",
	}
	updated, err := s.client.Update(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(updateReq.Caption, updated.Data.Caption)

	// 4. Trash
	_, err = s.client.TrashedReviewDetail(ctx, &pb.FindByIdReviewDetailRequest{Id: detailID})
	s.Require().NoError(err)

	// 5. Restore
	_, err = s.client.RestoreReviewDetail(ctx, &pb.FindByIdReviewDetailRequest{Id: detailID})
	s.Require().NoError(err)

	// 6. Delete Permanent
	_, err = s.client.TrashedReviewDetail(ctx, &pb.FindByIdReviewDetailRequest{Id: detailID})
	s.Require().NoError(err)

	// Create dummy file for new image before permanent delete
	newPath := "uploads/products/gapi-updated.jpg"
	_ = os.WriteFile(newPath, []byte("updated"), 0644)

	delRes, err := s.client.DeleteReviewDetailPermanent(ctx, &pb.FindByIdReviewDetailRequest{Id: detailID})
	s.Require().NoError(err)
	s.NotNil(delRes)
	s.Equal("success", delRes.Status)

	// Cleanup
	_ = os.Remove(oldPath)
	_ = os.Remove(newPath)
}

func TestReviewDetailGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewDetailGapiTestSuite))
}
