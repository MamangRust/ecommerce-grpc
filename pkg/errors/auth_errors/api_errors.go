package auth_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiLogin = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "login failed: invalid argument provided", http.StatusBadRequest)
	}

	ErrApiRefreshToken = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "refresh-token failed: invalid access token", http.StatusBadRequest)
	}

	ErrApiGetMe = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "get user info failed: unauthenticated", http.StatusUnauthorized)
	}

	ErrValidateLogin = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid login request", http.StatusBadRequest)
	}

	ErrValidateRegister = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid register request", http.StatusBadRequest)
	}

	ErrValidateRefreshToken = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid refresh-token request", http.StatusBadRequest)
	}

	ErrBindLogin = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "binding failed: invalid login request payload", http.StatusBadRequest)
	}

	ErrBindRefreshToken = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "binding failed: invalid refresh token request payload", http.StatusBadRequest)
	}

	ErrBindRegister = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "binding failed: invalid register request payload", http.StatusBadRequest)
	}

	ErrInvalidLogin = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid email or password", http.StatusUnauthorized)
	}

	ErrInvalidAccessToken = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid access token", http.StatusInternalServerError)
	}

	ErrApiRegister = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "register failed: invalid argument", http.StatusBadRequest)
	}
)
