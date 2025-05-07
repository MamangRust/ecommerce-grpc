package repository

import (
	"context"
	"database/sql"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/category_errors"
	"time"
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
func (r *categoryRepository) FindAllCategory(req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCategoriesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategories(r.ctx, reqDb)

	if err != nil {
		return nil, nil, category_errors.ErrFindAllCategory
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCategoriesRecordPagination(res), &totalCount, nil
}

func (r *categoryRepository) FindByActive(req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCategoriesActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategoriesActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, category_errors.ErrFindByActiveCategory
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCategoriesRecordActivePagination(res), &totalCount, nil
}

func (r *categoryRepository) FindByTrashed(req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCategoriesTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategoriesTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, category_errors.ErrFindByTrashedCategory
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCategoriesRecordTrashedPagination(res), &totalCount, nil
}

func (r *categoryRepository) GetMonthlyTotalPrice(req *requests.MonthTotalPrice) ([]*record.CategoriesMonthlyTotalPriceRecord, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalPrice(r.ctx, db.GetMonthlyTotalPriceParams{
		Extract:     currentMonthStart,
		CreatedAt:   sql.NullTime{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: sql.NullTime{Time: prevMonthStart, Valid: true},
		CreatedAt_3: sql.NullTime{Time: prevMonthEnd, Valid: true},
	})

	if err != nil {
		return nil, category_errors.ErrGetMonthlyTotalPrice
	}

	so := r.mapping.ToCategoryMonthlyTotalPrices(res)

	return so, nil
}

func (r *categoryRepository) GetYearlyTotalPrices(year int) ([]*record.CategoriesYearlyTotalPriceRecord, error) {
	res, err := r.db.GetYearlyTotalPrice(r.ctx, int32(year))

	if err != nil {
		return nil, category_errors.ErrGetYearlyTotalPrices
	}

	so := r.mapping.ToCategoryYearlyTotalPrices(res)

	return so, nil
}

func (r *categoryRepository) GetMonthlyTotalPriceById(req *requests.MonthTotalPriceCategory) ([]*record.CategoriesMonthlyTotalPriceRecord, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalPriceById(r.ctx, db.GetMonthlyTotalPriceByIdParams{
		Extract:     currentMonthStart,
		CreatedAt:   sql.NullTime{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: sql.NullTime{Time: prevMonthStart, Valid: true},
		CreatedAt_3: sql.NullTime{Time: prevMonthEnd, Valid: true},
		CategoryID:  int32(req.CategoryID),
	})

	if err != nil {
		return nil, category_errors.ErrGetMonthlyTotalPriceById
	}

	so := r.mapping.ToCategoryMonthlyTotalPricesById(res)

	return so, nil
}

func (r *categoryRepository) GetYearlyTotalPricesById(req *requests.YearTotalPriceCategory) ([]*record.CategoriesYearlyTotalPriceRecord, error) {
	res, err := r.db.GetYearlyTotalPriceById(r.ctx, db.GetYearlyTotalPriceByIdParams{
		Column1:    int32(req.Year),
		CategoryID: int32(req.CategoryID),
	})

	if err != nil {
		return nil, category_errors.ErrGetYearlyTotalPricesById
	}

	so := r.mapping.ToCategoryYearlyTotalPricesById(res)

	return so, nil
}

func (r *categoryRepository) GetMonthlyTotalPriceByMerchant(req *requests.MonthTotalPriceMerchant) ([]*record.CategoriesMonthlyTotalPriceRecord, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalPriceByMerchant(r.ctx, db.GetMonthlyTotalPriceByMerchantParams{
		Extract:     currentMonthStart,
		CreatedAt:   sql.NullTime{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: sql.NullTime{Time: prevMonthStart, Valid: true},
		CreatedAt_3: sql.NullTime{Time: prevMonthEnd, Valid: true},
		MerchantID:  int32(req.MerchantID),
	})

	if err != nil {
		return nil, category_errors.ErrGetMonthlyTotalPriceByMerchant
	}

	so := r.mapping.ToCategoryMonthlyTotalPricesByMerchant(res)

	return so, nil
}

func (r *categoryRepository) GetYearlyTotalPricesByMerchant(req *requests.YearTotalPriceMerchant) ([]*record.CategoriesYearlyTotalPriceRecord, error) {
	res, err := r.db.GetYearlyTotalPriceByMerchant(r.ctx, db.GetYearlyTotalPriceByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, category_errors.ErrGetYearlyTotalPricesByMerchant
	}

	so := r.mapping.ToCategoryYearlyTotalPricesByMerchant(res)

	return so, nil
}

func (r *categoryRepository) GetMonthPrice(year int) ([]*record.CategoriesMonthPriceRecord, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyCategory(r.ctx, yearStart)

	if err != nil {
		return nil, category_errors.ErrGetMonthPrice
	}

	return r.mapping.ToCategoryMonthlyPrices(res), nil
}

func (r *categoryRepository) GetYearPrice(year int) ([]*record.CategoriesYearPriceRecord, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyCategory(r.ctx, yearStart)

	if err != nil {
		return nil, category_errors.ErrGetYearPrice
	}

	return r.mapping.ToCategoryYearlyPrices(res), nil
}

func (r *categoryRepository) GetMonthPriceByMerchant(req *requests.MonthPriceMerchant) ([]*record.CategoriesMonthPriceRecord, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyCategoryByMerchant(r.ctx, db.GetMonthlyCategoryByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})
	if err != nil {
		return nil, category_errors.ErrGetMonthPriceByMerchant
	}

	return r.mapping.ToCategoryMonthlyPricesByMerchant(res), nil
}

func (r *categoryRepository) GetYearPriceByMerchant(req *requests.YearPriceMerchant) ([]*record.CategoriesYearPriceRecord, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyCategoryByMerchant(r.ctx, db.GetYearlyCategoryByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, category_errors.ErrGetYearPriceByMerchant
	}

	return r.mapping.ToCategoryYearlyPricesByMerchant(res), nil
}

func (r *categoryRepository) GetMonthPriceById(req *requests.MonthPriceId) ([]*record.CategoriesMonthPriceRecord, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyCategoryById(r.ctx, db.GetMonthlyCategoryByIdParams{
		Column1:    yearStart,
		CategoryID: int32(req.CategoryID),
	})
	if err != nil {
		return nil, category_errors.ErrGetMonthPriceById
	}

	return r.mapping.ToCategoryMonthlyPricesById(res), nil
}

func (r *categoryRepository) GetYearPriceById(req *requests.YearPriceId) ([]*record.CategoriesYearPriceRecord, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyCategoryById(r.ctx, db.GetYearlyCategoryByIdParams{
		Column1:    yearStart,
		CategoryID: int32(req.CategoryID),
	})

	if err != nil {
		return nil, category_errors.ErrGetYearPriceById
	}

	return r.mapping.ToCategoryYearlyPricesById(res), nil
}

func (r *categoryRepository) FindById(category_id int) (*record.CategoriesRecord, error) {
	res, err := r.db.GetCategoryByID(r.ctx, int32(category_id))

	if err != nil {
		return nil, category_errors.ErrFindCategoryById
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryRepository) FindByIdTrashed(category_id int) (*record.CategoriesRecord, error) {
	res, err := r.db.GetCategoryByIDTrashed(r.ctx, int32(category_id))

	if err != nil {
		return nil, category_errors.ErrFindCategoryByIdTrashed
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryRepository) CreateCategory(request *requests.CreateCategoryRequest) (*record.CategoriesRecord, error) {
	req := db.CreateCategoryParams{
		Name: request.Name,
		Description: sql.NullString{
			String: request.Description,
			Valid:  true,
		},
		SlugCategory: sql.NullString{
			String: *request.SlugCategory,
			Valid:  true,
		},
		ImageCategory: sql.NullString{
			String: request.ImageCategory,
			Valid:  true,
		},
	}

	category, err := r.db.CreateCategory(r.ctx, req)
	if err != nil {
		return nil, category_errors.ErrCreateCategory
	}

	return r.mapping.ToCategoryRecord(category), nil
}

func (r *categoryRepository) UpdateCategory(request *requests.UpdateCategoryRequest) (*record.CategoriesRecord, error) {
	req := db.UpdateCategoryParams{
		CategoryID: int32(*request.CategoryID),
		Name:       request.Name,
		Description: sql.NullString{
			String: request.Description,
			Valid:  true,
		},
		SlugCategory: sql.NullString{
			String: *request.SlugCategory,
			Valid:  true,
		},
		ImageCategory: sql.NullString{
			String: request.ImageCategory,
			Valid:  true,
		},
	}

	res, err := r.db.UpdateCategory(r.ctx, req)

	if err != nil {
		return nil, category_errors.ErrUpdateCategory
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryRepository) TrashedCategory(category_id int) (*record.CategoriesRecord, error) {
	res, err := r.db.TrashCategory(r.ctx, int32(category_id))

	if err != nil {
		return nil, category_errors.ErrTrashedCategory
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryRepository) RestoreCategory(category_id int) (*record.CategoriesRecord, error) {
	res, err := r.db.RestoreCategory(r.ctx, int32(category_id))

	if err != nil {
		return nil, category_errors.ErrRestoreCategory
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryRepository) DeleteCategoryPermanently(category_id int) (bool, error) {
	err := r.db.DeleteCategoryPermanently(r.ctx, int32(category_id))

	if err != nil {
		return false, category_errors.ErrDeleteCategoryPermanently
	}

	return true, nil
}

func (r *categoryRepository) RestoreAllCategories() (bool, error) {
	err := r.db.RestoreAllCategories(r.ctx)

	if err != nil {
		return false, category_errors.ErrRestoreAllCategories
	}
	return true, nil
}

func (r *categoryRepository) DeleteAllPermanentCategories() (bool, error) {
	err := r.db.DeleteAllPermanentCategories(r.ctx)

	if err != nil {
		return false, category_errors.ErrDeleteAllPermanentCategories
	}
	return true, nil
}
