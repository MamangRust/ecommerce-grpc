package user_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrUserIDInValid     = errors.NewErrorResponse("Invalid user ID", http.StatusBadRequest)
	ErrUserNotFoundRes   = errors.NewErrorResponse("User not found", http.StatusNotFound)
	ErrUserEmailAlready  = errors.NewErrorResponse("User email already exists", http.StatusBadRequest)
	ErrUserPassword      = errors.NewErrorResponse("Failed invalid password", http.StatusBadRequest)
	ErrFailedFindAll     = errors.NewErrorResponse("Failed to fetch users", http.StatusInternalServerError)
	ErrFailedFindActive  = errors.NewErrorResponse("Failed to fetch active users", http.StatusInternalServerError)
	ErrFailedFindTrashed = errors.NewErrorResponse("Failed to fetch trashed users", http.StatusInternalServerError)

	ErrFailedCreateUser = errors.NewErrorResponse("Failed to create user", http.StatusInternalServerError)
	ErrFailedUpdateUser = errors.NewErrorResponse("Failed to update user", http.StatusInternalServerError)

	ErrFailedTrashedUser     = errors.NewErrorResponse("Failed to move user to trash", http.StatusInternalServerError)
	ErrFailedRestoreUser     = errors.NewErrorResponse("Failed to restore user", http.StatusInternalServerError)
	ErrFailedDeletePermanent = errors.NewErrorResponse("Failed to delete user permanently", http.StatusInternalServerError)

	ErrFailedRestoreAll = errors.NewErrorResponse("Failed to restore all users", http.StatusInternalServerError)
	ErrFailedDeleteAll  = errors.NewErrorResponse("Failed to delete all users permanently", http.StatusInternalServerError)
)
