package role_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcRoleNotFound  = response.NewGrpcError("error", "Role not found", int(codes.NotFound))
	ErrGrpcRoleInvalidId = response.NewGrpcError("error", "Invalid Role ID", int(codes.NotFound))

	ErrGrpcFailedFindAll     = response.NewGrpcError("error", "Failed to fetch Roles", int(codes.Internal))
	ErrGrpcFailedFindActive  = response.NewGrpcError("error", "Failed to fetch active Roles", int(codes.Internal))
	ErrGrpcFailedFindTrashed = response.NewGrpcError("error", "Failed to fetch trashed Roles", int(codes.Internal))

	ErrGrpcFailedCreateRole   = response.NewGrpcError("error", "Failed to create Role", int(codes.Internal))
	ErrGrpcFailedUpdateRole   = response.NewGrpcError("error", "Failed to update Role", int(codes.Internal))
	ErrGrpcValidateCreateRole = response.NewGrpcError("error", "validation failed: invalid create Role request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateRole = response.NewGrpcError("error", "validation failed: invalid update Role request", int(codes.InvalidArgument))

	ErrGrpcFailedTrashedRole     = response.NewGrpcError("error", "Failed to move Role to trash", int(codes.Internal))
	ErrGrpcFailedRestoreRole     = response.NewGrpcError("error", "Failed to restore Role", int(codes.Internal))
	ErrGrpcFailedDeletePermanent = response.NewGrpcError("error", "Failed to delete Role permanently", int(codes.Internal))

	ErrGrpcFailedRestoreAll = response.NewGrpcError("error", "Failed to restore all Roles", int(codes.Internal))
	ErrGrpcFailedDeleteAll  = response.NewGrpcError("error", "Failed to delete all Roles permanently", int(codes.Internal))
)
