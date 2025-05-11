package review_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrFailedReviewNotFound = response.NewErrorResponse("Review not found", http.StatusNotFound)

	ErrFailedFindAllReviews        = response.NewErrorResponse("Failed to fetch all reviews", http.StatusInternalServerError)
	ErrFailedFindActiveReviews     = response.NewErrorResponse("Failed to fetch active reviews", http.StatusInternalServerError)
	ErrFailedFindTrashedReviews    = response.NewErrorResponse("Failed to fetch trashed reviews", http.StatusInternalServerError)
	ErrFailedFindByProductReviews  = response.NewErrorResponse("Failed to fetch reviews by product", http.StatusInternalServerError)
	ErrFailedFindByMerchantReviews = response.NewErrorResponse("Failed to fetch reviews by merchant", http.StatusInternalServerError)

	ErrFailedCreateReview = response.NewErrorResponse("Failed to create review", http.StatusInternalServerError)
	ErrFailedUpdateReview = response.NewErrorResponse("Failed to update review", http.StatusInternalServerError)

	ErrFailedTrashedReview         = response.NewErrorResponse("Failed to move review to trash", http.StatusInternalServerError)
	ErrFailedRestoreReview         = response.NewErrorResponse("Failed to restore review from trash", http.StatusInternalServerError)
	ErrFailedDeletePermanentReview = response.NewErrorResponse("Failed to permanently delete review", http.StatusInternalServerError)

	ErrFailedRestoreAllReviews         = response.NewErrorResponse("Failed to restore all reviews", http.StatusInternalServerError)
	ErrFailedDeleteAllPermanentReviews = response.NewErrorResponse("Failed to permanently delete all reviews", http.StatusInternalServerError)
)
