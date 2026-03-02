package merchantpolicy_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrFailedFindAllMerchantPolicies            = errors.NewErrorResponse("Failed to fetch all merchant policies", http.StatusInternalServerError)
	ErrFailedFindActiveMerchantPolicies         = errors.NewErrorResponse("Failed to fetch active merchant policies", http.StatusInternalServerError)
	ErrFailedFindTrashedMerchantPolicies        = errors.NewErrorResponse("Failed to fetch trashed merchant policies", http.StatusInternalServerError)
	ErrFailedFindMerchantPolicyById             = errors.NewErrorResponse("Failed to find merchant policy by ID", http.StatusInternalServerError)
	ErrFailedCreateMerchantPolicy               = errors.NewErrorResponse("Failed to create merchant policy", http.StatusInternalServerError)
	ErrFailedUpdateMerchantPolicy               = errors.NewErrorResponse("Failed to update merchant policy", http.StatusInternalServerError)
	ErrFailedTrashedMerchantPolicy              = errors.NewErrorResponse("Failed to trash merchant policy", http.StatusInternalServerError)
	ErrFailedRestoreMerchantPolicy              = errors.NewErrorResponse("Failed to restore merchant policy", http.StatusInternalServerError)
	ErrFailedDeleteMerchantPolicyPermanent      = errors.NewErrorResponse("Failed to permanently delete merchant policy", http.StatusInternalServerError)
	ErrFailedRestoreAllMerchantPolicies         = errors.NewErrorResponse("Failed to restore all merchant policies", http.StatusInternalServerError)
	ErrFailedDeleteAllMerchantPoliciesPermanent = errors.NewErrorResponse("Failed to permanently delete all merchant policies", http.StatusInternalServerError)
)
