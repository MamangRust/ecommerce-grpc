package category_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcValidateCreateCategory = response.NewGrpcError("error", "Validation failed: invalid create category request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateCategory = response.NewGrpcError("error", "Validation failed: invalid update category request", int(codes.InvalidArgument))

	ErrGrpcCategoryNotFound          = response.NewGrpcError("error", "Category not found", int(codes.NotFound))
	ErrGrpcCategoryInvalidId         = response.NewGrpcError("error", "Invalid category ID", int(codes.InvalidArgument))
	ErrGrpcCategoryInvalidYear       = response.NewGrpcError("error", "Invalid year", int(codes.InvalidArgument))
	ErrGrpcCategoryInvalidMonth      = response.NewGrpcError("error", "Invalid month", int(codes.InvalidArgument))
	ErrGrpcCategoryInvalidMerchantId = response.NewGrpcError("error", "Invalid merchant ID", int(codes.InvalidArgument))
)
