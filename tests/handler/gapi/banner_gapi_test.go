package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	banner_cache "ecommerce/internal/cache/banner"
	"ecommerce/internal/handler/gapi"
	"ecommerce/internal/pb"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"ecommerce/tests"
	"net"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BannerGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.BannerServiceClient
	conn        *grpc.ClientConn
}

func (s *BannerGapiTestSuite) SetupSuite() {
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

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test-banner-gapi", lp)
	obs, _ := observability.NewObservability("test-banner-gapi", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-banner-gapi")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	bannerCache := banner_cache.NewBannerMencache(cacheStore)

	bannerService := service.NewBannerService(service.BannerServiceDeps{
		BannerRepository: repos.Banner,
		Logger:           log,
		Observability:      obs,
		Cache:              bannerCache,
	})

	// Start gRPC Server
	bannerHandler := gapi.NewBannerHandleGrpc(bannerService)
	server := grpc.NewServer()
	pb.RegisterBannerServiceServer(server, bannerHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	// Create Client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewBannerServiceClient(conn)
}

func (s *BannerGapiTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
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

func (s *BannerGapiTestSuite) TestBannerLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pb.CreateBannerRequest{
		Name:      "Summer Sale Gapi",
		StartDate: "2026-06-01",
		EndDate:   "2026-08-31",
		StartTime: "00:00",
		EndTime:   "23:59",
		IsActive:  true,
	}

	res, err := s.client.Create(ctx, createReq)
	s.NoError(err)
	s.Equal(createReq.Name, res.Data.Name)
	bannerID := res.Data.BannerId

	// 2. Find By ID
	found, err := s.client.FindById(ctx, &pb.FindByIdBannerRequest{Id: bannerID})
	s.NoError(err)
	s.Equal(bannerID, found.Data.BannerId)

	// 3. Update
	updateReq := &pb.UpdateBannerRequest{
		BannerId:  bannerID,
		Name:      "Summer Sale Gapi Updated",
		StartDate: "2026-06-01",
		EndDate:   "2026-08-31",
		StartTime: "00:00",
		EndTime:   "23:59",
		IsActive:  true,
	}
	updated, err := s.client.Update(ctx, updateReq)
	s.NoError(err)
	s.Equal(updateReq.Name, updated.Data.Name)

	// 4. Trash
	_, err = s.client.TrashedBanner(ctx, &pb.FindByIdBannerRequest{Id: bannerID})
	s.NoError(err)

	// 5. Restore
	_, err = s.client.RestoreBanner(ctx, &pb.FindByIdBannerRequest{Id: bannerID})
	s.NoError(err)

	// 6. Delete Permanent
	// Need to trash again before permanent delete
	_, err = s.client.TrashedBanner(ctx, &pb.FindByIdBannerRequest{Id: bannerID})
	s.NoError(err)

	delRes, err := s.client.DeleteBannerPermanent(ctx, &pb.FindByIdBannerRequest{Id: bannerID})
	s.NoError(err)
	s.Equal("success", delRes.Status)
}

func TestBannerGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(BannerGapiTestSuite))
}
