package service_test

import (
	"context"
	"ecommerce/internal/cache"
	shippingaddress_cache "ecommerce/internal/cache/shipping_address"
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

type ShippingAddressServiceTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	service     service.ShippingAddressService
	repos       *repository.Repositories
	userID      int
	merchantID  int
	orderID     int
}

func (s *ShippingAddressServiceTestSuite) SetupSuite() {
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
	s.repos = repository.NewRepositories(queries)

	// Logging & Observability
	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	l, err := logger.NewLogger("test-shipping-service", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-shipping-service", l)
	s.Require().NoError(err)

	// Cache Setup
	cacheMetrics, err := observability.NewCacheMetrics("test-shipping-service")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.redisClient, l, cacheMetrics)
	shipCache := shippingaddress_cache.NewShippingAddressMencache(cacheStore)

	// Service
	s.service = service.NewShippingAddressService(service.ShippingAddressServiceDeps{
		ShippingRepository: s.repos.Shipping,
		Logger:             l,
		Observability:      obs,
		Cache:              shipCache,
	})

	// Create Prerequisites
	ctx := context.Background()
	user, err := s.repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Service",
		LastName:  "User",
		Email:     "service.shipping@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      s.userID,
		Name:        "Service Merchant",
		Description: "A test merchant for service tests",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	order, err := s.repos.Order.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		UserID:     s.userID,
		MerchantID: s.merchantID,
		TotalPrice: 2000,
	})
	s.Require().NoError(err)
	s.orderID = int(order.OrderID)
}

func (s *ShippingAddressServiceTestSuite) TearDownSuite() {
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

func (s *ShippingAddressServiceTestSuite) TestShippingAddressLifecycle() {
	ctx := context.Background()

	// 1. Create Shipping Address via Repository (since service uses it for management)
	createReq := &requests.CreateShippingAddressRequest{
		OrderID:        &s.orderID,
		Alamat:         "Jl. Service No. 1",
		Provinsi:       "West Java",
		Kota:           "Bandung",
		Negara:         "Indonesia",
		Courier:        "SiCepat",
		ShippingMethod: "Best",
		ShippingCost:   12000,
	}

	address, err := s.repos.Shipping.CreateShippingAddress(ctx, createReq)
	s.NoError(err)
	s.NotNil(address)
	shippingID := int(address.ShippingAddressID)

	// 2. Find By ID
	found, err := s.service.FindById(ctx, shippingID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(createReq.Alamat, found.Alamat)

	// 3. Find By Order
	foundByOrder, err := s.service.FindByOrder(ctx, s.orderID)
	s.NoError(err)
	s.NotNil(foundByOrder)
	s.Equal(createReq.Alamat, foundByOrder.Alamat)

	// 4. Find All
	all, total, err := s.service.FindAllShippingAddress(ctx, &requests.FindAllShippingAddress{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(all)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash Shipping Address
	trashed, err := s.service.TrashShippingAddress(ctx, shippingID)
	s.NoError(err)
	s.NotNil(trashed)
	s.True(trashed.DeletedAt.Valid)

	// 6. Find By Trashed
	trashedList, totalTrashed, err := s.service.FindByTrashed(ctx, &requests.FindAllShippingAddress{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(trashedList)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 7. Restore Shipping Address
	restored, err := s.service.RestoreShippingAddress(ctx, shippingID)
	s.NoError(err)
	s.NotNil(restored)
	s.False(restored.DeletedAt.Valid)

	// 8. Delete Permanent
	// Trash again
	_, err = s.service.TrashShippingAddress(ctx, shippingID)
	s.NoError(err)

	success, err := s.service.DeleteShippingAddressPermanently(ctx, shippingID)
	s.NoError(err)
	s.True(success)

	// 9. Verify it's gone
	_, err = s.service.FindById(ctx, shippingID)
	s.Error(err)
}

func TestShippingAddressServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ShippingAddressServiceTestSuite))
}
