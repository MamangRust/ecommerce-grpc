package repository

import (
	"context"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	review_errors "ecommerce/pkg/errors/review"
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

func (r *reviewRepository) FindAllReview(req *requests.FindAllReview) ([]*record.ReviewRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviews(r.ctx, reqDb)

	if err != nil {
		return nil, nil, review_errors.ErrFindAllReviews
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewsRecordPagination(res), &totalCount, nil
}

func (r *reviewRepository) FindByProduct(req *requests.FindAllReviewByProduct) ([]*record.ReviewsDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewByProductIdParams{
		ProductID: int32(req.ProductID),
		Column2:   int32(req.Rating),
		Limit:     int32(req.PageSize),
		Offset:    int32(offset),
	}

	res, err := r.db.GetReviewByProductId(r.ctx, reqDb)

	if err != nil {
		return nil, nil, review_errors.ErrFindReviewsByProduct
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewsProductRecordPagination(res), &totalCount, nil
}

func (r *reviewRepository) FindByMerchant(req *requests.FindAllReviewByMerchant) ([]*record.ReviewsDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewByMerchantIdParams{
		MerchantID: int32(req.MerchantID),
		Column2:    int32(req.Rating),
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetReviewByMerchantId(r.ctx, reqDb)

	if err != nil {
		return nil, nil, review_errors.ErrFindReviewsByMerchant
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewsMerchantRecordPagination(res), &totalCount, nil
}

func (r *reviewRepository) FindByActive(req *requests.FindAllReview) ([]*record.ReviewRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewsActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, review_errors.ErrFindActiveReviews
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewsRecordActivePagination(res), &totalCount, nil
}

func (r *reviewRepository) FindByTrashed(req *requests.FindAllReview) ([]*record.ReviewRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewsTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, review_errors.ErrFindTrashedReviews
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewsRecordTrashedPagination(res), &totalCount, nil
}

func (r *reviewRepository) FindById(id int) (*record.ReviewRecord, error) {
	res, err := r.db.GetReviewByID(r.ctx, int32(id))

	if err != nil {
		return nil, review_errors.ErrFindReviewByID
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
		return nil, review_errors.ErrCreateReview
	}

	return r.mapping.ToReviewRecord(review), nil
}

func (r *reviewRepository) UpdateReview(request *requests.UpdateReviewRequest) (*record.ReviewRecord, error) {
	req := db.UpdateReviewParams{
		ReviewID: int32(*request.ReviewID),
		Name:     request.Name,
		Rating:   int32(request.Rating),
		Comment:  request.Comment,
	}

	res, err := r.db.UpdateReview(r.ctx, req)

	if err != nil {
		return nil, review_errors.ErrUpdateReview
	}

	return r.mapping.ToReviewRecord(res), nil
}

func (r *reviewRepository) TrashReview(shipping_id int) (*record.ReviewRecord, error) {
	res, err := r.db.TrashReview(r.ctx, int32(shipping_id))

	if err != nil {
		return nil, review_errors.ErrTrashReview
	}

	return r.mapping.ToReviewRecord(res), nil
}

func (r *reviewRepository) RestoreReview(category_id int) (*record.ReviewRecord, error) {
	res, err := r.db.RestoreReview(r.ctx, int32(category_id))

	if err != nil {
		return nil, review_errors.ErrRestoreReview
	}

	return r.mapping.ToReviewRecord(res), nil
}

func (r *reviewRepository) DeleteReviewPermanently(category_id int) (bool, error) {
	err := r.db.DeleteReviewPermanently(r.ctx, int32(category_id))

	if err != nil {
		return false, review_errors.ErrDeleteReviewPermanent
	}

	return true, nil
}

func (r *reviewRepository) RestoreAllReview() (bool, error) {
	err := r.db.RestoreAllReviews(r.ctx)

	if err != nil {
		return false, review_errors.ErrRestoreAllReviews
	}
	return true, nil
}

func (r *reviewRepository) DeleteAllPermanentReview() (bool, error) {
	err := r.db.DeleteAllPermanentReviews(r.ctx)

	if err != nil {
		return false, review_errors.ErrDeleteAllPermanentReview
	}
	return true, nil
}
