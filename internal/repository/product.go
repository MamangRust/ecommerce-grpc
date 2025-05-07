package repository

import (
	"context"
	"database/sql"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/product_errors"
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

func (r *productRepository) FindAllProducts(req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProducts(r.ctx, reqDb)

	if err != nil {
		return nil, nil, product_errors.ErrFindAllProducts
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToProductsRecordPagination(res), &totalCount, nil
}

func (r *productRepository) FindByActive(req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, product_errors.ErrFindByActive
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToProductsRecordActivePagination(res), &totalCount, nil
}

func (r *productRepository) FindByTrashed(req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, product_errors.ErrFindByTrashed
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToProductsRecordTrashedPagination(res), &totalCount, nil
}

func (r *productRepository) FindByMerchant(req *requests.FindAllProductByMerchant) ([]*record.ProductRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsByMerchantParams{
		MerchantID: int32(req.MerchantID),
		Column2:    sql.NullString{String: req.Search, Valid: true},
		Column3:    req.CategoryID,
		Column4:    req.MinPrice,
		Column5:    req.MaxPrice,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetProductsByMerchant(r.ctx, reqDb)

	if err != nil {
		return nil, nil, product_errors.ErrFindByMerchant
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToProductsRecordMerchantPagination(res), &totalCount, nil
}

func (r *productRepository) FindByCategory(req *requests.FindAllProductByCategory) ([]*record.ProductRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsByCategoryNameParams{
		Name:    req.CategoryName,
		Column2: req.Search,
		Column3: req.MinPrice,
		Column4: req.MaxPrice,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsByCategoryName(r.ctx, reqDb)

	if err != nil {
		return nil, nil, product_errors.ErrFindByCategory
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToProductsRecordCategoryPagination(res), &totalCount, nil
}

func (r *productRepository) FindById(product_id int) (*record.ProductRecord, error) {
	res, err := r.db.GetProductByID(r.ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrFindById
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productRepository) FindByIdTrashed(product_id int) (*record.ProductRecord, error) {
	res, err := r.db.GetProductByIdTrashed(r.ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrFindByIdTrashed
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productRepository) CreateProduct(request *requests.CreateProductRequest) (*record.ProductRecord, error) {
	req := db.CreateProductParams{
		MerchantID:   int32(request.MerchantID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  sql.NullString{String: request.Description, Valid: true},
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        sql.NullString{String: request.Brand, Valid: true},
		Weight:       sql.NullInt32{Int32: int32(request.Weight), Valid: true},
		SlugProduct: sql.NullString{
			String: *request.SlugProduct,
			Valid:  true,
		},
		ImageProduct: sql.NullString{String: request.ImageProduct, Valid: true},
	}

	product, err := r.db.CreateProduct(r.ctx, req)

	if err != nil {
		return nil, product_errors.ErrCreateProduct
	}

	return r.mapping.ToProductRecord(product), nil
}

func (r *productRepository) UpdateProduct(request *requests.UpdateProductRequest) (*record.ProductRecord, error) {
	req := db.UpdateProductParams{
		ProductID:    int32(*request.ProductID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  sql.NullString{String: request.Description, Valid: true},
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        sql.NullString{String: request.Brand, Valid: true},
		Weight:       sql.NullInt32{Int32: int32(request.Weight), Valid: true},
		ImageProduct: sql.NullString{String: request.ImageProduct, Valid: true},
	}

	res, err := r.db.UpdateProduct(r.ctx, req)

	if err != nil {
		return nil, product_errors.ErrUpdateProduct
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productRepository) UpdateProductCountStock(product_id int, stock int) (*record.ProductRecord, error) {
	res, err := r.db.UpdateProductCountStock(r.ctx, db.UpdateProductCountStockParams{
		ProductID:    int32(product_id),
		CountInStock: int32(stock),
	})

	if err != nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productRepository) TrashedProduct(product_id int) (*record.ProductRecord, error) {
	res, err := r.db.TrashProduct(r.ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrTrashedProduct
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productRepository) RestoreProduct(product_id int) (*record.ProductRecord, error) {
	res, err := r.db.RestoreProduct(r.ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrRestoreProduct
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productRepository) DeleteProductPermanent(product_id int) (bool, error) {
	err := r.db.DeleteProductPermanently(r.ctx, int32(product_id))

	if err != nil {
		return false, product_errors.ErrDeleteProductPermanent
	}

	return true, nil
}

func (r *productRepository) RestoreAllProducts() (bool, error) {
	err := r.db.RestoreAllProducts(r.ctx)

	if err != nil {
		return false, product_errors.ErrRestoreAllProducts
	}

	return true, nil
}

func (r *productRepository) DeleteAllProductPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentProducts(r.ctx)

	if err != nil {
		return false, product_errors.ErrDeleteAllProductPermanent
	}

	return true, nil
}
