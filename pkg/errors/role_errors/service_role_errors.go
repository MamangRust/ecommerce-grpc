package role_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrRoleNotFoundRes   = errors.NewErrorResponse("Role not found", http.StatusNotFound)
	ErrFailedFindAll     = errors.NewErrorResponse("Failed to fetch Roles", http.StatusInternalServerError)
	ErrFailedFindActive  = errors.NewErrorResponse("Failed to fetch active Roles", http.StatusInternalServerError)
	ErrFailedFindTrashed = errors.NewErrorResponse("Failed to fetch trashed Roles", http.StatusInternalServerError)

	ErrFailedCreateRole = errors.NewErrorResponse("Failed to create Role", http.StatusInternalServerError)
	ErrFailedUpdateRole = errors.NewErrorResponse("Failed to update Role", http.StatusInternalServerError)

	ErrFailedTrashedRole     = errors.NewErrorResponse("Failed to move Role to trash", http.StatusInternalServerError)
	ErrFailedRestoreRole     = errors.NewErrorResponse("Failed to restore Role", http.StatusInternalServerError)
	ErrFailedDeletePermanent = errors.NewErrorResponse("Failed to delete Role permanently", http.StatusInternalServerError)

	ErrFailedRestoreAll = errors.NewErrorResponse("Failed to restore all Roles", http.StatusInternalServerError)
	ErrFailedDeleteAll  = errors.NewErrorResponse("Failed to delete all Roles permanently", http.StatusInternalServerError)
)
