package userrole_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrFailedAssignRoleToUser = response.NewErrorResponse("Failed to assign role to user", http.StatusInternalServerError)
	ErrFailedRemoveRole       = response.NewErrorResponse("Failed to remove role from user", http.StatusInternalServerError)
)
