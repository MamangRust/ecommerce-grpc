package user_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiUserNotFound = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "User not found", http.StatusNotFound)
	}

	ErrApiUserInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid User id", http.StatusNotFound)
	}

	ErrApiFailedFindAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch Users", http.StatusInternalServerError)
	}

	ErrApiFailedFindActive = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch active Users", http.StatusInternalServerError)
	}

	ErrApiFailedFindTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch trashed Users", http.StatusInternalServerError)
	}
	ErrInvalidUserId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid User id", http.StatusBadRequest)
	}

	ErrApiFailedCreateUser = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create User", http.StatusInternalServerError)
	}

	ErrApiFailedUpdateUser = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update User", http.StatusInternalServerError)
	}

	ErrApiValidateCreateUser = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create User request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateUser = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update User request", http.StatusBadRequest)
	}

	ErrApiBindCreateUser = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create User request", http.StatusBadRequest)
	}

	ErrApiBindUpdateUser = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update User request", http.StatusBadRequest)
	}

	ErrApiFailedTrashedUser = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to move User to trash", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreUser = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore User", http.StatusInternalServerError)
	}

	ErrApiFailedDeletePermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to delete User permanently", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all Users", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to delete all Users permanently", http.StatusInternalServerError)
	}
)
