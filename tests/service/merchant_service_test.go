package service_test

import (
	"context"
	"ecommerce/internal/cache"
	merchant_cache "ecommerce/internal/cache/merchant"
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

type MerchantServiceTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	rdb    *redis.Client
	srv    service.MerchantService
	userRepo repository.UserRepository
}

func (s *MerchantServiceTestSuite) SetupSuite() {
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
	repo := repository.NewMerchantRepository(queries)
	s.userRepo = repository.NewUserRepository(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	l, err := logger.NewLogger("test-merchant-service", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-merchant-service", l)
	s.Require().NoError(err)

	// Cache Setup
	cacheMetrics, err := observability.NewCacheMetrics("test-merchant-service")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.rdb, l, cacheMetrics)
	merchCache := merchant_cache.NewMerchantMencache(cacheStore)

	s.srv = service.NewMerchantService(service.MerchantServiceDeps{
		MerchantRepository: repo,
		Logger:             l,
		Observability:      obs,
		Cache:              merchCache,
	})
}

func (s *MerchantServiceTestSuite) TearDownSuite() {
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

func (s *MerchantServiceTestSuite) TestMerchantLifecycle() {
	ctx := context.Background()

	// 0. Create User first
	userReq := &requests.CreateUserRequest{
		FirstName: "Service",
		LastName:  "Owner",
		Email:     "service@example.com",
		Password:  "password123",
	}
	user, err := s.userRepo.CreateUser(ctx, userReq)
	s.NoError(err)
	s.NotNil(user)
	userID := int(user.UserID)

	// 1. Create Merchant
	createReq := &requests.CreateMerchantRequest{
		UserID:       userID,
		Name:         "Service Merchant",
		Description:  "Service description",
		Address:      "Jakarta",
		ContactEmail: "service@merchant.com",
		ContactPhone: "08123456789",
		Status:       "active",
	}

	merchant, err := s.srv.CreateMerchant(ctx, createReq)
	s.NoError(err)
	s.NotNil(merchant)

	merchantID := int(merchant.MerchantID)

	// 2. Find By ID
	found, err := s.srv.FindById(ctx, merchantID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(createReq.Name, found.Name)

	// 3. Update Merchant
	updateReq := &requests.UpdateMerchantRequest{
		MerchantID:   &merchantID,
		UserID:       userID,
		Name:         "Service Merchant Updated",
		Description:  "Service description updated",
		Address:      "Bandung",
		ContactEmail: "service-updated@merchant.com",
		ContactPhone: "08987654321",
		Status:       "active",
	}

	updated, err := s.srv.UpdateMerchant(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)

	// 4. Find All
	list, total, err := s.srv.FindAllMerchants(ctx, &requests.FindAllMerchant{
		Page:     1,
		PageSize: 10,
		Search:   "Service",
	})
	s.NoError(err)
	s.NotNil(list)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	trashed, err := s.srv.TrashedMerchant(ctx, merchantID)
	s.NoError(err)
	s.NotNil(trashed)

	// 6. Restore
	restored, err := s.srv.RestoreMerchant(ctx, merchantID)
	s.NoError(err)
	s.NotNil(restored)

	// 7. Delete Permanent
	_, _ = s.srv.TrashedMerchant(ctx, merchantID)
	success, err := s.srv.DeleteMerchantPermanent(ctx, merchantID)
	s.NoError(err)
	s.True(success)

	// 8. Verify
	_, err = s.srv.FindById(ctx, merchantID)
	s.Error(err)
}

func TestMerchantServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantServiceTestSuite))
}
