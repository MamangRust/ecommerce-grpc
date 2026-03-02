package merchant_errors

import (
	"ecommerce/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidMerchantId = errors.NewGrpcError("invalid merchant ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateMerchant = errors.NewGrpcError("Validation failed: invalid create merchant request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchant = errors.NewGrpcError("Validation failed: invalid update merchant request", int(codes.InvalidArgument))
)
