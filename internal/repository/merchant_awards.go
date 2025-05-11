package repository

import (
	"context"
	"database/sql"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	merchantaward_errors "ecommerce/pkg/errors/merchant_award"
	"time"
)

type merchantAwardRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantAwardMapping
}

func NewMerchantAwardRepository(
	db *db.Queries,
	ctx context.Context,
	mapping recordmapper.MerchantAwardMapping,
) *merchantAwardRepository {
	return &merchantAwardRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantAwardRepository) FindAllMerchants(req *requests.FindAllMerchant) ([]*record.MerchantAwardRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantCertificationsAndAwardsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantCertificationsAndAwards(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchantaward_errors.ErrFindAllMerchantAwards
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantAwardsRecordPagination(res), &totalCount, nil
}

func (r *merchantAwardRepository) FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantAwardRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantCertificationsAndAwardsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantCertificationsAndAwardsActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchantaward_errors.ErrFindByActiveMerchantAwards
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantAwardsRecordActivePagination(res), &totalCount, nil
}

func (r *merchantAwardRepository) FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantAwardRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantCertificationsAndAwardsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantCertificationsAndAwardsTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchantaward_errors.ErrFindByTrashedMerchantAwards
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantAwardsRecordTrashedPagination(res), &totalCount, nil
}

func (r *merchantAwardRepository) FindById(user_id int) (*record.MerchantAwardRecord, error) {
	res, err := r.db.GetMerchantCertificationOrAward(r.ctx, int32(user_id))

	if err != nil {
		return nil, merchantaward_errors.ErrFindByIdMerchantAward
	}

	return r.mapping.ToMerchantAwardRecord(res), nil
}

func (r *merchantAwardRepository) CreateMerchantAward(request *requests.CreateMerchantCertificationOrAwardRequest) (*record.MerchantAwardRecord, error) {
	req := db.CreateMerchantCertificationOrAwardParams{
		MerchantID:     int32(request.MerchantID),
		Title:          request.Title,
		Description:    sql.NullString{String: request.Description, Valid: request.Description != ""},
		IssuedBy:       sql.NullString{String: request.IssuedBy, Valid: request.IssuedBy != ""},
		IssueDate:      parseDateToNullTime(request.IssueDate),
		ExpiryDate:     parseDateToNullTime(request.ExpiryDate),
		CertificateUrl: sql.NullString{String: request.CertificateUrl, Valid: request.CertificateUrl != ""},
	}

	award, err := r.db.CreateMerchantCertificationOrAward(r.ctx, req)
	if err != nil {
		return nil, merchantaward_errors.ErrCreateMerchantAward
	}

	return r.mapping.ToMerchantAwardRecord(award), nil
}

func (r *merchantAwardRepository) UpdateMerchantAward(request *requests.UpdateMerchantCertificationOrAwardRequest) (*record.MerchantAwardRecord, error) {
	req := db.UpdateMerchantCertificationOrAwardParams{
		MerchantCertificationID: int32(*request.MerchantCertificationID),
		Title:                   request.Title,
		Description:             sql.NullString{String: request.Description, Valid: request.Description != ""},
		IssuedBy:                sql.NullString{String: request.IssuedBy, Valid: request.IssuedBy != ""},
		IssueDate:               parseDateToNullTime(request.IssueDate),
		ExpiryDate:              parseDateToNullTime(request.ExpiryDate),
		CertificateUrl:          sql.NullString{String: request.CertificateUrl, Valid: request.CertificateUrl != ""},
	}

	res, err := r.db.UpdateMerchantCertificationOrAward(r.ctx, req)
	if err != nil {
		return nil, merchantaward_errors.ErrUpdateMerchantAward
	}

	return r.mapping.ToMerchantAwardRecord(res), nil
}

func (r *merchantAwardRepository) TrashedMerchantAward(merchant_id int) (*record.MerchantAwardRecord, error) {
	res, err := r.db.TrashMerchantCertificationOrAward(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantaward_errors.ErrTrashedMerchantAward
	}

	return r.mapping.ToMerchantAwardRecord(res), nil
}

func (r *merchantAwardRepository) RestoreMerchantAward(merchant_id int) (*record.MerchantAwardRecord, error) {
	res, err := r.db.RestoreMerchantCertificationOrAward(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantaward_errors.ErrRestoreMerchantAward
	}

	return r.mapping.ToMerchantAwardRecord(res), nil
}

func (r *merchantAwardRepository) DeleteMerchantPermanent(Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantCertificationOrAwardPermanently(r.ctx, int32(Merchant_id))

	if err != nil {
		return false, merchantaward_errors.ErrDeleteMerchantAwardPermanent
	}

	return true, nil
}

func (r *merchantAwardRepository) RestoreAllMerchantAward() (bool, error) {
	err := r.db.RestoreAllMerchantCertificationsAndAwards(r.ctx)

	if err != nil {
		return false, merchantaward_errors.ErrRestoreAllMerchantAwards
	}
	return true, nil
}

func (r *merchantAwardRepository) DeleteAllMerchantAwardPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentMerchantCertificationsAndAwards(r.ctx)

	if err != nil {
		return false, merchantaward_errors.ErrDeleteAllMerchantAwardsPermanent
	}
	return true, nil
}

func parseDateToNullTime(dateStr string) sql.NullTime {
	if dateStr == "" {
		return sql.NullTime{Valid: false}
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return sql.NullTime{Valid: false}
	}

	return sql.NullTime{Time: t, Valid: true}
}
