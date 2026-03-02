package repository

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
	review_errors "ecommerce/pkg/errors/review"
)

type reviewRepository struct {
	db *db.Queries
}

func NewReviewRepository(
	db *db.Queries,
) *reviewRepository {
	return &reviewRepository{
		db: db,
	}
}

func (r *reviewRepository) FindAllReview(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviews(ctx, reqDb)

	if err != nil {
		return nil, review_errors.ErrFindAllReviews
	}

	return res, nil
}

func (r *reviewRepository) FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*db.GetReviewByProductIdRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewByProductIdParams{
		ProductID: int32(req.ProductID),
		Column2:   int32(req.Rating),
		Limit:     int32(req.PageSize),
		Offset:    int32(offset),
	}

	res, err := r.db.GetReviewByProductId(ctx, reqDb)

	if err != nil {
		return nil, review_errors.ErrFindReviewsByProduct
	}

	return res, nil
}

func (r *reviewRepository) FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*db.GetReviewByMerchantIdRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewByMerchantIdParams{
		MerchantID: int32(req.MerchantID),
		Column2:    int32(req.Rating),
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetReviewByMerchantId(ctx, reqDb)

	if err != nil {
		return nil, review_errors.ErrFindReviewsByMerchant
	}

	return res, nil
}

func (r *reviewRepository) FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewsActive(ctx, reqDb)

	if err != nil {
		return nil, review_errors.ErrFindActiveReviews
	}

	return res, nil
}

func (r *reviewRepository) FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewsTrashed(ctx, reqDb)

	if err != nil {
		return nil, review_errors.ErrFindTrashedReviews
	}

	return res, nil
}

func (r *reviewRepository) FindById(ctx context.Context, id int) (*db.GetReviewByIDRow, error) {
	res, err := r.db.GetReviewByID(ctx, int32(id))

	if err != nil {
		return nil, review_errors.ErrFindReviewByID
	}

	return res, nil
}

func (r *reviewRepository) CreateReview(ctx context.Context, request *requests.CreateReviewRequest) (*db.CreateReviewRow, error) {
	req := db.CreateReviewParams{
		UserID:    int32(request.UserID),
		ProductID: int32(request.ProductID),
		Rating:    int32(request.Rating),
		Comment:   request.Comment,
	}

	review, err := r.db.CreateReview(ctx, req)

	if err != nil {
		return nil, review_errors.ErrCreateReview
	}

	return review, nil
}

func (r *reviewRepository) UpdateReview(ctx context.Context, request *requests.UpdateReviewRequest) (*db.UpdateReviewRow, error) {
	req := db.UpdateReviewParams{
		ReviewID: int32(*request.ReviewID),
		Name:     request.Name,
		Rating:   int32(request.Rating),
		Comment:  request.Comment,
	}

	res, err := r.db.UpdateReview(ctx, req)

	if err != nil {
		return nil, review_errors.ErrUpdateReview
	}

	return res, nil
}

func (r *reviewRepository) TrashReview(ctx context.Context, shipping_id int) (*db.Review, error) {
	res, err := r.db.TrashReview(ctx, int32(shipping_id))

	if err != nil {
		return nil, review_errors.ErrTrashReview
	}

	return res, nil
}

func (r *reviewRepository) RestoreReview(ctx context.Context, category_id int) (*db.Review, error) {
	res, err := r.db.RestoreReview(ctx, int32(category_id))

	if err != nil {
		return nil, review_errors.ErrRestoreReview
	}

	return res, nil
}

func (r *reviewRepository) DeleteReviewPermanently(ctx context.Context, category_id int) (bool, error) {
	err := r.db.DeleteReviewPermanently(ctx, int32(category_id))

	if err != nil {
		return false, review_errors.ErrDeleteReviewPermanent
	}

	return true, nil
}

func (r *reviewRepository) RestoreAllReview(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllReviews(ctx)

	if err != nil {
		return false, review_errors.ErrRestoreAllReviews
	}
	return true, nil
}

func (r *reviewRepository) DeleteAllPermanentReview(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentReviews(ctx)

	if err != nil {
		return false, review_errors.ErrDeleteAllPermanentReview
	}
	return true, nil
}
