package reviewdetail_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidID = response.NewGrpcError("error", "invalid ID", int(codes.InvalidArgument))

	ErrGrpcFailedFindAllReviewDetails = response.NewGrpcError("error", "Failed to fetch all review details", int(codes.Internal))
	ErrGrpcFailedFindReviewDetailById = response.NewGrpcError("error", "Failed to fetch review detail by ID", int(codes.Internal))
	ErrGrpcFailedCreateReviewDetail   = response.NewGrpcError("error", "Failed to create review detail", int(codes.Internal))
	ErrGrpcFailedUpdateReviewDetail   = response.NewGrpcError("error", "Failed to update review detail", int(codes.Internal))
	ErrGrpcValidateCreateReviewDetail = response.NewGrpcError("error", "Validation failed: invalid create review detail request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateReviewDetail = response.NewGrpcError("error", "Validation failed: invalid update review detail request", int(codes.InvalidArgument))

	ErrGrpcFailedTrashReviewDetail      = response.NewGrpcError("error", "Failed to move review detail to trash", int(codes.Internal))
	ErrGrpcFailedRestoreReviewDetail    = response.NewGrpcError("error", "Failed to restore review detail from trash", int(codes.Internal))
	ErrGrpcFailedDeleteReviewDetail     = response.NewGrpcError("error", "Failed to delete review detail permanently", int(codes.Internal))
	ErrGrpcFailedRestoreAllReviewDetail = response.NewGrpcError("error", "Failed to restore all review details", int(codes.Internal))
	ErrGrpcFailedDeleteAllReviewDetail  = response.NewGrpcError("error", "Failed to delete all review details permanently", int(codes.Internal))
)
