package response_api

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"
)

type reviewResponseMapper struct {
}

func NewReviewResponseMapper() *reviewResponseMapper {
	return &reviewResponseMapper{}
}

func (r *reviewResponseMapper) ToResponseReview(pbResponse *pb.ReviewResponse) *response.ReviewResponse {
	return &response.ReviewResponse{
		ID:        int(pbResponse.Id),
		UserID:    int(pbResponse.UserId),
		ProductID: int(pbResponse.ProductId),
		Rating:    int(pbResponse.Rating),
		Comment:   pbResponse.Comment,
		CreatedAt: pbResponse.CreatedAt,
		UpdatedAt: pbResponse.UpdatedAt,
	}
}

func (r *reviewResponseMapper) ToResponsesReview(pbResponses []*pb.ReviewResponse) []*response.ReviewResponse {
	var reviews []*response.ReviewResponse
	for _, review := range pbResponses {
		reviews = append(reviews, r.ToResponseReview(review))
	}
	return reviews
}

func (r *reviewResponseMapper) ToResponseReviewDeleteAt(pbResponse *pb.ReviewResponseDeleteAt) *response.ReviewResponseDeleteAt {
	var deletedAt string
	if pbResponse.DeletedAt != nil {
		deletedAt = pbResponse.DeletedAt.Value
	}

	return &response.ReviewResponseDeleteAt{
		ID:        int(pbResponse.Id),
		UserID:    int(pbResponse.UserId),
		ProductID: int(pbResponse.ProductId),
		Rating:    int(pbResponse.Rating),
		Comment:   pbResponse.Comment,
		CreatedAt: pbResponse.CreatedAt,
		UpdatedAt: pbResponse.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

func (r *reviewResponseMapper) ToResponsesReviewDeleteAt(pbResponses []*pb.ReviewResponseDeleteAt) []*response.ReviewResponseDeleteAt {
	var reviews []*response.ReviewResponseDeleteAt
	for _, review := range pbResponses {
		reviews = append(reviews, r.ToResponseReviewDeleteAt(review))
	}
	return reviews
}

func (r *reviewResponseMapper) ToResponseReviewsDetail(pbResponse *pb.ReviewsDetailResponse) *response.ReviewsDetailResponse {
	if pbResponse == nil {
		return nil
	}

	var deletedAt *string
	if pbResponse.DeletedAt != "" {
		deletedAt = &pbResponse.DeletedAt
	}

	var reviewDetail *response.ReviewDetailResponse
	if pbResponse.ReviewDetail != nil {
		reviewDetail = &response.ReviewDetailResponse{
			ID:        int(pbResponse.ReviewDetail.Id),
			Type:      pbResponse.ReviewDetail.Type,
			Url:       pbResponse.ReviewDetail.Url,
			Caption:   pbResponse.ReviewDetail.Caption,
			CreatedAt: pbResponse.ReviewDetail.CreatedAt,
		}
	}

	return &response.ReviewsDetailResponse{
		ID:           int(pbResponse.Id),
		UserID:       int(pbResponse.UserId),
		ProductID:    int(pbResponse.ProductId),
		Name:         pbResponse.Name,
		Comment:      pbResponse.Comment,
		Rating:       int(pbResponse.Rating),
		ReviewDetail: reviewDetail,
		CreatedAt:    pbResponse.CreatedAt,
		UpdatedAt:    pbResponse.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func (r *reviewResponseMapper) ToResponsesReviewsDetail(pbResponses []*pb.ReviewsDetailResponse) []*response.ReviewsDetailResponse {
	var reviews []*response.ReviewsDetailResponse
	for _, review := range pbResponses {
		reviews = append(reviews, r.ToResponseReviewsDetail(review))
	}
	return reviews
}

func (r *reviewResponseMapper) ToApiResponseReview(pbResponse *pb.ApiResponseReview) *response.ApiResponseReview {
	return &response.ApiResponseReview{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    r.ToResponseReview(pbResponse.Data),
	}
}

func (r *reviewResponseMapper) ToApiResponseReviewDeleteAt(pbResponse *pb.ApiResponseReviewDeleteAt) *response.ApiResponseReviewDeleteAt {
	return &response.ApiResponseReviewDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    r.ToResponseReviewDeleteAt(pbResponse.Data),
	}
}

func (r *reviewResponseMapper) ToApiResponsesReview(pbResponse *pb.ApiResponsesReview) *response.ApiResponsesReview {
	return &response.ApiResponsesReview{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    r.ToResponsesReview(pbResponse.Data),
	}
}

func (r *reviewResponseMapper) ToApiResponseReviewDelete(pbResponse *pb.ApiResponseReviewDelete) *response.ApiResponseReviewDelete {
	return &response.ApiResponseReviewDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (r *reviewResponseMapper) ToApiResponseReviewAll(pbResponse *pb.ApiResponseReviewAll) *response.ApiResponseReviewAll {
	return &response.ApiResponseReviewAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (r *reviewResponseMapper) ToApiResponsePaginationReviewDeleteAt(pbResponse *pb.ApiResponsePaginationReviewDeleteAt) *response.ApiResponsePaginationReviewDeleteAt {
	return &response.ApiResponsePaginationReviewDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       r.ToResponsesReviewDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (r *reviewResponseMapper) ToApiResponsePaginationReview(pbResponse *pb.ApiResponsePaginationReview) *response.ApiResponsePaginationReview {
	return &response.ApiResponsePaginationReview{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       r.ToResponsesReview(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (r *reviewResponseMapper) ToApiResponsePaginationReviewsDetail(pbResponse *pb.ApiResponsePaginationReviewDetail) *response.ApiResponsePaginationReviewsDetail {
	return &response.ApiResponsePaginationReviewsDetail{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    r.ToResponsesReviewsDetail(pbResponse.Data),
	}
}
