package repository

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
	reviewdetail_errors "ecommerce/pkg/errors/review_detail"
)

type reviewDetailRepository struct {
	db *db.Queries
}

func NewReviewDetailRepository(db *db.Queries) *reviewDetailRepository {
	return &reviewDetailRepository{
		db: db,
	}
}

func (r *reviewDetailRepository) FindAllReviews(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewDetailsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewDetails(ctx, reqDb)

	if err != nil {
		return nil, reviewdetail_errors.ErrFindAllReviewDetails
	}

	return res, nil
}

func (r *reviewDetailRepository) FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewDetailsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewDetailsActive(ctx, reqDb)

	if err != nil {
		return nil, reviewdetail_errors.ErrFindActiveReviewDetails
	}

	return res, nil
}

func (r *reviewDetailRepository) FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewDetailsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewDetailsTrashed(ctx, reqDb)

	if err != nil {
		return nil, reviewdetail_errors.ErrFindTrashedReviewDetails
	}

	return res, nil
}

func (r *reviewDetailRepository) FindById(ctx context.Context, user_id int) (*db.GetReviewDetailRow, error) {
	res, err := r.db.GetReviewDetail(ctx, int32(user_id))

	if err != nil {
		return nil, reviewdetail_errors.ErrFindByIdReviewDetail
	}

	return res, nil
}

func (r *reviewDetailRepository) FindByIdTrashed(ctx context.Context, user_id int) (*db.ReviewDetail, error) {
	res, err := r.db.GetReviewDetailTrashed(ctx, int32(user_id))

	if err != nil {
		return nil, reviewdetail_errors.ErrFindByIdTrashedReviewDetail
	}

	return res, nil
}

func (r *reviewDetailRepository) CreateReviewDetail(ctx context.Context, request *requests.CreateReviewDetailRequest) (*db.CreateReviewDetailRow, error) {
	req := db.CreateReviewDetailParams{
		ReviewID: int32(request.ReviewID),
		Type:     request.Type,
		Url:      request.Url,
		Caption:  stringPtr(request.Caption),
	}

	reviewDetail, err := r.db.CreateReviewDetail(ctx, req)
	if err != nil {
		return nil, reviewdetail_errors.ErrCreateReviewDetail
	}

	return reviewDetail, nil
}

func (r *reviewDetailRepository) UpdateReviewDetail(ctx context.Context, request *requests.UpdateReviewDetailRequest) (*db.UpdateReviewDetailRow, error) {
	req := db.UpdateReviewDetailParams{
		ReviewDetailID: int32(*request.ReviewDetailID),
		Type:           request.Type,
		Url:            request.Url,
		Caption:        stringPtr(request.Caption),
	}

	res, err := r.db.UpdateReviewDetail(ctx, req)
	if err != nil {
		return nil, reviewdetail_errors.ErrUpdateReviewDetail
	}

	return res, nil
}

func (r *reviewDetailRepository) TrashedReviewDetail(ctx context.Context, ReviewDetail_id int) (*db.ReviewDetail, error) {
	res, err := r.db.TrashReviewDetail(ctx, int32(ReviewDetail_id))

	if err != nil {
		return nil, reviewdetail_errors.ErrTrashedReviewDetail
	}

	return res, nil
}

func (r *reviewDetailRepository) RestoreReviewDetail(ctx context.Context, ReviewDetail_id int) (*db.ReviewDetail, error) {
	res, err := r.db.RestoreReviewDetail(ctx, int32(ReviewDetail_id))

	if err != nil {
		return nil, reviewdetail_errors.ErrRestoreReviewDetail
	}

	return res, nil
}

func (r *reviewDetailRepository) DeleteReviewDetailPermanent(ctx context.Context, ReviewDetail_id int) (bool, error) {
	err := r.db.DeletePermanentReviewDetail(ctx, int32(ReviewDetail_id))

	if err != nil {
		return false, reviewdetail_errors.ErrDeleteReviewDetailPermanent
	}

	return true, nil
}

func (r *reviewDetailRepository) RestoreAllReviewDetail(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllReviewDetails(ctx)

	if err != nil {
		return false, reviewdetail_errors.ErrRestoreAllReviewDetails
	}
	return true, nil
}

func (r *reviewDetailRepository) DeleteAllReviewDetailPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentReviewDetails(ctx)

	if err != nil {
		return false, reviewdetail_errors.ErrDeleteAllReviewDetails
	}
	return true, nil
}
