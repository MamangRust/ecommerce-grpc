package merchantdetail_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrFailedImageNotFound             = response.NewErrorResponse("Image not found", http.StatusNotFound)
	ErrFailedRemoveImageMerchantDetail = response.NewErrorResponse("Failed to remove image merchant detail", http.StatusInternalServerError)
	ErrFailedLogoNotFound              = response.NewErrorResponse("Failed to upload logo merchant detail", http.StatusInternalServerError)
	ErrFailedRemoveLogoMerchantDetail  = response.NewErrorResponse("Failed to remove logo merchant detail", http.StatusInternalServerError)

	ErrFailedFindAllMerchantDetail            = response.NewErrorResponse("Failed to find all merchant details", http.StatusInternalServerError)
	ErrFailedFindActiveMerchantDetail         = response.NewErrorResponse("Failed to find active merchant details", http.StatusInternalServerError)
	ErrFailedFindTrashedMerchantDetail        = response.NewErrorResponse("Failed to find trashed merchant details", http.StatusInternalServerError)
	ErrFailedFindMerchantDetailById           = response.NewErrorResponse("Failed to find merchant detail by ID", http.StatusInternalServerError)
	ErrFailedCreateMerchantDetail             = response.NewErrorResponse("Failed to create merchant detail", http.StatusInternalServerError)
	ErrFailedUpdateMerchantDetail             = response.NewErrorResponse("Failed to update merchant detail", http.StatusInternalServerError)
	ErrFailedTrashedMerchantDetail            = response.NewErrorResponse("Failed to trash merchant detail", http.StatusInternalServerError)
	ErrFailedRestoreMerchantDetail            = response.NewErrorResponse("Failed to restore merchant detail", http.StatusInternalServerError)
	ErrFailedDeleteMerchantDetailPermanent    = response.NewErrorResponse("Failed to permanently delete merchant detail", http.StatusInternalServerError)
	ErrFailedRestoreAllMerchantDetail         = response.NewErrorResponse("Failed to restore all merchant details", http.StatusInternalServerError)
	ErrFailedDeleteAllMerchantDetailPermanent = response.NewErrorResponse("Failed to permanently delete all merchant details", http.StatusInternalServerError)
)
