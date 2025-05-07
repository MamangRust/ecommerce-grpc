package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/errors/category_errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/utils"
	"os"

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

func (s *categoryService) FindAll(req *requests.FindAllCategory) ([]*response.CategoryResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all categories",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	category, totalRecords, err := s.categoryRepository.FindAllCategory(req)

	if err != nil {
		s.logger.Error("Failed to fetch category",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		return nil, nil, category_errors.ErrFailedFindAllCategories
	}

	categoriesResponse := s.mapping.ToCategorysResponse(category)

	s.logger.Debug("Successfully fetched category",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return categoriesResponse, totalRecords, nil
}

func (s *categoryService) FindByActive(req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching categories active",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	category, totalRecords, err := s.categoryRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active categories",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		return nil, nil, category_errors.ErrFailedFindActiveCategories
	}

	s.logger.Debug("Successfully fetched categories",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToCategorysResponseDeleteAt(category), totalRecords, nil
}

func (s *categoryService) FindByTrashed(req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching categories trashed",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	categories, totalRecords, err := s.categoryRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to fetch categories",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		return nil, nil, category_errors.ErrFailedFindTrashedCategories
	}

	s.logger.Debug("Successfully fetched categories",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToCategorysResponseDeleteAt(categories), totalRecords, nil
}

func (s *categoryService) FindById(category_id int) (*response.CategoryResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching category by ID", zap.Int("category_id", category_id))

	category, err := s.categoryRepository.FindById(category_id)

	if err != nil {
		s.logger.Error("Failed to retrieve category details",
			zap.Error(err),
			zap.Int("category_id", category_id))

		return nil, category_errors.ErrFailedFindCategoryById
	}

	return s.mapping.ToCategoryResponse(category), nil
}

func (s *categoryService) FindMonthlyTotalPrice(req *requests.MonthTotalPrice) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	res, err := s.categoryRepository.GetMonthlyTotalPrice(req)

	if err != nil {
		s.logger.Error("failed to get monthly total sales",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Error(err))
		return nil, category_errors.ErrFailedFindMonthlyTotalPrice
	}

	return s.mapping.ToCategoryMonthlyTotalPrices(res), nil
}

func (s *categoryService) FindYearlyTotalPrice(year int) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse) {
	res, err := s.categoryRepository.GetYearlyTotalPrices(year)

	if err != nil {
		s.logger.Error("failed to get yearly total sales",
			zap.Int("year", year),
			zap.Error(err))
		return nil, category_errors.ErrFailedFindYearlyTotalPrice
	}

	return s.mapping.ToCategoryYearlyTotalPrices(res), nil
}

func (s *categoryService) FindMonthlyTotalPriceById(req *requests.MonthTotalPriceCategory) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	res, err := s.categoryRepository.GetMonthlyTotalPriceById(req)

	if err != nil {
		s.logger.Error("failed to get monthly total price",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Error(err))
		return nil, category_errors.ErrFailedFindMonthlyTotalPriceById
	}

	return s.mapping.ToCategoryMonthlyTotalPrices(res), nil
}

func (s *categoryService) FindYearlyTotalPriceById(req *requests.YearTotalPriceCategory) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse) {
	year := req.Year

	res, err := s.categoryRepository.GetYearlyTotalPricesById(req)

	if err != nil {
		s.logger.Error("failed to get yearly total price",
			zap.Int("year", year),
			zap.Error(err))

		return nil, category_errors.ErrFailedFindYearlyTotalPriceById
	}

	return s.mapping.ToCategoryYearlyTotalPrices(res), nil
}

func (s *categoryService) FindMonthlyTotalPriceByMerchant(req *requests.MonthTotalPriceMerchant) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	res, err := s.categoryRepository.GetMonthlyTotalPriceByMerchant(req)

	if err != nil {
		s.logger.Error("failed to get monthly total price",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Error(err))

		return nil, category_errors.ErrFailedFindMonthlyTotalPriceByMerchant
	}

	return s.mapping.ToCategoryMonthlyTotalPrices(res), nil
}

func (s *categoryService) FindYearlyTotalPriceByMerchant(req *requests.YearTotalPriceMerchant) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse) {
	year := req.Year

	res, err := s.categoryRepository.GetYearlyTotalPricesByMerchant(req)

	if err != nil {
		s.logger.Error("failed to get yearly total price",
			zap.Int("year", year),
			zap.Error(err))

		return nil, category_errors.ErrFailedFindYearlyTotalPriceByMerchant
	}

	return s.mapping.ToCategoryYearlyTotalPrices(res), nil
}

func (s *categoryService) FindMonthPrice(year int) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse) {
	res, err := s.categoryRepository.GetMonthPrice(year)

	if err != nil {
		s.logger.Error("failed to get monthly category prices",
			zap.Int("year", year),
			zap.Error(err))

		return nil, category_errors.ErrFailedFindMonthPrice
	}

	return s.mapping.ToCategoryMonthlyPrices(res), nil
}

func (s *categoryService) FindYearPrice(year int) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse) {
	res, err := s.categoryRepository.GetYearPrice(year)

	if err != nil {
		s.logger.Error("failed to get yearly category prices",
			zap.Int("year", year),
			zap.Error(err))

		return nil, category_errors.ErrFailedFindYearPrice
	}

	return s.mapping.ToCategoryYearlyPrices(res), nil
}

func (s *categoryService) FindMonthPriceByMerchant(req *requests.MonthPriceMerchant) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse) {
	year := req.Year
	merchant_id := req.MerchantID

	res, err := s.categoryRepository.GetMonthPriceByMerchant(req)

	if err != nil {
		s.logger.Error("failed to get monthly category prices by merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchant_id),
			zap.Error(err))

		return nil, category_errors.ErrFailedFindMonthPriceByMerchant
	}

	return s.mapping.ToCategoryMonthlyPrices(res), nil
}

func (s *categoryService) FindYearPriceByMerchant(req *requests.YearPriceMerchant) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse) {
	year := req.Year
	merchant_id := req.MerchantID

	res, err := s.categoryRepository.GetYearPriceByMerchant(req)

	if err != nil {
		s.logger.Error("failed to get yearly category prices by merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchant_id),
			zap.Error(err))

		return nil, category_errors.ErrFailedFindYearPriceByMerchant
	}

	return s.mapping.ToCategoryYearlyPrices(res), nil
}

func (s *categoryService) FindMonthPriceById(req *requests.MonthPriceId) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse) {
	year := req.Year
	category_id := req.CategoryID

	res, err := s.categoryRepository.GetMonthPriceById(req)

	if err != nil {
		s.logger.Error("failed to get monthly category prices by ID",
			zap.Int("year", year),
			zap.Int("category_id", category_id),
			zap.Error(err))
		return nil, category_errors.ErrFailedFindMonthPriceById
	}

	return s.mapping.ToCategoryMonthlyPrices(res), nil
}

func (s *categoryService) FindYearPriceById(req *requests.YearPriceId) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse) {
	year := req.Year
	category_id := req.CategoryID

	res, err := s.categoryRepository.GetYearPriceById(req)

	if err != nil {
		s.logger.Error("failed to get yearly category prices by ID",
			zap.Int("year", year),
			zap.Int("category_id", category_id),
			zap.Error(err))
		return nil, category_errors.ErrFailedFindYearPriceById
	}

	return s.mapping.ToCategoryYearlyPrices(res), nil
}

func (s *categoryService) CreateCategory(req *requests.CreateCategoryRequest) (*response.CategoryResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new category")

	slug := utils.GenerateSlug(req.Name)

	req.Name = slug

	cashier, err := s.categoryRepository.CreateCategory(req)

	if err != nil {
		s.logger.Error("Failed to create category", zap.Error(err))

		return nil, category_errors.ErrFailedCreateCategory
	}

	return s.mapping.ToCategoryResponse(cashier), nil
}

func (s *categoryService) UpdateCategory(req *requests.UpdateCategoryRequest) (*response.CategoryResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating category", zap.Int("category_id", *req.CategoryID))

	slug := utils.GenerateSlug(req.Name)

	req.Name = slug

	category, err := s.categoryRepository.UpdateCategory(req)

	if err != nil {
		s.logger.Error("Failed to update category",
			zap.Error(err),
			zap.Any("request", req))

		return nil, category_errors.ErrFailedUpdateCategory
	}

	return s.mapping.ToCategoryResponse(category), nil
}

func (s *categoryService) TrashedCategory(category_id int) (*response.CategoryResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing category", zap.Int("category", category_id))

	category, err := s.categoryRepository.TrashedCategory(category_id)

	if err != nil {
		s.logger.Error("Failed to move category to trash",
			zap.Error(err),
			zap.Int("category_id", category_id))

		return nil, category_errors.ErrFailedTrashedCategory
	}

	return s.mapping.ToCategoryResponseDeleteAt(category), nil
}

func (s *categoryService) RestoreCategory(categoryID int) (*response.CategoryResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring category", zap.Int("categoryID", categoryID))

	category, err := s.categoryRepository.RestoreCategory(categoryID)

	if err != nil {
		s.logger.Error("Failed to restore category from trash",
			zap.Error(err),
			zap.Int("category_id", categoryID))

		return nil, category_errors.ErrFailedRestoreCategory
	}

	return s.mapping.ToCategoryResponseDeleteAt(category), nil
}

func (s *categoryService) DeleteCategoryPermanent(categoryID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting category", zap.Int("categoryID", categoryID))

	res, err := s.categoryRepository.FindByIdTrashed(categoryID)

	if err != nil {
		s.logger.Error("Failed to find category",
			zap.Int("category_id", categoryID),
			zap.Error(err))

		return false, category_errors.ErrFailedFindCategoryIdTrashed
	}

	if res.ImageCategory != "" {
		err := os.Remove(res.ImageCategory)
		if err != nil {
			if os.IsNotExist(err) {
				s.logger.Debug("Category image file not found, continuing with category deletion",
					zap.String("image_path", res.ImageCategory))
			} else {
				s.logger.Debug("Failed to delete category image",
					zap.String("image_path", res.ImageCategory),
					zap.Error(err))

				return false, category_errors.ErrFailedRemoveImageCategory
			}
		} else {
			s.logger.Debug("Successfully deleted category image",
				zap.String("image_path", res.ImageCategory))
		}
	}

	success, err := s.categoryRepository.DeleteCategoryPermanently(categoryID)

	if err != nil {
		s.logger.Error("Failed to permanently delete category",
			zap.Error(err),
			zap.Int("category_id", categoryID))

		return false, category_errors.ErrFailedDeleteCategoryPermanent
	}

	return success, nil
}

func (s *categoryService) RestoreAllCategories() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed categories")

	success, err := s.categoryRepository.RestoreAllCategories()

	if err != nil {
		s.logger.Error("Failed to restore all trashed categories",
			zap.Error(err))

		return false, category_errors.ErrFailedRestoreAllCategories
	}

	return success, nil
}

func (s *categoryService) DeleteAllCategoriesPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all categories")

	success, err := s.categoryRepository.DeleteAllPermanentCategories()

	if err != nil {
		s.logger.Error("Failed to permanently delete all categories", zap.Error(err))

		return false, category_errors.ErrFailedDeleteAllCategoriesPermanent
	}

	return success, nil
}
