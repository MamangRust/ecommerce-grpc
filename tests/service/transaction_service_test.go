package service_test

import (
	"context"
	"ecommerce/internal/cache"
	transaction_cache "ecommerce/internal/cache/transaction"
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

type TransactionServiceTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	repos       *repository.Repositories
	service     service.TransactionService
	merchantID  int
	userID      int
	orderID     int
}

func (s *TransactionServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opt, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opt)

	queries := db.New(pool)
	s.repos = repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	l, err := logger.NewLogger("test-transaction-service", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-transaction-service", l)
	s.Require().NoError(err)

	cacheMetrics, err := observability.NewCacheMetrics("test-transaction-service")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.redisClient, l, cacheMetrics)
	transCache := transaction_cache.NewTransactionMencache(cacheStore)

	s.service = service.NewTransactionService(service.TransactionServiceDeps{
		MerchantRepository:    s.repos.Merchant,
		TransactionRepository: s.repos.Transaction,
		OrderRepository:       s.repos.Order,
		OrderItemRepository:   s.repos.OrderItem,
		ShippingRepository:    s.repos.Shipping,
		Logger:                l,
		Cache:                 transCache,
		Observability:         obs,
	})

	ctx := context.Background()

	// 1. Create User
	user, err := s.repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "TransService",
		LastName:  "User",
		Email:     "trans.service@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	// 2. Create Merchant
	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      s.userID,
		Name:        "TransService Merchant",
		Description: "Merchant for service testing",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	// 3. Create Category
	slugCat := "test-category"
	cat, err := s.repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:          "Test Category",
		Description:   "Test Description",
		SlugCategory:  &slugCat,
		ImageCategory: "category.jpg",
	})
	s.Require().NoError(err)

	// 4. Create Product
	slugProd := "test-product"
	rating := 5
	prod, err := s.repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:   s.merchantID,
		CategoryID:   int(cat.CategoryID),
		Name:         "Test Product",
		Description:  "Test Product Description",
		Price:        1000,
		CountInStock: 100,
		Brand:        "Test Brand",
		Weight:       1,
		Rating:       &rating,
		SlugProduct:  &slugProd,
		ImageProduct: "product.jpg",
	})
	s.Require().NoError(err)

	// 5. Create Order
	order, err := s.repos.Order.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: s.merchantID,
		UserID:     s.userID,
		TotalPrice: 1000,
	})
	s.Require().NoError(err)
	s.orderID = int(order.OrderID)

	// 6. Create Order Item (required for verification)
	_, err = s.repos.OrderItem.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{
		OrderID:   s.orderID,
		ProductID: int(prod.ProductID),
		Quantity:  1,
		Price:     1000,
	})
	s.Require().NoError(err)

	// 7. Create Shipping Address (required for verification)
	_, err = s.repos.Shipping.CreateShippingAddress(ctx, &requests.CreateShippingAddressRequest{
		OrderID:        &s.orderID,
		Alamat:         "Test Alamat Lengkap",
		Provinsi:       "Test Provinsi",
		Kota:           "Test Kota",
		Courier:        "JNE",
		ShippingMethod: "REG",
		ShippingCost:   100,
		Negara:         "Indonesia",
	})
	s.Require().NoError(err)
}

func (s *TransactionServiceTestSuite) TearDownSuite() {
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

func (s *TransactionServiceTestSuite) TestTransactionServiceLifecycle() {
	ctx := context.Background()
	// Total amount = (1000 * 1) + 100 = 1100
	// PPN = 1100 * 11% = 121
	// Total with Tax = 1221
	providedAmount := 2000

	// 1. Create Transaction
	createReq := &requests.CreateTransactionRequest{
		OrderID:       s.orderID,
		MerchantID:    s.merchantID,
		PaymentMethod: "credit_card",
		Amount:        providedAmount,
	}

	trans, err := s.service.CreateTransaction(ctx, createReq)
	s.NoError(err)
	s.NotNil(trans)
	s.Equal("success", trans.PaymentStatus)
	s.Equal(int32(1221), trans.Amount)
	transID := int(trans.TransactionID)

	// 2. Find All
	all, total, err := s.service.FindAllTransactions(ctx, &requests.FindAllTransaction{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.NotEmpty(all)
	s.NotNil(total)
	s.GreaterOrEqual(*total, 1)

	// 3. Find By ID
	found, err := s.service.FindById(ctx, transID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(transID), found.TransactionID)

	// 4. Update Transaction (usually fails if status is success, but let's check logic)
	// Service logic says: if existingTx.PaymentStatus == "success" || existingTx.PaymentStatus == "refunded" -> ERROR
	updateReq := &requests.UpdateTransactionRequest{
		TransactionID: &transID,
		OrderID:       s.orderID,
		MerchantID:    s.merchantID,
		PaymentMethod: "debit_card",
		Amount:        2000,
	}
	_, err = s.service.UpdateTransaction(ctx, updateReq)
	s.Error(err) // Should error because status is already success

	// 5. Trash
	trashed, err := s.service.TrashedTransaction(ctx, transID)
	s.NoError(err)
	s.NotNil(trashed)

	// 6. Restore
	restored, err := s.service.RestoreTransaction(ctx, transID)
	s.NoError(err)
	s.NotNil(restored)

	// 7. Delete Permanent
	// Trash again
	_, err = s.service.TrashedTransaction(ctx, transID)
	s.NoError(err)

	success, err := s.service.DeleteTransactionPermanently(ctx, transID)
	s.NoError(err)
	s.True(success)
}

func TestTransactionServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionServiceTestSuite))
}
