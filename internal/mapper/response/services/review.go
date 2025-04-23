package response_service

import (
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/response"
)

type reviewResponseMapper struct {
}

func NewReviewResponseMapper() *reviewResponseMapper {
	return &reviewResponseMapper{}
}

func (r *reviewResponseMapper) ToReviewResponse(review *record.ReviewRecord) *response.ReviewResponse {
	return &response.ReviewResponse{
		ID:        review.ID,
		UserID:    review.UserID,
		ProductID: review.ProductID,
		Name:      review.Name,
		Comment:   review.Comment,
		Rating:    review.Rating,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}
}

func (r *reviewResponseMapper) ToReviewsResponse(reviews []*record.ReviewRecord) []*response.ReviewResponse {
	var responses []*response.ReviewResponse

	for _, review := range reviews {
		responses = append(responses, r.ToReviewResponse(review))
	}

	return responses
}

func (r *reviewResponseMapper) ToReviewResponseDeleteAt(review *record.ReviewRecord) *response.ReviewResponseDeleteAt {
	return &response.ReviewResponseDeleteAt{
		ID:        review.ID,
		UserID:    review.UserID,
		ProductID: review.ProductID,
		Name:      review.Name,
		Comment:   review.Comment,
		Rating:    review.Rating,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
		DeletedAt: review.DeletedAt,
	}
}

func (r *reviewResponseMapper) ToReviewsResponseDeleteAt(reviews []*record.ReviewRecord) []*response.ReviewResponseDeleteAt {
	var responses []*response.ReviewResponseDeleteAt

	for _, review := range reviews {
		responses = append(responses, r.ToReviewResponseDeleteAt(review))
	}

	return responses
}

func (r *reviewResponseMapper) ToReviewDetailResponse(review *record.ReviewsDetailRecord) *response.ReviewsDetailResponse {
	return &response.ReviewsDetailResponse{
		ID:        review.ID,
		UserID:    review.UserID,
		ProductID: review.ProductID,
		Name:      review.Name,
		Comment:   review.Comment,
		Rating:    review.Rating,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}
}

func (r *reviewResponseMapper) ToReviewsDetailResponse(reviews []*record.ReviewsDetailRecord) []*response.ReviewsDetailResponse {
	var responses []*response.ReviewsDetailResponse

	for _, review := range reviews {
		responses = append(responses, r.ToReviewDetailResponse(review))
	}

	return responses
}
