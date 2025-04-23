package repository

import (
	"context"
	"database/sql"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"fmt"
)

type reviewDetailRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.ReviewDetailRecordMapping
}

func NewReviewDetailRepository(db *db.Queries, ctx context.Context, mapping recordmapper.ReviewDetailRecordMapping) *reviewDetailRepository {
	return &reviewDetailRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *reviewDetailRepository) FindAllReviews(req *requests.FindAllReview) ([]*record.ReviewDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewDetailsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewDetails(r.ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch ReviewDetails: invalid pagination (page %d, size %d) or search query '%s'", req.Page, req.PageSize, req.Search)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewDetailsRecordPagination(res), &totalCount, nil
}

func (r *reviewDetailRepository) FindByActive(req *requests.FindAllReview) ([]*record.ReviewDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewDetailsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewDetailsActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch ReviewDetails active: invalid pagination (page %d, size %d) or search query '%s'", req.Page, req.PageSize, req.Search)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewDetailsRecordActivePagination(res), &totalCount, nil
}

func (r *reviewDetailRepository) FindByTrashed(req *requests.FindAllReview) ([]*record.ReviewDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewDetailsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewDetailsTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch ReviewDetails trashed: invalid pagination (page %d, size %d) or search query '%s'", req.Page, req.PageSize, req.Search)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewDetailsRecordTrashedPagination(res), &totalCount, nil
}

func (r *reviewDetailRepository) FindById(user_id int) (*record.ReviewDetailRecord, error) {
	res, err := r.db.GetReviewDetail(r.ctx, int32(user_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find ReviewDetail: %w", err)
	}

	return r.mapping.ToReviewDetailRecord(res), nil
}

func (r *reviewDetailRepository) FindByIdTrashed(user_id int) (*record.ReviewDetailRecord, error) {
	res, err := r.db.GetReviewDetailTrashed(r.ctx, int32(user_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find ReviewDetail: %w", err)
	}

	return r.mapping.ToReviewDetailRecord(res), nil
}

func (r *reviewDetailRepository) CreateReviewDetail(request *requests.CreateReviewDetailRequest) (*record.ReviewDetailRecord, error) {
	req := db.CreateReviewDetailParams{
		ReviewID: int32(request.ReviewID),
		Type:     request.Type,
		Url:      request.Url,
		Caption:  sql.NullString{String: request.Caption, Valid: request.Caption != ""},
	}

	reviewDetail, err := r.db.CreateReviewDetail(r.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create review detail: %w", err)
	}

	return r.mapping.ToReviewDetailRecord(reviewDetail), nil
}

func (r *reviewDetailRepository) UpdateReviewDetail(request *requests.UpdateReviewDetailRequest) (*record.ReviewDetailRecord, error) {
	req := db.UpdateReviewDetailParams{
		ReviewDetailID: int32(*request.ReviewDetailID),
		Type:           request.Type,
		Url:            request.Url,
		Caption:        sql.NullString{String: request.Caption, Valid: request.Caption != ""},
	}

	res, err := r.db.UpdateReviewDetail(r.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update review detail: %w", err)
	}

	return r.mapping.ToReviewDetailRecord(res), nil
}

func (r *reviewDetailRepository) TrashedReviewDetail(ReviewDetail_id int) (*record.ReviewDetailRecord, error) {
	res, err := r.db.TrashReviewDetail(r.ctx, int32(ReviewDetail_id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash ReviewDetail: %w", err)
	}

	return r.mapping.ToReviewDetailRecord(res), nil
}

func (r *reviewDetailRepository) RestoreReviewDetail(ReviewDetail_id int) (*record.ReviewDetailRecord, error) {
	res, err := r.db.RestoreReviewDetail(r.ctx, int32(ReviewDetail_id))

	if err != nil {
		return nil, fmt.Errorf("failed to restore ReviewDetails: %w", err)
	}

	return r.mapping.ToReviewDetailRecord(res), nil
}

func (r *reviewDetailRepository) DeleteReviewDetailPermanent(ReviewDetail_id int) (bool, error) {
	err := r.db.DeletePermanentReviewDetail(r.ctx, int32(ReviewDetail_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete ReviewDetail: %w", err)
	}

	return true, nil
}

func (r *reviewDetailRepository) RestoreAllReviewDetail() (bool, error) {
	err := r.db.RestoreAllReviewDetails(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all ReviewDetails: %w", err)
	}
	return true, nil
}

func (r *reviewDetailRepository) DeleteAllReviewDetailPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentReviewDetails(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all ReviewDetails permanently: %w", err)
	}
	return true, nil
}
