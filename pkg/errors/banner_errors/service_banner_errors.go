package banner_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrBannerNotFoundRes = errors.NewErrorResponse("Banner not found", http.StatusNotFound)
	ErrBannerInvalidData = errors.NewErrorResponse("Invalid banner data", http.StatusBadRequest)

	ErrFailedFindAllBanners     = errors.NewErrorResponse("Failed to fetch banners", http.StatusInternalServerError)
	ErrFailedFindActiveBanners  = errors.NewErrorResponse("Failed to fetch active banners", http.StatusInternalServerError)
	ErrFailedFindTrashedBanners = errors.NewErrorResponse("Failed to fetch trashed banners", http.StatusInternalServerError)

	ErrFailedCreateBanner = errors.NewErrorResponse("Failed to create banner", http.StatusInternalServerError)
	ErrFailedUpdateBanner = errors.NewErrorResponse("Failed to update banner", http.StatusInternalServerError)

	ErrFailedTrashedBanner = errors.NewErrorResponse("Failed to move banner to trash", http.StatusInternalServerError)
	ErrFailedRestoreBanner = errors.NewErrorResponse("Failed to restore banner", http.StatusInternalServerError)
	ErrFailedDeleteBanner  = errors.NewErrorResponse("Failed to permanently delete banner", http.StatusInternalServerError)

	ErrFailedRestoreAllBanners = errors.NewErrorResponse("Failed to restore all banners", http.StatusInternalServerError)
	ErrFailedDeleteAllBanners  = errors.NewErrorResponse("Failed to permanently delete all banners", http.StatusInternalServerError)
)
