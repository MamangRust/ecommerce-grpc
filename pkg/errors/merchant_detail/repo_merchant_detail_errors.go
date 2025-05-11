package merchantdetail_errors

import "errors"

var (
	ErrFindAllMerchantDetails            = errors.New("failed to find all merchant details")
	ErrFindByActiveMerchantDetails       = errors.New("failed to find active merchant details")
	ErrFindByTrashedMerchantDetails      = errors.New("failed to find trashed merchant details")
	ErrFindByIdMerchantDetail            = errors.New("failed to find merchant detail by ID")
	ErrFindByIdTrashedMerchantDetail     = errors.New("failed to find trashed merchant detail by ID")
	ErrCreateMerchantDetail              = errors.New("failed to create merchant detail")
	ErrUpdateMerchantDetail              = errors.New("failed to update merchant detail")
	ErrTrashedMerchantDetail             = errors.New("failed to trash merchant detail")
	ErrRestoreMerchantDetail             = errors.New("failed to restore merchant detail")
	ErrDeleteMerchantDetailPermanent     = errors.New("failed to permanently delete merchant detail")
	ErrRestoreAllMerchantDetails         = errors.New("failed to restore all merchant details")
	ErrDeleteAllMerchantDetailsPermanent = errors.New("failed to permanently delete all merchant details")
)
