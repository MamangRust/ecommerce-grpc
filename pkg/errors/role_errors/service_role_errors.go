package role_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrRoleNotFoundRes   = response.NewErrorResponse("Role not found", http.StatusNotFound)
	ErrFailedFindAll     = response.NewErrorResponse("Failed to fetch Roles", http.StatusInternalServerError)
	ErrFailedFindActive  = response.NewErrorResponse("Failed to fetch active Roles", http.StatusInternalServerError)
	ErrFailedFindTrashed = response.NewErrorResponse("Failed to fetch trashed Roles", http.StatusInternalServerError)

	ErrFailedCreateRole = response.NewErrorResponse("Failed to create Role", http.StatusInternalServerError)
	ErrFailedUpdateRole = response.NewErrorResponse("Failed to update Role", http.StatusInternalServerError)

	ErrFailedTrashedRole     = response.NewErrorResponse("Failed to move Role to trash", http.StatusInternalServerError)
	ErrFailedRestoreRole     = response.NewErrorResponse("Failed to restore Role", http.StatusInternalServerError)
	ErrFailedDeletePermanent = response.NewErrorResponse("Failed to delete Role permanently", http.StatusInternalServerError)

	ErrFailedRestoreAll = response.NewErrorResponse("Failed to restore all Roles", http.StatusInternalServerError)
	ErrFailedDeleteAll  = response.NewErrorResponse("Failed to delete all Roles permanently", http.StatusInternalServerError)
)
