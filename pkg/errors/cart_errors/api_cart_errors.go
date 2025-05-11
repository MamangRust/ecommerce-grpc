package cart_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiFailedInvalidUserId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid user ID", http.StatusBadRequest)
	}

	ErrApiInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid cart ID", http.StatusBadRequest)
	}

	ErrApiFailedFindAllCarts = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch Carts", http.StatusInternalServerError)
	}

	ErrApiFailedCreateCart = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create Cart", http.StatusInternalServerError)
	}

	ErrApiValidateCreateCart = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create Cart request", http.StatusBadRequest)
	}

	ErrApiBindCreateCart = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create Cart request", http.StatusBadRequest)
	}

	ErrApiFailedDeleteCart = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to delete Cart", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteAllCarts = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to delete all Carts", http.StatusInternalServerError)
	}

	ErrApiBindDeleteAllCart = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid delete all Cart request", http.StatusBadRequest)
	}

	ErrApiValidateDeleteAllCart = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid delete all Cart request", http.StatusBadRequest)
	}
)
