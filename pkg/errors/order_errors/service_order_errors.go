package order_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrFailedNotDeleteAtOrder = response.NewErrorResponse("Failed to delete at order", http.StatusInternalServerError)

	ErrInsufficientProductStock  = response.NewErrorResponse("Insufficient product stock", http.StatusBadRequest)
	ErrFailedInvalidCountInStock = response.NewErrorResponse("Failed to find invalid count in stock", http.StatusInternalServerError)

	ErrFailedFindMonthlyTotalRevenue           = response.NewErrorResponse("Failed to find monthly total revenue", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalRevenue            = response.NewErrorResponse("Failed to find yearly total revenue", http.StatusInternalServerError)
	ErrFailedFindMonthlyTotalRevenueById       = response.NewErrorResponse("Failed to find monthly total revenue by order ID", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalRevenueById        = response.NewErrorResponse("Failed to find yearly total revenue by order ID", http.StatusInternalServerError)
	ErrFailedFindMonthlyTotalRevenueByMerchant = response.NewErrorResponse("Failed to find monthly total revenue by merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalRevenueByMerchant  = response.NewErrorResponse("Failed to find yearly total revenue by merchant", http.StatusInternalServerError)

	ErrFailedFindMonthlyOrder           = response.NewErrorResponse("Failed to find monthly order", http.StatusInternalServerError)
	ErrFailedFindYearlyOrder            = response.NewErrorResponse("Failed to find yearly order", http.StatusInternalServerError)
	ErrFailedFindMonthlyOrderByMerchant = response.NewErrorResponse("Failed to find monthly order by merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyOrderByMerchant  = response.NewErrorResponse("Failed to find yearly order by merchant", http.StatusInternalServerError)

	ErrFailedFindAllOrders           = response.NewErrorResponse("Failed to find all orders", http.StatusInternalServerError)
	ErrFailedFindOrderById           = response.NewErrorResponse("Failed to find order by ID", http.StatusInternalServerError)
	ErrFailedFindOrdersByActive      = response.NewErrorResponse("Failed to find active orders", http.StatusInternalServerError)
	ErrFailedFindOrdersByTrashed     = response.NewErrorResponse("Failed to find trashed orders", http.StatusInternalServerError)
	ErrFailedFindOrdersByMerchant    = response.NewErrorResponse("Failed to find orders by merchant", http.StatusInternalServerError)
	ErrFailedCreateOrder             = response.NewErrorResponse("Failed to create order", http.StatusInternalServerError)
	ErrFailedUpdateOrder             = response.NewErrorResponse("Failed to update order", http.StatusInternalServerError)
	ErrFailedTrashOrder              = response.NewErrorResponse("Failed to trash order", http.StatusInternalServerError)
	ErrFailedRestoreOrder            = response.NewErrorResponse("Failed to restore order", http.StatusInternalServerError)
	ErrFailedDeleteOrderPermanent    = response.NewErrorResponse("Failed to permanently delete order", http.StatusInternalServerError)
	ErrFailedRestoreAllOrder         = response.NewErrorResponse("Failed to restore all orders", http.StatusInternalServerError)
	ErrFailedDeleteAllOrderPermanent = response.NewErrorResponse("Failed to permanently delete all orders", http.StatusInternalServerError)
)
