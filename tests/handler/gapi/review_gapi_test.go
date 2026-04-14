package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	review_cache "ecommerce/internal/cache/review"
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
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReviewGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.ReviewServiceClient
	conn        *grpc.ClientConn
	userID      int
	productID   int
}

func (s *ReviewGapiTestSuite) SetupSuite() {
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
	log, _ := logger.NewLogger("test-review-gapi", lp)
	obs, _ := observability.NewObservability("test-review-gapi", log)

	cacheMetrics, _ := observability.NewCacheMetrics("test-review-gapi")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	revCache := review_cache.NewReviewMencache(cacheStore)

	reviewService := service.NewReviewService(service.ReviewServiceDeps{
		ReviewRepository:  repos.Review,
		ProductRepository: repos.Product,
		UserRepository:    repos.User,
		Logger:            log,
		Observability:     obs,
		Cache:             revCache,
	})

	// Start gRPC Server
	reviewHandler := gapi.NewReviewHandleGrpc(reviewService)
	server := grpc.NewServer()
	pb.RegisterReviewServiceServer(server, reviewHandler)
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
	s.client = pb.NewReviewServiceClient(conn)

	// Create deps
	ctx := context.Background()
	user, _ := repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "GapiUser",
		Email:     "gapi@example.com",
		Password:  "password123",
	})
	s.userID = int(user.UserID)

	merchant, _ := repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: s.userID,
		Name:   "Gapi Merchant",
	})

	catSlug := "gapi-cat"
	category, _ := repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "Gapi Category",
		SlugCategory: &catSlug,
	})

	prodSlug := "gapi-prod"
	product, _ := repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:  int(merchant.MerchantID),
		CategoryID:  int(category.CategoryID),
		Name:        "Gapi Product",
		SlugProduct: &prodSlug,
		Price:       100,
	})
	s.productID = int(product.ProductID)
}

func (s *ReviewGapiTestSuite) TearDownSuite() {
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

func (s *ReviewGapiTestSuite) TestReviewLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pb.CreateReviewRequest{
		UserId:    int32(s.userID),
		ProductId: int32(s.productID),
		Rating:    5,
		Comment:   "Good product via Gapi",
	}
	res, err := s.client.Create(ctx, createReq)
	s.NoError(err)
	s.Equal(createReq.Comment, res.Data.Comment)
	reviewID := res.Data.Id

	// 2. Find All
	all, err := s.client.FindAll(ctx, &pb.FindAllReviewRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.GreaterOrEqual(len(all.Data), 1)
	found := false
	for _, r := range all.Data {
		if r.Id == reviewID {
			found = true
			break
		}
	}
	s.True(found, "Created review should be found in FindAll")

	// 3. Update
	updateReq := &pb.UpdateReviewRequest{
		ReviewId: reviewID,
		Name:     "GapiUser",
		Rating:   4,
		Comment:  "Updated via Gapi",
	}
	updated, err := s.client.Update(ctx, updateReq)
	s.NoError(err)
	s.Equal(updateReq.Comment, updated.Data.Comment)

	// 4. Trash
	_, err = s.client.TrashedReview(ctx, &pb.FindByIdReviewRequest{Id: reviewID})
	s.NoError(err)

	// 5. Restore
	_, err = s.client.RestoreReview(ctx, &pb.FindByIdReviewRequest{Id: reviewID})
	s.NoError(err)

	// 6. Delete Permanent
	_, err = s.client.TrashedReview(ctx, &pb.FindByIdReviewRequest{Id: reviewID})
	s.NoError(err)

	delRes, err := s.client.DeleteReviewPermanent(ctx, &pb.FindByIdReviewRequest{Id: reviewID})
	s.NoError(err)
	s.Equal("success", delRes.Status)
}

func TestReviewGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewGapiTestSuite))
}
