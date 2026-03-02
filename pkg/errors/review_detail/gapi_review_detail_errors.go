package reviewdetail_errors

import (
	"ecommerce/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidID = errors.NewGrpcError("invalid ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateReviewDetail = errors.NewGrpcError("Validation failed: invalid create review detail request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateReviewDetail = errors.NewGrpcError("Validation failed: invalid update review detail request", int(codes.InvalidArgument))
)
