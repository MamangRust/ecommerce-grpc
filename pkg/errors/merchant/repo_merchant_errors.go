package merchant_errors

import "errors"

var (
	ErrFindAllMerchants           = errors.New("failed to find all merchants")
	ErrFindByActiveMerchant       = errors.New("failed to find active merchants")
	ErrFindByTrashedMerchant      = errors.New("failed to find trashed merchants")
	ErrFindByIdMerchant           = errors.New("failed to find merchant by ID")
	ErrCreateMerchant             = errors.New("failed to create merchant")
	ErrUpdateMerchant             = errors.New("failed to update merchant")
	ErrTrashedMerchant            = errors.New("failed to trash merchant")
	ErrRestoreMerchant            = errors.New("failed to restore merchant")
	ErrDeleteMerchantPermanent    = errors.New("failed to permanently delete merchant")
	ErrRestoreAllMerchant         = errors.New("failed to restore all merchants")
	ErrDeleteAllMerchantPermanent = errors.New("failed to permanently delete all merchants")
)
