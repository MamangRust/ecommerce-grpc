package review_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidID = response.NewGrpcError("error", "invalid ID", int(codes.InvalidArgument))

	ErrGrpcFailedFindAllReviews        = response.NewGrpcError("error", "Failed to fetch all reviews", int(codes.Internal))
	ErrGrpcFailedFindByProductReviews  = response.NewGrpcError("error", "Failed to fetch reviews by product", int(codes.Internal))
	ErrGrpcFailedFindByMerchantReviews = response.NewGrpcError("error", "Failed to fetch reviews by merchant", int(codes.Internal))
	ErrGrpcFailedFindTrashedReviews    = response.NewGrpcError("error", "Failed to fetch trashed reviews", int(codes.Internal))
	ErrGrpcFailedFindActiveReviews     = response.NewGrpcError("error", "Failed to fetch active reviews", int(codes.Internal))

	ErrGrpcFailedCreateReview = response.NewGrpcError("error", "Failed to create review", int(codes.Internal))
	ErrGrpcFailedUpdateReview = response.NewGrpcError("error", "Failed to update review", int(codes.Internal))

	ErrGrpcFailedTrashedReview         = response.NewGrpcError("error", "Failed to move review to trash", int(codes.Internal))
	ErrGrpcFailedRestoreReview         = response.NewGrpcError("error", "Failed to restore review from trash", int(codes.Internal))
	ErrGrpcFailedDeletePermanentReview = response.NewGrpcError("error", "Failed to permanently delete review", int(codes.Internal))

	ErrGrpcFailedRestoreAllReviews         = response.NewGrpcError("error", "Failed to restore all reviews", int(codes.Internal))
	ErrGrpcFailedDeleteAllPermanentReviews = response.NewGrpcError("error", "Failed to permanently delete all reviews", int(codes.Internal))

	ErrGrpcValidateCreateReview = response.NewGrpcError("error", "validation failed: invalid create review request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateReview = response.NewGrpcError("error", "validation failed: invalid update review request", int(codes.InvalidArgument))
)
