package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	merchantbusiness_cache "ecommerce/internal/cache/merchant_business"
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

type MerchantBusinessGapiTestSuite struct {
	suite.Suite
	ts           *tests.TestSuite
	dbPool       *pgxpool.Pool
	redisClient  *redis.Client
	handler      pb.MerchantBusinessServiceServer
	merchantRepo repository.MerchantRepository
	userRepo     repository.UserRepository
}

func (s *MerchantBusinessGapiTestSuite) SetupSuite() {
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
	businessRepo := repository.NewMerchantBusinessRepository(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test-merchant-business-gapi", lp)
	obs, _ := observability.NewObservability("test-merchant-business-gapi", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-merchant-business-gapi")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	businessCache := merchantbusiness_cache.NewMerchantBusinessMencache(cacheStore)

	businessService := service.NewMerchantBusinessService(service.MerchantBusinessServiceDeps{
		MerchantBusinessRepository: businessRepo,
		Logger:                     log,
		Observability:              obs,
		Cache:                      businessCache,
	})

	s.handler = gapi.NewMerchantBusinessHandleGrpc(businessService)
}

func (s *MerchantBusinessGapiTestSuite) TearDownSuite() {
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

func (s *MerchantBusinessGapiTestSuite) TestMerchantBusinessGapiLifecycle() {
	ctx := context.Background()

	// 0. Setup: Create User and Merchant
	userReq := &requests.CreateUserRequest{
		FirstName: "Business",
		LastName:  "Gapi",
		Email:     "business-gapi@merchant.com",
		Password:  "password123",
	}
	user, _ := s.userRepo.CreateUser(ctx, userReq)
	
	merchantReq := &requests.CreateMerchantRequest{
		UserID:       int(user.UserID),
		Name:         "Business Gapi Merchant",
		ContactEmail: "business-gapi@email.com",
	}
	merchant, _ := s.merchantRepo.CreateMerchant(ctx, merchantReq)
	merchantID := int(merchant.MerchantID)

	// 1. Create
	createReq := &pb.CreateMerchantBusinessRequest{
		MerchantId:        int32(merchantID),
		BusinessType:      "PT",
		TaxId:             "123-456",
		EstablishedYear:   2020,
		NumberOfEmployees: 50,
		WebsiteUrl:        "https://business-gapi.com",
	}
	res, err := s.handler.Create(ctx, createReq)
	s.NoError(err)
	s.Equal("success", res.Status)
	
	infoID := res.Data.Id

	// 2. FindById
	findRes, err := s.handler.FindById(ctx, &pb.FindByIdMerchantBusinessRequest{Id: infoID})
	s.NoError(err)
	s.Equal(createReq.BusinessType, findRes.Data.BusinessType)

	// 3. Update
	updateReq := &pb.UpdateMerchantBusinessRequest{
		MerchantBusinessInfoId: infoID,
		BusinessType:           "CV",
		TaxId:                  "987-654",
		EstablishedYear:        2021,
		NumberOfEmployees:      60,
		WebsiteUrl:             "https://business-gapi-updated.com",
	}
	updateRes, err := s.handler.Update(ctx, updateReq)
	s.NoError(err)
	s.Equal(updateReq.BusinessType, updateRes.Data.BusinessType)

	// 4. FindAll
	allRes, err := s.handler.FindAll(ctx, &pb.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. Trashed
	trashRes, err := s.handler.TrashedMerchantBusiness(ctx, &pb.FindByIdMerchantBusinessRequest{Id: infoID})
	s.NoError(err)
	s.NotNil(trashRes.Data.DeletedAt)

	// 6. Restore
	_, err = s.handler.RestoreMerchantBusiness(ctx, &pb.FindByIdMerchantBusinessRequest{Id: infoID})
	s.NoError(err)
}

func TestMerchantBusinessGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantBusinessGapiTestSuite))
}
