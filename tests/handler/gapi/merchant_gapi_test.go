package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	merchant_cache "ecommerce/internal/cache/merchant"
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

type MerchantGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	handler     pb.MerchantServiceServer
	userRepo    repository.UserRepository
}

func (s *MerchantGapiTestSuite) SetupSuite() {
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
	repos := repository.NewRepositories(queries)
	s.userRepo = repos.User

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test-merchant-gapi", lp)
	obs, _ := observability.NewObservability("test-merchant-gapi", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-merchant-gapi")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	merchCache := merchant_cache.NewMerchantMencache(cacheStore)

	merchantService := service.NewMerchantService(service.MerchantServiceDeps{
		MerchantRepository: repos.Merchant,
		Logger:             log,
		Observability:      obs,
		Cache:              merchCache,
	})

	s.handler = gapi.NewMerchantHandleGrpc(merchantService)
}

func (s *MerchantGapiTestSuite) TearDownSuite() {
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

func (s *MerchantGapiTestSuite) TestMerchantGapiLifecycle() {
	ctx := context.Background()

	// 0. Create User
	userReq := &requests.CreateUserRequest{
		FirstName: "Gapi",
		LastName:  "Merchant",
		Email:     "gapi@merchant.com",
		Password:  "password123",
	}
	user, err := s.userRepo.CreateUser(ctx, userReq)
	s.NoError(err)

	// 1. Create
	createReq := &pb.CreateMerchantRequest{
		UserId:       int32(user.UserID),
		Name:         "Gapi Merchant",
		Description:  "Gapi desc",
		Address:      "Gapi Addr",
		ContactEmail: "gapi@email.com",
		ContactPhone: "08123456789",
		Status:       "active",
	}
	res, err := s.handler.Create(ctx, createReq)
	s.NoError(err)
	s.NotNil(res)
	s.Equal("success", res.Status)
	
	merchantID := res.Data.Id

	// 2. FindById
	findRes, err := s.handler.FindById(ctx, &pb.FindByIdMerchantRequest{Id: merchantID})
	s.NoError(err)
	s.Equal(createReq.Name, findRes.Data.Name)

	// 3. Update
	updateReq := &pb.UpdateMerchantRequest{
		MerchantId:   merchantID,
		UserId:       int32(user.UserID),
		Name:         "Gapi Merchant Updated",
		Description:  "Gapi desc updated",
		Address:      "Gapi Addr Updated",
		ContactEmail: "gapi-updated@email.com",
		ContactPhone: "08987654321",
		Status:       "active",
	}
	updateRes, err := s.handler.Update(ctx, updateReq)
	s.NoError(err)
	s.Equal(updateReq.Name, updateRes.Data.Name)

	// 4. FindAll
	allRes, err := s.handler.FindAll(ctx, &pb.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. Trashed
	trashRes, err := s.handler.TrashedMerchant(ctx, &pb.FindByIdMerchantRequest{Id: merchantID})
	s.NoError(err)
	s.NotNil(trashRes.Data.DeletedAt)

	// 6. Restore
	restoreRes, err := s.handler.RestoreMerchant(ctx, &pb.FindByIdMerchantRequest{Id: merchantID})
	s.NoError(err)
	s.Nil(restoreRes.Data.DeletedAt)

	// 7. DeletePermanent
	_, _ = s.handler.TrashedMerchant(ctx, &pb.FindByIdMerchantRequest{Id: merchantID})
	delRes, err := s.handler.DeleteMerchantPermanent(ctx, &pb.FindByIdMerchantRequest{Id: merchantID})
	s.NoError(err)
	s.Equal("success", delRes.Status)
}

func TestMerchantGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantGapiTestSuite))
}
