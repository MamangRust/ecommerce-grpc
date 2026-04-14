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

type ProductRepositoryTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	dbPool     *pgxpool.Pool
	repos      *repository.Repositories
	merchantID int
	categoryID int
}

func (s *ProductRepositoryTestSuite) SetupSuite() {
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
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)

	// 2. Create Merchant
	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      int(user.UserID),
		Name:        "Test Merchant",
		Description: "A test merchant",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	// 3. Create Category
	slug := "test-category"
	category, err := s.repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "Test Category",
		Description:  "A test category",
		SlugCategory: &slug,
	})
	s.Require().NoError(err)
	s.categoryID = int(category.CategoryID)
}

func (s *ProductRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *ProductRepositoryTestSuite) TestProductLifecycle() {
	ctx := context.Background()

	// 1. Create Product
	slug := "test-product"
	createReq := &requests.CreateProductRequest{
		MerchantID:   s.merchantID,
		CategoryID:   s.categoryID,
		Name:         "Test Product",
		Description:  "A test product description",
		Price:        100,
		CountInStock: 50,
		Brand:        "Test Brand",
		Weight:       1000,
		SlugProduct:  &slug,
		ImageProduct: "test-product.jpg",
	}

	product, err := s.repos.Product.CreateProduct(ctx, createReq)
	s.NoError(err)
	s.NotNil(product)
	s.Equal(createReq.Name, product.Name)
	productID := int(product.ProductID)

	// 2. Find By ID
	found, err := s.repos.Product.FindById(ctx, productID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(product.Name, found.Name)

	// 3. Find All
	products, err := s.repos.Product.FindAllProducts(ctx, &requests.FindAllProduct{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(products)

	// 4. Update Product
	updateReq := &requests.UpdateProductRequest{
		ProductID:    &productID,
		MerchantID:   s.merchantID,
		CategoryID:   s.categoryID,
		Name:         "Updated Product",
		Description:  "Updated description",
		Price:        150,
		CountInStock: 40,
		Brand:        "Updated Brand",
		Weight:       1200,
		SlugProduct:  &slug,
		ImageProduct: "updated-product.jpg",
	}

	updated, err := s.repos.Product.UpdateProduct(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(updateReq.Name, updated.Name)

	// 5. Update Stock
	updatedStock, err := s.repos.Product.UpdateProductCountStock(ctx, productID, 30)
	s.NoError(err)
	s.NotNil(updatedStock)
	s.Equal(int32(30), updatedStock.CountInStock)

	// 6. Trash Product
	trashed, err := s.repos.Product.TrashedProduct(ctx, productID)
	s.NoError(err)
	s.NotNil(trashed)
	s.True(trashed.DeletedAt.Valid)

	// 7. Find By Trashed
	trashedList, err := s.repos.Product.FindByTrashed(ctx, &requests.FindAllProduct{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(trashedList)

	// 8. Restore Product
	restored, err := s.repos.Product.RestoreProduct(ctx, productID)
	s.NoError(err)
	s.NotNil(restored)
	s.False(restored.DeletedAt.Valid)

	// 9. Delete Permanent
	// Trash again first
	_, err = s.repos.Product.TrashedProduct(ctx, productID)
	s.NoError(err)

	success, err := s.repos.Product.DeleteProductPermanent(ctx, productID)
	s.NoError(err)
	s.True(success)

	// 10. Verify it's gone
	_, err = s.repos.Product.FindById(ctx, productID)
	s.Error(err)
}

func TestProductRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductRepositoryTestSuite))
}
