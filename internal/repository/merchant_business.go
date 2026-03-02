package repository

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
	merchantbusiness_errors "ecommerce/pkg/errors/merchant_business"
)

type merchantBusinessRepository struct {
	db *db.Queries
}

func NewMerchantBusinessRepository(
	db *db.Queries,
) *merchantBusinessRepository {
	return &merchantBusinessRepository{
		db: db,
	}
}

func (r *merchantBusinessRepository) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsBusinessInformationParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsBusinessInformation(ctx, reqDb)

	if err != nil {
		return nil, merchantbusiness_errors.ErrFindAllMerchantBusinesses
	}

	return res, nil
}

func (r *merchantBusinessRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsBusinessInformationActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsBusinessInformationActive(ctx, reqDb)

	if err != nil {
		return nil, merchantbusiness_errors.ErrFindActiveMerchantBusinesses
	}

	return res, nil
}

func (r *merchantBusinessRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsBusinessInformationTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsBusinessInformationTrashed(ctx, reqDb)

	if err != nil {
		return nil, merchantbusiness_errors.ErrFindTrashedMerchantBusinesses
	}

	return res, nil
}

func (r *merchantBusinessRepository) FindById(ctx context.Context, user_id int) (*db.GetMerchantBusinessInformationRow, error) {
	res, err := r.db.GetMerchantBusinessInformation(ctx, int32(user_id))

	if err != nil {
		return nil, merchantbusiness_errors.ErrMerchantBusinessNotFound
	}

	return res, nil
}

func (r *merchantBusinessRepository) CreateMerchantBusiness(
	ctx context.Context,
	request *requests.CreateMerchantBusinessInformationRequest,
) (*db.CreateMerchantBusinessInformationRow, error) {

	req := db.CreateMerchantBusinessInformationParams{
		MerchantID: int32(request.MerchantID),

		BusinessType: stringPtr(request.BusinessType),
		TaxID:        stringPtr(request.TaxID),
		WebsiteUrl:   stringPtr(request.WebsiteUrl),

		EstablishedYear:   int32Ptr(request.EstablishedYear),
		NumberOfEmployees: int32Ptr(request.NumberOfEmployees),
	}

	merchant, err := r.db.CreateMerchantBusinessInformation(ctx, req)
	if err != nil {
		return nil, merchantbusiness_errors.ErrCreateMerchantBusiness
	}

	return merchant, nil
}

func (r *merchantBusinessRepository) UpdateMerchantBusiness(ctx context.Context, request *requests.UpdateMerchantBusinessInformationRequest) (*db.UpdateMerchantBusinessInformationRow, error) {
	req := db.UpdateMerchantBusinessInformationParams{
		MerchantBusinessInfoID: int32(*request.MerchantBusinessInfoID),
		BusinessType:           stringPtr(request.BusinessType),
		TaxID:                  stringPtr(request.TaxID),
		WebsiteUrl:             stringPtr(request.WebsiteUrl),
		EstablishedYear:        int32Ptr(request.EstablishedYear),
		NumberOfEmployees:      int32Ptr(request.NumberOfEmployees),
	}

	merchant, err := r.db.UpdateMerchantBusinessInformation(ctx, req)
	if err != nil {
		return nil, merchantbusiness_errors.ErrUpdateMerchantBusiness
	}

	return merchant, nil
}

func (r *merchantBusinessRepository) TrashedMerchantBusiness(ctx context.Context, merchant_id int) (*db.MerchantBusinessInformation, error) {
	res, err := r.db.TrashMerchantBusinessInformation(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantbusiness_errors.ErrTrashMerchantBusiness
	}

	return res, nil
}

func (r *merchantBusinessRepository) RestoreMerchantBusiness(ctx context.Context, merchant_id int) (*db.MerchantBusinessInformation, error) {
	res, err := r.db.RestoreMerchantBusinessInformation(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantbusiness_errors.ErrRestoreMerchantBusiness
	}

	return res, nil
}

func (r *merchantBusinessRepository) DeleteMerchantBusinessPermanent(ctx context.Context, Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantBusinessInformationPermanently(ctx, int32(Merchant_id))

	if err != nil {
		return false, merchantbusiness_errors.ErrDeletePermanentMerchantBusiness
	}

	return true, nil
}

func (r *merchantBusinessRepository) RestoreAllMerchantBusiness(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchants(ctx)

	if err != nil {
		return false, merchantbusiness_errors.ErrRestoreAllMerchantBusinesses
	}
	return true, nil
}

func (r *merchantBusinessRepository) DeleteAllMerchantBusinessPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(ctx)

	if err != nil {
		return false, merchantbusiness_errors.ErrDeleteAllPermanentMerchantBusinesses
	}
	return true, nil
}

func int32Ptr(v int) *int32 {
	if v == 0 {
		return nil
	}
	val := int32(v)
	return &val
}
