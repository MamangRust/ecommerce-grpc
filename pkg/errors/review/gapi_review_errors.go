package review_errors

import (
	"ecommerce/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidID = errors.NewGrpcError("invalid ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateReview = errors.NewGrpcError("validation failed: invalid create review request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateReview = errors.NewGrpcError("validation failed: invalid update review request", int(codes.InvalidArgument))
)
