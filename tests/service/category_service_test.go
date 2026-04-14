package service_test

import (
	"ecommerce/internal/cache"
	category_cache "ecommerce/internal/cache/category"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"ecommerce/tests"
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type CategoryServiceTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	rdb    *redis.Client
	srv    service.CategoryService
}

func (s *CategoryServiceTestSuite) SetupSuite() {
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
	repo := repository.NewCategoryRepository(queries)


	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	l, err := logger.NewLogger("test-service", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-service", l)
	s.Require().NoError(err)

	// Cache Setup
	cacheMetrics, err := observability.NewCacheMetrics("test-service")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.rdb, l, cacheMetrics)
	catCache := category_cache.NewCategoryMencache(cacheStore)


	s.srv = service.NewCategoryService(service.CategoryServiceDeps{
		CategoryRepository: repo,
		Logger:             l,
		Observability:      obs,
		Cache:              catCache,
	})
}

func (s *CategoryServiceTestSuite) TearDownSuite() {
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

func (s *CategoryServiceTestSuite) TestCategoryLifecycle() {
	ctx := context.Background()

	// 1. Create Category
	slug := "electronics-service"
	createReq := &requests.CreateCategoryRequest{
		Name:          "Electronics",
		Description:   "Electronic gadgets",
		SlugCategory:  &slug,
		ImageCategory: "electronics.jpg",
	}

	category, err := s.srv.CreateCategory(ctx, createReq)
	s.NoError(err)
	s.NotNil(category)

	categoryID := int(category.CategoryID)

	// 2. Find By ID
	found, err := s.srv.FindById(ctx, categoryID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(category.Name, found.Name)

	// 3. Update Category
	slugUpdate := "electronics-v2-service"
	updateReq := &requests.UpdateCategoryRequest{
		CategoryID:    &categoryID,
		Name:          "Electronics v2",
		Description:   "Updated electronic gadgets",
		SlugCategory:  &slugUpdate,
		ImageCategory: "electronics-v2.jpg",
	}

	updated, err := s.srv.UpdateCategory(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)

	// 4. Find All
	list, total, err := s.srv.FindAll(ctx, &requests.FindAllCategory{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(list)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash Category
	trashed, err := s.srv.TrashedCategory(ctx, categoryID)
	s.NoError(err)
	s.NotNil(trashed)
	s.True(trashed.DeletedAt.Valid)

	// 6. Restore Category
	restored, err := s.srv.RestoreCategory(ctx, categoryID)
	s.NoError(err)
	s.NotNil(restored)
	s.False(restored.DeletedAt.Valid)

	// 7. Delete Permanent
	// Re-trash first because DeleteCategoryPermanent requires the category to be trashed
	_, err = s.srv.TrashedCategory(ctx, categoryID)
	s.NoError(err)

	success, err := s.srv.DeleteCategoryPermanent(ctx, categoryID)
	s.NoError(err)
	s.True(success)

	// 8. Verify it's gone
	_, err = s.srv.FindById(ctx, categoryID)
	s.Error(err)
}


func TestCategoryServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryServiceTestSuite))
}
