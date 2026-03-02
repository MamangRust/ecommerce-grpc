package user_errors

import (
	"ecommerce/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcUserNotFound  = errors.NewGrpcError("User not found", int(codes.NotFound))
	ErrGrpcUserInvalidId = errors.NewGrpcError("Invalid User ID", int(codes.NotFound))

	ErrGrpcValidateCreateUser = errors.NewGrpcError("validation failed: invalid create User request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateUser = errors.NewGrpcError("validation failed: invalid update User request", int(codes.InvalidArgument))
)
