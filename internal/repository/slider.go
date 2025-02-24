package repository

import (
	"context"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"fmt"
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

func (r *sliderRepository) FindAllSlider(search string, page, pageSize int) ([]*record.SliderRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetSlidersParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSliders(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find shipping address: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToSlidersRecordPagination(res), totalCount, nil
}

func (r *sliderRepository) FindByActive(search string, page, pageSize int) ([]*record.SliderRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetSlidersActiveParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSlidersActive(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find shipping address: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToSlidersRecordActivePagination(res), totalCount, nil
}

func (r *sliderRepository) FindByTrashed(search string, page, pageSize int) ([]*record.SliderRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetSlidersTrashedParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSlidersTrashed(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find shipping address: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToSlidersRecordTrashedPagination(res), totalCount, nil
}

func (r *sliderRepository) CreateSlider(request *requests.CreateSliderRequest) (*record.SliderRecord, error) {
	req := db.CreateSliderParams{
		Name:  request.Nama,
		Image: request.FilePath,
	}

	slider, err := r.db.CreateSlider(r.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create slider: %w", err)
	}

	return r.mapping.ToSliderRecord(slider), nil
}

func (r *sliderRepository) UpdateSlider(request *requests.UpdateSliderRequest) (*record.SliderRecord, error) {
	req := db.UpdateSliderParams{
		SliderID: int32(request.ID),
		Name:     request.Nama,
		Image:    request.FilePath,
	}

	res, err := r.db.UpdateSlider(r.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update slider: %w", err)
	}

	return r.mapping.ToSliderRecord(res), nil
}

func (r *sliderRepository) TrashSlider(slider_id int) (*record.SliderRecord, error) {
	res, err := r.db.TrashSlider(r.ctx, int32(slider_id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash shipping address: %w", err)
	}

	return r.mapping.ToSliderRecord(res), nil
}

func (r *sliderRepository) RestoreSlider(slider_id int) (*record.SliderRecord, error) {
	res, err := r.db.RestoreSlider(r.ctx, int32(slider_id))

	if err != nil {
		return nil, fmt.Errorf("failed to shipping address: %w", err)
	}

	return r.mapping.ToSliderRecord(res), nil
}

func (r *sliderRepository) DeleteSliderPermanently(slider_id int) (bool, error) {
	err := r.db.DeleteSliderPermanently(r.ctx, int32(slider_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete shipping address: %w", err)
	}

	return true, nil
}

func (r *sliderRepository) RestoreAllSlider() (bool, error) {
	err := r.db.RestoreAllSliders(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all shipping address: %w", err)
	}
	return true, nil
}

func (r *sliderRepository) DeleteAllPermanentSlider() (bool, error) {
	err := r.db.DeleteAllPermanentSliders(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all shipping address permanently: %w", err)
	}
	return true, nil
}
