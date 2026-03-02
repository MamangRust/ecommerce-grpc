package repository

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
	merchantdetail_errors "ecommerce/pkg/errors/merchant_detail"
)

type merchantDetailRepository struct {
	db *db.Queries
}

func NewMerchantDetailRepository(
	db *db.Queries,
) *merchantDetailRepository {
	return &merchantDetailRepository{
		db: db,
	}
}

func (r *merchantDetailRepository) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantDetailsParams{
		Column1: &req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantDetails(ctx, reqDb)

	if err != nil {
		return nil, merchantdetail_errors.ErrFindAllMerchantDetails
	}

	return res, nil
}

func (r *merchantDetailRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantDetailsActiveParams{
		Column1: &req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantDetailsActive(ctx, reqDb)

	if err != nil {
		return nil, merchantdetail_errors.ErrFindByActiveMerchantDetails
	}

	return res, nil
}

func (r *merchantDetailRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantDetailsTrashedParams{
		Column1: &req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantDetailsTrashed(ctx, reqDb)

	if err != nil {
		return nil, merchantdetail_errors.ErrFindByTrashedMerchantDetails
	}

	return res, nil
}

func (r *merchantDetailRepository) FindById(ctx context.Context, user_id int) (*db.GetMerchantDetailRow, error) {
	res, err := r.db.GetMerchantDetail(ctx, int32(user_id))

	if err != nil {
		return nil, merchantdetail_errors.ErrFindByIdMerchantDetail
	}

	return res, nil
}

func (r *merchantDetailRepository) FindByIdTrashed(ctx context.Context, user_id int) (*db.MerchantDetail, error) {
	res, err := r.db.GetMerchantDetailTrashed(ctx, int32(user_id))

	if err != nil {
		return nil, merchantdetail_errors.ErrFindByIdTrashedMerchantDetail
	}

	return res, nil
}

func (r *merchantDetailRepository) CreateMerchantDetail(
	ctx context.Context,
	request *requests.CreateMerchantDetailRequest,
) (*db.CreateMerchantDetailRow, error) {

	req := db.CreateMerchantDetailParams{
		MerchantID:       int32(request.MerchantID),
		DisplayName:      stringPtr(request.DisplayName),
		CoverImageUrl:    stringPtr(request.CoverImageUrl),
		LogoUrl:          stringPtr(request.LogoUrl),
		ShortDescription: stringPtr(request.ShortDescription),
		WebsiteUrl:       stringPtr(request.WebsiteUrl),
	}

	merchant, err := r.db.CreateMerchantDetail(ctx, req)
	if err != nil {
		return nil, merchantdetail_errors.ErrCreateMerchantDetail
	}

	return merchant, nil
}

func (r *merchantDetailRepository) UpdateMerchantDetail(ctx context.Context, request *requests.UpdateMerchantDetailRequest) (*db.UpdateMerchantDetailRow, error) {
	req := db.UpdateMerchantDetailParams{
		MerchantDetailID: int32(*request.MerchantDetailID),
		DisplayName:      stringPtr(request.DisplayName),
		CoverImageUrl:    stringPtr(request.CoverImageUrl),
		LogoUrl:          stringPtr(request.LogoUrl),
		ShortDescription: stringPtr(request.ShortDescription),
		WebsiteUrl:       stringPtr(request.WebsiteUrl),
	}

	res, err := r.db.UpdateMerchantDetail(ctx, req)
	if err != nil {
		return nil, merchantdetail_errors.ErrUpdateMerchantDetail
	}

	return res, nil
}

func (r *merchantDetailRepository) TrashedMerchantDetail(ctx context.Context, merchant_id int) (*db.MerchantDetail, error) {
	res, err := r.db.TrashMerchantDetail(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantdetail_errors.ErrTrashedMerchantDetail
	}

	return res, nil
}

func (r *merchantDetailRepository) RestoreMerchantDetail(ctx context.Context, merchant_id int) (*db.MerchantDetail, error) {
	res, err := r.db.RestoreMerchantDetail(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantdetail_errors.ErrRestoreMerchantDetail
	}

	return res, nil
}

func (r *merchantDetailRepository) DeleteMerchantDetailPermanent(ctx context.Context, Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantDetailPermanently(ctx, int32(Merchant_id))

	if err != nil {
		return false, merchantdetail_errors.ErrDeleteMerchantDetailPermanent
	}

	return true, nil
}

func (r *merchantDetailRepository) RestoreAllMerchantDetail(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchantDetails(ctx)

	if err != nil {
		return false, merchantdetail_errors.ErrRestoreAllMerchantDetails
	}
	return true, nil
}

func (r *merchantDetailRepository) DeleteAllMerchantDetailPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchantDetails(ctx)

	if err != nil {
		return false, merchantdetail_errors.ErrDeleteAllMerchantDetailsPermanent
	}
	return true, nil
}
