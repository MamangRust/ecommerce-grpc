package merchantpolicy_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrFailedFindAllMerchantPolicies            = response.NewErrorResponse("Failed to fetch all merchant policies", http.StatusInternalServerError)
	ErrFailedFindActiveMerchantPolicies         = response.NewErrorResponse("Failed to fetch active merchant policies", http.StatusInternalServerError)
	ErrFailedFindTrashedMerchantPolicies        = response.NewErrorResponse("Failed to fetch trashed merchant policies", http.StatusInternalServerError)
	ErrFailedFindMerchantPolicyById             = response.NewErrorResponse("Failed to find merchant policy by ID", http.StatusInternalServerError)
	ErrFailedCreateMerchantPolicy               = response.NewErrorResponse("Failed to create merchant policy", http.StatusInternalServerError)
	ErrFailedUpdateMerchantPolicy               = response.NewErrorResponse("Failed to update merchant policy", http.StatusInternalServerError)
	ErrFailedTrashedMerchantPolicy              = response.NewErrorResponse("Failed to trash merchant policy", http.StatusInternalServerError)
	ErrFailedRestoreMerchantPolicy              = response.NewErrorResponse("Failed to restore merchant policy", http.StatusInternalServerError)
	ErrFailedDeleteMerchantPolicyPermanent      = response.NewErrorResponse("Failed to permanently delete merchant policy", http.StatusInternalServerError)
	ErrFailedRestoreAllMerchantPolicies         = response.NewErrorResponse("Failed to restore all merchant policies", http.StatusInternalServerError)
	ErrFailedDeleteAllMerchantPoliciesPermanent = response.NewErrorResponse("Failed to permanently delete all merchant policies", http.StatusInternalServerError)
)
