package cart_errors

import (
	"ecommerce/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcCartNotFound  = errors.NewGrpcError("Cart not found", int(codes.NotFound))
	ErrGrpcCartInvalidId = errors.NewGrpcError("Invalid cart ID", int(codes.InvalidArgument))

	ErrGrpcFailedCreateCart   = errors.NewGrpcError("Failed to create cart", int(codes.Internal))
	ErrGrpcValidateCreateCart = errors.NewGrpcError("Validation failed: invalid create cart request", int(codes.InvalidArgument))
)
