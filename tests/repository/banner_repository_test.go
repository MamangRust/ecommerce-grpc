package repository_test

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	"ecommerce/tests"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type BannerRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   repository.BannerRepository
}

func (s *BannerRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = repository.NewBannerRepository(queries)
}

func (s *BannerRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *BannerRepositoryTestSuite) TestBannerLifecycle() {
	ctx := context.Background()

	// 1. Create Banner
	createReq := &requests.CreateBannerRequest{
		Name:      "Spring Sale",
		StartDate: "2026-04-10",
		EndDate:   "2026-05-10",
		StartTime: "00:00",
		EndTime:   "23:59",
		IsActive:  true,
	}

	banner, err := s.repo.CreateBanner(ctx, createReq)
	s.NoError(err)
	s.NotNil(banner)
	s.Equal(createReq.Name, banner.Name)

	bannerID := int(banner.BannerID)

	// 2. Find By ID
	found, err := s.repo.FindById(ctx, bannerID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(banner.Name, found.Name)

	// 3. Update Banner
	updateReq := &requests.UpdateBannerRequest{
		BannerID:  &bannerID,
		Name:      "Spring Sale Updated",
		StartDate: "2026-04-10",
		EndDate:   "2026-05-10",
		StartTime: "00:00",
		EndTime:   "23:59",
		IsActive:  true,
	}

	updated, err := s.repo.UpdateBanner(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(updateReq.Name, updated.Name)

	// 4. Trash Banner
	trashed, err := s.repo.TrashedBanner(ctx, bannerID)
	s.NoError(err)
	s.NotNil(trashed)
	s.True(trashed.DeletedAt.Valid)

	// 5. Find By Trashed
	trashedList, err := s.repo.FindByTrashed(ctx, &requests.FindAllBanner{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.NotEmpty(trashedList)

	// 6. Restore Banner
	restored, err := s.repo.RestoreBanner(ctx, bannerID)
	s.NoError(err)
	s.NotNil(restored)
	s.False(restored.DeletedAt.Valid)

	// 7. Delete Permanent
	// Re-trash first
	_, err = s.repo.TrashedBanner(ctx, bannerID)
	s.NoError(err)

	success, err := s.repo.DeleteBannerPermanent(ctx, bannerID)
	s.NoError(err)
	s.True(success)

	// 8. Verify it's gone
	_, err = s.repo.FindById(ctx, bannerID)
	s.Error(err)
}

func TestBannerRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(BannerRepositoryTestSuite))
}
