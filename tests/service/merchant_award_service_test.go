package service_test

import (
	"context"
	"ecommerce/internal/cache"
	merchantawards_cache "ecommerce/internal/cache/merchant_awards"
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

type MerchantAwardServiceTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	rdb    *redis.Client
	srv    service.MerchantAwardService
	merchantRepo repository.MerchantRepository
	userRepo repository.UserRepository
}

func (s *MerchantAwardServiceTestSuite) SetupSuite() {
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
	repo := repository.NewMerchantAwardRepository(queries)
	s.merchantRepo = repository.NewMerchantRepository(queries)
	s.userRepo = repository.NewUserRepository(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	l, err := logger.NewLogger("test-merchant-award-service", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-merchant-award-service", l)
	s.Require().NoError(err)

	cacheMetrics, err := observability.NewCacheMetrics("test-merchant-award-service")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.rdb, l, cacheMetrics)
	awardCache := merchantawards_cache.NewMerchantAward(cacheStore)

	s.srv = service.NewMerchantAwardService(service.MerchantAwardServiceDeps{
		MerchantAwardRepository: repo,
		Logger:                  l,
		Observability:           obs,
		Cache:                   awardCache,
	})
}

func (s *MerchantAwardServiceTestSuite) TearDownSuite() {
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

func (s *MerchantAwardServiceTestSuite) TestMerchantAwardLifecycle() {
	ctx := context.Background()

	// 0. Setup: Create User and Merchant
	userReq := &requests.CreateUserRequest{
		FirstName: "Award",
		LastName:  "Service",
		Email:     "award-service@example.com",
		Password:  "password123",
	}
	user, _ := s.userRepo.CreateUser(ctx, userReq)
	
	merchantReq := &requests.CreateMerchantRequest{
		UserID:       int(user.UserID),
		Name:         "Award Merchant",
		Description:  "Award description",
		Address:      "Address",
		ContactEmail: "award-service@email.com",
		ContactPhone: "08123456789",
		Status:       "active",
	}
	merchant, _ := s.merchantRepo.CreateMerchant(ctx, merchantReq)
	merchantID := int(merchant.MerchantID)

	// 1. Create Award
	createReq := &requests.CreateMerchantCertificationOrAwardRequest{
		MerchantID:     merchantID,
		Title:          "Outstanding Service",
		Description:    "Excellence in service",
		IssuedBy:       "Association",
		IssueDate:      "2024-01-01",
		ExpiryDate:     "2025-01-01",
		CertificateUrl: "https://award.com/cert.pdf",
	}

	award, err := s.srv.CreateMerchantAward(ctx, createReq)
	s.NoError(err)
	s.NotNil(award)

	awardID := int(award.MerchantCertificationID)

	// 2. Find By ID
	found, err := s.srv.FindById(ctx, awardID)
	s.NoError(err)
	s.NotNil(found)

	// 3. Update Award
	updateReq := &requests.UpdateMerchantCertificationOrAwardRequest{
		MerchantCertificationID: &awardID,
		Title:                   "Outstanding Service Updated",
		Description:             "Excellence in service updated",
		IssuedBy:                "Association",
		IssueDate:               "2024-02-01",
		ExpiryDate:              "2025-02-01",
		CertificateUrl:          "https://award.com/cert-updated.pdf",
	}

	updated, err := s.srv.UpdateMerchantAward(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)

	// 4. Find All
	list, total, err := s.srv.FindAllMerchants(ctx, &requests.FindAllMerchant{
		Page:     1,
		PageSize: 10,
		Search:   "Award",
	})
	s.NoError(err)
	s.NotNil(list)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	trashed, err := s.srv.TrashedMerchantAward(ctx, awardID)
	s.NoError(err)
	s.NotNil(trashed)

	// 6. Restore
	restored, err := s.srv.RestoreMerchantAward(ctx, awardID)
	s.NoError(err)
	s.NotNil(restored)

	// 7. Delete Permanent
	_, _ = s.srv.TrashedMerchantAward(ctx, awardID)
	success, err := s.srv.DeleteMerchantPermanent(ctx, awardID)
	s.NoError(err)
	s.True(success)
}

func TestMerchantAwardServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantAwardServiceTestSuite))
}
