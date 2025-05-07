package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/errors/category_errors"
	merchant_errors "ecommerce/pkg/errors/merchant"
	"ecommerce/pkg/errors/product_errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/utils"
	"os"

	"go.uber.org/zap"
)

type productService struct {
	categoryRepository repository.CategoryRepository
	merchantRepository repository.MerchantRepository
	productRepository  repository.ProductRepository
	logger             logger.LoggerInterface
	mapping            response_service.ProductResponseMapper
}

func NewProductService(
	categoryRepository repository.CategoryRepository,
	merchantRepository repository.MerchantRepository,
	productRepository repository.ProductRepository,
	logger logger.LoggerInterface,
	mapping response_service.ProductResponseMapper,
) *productService {
	return &productService{
		categoryRepository: categoryRepository,
		merchantRepository: merchantRepository,
		productRepository:  productRepository,
		logger:             logger,
		mapping:            mapping,
	}
}

func (s *productService) FindAll(req *requests.FindAllProduct) ([]*response.ProductResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all products",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	products, totalRecords, err := s.productRepository.FindAllProducts(req)

	if err != nil {
		s.logger.Error("Failed to retrieve product list",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, product_errors.ErrFailedFindAllProducts
	}

	s.logger.Debug("Successfully fetched products",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToProductsResponse(products), totalRecords, nil
}

func (s *productService) FindByMerchant(req *requests.FindAllProductByMerchant) ([]*response.ProductResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search
	merchantId := req.MerchantID

	s.logger.Debug("Fetching all products by merchant",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search),
		zap.Int("merchant_id", merchantId),
	)

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	products, totalRecords, err := s.productRepository.FindByMerchant(req)

	if err != nil {
		s.logger.Error("Failed to retrieve order list",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, product_errors.ErrFailedFindProductsByMerchant
	}

	s.logger.Debug("Successfully fetched products",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToProductsResponse(products), totalRecords, nil
}

func (s *productService) FindByCategory(req *requests.FindAllProductByCategory) ([]*response.ProductResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search
	category_name := req.CategoryName

	s.logger.Debug("Fetching all products by category name",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search),
		zap.String("category_name", category_name),
	)

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	products, totalRecords, err := s.productRepository.FindByCategory(req)

	if err != nil {
		s.logger.Error("Failed to retrieve category product list",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, product_errors.ErrFailedFindProductsByCategory
	}

	s.logger.Debug("Successfully fetched products",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToProductsResponse(products), totalRecords, nil
}

func (s *productService) FindById(productID int) (*response.ProductResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching product by ID", zap.Int("productID", productID))

	product, err := s.productRepository.FindById(productID)

	if err != nil {
		s.logger.Error("Failed to retrieve product details",
			zap.Error(err),
			zap.Int("productID", productID))

		return nil, product_errors.ErrFailedFindProductById
	}

	return s.mapping.ToProductResponse(product), nil
}

func (s *productService) FindByActive(req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all products active",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	products, totalRecords, err := s.productRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active products",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, product_errors.ErrFailedFindProductsByActive
	}

	s.logger.Debug("Successfully fetched products",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToProductsResponseDeleteAt(products), totalRecords, nil
}

func (s *productService) FindByTrashed(req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all products trashed",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	products, totalRecords, err := s.productRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed products",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, product_errors.ErrFailedFindProductsByTrashed
	}

	s.logger.Debug("Successfully fetched products",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToProductsResponseDeleteAt(products), totalRecords, nil
}

func (s *productService) CreateProduct(req *requests.CreateProductRequest) (*response.ProductResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new product")

	_, err := s.categoryRepository.FindById(req.CategoryID)

	if err != nil {
		s.logger.Error("Category not found for product creation",
			zap.Int("categoryID", req.CategoryID),
			zap.Error(err))

		return nil, category_errors.ErrFailedFindCategoryById
	}

	_, err = s.merchantRepository.FindById(req.MerchantID)

	if err != nil {
		s.logger.Error("Merchant not found for product creation",
			zap.Int("merchantID", req.MerchantID),
			zap.Error(err))

		return nil, merchant_errors.ErrFailedFindMerchantById
	}

	slug := utils.GenerateSlug(req.Name)

	req.SlugProduct = &slug

	product, err := s.productRepository.CreateProduct(req)

	if err != nil {
		s.logger.Error("Failed to create product record",
			zap.Error(err))

		return nil, product_errors.ErrFailedCreateProduct
	}

	s.logger.Debug("Product created successfully", zap.Int("productID", product.ID))

	return s.mapping.ToProductResponse(product), nil
}

func (s *productService) UpdateProduct(req *requests.UpdateProductRequest) (*response.ProductResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating product", zap.Int("productID", *req.ProductID))

	_, err := s.categoryRepository.FindById(req.CategoryID)

	if err != nil {
		s.logger.Error("Category not found for product update",
			zap.Int("categoryID", req.CategoryID),
			zap.Error(err))

		return nil, category_errors.ErrFailedFindCategoryById
	}

	_, err = s.merchantRepository.FindById(req.MerchantID)

	if err != nil {
		s.logger.Error("Merchant not found for product update",
			zap.Int("merchantID", req.MerchantID),
			zap.Error(err))

		return nil, merchant_errors.ErrFailedFindMerchantById
	}

	slug := utils.GenerateSlug(req.Name)

	req.SlugProduct = &slug

	product, err := s.productRepository.UpdateProduct(req)

	if err != nil {
		s.logger.Error("Failed to update product record",
			zap.Error(err))

		return nil, product_errors.ErrFailedUpdateProduct
	}

	s.logger.Debug("Product updated successfully", zap.Int("productID", product.ID))
	return s.mapping.ToProductResponse(product), nil
}

func (s *productService) TrashProduct(productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing product", zap.Int("productID", productID))

	product, err := s.productRepository.TrashedProduct(productID)

	if err != nil {
		s.logger.Error("Failed to move product to trash",
			zap.Int("product_id", productID),
			zap.Error(err))

		return nil, product_errors.ErrFailedTrashProduct
	}

	return s.mapping.ToProductResponseDeleteAt(product), nil
}

func (s *productService) RestoreProduct(productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring product", zap.Int("productID", productID))

	product, err := s.productRepository.RestoreProduct(productID)

	if err != nil {
		s.logger.Error("Failed to restore product from trash",
			zap.Int("product_id", productID),
			zap.Error(err))

		return nil, product_errors.ErrFailedRestoreProduct
	}

	return s.mapping.ToProductResponseDeleteAt(product), nil
}

func (s *productService) DeleteProductPermanent(productID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting product", zap.Int("productID", productID))

	res, err := s.productRepository.FindByIdTrashed(productID)

	if err != nil {
		s.logger.Error("Failed to find product",
			zap.Int("product_id", productID),
			zap.Error(err))

		return false, product_errors.ErrFailedFindProductTrashedById
	}

	if res.ImageProduct != "" {
		err := os.Remove(res.ImageProduct)
		if err != nil {
			if os.IsNotExist(err) {
				s.logger.Debug("Product image file not found, continuing with product deletion",
					zap.String("image_path", res.ImageProduct))
			} else {
				s.logger.Debug("Failed to delete product image",
					zap.String("image_path", res.ImageProduct),
					zap.Error(err))

				return false, product_errors.ErrFailedDeleteImageProduct
			}
		} else {
			s.logger.Debug("Successfully deleted product image",
				zap.String("image_path", res.ImageProduct))
		}
	}

	_, err = s.productRepository.DeleteProductPermanent(productID)

	if err != nil {
		s.logger.Error("Failed to permanently delete product",
			zap.Int("product_id", productID),
			zap.Error(err))

		return false, product_errors.ErrFailedDeleteProductPermanent
	}

	s.logger.Debug("Product permanently deleted successfully",
		zap.Int("product_id", productID))

	return true, nil
}

func (s *productService) RestoreAllProducts() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed products")

	success, err := s.productRepository.RestoreAllProducts()

	if err != nil {
		s.logger.Error("Failed to restore all trashed products",
			zap.Error(err))

		return false, product_errors.ErrFailedRestoreAllProducts
	}

	s.logger.Debug("All trashed products restored successfully",
		zap.Bool("success", success))

	return success, nil
}

func (s *productService) DeleteAllProductsPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all products")

	success, err := s.productRepository.DeleteAllProductPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed products",
			zap.Error(err))

		return false, product_errors.ErrFailedDeleteAllProductsPermanent
	}

	s.logger.Debug("All trashed products permanently deleted successfully",
		zap.Bool("success", success))

	return success, nil
}
