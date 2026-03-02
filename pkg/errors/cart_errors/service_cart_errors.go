package cart_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrCartNotFoundRes   = errors.NewErrorResponse("Cart not found", http.StatusNotFound)
	ErrCartAlreadyExists = errors.NewErrorResponse("Cart already exists", http.StatusBadRequest)
	ErrCartInvalidData   = errors.NewErrorResponse("Invalid cart data", http.StatusBadRequest)

	ErrFailedFindAllCarts = errors.NewErrorResponse("Failed to fetch carts", http.StatusInternalServerError)
	ErrFailedCreateCart   = errors.NewErrorResponse("Failed to create cart", http.StatusInternalServerError)

	ErrFailedDeleteCart     = errors.NewErrorResponse("Failed to permanently delete cart", http.StatusInternalServerError)
	ErrFailedDeleteAllCarts = errors.NewErrorResponse("Failed to permanently delete all carts", http.StatusInternalServerError)
)
