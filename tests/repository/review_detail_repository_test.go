package repository_test

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	"ecommerce/tests"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type ReviewDetailRepositoryTestSuite struct {
	suite.Suite
	ts       *tests.TestSuite
	dbPool   *pgxpool.Pool
	repos    *repository.Repositories
	reviewID int
}

func (s *ReviewDetailRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repos = repository.NewRepositories(queries)

	ctx := context.Background()

	// 1. Create User
	user, err := s.repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Detail",
		LastName:  "User",
		Email:     "detail@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)

	// 2. Create Merchant
	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      int(user.UserID),
		Name:        "Detail Merchant",
		Description: "A merchant for details",
	})
	s.Require().NoError(err)

	// 3. Create Category
	catSlug := "cat-detail"
	category, err := s.repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "Detail Category",
		Description:  "A category for details",
		SlugCategory: &catSlug,
	})
	s.Require().NoError(err)

	// 4. Create Product
	prodSlug := "prod-detail"
	product, err := s.repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:   int(merchant.MerchantID),
		CategoryID:   int(category.CategoryID),
		Name:         "Detail Product",
		Description:  "A product for details",
		Price:        100,
		CountInStock: 50,
		Brand:        "Detail Brand",
		Weight:       1000,
		SlugProduct:  &prodSlug,
		ImageProduct: "detail-product.jpg",
	})
	s.Require().NoError(err)

	// 5. Create Review
	review, err := s.repos.Review.CreateReview(ctx, &requests.CreateReviewRequest{
		UserID:    int(user.UserID),
		ProductID: int(product.ProductID),
		Rating:    5,
		Comment:   "Base review for detail",
	})
	s.Require().NoError(err)
	s.reviewID = int(review.ReviewID)
}

func (s *ReviewDetailRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *ReviewDetailRepositoryTestSuite) TestReviewDetailLifecycle() {
	ctx := context.Background()

	// 1. Create Review Detail
	createReq := &requests.CreateReviewDetailRequest{
		ReviewID: s.reviewID,
		Type:     "photo",
		Url:      "http://example.com/review.jpg",
		Caption:  "My review image",
	}

	detail, err := s.repos.ReviewDetail.CreateReviewDetail(ctx, createReq)
	s.NoError(err)
	s.NotNil(detail)
	s.Equal(createReq.Url, detail.Url)
	s.Require().NotNil(detail.Caption)
	s.Equal(createReq.Caption, *detail.Caption)

	detailID := int(detail.ReviewDetailID)

	// 2. Find By ID
	found, err := s.repos.ReviewDetail.FindById(ctx, detailID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(detail.Url, found.Url)

	// 3. Find All
	details, err := s.repos.ReviewDetail.FindAllReviews(ctx, &requests.FindAllReview{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(details)

	// 4. Update Review Detail
	updateReq := &requests.UpdateReviewDetailRequest{
		ReviewDetailID: &detailID,
		Type:           "video",
		Url:            "http://example.com/review.mp4",
		Caption:        "My review video",
	}

	updated, err := s.repos.ReviewDetail.UpdateReviewDetail(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(updateReq.Url, updated.Url)
	s.Require().NotNil(updated.Caption)
	s.Equal(updateReq.Caption, *updated.Caption)

	// 5. Trash Review Detail
	trashed, err := s.repos.ReviewDetail.TrashedReviewDetail(ctx, detailID)
	s.NoError(err)
	s.NotNil(trashed)
	s.True(trashed.DeletedAt.Valid)

	// 6. Find By Trashed
	trashedList, err := s.repos.ReviewDetail.FindByTrashed(ctx, &requests.FindAllReview{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(trashedList)

	// 7. Restore Review Detail
	restored, err := s.repos.ReviewDetail.RestoreReviewDetail(ctx, detailID)
	s.NoError(err)
	s.NotNil(restored)
	s.False(restored.DeletedAt.Valid)

	// 8. Delete Permanent
	// Trash again first
	_, err = s.repos.ReviewDetail.TrashedReviewDetail(ctx, detailID)
	s.NoError(err)

	success, err := s.repos.ReviewDetail.DeleteReviewDetailPermanent(ctx, detailID)
	s.NoError(err)
	s.True(success)

	// 9. Verify it's gone
	_, err = s.repos.ReviewDetail.FindById(ctx, detailID)
	s.Error(err)
}

func TestReviewDetailRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewDetailRepositoryTestSuite))
}
