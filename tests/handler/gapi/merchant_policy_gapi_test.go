package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	merchantpolicies_cache "ecommerce/internal/cache/merchant_policies"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/handler/gapi"
	"ecommerce/internal/pb"
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

type MerchantPolicyGapiTestSuite struct {
	suite.Suite
	ts           *tests.TestSuite
	dbPool       *pgxpool.Pool
	redisClient  *redis.Client
	handler      pb.MerchantPoliciesServiceServer
	merchantRepo repository.MerchantRepository
	userRepo     repository.UserRepository
}

func (s *MerchantPolicyGapiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)

	queries := db.New(pool)
	s.merchantRepo = repository.NewMerchantRepository(queries)
	s.userRepo = repository.NewUserRepository(queries)
	repo := repository.NewMerchantPoliciesRepository(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test-merchant-policy-gapi", lp)
	obs, _ := observability.NewObservability("test-merchant-policy-gapi", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-merchant-policy-gapi")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	policyCache := merchantpolicies_cache.NewMerchantPoliciesMencache(cacheStore)

	policyService := service.NewMerchantPoliciesService(service.MerchantPoliciesServiceDeps{
		MerchantPoliciesRepository: repo,
		Logger:                     log,
		Observability:              obs,
		Cache:                      policyCache,
	})

	s.handler = gapi.NewMerchantPolicyHandleGrpc(policyService)
}

func (s *MerchantPolicyGapiTestSuite) TearDownSuite() {
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *MerchantPolicyGapiTestSuite) TestMerchantPolicyGapiLifecycle() {
	ctx := context.Background()

	// 0. Setup: Create User and Merchant
	userReq := &requests.CreateUserRequest{
		FirstName: "Policy",
		LastName:  "Gapi",
		Email:     "policy-gapi@merchant.com",
		Password:  "password123",
	}
	user, _ := s.userRepo.CreateUser(ctx, userReq)
	
	merchantReq := &requests.CreateMerchantRequest{
		UserID:       int(user.UserID),
		Name:         "Policy Gapi Merchant",
		ContactEmail: "policy-gapi@email.com",
	}
	merchant, _ := s.merchantRepo.CreateMerchant(ctx, merchantReq)
	merchantID := int(merchant.MerchantID)

	// 1. Create
	createReq := &pb.CreateMerchantPoliciesRequest{
		MerchantId:  int32(merchantID),
		PolicyType:  "Shipping",
		Title:       "Free Shipping Gapi",
		Description: "Free shipping desc",
	}
	res, err := s.handler.Create(ctx, createReq)
	s.NoError(err)
	s.Equal("success", res.Status)
	
	policyID := res.Data.Id

	// 2. FindById
	findRes, err := s.handler.FindById(ctx, &pb.FindByIdMerchantPoliciesRequest{Id: policyID})
	s.NoError(err)
	s.Equal(createReq.Title, findRes.Data.Title)

	// 3. Update
	updateReq := &pb.UpdateMerchantPoliciesRequest{
		MerchantPolicyId: policyID,
		PolicyType:       "Shipping",
		Title:            "Free Shipping Gapi Updated",
		Description:      "Free shipping desc updated",
	}
	updateRes, err := s.handler.Update(ctx, updateReq)
	s.NoError(err)
	s.Equal(updateReq.Title, updateRes.Data.Title)

	// 4. FindAll
	allRes, err := s.handler.FindAll(ctx, &pb.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. Trashed
	trashRes, err := s.handler.TrashedMerchantPolicies(ctx, &pb.FindByIdMerchantPoliciesRequest{Id: policyID})
	s.NoError(err)
	s.NotNil(trashRes.Data.DeletedAt)

	// 6. Restore
	_, err = s.handler.RestoreMerchantPolicies(ctx, &pb.FindByIdMerchantPoliciesRequest{Id: policyID})
	s.NoError(err)
}

func TestMerchantPolicyGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantPolicyGapiTestSuite))
}
