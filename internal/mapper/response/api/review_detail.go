package response_api

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"
)

type reviewDetailResponseMapper struct {
}

func NewReviewDetailResponseMapper() *reviewDetailResponseMapper {
	return &reviewDetailResponseMapper{}
}

func (m *reviewDetailResponseMapper) ToResponseReviewDetail(reviewDetail *pb.ReviewDetailsResponse) *response.ReviewDetailsResponse {
	return &response.ReviewDetailsResponse{
		ID:        int(reviewDetail.Id),
		ReviewID:  int(reviewDetail.ReviewId),
		Type:      reviewDetail.Type,
		Url:       reviewDetail.Url,
		Caption:   reviewDetail.Caption,
		CreatedAt: reviewDetail.CreatedAt,
		UpdatedAt: reviewDetail.UpdatedAt,
	}
}

func (m *reviewDetailResponseMapper) ToResponsesReviewDetail(ReviewDetails []*pb.ReviewDetailsResponse) []*response.ReviewDetailsResponse {
	var mappedReviewDetails []*response.ReviewDetailsResponse

	for _, ReviewDetail := range ReviewDetails {
		mappedReviewDetails = append(mappedReviewDetails, m.ToResponseReviewDetail(ReviewDetail))
	}

	return mappedReviewDetails
}

func (m *reviewDetailResponseMapper) ToResponseReviewDetailDeleteAt(reviewDetail *pb.ReviewDetailsResponseDeleteAt) *response.ReviewDetailsResponseDeleteAt {
	var deletedAt string
	if reviewDetail.DeletedAt != nil {
		deletedAt = reviewDetail.DeletedAt.Value
	}

	return &response.ReviewDetailsResponseDeleteAt{
		ID:        int(reviewDetail.Id),
		ReviewID:  int(reviewDetail.ReviewId),
		Type:      reviewDetail.Type,
		Url:       reviewDetail.Url,
		Caption:   reviewDetail.Caption,
		CreatedAt: reviewDetail.CreatedAt,
		UpdatedAt: reviewDetail.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

func (m *reviewDetailResponseMapper) ToResponsesReviewDetailDeleteAt(ReviewDetails []*pb.ReviewDetailsResponseDeleteAt) []*response.ReviewDetailsResponseDeleteAt {
	var mappedReviewDetails []*response.ReviewDetailsResponseDeleteAt

	for _, ReviewDetail := range ReviewDetails {
		mappedReviewDetails = append(mappedReviewDetails, m.ToResponseReviewDetailDeleteAt(ReviewDetail))
	}

	return mappedReviewDetails
}

func (m *reviewDetailResponseMapper) ToApiResponseReviewDetail(pbResponse *pb.ApiResponseReviewDetail) *response.ApiResponseReviewDetail {
	return &response.ApiResponseReviewDetail{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseReviewDetail(pbResponse.Data),
	}
}

func (m *reviewDetailResponseMapper) ToApiResponseReviewDetailDeleteAt(pbResponse *pb.ApiResponseReviewDetailDeleteAt) *response.ApiResponseReviewDetailDeleteAt {
	return &response.ApiResponseReviewDetailDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseReviewDetailDeleteAt(pbResponse.Data),
	}
}

func (m *reviewDetailResponseMapper) ToApiResponsesReviewDetail(pbResponse *pb.ApiResponsesReviewDetails) *response.ApiResponsesReviewDetails {
	return &response.ApiResponsesReviewDetails{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponsesReviewDetail(pbResponse.Data),
	}
}

func (m *reviewDetailResponseMapper) ToApiResponsePaginationReviewDetailDeleteAt(pbResponse *pb.ApiResponsePaginationReviewDetailsDeleteAt) *response.ApiResponsePaginationReviewDetailsDeleteAt {
	return &response.ApiResponsePaginationReviewDetailsDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesReviewDetailDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *reviewDetailResponseMapper) ToApiResponsePaginationReviewDetail(pbResponse *pb.ApiResponsePaginationReviewDetails) *response.ApiResponsePaginationReviewDetails {
	return &response.ApiResponsePaginationReviewDetails{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesReviewDetail(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}
