package merchantdetail_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrFailedImageNotFound             = errors.NewErrorResponse("Image not found", http.StatusNotFound)
	ErrFailedRemoveImageMerchantDetail = errors.NewErrorResponse("Failed to remove image merchant detail", http.StatusInternalServerError)
	ErrFailedLogoNotFound              = errors.NewErrorResponse("Failed to upload logo merchant detail", http.StatusInternalServerError)
	ErrFailedRemoveLogoMerchantDetail  = errors.NewErrorResponse("Failed to remove logo merchant detail", http.StatusInternalServerError)

	ErrFailedFindAllMerchantDetail            = errors.NewErrorResponse("Failed to find all merchant details", http.StatusInternalServerError)
	ErrFailedFindActiveMerchantDetail         = errors.NewErrorResponse("Failed to find active merchant details", http.StatusInternalServerError)
	ErrFailedFindTrashedMerchantDetail        = errors.NewErrorResponse("Failed to find trashed merchant details", http.StatusInternalServerError)
	ErrFailedFindMerchantDetailById           = errors.NewErrorResponse("Failed to find merchant detail by ID", http.StatusInternalServerError)
	ErrFailedCreateMerchantDetail             = errors.NewErrorResponse("Failed to create merchant detail", http.StatusInternalServerError)
	ErrFailedUpdateMerchantDetail             = errors.NewErrorResponse("Failed to update merchant detail", http.StatusInternalServerError)
	ErrFailedTrashedMerchantDetail            = errors.NewErrorResponse("Failed to trash merchant detail", http.StatusInternalServerError)
	ErrFailedRestoreMerchantDetail            = errors.NewErrorResponse("Failed to restore merchant detail", http.StatusInternalServerError)
	ErrFailedDeleteMerchantDetailPermanent    = errors.NewErrorResponse("Failed to permanently delete merchant detail", http.StatusInternalServerError)
	ErrFailedRestoreAllMerchantDetail         = errors.NewErrorResponse("Failed to restore all merchant details", http.StatusInternalServerError)
	ErrFailedDeleteAllMerchantDetailPermanent = errors.NewErrorResponse("Failed to permanently delete all merchant details", http.StatusInternalServerError)
)
