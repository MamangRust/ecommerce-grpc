package merchantpolicy_errors

import (
	"ecommerce/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidMerchantPolicyID = errors.NewGrpcError("Invalid merchant policy ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateMerchantPolicy = errors.NewGrpcError("Validation failed: invalid create merchant policy request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchantPolicy = errors.NewGrpcError("Validation failed: invalid update merchant policy request", int(codes.InvalidArgument))
)
