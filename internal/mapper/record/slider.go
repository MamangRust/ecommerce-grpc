package recordmapper

import (
	"ecommerce/internal/domain/record"
	db "ecommerce/pkg/database/schema"
)

type sliderRecordMapper struct {
}

func NewSliderRecordMapper() *sliderRecordMapper {
	return &sliderRecordMapper{}
}

func (s *sliderRecordMapper) ToSliderRecord(slider *db.Slider) *record.SliderRecord {
	var deletedAt *string
	if slider.DeletedAt.Valid {
		deletedAtStr := slider.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.SliderRecord{
		ID:        int(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *sliderRecordMapper) ToSliderRecordPagination(slider *db.GetSlidersRow) *record.SliderRecord {
	var deletedAt *string
	if slider.DeletedAt.Valid {
		deletedAtStr := slider.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.SliderRecord{
		ID:        int(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *sliderRecordMapper) ToSlidersRecordPagination(sliders []*db.GetSlidersRow) []*record.SliderRecord {
	var result []*record.SliderRecord

	for _, slider := range sliders {
		result = append(result, s.ToSliderRecordPagination(slider))
	}

	return result
}

func (s *sliderRecordMapper) ToSliderRecordActivePagination(slider *db.GetSlidersActiveRow) *record.SliderRecord {
	var deletedAt *string
	if slider.DeletedAt.Valid {
		deletedAtStr := slider.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.SliderRecord{
		ID:        int(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *sliderRecordMapper) ToSlidersRecordActivePagination(sliders []*db.GetSlidersActiveRow) []*record.SliderRecord {
	var result []*record.SliderRecord

	for _, slider := range sliders {
		result = append(result, s.ToSliderRecordActivePagination(slider))
	}

	return result
}

func (s *sliderRecordMapper) ToSliderRecordTrashedPagination(slider *db.GetSlidersTrashedRow) *record.SliderRecord {
	var deletedAt *string
	if slider.DeletedAt.Valid {
		deletedAtStr := slider.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.SliderRecord{
		ID:        int(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *sliderRecordMapper) ToSlidersRecordTrashedPagination(sliders []*db.GetSlidersTrashedRow) []*record.SliderRecord {
	var result []*record.SliderRecord

	for _, slider := range sliders {
		result = append(result, s.ToSliderRecordTrashedPagination(slider))
	}

	return result
}
