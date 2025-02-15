package repository

import (
	"context"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"errors"
	"fmt"
)

type reviewRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.ReviewRecordMapping
}

func NewReviewRepository(
	db *db.Queries,
	ctx context.Context,
	mapping recordmapper.ReviewRecordMapping,
) *reviewRepository {
	return &reviewRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *reviewRepository) FindAllReview(search string, page, pageSize int) ([]*record.ReviewRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetReviewsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviews(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find shipping address: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewsRecordPagination(res), totalCount, nil
}

func (r *reviewRepository) FindByProduct(product_id int, search string, page, pageSize int) ([]*record.ReviewRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetReviewsByProductIDParams{
		ProductID: int32(product_id),
		Column2:   search,
		Limit:     int32(pageSize),
		Offset:    int32(offset),
	}

	res, err := r.db.GetReviewsByProductID(r.ctx, req)

	if err != nil {

		return nil, 0, fmt.Errorf("failed to find shipping address: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewsProductRecordPagination(res), totalCount, nil
}

func (r *reviewRepository) FindByActive(search string, page, pageSize int) ([]*record.ReviewRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetReviewsActiveParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewsActive(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find shipping address: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewsRecordActivePagination(res), totalCount, nil
}

func (r *reviewRepository) FindByTrashed(search string, page, pageSize int) ([]*record.ReviewRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetReviewsTrashedParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewsTrashed(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find shipping address: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewsRecordTrashedPagination(res), totalCount, nil
}

func (r *reviewRepository) FindById(id int) (*record.ReviewRecord, error) {
	res, err := r.db.GetReviewByID(r.ctx, int32(id))

	if err != nil {
		fmt.Printf("Error fetching review: %v\n", err)

		return nil, fmt.Errorf("failed to find review: %w", err)
	}

	return r.mapping.ToReviewRecord(res), nil
}

func (r *reviewRepository) CreateReview(request *requests.CreateReviewRequest) (*record.ReviewRecord, error) {
	req := db.CreateReviewParams{
		UserID:    int32(request.UserID),
		ProductID: int32(request.ProductID),
		Rating:    int32(request.Rating),
		Comment:   request.Comment,
	}

	review, err := r.db.CreateReview(r.ctx, req)
	if err != nil {
		return nil, errors.New("failed to create review")
	}

	return r.mapping.ToReviewRecord(review), nil
}

func (r *reviewRepository) UpdateReview(request *requests.UpdateReviewRequest) (*record.ReviewRecord, error) {
	req := db.UpdateReviewParams{
		ReviewID: int32(request.ReviewID),
		Name:     request.Name,
		Rating:   int32(request.Rating),
		Comment:  request.Comment,
	}

	res, err := r.db.UpdateReview(r.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	return r.mapping.ToReviewRecord(res), nil
}

func (r *reviewRepository) TrashReview(shipping_id int) (*record.ReviewRecord, error) {
	res, err := r.db.TrashReview(r.ctx, int32(shipping_id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash shipping address: %w", err)
	}

	return r.mapping.ToReviewRecord(res), nil
}

func (r *reviewRepository) RestoreReview(category_id int) (*record.ReviewRecord, error) {
	res, err := r.db.RestoreReview(r.ctx, int32(category_id))

	if err != nil {
		return nil, fmt.Errorf("failed to shipping address: %w", err)
	}

	return r.mapping.ToReviewRecord(res), nil
}

func (r *reviewRepository) DeleteReviewPermanently(category_id int) (bool, error) {
	err := r.db.DeleteReviewPermanently(r.ctx, int32(category_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete shipping address: %w", err)
	}

	return true, nil
}

func (r *reviewRepository) RestoreAllReview() (bool, error) {
	err := r.db.RestoreAllReviews(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all shipping address: %w", err)
	}
	return true, nil
}

func (r *reviewRepository) DeleteAllPermanentReview() (bool, error) {
	err := r.db.DeleteAllPermanentReviews(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all shipping address permanently: %w", err)
	}
	return true, nil
}
