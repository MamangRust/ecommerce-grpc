package merchantsociallink_errors

import "errors"

var (
	ErrCreateMerchantSocialLink              = errors.New("failed to create merchant social link")
	ErrUpdateMerchantSocialLink              = errors.New("failed to update merchant social link")
	ErrTrashMerchantSocialLink               = errors.New("failed to trash merchant social link")
	ErrRestoreMerchantSocialLink             = errors.New("failed to restore merchant social link")
	ErrDeletePermanentMerchantSocialLink     = errors.New("failed to permanently delete merchant social link")
	ErrRestoreAllMerchantSocialLinks         = errors.New("failed to restore all merchant social links")
	ErrDeleteAllPermanentMerchantSocialLinks = errors.New("failed to permanently delete all merchant social links")
)
