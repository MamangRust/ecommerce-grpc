package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	merchantdetail_cache "ecommerce/internal/cache/merchant_detail"
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

type MerchantDetailGapiTestSuite struct {
	suite.Suite
	ts           *tests.TestSuite
	dbPool       *pgxpool.Pool
	redisClient  *redis.Client
	handler      pb.MerchantDetailServiceServer
	merchantRepo repository.MerchantRepository
	userRepo     repository.UserRepository
}

func (s *MerchantDetailGapiTestSuite) SetupSuite() {
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
	detailRepo := repository.NewMerchantDetailRepository(queries)
	socialRepo := repository.NewMerchantSocialLinkRepository(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test-merchant-detail-gapi", lp)
	obs, _ := observability.NewObservability("test-merchant-detail-gapi", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-merchant-detail-gapi")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	detailCache := merchantdetail_cache.NewMerchantDetailMencache(cacheStore)

	detailService := service.NewMerchantDetailService(service.MerchantDetailServiceDeps{
		MerchantDetailRepository:     detailRepo,
		MerchantSocialLinkRepository: socialRepo,
		Logger:                       log,
		Cache:                        detailCache,
		Observability:                obs,
	})

	s.handler = gapi.NewMerchantDetailHandleGrpc(detailService)
}

func (s *MerchantDetailGapiTestSuite) TearDownSuite() {
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

func (s *MerchantDetailGapiTestSuite) TestMerchantDetailGapiLifecycle() {
	ctx := context.Background()

	// 0. Setup: Create User and Merchant
	userReq := &requests.CreateUserRequest{
		FirstName: "Detail",
		LastName:  "Gapi",
		Email:     "detail-gapi@merchant.com",
		Password:  "password123",
	}
	user, _ := s.userRepo.CreateUser(ctx, userReq)
	
	merchantReq := &requests.CreateMerchantRequest{
		UserID:       int(user.UserID),
		Name:         "Detail Gapi Merchant",
		ContactEmail: "detail-gapi@email.com",
	}
	merchant, _ := s.merchantRepo.CreateMerchant(ctx, merchantReq)
	merchantID := int(merchant.MerchantID)

	// 1. Create
	createReq := &pb.CreateMerchantDetailRequest{
		MerchantId:       int32(merchantID),
		DisplayName:      "Gapi Display",
		CoverImageUrl:    "https://example.com/cover.jpg",
		LogoUrl:          "https://example.com/logo.jpg",
		ShortDescription: "Gapi short desc",
		WebsiteUrl:       "https://detail-gapi.com",
		SocialLinks: []*pb.CreateMerchantSocialRequest{
			{
				Platform: "Instagram",
				Url:      "https://instagr.am/gapi",
			},
		},
	}
	res, err := s.handler.Create(ctx, createReq)
	s.NoError(err)
	s.Equal("success", res.Status)
	
	detailID := res.Data.Id

	// 2. FindById
	findRes, err := s.handler.FindById(ctx, &pb.FindByIdMerchantDetailRequest{Id: detailID})
	s.NoError(err)
	s.Equal(createReq.DisplayName, findRes.Data.DisplayName)

	// 3. Update
	updateReq := &pb.UpdateMerchantDetailRequest{
		MerchantDetailId: detailID,
		DisplayName:      "Gapi Display Updated",
		CoverImageUrl:    "https://example.com/cover-updated.jpg",
		LogoUrl:          "https://example.com/logo-updated.jpg",
		ShortDescription: "Gapi short desc updated",
		WebsiteUrl:       "https://detail-gapi-updated.com",
		// SocialLinks: ...
	}
	updateRes, err := s.handler.Update(ctx, updateReq)
	s.NoError(err)
	s.Equal(updateReq.DisplayName, updateRes.Data.DisplayName)

	// 4. FindAll
	allRes, err := s.handler.FindAll(ctx, &pb.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. Trashed
	trashRes, err := s.handler.TrashedMerchantDetail(ctx, &pb.FindByIdMerchantDetailRequest{Id: detailID})
	s.NoError(err)
	s.NotNil(trashRes.Data.DeletedAt)

	// 6. Restore
	_, err = s.handler.RestoreMerchantDetail(ctx, &pb.FindByIdMerchantDetailRequest{Id: detailID})
	s.NoError(err)
}

func TestMerchantDetailGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantDetailGapiTestSuite))
}
