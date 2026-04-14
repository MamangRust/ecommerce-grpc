package service_test

import (
	"context"
	"ecommerce/internal/cache"
	reviewdetail_cache "ecommerce/internal/cache/review_detail"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"ecommerce/tests"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type ReviewDetailServiceTestSuite struct {
	suite.Suite
	ts       *tests.TestSuite
	dbPool   *pgxpool.Pool
	rdb      *redis.Client
	srv      service.ReviewDetailService
	reviewID int
}

func (s *ReviewDetailServiceTestSuite) SetupSuite() {
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
	l, err := logger.NewLogger("test-review-detail-service", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-review-detail-service", l)
	s.Require().NoError(err)

	cacheMetrics, err := observability.NewCacheMetrics("test-review-detail-service")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.rdb, l, cacheMetrics)
	revDetCache := reviewdetail_cache.NewReviewDetailMencache(cacheStore)

	s.srv = service.NewReviewDetailService(service.ReviewDetailServiceDeps{
		ReviewDetailRepository: repos.ReviewDetail,
		Logger:                 l,
		Observability:          obs,
		Cache:                  revDetCache,
	})

	ctx := context.Background()

	// Create dependencies (User, Merchant, Category, Product, Review)
	user, err := repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "SrvDetail",
		LastName:  "User",
		Email:     "srvdetail@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)

	merchant, err := repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      int(user.UserID),
		Name:        "SrvDetail Merchant",
		Description: "Merchant for service detail tests",
	})
	s.Require().NoError(err)

	catSlug := "cat-srvdetail"
	category, err := repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "SrvDetail Category",
		SlugCategory: &catSlug,
	})
	s.Require().NoError(err)

	prodSlug := "prod-srvdetail"
	product, err := repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:   int(merchant.MerchantID),
		CategoryID:   int(category.CategoryID),
		Name:         "SrvDetail Product",
		SlugProduct:  &prodSlug,
		Price:        100,
		CountInStock: 50,
	})
	s.Require().NoError(err)

	review, err := repos.Review.CreateReview(ctx, &requests.CreateReviewRequest{
		UserID:    int(user.UserID),
		ProductID: int(product.ProductID),
		Rating:    5,
		Comment:   "Review for service detail test",
	})
	s.Require().NoError(err)
	s.reviewID = int(review.ReviewID)
}

func (s *ReviewDetailServiceTestSuite) TearDownSuite() {
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

func (s *ReviewDetailServiceTestSuite) TestReviewDetailLifecycle() {
	ctx := context.Background()

	// Create a dummy file for testing DeleteReviewDetailPermanent
	dummyFile, err := os.CreateTemp("", "review-detail-test-*.jpg")
	s.NoError(err)
	dummyPath := dummyFile.Name()
	dummyFile.Close()

	// 1. Create Review Detail
	createReq := &requests.CreateReviewDetailRequest{
		ReviewID: s.reviewID,
		Type:     "photo",
		Url:      dummyPath,
		Caption:  "Service detail image",
	}

	detail, err := s.srv.CreateReviewDetail(ctx, createReq)
	s.NoError(err)
	s.NotNil(detail)
	s.Equal(createReq.Url, detail.Url)

	detailID := int(detail.ReviewDetailID)

	// 2. Find By ID
	found, err := s.srv.FindById(ctx, detailID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(detail.Url, found.Url)

	// 3. Update Review Detail
	updateReq := &requests.UpdateReviewDetailRequest{
		ReviewDetailID: &detailID,
		Type:           "photo",
		Url:            dummyPath,
		Caption:        "Updated service detail image",
	}

	updated, err := s.srv.UpdateReviewDetail(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Require().NotNil(updated.Caption)
	s.Equal(updateReq.Caption, *updated.Caption)

	// 4. Find All
	details, total, err := s.srv.FindAllReviews(ctx, &requests.FindAllReview{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(details)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash Review Detail
	trashed, err := s.srv.TrashedReviewDetail(ctx, detailID)
	s.NoError(err)
	s.NotNil(trashed)
	s.True(trashed.DeletedAt.Valid)

	// 6. Restore Review Detail
	restored, err := s.srv.RestoreReviewDetail(ctx, detailID)
	s.NoError(err)
	s.NotNil(restored)
	s.False(restored.DeletedAt.Valid)

	// 7. Delete Permanent
	// Trash again first
	_, err = s.srv.TrashedReviewDetail(ctx, detailID)
	s.NoError(err)

	success, err := s.srv.DeleteReviewDetailPermanent(ctx, detailID)
	s.NoError(err)
	s.True(success)

	// 8. Verify it's gone
	_, err = s.srv.FindById(ctx, detailID)
	s.Error(err)

	// Verify dummy file is also gone
	_, err = os.Stat(dummyPath)
	s.True(os.IsNotExist(err))
}

func TestReviewDetailServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewDetailServiceTestSuite))
}
