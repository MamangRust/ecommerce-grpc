package category_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrFailedFindMonthlyTotalPrice           = errors.NewErrorResponse("Failed to find monthly total price", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalPrice            = errors.NewErrorResponse("Failed to find yearly total price", http.StatusInternalServerError)
	ErrFailedFindMonthlyTotalPriceById       = errors.NewErrorResponse("Failed to find monthly total price by category ID", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalPriceById        = errors.NewErrorResponse("Failed to find yearly total price by category ID", http.StatusInternalServerError)
	ErrFailedFindMonthlyTotalPriceByMerchant = errors.NewErrorResponse("Failed to find monthly total price by merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalPriceByMerchant  = errors.NewErrorResponse("Failed to find yearly total price by merchant", http.StatusInternalServerError)

	ErrFailedFindMonthPrice           = errors.NewErrorResponse("Failed to find monthly price", http.StatusInternalServerError)
	ErrFailedFindYearPrice            = errors.NewErrorResponse("Failed to find yearly price", http.StatusInternalServerError)
	ErrFailedFindMonthPriceByMerchant = errors.NewErrorResponse("Failed to find monthly price by merchant", http.StatusInternalServerError)
	ErrFailedFindYearPriceByMerchant  = errors.NewErrorResponse("Failed to find yearly price by merchant", http.StatusInternalServerError)
	ErrFailedFindMonthPriceById       = errors.NewErrorResponse("Failed to find monthly price by category ID", http.StatusInternalServerError)
	ErrFailedFindYearPriceById        = errors.NewErrorResponse("Failed to find yearly price by category ID", http.StatusInternalServerError)

	ErrFailedFindAllCategories     = errors.NewErrorResponse("Failed to find all categories", http.StatusInternalServerError)
	ErrFailedFindActiveCategories  = errors.NewErrorResponse("Failed to find active categories", http.StatusInternalServerError)
	ErrFailedFindTrashedCategories = errors.NewErrorResponse("Failed to find trashed categories", http.StatusInternalServerError)
	ErrFailedFindCategoryById      = errors.NewErrorResponse("Failed to find category by ID", http.StatusInternalServerError)
	ErrFailedFindCategoryIdTrashed = errors.NewErrorResponse("Failed to find category ID trashed", http.StatusInternalServerError)
	ErrFailedRemoveImageCategory   = errors.NewErrorResponse("Failed to remove image category", http.StatusInternalServerError)

	ErrFailedCreateCategory               = errors.NewErrorResponse("Failed to create category", http.StatusInternalServerError)
	ErrFailedUpdateCategory               = errors.NewErrorResponse("Failed to update category", http.StatusInternalServerError)
	ErrFailedTrashedCategory              = errors.NewErrorResponse("Failed to move category to trash", http.StatusInternalServerError)
	ErrFailedRestoreCategory              = errors.NewErrorResponse("Failed to restore category", http.StatusInternalServerError)
	ErrFailedDeleteCategoryPermanent      = errors.NewErrorResponse("Failed to permanently delete category", http.StatusInternalServerError)
	ErrFailedRestoreAllCategories         = errors.NewErrorResponse("Failed to restore all categories", http.StatusInternalServerError)
	ErrFailedDeleteAllCategoriesPermanent = errors.NewErrorResponse("Failed to permanently delete all categories", http.StatusInternalServerError)
)
