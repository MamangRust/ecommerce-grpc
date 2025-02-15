package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"

	"go.uber.org/zap"
)

type categoryService struct {
	categoryRepository repository.CategoryRepository
	logger             logger.LoggerInterface
	mapping            response_service.CategoryResponseMapper
}

func NewCategoryService(
	categoryRepository repository.CategoryRepository,
	logger logger.LoggerInterface,
	mapping response_service.CategoryResponseMapper,
) *categoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
		logger:             logger,
		mapping:            mapping,
	}
}

func (s *categoryService) FindAll(page int, pageSize int, search string) ([]*response.CategoryResponse, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching category",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	category, totalRecords, err := s.categoryRepository.FindAllCategory(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch category",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch category",
		}
	}

	categoriesResponse := s.mapping.ToCategorysResponse(category)

	s.logger.Debug("Successfully fetched category",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return categoriesResponse, int(totalRecords), nil
}

func (s *categoryService) FindById(category_id int) (*response.CategoryResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching category by ID", zap.Int("category_id", category_id))

	category, err := s.categoryRepository.FindById(category_id)
	if err != nil {
		s.logger.Error("Failed to fetch category", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "category not found"}
	}

	return s.mapping.ToCategoryResponse(category), nil
}

func (s *categoryService) FindByActive(search string, page, pageSize int) ([]*response.CategoryResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching categories",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	cashiers, totalRecords, err := s.categoryRepository.FindByActive(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch categories",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch active categories"}
	}

	s.logger.Debug("Successfully fetched categories",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToCategorysResponseDeleteAt(cashiers), totalRecords, nil
}

func (s *categoryService) FindByTrashed(search string, page, pageSize int) ([]*response.CategoryResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching categories",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	categories, totalRecords, err := s.categoryRepository.FindByTrashed(search, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch cashier",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch trashed categories"}
	}

	s.logger.Debug("Successfully fetched categories",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToCategorysResponseDeleteAt(categories), totalRecords, nil
}

func (s *categoryService) CreateCategory(req *requests.CreateCategoryRequest) (*response.CategoryResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new cashier")

	cashier, err := s.categoryRepository.CreateCategory(req)
	if err != nil {
		s.logger.Error("Failed to create cashier", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to create cashier"}
	}

	return s.mapping.ToCategoryResponse(cashier), nil
}

func (s *categoryService) UpdateCategory(req *requests.UpdateCategoryRequest) (*response.CategoryResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating category", zap.Int("category_id", req.CategoryID))

	category, err := s.categoryRepository.UpdateCategory(req)
	if err != nil {
		s.logger.Error("Failed to update category", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to update category"}
	}

	return s.mapping.ToCategoryResponse(category), nil
}

func (s *categoryService) TrashedCategory(category_id int) (*response.CategoryResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing category", zap.Int("category", category_id))

	category, err := s.categoryRepository.TrashedCategory(category_id)
	if err != nil {
		s.logger.Error("Failed to trash category", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to trash category"}
	}

	return s.mapping.ToCategoryResponseDeleteAt(category), nil
}

func (s *categoryService) RestoreCategory(categoryID int) (*response.CategoryResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring category", zap.Int("categoryID", categoryID))

	category, err := s.categoryRepository.RestoreCategory(categoryID)
	if err != nil {
		s.logger.Error("Failed to restore category", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to restore category"}
	}

	return s.mapping.ToCategoryResponseDeleteAt(category), nil
}

func (s *categoryService) DeleteCategoryPermanent(categoryID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting category", zap.Int("categoryID", categoryID))

	success, err := s.categoryRepository.DeleteCategoryPermanently(categoryID)
	if err != nil {
		s.logger.Error("Failed to permanently delete category", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete category"}
	}

	return success, nil
}

func (s *categoryService) RestoreAllCategories() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed categories")

	success, err := s.categoryRepository.RestoreAllCategories()
	if err != nil {
		s.logger.Error("Failed to restore all categories", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to restore all categories"}
	}

	return success, nil
}

func (s *categoryService) DeleteAllCategoriesPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all categories")

	success, err := s.categoryRepository.DeleteAllPermanentCategories()
	if err != nil {
		s.logger.Error("Failed to permanently delete all categories", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete all categories"}
	}

	return success, nil
}
