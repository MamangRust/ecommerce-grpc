package merchantsociallink_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrFailedCreateMerchantSocialLink              = response.NewErrorResponse("failed to create merchant social link", http.StatusInternalServerError)
	ErrFailedUpdateMerchantSocialLink              = response.NewErrorResponse("failed to update merchant social link", http.StatusInternalServerError)
	ErrFailedTrashMerchantSocialLink               = response.NewErrorResponse("failed to trash merchant social link", http.StatusInternalServerError)
	ErrFailedRestoreMerchantSocialLink             = response.NewErrorResponse("failed to restore merchant social link", http.StatusInternalServerError)
	ErrFailedDeletePermanentMerchantSocialLink     = response.NewErrorResponse("failed to permanently delete merchant social link", http.StatusInternalServerError)
	ErrFailedRestoreAllMerchantSocialLinks         = response.NewErrorResponse("failed to restore all merchant social links", http.StatusInternalServerError)
	ErrFailedDeleteAllPermanentMerchantSocialLinks = response.NewErrorResponse("failed to permanently delete all merchant social links", http.StatusInternalServerError)
	ErrFailedFindMonthlyTotalPriceByMerchant       = response.NewErrorResponse("Failed to find monthly total price by merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalPriceByMerchant        = response.NewErrorResponse("Failed to find yearly total price by merchant", http.StatusInternalServerError)
)
