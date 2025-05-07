package product_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrFailedCountStock = response.NewErrorResponse("Failed to count stock", http.StatusInternalServerError)

	ErrFailedDeletingNotFoundProduct = response.NewErrorResponse("Product not found", http.StatusNotFound)
	ErrFailedDeleteImageProduct      = response.NewErrorResponse("Failed to delete image product", http.StatusInternalServerError)

	ErrFailedFindAllProducts        = response.NewErrorResponse("Failed to find all products", http.StatusInternalServerError)
	ErrFailedFindProductsByMerchant = response.NewErrorResponse("Failed to find products by merchant", http.StatusInternalServerError)
	ErrFailedFindProductsByCategory = response.NewErrorResponse("Failed to find products by category", http.StatusInternalServerError)
	ErrFailedFindProductById        = response.NewErrorResponse("Failed to find product by ID", http.StatusInternalServerError)
	ErrFailedFindProductTrashedById = response.NewErrorResponse("Failed to find trashed product by ID", http.StatusInternalServerError)
	ErrFailedFindProductByTrashed   = response.NewErrorResponse("Failed to find product by trashed", http.StatusInternalServerError)

	ErrFailedFindProductsByActive  = response.NewErrorResponse("Failed to find active products", http.StatusInternalServerError)
	ErrFailedFindProductsByTrashed = response.NewErrorResponse("Failed to find trashed products", http.StatusInternalServerError)
	ErrFailedCreateProduct         = response.NewErrorResponse("Failed to create product", http.StatusInternalServerError)
	ErrFailedUpdateProduct         = response.NewErrorResponse("Failed to update product", http.StatusInternalServerError)

	ErrFailedTrashProduct               = response.NewErrorResponse("Failed to trash product", http.StatusInternalServerError)
	ErrFailedRestoreProduct             = response.NewErrorResponse("Failed to restore product", http.StatusInternalServerError)
	ErrFailedDeleteProductPermanent     = response.NewErrorResponse("Failed to permanently delete product", http.StatusInternalServerError)
	ErrFailedRestoreAllProducts         = response.NewErrorResponse("Failed to restore all products", http.StatusInternalServerError)
	ErrFailedDeleteAllProductsPermanent = response.NewErrorResponse("Failed to permanently delete all products", http.StatusInternalServerError)
)
