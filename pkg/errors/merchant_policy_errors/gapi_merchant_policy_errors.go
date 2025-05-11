package merchantpolicy_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidMerchantPolicyID = response.NewGrpcError("error", "Invalid merchant policy ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateMerchantPolicy = response.NewGrpcError("error", "Validation failed: invalid create merchant policy request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchantPolicy = response.NewGrpcError("error", "Validation failed: invalid update merchant policy request", int(codes.InvalidArgument))
)
