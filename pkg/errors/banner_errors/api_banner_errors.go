package banner_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiBannerInvalidUserId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid user ID", http.StatusBadRequest)
	}

	ErrApiBannerNotFound = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Banner not found", http.StatusNotFound)
	}

	ErrApiBannerInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid banner ID", http.StatusBadRequest)
	}

	ErrApiFailedFindAllBanner = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch banners", http.StatusInternalServerError)
	}

	ErrApiFailedFindById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch banner by ID", http.StatusInternalServerError)
	}

	ErrApiFailedFindByActive = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch active banners", http.StatusInternalServerError)
	}

	ErrApiFailedFindByTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch trashed banners", http.StatusInternalServerError)
	}

	ErrApiFailedCreateBanner = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to create banner", http.StatusInternalServerError)
	}

	ErrApiFailedUpdateBanner = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to update banner", http.StatusInternalServerError)
	}

	ErrApiValidateCreateBanner = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create banner request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateBanner = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update banner request", http.StatusBadRequest)
	}

	ErrApiBindCreateBanner = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create banner request", http.StatusBadRequest)
	}

	ErrApiBindUpdateBanner = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update banner request", http.StatusBadRequest)
	}

	ErrApiFailedTrashedBanner = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to trash banner", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreBanner = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore banner", http.StatusInternalServerError)
	}

	ErrApiFailedDeletePermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete banner", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllBanner = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore all banners", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteAllPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete all banners", http.StatusInternalServerError)
	}
)
