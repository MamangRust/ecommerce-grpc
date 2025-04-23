package repository

import (
	"context"
	"database/sql"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"fmt"
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
		return nil, nil, fmt.Errorf("failed to fetch Banners: invalid pagination (page %d, size %d) or search query '%s'", req.Page, req.PageSize, req.Search)
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
		return nil, nil, fmt.Errorf("failed to fetch Banners active: invalid pagination (page %d, size %d) or search query '%s'", req.Page, req.PageSize, req.Search)
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
		return nil, nil, fmt.Errorf("failed to fetch Banners trashed: invalid pagination (page %d, size %d) or search query '%s'", req.Page, req.PageSize, req.Search)
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
		return nil, fmt.Errorf("failed to find Banner: %w", err)
	}

	return r.mapping.ToBannerRecord(res), nil
}

func (r *bannerRepository) CreateBanner(request *requests.CreateBannerRequest) (*record.BannerRecord, error) {
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format: %w", err)
	}

	startTime, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time format: %w", err)
	}

	endTime, err := time.Parse("15:04", request.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end_time format: %w", err)
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
		return nil, fmt.Errorf("failed to create banner: %w", err)
	}

	return r.mapping.ToBannerRecord(result), nil
}

func (r *bannerRepository) UpdateBanner(request *requests.UpdateBannerRequest) (*record.BannerRecord, error) {
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format: %w", err)
	}

	startTime, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time format: %w", err)
	}

	endTime, err := time.Parse("15:04", request.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end_time format: %w", err)
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
		return nil, fmt.Errorf("failed to update banner: %w", err)
	}

	return r.mapping.ToBannerRecord(result), nil
}

func (r *bannerRepository) TrashedBanner(Banner_id int) (*record.BannerRecord, error) {
	res, err := r.db.TrashBanner(r.ctx, int32(Banner_id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash Banner: %w", err)
	}

	return r.mapping.ToBannerRecord(res), nil
}

func (r *bannerRepository) RestoreBanner(Banner_id int) (*record.BannerRecord, error) {
	res, err := r.db.RestoreBanner(r.ctx, int32(Banner_id))

	if err != nil {
		return nil, fmt.Errorf("failed to restore Banners: %w", err)
	}

	return r.mapping.ToBannerRecord(res), nil
}

func (r *bannerRepository) DeleteBannerPermanent(Banner_id int) (bool, error) {
	err := r.db.DeleteBannerPermanently(r.ctx, int32(Banner_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete Banner: %w", err)
	}

	return true, nil
}

func (r *bannerRepository) RestoreAllBanner() (bool, error) {
	err := r.db.RestoreAllBanners(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all Banners: %w", err)
	}
	return true, nil
}

func (r *bannerRepository) DeleteAllBannerPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentBanners(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all Banners permanently: %w", err)
	}
	return true, nil
}
