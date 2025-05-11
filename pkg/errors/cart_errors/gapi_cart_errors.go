package cart_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcCartNotFound  = response.NewGrpcError("error", "Cart not found", int(codes.NotFound))
	ErrGrpcCartInvalidId = response.NewGrpcError("error", "Invalid cart ID", int(codes.InvalidArgument))

	ErrGrpcFailedCreateCart   = response.NewGrpcError("error", "Failed to create cart", int(codes.Internal))
	ErrGrpcValidateCreateCart = response.NewGrpcError("error", "Validation failed: invalid create cart request", int(codes.InvalidArgument))
)
