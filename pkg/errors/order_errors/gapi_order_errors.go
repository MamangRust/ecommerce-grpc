package order_errors

import (
	"ecommerce/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidYear             = errors.NewGrpcError("Invalid year", int(codes.InvalidArgument))
	ErrGrpcInvalidMonth            = errors.NewGrpcError("Invalid month", int(codes.InvalidArgument))
	ErrGrpcFailedInvalidMerchantId = errors.NewGrpcError("Invalid merchant ID", int(codes.InvalidArgument))
	ErrGrpcFailedInvalidId         = errors.NewGrpcError("Invalid ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateOrder = errors.NewGrpcError("validation failed: invalid create order request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateOrder = errors.NewGrpcError("validation failed: invalid update order request", int(codes.InvalidArgument))
)
