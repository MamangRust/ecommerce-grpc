package banner_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrBannerNotFoundRes = response.NewErrorResponse("Banner not found", http.StatusNotFound)
	ErrBannerInvalidData = response.NewErrorResponse("Invalid banner data", http.StatusBadRequest)

	ErrFailedFindAllBanners     = response.NewErrorResponse("Failed to fetch banners", http.StatusInternalServerError)
	ErrFailedFindActiveBanners  = response.NewErrorResponse("Failed to fetch active banners", http.StatusInternalServerError)
	ErrFailedFindTrashedBanners = response.NewErrorResponse("Failed to fetch trashed banners", http.StatusInternalServerError)

	ErrFailedCreateBanner = response.NewErrorResponse("Failed to create banner", http.StatusInternalServerError)
	ErrFailedUpdateBanner = response.NewErrorResponse("Failed to update banner", http.StatusInternalServerError)

	ErrFailedTrashedBanner = response.NewErrorResponse("Failed to move banner to trash", http.StatusInternalServerError)
	ErrFailedRestoreBanner = response.NewErrorResponse("Failed to restore banner", http.StatusInternalServerError)
	ErrFailedDeleteBanner  = response.NewErrorResponse("Failed to permanently delete banner", http.StatusInternalServerError)

	ErrFailedRestoreAllBanners = response.NewErrorResponse("Failed to restore all banners", http.StatusInternalServerError)
	ErrFailedDeleteAllBanners  = response.NewErrorResponse("Failed to permanently delete all banners", http.StatusInternalServerError)
)
