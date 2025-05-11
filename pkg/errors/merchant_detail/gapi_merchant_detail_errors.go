package merchantdetail_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidMerchantDetailId = response.NewGrpcError("error", "invalid merchant detail ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateMerchantDetail = response.NewGrpcError("error", "Validation failed: invalid create merchant detail request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchantDetail = response.NewGrpcError("error", "Validation failed: invalid update merchant detail request", int(codes.InvalidArgument))
)
