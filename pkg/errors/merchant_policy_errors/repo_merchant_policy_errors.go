package merchantpolicy_errors

import "errors"

var (
	ErrFindAllMerchantPolicy            = errors.New("failed to find all merchant policies")
	ErrFindByActiveMerchantPolicy       = errors.New("failed to find active merchant policies")
	ErrFindByTrashedMerchantPolicy      = errors.New("failed to find trashed merchant policies")
	ErrFindByIdMerchantPolicy           = errors.New("failed to find merchant policy by ID")
	ErrCreateMerchantPolicy             = errors.New("failed to create merchant policy")
	ErrUpdateMerchantPolicy             = errors.New("failed to update merchant policy")
	ErrTrashedMerchantPolicy            = errors.New("failed to trash merchant policy")
	ErrRestoreMerchantPolicy            = errors.New("failed to restore merchant policy")
	ErrDeleteMerchantPolicyPermanent    = errors.New("failed to permanently delete merchant policy")
	ErrRestoreAllMerchantPolicy         = errors.New("failed to restore all merchant policies")
	ErrDeleteAllMerchantPolicyPermanent = errors.New("failed to permanently delete all merchant policies")
)
