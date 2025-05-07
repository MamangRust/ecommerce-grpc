package order_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidYear             = response.NewGrpcError("error", "Invalid year", int(codes.InvalidArgument))
	ErrGrpcInvalidMonth            = response.NewGrpcError("error", "Invalid month", int(codes.InvalidArgument))
	ErrGrpcFailedInvalidMerchantId = response.NewGrpcError("error", "Invalid merchant ID", int(codes.InvalidArgument))
	ErrGrpcFailedInvalidId         = response.NewGrpcError("error", "Invalid ID", int(codes.InvalidArgument))

	ErrGrpcFailedFindMonthlyTotalRevenue           = response.NewGrpcError("error", "Failed to find monthly total revenue", int(codes.Internal))
	ErrGrpcFailedFindYearlyTotalRevenue            = response.NewGrpcError("error", "Failed to find yearly total revenue", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTotalRevenueById       = response.NewGrpcError("error", "Failed to find monthly total revenue by ID", int(codes.Internal))
	ErrGrpcFailedFindYearlyTotalRevenueById        = response.NewGrpcError("error", "Failed to find yearly total revenue by ID", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTotalRevenueByMerchant = response.NewGrpcError("error", "Failed to find monthly total revenue by merchant", int(codes.Internal))
	ErrGrpcFailedFindYearlyTotalRevenueByMerchant  = response.NewGrpcError("error", "Failed to find yearly total revenue by merchant", int(codes.Internal))

	ErrGrpcFailedFindAllOrder        = response.NewGrpcError("error", "Failed to find all orders", int(codes.Internal))
	ErrGrpcFailedFindOrderByMerchant = response.NewGrpcError("error", "Failed to find orders by merchant", int(codes.Internal))
	ErrGrpcFailedFindOrderById       = response.NewGrpcError("error", "Failed to find order by ID", int(codes.Internal))

	ErrGrpcFailedFindMonthlyRevenue           = response.NewGrpcError("error", "Failed to find monthly revenue", int(codes.Internal))
	ErrGrpcFailedFindYearlyRevenue            = response.NewGrpcError("error", "Failed to find yearly revenue", int(codes.Internal))
	ErrGrpcFailedFindMonthlyRevenueByMerchant = response.NewGrpcError("error", "Failed to find monthly revenue by merchant", int(codes.Internal))
	ErrGrpcFailedFindYearlyRevenueByMerchant  = response.NewGrpcError("error", "Failed to find yearly revenue by merchant", int(codes.Internal))

	ErrGrpcFailedFindOrderByActive  = response.NewGrpcError("error", "Failed to find active orders", int(codes.Internal))
	ErrGrpcFailedFindOrderByTrashed = response.NewGrpcError("error", "Failed to find trashed orders", int(codes.Internal))

	ErrGrpcFailedCreateOrder          = response.NewGrpcError("error", "Failed to create order", int(codes.Internal))
	ErrGrpcFailedUpdateOrder          = response.NewGrpcError("error", "Failed to update order", int(codes.Internal))
	ErrGrpcFailedTrashedOrder         = response.NewGrpcError("error", "Failed to trash order", int(codes.Internal))
	ErrGrpcFailedRestoreOrder         = response.NewGrpcError("error", "Failed to restore order", int(codes.Internal))
	ErrGrpcFailedDeleteOrderPermanent = response.NewGrpcError("error", "Failed to permanently delete order", int(codes.Internal))

	ErrGrpcFailedRestoreAllOrder         = response.NewGrpcError("error", "Failed to restore all orders", int(codes.Internal))
	ErrGrpcFailedDeleteAllOrderPermanent = response.NewGrpcError("error", "Failed to permanently delete all orders", int(codes.Internal))

	ErrGrpcValidateCreateOrder = response.NewGrpcError("error", "validation failed: invalid create order request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateOrder = response.NewGrpcError("error", "validation failed: invalid update order request", int(codes.InvalidArgument))
)
