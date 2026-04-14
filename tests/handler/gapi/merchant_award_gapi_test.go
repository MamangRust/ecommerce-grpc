package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	merchantawards_cache "ecommerce/internal/cache/merchant_awards"
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

type MerchantAwardGapiTestSuite struct {
	suite.Suite
	ts           *tests.TestSuite
	dbPool       *pgxpool.Pool
	redisClient  *redis.Client
	handler      pb.MerchantAwardServiceServer
	merchantRepo repository.MerchantRepository
	userRepo     repository.UserRepository
}

func (s *MerchantAwardGapiTestSuite) SetupSuite() {
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
	awardRepo := repository.NewMerchantAwardRepository(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test-merchant-award-gapi", lp)
	obs, _ := observability.NewObservability("test-merchant-award-gapi", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-merchant-award-gapi")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	awardCache := merchantawards_cache.NewMerchantAward(cacheStore)

	awardService := service.NewMerchantAwardService(service.MerchantAwardServiceDeps{
		MerchantAwardRepository: awardRepo,
		Logger:                  log,
		Observability:           obs,
		Cache:                   awardCache,
	})

	s.handler = gapi.NewMerchantAwardHandleGrpc(awardService)
}

func (s *MerchantAwardGapiTestSuite) TearDownSuite() {
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

func (s *MerchantAwardGapiTestSuite) TestMerchantAwardGapiLifecycle() {
	ctx := context.Background()

	// 0. Setup: Create User and Merchant
	userReq := &requests.CreateUserRequest{
		FirstName: "Award",
		LastName:  "Gapi",
		Email:     "award-gapi@merchant.com",
		Password:  "password123",
	}
	user, _ := s.userRepo.CreateUser(ctx, userReq)
	
	merchantReq := &requests.CreateMerchantRequest{
		UserID:       int(user.UserID),
		Name:         "Award Gapi Merchant",
		ContactEmail: "award-gapi@email.com",
	}
	merchant, _ := s.merchantRepo.CreateMerchant(ctx, merchantReq)
	merchantID := int(merchant.MerchantID)

	// 1. Create
	createReq := &pb.CreateMerchantAwardRequest{
		MerchantId:     int32(merchantID),
		Title:          "Best Gapi Service",
		Description:    "Best service desc",
		IssuedBy:       "Association",
		IssueDate:      "2024-01-01",
		ExpiryDate:     "2025-01-01",
		CertificateUrl: "https://cert.com/1",
	}
	res, err := s.handler.Create(ctx, createReq)
	s.NoError(err)
	s.Equal("success", res.Status)
	
	awardID := res.Data.Id

	// 2. FindById
	findRes, err := s.handler.FindById(ctx, &pb.FindByIdMerchantAwardRequest{Id: awardID})
	s.NoError(err)
	s.Equal(createReq.Title, findRes.Data.Title)

	// 3. Update
	updateReq := &pb.UpdateMerchantAwardRequest{
		MerchantCertificationId: awardID,
		Title:                   "Best Gapi Service Updated",
		Description:             "Best service desc updated",
		IssuedBy:                "Association",
		IssueDate:               "2024-02-01",
		ExpiryDate:              "2025-02-01",
		CertificateUrl:          "https://cert.com/1-updated",
	}
	updateRes, err := s.handler.Update(ctx, updateReq)
	s.NoError(err)
	s.Equal(updateReq.Title, updateRes.Data.Title)

	// 4. FindAll
	allRes, err := s.handler.FindAll(ctx, &pb.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. Trashed
	trashRes, err := s.handler.TrashedMerchantAward(ctx, &pb.FindByIdMerchantAwardRequest{Id: awardID})
	s.NoError(err)
	s.NotNil(trashRes.Data.DeletedAt)

	// 6. Restore
	_, err = s.handler.RestoreMerchantAward(ctx, &pb.FindByIdMerchantAwardRequest{Id: awardID})
	s.NoError(err)
}

func TestMerchantAwardGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantAwardGapiTestSuite))
}
