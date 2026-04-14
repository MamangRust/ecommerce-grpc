package service_test

import (
	"context"
	"ecommerce/internal/cache"
	slider_cache "ecommerce/internal/cache/slider"
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

type SliderServiceTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	service     service.SliderService
	repos       *repository.Repositories
}

func (s *SliderServiceTestSuite) SetupSuite() {
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
	s.redisClient = redis.NewClient(opt)

	// Repositories
	queries := db.New(pool)
	s.repos = repository.NewRepositories(queries)

	// Logging & Observability
	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	l, err := logger.NewLogger("test-slider-service", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-slider-service", l)
	s.Require().NoError(err)

	// Cache Setup
	cacheMetrics, err := observability.NewCacheMetrics("test-slider-service")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.redisClient, l, cacheMetrics)
	sliderCache := slider_cache.NewSliderMencache(cacheStore)

	// Service
	s.service = service.NewSliderService(service.SliderServiceDeps{
		SliderRepository: s.repos.Slider,
		Logger:           l,
		Observability:      obs,
		Cache:              sliderCache,
	})
}

func (s *SliderServiceTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *SliderServiceTestSuite) TestSliderLifecycle() {
	ctx := context.Background()

	// 1. Create Slider
	createReq := &requests.CreateSliderRequest{
		Nama:     "New Arrivals",
		FilePath: "https://example.com/new-arrivals.jpg",
	}

	slider, err := s.service.CreateSlider(ctx, createReq)
	s.NoError(err)
	s.NotNil(slider)
	s.Equal(createReq.Nama, slider.Name)

	sliderID := int(slider.SliderID)

	// 2. Find By ID
	found, err := s.service.FindById(ctx, sliderID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(createReq.Nama, found.Name)

	// 3. Find All
	all, total, err := s.service.FindAllSlider(ctx, &requests.FindAllSlider{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(all)
	s.GreaterOrEqual(*total, 1)

	// 4. Update Slider
	updateReq := &requests.UpdateSliderRequest{
		ID:       &sliderID,
		Nama:     "Updated Arrivals",
		FilePath: "https://example.com/updated-arrivals.jpg",
	}
	updated, err := s.service.UpdateSlider(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(updateReq.Nama, updated.Name)

	// 5. Trash Slider
	trashed, err := s.service.TrashSlider(ctx, sliderID)
	s.NoError(err)
	s.NotNil(trashed)

	// 6. Find By Trashed
	trashedList, totalTrashed, err := s.service.FindByTrashed(ctx, &requests.FindAllSlider{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(trashedList)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 7. Restore Slider
	restored, err := s.service.RestoreSlider(ctx, sliderID)
	s.NoError(err)
	s.NotNil(restored)

	// 8. Delete Permanent
	// Trash again
	_, err = s.service.TrashSlider(ctx, sliderID)
	s.NoError(err)

	success, err := s.service.DeleteSliderPermanently(ctx, sliderID)
	s.NoError(err)
	s.True(success)

	// 9. Verify it's gone
	_, err = s.service.FindById(ctx, sliderID)
	s.Error(err)
}

func TestSliderServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SliderServiceTestSuite))
}
