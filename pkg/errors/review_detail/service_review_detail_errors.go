package reviewdetail_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrFailedImageNotFound = response.NewErrorResponse("Image not found", http.StatusNotFound)
	ErrFailedRemoveImage   = response.NewErrorResponse("Failed to remove image", http.StatusInternalServerError)

	ErrReviewDetailNotFoundRes = response.NewErrorResponse("Review Detail not found", http.StatusNotFound)
	ErrFailedFindAllReview     = response.NewErrorResponse("Failed to fetch Review Details", http.StatusInternalServerError)
	ErrFailedFindActiveReview  = response.NewErrorResponse("Failed to fetch active Review Details", http.StatusInternalServerError)
	ErrFailedFindTrashedReview = response.NewErrorResponse("Failed to fetch trashed Review Details", http.StatusInternalServerError)

	ErrFailedCreateReviewDetail = response.NewErrorResponse("Failed to create Review Detail", http.StatusInternalServerError)
	ErrFailedUpdateReviewDetail = response.NewErrorResponse("Failed to update Review Detail", http.StatusInternalServerError)

	ErrFailedTrashedReviewDetail   = response.NewErrorResponse("Failed to move Review Detail to trash", http.StatusInternalServerError)
	ErrFailedRestoreReviewDetail   = response.NewErrorResponse("Failed to restore Review Detail", http.StatusInternalServerError)
	ErrFailedDeletePermanentReview = response.NewErrorResponse("Failed to delete Review Detail permanently", http.StatusInternalServerError)

	ErrFailedRestoreAllReviewDetail = response.NewErrorResponse("Failed to restore all Review Details", http.StatusInternalServerError)
	ErrFailedDeleteAllReviewDetail  = response.NewErrorResponse("Failed to delete all Review Details permanently", http.StatusInternalServerError)
)
