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

type merchantDetailRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantDetailMapping
}

func NewMerchantDetailRepository(
	db *db.Queries, ctx context.Context, mapping recordmapper.MerchantDetailMapping,
) *merchantDetailRepository {
	return &merchantDetailRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantDetailRepository) FindAllMerchants(req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantDetailsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantDetails(r.ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch merchants: invalid pagination (page %d, size %d) or search query '%s'", req.Page, req.PageSize, req.Search)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantDetailsRecordPagination(res), &totalCount, nil
}

func (r *merchantDetailRepository) FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantDetailsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantDetailsActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch merchants active: invalid pagination (page %d, size %d) or search query '%s'", req.Page, req.PageSize, req.Search)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantDetailsRecordActivePagination(res), &totalCount, nil
}

func (r *merchantDetailRepository) FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantDetailsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantDetailsTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch merchants trashed: invalid pagination (page %d, size %d) or search query '%s'", req.Page, req.PageSize, req.Search)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantDetailsRecordTrashedPagination(res), &totalCount, nil
}

func (r *merchantDetailRepository) FindById(user_id int) (*record.MerchantDetailRecord, error) {
	res, err := r.db.GetMerchantDetail(r.ctx, int32(user_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find merchant: %w", err)
	}

	return r.mapping.ToMerchantDetailRelationRecord(res), nil
}

func (r *merchantDetailRepository) FindByIdTrashed(user_id int) (*record.MerchantDetailRecord, error) {
	res, err := r.db.GetMerchantDetailTrashed(r.ctx, int32(user_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find merchant: %w", err)
	}

	return r.mapping.ToMerchantDetailRecord(res), nil
}

func (r *merchantDetailRepository) CreateMerchantDetail(request *requests.CreateMerchantDetailRequest) (*record.MerchantDetailRecord, error) {
	req := db.CreateMerchantDetailParams{
		MerchantID:       int32(request.MerchantID),
		DisplayName:      sql.NullString{String: request.DisplayName, Valid: true},
		CoverImageUrl:    sql.NullString{String: request.CoverImageUrl, Valid: true},
		LogoUrl:          sql.NullString{String: request.LogoUrl, Valid: true},
		ShortDescription: sql.NullString{String: request.ShortDescription, Valid: true},
		WebsiteUrl:       sql.NullString{String: request.WebsiteUrl, Valid: request.WebsiteUrl != ""},
	}

	merchant, err := r.db.CreateMerchantDetail(r.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create merchant: %w", err)
	}

	return r.mapping.ToMerchantDetailRecord(merchant), nil
}

func (r *merchantDetailRepository) UpdateMerchantDetail(request *requests.UpdateMerchantDetailRequest) (*record.MerchantDetailRecord, error) {
	req := db.UpdateMerchantDetailParams{
		MerchantDetailID: int32(*request.MerchantDetailID),
		DisplayName:      sql.NullString{String: request.DisplayName, Valid: true},
		CoverImageUrl:    sql.NullString{String: request.CoverImageUrl, Valid: true},
		LogoUrl:          sql.NullString{String: request.LogoUrl, Valid: true},
		ShortDescription: sql.NullString{String: request.ShortDescription, Valid: true},
		WebsiteUrl:       sql.NullString{String: request.WebsiteUrl, Valid: request.WebsiteUrl != ""},
	}

	res, err := r.db.UpdateMerchantDetail(r.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update merchant: %w", err)
	}

	return r.mapping.ToMerchantDetailRecord(res), nil
}

func (r *merchantDetailRepository) TrashedMerchantDetail(merchant_id int) (*record.MerchantDetailRecord, error) {
	res, err := r.db.TrashMerchantDetail(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash Merchant: %w", err)
	}

	return r.mapping.ToMerchantDetailRecord(res), nil
}

func (r *merchantDetailRepository) RestoreMerchantDetail(merchant_id int) (*record.MerchantDetailRecord, error) {
	res, err := r.db.RestoreMerchantDetail(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, fmt.Errorf("failed to restore Merchants: %w", err)
	}

	return r.mapping.ToMerchantDetailRecord(res), nil
}

func (r *merchantDetailRepository) DeleteMerchantDetailPermanent(Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantDetailPermanently(r.ctx, int32(Merchant_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete Merchant: %w", err)
	}

	return true, nil
}

func (r *merchantDetailRepository) RestoreAllMerchantDetail() (bool, error) {
	err := r.db.RestoreAllMerchantDetails(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all Merchants: %w", err)
	}
	return true, nil
}

func (r *merchantDetailRepository) DeleteAllMerchantDetailPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentMerchantDetails(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all Merchants permanently: %w", err)
	}
	return true, nil
}
