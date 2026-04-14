package service_test

import (
	"context"
	"ecommerce/internal/cache"
	order_cache "ecommerce/internal/cache/order"
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

type OrderServiceTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	dbPool     *pgxpool.Pool
	rdb        *redis.Client
	srv        service.OrderService
	repos      *repository.Repositories
	merchantID int
	userID     int
	categoryID int
	productID  int
}

func (s *OrderServiceTestSuite) SetupSuite() {
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

	s.srv = service.NewOrderService(service.OrderServiceDeps{
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
		FirstName: "OrderSvc",
		LastName:  "User",
		Email:     "ordersvc.user@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	// 2. Create Merchant
	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      s.userID,
		Name:        "OrderSvc Merchant",
		Description: "A merchant for testing order services",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	// 3. Create Category
	slugCat := "test-cat-svc"
	category, err := s.repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "Test Cat Svc",
		Description:  "Test Category for Svc",
		SlugCategory: &slugCat,
	})
	s.Require().NoError(err)
	s.categoryID = int(category.CategoryID)

	// 4. Create Product
	slugProd := "test-prod-svc"
	product, err := s.repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:   s.merchantID,
		CategoryID:   s.categoryID,
		Name:         "Test Prod Svc",
		Description:  "Test Product for Svc",
		Price:        100,
		CountInStock: 100,
		Brand:        "Test Brand",
		Weight:       1000,
		SlugProduct:  &slugProd,
	})
	s.Require().NoError(err)
	s.productID = int(product.ProductID)
}

func (s *OrderServiceTestSuite) TearDownSuite() {
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

func (s *OrderServiceTestSuite) TestOrderServiceLifecycle() {
	ctx := context.Background()

	// 1. Create Order
	createReq := &requests.CreateOrderRequest{
		MerchantID: s.merchantID,
		UserID:     s.userID,
		TotalPrice: 200,
		Items: []requests.CreateOrderItemRequest{
			{
				ProductID: s.productID,
				Quantity:  2,
				Price:     100,
			},
		},
		ShippingAddress: requests.CreateShippingAddressRequest{
			Alamat:         "Jl. Test",
			Provinsi:       "Test Prov",
			Kota:           "Test Kota",
			Courier:        "Test Courier",
			ShippingMethod: "Test Method",
			ShippingCost:   10000,
			Negara:         "Test Negara",
		},
	}

	order, err := s.srv.CreateOrder(ctx, createReq)
	s.NoError(err)
	s.NotNil(order)
	orderID := int(order.OrderID)

	// Fetch ShippingID for Update
	shipping, err := s.repos.Shipping.FindByOrder(ctx, orderID)
	s.NoError(err)
	shippingID := int(shipping.ShippingAddressID)

	// 2. Find By ID
	found, err := s.srv.FindById(ctx, orderID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(orderID), found.OrderID)

	// 3. Find All
	list, total, err := s.srv.FindAllOrders(ctx, &requests.FindAllOrder{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(list)
	s.GreaterOrEqual(*total, 1)

	// 4. Update Order
	updateReq := &requests.UpdateOrderRequest{
		OrderID:    &orderID,
		UserID:     s.userID,
		TotalPrice: 300,
		Items: []requests.UpdateOrderItemRequest{
			// Note: update item logic might need valid order_item_id
		},
		ShippingAddress: requests.UpdateShippingAddressRequest{
			ShippingID:     &shippingID,
			Alamat:         "Jl. Updated",
			Courier:        "Test Courier",
			ShippingMethod: "Test Method",
			ShippingCost:   10000,
		},
	}

	updated, err := s.srv.UpdateOrder(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)

	// 5. Trash Order
	trashed, err := s.srv.TrashedOrder(ctx, orderID)
	s.NoError(err)
	s.NotNil(trashed)
	s.True(trashed.DeletedAt.Valid)

	// 6. Restore Order
	restored, err := s.srv.RestoreOrder(ctx, orderID)
	s.NoError(err)
	s.NotNil(restored)
	s.False(restored.DeletedAt.Valid)

	// 7. Delete Permanent
	_, err = s.srv.TrashedOrder(ctx, orderID)
	s.NoError(err)

	success, err := s.srv.DeleteOrderPermanent(ctx, orderID)
	s.NoError(err)
	s.True(success)

	// 8. Verify it's gone
	_, err = s.srv.FindById(ctx, orderID)
	s.Error(err)
}

func TestOrderServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderServiceTestSuite))
}
