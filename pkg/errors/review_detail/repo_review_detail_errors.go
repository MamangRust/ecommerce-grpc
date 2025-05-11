package reviewdetail_errors

import "errors"

var (
	ErrReviewDetailNotFound     = errors.New("review detail not found")
	ErrFindAllReviewDetails     = errors.New("failed to find all review details")
	ErrFindActiveReviewDetails  = errors.New("failed to find active review details")
	ErrFindTrashedReviewDetails = errors.New("failed to find trashed review details")
	ErrReviewDetailConflict     = errors.New("failed, review detail already exists")

	ErrFindByIdReviewDetail        = errors.New("failed to find review detail by ID")
	ErrFindByIdTrashedReviewDetail = errors.New("failed to find trashed review detail by ID")

	ErrCreateReviewDetail = errors.New("failed to create review detail")
	ErrUpdateReviewDetail = errors.New("failed to update review detail")

	ErrTrashedReviewDetail         = errors.New("failed to move review detail to trash")
	ErrRestoreReviewDetail         = errors.New("failed to restore review detail from trash")
	ErrDeleteReviewDetailPermanent = errors.New("failed to permanently delete review detail")

	ErrRestoreAllReviewDetails = errors.New("failed to restore all review details")
	ErrDeleteAllReviewDetails  = errors.New("failed to permanently delete all review details")
)
