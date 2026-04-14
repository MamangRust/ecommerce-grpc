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

type ReviewRepositoryTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	dbPool     *pgxpool.Pool
	repos      *repository.Repositories
	userID     int
	productID  int
}

func (s *ReviewRepositoryTestSuite) SetupSuite() {
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
		FirstName: "Reviewer",
		LastName:  "User",
		Email:     "reviewer@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	// 2. Create Merchant
	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      s.userID,
		Name:        "Review Merchant",
		Description: "A merchant for reviews",
	})
	s.Require().NoError(err)

	// 3. Create Category
	catSlug := "cat-review"
	category, err := s.repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "Review Category",
		Description:  "A category for reviews",
		SlugCategory: &catSlug,
	})
	s.Require().NoError(err)

	// 4. Create Product
	prodSlug := "prod-review"
	product, err := s.repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:   int(merchant.MerchantID),
		CategoryID:   int(category.CategoryID),
		Name:         "Review Product",
		Description:  "A product for reviews",
		Price:        100,
		CountInStock: 50,
		Brand:        "Review Brand",
		Weight:       1000,
		SlugProduct:  &prodSlug,
		ImageProduct: "review-product.jpg",
	})
	s.Require().NoError(err)
	s.productID = int(product.ProductID)
}

func (s *ReviewRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *ReviewRepositoryTestSuite) TestReviewLifecycle() {
	ctx := context.Background()

	// 1. Create Review
	createReq := &requests.CreateReviewRequest{
		UserID:    s.userID,
		ProductID: s.productID,
		Rating:    5,
		Comment:   "Excellent product!",
	}

	review, err := s.repos.Review.CreateReview(ctx, createReq)
	s.NoError(err)
	s.NotNil(review)
	s.Equal(int32(createReq.Rating), review.Rating)
	s.Equal(createReq.Comment, review.Comment)

	reviewID := int(review.ReviewID)

	// 2. Find By ID
	found, err := s.repos.Review.FindById(ctx, reviewID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(review.Comment, found.Comment)

	// 3. Find All
	reviews, err := s.repos.Review.FindAllReview(ctx, &requests.FindAllReview{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(reviews)

	// 4. Update Review
	updateReq := &requests.UpdateReviewRequest{
		ReviewID: &reviewID,
		Name:     "Updated Reviewer",
		Rating:   4,
		Comment:  "Pretty good, but could be better.",
	}

	updated, err := s.repos.Review.UpdateReview(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(updateReq.Rating), updated.Rating)
	s.Equal(updateReq.Comment, updated.Comment)

	// 5. Trash Review
	trashed, err := s.repos.Review.TrashReview(ctx, reviewID)
	s.NoError(err)
	s.NotNil(trashed)
	s.True(trashed.DeletedAt.Valid)

	// 6. Find By Trashed
	trashedList, err := s.repos.Review.FindByTrashed(ctx, &requests.FindAllReview{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(trashedList)

	// 7. Restore Review
	restored, err := s.repos.Review.RestoreReview(ctx, reviewID)
	s.NoError(err)
	s.NotNil(restored)
	s.False(restored.DeletedAt.Valid)

	// 8. Delete Permanent
	// Trash again first
	_, err = s.repos.Review.TrashReview(ctx, reviewID)
	s.NoError(err)

	success, err := s.repos.Review.DeleteReviewPermanently(ctx, reviewID)
	s.NoError(err)
	s.True(success)

	// 9. Verify it's gone
	_, err = s.repos.Review.FindById(ctx, reviewID)
	s.Error(err)
}

func TestReviewRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewRepositoryTestSuite))
}
