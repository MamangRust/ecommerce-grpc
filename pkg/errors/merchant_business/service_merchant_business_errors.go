package merchantbusiness_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrFailedFindAllMerchantBusiness            = response.NewErrorResponse("Failed to fetch all merchant businesses", http.StatusInternalServerError)
	ErrFailedFindActiveMerchantBusiness         = response.NewErrorResponse("Failed to fetch active merchant businesses", http.StatusInternalServerError)
	ErrFailedFindTrashedMerchantBusiness        = response.NewErrorResponse("Failed to fetch trashed merchant businesses", http.StatusInternalServerError)
	ErrFailedFindMerchantBusinessById           = response.NewErrorResponse("Failed to find merchant business by ID", http.StatusInternalServerError)
	ErrFailedCreateMerchantBusiness             = response.NewErrorResponse("Failed to create merchant business", http.StatusInternalServerError)
	ErrFailedUpdateMerchantBusiness             = response.NewErrorResponse("Failed to update merchant business", http.StatusInternalServerError)
	ErrFailedTrashedMerchantBusiness            = response.NewErrorResponse("Failed to trash merchant business", http.StatusInternalServerError)
	ErrFailedRestoreMerchantBusiness            = response.NewErrorResponse("Failed to restore merchant business", http.StatusInternalServerError)
	ErrFailedDeleteMerchantBusinessPermanent    = response.NewErrorResponse("Failed to permanently delete merchant business", http.StatusInternalServerError)
	ErrFailedRestoreAllMerchantBusiness         = response.NewErrorResponse("Failed to restore all merchant businesses", http.StatusInternalServerError)
	ErrFailedDeleteAllMerchantBusinessPermanent = response.NewErrorResponse("Failed to permanently delete all merchant businesses", http.StatusInternalServerError)
)
