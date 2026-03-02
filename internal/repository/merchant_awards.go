package repository

import (
	"context"
	"database/sql"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
	merchantaward_errors "ecommerce/pkg/errors/merchant_award"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type merchantAwardRepository struct {
	db *db.Queries
}

func NewMerchantAwardRepository(
	db *db.Queries,
) *merchantAwardRepository {
	return &merchantAwardRepository{
		db: db,
	}
}

func (r *merchantAwardRepository) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantCertificationsAndAwardsParams{
		Column1: &req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantCertificationsAndAwards(ctx, reqDb)

	if err != nil {
		return nil, merchantaward_errors.ErrFindAllMerchantAwards
	}

	return res, nil
}

func (r *merchantAwardRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantCertificationsAndAwardsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantCertificationsAndAwardsActive(ctx, reqDb)

	if err != nil {
		return nil, merchantaward_errors.ErrFindByActiveMerchantAwards
	}

	return res, nil
}

func (r *merchantAwardRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantCertificationsAndAwardsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantCertificationsAndAwardsTrashed(ctx, reqDb)

	if err != nil {
		return nil, merchantaward_errors.ErrFindByTrashedMerchantAwards
	}

	return res, nil
}

func (r *merchantAwardRepository) FindById(ctx context.Context, user_id int) (*db.GetMerchantCertificationOrAwardRow, error) {
	res, err := r.db.GetMerchantCertificationOrAward(ctx, int32(user_id))

	if err != nil {
		return nil, merchantaward_errors.ErrFindByIdMerchantAward
	}

	return res, nil
}

func (r *merchantAwardRepository) CreateMerchantAward(
	ctx context.Context,
	request *requests.CreateMerchantCertificationOrAwardRequest,
) (*db.CreateMerchantCertificationOrAwardRow, error) {

	req := db.CreateMerchantCertificationOrAwardParams{
		MerchantID: int32(request.MerchantID),
		Title:      request.Title,

		Description:    stringPtr(request.Description),
		IssuedBy:       stringPtr(request.IssuedBy),
		CertificateUrl: stringPtr(request.CertificateUrl),

		IssueDate:  parseDateToPgDate(request.IssueDate),
		ExpiryDate: parseDateToPgDate(request.ExpiryDate),
	}

	award, err := r.db.CreateMerchantCertificationOrAward(ctx, req)
	if err != nil {
		return nil, merchantaward_errors.ErrCreateMerchantAward
	}

	return award, nil
}

func (r *merchantAwardRepository) UpdateMerchantAward(ctx context.Context, request *requests.UpdateMerchantCertificationOrAwardRequest) (*db.UpdateMerchantCertificationOrAwardRow, error) {
	req := db.UpdateMerchantCertificationOrAwardParams{
		MerchantCertificationID: int32(*request.MerchantCertificationID),
		Title:                   request.Title,
		Description:             stringPtr(request.Description),
		IssuedBy:                stringPtr(request.IssuedBy),
		CertificateUrl:          stringPtr(request.CertificateUrl),
		IssueDate:               parseDateToPgDate(request.IssueDate),
		ExpiryDate:              parseDateToPgDate(request.ExpiryDate),
	}

	res, err := r.db.UpdateMerchantCertificationOrAward(ctx, req)
	if err != nil {
		return nil, merchantaward_errors.ErrUpdateMerchantAward
	}

	return res, nil
}

func (r *merchantAwardRepository) TrashedMerchantAward(ctx context.Context, merchant_id int) (*db.MerchantCertificationsAndAward, error) {
	res, err := r.db.TrashMerchantCertificationOrAward(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantaward_errors.ErrTrashedMerchantAward
	}

	return res, nil
}

func (r *merchantAwardRepository) RestoreMerchantAward(ctx context.Context, merchant_id int) (*db.MerchantCertificationsAndAward, error) {
	res, err := r.db.RestoreMerchantCertificationOrAward(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantaward_errors.ErrRestoreMerchantAward
	}

	return res, nil
}

func (r *merchantAwardRepository) DeleteMerchantPermanent(ctx context.Context, Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantCertificationOrAwardPermanently(ctx, int32(Merchant_id))

	if err != nil {
		return false, merchantaward_errors.ErrDeleteMerchantAwardPermanent
	}

	return true, nil
}

func (r *merchantAwardRepository) RestoreAllMerchantAward(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchantCertificationsAndAwards(ctx)

	if err != nil {
		return false, merchantaward_errors.ErrRestoreAllMerchantAwards
	}
	return true, nil
}

func (r *merchantAwardRepository) DeleteAllMerchantAwardPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchantCertificationsAndAwards(ctx)

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

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func parseDateToPgDate(dateStr string) pgtype.Date {
	if dateStr == "" {
		return pgtype.Date{Valid: false}
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return pgtype.Date{Valid: false}
	}

	return pgtype.Date{
		Time:  t,
		Valid: true,
	}
}
