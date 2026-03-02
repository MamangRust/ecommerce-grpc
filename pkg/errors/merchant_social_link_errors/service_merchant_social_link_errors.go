package merchantsociallink_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrFailedCreateMerchantSocialLink              = errors.NewErrorResponse("failed to create merchant social link", http.StatusInternalServerError)
	ErrFailedUpdateMerchantSocialLink              = errors.NewErrorResponse("failed to update merchant social link", http.StatusInternalServerError)
	ErrFailedTrashMerchantSocialLink               = errors.NewErrorResponse("failed to trash merchant social link", http.StatusInternalServerError)
	ErrFailedRestoreMerchantSocialLink             = errors.NewErrorResponse("failed to restore merchant social link", http.StatusInternalServerError)
	ErrFailedDeletePermanentMerchantSocialLink     = errors.NewErrorResponse("failed to permanently delete merchant social link", http.StatusInternalServerError)
	ErrFailedRestoreAllMerchantSocialLinks         = errors.NewErrorResponse("failed to restore all merchant social links", http.StatusInternalServerError)
	ErrFailedDeleteAllPermanentMerchantSocialLinks = errors.NewErrorResponse("failed to permanently delete all merchant social links", http.StatusInternalServerError)
	ErrFailedFindMonthlyTotalPriceByMerchant       = errors.NewErrorResponse("Failed to find monthly total price by merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalPriceByMerchant        = errors.NewErrorResponse("Failed to find yearly total price by merchant", http.StatusInternalServerError)
)
