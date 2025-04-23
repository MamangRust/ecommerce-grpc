package protomapper

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type reviewDetailProtoMapper struct {
}

func NewReviewDetailProtoMapper() *reviewDetailProtoMapper {
	return &reviewDetailProtoMapper{}
}

func (m *reviewDetailProtoMapper) ToProtoResponseReviewDetail(status string, message string, pbResponse *response.ReviewDetailsResponse) *pb.ApiResponseReviewDetail {
	return &pb.ApiResponseReviewDetail{
		Status:  status,
		Message: message,
		Data:    m.mapResponseReviewDetail(pbResponse),
	}
}

func (m *reviewDetailProtoMapper) ToProtoResponseReviewDetailDeleteAt(status string, message string, pbResponse *response.ReviewDetailsResponseDeleteAt) *pb.ApiResponseReviewDetailDeleteAt {
	return &pb.ApiResponseReviewDetailDeleteAt{
		Status:  status,
		Message: message,
		Data:    m.mapResponseReviewDetailDeleteAt(pbResponse),
	}
}

func (m *reviewDetailProtoMapper) ToProtoResponsesReviewDetail(status string, message string, pbResponse []*response.ReviewDetailsResponse) *pb.ApiResponsesReviewDetails {
	return &pb.ApiResponsesReviewDetails{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesReviewDetail(pbResponse),
	}
}

func (m *reviewDetailProtoMapper) ToProtoResponsePaginationReviewDetailDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.ReviewDetailsResponseDeleteAt) *pb.ApiResponsePaginationReviewDetailsDeleteAt {
	return &pb.ApiResponsePaginationReviewDetailsDeleteAt{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesReviewDetailDeleteAt(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *reviewDetailProtoMapper) ToProtoResponsePaginationReviewDetail(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.ReviewDetailsResponse) *pb.ApiResponsePaginationReviewDetails {
	return &pb.ApiResponsePaginationReviewDetails{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesReviewDetail(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *reviewDetailProtoMapper) mapResponseReviewDetail(reviewDetail *response.ReviewDetailsResponse) *pb.ReviewDetailsResponse {
	return &pb.ReviewDetailsResponse{
		Id:        int32(reviewDetail.ID),
		ReviewId:  int32(reviewDetail.ReviewID),
		Type:      reviewDetail.Type,
		Url:       reviewDetail.Url,
		Caption:   reviewDetail.Caption,
		CreatedAt: reviewDetail.CreatedAt,
		UpdatedAt: reviewDetail.UpdatedAt,
	}
}

func (m *reviewDetailProtoMapper) mapResponsesReviewDetail(ReviewDetails []*response.ReviewDetailsResponse) []*pb.ReviewDetailsResponse {
	var mappedReviewDetails []*pb.ReviewDetailsResponse

	for _, ReviewDetail := range ReviewDetails {
		mappedReviewDetails = append(mappedReviewDetails, m.mapResponseReviewDetail(ReviewDetail))
	}

	return mappedReviewDetails
}

func (m *reviewDetailProtoMapper) mapResponseReviewDetailDeleteAt(reviewDetail *response.ReviewDetailsResponseDeleteAt) *pb.ReviewDetailsResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if reviewDetail.DeletedAt != nil {
		deletedAt = wrapperspb.String(*reviewDetail.DeletedAt)
	}

	return &pb.ReviewDetailsResponseDeleteAt{
		Id:        int32(reviewDetail.ID),
		ReviewId:  int32(reviewDetail.ReviewID),
		Type:      reviewDetail.Type,
		Url:       reviewDetail.Url,
		Caption:   reviewDetail.Caption,
		CreatedAt: reviewDetail.CreatedAt,
		UpdatedAt: reviewDetail.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (m *reviewDetailProtoMapper) mapResponsesReviewDetailDeleteAt(ReviewDetails []*response.ReviewDetailsResponseDeleteAt) []*pb.ReviewDetailsResponseDeleteAt {
	var mappedReviewDetails []*pb.ReviewDetailsResponseDeleteAt

	for _, ReviewDetail := range ReviewDetails {
		mappedReviewDetails = append(mappedReviewDetails, m.mapResponseReviewDetailDeleteAt(ReviewDetail))
	}

	return mappedReviewDetails
}
