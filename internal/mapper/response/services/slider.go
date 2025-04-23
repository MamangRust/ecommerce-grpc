package response_service

import (
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/response"
)

type sliderResponseMapper struct {
}

func NewSliderResponseMapper() *sliderResponseMapper {
	return &sliderResponseMapper{}
}

func (s *sliderResponseMapper) ToSliderResponse(slider *record.SliderRecord) *response.SliderResponse {
	return &response.SliderResponse{
		ID:        slider.ID,
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt,
		UpdatedAt: slider.UpdatedAt,
	}
}

func (s *sliderResponseMapper) ToSlidersResponse(sliders []*record.SliderRecord) []*response.SliderResponse {
	var responses []*response.SliderResponse

	for _, slider := range sliders {
		responses = append(responses, s.ToSliderResponse(slider))
	}

	return responses
}

func (s *sliderResponseMapper) ToSliderResponseDeleteAt(slider *record.SliderRecord) *response.SliderResponseDeleteAt {
	return &response.SliderResponseDeleteAt{
		ID:        slider.ID,
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt,
		UpdatedAt: slider.UpdatedAt,
		DeletedAt: slider.DeletedAt,
	}
}

func (s *sliderResponseMapper) ToSlidersResponseDeleteAt(sliders []*record.SliderRecord) []*response.SliderResponseDeleteAt {
	var responses []*response.SliderResponseDeleteAt

	for _, slider := range sliders {
		responses = append(responses, s.ToSliderResponseDeleteAt(slider))
	}

	return responses
}
