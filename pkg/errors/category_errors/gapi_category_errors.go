package category_errors

import (
	"ecommerce/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcValidateCreateCategory = errors.NewGrpcError("Validation failed: invalid create category request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateCategory = errors.NewGrpcError("Validation failed: invalid update category request", int(codes.InvalidArgument))

	ErrGrpcCategoryNotFound          = errors.NewGrpcError("Category not found", int(codes.NotFound))
	ErrGrpcCategoryInvalidId         = errors.NewGrpcError("Invalid category ID", int(codes.InvalidArgument))
	ErrGrpcCategoryInvalidYear       = errors.NewGrpcError("Invalid year", int(codes.InvalidArgument))
	ErrGrpcCategoryInvalidMonth      = errors.NewGrpcError("Invalid month", int(codes.InvalidArgument))
	ErrGrpcCategoryInvalidMerchantId = errors.NewGrpcError("Invalid merchant ID", int(codes.InvalidArgument))
)
