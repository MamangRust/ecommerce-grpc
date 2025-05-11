package slider_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidID = response.NewGrpcError("error", "invalid ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateSlider = response.NewGrpcError("error", "validation failed: invalid create slider request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateSlider = response.NewGrpcError("error", "validation failed: invalid update slider request", int(codes.InvalidArgument))
)
