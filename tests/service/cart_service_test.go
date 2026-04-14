package service_test

import (
	"context"
	"ecommerce/internal/cache"
	cart_cache "ecommerce/internal/cache/cart"
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

type CartServiceTestSuite struct {
	suite.Suite
	ts        *tests.TestSuite
	dbPool    *pgxpool.Pool
	rdb       *redis.Client
	srv       service.CartService
	userID    int
	productID int
}

func (s *CartServiceTestSuite) SetupSuite() {
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
	s.rdb = redis.NewClient(opt)

	// Dependencies
	queries := db.New(pool)
	repos := repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	l, err := logger.NewLogger("test-cart-service", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-cart-service", l)
	s.Require().NoError(err)

	// Cache Setup
	cacheMetrics, err := observability.NewCacheMetrics("test-cart-service")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.rdb, l, cacheMetrics)
	cartCache := cart_cache.NewCartMencache(cacheStore)

	s.srv = service.NewCartService(service.CartServiceDeps{
		ProductRepository: repos.Product,
		UserRepository:    repos.User,
		CartRepository:    repos.Cart,
		Logger:            l,
		Cache:             cartCache,
		Observability:     obs,
	})

	ctx := context.Background()

	// Setup User and Product
	user, _ := repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Service",
		LastName:  "User",
		Email:     "service.user@example.com",
		Password:  "password123",
	})
	s.userID = int(user.UserID)

	merchant, _ := repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: s.userID,
		Name:   "Service Merchant",
	})

	slugCat := "service-category"
	category, _ := repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "Service Category",
		SlugCategory: &slugCat,
	})

	slugProd := "service-product"
	product, _ := repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:   int(merchant.MerchantID),
		CategoryID:   int(category.CategoryID),
		Name:         "Service Product",
		Description:  "Service Product Description",
		Price:        2000,
		CountInStock: 20,
		Brand:        "Service Brand",
		Weight:       100,
		Rating:       &[]int{5}[0],
		SlugProduct:  &slugProd,
		ImageProduct: "service-product.jpg",
	})
	s.productID = int(product.ProductID)
}

func (s *CartServiceTestSuite) TearDownSuite() {
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

func (s *CartServiceTestSuite) TestCartLifecycle() {
	ctx := context.Background()

	// 1. Create Cart
	createReq := &requests.CreateCartRequest{
		Quantity:  3,
		ProductID: s.productID,
		UserID:    s.userID,
	}

	cart, err := s.srv.CreateCart(ctx, createReq)
	s.NoError(err)
	s.NotNil(cart)
	s.Equal(int32(createReq.Quantity), cart.Quantity)

	cartID := int(cart.CartID)

	// 2. Find All
	list, total, err := s.srv.FindAll(ctx, &requests.FindAllCarts{
		UserID:   s.userID,
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(list)
	s.GreaterOrEqual(*total, 1)

	// 3. Delete Permanent
	success, err := s.srv.DeletePermanent(ctx, cartID)
	s.NoError(err)
	s.True(success)

	// 4. Delete All Permanently
	cart2, _ := s.srv.CreateCart(ctx, createReq)
	cart3, _ := s.srv.CreateCart(ctx, createReq)

	success, err = s.srv.DeleteAllPermanently(ctx, &requests.DeleteCartRequest{
		CartIds: []int{int(cart2.CartID), int(cart3.CartID)},
	})
	s.NoError(err)
	s.True(success)
}

func TestCartServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CartServiceTestSuite))
}
