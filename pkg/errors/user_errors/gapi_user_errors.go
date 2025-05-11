package user_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcUserNotFound  = response.NewGrpcError("error", "User not found", int(codes.NotFound))
	ErrGrpcUserInvalidId = response.NewGrpcError("error", "Invalid User ID", int(codes.NotFound))

	ErrGrpcValidateCreateUser = response.NewGrpcError("error", "validation failed: invalid create User request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateUser = response.NewGrpcError("error", "validation failed: invalid update User request", int(codes.InvalidArgument))
)
