package merchantbusiness_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiFailedInvalidIdMerchantBusiness = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid merchant business ID", http.StatusBadRequest)
	}

	ErrApiFailedFindAllMerchantBusiness = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch merchant businesses", http.StatusInternalServerError)
	}

	ErrApiFailedFindByIdMerchantBusiness = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find merchant business by ID", http.StatusInternalServerError)
	}

	ErrApiFailedFindByActiveMerchantBusiness = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch active merchant businesses", http.StatusInternalServerError)
	}

	ErrApiFailedFindByTrashedMerchantBusiness = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch trashed merchant businesses", http.StatusInternalServerError)
	}

	ErrApiFailedCreateMerchantBusiness = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create merchant business", http.StatusInternalServerError)
	}

	ErrApiFailedUpdateMerchantBusiness = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update merchant business", http.StatusInternalServerError)
	}

	ErrApiBindCreateMerchantBusiness = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to bind merchant business", http.StatusInternalServerError)
	}

	ErrApiBindUpdateMerchantBusiness = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to bind merchant business", http.StatusInternalServerError)
	}

	ErrApiFailedTrashMerchantBusiness = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to trash merchant business", http.StatusInternalServerError)
	}

	ErrApiValidateCreateMerchantBusiness = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to validate merchant business", http.StatusInternalServerError)
	}

	ErrApiValidateUpdateMerchantBusiness = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to validate merchant business", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreMerchantBusiness = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore merchant business", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteMerchantBusinessPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete merchant business", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllMerchantBusiness = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all merchant businesses", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteAllMerchantBusinessPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete all merchant businesses", http.StatusInternalServerError)
	}

	ErrApiInvalidMerchantBusinessData = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid merchant business data", http.StatusBadRequest)
	}

	ErrApiMerchantBusinessNotFound = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "merchant business not found", http.StatusNotFound)
	}
)
