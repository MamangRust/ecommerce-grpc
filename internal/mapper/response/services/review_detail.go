package response_service

import (
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/response"
)

type reviewDetailResponseMapper struct {
}

func NewReviewDetailResponseMapper() *reviewDetailResponseMapper {
	return &reviewDetailResponseMapper{}
}

func (s *reviewDetailResponseMapper) ToReviewDetailsResponse(record *record.ReviewDetailRecord) *response.ReviewDetailsResponse {
	return &response.ReviewDetailsResponse{
		ID:        record.ID,
		ReviewID:  *record.ReviewID,
		Type:      record.Type,
		Url:       record.Url,
		Caption:   record.Caption,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
	}
}

func (s *reviewDetailResponseMapper) ToReviewsDetailsResponse(records []*record.ReviewDetailRecord) []*response.ReviewDetailsResponse {
	var responses []*response.ReviewDetailsResponse
	for _, record := range records {
		responses = append(responses, s.ToReviewDetailsResponse(record))
	}
	return responses
}

func (s *reviewDetailResponseMapper) ToReviewDetailResponseDeleteAt(record *record.ReviewDetailRecord) *response.ReviewDetailsResponseDeleteAt {
	return &response.ReviewDetailsResponseDeleteAt{
		ID:        record.ID,
		ReviewID:  *record.ReviewID,
		Type:      record.Type,
		Url:       record.Url,
		Caption:   record.Caption,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
		DeletedAt: record.DeletedAt,
	}
}

func (s *reviewDetailResponseMapper) ToReviewDetailsResponseDeleteAt(merchants []*record.ReviewDetailRecord) []*response.ReviewDetailsResponseDeleteAt {
	var responses []*response.ReviewDetailsResponseDeleteAt

	for _, merchant := range merchants {
		responses = append(responses, s.ToReviewDetailResponseDeleteAt(merchant))
	}

	return responses
}
