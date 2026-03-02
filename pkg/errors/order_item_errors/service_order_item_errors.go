package orderitem_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrFailedOrderItemEmptyId     = errors.NewErrorResponse("Order item ID is empty", http.StatusBadRequest)
	ErrFailedNotDeleteAtOrderItem = errors.NewErrorResponse("Failed to delete at order item", http.StatusInternalServerError)
	ErrFailedOrderItemEmpty       = errors.NewErrorResponse("Failed to find order item", http.StatusInternalServerError)
	ErrFailedInvalidQuantity      = errors.NewErrorResponse("Invalid quantity", http.StatusBadRequest)

	ErrFailedOrderItemNotFound = errors.NewErrorResponse("Order item not found", http.StatusNotFound)
	ErrFailedTrashedOrderItem  = errors.NewErrorResponse("Order item is already trashed", http.StatusBadRequest)
	ErrFailedRestoreOrderItem  = errors.NewErrorResponse("Failed to restore order item", http.StatusInternalServerError)
	ErrFailedDeleteOrderItem   = errors.NewErrorResponse("Failed to delete order item", http.StatusInternalServerError)

	ErrFailedCreateOrderItem = errors.NewErrorResponse("Failed to create order item", http.StatusInternalServerError)
	ErrFailedUpdateOrderItem = errors.NewErrorResponse("Failed to update order item", http.StatusInternalServerError)
	ErrFailedCalculateTotal  = errors.NewErrorResponse("Failed to calculate total", http.StatusInternalServerError)

	ErrFailedFindAllOrderItems       = errors.NewErrorResponse("Failed to find all order items", http.StatusInternalServerError)
	ErrFailedFindOrderItemsByActive  = errors.NewErrorResponse("Failed to find active order items", http.StatusInternalServerError)
	ErrFailedFindOrderItemsByTrashed = errors.NewErrorResponse("Failed to find trashed order items", http.StatusInternalServerError)
	ErrFailedFindOrderItemByOrder    = errors.NewErrorResponse("Failed to find order items by order ID", http.StatusInternalServerError)

	ErrFailedRestoreAllOrderItem = errors.NewErrorResponse("Failed to restore all order items", http.StatusInternalServerError)
	ErrFailedDeleteAllOrderItem  = errors.NewErrorResponse("Failed to delete all order items", http.StatusInternalServerError)
)
