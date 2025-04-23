package response_service

import (
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/response"
)

type merchantAwardResponseMapper struct {
}

func NewMerchantAwardResponseMapper() *merchantAwardResponseMapper {
	return &merchantAwardResponseMapper{}
}

func (s *merchantAwardResponseMapper) ToMerchantAwardResponse(merchant *record.MerchantAwardRecord) *response.MerchantAwardResponse {
	return &response.MerchantAwardResponse{
		ID:             merchant.ID,
		MerchantID:     merchant.MerchantID,
		Title:          merchant.Title,
		Description:    merchant.Description,
		IssuedBy:       merchant.IssuedBy,
		IssueDate:      merchant.IssueDate,
		ExpiryDate:     merchant.ExpiryDate,
		CertificateUrl: merchant.CertificateUrl,
		CreatedAt:      merchant.CreatedAt,
		UpdatedAt:      merchant.UpdatedAt,
	}
}

func (s *merchantAwardResponseMapper) ToMerchantsAwardResponse(merchants []*record.MerchantAwardRecord) []*response.MerchantAwardResponse {
	var responses []*response.MerchantAwardResponse

	for _, merchant := range merchants {
		responses = append(responses, s.ToMerchantAwardResponse(merchant))
	}

	return responses
}

func (s *merchantAwardResponseMapper) ToMerchantAwardResponseDeleteAt(merchant *record.MerchantAwardRecord) *response.MerchantAwardResponseDeleteAt {
	return &response.MerchantAwardResponseDeleteAt{
		ID:             merchant.ID,
		MerchantID:     merchant.MerchantID,
		Title:          merchant.Title,
		Description:    merchant.Description,
		IssuedBy:       merchant.IssuedBy,
		IssueDate:      merchant.IssueDate,
		ExpiryDate:     merchant.ExpiryDate,
		CertificateUrl: merchant.CertificateUrl,
		CreatedAt:      merchant.CreatedAt,
		UpdatedAt:      merchant.UpdatedAt,
		DeletedAt:      merchant.DeletedAt,
	}
}

func (s *merchantAwardResponseMapper) ToMerchantsAwardResponseDeleteAt(merchants []*record.MerchantAwardRecord) []*response.MerchantAwardResponseDeleteAt {
	var responses []*response.MerchantAwardResponseDeleteAt

	for _, merchant := range merchants {
		responses = append(responses, s.ToMerchantAwardResponseDeleteAt(merchant))
	}

	return responses
}
