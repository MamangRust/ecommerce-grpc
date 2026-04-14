package repository

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/product_errors"
	"ecommerce/pkg/utils"
)

type productRepository struct {
	db *db.Queries
}

func NewProductRepository(db *db.Queries) *productRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) FindAllProducts(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProducts(ctx, reqDb)

	if err != nil {
		return nil, product_errors.ErrFindAllProducts
	}

	return res, nil
}

func (r *productRepository) FindByActive(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsActive(ctx, reqDb)

	if err != nil {
		return nil, product_errors.ErrFindByActive
	}

	return res, nil
}

func (r *productRepository) FindByTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsTrashed(ctx, reqDb)

	if err != nil {
		return nil, product_errors.ErrFindByTrashed
	}

	return res, nil
}

func (r *productRepository) FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*db.GetProductsByMerchantRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsByMerchantParams{
		MerchantID: int32(req.MerchantID),
		Column2:    req.Search,
		Column3:    int32(req.CategoryID),
		Column4:    int32(utils.IntPtrToInt(req.MinPrice)),
		Column5:    int32(utils.IntPtrToInt(req.MaxPrice)),
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetProductsByMerchant(ctx, reqDb)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *productRepository) FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*db.GetProductsByCategoryNameRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsByCategoryNameParams{
		Name:    req.CategoryName,
		Column2: req.Search,
		Column3: int32(utils.IntPtrToInt(req.MinPrice)),
		Column4: int32(utils.IntPtrToInt(req.MaxPrice)),
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsByCategoryName(ctx, reqDb)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *productRepository) FindById(ctx context.Context, product_id int) (*db.GetProductByIDRow, error) {
	res, err := r.db.GetProductByID(ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrFindById
	}

	return res, nil
}

func (r *productRepository) FindByIdTrashed(ctx context.Context, product_id int) (*db.GetProductByIdTrashedRow, error) {
	res, err := r.db.GetProductByIdTrashed(ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrFindByIdTrashed
	}

	return res, nil
}

func (r *productRepository) CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*db.CreateProductRow, error) {
	req := db.CreateProductParams{
		MerchantID:   int32(request.MerchantID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  stringPtr(request.Description),
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        stringPtr(request.Brand),
		Weight:       int32Ptr(request.Weight),
		SlugProduct:  request.SlugProduct,
		ImageProduct: stringPtr(request.ImageProduct),
	}

	product, err := r.db.CreateProduct(ctx, req)

	if err != nil {
		return nil, product_errors.ErrCreateProduct
	}

	return product, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*db.UpdateProductRow, error) {
	req := db.UpdateProductParams{
		ProductID:    int32(*request.ProductID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  stringPtr(request.Description),
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        stringPtr(request.Brand),
		Weight:       int32Ptr(request.Weight),
		SlugProduct:  request.SlugProduct,
		ImageProduct: stringPtr(request.ImageProduct),
	}

	res, err := r.db.UpdateProduct(ctx, req)

	if err != nil {
		return nil, product_errors.ErrUpdateProduct
	}

	return res, nil
}

func (r *productRepository) UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*db.UpdateProductCountStockRow, error) {
	res, err := r.db.UpdateProductCountStock(ctx, db.UpdateProductCountStockParams{
		ProductID:    int32(product_id),
		CountInStock: int32(stock),
	})

	if err != nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}

	return res, nil
}

func (r *productRepository) TrashedProduct(ctx context.Context, product_id int) (*db.Product, error) {
	res, err := r.db.TrashProduct(ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrTrashedProduct
	}

	return res, nil
}

func (r *productRepository) RestoreProduct(ctx context.Context, product_id int) (*db.Product, error) {
	res, err := r.db.RestoreProduct(ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrRestoreProduct
	}

	return res, nil
}

func (r *productRepository) DeleteProductPermanent(ctx context.Context, product_id int) (bool, error) {
	err := r.db.DeleteProductPermanently(ctx, int32(product_id))

	if err != nil {
		return false, product_errors.ErrDeleteProductPermanent
	}

	return true, nil
}

func (r *productRepository) RestoreAllProducts(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllProducts(ctx)

	if err != nil {
		return false, product_errors.ErrRestoreAllProducts
	}

	return true, nil
}

func (r *productRepository) DeleteAllProductPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentProducts(ctx)

	if err != nil {
		return false, product_errors.ErrDeleteAllProductPermanent
	}

	return true, nil
}
