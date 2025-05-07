package category_errors

import "errors"

var (
	ErrGetMonthlyTotalPrice           = errors.New("failed to get monthly total price for categories")
	ErrGetYearlyTotalPrices           = errors.New("failed to get yearly total prices for categories")
	ErrGetMonthlyTotalPriceById       = errors.New("failed to get monthly total price by category ID")
	ErrGetYearlyTotalPricesById       = errors.New("failed to get yearly total prices by category ID")
	ErrGetMonthlyTotalPriceByMerchant = errors.New("failed to get monthly total price by merchant")
	ErrGetYearlyTotalPricesByMerchant = errors.New("failed to get yearly total prices by merchant")

	ErrGetMonthPrice           = errors.New("failed to get month price for categories")
	ErrGetYearPrice            = errors.New("failed to get year price for categories")
	ErrGetMonthPriceByMerchant = errors.New("failed to get month price by merchant")
	ErrGetYearPriceByMerchant  = errors.New("failed to get year price by merchant")
	ErrGetMonthPriceById       = errors.New("failed to get month price by category ID")
	ErrGetYearPriceById        = errors.New("failed to get year price by category ID")

	ErrFindAllCategory         = errors.New("failed to find all categories")
	ErrFindByActiveCategory    = errors.New("failed to find active categories")
	ErrFindByTrashedCategory   = errors.New("failed to find trashed categories")
	ErrFindCategoryById        = errors.New("failed to find category by ID")
	ErrFindCategoryByIdTrashed = errors.New("failed to find trashed category by ID")

	ErrCreateCategory            = errors.New("failed to create category")
	ErrUpdateCategory            = errors.New("failed to update category")
	ErrTrashedCategory           = errors.New("failed to move category to trash")
	ErrRestoreCategory           = errors.New("failed to restore category")
	ErrDeleteCategoryPermanently = errors.New("failed to permanently delete category")

	ErrRestoreAllCategories         = errors.New("failed to restore all categories")
	ErrDeleteAllPermanentCategories = errors.New("failed to permanently delete all categories")
)
