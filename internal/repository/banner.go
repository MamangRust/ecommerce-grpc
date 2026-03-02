package repository

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/banner_errors"
	"ecommerce/pkg/utils"
)

type bannerRepository struct {
	db *db.Queries
}

func NewBannerRepository(db *db.Queries) *bannerRepository {
	return &bannerRepository{
		db: db,
	}
}

func (r *bannerRepository) FindAllBanners(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetBannersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetBanners(ctx, reqDb)

	if err != nil {
		return nil, banner_errors.ErrFindAllBanners
	}

	return res, nil
}

func (r *bannerRepository) FindByActive(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetBannersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetBannersActive(ctx, reqDb)

	if err != nil {
		return nil, banner_errors.ErrFindActiveBanners
	}

	return res, nil
}

func (r *bannerRepository) FindByTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetBannersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetBannersTrashed(ctx, reqDb)

	if err != nil {
		return nil, banner_errors.ErrFindTrashedBanners
	}

	return res, nil
}

func (r *bannerRepository) FindById(ctx context.Context, user_id int) (*db.GetBannerRow, error) {
	res, err := r.db.GetBanner(ctx, int32(user_id))

	if err != nil {
		return nil, banner_errors.ErrBannerNotFound
	}

	return res, nil
}

func (r *bannerRepository) CreateBanner(ctx context.Context, request *requests.CreateBannerRequest) (*db.CreateBannerRow, error) {
	startDate, err := utils.ParseDate(request.StartDate)
	if err != nil {
		return nil, banner_errors.ErrBannerStartDate
	}

	endDate, err := utils.ParseDate(request.EndDate)
	if err != nil {
		return nil, banner_errors.ErrBannerEndDate
	}

	startTime, err := utils.ParseTime(request.StartTime)
	if err != nil {
		return nil, banner_errors.ErrBannerStartTime
	}

	endTime, err := utils.ParseTime(request.EndTime)
	if err != nil {
		return nil, banner_errors.ErrBannerEndTime
	}

	req := db.CreateBannerParams{
		Name:      request.Name,
		StartDate: startDate,
		EndDate:   endDate,
		StartTime: startTime,
		EndTime:   endTime,
		IsActive:  &request.IsActive,
	}

	result, err := r.db.CreateBanner(ctx, req)
	if err != nil {
		return nil, banner_errors.ErrCreateBanner
	}

	return result, nil
}

func (r *bannerRepository) UpdateBanner(ctx context.Context, request *requests.UpdateBannerRequest) (*db.UpdateBannerRow, error) {
	startDate, err := utils.ParseDate(request.StartDate)
	if err != nil {
		return nil, banner_errors.ErrBannerStartDate
	}

	endDate, err := utils.ParseDate(request.EndDate)
	if err != nil {
		return nil, banner_errors.ErrBannerEndDate
	}

	startTime, err := utils.ParseTime(request.StartTime)
	if err != nil {
		return nil, banner_errors.ErrBannerStartTime
	}

	endTime, err := utils.ParseTime(request.EndTime)
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
		IsActive:  &request.IsActive,
	}

	result, err := r.db.UpdateBanner(ctx, req)
	if err != nil {
		return nil, banner_errors.ErrUpdateBanner
	}

	return result, nil
}

func (r *bannerRepository) TrashedBanner(ctx context.Context, Banner_id int) (*db.Banner, error) {
	res, err := r.db.TrashBanner(ctx, int32(Banner_id))

	if err != nil {
		return nil, banner_errors.ErrTrashedBanner
	}

	return res, nil
}

func (r *bannerRepository) RestoreBanner(ctx context.Context, Banner_id int) (*db.Banner, error) {
	res, err := r.db.RestoreBanner(ctx, int32(Banner_id))

	if err != nil {
		return nil, banner_errors.ErrRestoreBanner
	}

	return res, nil
}

func (r *bannerRepository) DeleteBannerPermanent(ctx context.Context, Banner_id int) (bool, error) {
	err := r.db.DeleteBannerPermanently(ctx, int32(Banner_id))

	if err != nil {
		return false, banner_errors.ErrDeleteBannerPermanent
	}

	return true, nil
}

func (r *bannerRepository) RestoreAllBanner(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllBanners(ctx)

	if err != nil {
		return false, banner_errors.ErrRestoreAllBanners
	}
	return true, nil
}

func (r *bannerRepository) DeleteAllBannerPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentBanners(ctx)

	if err != nil {
		return false, banner_errors.ErrDeleteAllBanners
	}
	return true, nil
}
