package merchantaward_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiFailedInvalidMerchantAwardId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid merchant award ID", http.StatusBadRequest)
	}

	ErrApiFailedFindAllMerchantAward = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch all merchant awards", http.StatusInternalServerError)
	}

	ErrApiFailedFindMerchantAwardById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find merchant award by ID", http.StatusInternalServerError)
	}

	ErrApiFailedFindActiveMerchantAward = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch active merchant awards", http.StatusInternalServerError)
	}

	ErrApiFailedFindTrashedMerchantAward = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch trashed merchant awards", http.StatusInternalServerError)
	}

	ErrApiFailedCreateMerchantAward = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create merchant award", http.StatusInternalServerError)
	}

	ErrApiFailedUpdateMerchantAward = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update merchant award", http.StatusInternalServerError)
	}

	ErrApiBindCreateMerchantAward = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to bind merchant award", http.StatusInternalServerError)
	}

	ErrApiBindUpdateMerchantAward = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to bind merchant award", http.StatusInternalServerError)
	}

	ErrApiValidateCreateMerchantAward = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to validate merchant award", http.StatusInternalServerError)
	}

	ErrApiValidateUpdateMerchantAward = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to validate merchant award", http.StatusInternalServerError)
	}

	ErrApiFailedTrashedMerchantAward = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to trash merchant award", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreMerchantAward = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore merchant award", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteMerchantAwardPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete merchant award", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllMerchantAward = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all merchant awards", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteAllMerchantAwardPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete all merchant awards", http.StatusInternalServerError)
	}
)
