package shippingaddress_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiInvalidIdShippingAddress = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid shipping address ID", http.StatusBadRequest)
	}

	ErrApiInvalidOrderIdShippingAddress = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid order ID", http.StatusBadRequest)
	}

	ErrApiFailedFindAllShippingAddresses = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch shipping addresses", http.StatusInternalServerError)
	}

	ErrApiFailedFindShippingAddressById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find shipping address by ID", http.StatusInternalServerError)
	}

	ErrApiFailedFindShippingAddressByOrder = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find shipping address by order ID", http.StatusInternalServerError)
	}

	ErrApiFailedFindActiveShippingAddresses = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch active shipping addresses", http.StatusInternalServerError)
	}

	ErrApiFailedFindTrashedShippingAddresses = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch trashed shipping addresses", http.StatusInternalServerError)
	}

	ErrApiFailedTrashShippingAddress = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to trash shipping address", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreShippingAddress = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore shipping address", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteShippingAddressPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete shipping address", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllShippingAddresses = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all shipping addresses", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteAllPermanentShippingAddresses = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete all shipping addresses", http.StatusInternalServerError)
	}
)
