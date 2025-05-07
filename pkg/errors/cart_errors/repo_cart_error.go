package cart_errors

import "errors"

var (
	ErrCartNotFound = errors.New("cart not found")
	ErrFindAllCarts = errors.New("failed to find all carts")
	ErrCartConflict = errors.New("cart already exists")

	ErrCreateCart = errors.New("failed to create cart")

	ErrDeleteCartPermanent = errors.New("failed to permanently delete cart")
	ErrDeleteAllCarts      = errors.New("failed to permanently delete all carts")
)
