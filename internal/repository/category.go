package repository

import (
	"context"
	"database/sql"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"fmt"
)

type categoryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.CategoryRecordMapper
}

func NewCategoryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CategoryRecordMapper) *categoryRepository {
	return &categoryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *categoryRepository) FindAllCategory(search string, page, pageSize int) ([]*record.CategoriesRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetCategoriesParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategories(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find categories: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCategoriesRecordPagination(res), totalCount, nil
}

func (r *categoryRepository) FindById(category_id int) (*record.CategoriesRecord, error) {
	fmt.Printf("Searching for cashier with ID: %d\n", category_id)
	res, err := r.db.GetCategoryByID(r.ctx, int32(category_id))

	if err != nil {
		fmt.Printf("Error fetching cashier: %v\n", err)

		return nil, fmt.Errorf("failed to find cashier: %w", err)
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryRepository) FindByActive(search string, page, pageSize int) ([]*record.CategoriesRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetCategoriesActiveParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategoriesActive(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find cashier: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCategoriesRecordActivePagination(res), totalCount, nil
}

func (r *categoryRepository) FindByTrashed(search string, page, pageSize int) ([]*record.CategoriesRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetCategoriesTrashedParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategoriesTrashed(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find cashiers: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCategoriesRecordTrashedPagination(res), totalCount, nil
}

func (r *categoryRepository) CreateCategory(request *requests.CreateCategoryRequest) (*record.CategoriesRecord, error) {
	req := db.CreateCategoryParams{
		Name: request.Name,
		Description: sql.NullString{
			String: request.Description,
		},
		SlugCategory: sql.NullString{
			String: request.SlugCategory,
		},
		ImageCategory: sql.NullString{
			String: request.ImageCategory,
		},
	}

	category, err := r.db.CreateCategory(r.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return r.mapping.ToCategoryRecord(category), nil
}

func (r *categoryRepository) UpdateCategory(request *requests.UpdateCategoryRequest) (*record.CategoriesRecord, error) {
	req := db.UpdateCategoryParams{
		CategoryID: int32(request.CategoryID),
		Name:       request.Name,
		Description: sql.NullString{
			String: request.Description,
		},
		SlugCategory: sql.NullString{
			String: request.SlugCategory,
		},
		ImageCategory: sql.NullString{
			String: request.ImageCategory,
		},
	}

	res, err := r.db.UpdateCategory(r.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryRepository) TrashedCategory(category_id int) (*record.CategoriesRecord, error) {
	res, err := r.db.TrashCategory(r.ctx, int32(category_id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash category: %w", err)
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryRepository) RestoreCategory(category_id int) (*record.CategoriesRecord, error) {
	res, err := r.db.RestoreCategory(r.ctx, int32(category_id))

	if err != nil {
		return nil, fmt.Errorf("failed to restore category: %w", err)
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryRepository) DeleteCategoryPermanently(category_id int) (bool, error) {
	err := r.db.DeleteCategoryPermanently(r.ctx, int32(category_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete category: %w", err)
	}

	return true, nil
}

func (r *categoryRepository) RestoreAllCategories() (bool, error) {
	err := r.db.RestoreAllCategories(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all category: %w", err)
	}
	return true, nil
}

func (r *categoryRepository) DeleteAllPermanentCategories() (bool, error) {
	err := r.db.DeleteAllPermanentCategories(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all category permanently: %w", err)
	}
	return true, nil
}
