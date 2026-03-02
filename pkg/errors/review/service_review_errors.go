package review_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrFailedReviewNotFound = errors.NewErrorResponse("Review not found", http.StatusNotFound)

	ErrFailedFindAllReviews        = errors.NewErrorResponse("Failed to fetch all reviews", http.StatusInternalServerError)
	ErrFailedFindActiveReviews     = errors.NewErrorResponse("Failed to fetch active reviews", http.StatusInternalServerError)
	ErrFailedFindTrashedReviews    = errors.NewErrorResponse("Failed to fetch trashed reviews", http.StatusInternalServerError)
	ErrFailedFindByProductReviews  = errors.NewErrorResponse("Failed to fetch reviews by product", http.StatusInternalServerError)
	ErrFailedFindByMerchantReviews = errors.NewErrorResponse("Failed to fetch reviews by merchant", http.StatusInternalServerError)

	ErrFailedCreateReview = errors.NewErrorResponse("Failed to create review", http.StatusInternalServerError)
	ErrFailedUpdateReview = errors.NewErrorResponse("Failed to update review", http.StatusInternalServerError)

	ErrFailedTrashedReview         = errors.NewErrorResponse("Failed to move review to trash", http.StatusInternalServerError)
	ErrFailedRestoreReview         = errors.NewErrorResponse("Failed to restore review from trash", http.StatusInternalServerError)
	ErrFailedDeletePermanentReview = errors.NewErrorResponse("Failed to permanently delete review", http.StatusInternalServerError)

	ErrFailedRestoreAllReviews         = errors.NewErrorResponse("Failed to restore all reviews", http.StatusInternalServerError)
	ErrFailedDeleteAllPermanentReviews = errors.NewErrorResponse("Failed to permanently delete all reviews", http.StatusInternalServerError)
)
