package merchantaward_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrFailedFindAllMerchantAwards            = errors.NewErrorResponse("Failed to find all merchant awards", http.StatusInternalServerError)
	ErrFailedFindActiveMerchantAwards         = errors.NewErrorResponse("Failed to find active merchant awards", http.StatusInternalServerError)
	ErrFailedFindTrashedMerchantAwards        = errors.NewErrorResponse("Failed to find trashed merchant awards", http.StatusInternalServerError)
	ErrFailedFindMerchantAwardById            = errors.NewErrorResponse("Failed to find merchant award by ID", http.StatusInternalServerError)
	ErrFailedCreateMerchantAward              = errors.NewErrorResponse("Failed to create merchant award", http.StatusInternalServerError)
	ErrFailedUpdateMerchantAward              = errors.NewErrorResponse("Failed to update merchant award", http.StatusInternalServerError)
	ErrFailedTrashedMerchantAward             = errors.NewErrorResponse("Failed to trash merchant award", http.StatusInternalServerError)
	ErrFailedRestoreMerchantAward             = errors.NewErrorResponse("Failed to restore merchant award", http.StatusInternalServerError)
	ErrFailedDeleteMerchantAwardPermanent     = errors.NewErrorResponse("Failed to permanently delete merchant award", http.StatusInternalServerError)
	ErrFailedRestoreAllMerchantAwards         = errors.NewErrorResponse("Failed to restore all merchant awards", http.StatusInternalServerError)
	ErrFailedDeleteAllMerchantAwardsPermanent = errors.NewErrorResponse("Failed to permanently delete all merchant awards", http.StatusInternalServerError)
)
