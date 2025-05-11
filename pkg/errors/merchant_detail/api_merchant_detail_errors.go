package merchantdetail_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiInvalidMerchantId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "invalid_merchant_id", "Please provide a valid merchant ID", http.StatusBadRequest)
	}

	ErrApiDisplayNameRequired = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "validation_error", "Display name is required", http.StatusBadRequest)
	}

	ErrApiShortDescriptionRequired = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "validation_error", "Short description is required", http.StatusBadRequest)
	}

	ErrApiCoverImageRequired = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "cover_image_required", "Cover image is required", http.StatusBadRequest)
	}

	ErrApiLogoRequired = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "logo_required", "Logo is required", http.StatusBadRequest)
	}

	ErrApiSocialLinksRequired = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "social_links_required", "Social links are required", http.StatusBadRequest)
	}

	ErrApiInvalidSocialLinks = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "invalid_social_links", "Invalid format for social links", http.StatusBadRequest)
	}

	ErrApiInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid merchant detail ID", http.StatusBadRequest)
	}
	ErrApiInvalidBody = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid request body", http.StatusBadRequest)
	}

	ErrApiFailedFindAllMerchantDetail = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch merchant details", http.StatusInternalServerError)
	}
	ErrApiFailedFindActiveMerchantDetail = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch active merchant details", http.StatusInternalServerError)
	}
	ErrApiFailedFindTrashedMerchantDetail = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch trashed merchant details", http.StatusInternalServerError)
	}
	ErrApiFailedFindMerchantDetailById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find merchant detail by ID", http.StatusInternalServerError)
	}
	ErrApiFailedCreateMerchantDetail = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create merchant detail", http.StatusInternalServerError)
	}
	ErrApiFailedUpdateMerchantDetail = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update merchant detail", http.StatusInternalServerError)
	}
	ErrApiFailedTrashedMerchantDetail = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to trash merchant detail", http.StatusInternalServerError)
	}
	ErrApiFailedRestoreMerchantDetail = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore merchant detail", http.StatusInternalServerError)
	}
	ErrApiFailedDeleteMerchantDetailPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete merchant detail", http.StatusInternalServerError)
	}
	ErrApiFailedRestoreAllMerchantDetail = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all merchant details", http.StatusInternalServerError)
	}
	ErrApiFailedDeleteAllMerchantDetailPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete all merchant details", http.StatusInternalServerError)
	}
)
