package merchantaward_errors

import "errors"

var (
	ErrFindAllMerchantAwards            = errors.New("failed to find all merchant awards")
	ErrFindByActiveMerchantAwards       = errors.New("failed to find active merchant awards")
	ErrFindByTrashedMerchantAwards      = errors.New("failed to find trashed merchant awards")
	ErrFindByIdMerchantAward            = errors.New("failed to find merchant award by ID")
	ErrCreateMerchantAward              = errors.New("failed to create merchant award")
	ErrUpdateMerchantAward              = errors.New("failed to update merchant award")
	ErrTrashedMerchantAward             = errors.New("failed to trash merchant award")
	ErrRestoreMerchantAward             = errors.New("failed to restore merchant award")
	ErrDeleteMerchantAwardPermanent     = errors.New("failed to permanently delete merchant award")
	ErrRestoreAllMerchantAwards         = errors.New("failed to restore all merchant awards")
	ErrDeleteAllMerchantAwardsPermanent = errors.New("failed to permanently delete all merchant awards")
)
