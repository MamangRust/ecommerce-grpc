package user_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrUserIDInValid     = response.NewErrorResponse("Invalid user ID", http.StatusBadRequest)
	ErrUserNotFoundRes   = response.NewErrorResponse("User not found", http.StatusNotFound)
	ErrUserEmailAlready  = response.NewErrorResponse("User email already exists", http.StatusBadRequest)
	ErrUserPassword      = response.NewErrorResponse("Failed invalid password", http.StatusBadRequest)
	ErrFailedFindAll     = response.NewErrorResponse("Failed to fetch users", http.StatusInternalServerError)
	ErrFailedFindActive  = response.NewErrorResponse("Failed to fetch active users", http.StatusInternalServerError)
	ErrFailedFindTrashed = response.NewErrorResponse("Failed to fetch trashed users", http.StatusInternalServerError)

	ErrFailedCreateUser = response.NewErrorResponse("Failed to create user", http.StatusInternalServerError)
	ErrFailedUpdateUser = response.NewErrorResponse("Failed to update user", http.StatusInternalServerError)

	ErrFailedTrashedUser     = response.NewErrorResponse("Failed to move user to trash", http.StatusInternalServerError)
	ErrFailedRestoreUser     = response.NewErrorResponse("Failed to restore user", http.StatusInternalServerError)
	ErrFailedDeletePermanent = response.NewErrorResponse("Failed to delete user permanently", http.StatusInternalServerError)

	ErrFailedRestoreAll = response.NewErrorResponse("Failed to restore all users", http.StatusInternalServerError)
	ErrFailedDeleteAll  = response.NewErrorResponse("Failed to delete all users permanently", http.StatusInternalServerError)
)
