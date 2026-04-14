package service_test

import (
	"context"
	"ecommerce/internal/cache"
	review_cache "ecommerce/internal/cache/review"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"ecommerce/tests"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type ReviewServiceTestSuite struct {
	suite.Suite
	ts        *tests.TestSuite
	dbPool    *pgxpool.Pool
	rdb       *redis.Client
	srv       service.ReviewService
	userID    int
	productID int
}

func (s *ReviewServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opt, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.rdb = redis.NewClient(opt)

	queries := db.New(pool)
	repos := repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	l, err := logger.NewLogger("test-review-service", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-review-service", l)
	s.Require().NoError(err)

	cacheMetrics, err := observability.NewCacheMetrics("test-review-service")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.rdb, l, cacheMetrics)
	revCache := review_cache.NewReviewMencache(cacheStore)

	s.srv = service.NewReviewService(service.ReviewServiceDeps{
		ReviewRepository:  repos.Review,
		ProductRepository: repos.Product,
		UserRepository:    repos.User,
		Logger:            l,
		Observability:     obs,
		Cache:             revCache,
	})

	ctx := context.Background()

	// Create User
	user, err := repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Service",
		LastName:  "Reviewer",
		Email:     "service@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	// Create Merchant
	merchant, err := repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      s.userID,
		Name:        "Service Merchant",
		Description: "A merchant for service tests",
	})
	s.Require().NoError(err)

	// Create Category
	catSlug := "cat-service"
	category, err := repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "Service Category",
		Description:  "A category for service tests",
		SlugCategory: &catSlug,
	})
	s.Require().NoError(err)

	// Create Product
	prodSlug := "prod-service"
	product, err := repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:   int(merchant.MerchantID),
		CategoryID:   int(category.CategoryID),
		Name:         "Service Product",
		Description:  "A product for service tests",
		Price:        100,
		CountInStock: 50,
		Brand:        "Service Brand",
		Weight:       1000,
		SlugProduct:  &prodSlug,
		ImageProduct: "service-product.jpg",
	})
	s.Require().NoError(err)
	s.productID = int(product.ProductID)
}

func (s *ReviewServiceTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.rdb != nil {
		s.rdb.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *ReviewServiceTestSuite) TestReviewLifecycle() {
	ctx := context.Background()

	// 1. Create Review
	createReq := &requests.CreateReviewRequest{
		UserID:    s.userID,
		ProductID: s.productID,
		Rating:    5,
		Comment:   "Outstanding service!",
	}

	review, err := s.srv.CreateReview(ctx, createReq)
	s.NoError(err)
	s.NotNil(review)
	s.Equal(int32(createReq.Rating), review.Rating)

	reviewID := int(review.ReviewID)

	// 2. Find By ID
	found, err := s.srv.FindById(ctx, reviewID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(review.Comment, found.Comment)

	// 3. Update Review
	updateReq := &requests.UpdateReviewRequest{
		ReviewID: &reviewID,
		Name:     "Updated Service Reviewer",
		Rating:   4,
		Comment:  "Excellent, but the update test needs it.",
	}

	updated, err := s.srv.UpdateReview(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(updateReq.Rating), updated.Rating)

	// 4. Find All
	reviews, total, err := s.srv.FindAllReview(ctx, &requests.FindAllReview{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(reviews)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash Review
	trashed, err := s.srv.TrashReview(ctx, reviewID)
	s.NoError(err)
	s.NotNil(trashed)
	s.True(trashed.DeletedAt.Valid)

	// 6. Restore Review
	restored, err := s.srv.RestoreReview(ctx, reviewID)
	s.NoError(err)
	s.NotNil(restored)
	s.False(restored.DeletedAt.Valid)

	// 7. Delete Permanent
	// Trash again first
	_, err = s.srv.TrashReview(ctx, reviewID)
	s.NoError(err)

	success, err := s.srv.DeleteReviewPermanently(ctx, reviewID)
	s.NoError(err)
	s.True(success)

	// 8. Verify it's gone
	_, err = s.srv.FindById(ctx, reviewID)
	s.Error(err)
}

func TestReviewServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewServiceTestSuite))
}
