package product_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrFailedCountStock = errors.NewErrorResponse("Failed to count stock", http.StatusInternalServerError)

	ErrFailedDeletingNotFoundProduct = errors.NewErrorResponse("Product not found", http.StatusNotFound)
	ErrFailedDeleteImageProduct      = errors.NewErrorResponse("Failed to delete image product", http.StatusInternalServerError)

	ErrFailedFindAllProducts        = errors.NewErrorResponse("Failed to find all products", http.StatusInternalServerError)
	ErrFailedFindProductsByMerchant = errors.NewErrorResponse("Failed to find products by merchant", http.StatusInternalServerError)
	ErrFailedFindProductsByCategory = errors.NewErrorResponse("Failed to find products by category", http.StatusInternalServerError)
	ErrFailedFindProductById        = errors.NewErrorResponse("Failed to find product by ID", http.StatusInternalServerError)
	ErrFailedFindProductTrashedById = errors.NewErrorResponse("Failed to find trashed product by ID", http.StatusInternalServerError)
	ErrFailedFindProductByTrashed   = errors.NewErrorResponse("Failed to find product by trashed", http.StatusInternalServerError)

	ErrFailedFindProductsByActive  = errors.NewErrorResponse("Failed to find active products", http.StatusInternalServerError)
	ErrFailedFindProductsByTrashed = errors.NewErrorResponse("Failed to find trashed products", http.StatusInternalServerError)
	ErrFailedCreateProduct         = errors.NewErrorResponse("Failed to create product", http.StatusInternalServerError)
	ErrFailedUpdateProduct         = errors.NewErrorResponse("Failed to update product", http.StatusInternalServerError)

	ErrFailedTrashProduct               = errors.NewErrorResponse("Failed to trash product", http.StatusInternalServerError)
	ErrFailedRestoreProduct             = errors.NewErrorResponse("Failed to restore product", http.StatusInternalServerError)
	ErrFailedDeleteProductPermanent     = errors.NewErrorResponse("Failed to permanently delete product", http.StatusInternalServerError)
	ErrFailedRestoreAllProducts         = errors.NewErrorResponse("Failed to restore all products", http.StatusInternalServerError)
	ErrFailedDeleteAllProductsPermanent = errors.NewErrorResponse("Failed to permanently delete all products", http.StatusInternalServerError)
)
