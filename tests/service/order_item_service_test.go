package service_test

import (
	"context"
	"ecommerce/internal/cache"
	order_cache "ecommerce/internal/cache/order"
	orderitem_cache "ecommerce/internal/cache/order_item"
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

type OrderItemServiceTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	dbPool     *pgxpool.Pool
	rdb        *redis.Client
	srv        service.OrderItemService
	repos      *repository.Repositories
	orderSrv   service.OrderService
	merchantID int
	userID     int
	categoryID int
	productID  int
	orderID    int
}

func (s *OrderItemServiceTestSuite) SetupSuite() {
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
	s.repos = repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	l, err := logger.NewLogger("test-service", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-service", l)
	s.Require().NoError(err)

	// Cache Setup
	cacheMetrics, err := observability.NewCacheMetrics("test-service")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.rdb, l, cacheMetrics)
	orderCache := order_cache.OrderNewMencache(cacheStore)
	orderItemCache := orderitem_cache.NewOrderItemMencache(cacheStore)

	s.srv = service.NewOrderItemService(service.OrderItemServiceDeps{
		OrderItemRepository: s.repos.OrderItem,
		Logger:              l,
		Observability:       obs,
		Cache:               orderItemCache,
	})
	s.orderSrv = service.NewOrderService(service.OrderServiceDeps{
		OrderRepository:     s.repos.Order,
		OrderItemRepository: s.repos.OrderItem,
		ProductRepository:   s.repos.Product,
		UserRepository:      s.repos.User,
		MerchantRepository:  s.repos.Merchant,
		ShippingRepository:  s.repos.Shipping,
		Logger:              l,
		Observability:       obs,
		Cache:               orderCache,
	})

	ctx := context.Background()

	// 1. Create User
	user, err := s.repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "OrderItemSvc",
		LastName:  "User",
		Email:     "orderitemsvc.user@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	// 2. Create Merchant
	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      s.userID,
		Name:        "OrderItemSvc Merchant",
		Description: "A merchant for testing order item services",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	// 3. Create Category
	slugCat := "test-cat-svc-oi"
	category, err := s.repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "Test Cat Svc OI",
		Description:  "Test Category for OI Svc",
		SlugCategory: &slugCat,
	})
	s.Require().NoError(err)
	s.categoryID = int(category.CategoryID)

	// 4. Create Product
	slugProd := "test-prod-svc-oi"
	product, err := s.repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:   s.merchantID,
		CategoryID:   s.categoryID,
		Name:         "Test Prod Svc OI",
		Description:  "Test Product for OI Svc",
		Price:        100,
		CountInStock: 100,
		Brand:        "Test Brand",
		Weight:       1000,
		SlugProduct:  &slugProd,
	})
	s.Require().NoError(err)
	s.productID = int(product.ProductID)

	// 5. Create Order via Service
	createReq := &requests.CreateOrderRequest{
		MerchantID: s.merchantID,
		UserID:     s.userID,
		TotalPrice: 100,
		Items: []requests.CreateOrderItemRequest{
			{
				ProductID: s.productID,
				Quantity:  1,
				Price:     100,
			},
		},
		ShippingAddress: requests.CreateShippingAddressRequest{
			Alamat: "Jl. Items",
		},
	}
	order, err := s.orderSrv.CreateOrder(ctx, createReq)
	s.Require().NoError(err)
	s.orderID = int(order.OrderID)
}

func (s *OrderItemServiceTestSuite) TearDownSuite() {
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

func (s *OrderItemServiceTestSuite) TestOrderItemServiceLifecycle() {
	ctx := context.Background()

	// 1. Find By Order
	items, err := s.srv.FindOrderItemByOrder(ctx, s.orderID)
	s.NoError(err)
	s.NotEmpty(items)

	// 2. Find All
	list, total, err := s.srv.FindAllOrderItems(ctx, &requests.FindAllOrderItems{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(list)
	s.GreaterOrEqual(*total, 1)

	// 3. Find By Active
	activeList, totalActive, err := s.srv.FindByActive(ctx, &requests.FindAllOrderItems{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(activeList)
	s.GreaterOrEqual(*totalActive, 1)

	// Trash the order to see the items as trashed?
	// Actually, trashing item usually is done through repository if needed, but service only has finders.
	// Let's manually trash one via repo for testing FindByTrashed
	for _, item := range list {
		_, err = s.repos.OrderItem.TrashedOrderItem(ctx, int(item.OrderItemID))
		s.NoError(err)
	}

	// 4. Find By Trashed
	trashedList, totalTrashed, err := s.srv.FindByTrashed(ctx, &requests.FindAllOrderItems{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(trashedList)
	s.GreaterOrEqual(*totalTrashed, 1)
}

func TestOrderItemServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemServiceTestSuite))
}
