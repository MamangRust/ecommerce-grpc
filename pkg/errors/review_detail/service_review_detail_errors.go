package reviewdetail_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrFailedImageNotFound = errors.NewErrorResponse("Image not found", http.StatusNotFound)
	ErrFailedRemoveImage   = errors.NewErrorResponse("Failed to remove image", http.StatusInternalServerError)

	ErrReviewDetailNotFoundRes = errors.NewErrorResponse("Review Detail not found", http.StatusNotFound)
	ErrFailedFindAllReview     = errors.NewErrorResponse("Failed to fetch Review Details", http.StatusInternalServerError)
	ErrFailedFindActiveReview  = errors.NewErrorResponse("Failed to fetch active Review Details", http.StatusInternalServerError)
	ErrFailedFindTrashedReview = errors.NewErrorResponse("Failed to fetch trashed Review Details", http.StatusInternalServerError)

	ErrFailedCreateReviewDetail = errors.NewErrorResponse("Failed to create Review Detail", http.StatusInternalServerError)
	ErrFailedUpdateReviewDetail = errors.NewErrorResponse("Failed to update Review Detail", http.StatusInternalServerError)

	ErrFailedTrashedReviewDetail   = errors.NewErrorResponse("Failed to move Review Detail to trash", http.StatusInternalServerError)
	ErrFailedRestoreReviewDetail   = errors.NewErrorResponse("Failed to restore Review Detail", http.StatusInternalServerError)
	ErrFailedDeletePermanentReview = errors.NewErrorResponse("Failed to delete Review Detail permanently", http.StatusInternalServerError)

	ErrFailedRestoreAllReviewDetail = errors.NewErrorResponse("Failed to restore all Review Details", http.StatusInternalServerError)
	ErrFailedDeleteAllReviewDetail  = errors.NewErrorResponse("Failed to delete all Review Details permanently", http.StatusInternalServerError)
)
