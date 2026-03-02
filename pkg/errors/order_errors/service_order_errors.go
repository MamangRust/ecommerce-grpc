package order_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrFailedNotDeleteAtOrder = errors.NewErrorResponse("Failed to delete at order", http.StatusInternalServerError)

	ErrInsufficientProductStock  = errors.NewErrorResponse("Insufficient product stock", http.StatusBadRequest)
	ErrFailedInvalidCountInStock = errors.NewErrorResponse("Failed to find invalid count in stock", http.StatusInternalServerError)

	ErrFailedFindMonthlyTotalRevenue           = errors.NewErrorResponse("Failed to find monthly total revenue", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalRevenue            = errors.NewErrorResponse("Failed to find yearly total revenue", http.StatusInternalServerError)
	ErrFailedFindMonthlyTotalRevenueById       = errors.NewErrorResponse("Failed to find monthly total revenue by order ID", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalRevenueById        = errors.NewErrorResponse("Failed to find yearly total revenue by order ID", http.StatusInternalServerError)
	ErrFailedFindMonthlyTotalRevenueByMerchant = errors.NewErrorResponse("Failed to find monthly total revenue by merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalRevenueByMerchant  = errors.NewErrorResponse("Failed to find yearly total revenue by merchant", http.StatusInternalServerError)

	ErrFailedFindMonthlyOrder           = errors.NewErrorResponse("Failed to find monthly order", http.StatusInternalServerError)
	ErrFailedFindYearlyOrder            = errors.NewErrorResponse("Failed to find yearly order", http.StatusInternalServerError)
	ErrFailedFindMonthlyOrderByMerchant = errors.NewErrorResponse("Failed to find monthly order by merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyOrderByMerchant  = errors.NewErrorResponse("Failed to find yearly order by merchant", http.StatusInternalServerError)

	ErrFailedFindAllOrders           = errors.NewErrorResponse("Failed to find all orders", http.StatusInternalServerError)
	ErrFailedFindOrderById           = errors.NewErrorResponse("Failed to find order by ID", http.StatusInternalServerError)
	ErrFailedFindOrdersByActive      = errors.NewErrorResponse("Failed to find active orders", http.StatusInternalServerError)
	ErrFailedFindOrdersByTrashed     = errors.NewErrorResponse("Failed to find trashed orders", http.StatusInternalServerError)
	ErrFailedFindOrdersByMerchant    = errors.NewErrorResponse("Failed to find orders by merchant", http.StatusInternalServerError)
	ErrFailedCreateOrder             = errors.NewErrorResponse("Failed to create order", http.StatusInternalServerError)
	ErrFailedUpdateOrder             = errors.NewErrorResponse("Failed to update order", http.StatusInternalServerError)
	ErrFailedTrashOrder              = errors.NewErrorResponse("Failed to trash order", http.StatusInternalServerError)
	ErrFailedRestoreOrder            = errors.NewErrorResponse("Failed to restore order", http.StatusInternalServerError)
	ErrFailedDeleteOrderPermanent    = errors.NewErrorResponse("Failed to permanently delete order", http.StatusInternalServerError)
	ErrFailedRestoreAllOrder         = errors.NewErrorResponse("Failed to restore all orders", http.StatusInternalServerError)
	ErrFailedDeleteAllOrderPermanent = errors.NewErrorResponse("Failed to permanently delete all orders", http.StatusInternalServerError)
)
