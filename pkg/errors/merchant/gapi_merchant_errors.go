package merchant_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidMerchantId = response.NewGrpcError("error", "invalid merchant ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateMerchant = response.NewGrpcError("error", "Validation failed: invalid create merchant request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchant = response.NewGrpcError("error", "Validation failed: invalid update merchant request", int(codes.InvalidArgument))
)
