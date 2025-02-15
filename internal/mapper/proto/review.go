package protomapper

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"
)

type reviewProtoMapper struct {
}

func NewReviewProtoMapper() *reviewProtoMapper {
	return &reviewProtoMapper{}
}

func (r *reviewProtoMapper) ToProtoResponseReview(status string, message string, pbResponse *response.ReviewResponse) *pb.ApiResponseReview {
	return &pb.ApiResponseReview{
		Status:  status,
		Message: message,
		Data:    r.mapResponseReview(pbResponse),
	}
}

func (r *reviewProtoMapper) ToProtoResponseReviewDeleteAt(status string, message string, pbResponse *response.ReviewResponseDeleteAt) *pb.ApiResponseReviewDeleteAt {
	return &pb.ApiResponseReviewDeleteAt{
		Status:  status,
		Message: message,
		Data:    r.mapResponseReviewDeleteAt(pbResponse),
	}
}

func (r *reviewProtoMapper) ToProtoResponsesReview(status string, message string, pbResponse []*response.ReviewResponse) *pb.ApiResponsesReview {
	return &pb.ApiResponsesReview{
		Status:  status,
		Message: message,
		Data:    r.mapResponsesReview(pbResponse),
	}
}

func (r *reviewProtoMapper) ToProtoResponseReviewDelete(status string, message string) *pb.ApiResponseReviewDelete {
	return &pb.ApiResponseReviewDelete{
		Status:  status,
		Message: message,
	}
}

func (r *reviewProtoMapper) ToProtoResponsePaginationReviewDeleteAt(pagination *pb.PaginationMeta, status string, message string, reviews []*response.ReviewResponseDeleteAt) *pb.ApiResponsePaginationReviewDeleteAt {
	return &pb.ApiResponsePaginationReviewDeleteAt{
		Status:     status,
		Message:    message,
		Data:       r.mapResponsesReviewDeleteAt(reviews),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (c *reviewProtoMapper) ToProtoResponseReviewAll(status string, message string) *pb.ApiResponseReviewAll {
	return &pb.ApiResponseReviewAll{
		Status:  status,
		Message: message,
	}
}

func (r *reviewProtoMapper) ToProtoResponsePaginationReview(pagination *pb.PaginationMeta, status string, message string, reviews []*response.ReviewResponse) *pb.ApiResponsePaginationReview {
	return &pb.ApiResponsePaginationReview{
		Status:     status,
		Message:    message,
		Data:       r.mapResponsesReview(reviews),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (r *reviewProtoMapper) mapResponseReview(review *response.ReviewResponse) *pb.ReviewResponse {
	return &pb.ReviewResponse{
		Id:        int32(review.ID),
		UserId:    int32(review.UserID),
		ProductId: int32(review.ProductID),
		Name:      review.Name,
		Comment:   review.Comment,
		Rating:    int32(review.Rating),
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}
}

func (r *reviewProtoMapper) mapResponsesReview(reviews []*response.ReviewResponse) []*pb.ReviewResponse {
	var mappedReviews []*pb.ReviewResponse

	for _, review := range reviews {
		mappedReviews = append(mappedReviews, r.mapResponseReview(review))
	}

	return mappedReviews
}

func (r *reviewProtoMapper) mapResponseReviewDeleteAt(review *response.ReviewResponseDeleteAt) *pb.ReviewResponseDeleteAt {
	return &pb.ReviewResponseDeleteAt{
		Id:        int32(review.ID),
		UserId:    int32(review.UserID),
		ProductId: int32(review.ProductID),
		Name:      review.Name,
		Comment:   review.Comment,
		Rating:    int32(review.Rating),
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
		DeletedAt: review.DeletedAt,
	}
}

func (r *reviewProtoMapper) mapResponsesReviewDeleteAt(reviews []*response.ReviewResponseDeleteAt) []*pb.ReviewResponseDeleteAt {
	var mappedReviews []*pb.ReviewResponseDeleteAt

	for _, review := range reviews {
		mappedReviews = append(mappedReviews, r.mapResponseReviewDeleteAt(review))
	}

	return mappedReviews
}
