package merchantpolicy_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid merchant policy ID", http.StatusBadRequest)
	}

	ErrApiFailedFindAllMerchantPolicy = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch merchant policies", http.StatusInternalServerError)
	}

	ErrApiFailedFindByIdMerchantPolicy = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find merchant policy by ID", http.StatusInternalServerError)
	}

	ErrApiFailedFindByActiveMerchantPolicy = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch active merchant policies", http.StatusInternalServerError)
	}

	ErrApiFailedFindByTrashedMerchantPolicy = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch trashed merchant policies", http.StatusInternalServerError)
	}

	ErrApiFailedCreateMerchantPolicy = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create merchant policy", http.StatusInternalServerError)
	}

	ErrApiFailedUpdateMerchantPolicy = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update merchant policy", http.StatusInternalServerError)
	}

	ErrApiBindCreateMerchantPolicy = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to bind merchant policy", http.StatusInternalServerError)
	}

	ErrApiBindUpdateMerchantPolicy = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to bind merchant policy", http.StatusInternalServerError)
	}

	ErrValidateCreateMerchantPolicy = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to validate merchant policy", http.StatusInternalServerError)
	}

	ErrValidateUpdateMerchantPolicy = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to validate merchant policy", http.StatusInternalServerError)
	}

	ErrApiFailedTrashMerchantPolicy = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to trash merchant policy", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreMerchantPolicy = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore merchant policy", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteMerchantPolicy = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete merchant policy", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllMerchantPolicies = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all merchant policies", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteAllMerchantPolicies = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete all merchant policies", http.StatusInternalServerError)
	}
)
