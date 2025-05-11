package repository

import (
	"context"
	"database/sql"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	merchantbusiness_errors "ecommerce/pkg/errors/merchant_business"
)

type merchantBusinessRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantBusinessMapping
}

func NewMerchantBusinessRepository(
	db *db.Queries,
	ctx context.Context,
	mapping recordmapper.MerchantBusinessMapping,
) *merchantBusinessRepository {
	return &merchantBusinessRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantBusinessRepository) FindAllMerchants(req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsBusinessInformationParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsBusinessInformation(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchantbusiness_errors.ErrFindAllMerchantBusinesses
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantBusinesssRecordPagination(res), &totalCount, nil
}

func (r *merchantBusinessRepository) FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsBusinessInformationActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsBusinessInformationActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchantbusiness_errors.ErrFindActiveMerchantBusinesses
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantBusinesssRecordActivePagination(res), &totalCount, nil
}

func (r *merchantBusinessRepository) FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsBusinessInformationTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsBusinessInformationTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchantbusiness_errors.ErrFindTrashedMerchantBusinesses
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantBusinesssRecordTrashedPagination(res), &totalCount, nil
}

func (r *merchantBusinessRepository) FindById(user_id int) (*record.MerchantBusinessRecord, error) {
	res, err := r.db.GetMerchantBusinessInformation(r.ctx, int32(user_id))

	if err != nil {
		return nil, merchantbusiness_errors.ErrMerchantBusinessNotFound
	}

	return r.mapping.ToMerchantBusinessRecord(res), nil
}

func (r *merchantBusinessRepository) CreateMerchantBusiness(request *requests.CreateMerchantBusinessInformationRequest) (*record.MerchantBusinessRecord, error) {
	req := db.CreateMerchantBusinessInformationParams{
		MerchantID:        int32(request.MerchantID),
		BusinessType:      sql.NullString{String: request.BusinessType, Valid: request.BusinessType != ""},
		TaxID:             sql.NullString{String: request.TaxID, Valid: request.TaxID != ""},
		EstablishedYear:   sql.NullInt32{Int32: int32(request.EstablishedYear), Valid: request.EstablishedYear != 0},
		NumberOfEmployees: sql.NullInt32{Int32: int32(request.NumberOfEmployees), Valid: request.NumberOfEmployees != 0},
		WebsiteUrl:        sql.NullString{String: request.WebsiteUrl, Valid: request.WebsiteUrl != ""},
	}

	merchant, err := r.db.CreateMerchantBusinessInformation(r.ctx, req)
	if err != nil {
		return nil, merchantbusiness_errors.ErrCreateMerchantBusiness
	}

	return r.mapping.ToMerchantBusinessRecord(merchant), nil
}

func (r *merchantBusinessRepository) UpdateMerchantBusiness(request *requests.UpdateMerchantBusinessInformationRequest) (*record.MerchantBusinessRecord, error) {
	req := db.UpdateMerchantBusinessInformationParams{
		MerchantBusinessInfoID: int32(*request.MerchantBusinessInfoID),
		BusinessType:           sql.NullString{String: request.BusinessType, Valid: request.BusinessType != ""},
		TaxID:                  sql.NullString{String: request.TaxID, Valid: request.TaxID != ""},
		EstablishedYear:        sql.NullInt32{Int32: int32(request.EstablishedYear), Valid: request.EstablishedYear != 0},
		NumberOfEmployees:      sql.NullInt32{Int32: int32(request.NumberOfEmployees), Valid: request.NumberOfEmployees != 0},
		WebsiteUrl:             sql.NullString{String: request.WebsiteUrl, Valid: request.WebsiteUrl != ""},
	}

	merchant, err := r.db.UpdateMerchantBusinessInformation(r.ctx, req)
	if err != nil {
		return nil, merchantbusiness_errors.ErrUpdateMerchantBusiness
	}

	return r.mapping.ToMerchantBusinessRecord(merchant), nil
}

func (r *merchantBusinessRepository) TrashedMerchantBusiness(merchant_id int) (*record.MerchantBusinessRecord, error) {
	res, err := r.db.TrashMerchantBusinessInformation(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantbusiness_errors.ErrTrashMerchantBusiness
	}

	return r.mapping.ToMerchantBusinessRecord(res), nil
}

func (r *merchantBusinessRepository) RestoreMerchantBusiness(merchant_id int) (*record.MerchantBusinessRecord, error) {
	res, err := r.db.RestoreMerchantBusinessInformation(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantbusiness_errors.ErrRestoreMerchantBusiness
	}

	return r.mapping.ToMerchantBusinessRecord(res), nil
}

func (r *merchantBusinessRepository) DeleteMerchantBusinessPermanent(Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantBusinessInformationPermanently(r.ctx, int32(Merchant_id))

	if err != nil {
		return false, merchantbusiness_errors.ErrDeletePermanentMerchantBusiness
	}

	return true, nil
}

func (r *merchantBusinessRepository) RestoreAllMerchantBusiness() (bool, error) {
	err := r.db.RestoreAllMerchants(r.ctx)

	if err != nil {
		return false, merchantbusiness_errors.ErrRestoreAllMerchantBusinesses
	}
	return true, nil
}

func (r *merchantBusinessRepository) DeleteAllMerchantBusinessPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(r.ctx)

	if err != nil {
		return false, merchantbusiness_errors.ErrDeleteAllPermanentMerchantBusinesses
	}
	return true, nil
}
