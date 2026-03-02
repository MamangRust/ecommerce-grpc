package repository

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/slider_errors"
)

type sliderRepository struct {
	db *db.Queries
}

func NewSliderRepository(
	db *db.Queries,
) *sliderRepository {
	return &sliderRepository{
		db: db,
	}
}

func (r *sliderRepository) FindAllSlider(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSlidersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSliders(ctx, reqDb)

	if err != nil {
		return nil, slider_errors.ErrFindAllSliders
	}

	return res, nil
}

func (r *sliderRepository) FindByActive(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSlidersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSlidersActive(ctx, reqDb)

	if err != nil {
		return nil, slider_errors.ErrFindActiveSliders
	}

	return res, nil
}

func (r *sliderRepository) FindByTrashed(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSlidersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSlidersTrashed(ctx, reqDb)

	if err != nil {
		return nil, slider_errors.ErrFindTrashedSliders
	}

	return res, nil
}

func (r *sliderRepository) CreateSlider(ctx context.Context, request *requests.CreateSliderRequest) (*db.CreateSliderRow, error) {
	req := db.CreateSliderParams{
		Name:  request.Nama,
		Image: request.FilePath,
	}

	slider, err := r.db.CreateSlider(ctx, req)

	if err != nil {
		return nil, slider_errors.ErrCreateSlider
	}

	return slider, nil
}

func (r *sliderRepository) UpdateSlider(ctx context.Context, request *requests.UpdateSliderRequest) (*db.UpdateSliderRow, error) {
	req := db.UpdateSliderParams{
		SliderID: int32(*request.ID),
		Name:     request.Nama,
		Image:    request.FilePath,
	}

	res, err := r.db.UpdateSlider(ctx, req)

	if err != nil {
		return nil, slider_errors.ErrUpdateSlider
	}

	return res, nil
}

func (r *sliderRepository) TrashSlider(ctx context.Context, slider_id int) (*db.Slider, error) {
	res, err := r.db.TrashSlider(ctx, int32(slider_id))

	if err != nil {
		return nil, slider_errors.ErrTrashSlider
	}

	return res, nil
}

func (r *sliderRepository) RestoreSlider(ctx context.Context, slider_id int) (*db.Slider, error) {
	res, err := r.db.RestoreSlider(ctx, int32(slider_id))

	if err != nil {
		return nil, slider_errors.ErrRestoreSlider
	}

	return res, nil
}

func (r *sliderRepository) DeleteSliderPermanently(ctx context.Context, slider_id int) (bool, error) {
	err := r.db.DeleteSliderPermanently(ctx, int32(slider_id))

	if err != nil {
		return false, slider_errors.ErrDeletePermanentSlider
	}

	return true, nil
}

func (r *sliderRepository) RestoreAllSlider(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllSliders(ctx)

	if err != nil {
		return false, slider_errors.ErrRestoreAllSlider
	}
	return true, nil
}

func (r *sliderRepository) DeleteAllPermanentSlider(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentSliders(ctx)

	if err != nil {
		return false, slider_errors.ErrDeleteAllPermanentSlider
	}
	return true, nil
}
