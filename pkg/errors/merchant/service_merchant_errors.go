package merchant_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrFailedFindAllMerchants            = errors.NewErrorResponse("Failed to find all merchants", http.StatusInternalServerError)
	ErrFailedFindActiveMerchants         = errors.NewErrorResponse("Failed to find active merchants", http.StatusInternalServerError)
	ErrFailedFindTrashedMerchants        = errors.NewErrorResponse("Failed to find trashed merchants", http.StatusInternalServerError)
	ErrFailedFindMerchantById            = errors.NewErrorResponse("Failed to find merchant by ID", http.StatusInternalServerError)
	ErrFailedCreateMerchant              = errors.NewErrorResponse("Failed to create merchant", http.StatusInternalServerError)
	ErrFailedUpdateMerchant              = errors.NewErrorResponse("Failed to update merchant", http.StatusInternalServerError)
	ErrFailedTrashedMerchant             = errors.NewErrorResponse("Failed to trash merchant", http.StatusInternalServerError)
	ErrFailedRestoreMerchant             = errors.NewErrorResponse("Failed to restore merchant", http.StatusInternalServerError)
	ErrFailedDeleteMerchantPermanent     = errors.NewErrorResponse("Failed to permanently delete merchant", http.StatusInternalServerError)
	ErrFailedRestoreAllMerchants         = errors.NewErrorResponse("Failed to restore all merchants", http.StatusInternalServerError)
	ErrFailedDeleteAllMerchantsPermanent = errors.NewErrorResponse("Failed to permanently delete all merchants", http.StatusInternalServerError)
)
