package refreshtoken_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrRefreshTokenNotFound = errors.NewErrorResponse("Refresh token not found", http.StatusNotFound)
	ErrFailedExpire         = errors.NewErrorResponse("Failed to find refresh token by token", http.StatusInternalServerError)
	ErrFailedFindByToken    = errors.NewErrorResponse("Failed to find refresh token by token", http.StatusInternalServerError)
	ErrFailedFindByUserID   = errors.NewErrorResponse("Failed to find refresh token by user ID", http.StatusInternalServerError)
	ErrFailedInValidToken   = errors.NewErrorResponse("Failed to invalid access token", http.StatusInternalServerError)
	ErrFailedInValidUserId  = errors.NewErrorResponse("Failed to invalid user id", http.StatusInternalServerError)

	ErrFailedCreateAccess  = errors.NewErrorResponse("Failed to create access token", http.StatusInternalServerError)
	ErrFailedCreateRefresh = errors.NewErrorResponse("Failed to create refresh token", http.StatusInternalServerError)

	ErrFailedCreateRefreshToken  = errors.NewErrorResponse("Failed to create refresh token", http.StatusInternalServerError)
	ErrFailedUpdateRefreshToken  = errors.NewErrorResponse("Failed to update refresh token", http.StatusInternalServerError)
	ErrFailedDeleteRefreshToken  = errors.NewErrorResponse("Failed to delete refresh token", http.StatusInternalServerError)
	ErrFailedDeleteByUserID      = errors.NewErrorResponse("Failed to delete refresh token by user ID", http.StatusInternalServerError)
	ErrFailedParseExpirationDate = errors.NewErrorResponse("Failed to parse expiration date", http.StatusBadRequest)
)
