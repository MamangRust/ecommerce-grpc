package repository

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/repository"
	"ecommerce/tests"
	db "ecommerce/pkg/database/schema"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type SliderRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   repository.SliderRepository
}

func (s *SliderRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = repository.NewSliderRepository(queries)
}

func (s *SliderRepositoryTestSuite) TearDownSuite() {
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *SliderRepositoryTestSuite) TestSliderLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &requests.CreateSliderRequest{
		Nama:     "Summer Sale",
		FilePath: "https://example.com/summer-sale.jpg",
	}
	slider, err := s.repo.CreateSlider(ctx, createReq)
	s.NoError(err)
	s.NotNil(slider)
	s.Equal(createReq.Nama, slider.Name)
	s.Equal(createReq.FilePath, slider.Image)

	sliderID := int(slider.SliderID)

	// 2. Find By ID
	found, err := s.repo.FindById(ctx, sliderID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(slider.Name, found.Name)

	// 3. Update
	updateReq := &requests.UpdateSliderRequest{
		ID:       &sliderID,
		Nama:     "Autumn Sale",
		FilePath: "https://example.com/autumn-sale.jpg",
	}
	updated, err := s.repo.UpdateSlider(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(updateReq.Nama, updated.Name)

	// 4. Find All
	allReq := &requests.FindAllSlider{
		Page:     1,
		PageSize: 10,
	}
	allSliders, err := s.repo.FindAllSlider(ctx, allReq)
	s.NoError(err)
	s.NotEmpty(allSliders)

	// 5. Trash
	trashed, err := s.repo.TrashSlider(ctx, sliderID)
	s.NoError(err)
	s.NotNil(trashed)

	// 6. Find By Trashed
	trashedSliders, err := s.repo.FindByTrashed(ctx, allReq)
	s.NoError(err)
	s.NotEmpty(trashedSliders)

	// 7. Restore
	restored, err := s.repo.RestoreSlider(ctx, sliderID)
	s.NoError(err)
	s.NotNil(restored)

	// 8. Delete Permanent
	_, err = s.repo.TrashSlider(ctx, sliderID) // Trash first
	s.NoError(err)
	success, err := s.repo.DeleteSliderPermanently(ctx, sliderID)
	s.NoError(err)
	s.True(success)

	// 9. Verify deletion
	_, err = s.repo.FindById(ctx, sliderID)
	s.Error(err)
}

func TestSliderRepositorySuite(t *testing.T) {
	suite.Run(t, new(SliderRepositoryTestSuite))
}
