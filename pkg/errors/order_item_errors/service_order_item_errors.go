package orderitem_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrFailedOrderItemEmptyId     = response.NewErrorResponse("Order item ID is empty", http.StatusBadRequest)
	ErrFailedNotDeleteAtOrderItem = response.NewErrorResponse("Failed to delete at order item", http.StatusInternalServerError)
	ErrFailedOrderItemEmpty       = response.NewErrorResponse("Failed to find order item", http.StatusInternalServerError)
	ErrFailedInvalidQuantity      = response.NewErrorResponse("Invalid quantity", http.StatusBadRequest)

	ErrFailedOrderItemNotFound = response.NewErrorResponse("Order item not found", http.StatusNotFound)
	ErrFailedTrashedOrderItem  = response.NewErrorResponse("Order item is already trashed", http.StatusBadRequest)
	ErrFailedRestoreOrderItem  = response.NewErrorResponse("Failed to restore order item", http.StatusInternalServerError)
	ErrFailedDeleteOrderItem   = response.NewErrorResponse("Failed to delete order item", http.StatusInternalServerError)

	ErrFailedCreateOrderItem = response.NewErrorResponse("Failed to create order item", http.StatusInternalServerError)
	ErrFailedUpdateOrderItem = response.NewErrorResponse("Failed to update order item", http.StatusInternalServerError)
	ErrFailedCalculateTotal  = response.NewErrorResponse("Failed to calculate total", http.StatusInternalServerError)

	ErrFailedFindAllOrderItems       = response.NewErrorResponse("Failed to find all order items", http.StatusInternalServerError)
	ErrFailedFindOrderItemsByActive  = response.NewErrorResponse("Failed to find active order items", http.StatusInternalServerError)
	ErrFailedFindOrderItemsByTrashed = response.NewErrorResponse("Failed to find trashed order items", http.StatusInternalServerError)
	ErrFailedFindOrderItemByOrder    = response.NewErrorResponse("Failed to find order items by order ID", http.StatusInternalServerError)

	ErrFailedRestoreAllOrderItem = response.NewErrorResponse("Failed to restore all order items", http.StatusInternalServerError)
	ErrFailedDeleteAllOrderItem  = response.NewErrorResponse("Failed to delete all order items", http.StatusInternalServerError)
)
