package reviewdetail_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidID = response.NewGrpcError("error", "invalid ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateReviewDetail = response.NewGrpcError("error", "Validation failed: invalid create review detail request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateReviewDetail = response.NewGrpcError("error", "Validation failed: invalid update review detail request", int(codes.InvalidArgument))
)
