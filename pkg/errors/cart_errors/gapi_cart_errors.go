package cart_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcCartNotFound  = response.NewGrpcError("error", "Cart not found", int(codes.NotFound))
	ErrGrpcCartInvalidId = response.NewGrpcError("error", "Invalid cart ID", int(codes.InvalidArgument))

	ErrGrpcFailedFindAll = response.NewGrpcError("error", "Failed to fetch carts", int(codes.Internal))

	ErrGrpcFailedCreateCart   = response.NewGrpcError("error", "Failed to create cart", int(codes.Internal))
	ErrGrpcValidateCreateCart = response.NewGrpcError("error", "Validation failed: invalid create cart request", int(codes.InvalidArgument))

	ErrGrpcFailedDeleteCart = response.NewGrpcError("error", "Failed to delete cart permanently", int(codes.Internal))
	ErrGrpcFailedDeleteAll  = response.NewGrpcError("error", "Failed to delete all carts permanently", int(codes.Internal))
)
