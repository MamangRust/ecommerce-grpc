package service_test

import (
	"context"
	"ecommerce/internal/cache"
	merchantpolicies_cache "ecommerce/internal/cache/merchant_policies"
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

type MerchantPoliciesServiceTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	rdb    *redis.Client
	srv    service.MerchantPoliciesService
	merchantRepo repository.MerchantRepository
	userRepo repository.UserRepository
}

func (s *MerchantPoliciesServiceTestSuite) SetupSuite() {
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
	repo := repository.NewMerchantPoliciesRepository(queries)
	s.merchantRepo = repository.NewMerchantRepository(queries)
	s.userRepo = repository.NewUserRepository(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	l, err := logger.NewLogger("test-merchant-policy-service", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-merchant-policy-service", l)
	s.Require().NoError(err)

	cacheMetrics, err := observability.NewCacheMetrics("test-merchant-policy-service")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.rdb, l, cacheMetrics)
	policyCache := merchantpolicies_cache.NewMerchantPoliciesMencache(cacheStore)

	s.srv = service.NewMerchantPoliciesService(service.MerchantPoliciesServiceDeps{
		MerchantPoliciesRepository: repo,
		Logger:                     l,
		Observability:              obs,
		Cache:                      policyCache,
	})
}

func (s *MerchantPoliciesServiceTestSuite) TearDownSuite() {
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

func (s *MerchantPoliciesServiceTestSuite) TestMerchantPolicyLifecycle() {
	ctx := context.Background()

	// 0. Setup: Create User and Merchant
	userReq := &requests.CreateUserRequest{
		FirstName: "Policy",
		LastName:  "Service",
		Email:     "policy-service@example.com",
		Password:  "password123",
	}
	user, _ := s.userRepo.CreateUser(ctx, userReq)
	
	merchantReq := &requests.CreateMerchantRequest{
		UserID:       int(user.UserID),
		Name:         "Policy Merchant",
		Description:  "Policy description",
		Address:      "Address",
		ContactEmail: "policy@email.com",
		ContactPhone: "08123456789",
		Status:       "active",
	}
	merchant, _ := s.merchantRepo.CreateMerchant(ctx, merchantReq)
	merchantID := int(merchant.MerchantID)

	// 1. Create Policy
	createReq := &requests.CreateMerchantPolicyRequest{
		MerchantID:  merchantID,
		PolicyType:  "Shipping",
		Title:       "Free Shipping",
		Description: "Free shipping for orders over $100",
	}

	policy, err := s.srv.CreateMerchantPolicy(ctx, createReq)
	s.NoError(err)
	s.NotNil(policy)

	policyID := int(policy.MerchantPolicyID)

	// 2. Find By ID
	found, err := s.srv.FindById(ctx, policyID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(createReq.Title, found.Title)

	// 3. Update Policy
	updateReq := &requests.UpdateMerchantPolicyRequest{
		MerchantPolicyID: &policyID,
		PolicyType:       "Shipping",
		Title:            "Free Shipping Updated",
		Description:      "Free shipping for orders over $150",
	}

	updated, err := s.srv.UpdateMerchantPolicy(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)

	// 4. Find All
	list, total, err := s.srv.FindAllMerchantPolicy(ctx, &requests.FindAllMerchant{
		Page:     1,
		PageSize: 10,
		Search:   "Policy",
	})
	s.NoError(err)
	s.NotNil(list)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	trashed, err := s.srv.TrashedMerchantPolicy(ctx, policyID)
	s.NoError(err)
	s.NotNil(trashed)

	// 6. Restore
	restored, err := s.srv.RestoreMerchantPolicy(ctx, policyID)
	s.NoError(err)
	s.NotNil(restored)

	// 7. Delete Permanent
	_, _ = s.srv.TrashedMerchantPolicy(ctx, policyID)
	success, err := s.srv.DeleteMerchantPolicyPermanent(ctx, policyID)
	s.NoError(err)
	s.True(success)
}

func TestMerchantPoliciesServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantPoliciesServiceTestSuite))
}
