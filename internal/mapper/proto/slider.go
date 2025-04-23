package protomapper

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type sliderProtoMapper struct {
}

func NewSliderProtoMapper() *sliderProtoMapper {
	return &sliderProtoMapper{}
}

func (s *sliderProtoMapper) ToProtoResponseSlider(status string, message string, pbResponse *response.SliderResponse) *pb.ApiResponseSlider {
	return &pb.ApiResponseSlider{
		Status:  status,
		Message: message,
		Data:    s.mapResponseSlider(pbResponse),
	}
}

func (s *sliderProtoMapper) ToProtoResponseSliderDeleteAt(status string, message string, pbResponse *response.SliderResponseDeleteAt) *pb.ApiResponseSliderDeleteAt {
	return &pb.ApiResponseSliderDeleteAt{
		Status:  status,
		Message: message,
		Data:    s.mapResponseSliderDeleteAt(pbResponse),
	}
}

func (s *sliderProtoMapper) ToProtoResponsesSlider(status string, message string, pbResponse []*response.SliderResponse) *pb.ApiResponsesSlider {
	return &pb.ApiResponsesSlider{
		Status:  status,
		Message: message,
		Data:    s.mapResponsesSlider(pbResponse),
	}
}

func (s *sliderProtoMapper) ToProtoResponseSliderDelete(status string, message string) *pb.ApiResponseSliderDelete {
	return &pb.ApiResponseSliderDelete{
		Status:  status,
		Message: message,
	}
}

func (s *sliderProtoMapper) ToProtoResponsePaginationSliderDeleteAt(pagination *pb.PaginationMeta, status string, message string, sliders []*response.SliderResponseDeleteAt) *pb.ApiResponsePaginationSliderDeleteAt {
	return &pb.ApiResponsePaginationSliderDeleteAt{
		Status:     status,
		Message:    message,
		Data:       s.mapResponsesSliderDeleteAt(sliders),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (c *sliderProtoMapper) ToProtoResponseSliderAll(status string, message string) *pb.ApiResponseSliderAll {
	return &pb.ApiResponseSliderAll{
		Status:  status,
		Message: message,
	}
}

func (s *sliderProtoMapper) ToProtoResponsePaginationSlider(pagination *pb.PaginationMeta, status string, message string, sliders []*response.SliderResponse) *pb.ApiResponsePaginationSlider {
	return &pb.ApiResponsePaginationSlider{
		Status:     status,
		Message:    message,
		Data:       s.mapResponsesSlider(sliders),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (s *sliderProtoMapper) mapResponseSlider(slider *response.SliderResponse) *pb.SliderResponse {
	return &pb.SliderResponse{
		Id:        int32(slider.ID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt,
		UpdatedAt: slider.UpdatedAt,
	}
}

func (s *sliderProtoMapper) mapResponsesSlider(sliders []*response.SliderResponse) []*pb.SliderResponse {
	var mappedSliders []*pb.SliderResponse

	for _, slider := range sliders {
		mappedSliders = append(mappedSliders, s.mapResponseSlider(slider))
	}

	return mappedSliders
}

func (s *sliderProtoMapper) mapResponseSliderDeleteAt(slider *response.SliderResponseDeleteAt) *pb.SliderResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue

	if slider.DeletedAt != nil {
		deletedAt = wrapperspb.String(*slider.DeletedAt)
	}

	return &pb.SliderResponseDeleteAt{
		Id:        int32(slider.ID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt,
		UpdatedAt: slider.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (s *sliderProtoMapper) mapResponsesSliderDeleteAt(sliders []*response.SliderResponseDeleteAt) []*pb.SliderResponseDeleteAt {
	var mappedSliders []*pb.SliderResponseDeleteAt

	for _, slider := range sliders {
		mappedSliders = append(mappedSliders, s.mapResponseSliderDeleteAt(slider))
	}

	return mappedSliders
}
