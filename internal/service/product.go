package service

import (
	"context"
	product_cache "ecommerce/internal/cache/product"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/category_errors"
	merchant_errors "ecommerce/pkg/errors/merchant"
	"ecommerce/pkg/errors/product_errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"ecommerce/pkg/utils"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type productService struct {
	categoryRepository repository.CategoryRepository
	merchantRepository repository.MerchantRepository
	productRepository  repository.ProductRepository
	logger             logger.LoggerInterface
	observability      observability.TraceLoggerObservability
	cache              product_cache.ProductMencache
}

type ProductServiceDeps struct {
	CategoryRepository repository.CategoryRepository
	MerchantRepository repository.MerchantRepository
	ProductRepository  repository.ProductRepository
	Logger             logger.LoggerInterface
	Observability      observability.TraceLoggerObservability
	Cache              product_cache.ProductMencache
}

func NewProductService(deps ProductServiceDeps) ProductService {
	return &productService{
		categoryRepository: deps.CategoryRepository,
		merchantRepository: deps.MerchantRepository,
		productRepository:  deps.ProductRepository,
		logger:             deps.Logger,
		observability:      deps.Observability,
		cache:              deps.Cache,
	}
}

func (s *productService) FindAllProducts(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsRow, *int, error) {
	const method = "FindAllProducts"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedProducts(ctx, req); found {
		logSuccess("Successfully retrieved all product records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	products, err := s.productRepository.FindAllProducts(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetProductsRow](
			s.logger,
			product_errors.ErrFailedFindAllProducts,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(products) > 0 {
		totalCount = int(products[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedProducts(ctx, req, products, &totalCount)

	logSuccess("Successfully fetched all products",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return products, &totalCount, nil
}

func (s *productService) FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*db.GetProductsByMerchantRow, *int, error) {
	const method = "FindByMerchantProducts"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search
	merchantId := req.MerchantID

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
		attribute.Int("merchant_id", merchantId))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedProductsByMerchant(ctx, req); found {
		logSuccess("Successfully retrieved merchant product records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.Int("merchant_id", merchantId))
		return data, total, nil
	}

	products, err := s.productRepository.FindByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetProductsByMerchantRow](
			s.logger,
			product_errors.ErrFailedFindProductsByMerchant,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.Int("merchant_id", merchantId),
		)
	}

	var totalCount int

	if len(products) > 0 {
		totalCount = int(products[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedProductsByMerchant(ctx, req, products, &totalCount)

	logSuccess("Successfully fetched merchant products",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.Int("merchant_id", merchantId))

	return products, &totalCount, nil
}

func (s *productService) FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*db.GetProductsByCategoryNameRow, *int, error) {
	const method = "FindByCategoryProducts"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search
	category_name := req.CategoryName

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
		attribute.String("category_name", category_name))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedProductsByCategory(ctx, req); found {
		logSuccess("Successfully retrieved category product records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("category_name", category_name))
		return data, total, nil
	}

	products, err := s.productRepository.FindByCategory(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetProductsByCategoryNameRow](
			s.logger,
			product_errors.ErrFailedFindProductsByCategory,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("category_name", category_name),
		)
	}

	var totalCount int

	if len(products) > 0 {
		totalCount = int(products[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedProductsByCategory(ctx, req, products, &totalCount)

	logSuccess("Successfully fetched category products",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("category_name", category_name))

	return products, &totalCount, nil
}

func (s *productService) FindById(ctx context.Context, productID int) (*db.GetProductByIDRow, error) {
	const method = "FindProductById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("productID", productID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedProduct(ctx, productID); found {
		logSuccess("Successfully retrieved product by ID from cache",
			zap.Int("productID", productID))
		return data, nil
	}

	product, err := s.productRepository.FindById(ctx, productID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetProductByIDRow](
			s.logger,
			product_errors.ErrFailedFindProductById,
			method,
			span,

			zap.Int("productID", productID),
		)
	}

	s.cache.SetCachedProduct(ctx, product)

	logSuccess("Successfully fetched product by ID",
		zap.Int("productID", productID))

	return product, nil
}

func (s *productService) FindByActive(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsActiveRow, *int, error) {
	const method = "FindByActiveProducts"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedProductActive(ctx, req); found {
		logSuccess("Successfully retrieved active product records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	products, err := s.productRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetProductsActiveRow](
			s.logger,
			product_errors.ErrFailedFindProductsByActive,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(products) > 0 {
		totalCount = int(products[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedProductActive(ctx, req, products, &totalCount)

	logSuccess("Successfully fetched active products",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return products, &totalCount, nil
}

func (s *productService) FindByTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsTrashedRow, *int, error) {
	const method = "FindByTrashedProducts"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedProductTrashed(ctx, req); found {
		logSuccess("Successfully retrieved trashed product records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	products, err := s.productRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetProductsTrashedRow](
			s.logger,
			product_errors.ErrFailedFindProductsByTrashed,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(products) > 0 {
		totalCount = int(products[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedProductTrashed(ctx, req, products, &totalCount)

	logSuccess("Successfully fetched trashed products",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return products, &totalCount, nil
}

func (s *productService) CreateProduct(ctx context.Context, req *requests.CreateProductRequest) (*db.CreateProductRow, error) {
	const method = "CreateProduct"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("categoryID", req.CategoryID),
		attribute.Int("merchantID", req.MerchantID),
		attribute.String("name", req.Name))

	defer func() {
		end(status)
	}()

	_, err := s.categoryRepository.FindById(ctx, req.CategoryID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateProductRow](
			s.logger,
			category_errors.ErrFailedFindCategoryById,
			method,
			span,

			zap.Int("categoryID", req.CategoryID),
		)
	}

	_, err = s.merchantRepository.FindById(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateProductRow](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById,
			method,
			span,

			zap.Int("merchantID", req.MerchantID),
		)
	}

	slug := utils.GenerateSlug(req.Name)
	req.SlugProduct = &slug

	product, err := s.productRepository.CreateProduct(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateProductRow](
			s.logger,
			product_errors.ErrFailedCreateProduct,
			method,
			span,
		)
	}

	s.cache.DeleteCachedProduct(ctx, int(product.ProductID))

	logSuccess("Successfully created product",
		zap.Int("productID", int(product.ProductID)),
		zap.String("slug", slug))

	return product, nil
}

func (s *productService) UpdateProduct(ctx context.Context, req *requests.UpdateProductRequest) (*db.UpdateProductRow, error) {
	const method = "UpdateProduct"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("productID", *req.ProductID),
		attribute.Int("categoryID", req.CategoryID),
		attribute.Int("merchantID", req.MerchantID),
		attribute.String("name", req.Name))

	defer func() {
		end(status)
	}()

	_, err := s.categoryRepository.FindById(ctx, req.CategoryID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateProductRow](
			s.logger,
			category_errors.ErrFailedFindCategoryById,
			method,
			span,

			zap.Int("categoryID", req.CategoryID),
		)
	}

	_, err = s.merchantRepository.FindById(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateProductRow](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById,
			method,
			span,

			zap.Int("merchantID", req.MerchantID),
		)
	}

	slug := utils.GenerateSlug(req.Name)
	req.SlugProduct = &slug

	product, err := s.productRepository.UpdateProduct(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateProductRow](
			s.logger,
			product_errors.ErrFailedUpdateProduct,
			method,
			span,
		)
	}

	s.cache.DeleteCachedProduct(ctx, int(product.ProductID))

	logSuccess("Successfully updated product",
		zap.Int("productID", int(product.ProductID)),
		zap.String("slug", slug))

	return product, nil
}

func (s *productService) UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*db.UpdateProductCountStockRow, error) {
	const method = "UpdateProductCountStock"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("product_id", product_id),
		attribute.Int("stock", stock))

	defer func() {
		end(status)
	}()

	product, err := s.productRepository.UpdateProductCountStock(ctx, product_id, stock)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateProductCountStockRow](
			s.logger,
			product_errors.ErrFailedUpdateProduct,
			method,
			span,

			zap.Int("product_id", product_id),
			zap.Int("stock", stock),
		)
	}

	s.cache.DeleteCachedProduct(ctx, product_id)

	logSuccess("Successfully updated product stock",
		zap.Int("product_id", product_id),
		zap.Int("new_stock", stock))

	return product, nil
}

func (s *productService) TrashedProduct(ctx context.Context, productID int) (*db.Product, error) {
	const method = "TrashedProduct"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("productID", productID))

	defer func() {
		end(status)
	}()

	product, err := s.productRepository.TrashedProduct(ctx, productID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Product](
			s.logger,
			product_errors.ErrFailedTrashProduct,
			method,
			span,

			zap.Int("product_id", productID),
		)
	}

	s.cache.DeleteCachedProduct(ctx, productID)

	logSuccess("Successfully trashed product",
		zap.Int("productID", productID))

	return product, nil
}

func (s *productService) RestoreProduct(ctx context.Context, productID int) (*db.Product, error) {
	const method = "RestoreProduct"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("productID", productID))

	defer func() {
		end(status)
	}()

	product, err := s.productRepository.RestoreProduct(ctx, productID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Product](
			s.logger,
			product_errors.ErrFailedRestoreProduct,
			method,
			span,

			zap.Int("product_id", productID),
		)
	}

	s.cache.DeleteCachedProduct(ctx, productID)

	logSuccess("Successfully restored product",
		zap.Int("productID", productID))

	return product, nil
}

func (s *productService) DeleteProductPermanent(ctx context.Context, productID int) (bool, error) {
	const method = "DeleteProductPermanent"

	ctx, span, end, status, logSuccess :=
		s.observability.StartTracingAndLogging(
			ctx,
			method,
			attribute.Int("productID", productID),
		)

	defer func() {
		end(status)
	}()

	product, err := s.productRepository.FindByIdTrashed(ctx, productID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			product_errors.ErrFailedFindProductTrashedById,
			method,
			span,

			zap.Int("product_id", productID),
		)
	}

	if product.ImageProduct != nil && *product.ImageProduct != "" {
		if err := os.Remove(*product.ImageProduct); err != nil {
			if !os.IsNotExist(err) {
				status = "error"
				return errorhandler.HandleError[bool](
					s.logger,
					product_errors.ErrFailedDeleteImageProduct,
					method,
					span,
					zap.String("image_path", *product.ImageProduct),
				)
			}

		}
	}

	_, err = s.productRepository.DeleteProductPermanent(ctx, productID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			product_errors.ErrFailedDeleteProductPermanent,
			method,
			span,

			zap.Int("product_id", productID),
		)
	}

	s.cache.DeleteCachedProduct(ctx, productID)

	logSuccess(
		"Successfully permanently deleted product",
		zap.Int("productID", productID),
	)

	return true, nil
}

func (s *productService) RestoreAllProducts(ctx context.Context) (bool, error) {
	const method = "RestoreAllProducts"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.productRepository.RestoreAllProducts(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			product_errors.ErrFailedRestoreAllProducts,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all trashed products")

	return success, nil
}

func (s *productService) DeleteAllProductPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllProductPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.productRepository.DeleteAllProductPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			product_errors.ErrFailedDeleteAllProductsPermanent,
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all trashed products")

	return success, nil
}
