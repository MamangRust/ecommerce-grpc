package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"

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

func (s *productService) FindAll(page, pageSize int, search string) ([]*response.ProductResponse, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching products",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	products, totalRecords, err := s.productRepository.FindAllProducts(search, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch products",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch products"}
	}

	s.logger.Debug("Successfully fetched products",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToProductsResponse(products), totalRecords, nil
}

func (s *productService) FindByMerchant(merchant_id int, page, pageSize int, search string) ([]*response.ProductResponse, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching products",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	products, totalRecords, err := s.productRepository.FindByMerchant(merchant_id, search, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch products",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch products"}
	}

	s.logger.Debug("Successfully fetched products",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToProductsResponse(products), totalRecords, nil
}

func (s *productService) FindByCategory(category_name string, page, pageSize int, search string) ([]*response.ProductResponse, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching products",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	products, totalRecords, err := s.productRepository.FindByCategory(category_name, search, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch products",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch products"}
	}

	s.logger.Debug("Successfully fetched products",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToProductsResponse(products), totalRecords, nil
}

func (s *productService) FindById(productID int) (*response.ProductResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching product by ID", zap.Int("productID", productID))

	product, err := s.productRepository.FindById(productID)
	if err != nil {
		s.logger.Error("Failed to fetch product", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Product not found"}
	}

	return s.mapping.ToProductResponse(product), nil
}

func (s *productService) FindByActive(search string, page, pageSize int) ([]*response.ProductResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching products",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	products, totalRecords, err := s.productRepository.FindByActive(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch products",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch active products"}
	}

	s.logger.Debug("Successfully fetched products",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToProductsResponseDeleteAt(products), totalRecords, nil
}

func (s *productService) FindByTrashed(search string, page, pageSize int) ([]*response.ProductResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching products",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	products, totalRecords, err := s.productRepository.FindByTrashed(search, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch products",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch trashed products"}
	}

	s.logger.Debug("Successfully fetched products",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToProductsResponseDeleteAt(products), totalRecords, nil
}

func (s *productService) CreateProduct(req *requests.CreateProductRequest) (*response.ProductResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new product")

	_, err := s.categoryRepository.FindById(req.CategoryID)
	if err != nil {
		s.logger.Error("Category not found", zap.Int("categoryID", req.CategoryID), zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Category not found"}
	}

	_, err = s.merchantRepository.FindById(req.MerchantID)
	if err != nil {
		s.logger.Error("Merchant not found", zap.Int("merchantID", req.MerchantID), zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Merchant not found"}
	}

	product, err := s.productRepository.CreateProduct(req)
	if err != nil {
		s.logger.Error("Failed to create product", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to create product"}
	}

	s.logger.Debug("Product created successfully", zap.Int("productID", product.ID))
	return s.mapping.ToProductResponse(product), nil
}

func (s *productService) UpdateProduct(req *requests.UpdateProductRequest) (*response.ProductResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating product", zap.Int("productID", req.ProductID))

	_, err := s.categoryRepository.FindById(req.CategoryID)
	if err != nil {
		s.logger.Error("Category not found", zap.Int("categoryID", req.CategoryID), zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Category not found"}
	}

	_, err = s.merchantRepository.FindById(req.MerchantID)
	if err != nil {
		s.logger.Error("Merchant not found", zap.Int("merchantID", req.MerchantID), zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Merchant not found"}
	}

	product, err := s.productRepository.UpdateProduct(req)
	if err != nil {
		s.logger.Error("Failed to update product", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to update product"}
	}

	s.logger.Debug("Product updated successfully", zap.Int("productID", product.ID))
	return s.mapping.ToProductResponse(product), nil
}

func (s *productService) TrashProduct(productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing product", zap.Int("productID", productID))

	product, err := s.productRepository.TrashedProduct(productID)
	if err != nil {
		s.logger.Error("Failed to trash product", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to trash product"}
	}

	return s.mapping.ToProductResponseDeleteAt(product), nil
}

func (s *productService) RestoreProduct(productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring product", zap.Int("productID", productID))

	product, err := s.productRepository.RestoreProduct(productID)
	if err != nil {
		s.logger.Error("Failed to restore product", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to restore product"}
	}

	return s.mapping.ToProductResponseDeleteAt(product), nil
}

func (s *productService) DeleteProductPermanent(productID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting product", zap.Int("productID", productID))

	success, err := s.productRepository.DeleteProductPermanent(productID)
	if err != nil {
		s.logger.Error("Failed to permanently delete product", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete product"}
	}

	return success, nil
}

func (s *productService) RestoreAllProducts() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed products")

	success, err := s.productRepository.RestoreAllProducts()
	if err != nil {
		s.logger.Error("Failed to restore all products", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to restore all products"}
	}

	return success, nil
}

func (s *productService) DeleteAllProductsPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all products")

	success, err := s.productRepository.DeleteAllProductPermanent()
	if err != nil {
		s.logger.Error("Failed to permanently delete all products", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete all products"}
	}

	return success, nil
}
