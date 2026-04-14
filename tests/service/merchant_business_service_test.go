package service_test

import (
	"context"
	"ecommerce/internal/cache"
	merchantbusiness_cache "ecommerce/internal/cache/merchant_business"
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

type MerchantBusinessServiceTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	rdb    *redis.Client
	srv    service.MerchantBusinessService
	merchantRepo repository.MerchantRepository
	userRepo repository.UserRepository
}

func (s *MerchantBusinessServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opt, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.rdb = redis.NewClient(opt)

	queries := db.New(pool)
	repo := repository.NewMerchantBusinessRepository(queries)
	s.merchantRepo = repository.NewMerchantRepository(queries)
	s.userRepo = repository.NewUserRepository(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	l, err := logger.NewLogger("test-merchant-business-service", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-merchant-business-service", l)
	s.Require().NoError(err)

	cacheMetrics, err := observability.NewCacheMetrics("test-merchant-business-service")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.rdb, l, cacheMetrics)
	businessCache := merchantbusiness_cache.NewMerchantBusinessMencache(cacheStore)

	s.srv = service.NewMerchantBusinessService(service.MerchantBusinessServiceDeps{
		MerchantBusinessRepository: repo,
		Logger:                     l,
		Observability:              obs,
		Cache:                      businessCache,
	})
}

func (s *MerchantBusinessServiceTestSuite) TearDownSuite() {
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

func (s *MerchantBusinessServiceTestSuite) TestMerchantBusinessLifecycle() {
	ctx := context.Background()

	// 0. Setup: Create User and Merchant
	userReq := &requests.CreateUserRequest{
		FirstName: "Business",
		LastName:  "Service",
		Email:     "business-service@example.com",
		Password:  "password123",
	}
	user, _ := s.userRepo.CreateUser(ctx, userReq)
	
	merchantReq := &requests.CreateMerchantRequest{
		UserID:       int(user.UserID),
		Name:         "Business Merchant",
		Description:  "Business description",
		Address:      "Address",
		ContactEmail: "business-service@email.com",
		ContactPhone: "08123456789",
		Status:       "active",
	}
	merchant, _ := s.merchantRepo.CreateMerchant(ctx, merchantReq)
	merchantID := int(merchant.MerchantID)

	// 1. Create Business
	createReq := &requests.CreateMerchantBusinessInformationRequest{
		MerchantID:        merchantID,
		BusinessType:      "PT",
		TaxID:             "123-456-789",
		EstablishedYear:   2020,
		NumberOfEmployees: 100,
		WebsiteUrl:        "https://business.com",
	}

	business, err := s.srv.CreateMerchantBusiness(ctx, createReq)
	s.NoError(err)
	s.NotNil(business)

	businessID := int(business.MerchantBusinessInfoID)

	// 2. Find By ID
	found, err := s.srv.FindById(ctx, businessID)
	s.NoError(err)
	s.NotNil(found)

	// 3. Update Business
	updateReq := &requests.UpdateMerchantBusinessInformationRequest{
		MerchantBusinessInfoID: &businessID,
		BusinessType:           "CV",
		TaxID:                  "987-654-321",
		EstablishedYear:        2021,
		NumberOfEmployees:      150,
		WebsiteUrl:             "https://business-updated.com",
	}

	updated, err := s.srv.UpdateMerchantBusiness(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)

	// 4. Find All
	list, total, err := s.srv.FindAllMerchants(ctx, &requests.FindAllMerchant{
		Page:     1,
		PageSize: 10,
		Search:   "Business",
	})
	s.NoError(err)
	s.NotNil(list)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	trashed, err := s.srv.TrashedMerchantBusiness(ctx, businessID)
	s.NoError(err)
	s.NotNil(trashed)

	// 6. Restore
	restored, err := s.srv.RestoreMerchantBusiness(ctx, businessID)
	s.NoError(err)
	s.NotNil(restored)

	// 7. Delete Permanent
	_, _ = s.srv.TrashedMerchantBusiness(ctx, businessID)
	success, err := s.srv.DeleteMerchantBusinessPermanent(ctx, businessID)
	s.NoError(err)
	s.True(success)
}

func TestMerchantBusinessServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantBusinessServiceTestSuite))
}
