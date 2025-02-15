package repository

import (
	"context"
	"database/sql"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"errors"
	"fmt"
)

type productRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.ProductRecordMapping
}

func NewProductRepository(db *db.Queries, ctx context.Context, mapping recordmapper.ProductRecordMapping) *productRepository {
	return &productRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *productRepository) FindAllProducts(search string, page, pageSize int) ([]*record.ProductRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetProductsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProducts(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find products: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToProductsRecordPagination(res), totalCount, nil
}

func (r *productRepository) FindByActive(search string, page, pageSize int) ([]*record.ProductRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetProductsActiveParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsActive(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find products: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToProductsRecordActivePagination(res), totalCount, nil
}

func (r *productRepository) FindByTrashed(search string, page, pageSize int) ([]*record.ProductRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetProductsTrashedParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsTrashed(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find users: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToProductsRecordTrashedPagination(res), totalCount, nil
}

func (r *productRepository) FindByMerchant(merchant_id int, search string, page, pageSize int) ([]*record.ProductRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetProductsByMerchantParams{
		MerchantID: int32(merchant_id),
		Column2:    search,
		Limit:      int32(pageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetProductsByMerchant(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find products: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToProductsRecordMerchantPagination(res), totalCount, nil
}

func (r *productRepository) FindByCategory(category_name string, search string, page, pageSize int) ([]*record.ProductRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetProductsByCategoryNameParams{
		Name:    category_name,
		Column2: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsByCategoryName(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find products: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToProductsRecordCategoryPagination(res), totalCount, nil
}

func (r *productRepository) FindById(user_id int) (*record.ProductRecord, error) {
	res, err := r.db.GetProductByID(r.ctx, int32(user_id))

	if err != nil {
		fmt.Printf("Error fetching products: %v\n", err)

		return nil, fmt.Errorf("failed to find products: %w", err)
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productRepository) CreateProduct(request *requests.CreateProductRequest) (*record.ProductRecord, error) {
	req := db.CreateProductParams{
		MerchantID:   int32(request.MerchantID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  sql.NullString{String: request.Description, Valid: request.Description != ""},
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        sql.NullString{String: request.Brand, Valid: request.Brand != ""},
		Weight:       sql.NullInt32{Int32: int32(request.Weight)},
		Rating:       sql.NullFloat64{Float64: float64(request.Rating), Valid: request.Rating != 0},
		SlugProduct: sql.NullString{
			String: request.SlugProduct,
		},
		ImageProduct: sql.NullString{String: request.ImageProduct, Valid: request.ImageProduct != ""},
	}

	product, err := r.db.CreateProduct(r.ctx, req)
	if err != nil {
		return nil, errors.New("failed to create product")
	}

	return r.mapping.ToProductRecord(product), nil
}

func (r *productRepository) UpdateProduct(request *requests.UpdateProductRequest) (*record.ProductRecord, error) {
	req := db.UpdateProductParams{
		ProductID:    int32(request.ProductID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  sql.NullString{String: request.Description, Valid: request.Description != ""},
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        sql.NullString{String: request.Brand, Valid: request.Brand != ""},
		Weight:       sql.NullInt32{Int32: int32(request.Weight)},
		Rating:       sql.NullFloat64{Float64: float64(request.Rating), Valid: request.Rating != 0},
		SlugProduct: sql.NullString{
			String: request.SlugProduct,
		},
		ImageProduct: sql.NullString{String: request.ImageProduct, Valid: request.ImageProduct != ""},
	}

	res, err := r.db.UpdateProduct(r.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productRepository) UpdateProductCountStock(product_id int, stock int) (*record.ProductRecord, error) {
	res, err := r.db.UpdateProductCountStock(r.ctx, db.UpdateProductCountStockParams{
		ProductID:    int32(product_id),
		CountInStock: int32(stock),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productRepository) TrashedProduct(user_id int) (*record.ProductRecord, error) {
	res, err := r.db.TrashProduct(r.ctx, int32(user_id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash user: %w", err)
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productRepository) RestoreProduct(user_id int) (*record.ProductRecord, error) {
	res, err := r.db.RestoreProduct(r.ctx, int32(user_id))

	if err != nil {
		return nil, fmt.Errorf("failed to restore topup: %w", err)
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productRepository) DeleteProductPermanent(user_id int) (bool, error) {
	err := r.db.DeleteProductPermanently(r.ctx, int32(user_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete product: %w", err)
	}

	return true, nil
}

func (r *productRepository) RestoreAllProducts() (bool, error) {
	err := r.db.RestoreAllProducts(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all products: %w", err)
	}
	return true, nil
}

func (r *productRepository) DeleteAllProductPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentProducts(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all products permanently: %w", err)
	}
	return true, nil
}
