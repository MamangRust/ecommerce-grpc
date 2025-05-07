package category_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcFailedFindMonthlyTotalPrices           = response.NewGrpcError("error", "Failed to fetch monthly total prices", int(codes.Internal))
	ErrGrpcFailedFindYearlyTotalPrices            = response.NewGrpcError("error", "Failed to fetch yearly total prices", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTotalPricesById       = response.NewGrpcError("error", "Failed to fetch monthly total prices by category ID", int(codes.Internal))
	ErrGrpcFailedFindYearlyTotalPricesById        = response.NewGrpcError("error", "Failed to fetch yearly total prices by category ID", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTotalPricesByMerchant = response.NewGrpcError("error", "Failed to fetch monthly total prices by merchant", int(codes.Internal))
	ErrGrpcFailedFindYearlyTotalPricesByMerchant  = response.NewGrpcError("error", "Failed to fetch yearly total prices by merchant", int(codes.Internal))

	ErrGrpcFailedFindMonthPrice           = response.NewGrpcError("error", "Failed to fetch month price", int(codes.Internal))
	ErrGrpcFailedFindYearPrice            = response.NewGrpcError("error", "Failed to fetch year price", int(codes.Internal))
	ErrGrpcFailedFindMonthPriceByMerchant = response.NewGrpcError("error", "Failed to fetch month price by merchant", int(codes.Internal))
	ErrGrpcFailedFindYearPriceByMerchant  = response.NewGrpcError("error", "Failed to fetch year price by merchant", int(codes.Internal))
	ErrGrpcFailedFindMonthPriceById       = response.NewGrpcError("error", "Failed to fetch month price by category ID", int(codes.Internal))
	ErrGrpcFailedFindYearPriceById        = response.NewGrpcError("error", "Failed to fetch year price by category ID", int(codes.Internal))

	ErrGrpcFailedFindAll       = response.NewGrpcError("error", "Failed to fetch all categories", int(codes.Internal))
	ErrGrpcFailedFindById      = response.NewGrpcError("error", "Failed to find category by ID", int(codes.Internal))
	ErrGrpcFailedFindByActive  = response.NewGrpcError("error", "Failed to fetch active categories", int(codes.Internal))
	ErrGrpcFailedFindByTrashed = response.NewGrpcError("error", "Failed to fetch trashed categories", int(codes.Internal))

	ErrGrpcFailedCreateCategory          = response.NewGrpcError("error", "Failed to create category", int(codes.Internal))
	ErrGrpcValidateCreateCategory        = response.NewGrpcError("error", "Validation failed: invalid create category request", int(codes.InvalidArgument))
	ErrGrpcFailedUpdateCategory          = response.NewGrpcError("error", "Failed to update category", int(codes.Internal))
	ErrGrpcFailedTrashedCategory         = response.NewGrpcError("error", "Failed to trash category", int(codes.Internal))
	ErrGrpcFailedRestoreCategory         = response.NewGrpcError("error", "Failed to restore category", int(codes.Internal))
	ErrGrpcFailedDeleteCategoryPermanent = response.NewGrpcError("error", "Failed to permanently delete category", int(codes.Internal))

	ErrGrpcFailedRestoreAllCategory         = response.NewGrpcError("error", "Failed to restore all categories", int(codes.Internal))
	ErrGrpcFailedDeleteAllCategoryPermanent = response.NewGrpcError("error", "Failed to permanently delete all categories", int(codes.Internal))

	ErrGrpcCategoryNotFound          = response.NewGrpcError("error", "Category not found", int(codes.NotFound))
	ErrGrpcCategoryInvalidId         = response.NewGrpcError("error", "Invalid category ID", int(codes.InvalidArgument))
	ErrGrpcCategoryInvalidYear       = response.NewGrpcError("error", "Invalid year", int(codes.InvalidArgument))
	ErrGrpcCategoryInvalidMonth      = response.NewGrpcError("error", "Invalid month", int(codes.InvalidArgument))
	ErrGrpcCategoryInvalidMerchantId = response.NewGrpcError("error", "Invalid merchant ID", int(codes.InvalidArgument))
)
