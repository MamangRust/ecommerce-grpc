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

type merchantRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantRecordMapping
}

func NewMerchantRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantRecordMapping) *merchantRepository {
	return &merchantRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantRepository) FindAllMerchants(req *requests.FindAllMerchant) ([]*record.MerchantRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchants(r.ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch merchants: invalid pagination (page %d, size %d) or search query '%s'", req.Page, req.PageSize, req.Search)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsRecordPagination(res), &totalCount, nil
}

func (r *merchantRepository) FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch merchants active: invalid pagination (page %d, size %d) or search query '%s'", req.Page, req.PageSize, req.Search)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsRecordActivePagination(res), &totalCount, nil
}

func (r *merchantRepository) FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch merchants trashed: invalid pagination (page %d, size %d) or search query '%s'", req.Page, req.PageSize, req.Search)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsRecordTrashedPagination(res), &totalCount, nil
}

func (r *merchantRepository) FindById(user_id int) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByID(r.ctx, int32(user_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find merchant: %w", err)
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) CreateMerchant(request *requests.CreateMerchantRequest) (*record.MerchantRecord, error) {
	req := db.CreateMerchantParams{
		UserID:       int32(request.UserID),
		Name:         request.Name,
		Description:  sql.NullString{String: request.Description, Valid: true},
		Address:      sql.NullString{String: request.Address, Valid: true},
		ContactEmail: sql.NullString{String: request.ContactEmail, Valid: true},
		ContactPhone: sql.NullString{String: request.ContactPhone, Valid: true},
		Status:       "active",
	}

	merchant, err := r.db.CreateMerchant(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to create merchant: %w", err)
	}

	return r.mapping.ToMerchantRecord(merchant), nil
}

func (r *merchantRepository) UpdateMerchant(request *requests.UpdateMerchantRequest) (*record.MerchantRecord, error) {
	req := db.UpdateMerchantParams{
		MerchantID:   int32(*request.MerchantID),
		Name:         request.Name,
		Description:  sql.NullString{String: request.Description, Valid: true},
		Address:      sql.NullString{String: request.Address, Valid: true},
		ContactEmail: sql.NullString{String: request.ContactEmail, Valid: true},
		ContactPhone: sql.NullString{String: request.ContactPhone, Valid: true},
		Status:       request.Status,
	}

	res, err := r.db.UpdateMerchant(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update Merchant: %w", err)
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) TrashedMerchant(merchant_id int) (*record.MerchantRecord, error) {
	res, err := r.db.TrashMerchant(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash Merchant: %w", err)
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) RestoreMerchant(merchant_id int) (*record.MerchantRecord, error) {
	res, err := r.db.RestoreMerchant(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, fmt.Errorf("failed to restore Merchants: %w", err)
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) DeleteMerchantPermanent(Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantPermanently(r.ctx, int32(Merchant_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete Merchant: %w", err)
	}

	return true, nil
}

func (r *merchantRepository) RestoreAllMerchant() (bool, error) {
	err := r.db.RestoreAllMerchants(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all Merchants: %w", err)
	}
	return true, nil
}

func (r *merchantRepository) DeleteAllMerchantPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all Merchants permanently: %w", err)
	}
	return true, nil
}
