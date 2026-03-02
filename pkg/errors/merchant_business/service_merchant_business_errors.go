package merchantbusiness_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrFailedFindAllMerchantBusiness            = errors.NewErrorResponse("Failed to fetch all merchant businesses", http.StatusInternalServerError)
	ErrFailedFindActiveMerchantBusiness         = errors.NewErrorResponse("Failed to fetch active merchant businesses", http.StatusInternalServerError)
	ErrFailedFindTrashedMerchantBusiness        = errors.NewErrorResponse("Failed to fetch trashed merchant businesses", http.StatusInternalServerError)
	ErrFailedFindMerchantBusinessById           = errors.NewErrorResponse("Failed to find merchant business by ID", http.StatusInternalServerError)
	ErrFailedCreateMerchantBusiness             = errors.NewErrorResponse("Failed to create merchant business", http.StatusInternalServerError)
	ErrFailedUpdateMerchantBusiness             = errors.NewErrorResponse("Failed to update merchant business", http.StatusInternalServerError)
	ErrFailedTrashedMerchantBusiness            = errors.NewErrorResponse("Failed to trash merchant business", http.StatusInternalServerError)
	ErrFailedRestoreMerchantBusiness            = errors.NewErrorResponse("Failed to restore merchant business", http.StatusInternalServerError)
	ErrFailedDeleteMerchantBusinessPermanent    = errors.NewErrorResponse("Failed to permanently delete merchant business", http.StatusInternalServerError)
	ErrFailedRestoreAllMerchantBusiness         = errors.NewErrorResponse("Failed to restore all merchant businesses", http.StatusInternalServerError)
	ErrFailedDeleteAllMerchantBusinessPermanent = errors.NewErrorResponse("Failed to permanently delete all merchant businesses", http.StatusInternalServerError)
)
