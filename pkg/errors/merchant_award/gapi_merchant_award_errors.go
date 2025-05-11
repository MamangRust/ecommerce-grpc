package merchantaward_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcMerchantInvalidId = response.NewGrpcError("error", "Invalid merchant ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateMerchantAward = response.NewGrpcError("error", "Validation failed: invalid create merchant award request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchantAward = response.NewGrpcError("error", "Validation failed: invalid update merchant award request", int(codes.InvalidArgument))
)
