package merchantbusiness_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcValidateCreateMerchantBusiness = response.NewGrpcError("error", "Validation failed: invalid create merchant business request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchantBusiness = response.NewGrpcError("error", "Validation failed: invalid update merchant business request", int(codes.InvalidArgument))

	ErrGrpcMerchantBusinessNotFound  = response.NewGrpcError("error", "Merchant business not found", int(codes.NotFound))
	ErrGrpcInvalidMerchantBusinessId = response.NewGrpcError("error", "Invalid merchant business ID", int(codes.InvalidArgument))
)
