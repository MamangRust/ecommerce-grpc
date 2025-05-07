package banner_errors

import "errors"

var (
	ErrBannerStartDate = errors.New("banner start date must be less than end date")
	ErrBannerEndDate   = errors.New("banner end date must be greater than start date")
	ErrBannerStartTime = errors.New("banner start time must be less than end time")
	ErrBannerEndTime   = errors.New("banner end time must be greater than start time")

	ErrBannerNotFound     = errors.New("banner not found")
	ErrFindAllBanners     = errors.New("failed to find all banners")
	ErrFindActiveBanners  = errors.New("failed to find active banners")
	ErrFindTrashedBanners = errors.New("failed to find trashed banners")
	ErrBannerConflict     = errors.New("failed banner already exists")

	ErrCreateBanner = errors.New("failed to create banner")
	ErrUpdateBanner = errors.New("failed to update banner")

	ErrTrashedBanner         = errors.New("failed to move banner to trash")
	ErrRestoreBanner         = errors.New("failed to restore banner from trash")
	ErrDeleteBannerPermanent = errors.New("failed to permanently delete banner")

	ErrRestoreAllBanners = errors.New("failed to restore all banners")
	ErrDeleteAllBanners  = errors.New("failed to permanently delete all banners")
)
