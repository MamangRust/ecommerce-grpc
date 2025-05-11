package banner_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcBannerNotFound  = response.NewGrpcError("error", "Banner not found", int(codes.NotFound))
	ErrGrpcBannerInvalidId = response.NewGrpcError("error", "Invalid Banner ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateBanner = response.NewGrpcError("error", "Validation failed: invalid create banner request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateBanner = response.NewGrpcError("error", "Validation failed: invalid update banner request", int(codes.InvalidArgument))
)
