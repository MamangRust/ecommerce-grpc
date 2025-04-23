package response_service

import (
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/response"
)

type merchantBusinessResponseMapper struct {
}

func NewMerchantBusinessResponseMapper() *merchantBusinessResponseMapper {
	return &merchantBusinessResponseMapper{}
}

func (s *merchantBusinessResponseMapper) ToMerchantBusinessResponse(merchant *record.MerchantBusinessRecord) *response.MerchantBusinessResponse {
	return &response.MerchantBusinessResponse{
		ID:                merchant.ID,
		MerchantID:        merchant.MerchantID,
		BusinessType:      merchant.BusinessType,
		TaxID:             merchant.TaxID,
		EstablishedYear:   merchant.EstablishedYear,
		NumberOfEmployees: merchant.NumberOfEmployees,
		WebsiteUrl:        merchant.WebsiteUrl,
		CreatedAt:         merchant.CreatedAt,
		UpdatedAt:         merchant.UpdatedAt,
	}
}

func (s *merchantBusinessResponseMapper) ToMerchantBusinessResponseRelation(merchant *record.MerchantBusinessRecord) *response.MerchantBusinessResponse {
	return &response.MerchantBusinessResponse{
		ID:                merchant.ID,
		MerchantID:        merchant.MerchantID,
		BusinessType:      merchant.BusinessType,
		TaxID:             merchant.TaxID,
		EstablishedYear:   merchant.EstablishedYear,
		NumberOfEmployees: merchant.NumberOfEmployees,
		MerchantName:      merchant.MerchantName,
		WebsiteUrl:        merchant.WebsiteUrl,
		CreatedAt:         merchant.CreatedAt,
		UpdatedAt:         merchant.UpdatedAt,
	}
}

func (s *merchantBusinessResponseMapper) ToMerchantsBusinessResponse(merchants []*record.MerchantBusinessRecord) []*response.MerchantBusinessResponse {
	var responses []*response.MerchantBusinessResponse

	for _, merchant := range merchants {
		responses = append(responses, s.ToMerchantBusinessResponseRelation(merchant))
	}

	return responses
}

func (s *merchantBusinessResponseMapper) ToMerchantBusinessResponseDeleteAt(merchant *record.MerchantBusinessRecord) *response.MerchantBusinessResponseDeleteAt {
	return &response.MerchantBusinessResponseDeleteAt{
		ID:                merchant.ID,
		MerchantID:        merchant.MerchantID,
		BusinessType:      merchant.BusinessType,
		TaxID:             merchant.TaxID,
		EstablishedYear:   merchant.EstablishedYear,
		NumberOfEmployees: merchant.NumberOfEmployees,
		WebsiteUrl:        merchant.WebsiteUrl,
		CreatedAt:         merchant.CreatedAt,
		UpdatedAt:         merchant.UpdatedAt,
		DeletedAt:         merchant.DeletedAt,
	}
}

func (s *merchantBusinessResponseMapper) ToMerchantsBusinessResponseDeleteAt(merchants []*record.MerchantBusinessRecord) []*response.MerchantBusinessResponseDeleteAt {
	var responses []*response.MerchantBusinessResponseDeleteAt

	for _, merchant := range merchants {
		responses = append(responses, s.ToMerchantBusinessResponseDeleteAt(merchant))
	}

	return responses
}
