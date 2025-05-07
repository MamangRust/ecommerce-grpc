package refreshtoken_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrRefreshTokenNotFound = response.NewErrorResponse("Refresh token not found", http.StatusNotFound)
	ErrFailedExpire         = response.NewErrorResponse("Failed to find refresh token by token", http.StatusInternalServerError)
	ErrFailedFindByToken    = response.NewErrorResponse("Failed to find refresh token by token", http.StatusInternalServerError)
	ErrFailedFindByUserID   = response.NewErrorResponse("Failed to find refresh token by user ID", http.StatusInternalServerError)
	ErrFailedInValidToken   = response.NewErrorResponse("Failed to invalid access token", http.StatusInternalServerError)
	ErrFailedInValidUserId  = response.NewErrorResponse("Failed to invalid user id", http.StatusInternalServerError)

	ErrFailedCreateAccess  = response.NewErrorResponse("Failed to create access token", http.StatusInternalServerError)
	ErrFailedCreateRefresh = response.NewErrorResponse("Failed to create refresh token", http.StatusInternalServerError)

	ErrFailedCreateRefreshToken  = response.NewErrorResponse("Failed to create refresh token", http.StatusInternalServerError)
	ErrFailedUpdateRefreshToken  = response.NewErrorResponse("Failed to update refresh token", http.StatusInternalServerError)
	ErrFailedDeleteRefreshToken  = response.NewErrorResponse("Failed to delete refresh token", http.StatusInternalServerError)
	ErrFailedDeleteByUserID      = response.NewErrorResponse("Failed to delete refresh token by user ID", http.StatusInternalServerError)
	ErrFailedParseExpirationDate = response.NewErrorResponse("Failed to parse expiration date", http.StatusBadRequest)
)
