package service_test

import (
	"context"
	"ecommerce/internal/cache"
	product_cache "ecommerce/internal/cache/product"
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

type ProductServiceTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	service     service.ProductService
	merchantID  int
	categoryID  int
	categoryName string
}

func (s *ProductServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	// DB Setup
	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	// Redis Setup
	opt, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opt)

	// Repositories
	queries := db.New(pool)
	repos := repository.NewRepositories(queries)

	// Logging & Observability
	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	l, err := logger.NewLogger("test-product-service", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-product-service", l)
	s.Require().NoError(err)

	// Cache Setup
	cacheMetrics, err := observability.NewCacheMetrics("test-product-service")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.redisClient, l, cacheMetrics)
	prodCache := product_cache.NewProductMencache(cacheStore)

	// Service
	s.service = service.NewProductService(service.ProductServiceDeps{
		CategoryRepository: repos.Category,
		MerchantRepository: repos.Merchant,
		ProductRepository:  repos.Product,
		Logger:             l,
		Observability:      obs,
		Cache:              prodCache,
	})

	// Create Prerequisites
	ctx := context.Background()
	user, err := repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Service",
		LastName:  "User",
		Email:     "service.user@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)

	merchant, err := repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      int(user.UserID),
		Name:        "Service Merchant",
		Description: "A test merchant for service tests",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	slug := "service-category"
	category, err := repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "Service Category",
		Description:  "A test category for service tests",
		SlugCategory: &slug,
	})
	s.Require().NoError(err)
	s.categoryID = int(category.CategoryID)
	s.categoryName = category.Name
}

func (s *ProductServiceTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *ProductServiceTestSuite) TestProductLifecycle() {
	ctx := context.Background()

	// 1. Create Product
	slug := "test-product-service"
	createReq := &requests.CreateProductRequest{
		MerchantID:   s.merchantID,
		CategoryID:   s.categoryID,
		Name:         "Test Product Service",
		Description:  "A test product for service layer",
		Price:        200,
		CountInStock: 100,
		Brand:        "Service Brand",
		Weight:       500,
		SlugProduct:  &slug,
		ImageProduct: "service-product.jpg",
	}

	product, err := s.service.CreateProduct(ctx, createReq)
	s.NoError(err)
	s.NotNil(product)
	s.Equal(createReq.Name, product.Name)
	productID := int(product.ProductID)

	// 2. Find By ID
	found, err := s.service.FindById(ctx, productID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(product.Name, found.Name)

	// 3. Find All
	products, total, err := s.service.FindAllProducts(ctx, &requests.FindAllProduct{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(products)
	s.GreaterOrEqual(*total, 1)

	// 4. Find By Category
	catProducts, catTotal, err := s.service.FindByCategory(ctx, &requests.FindAllProductByCategory{
		CategoryName: s.categoryName,
		Page:         1,
		PageSize:     10,
	})
	s.NoError(err)
	s.NotEmpty(catProducts)
	s.GreaterOrEqual(*catTotal, 1)

	// 5. Update Product
	updateReq := &requests.UpdateProductRequest{
		ProductID:    &productID,
		MerchantID:   s.merchantID,
		CategoryID:   s.categoryID,
		Name:         "Updated Product Service",
		Description:  "Updated description service",
		Price:        250,
		CountInStock: 80,
		Brand:        "Updated Service Brand",
		Weight:       600,
		SlugProduct:  &slug,
		ImageProduct: "updated-service-product.jpg",
	}

	updated, err := s.service.UpdateProduct(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(updateReq.Name, updated.Name)

	// 6. Trash Product
	trashed, err := s.service.TrashedProduct(ctx, productID)
	s.NoError(err)
	s.NotNil(trashed)
	s.True(trashed.DeletedAt.Valid)

	// 7. Find By Trashed
	trashedList, trashedTotal, err := s.service.FindByTrashed(ctx, &requests.FindAllProduct{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(trashedList)
	s.GreaterOrEqual(*trashedTotal, 1)

	// 8. Restore Product
	restored, err := s.service.RestoreProduct(ctx, productID)
	s.NoError(err)
	s.NotNil(restored)
	s.False(restored.DeletedAt.Valid)

	// 9. Delete Permanent
	// Trash again
	_, err = s.service.TrashedProduct(ctx, productID)
	s.NoError(err)

	success, err := s.service.DeleteProductPermanent(ctx, productID)
	s.NoError(err)
	s.True(success)

	// 10. Verify it's gone
	_, err = s.service.FindById(ctx, productID)
	s.Error(err)
}

func TestProductServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductServiceTestSuite))
}
