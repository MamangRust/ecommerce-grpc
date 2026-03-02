package merchantaward_errors

import (
	"ecommerce/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcMerchantInvalidId = errors.NewGrpcError("Invalid merchant ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateMerchantAward = errors.NewGrpcError("Validation failed: invalid create merchant award request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchantAward = errors.NewGrpcError("Validation failed: invalid update merchant award request", int(codes.InvalidArgument))
)
