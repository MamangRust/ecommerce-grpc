package review_errors

import "errors"

var (
	ErrFindAllReviews           = errors.New("failed to find all reviews")
	ErrFindReviewsByProduct     = errors.New("failed to find reviews by product")
	ErrFindReviewsByMerchant    = errors.New("failed to find reviews by merchant")
	ErrFindActiveReviews        = errors.New("failed to find active reviews")
	ErrFindTrashedReviews       = errors.New("failed to find trashed reviews")
	ErrFindReviewByID           = errors.New("failed to find review by ID")
	ErrCreateReview             = errors.New("failed to create review")
	ErrUpdateReview             = errors.New("failed to update review")
	ErrTrashReview              = errors.New("failed to move review to trash")
	ErrRestoreReview            = errors.New("failed to restore review from trash")
	ErrDeleteReviewPermanent    = errors.New("failed to permanently delete review")
	ErrRestoreAllReviews        = errors.New("failed to restore all reviews")
	ErrDeleteAllPermanentReview = errors.New("failed to permanently delete all reviews")
)
