package service_test

import (
	"context"
	"ecommerce/internal/cache"
	banner_cache "ecommerce/internal/cache/banner"
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

type BannerServiceTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	rdb    *redis.Client
	srv    service.BannerService
}

func (s *BannerServiceTestSuite) SetupSuite() {
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
	repo := repository.NewBannerRepository(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	l, err := logger.NewLogger("test-banner-service", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-banner-service", l)
	s.Require().NoError(err)

	// Cache Setup
	cacheMetrics, err := observability.NewCacheMetrics("test-banner-service")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.rdb, l, cacheMetrics)
	bannerCache := banner_cache.NewBannerMencache(cacheStore)

	s.srv = service.NewBannerService(service.BannerServiceDeps{
		BannerRepository: repo,
		Logger:           l,
		Observability:    obs,
		Cache:            bannerCache,
	})
}

func (s *BannerServiceTestSuite) TearDownSuite() {
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

func (s *BannerServiceTestSuite) TestBannerLifecycle() {
	ctx := context.Background()

	// 1. Create Banner
	createReq := &requests.CreateBannerRequest{
		Name:      "Holiday Sale",
		StartDate: "2026-12-01",
		EndDate:   "2026-12-31",
		StartTime: "00:00",
		EndTime:   "23:59",
		IsActive:  true,
	}

	banner, err := s.srv.CreateBanner(ctx, createReq)
	s.NoError(err)
	s.NotNil(banner)

	bannerID := int(banner.BannerID)

	// 2. Find By ID
	found, err := s.srv.FindById(ctx, bannerID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(createReq.Name, found.Name)

	// 3. Update Banner
	updateReq := &requests.UpdateBannerRequest{
		BannerID:  &bannerID,
		Name:      "Holiday Sale Extended",
		StartDate: "2026-12-01",
		EndDate:   "2026-12-31",
		StartTime: "00:00",
		EndTime:   "23:59",
		IsActive:  true,
	}

	updated, err := s.srv.UpdateBanner(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)

	// 4. Find All
	list, total, err := s.srv.FindAll(ctx, &requests.FindAllBanner{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(list)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash Banner
	trashed, err := s.srv.TrashedBanner(ctx, bannerID)
	s.NoError(err)
	s.NotNil(trashed)

	// 6. Restore Banner
	restored, err := s.srv.RestoreBanner(ctx, bannerID)
	s.NoError(err)
	s.NotNil(restored)

	// 7. Delete Permanent
	_, err = s.srv.TrashedBanner(ctx, bannerID)
	s.NoError(err)

	success, err := s.srv.DeleteBannerPermanent(ctx, bannerID)
	s.NoError(err)
	s.True(success)

	// 8. Verify it's gone
	_, err = s.srv.FindById(ctx, bannerID)
	s.Error(err)
}

func TestBannerServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(BannerServiceTestSuite))
}
