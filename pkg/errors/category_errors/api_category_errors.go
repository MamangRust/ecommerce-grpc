package category_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiCategoryNameRequired = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "validation_error", "Category name is required", http.StatusBadRequest)
	}

	ErrApiCategoryDescriptionRequired = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "validation_error", "Description is required", http.StatusBadRequest)
	}

	ErrApiCategorySlugRequired = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "validation_error", "Slug category is required", http.StatusBadRequest)
	}

	ErrApiCategoryImageRequired = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "image_required", "A category image is required", http.StatusBadRequest)
	}

	ErrApiCategoryInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid category ID", http.StatusBadRequest)
	}

	ErrApiCategoryInvalidYear = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid year parameter", http.StatusBadRequest)
	}

	ErrApiCategoryInvalidMonth = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid month parameter", http.StatusBadRequest)
	}

	ErrApiCategoryInvalidMerchantId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid merchant ID", http.StatusBadRequest)
	}

	ErrApiFailedFindAllCategory = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch categories", http.StatusInternalServerError)
	}

	ErrApiFailedFindByIdCategory = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find category by ID", http.StatusInternalServerError)
	}

	ErrApiFailedFindByActiveCategory = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch active categories", http.StatusInternalServerError)
	}

	ErrApiFailedFindByTrashedCategory = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch trashed categories", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthTotalPrice = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch monthly total price", http.StatusInternalServerError)
	}

	ErrApiFailedFindYearTotalPrice = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch yearly total price", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthTotalPriceByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch monthly total price by merchant", http.StatusInternalServerError)
	}

	ErrApiFailedFindYearTotalPriceByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch yearly total price by merchant", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthTotalPriceById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch monthly total price by category ID", http.StatusInternalServerError)
	}

	ErrApiFailedFindYearTotalPriceById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch yearly total price by category ID", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthPrice = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch month price", http.StatusInternalServerError)
	}

	ErrApiFailedFindYearPrice = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch year price", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthPriceByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch month price by merchant", http.StatusInternalServerError)
	}

	ErrApiFailedFindYearPriceByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch year price by merchant", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthPriceById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch month price by category ID", http.StatusInternalServerError)
	}

	ErrApiFailedFindYearPriceById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch year price by category ID", http.StatusInternalServerError)
	}

	ErrApiFailedCreateCategory = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create category", http.StatusInternalServerError)
	}

	ErrApiFailedUpdateCategory = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update category", http.StatusInternalServerError)
	}

	ErrApiInvalidBody = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid request body", http.StatusBadRequest)
	}

	ErrApiFailedTrashedCategory = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to trash category", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreCategory = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore category", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteCategoryPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete category", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllCategories = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all categories", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteAllCategoriesPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete all categories", http.StatusInternalServerError)
	}
)
