package repository_test

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	"ecommerce/tests"
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type CategoryRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   repository.CategoryRepository
}

func (s *CategoryRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = repository.NewCategoryRepository(queries)
}

func (s *CategoryRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *CategoryRepositoryTestSuite) TestCategoryLifecycle() {
	ctx := context.Background()

	// 1. Create Category
	slug := "electronics"
	createReq := &requests.CreateCategoryRequest{
		Name:          "Electronics",
		Description:   "Electronic gadgets",
		SlugCategory:  &slug,
		ImageCategory: "electronics.jpg",
	}

	category, err := s.repo.CreateCategory(ctx, createReq)
	s.NoError(err)
	s.NotNil(category)
	s.Equal(createReq.Name, category.Name)

	categoryID := int(category.CategoryID)

	// 2. Find By ID
	found, err := s.repo.FindById(ctx, categoryID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(category.Name, found.Name)

	// 3. Update Category
	slugUpdate := "electronics-v2"
	updateReq := &requests.UpdateCategoryRequest{
		CategoryID:    &categoryID,
		Name:          "Electronics v2",
		Description:   "Updated electronic gadgets",
		SlugCategory:  &slugUpdate,
		ImageCategory: "electronics-v2.jpg",
	}



	updated, err := s.repo.UpdateCategory(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(updateReq.Name, updated.Name)

	// 4. Trash Category
	trashed, err := s.repo.TrashedCategory(ctx, categoryID)
	s.NoError(err)
	s.NotNil(trashed)
	s.True(trashed.DeletedAt.Valid)

	// 5. Find By Trashed
	trashedList, err := s.repo.FindByTrashed(ctx, &requests.FindAllCategory{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(trashedList)

	// 6. Restore Category
	restored, err := s.repo.RestoreCategory(ctx, categoryID)
	s.NoError(err)
	s.NotNil(restored)
	s.False(restored.DeletedAt.Valid)


	// 7. Delete Permanent
	// Re-trash first because DeleteCategoryPermanently requires the category to be trashed
	_, err = s.repo.TrashedCategory(ctx, categoryID)
	s.NoError(err)

	success, err := s.repo.DeleteCategoryPermanently(ctx, categoryID)
	s.NoError(err)
	s.True(success)

	// 8. Verify it's gone
	_, err = s.repo.FindById(ctx, categoryID)
	s.Error(err)
}


func TestCategoryRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryRepositoryTestSuite))
}
