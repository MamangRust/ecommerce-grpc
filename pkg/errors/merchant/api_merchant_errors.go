package merchant_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid merchant ID", http.StatusBadRequest)
	}

	ErrApiFailedFindAllMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch all merchants", http.StatusInternalServerError)
	}
	ErrApiFailedFindMerchantById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch merchant by ID", http.StatusInternalServerError)
	}
	ErrApiFailedFindActiveMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch active merchants", http.StatusInternalServerError)
	}
	ErrApiFailedFindTrashedMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch trashed merchants", http.StatusInternalServerError)
	}

	ErrApiFailedCreateMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create merchant", http.StatusInternalServerError)
	}
	ErrApiFailedUpdateMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update merchant", http.StatusInternalServerError)
	}

	ErrApiBindCreateMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to bind merchant", http.StatusInternalServerError)
	}
	ErrApiBindUpdateMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to bind merchant", http.StatusInternalServerError)
	}

	ErrApiValidateCreateMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to validate merchant", http.StatusInternalServerError)
	}
	ErrApiValidateUpdateMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to validate merchant", http.StatusInternalServerError)
	}

	ErrApiFailedTrashMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to trash merchant", http.StatusInternalServerError)
	}
	ErrApiFailedRestoreMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore merchant", http.StatusInternalServerError)
	}
	ErrApiFailedDeleteMerchantPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete merchant", http.StatusInternalServerError)
	}
	ErrApiFailedRestoreAllMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all merchants", http.StatusInternalServerError)
	}
	ErrApiFailedDeleteAllMerchantPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete all merchants", http.StatusInternalServerError)
	}
)
