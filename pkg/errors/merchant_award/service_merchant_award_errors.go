package merchantaward_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrFailedFindAllMerchantAwards            = response.NewErrorResponse("Failed to find all merchant awards", http.StatusInternalServerError)
	ErrFailedFindActiveMerchantAwards         = response.NewErrorResponse("Failed to find active merchant awards", http.StatusInternalServerError)
	ErrFailedFindTrashedMerchantAwards        = response.NewErrorResponse("Failed to find trashed merchant awards", http.StatusInternalServerError)
	ErrFailedFindMerchantAwardById            = response.NewErrorResponse("Failed to find merchant award by ID", http.StatusInternalServerError)
	ErrFailedCreateMerchantAward              = response.NewErrorResponse("Failed to create merchant award", http.StatusInternalServerError)
	ErrFailedUpdateMerchantAward              = response.NewErrorResponse("Failed to update merchant award", http.StatusInternalServerError)
	ErrFailedTrashedMerchantAward             = response.NewErrorResponse("Failed to trash merchant award", http.StatusInternalServerError)
	ErrFailedRestoreMerchantAward             = response.NewErrorResponse("Failed to restore merchant award", http.StatusInternalServerError)
	ErrFailedDeleteMerchantAwardPermanent     = response.NewErrorResponse("Failed to permanently delete merchant award", http.StatusInternalServerError)
	ErrFailedRestoreAllMerchantAwards         = response.NewErrorResponse("Failed to restore all merchant awards", http.StatusInternalServerError)
	ErrFailedDeleteAllMerchantAwardsPermanent = response.NewErrorResponse("Failed to permanently delete all merchant awards", http.StatusInternalServerError)
)
