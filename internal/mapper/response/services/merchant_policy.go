package response_service

import (
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/response"
)

type merchantPolicyResponseMapper struct {
}

func NewMerchantPolicyResponseMapper() *merchantPolicyResponseMapper {
	return &merchantPolicyResponseMapper{}
}

func (s *merchantPolicyResponseMapper) ToMerchantPolicyResponse(merchant *record.MerchantPoliciesRecord) *response.MerchantPoliciesResponse {
	return &response.MerchantPoliciesResponse{
		ID:          merchant.ID,
		MerchantID:  merchant.MerchantID,
		PolicyType:  merchant.PolicyType,
		Title:       merchant.Title,
		Description: merchant.Description,
		CreatedAt:   merchant.CreatedAt,
		UpdatedAt:   merchant.UpdatedAt,
	}
}

func (s *merchantPolicyResponseMapper) ToMerchantsPolicyResponse(merchants []*record.MerchantPoliciesRecord) []*response.MerchantPoliciesResponse {
	var responses []*response.MerchantPoliciesResponse

	for _, merchant := range merchants {
		responses = append(responses, s.ToMerchantPolicyResponse(merchant))
	}

	return responses
}

func (s *merchantPolicyResponseMapper) ToMerchantPolicyResponseDeleteAt(merchant *record.MerchantPoliciesRecord) *response.MerchantPoliciesResponseDeleteAt {
	return &response.MerchantPoliciesResponseDeleteAt{
		ID:          merchant.ID,
		MerchantID:  merchant.MerchantID,
		PolicyType:  merchant.PolicyType,
		Title:       merchant.Title,
		Description: merchant.Description,
		CreatedAt:   merchant.CreatedAt,
		UpdatedAt:   merchant.UpdatedAt,
		DeletedAt:   merchant.DeletedAt,
	}
}

func (s *merchantPolicyResponseMapper) ToMerchantsPolicyResponseDeleteAt(merchants []*record.MerchantPoliciesRecord) []*response.MerchantPoliciesResponseDeleteAt {
	var responses []*response.MerchantPoliciesResponseDeleteAt

	for _, merchant := range merchants {
		responses = append(responses, s.ToMerchantPolicyResponseDeleteAt(merchant))
	}

	return responses
}
