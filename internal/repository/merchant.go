package repository

import (
	"context"
	"database/sql"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"errors"
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

func (r *merchantRepository) FindAllMerchants(search string, page, pageSize int) ([]*record.MerchantRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetMerchantsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchants(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find Merchants: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsRecordPagination(res), totalCount, nil
}

func (r *merchantRepository) FindByActive(search string, page, pageSize int) ([]*record.MerchantRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetMerchantsActiveParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsActive(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find Merchants: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsRecordActivePagination(res), totalCount, nil
}

func (r *merchantRepository) FindByTrashed(search string, page, pageSize int) ([]*record.MerchantRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetMerchantsTrashedParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsTrashed(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find Merchants: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsRecordTrashedPagination(res), totalCount, nil
}

func (r *merchantRepository) FindById(user_id int) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByID(r.ctx, int32(user_id))

	if err != nil {
		fmt.Printf("Error fetching user: %v\n", err)

		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) CreateMerchant(request *requests.CreateMerchantRequest) (*record.MerchantRecord, error) {
	req := db.CreateMerchantParams{
		UserID:       int32(request.UserID),
		Name:         request.Name,
		Description:  sql.NullString{String: request.Description, Valid: request.Description != ""},
		Address:      sql.NullString{String: request.Address, Valid: request.Address != ""},
		ContactEmail: sql.NullString{String: request.ContactEmail, Valid: request.ContactEmail != ""},
		ContactPhone: sql.NullString{String: request.ContactPhone, Valid: request.ContactPhone != ""},
		Status:       "active",
	}

	merchant, err := r.db.CreateMerchant(r.ctx, req)
	if err != nil {
		return nil, errors.New("failed to create Merchant")
	}

	return r.mapping.ToMerchantRecord(merchant), nil
}

func (r *merchantRepository) UpdateMerchant(request *requests.UpdateMerchantRequest) (*record.MerchantRecord, error) {
	req := db.UpdateMerchantParams{
		MerchantID:   int32(request.MerchantID),
		Name:         request.Name,
		Description:  sql.NullString{String: request.Description, Valid: request.Description != ""},
		Address:      sql.NullString{String: request.Address, Valid: request.Address != ""},
		ContactEmail: sql.NullString{String: request.ContactEmail, Valid: request.ContactEmail != ""},
		ContactPhone: sql.NullString{String: request.ContactPhone, Valid: request.ContactPhone != ""},
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
