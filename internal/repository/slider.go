package repository

import (
	"context"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/slider_errors"
)

type sliderRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.SliderMapping
}

func NewSliderRepository(
	db *db.Queries,
	ctx context.Context,
	mapping recordmapper.SliderMapping,
) *sliderRepository {
	return &sliderRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *sliderRepository) FindAllSlider(req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSlidersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSliders(r.ctx, reqDb)

	if err != nil {
		return nil, nil, slider_errors.ErrFindAllSliders
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToSlidersRecordPagination(res), &totalCount, nil
}

func (r *sliderRepository) FindByActive(req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSlidersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSlidersActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, slider_errors.ErrFindActiveSliders
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToSlidersRecordActivePagination(res), &totalCount, nil
}

func (r *sliderRepository) FindByTrashed(req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSlidersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSlidersTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, slider_errors.ErrFindTrashedSliders
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToSlidersRecordTrashedPagination(res), &totalCount, nil
}

func (r *sliderRepository) CreateSlider(request *requests.CreateSliderRequest) (*record.SliderRecord, error) {
	req := db.CreateSliderParams{
		Name:  request.Nama,
		Image: request.FilePath,
	}

	slider, err := r.db.CreateSlider(r.ctx, req)

	if err != nil {
		return nil, slider_errors.ErrCreateSlider
	}

	return r.mapping.ToSliderRecord(slider), nil
}

func (r *sliderRepository) UpdateSlider(request *requests.UpdateSliderRequest) (*record.SliderRecord, error) {
	req := db.UpdateSliderParams{
		SliderID: int32(*request.ID),
		Name:     request.Nama,
		Image:    request.FilePath,
	}

	res, err := r.db.UpdateSlider(r.ctx, req)

	if err != nil {
		return nil, slider_errors.ErrUpdateSlider
	}

	return r.mapping.ToSliderRecord(res), nil
}

func (r *sliderRepository) TrashSlider(slider_id int) (*record.SliderRecord, error) {
	res, err := r.db.TrashSlider(r.ctx, int32(slider_id))

	if err != nil {
		return nil, slider_errors.ErrTrashSlider
	}

	return r.mapping.ToSliderRecord(res), nil
}

func (r *sliderRepository) RestoreSlider(slider_id int) (*record.SliderRecord, error) {
	res, err := r.db.RestoreSlider(r.ctx, int32(slider_id))

	if err != nil {
		return nil, slider_errors.ErrRestoreSlider
	}

	return r.mapping.ToSliderRecord(res), nil
}

func (r *sliderRepository) DeleteSliderPermanently(slider_id int) (bool, error) {
	err := r.db.DeleteSliderPermanently(r.ctx, int32(slider_id))

	if err != nil {
		return false, slider_errors.ErrDeletePermanentSlider
	}

	return true, nil
}

func (r *sliderRepository) RestoreAllSlider() (bool, error) {
	err := r.db.RestoreAllSliders(r.ctx)

	if err != nil {
		return false, slider_errors.ErrRestoreAllSlider
	}
	return true, nil
}

func (r *sliderRepository) DeleteAllPermanentSlider() (bool, error) {
	err := r.db.DeleteAllPermanentSliders(r.ctx)

	if err != nil {
		return false, slider_errors.ErrDeleteAllPermanentSlider
	}
	return true, nil
}
