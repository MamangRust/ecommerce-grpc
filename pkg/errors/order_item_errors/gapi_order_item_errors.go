package orderitem_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidID = response.NewGrpcError("error", "invalid ID", int(codes.InvalidArgument))

	ErrGrpcFailedFindAllOrderItem       = response.NewGrpcError("error", "Failed to find all order items", int(codes.Internal))
	ErrGrpcFailedFindOrderItemByActive  = response.NewGrpcError("error", "Failed to find active order items", int(codes.Internal))
	ErrGrpcFailedFindOrderItemByTrashed = response.NewGrpcError("error", "Failed to find trashed order items", int(codes.Internal))
	ErrGrpcFailedFindOrderItemByOrder   = response.NewGrpcError("error", "Failed to find order items by order", int(codes.Internal))
)
