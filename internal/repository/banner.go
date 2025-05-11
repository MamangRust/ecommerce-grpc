package repository

import (
	"context"
	"database/sql"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/banner_errors"
	"time"
)

type bannerRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.BannerRecordMapping
}

func NewBannerRepository(db *db.Queries, ctx context.Context, mapping recordmapper.BannerRecordMapping) *bannerRepository {
	return &bannerRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *bannerRepository) FindAllBanners(req *requests.FindAllBanner) ([]*record.BannerRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetBannersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetBanners(r.ctx, reqDb)

	if err != nil {
		return nil, nil, banner_errors.ErrFindAllBanners
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToBannersRecordPagination(res), &totalCount, nil
}

func (r *bannerRepository) FindByActive(req *requests.FindAllBanner) ([]*record.BannerRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetBannersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetBannersActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, banner_errors.ErrFindActiveBanners
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToBannersRecordActivePagination(res), &totalCount, nil
}

func (r *bannerRepository) FindByTrashed(req *requests.FindAllBanner) ([]*record.BannerRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetBannersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetBannersTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, banner_errors.ErrFindTrashedBanners
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToBannersRecordTrashedPagination(res), &totalCount, nil
}

func (r *bannerRepository) FindById(user_id int) (*record.BannerRecord, error) {
	res, err := r.db.GetBanner(r.ctx, int32(user_id))

	if err != nil {
		return nil, banner_errors.ErrBannerNotFound
	}

	return r.mapping.ToBannerRecord(res), nil
}

func (r *bannerRepository) CreateBanner(request *requests.CreateBannerRequest) (*record.BannerRecord, error) {
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return nil, banner_errors.ErrBannerStartDate
	}

	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		return nil, banner_errors.ErrBannerEndDate
	}

	startTime, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		return nil, banner_errors.ErrBannerStartTime
	}

	endTime, err := time.Parse("15:04", request.EndTime)
	if err != nil {
		return nil, banner_errors.ErrBannerEndTime
	}

	req := db.CreateBannerParams{
		Name:      request.Name,
		StartDate: startDate,
		EndDate:   endDate,
		StartTime: startTime,
		EndTime:   endTime,
		IsActive:  sql.NullBool{Bool: request.IsActive, Valid: true},
	}

	result, err := r.db.CreateBanner(r.ctx, req)
	if err != nil {
		return nil, banner_errors.ErrCreateBanner
	}

	return r.mapping.ToBannerRecord(result), nil
}

func (r *bannerRepository) UpdateBanner(request *requests.UpdateBannerRequest) (*record.BannerRecord, error) {
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return nil, banner_errors.ErrBannerStartDate
	}

	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		return nil, banner_errors.ErrBannerEndDate
	}

	startTime, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		return nil, banner_errors.ErrBannerStartTime
	}

	endTime, err := time.Parse("15:04", request.EndTime)
	if err != nil {
		return nil, banner_errors.ErrBannerEndTime
	}

	req := db.UpdateBannerParams{
		BannerID:  int32(*request.BannerID),
		Name:      request.Name,
		StartDate: startDate,
		EndDate:   endDate,
		StartTime: startTime,
		EndTime:   endTime,
		IsActive:  sql.NullBool{Bool: request.IsActive, Valid: true},
	}

	result, err := r.db.UpdateBanner(r.ctx, req)
	if err != nil {
		return nil, banner_errors.ErrUpdateBanner
	}

	return r.mapping.ToBannerRecord(result), nil
}

func (r *bannerRepository) TrashedBanner(Banner_id int) (*record.BannerRecord, error) {
	res, err := r.db.TrashBanner(r.ctx, int32(Banner_id))

	if err != nil {
		return nil, banner_errors.ErrTrashedBanner
	}

	return r.mapping.ToBannerRecord(res), nil
}

func (r *bannerRepository) RestoreBanner(Banner_id int) (*record.BannerRecord, error) {
	res, err := r.db.RestoreBanner(r.ctx, int32(Banner_id))

	if err != nil {
		return nil, banner_errors.ErrRestoreBanner
	}

	return r.mapping.ToBannerRecord(res), nil
}

func (r *bannerRepository) DeleteBannerPermanent(Banner_id int) (bool, error) {
	err := r.db.DeleteBannerPermanently(r.ctx, int32(Banner_id))

	if err != nil {
		return false, banner_errors.ErrDeleteBannerPermanent
	}

	return true, nil
}

func (r *bannerRepository) RestoreAllBanner() (bool, error) {
	err := r.db.RestoreAllBanners(r.ctx)

	if err != nil {
		return false, banner_errors.ErrRestoreAllBanners
	}
	return true, nil
}

func (r *bannerRepository) DeleteAllBannerPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentBanners(r.ctx)

	if err != nil {
		return false, banner_errors.ErrDeleteAllBanners
	}
	return true, nil
}
