package merchantbusiness_errors

import "errors"

var (
	ErrFindAllMerchantBusinesses            = errors.New("failed to find all merchant businesses")
	ErrFindActiveMerchantBusinesses         = errors.New("failed to find active merchant businesses")
	ErrFindTrashedMerchantBusinesses        = errors.New("failed to find trashed merchant businesses")
	ErrMerchantBusinessNotFound             = errors.New("merchant business not found")
	ErrCreateMerchantBusiness               = errors.New("failed to create merchant business")
	ErrUpdateMerchantBusiness               = errors.New("failed to update merchant business")
	ErrTrashMerchantBusiness                = errors.New("failed to trash merchant business")
	ErrRestoreMerchantBusiness              = errors.New("failed to restore merchant business")
	ErrDeletePermanentMerchantBusiness      = errors.New("failed to permanently delete merchant business")
	ErrRestoreAllMerchantBusinesses         = errors.New("failed to restore all merchant businesses")
	ErrDeleteAllPermanentMerchantBusinesses = errors.New("failed to permanently delete all merchant businesses")
)
