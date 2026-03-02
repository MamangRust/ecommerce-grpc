package response_api

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"
)

type merchantPolicyResponseMapper struct {
}

func NewMerchantPolicyResponseMapper() *merchantPolicyResponseMapper {
	return &merchantPolicyResponseMapper{}
}

func (m *merchantPolicyResponseMapper) ToResponseMerchantPolicy(merchant *pb.MerchantPoliciesResponse) *response.MerchantPoliciesResponse {
	return &response.MerchantPoliciesResponse{
		ID:           int(merchant.Id),
		MerchantID:   int(merchant.MerchantId),
		PolicyType:   merchant.PolicyType,
		Title:        merchant.Title,
		Description:  merchant.Description,
		CreatedAt:    merchant.CreatedAt,
		UpdatedAt:    merchant.UpdatedAt,
		MerchantName: merchant.MerchantName,
	}
}

func (m *merchantPolicyResponseMapper) ToResponsesMerchantPolicy(merchants []*pb.MerchantPoliciesResponse) []*response.MerchantPoliciesResponse {
	var mappedMerchants []*response.MerchantPoliciesResponse

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.ToResponseMerchantPolicy(merchant))
	}

	return mappedMerchants
}

func (m *merchantPolicyResponseMapper) ToResponseMerchantPolicyDeleteAt(merchant *pb.MerchantPoliciesResponseDeleteAt) *response.MerchantPoliciesResponseDeleteAt {
	var deletedAt string
	if merchant.DeletedAt != nil {
		deletedAt = merchant.DeletedAt.Value
	}

	return &response.MerchantPoliciesResponseDeleteAt{
		ID:           int(merchant.Id),
		MerchantID:   int(merchant.MerchantId),
		PolicyType:   merchant.PolicyType,
		Title:        merchant.Title,
		Description:  merchant.Description,
		CreatedAt:    merchant.CreatedAt,
		UpdatedAt:    merchant.UpdatedAt,
		MerchantName: merchant.MerchantName,
		DeletedAt:    &deletedAt,
	}
}

func (m *merchantPolicyResponseMapper) ToResponsesMerchantPolicyDeleteAt(merchants []*pb.MerchantPoliciesResponseDeleteAt) []*response.MerchantPoliciesResponseDeleteAt {
	var mappedMerchants []*response.MerchantPoliciesResponseDeleteAt

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.ToResponseMerchantPolicyDeleteAt(merchant))
	}

	return mappedMerchants
}

func (m *merchantPolicyResponseMapper) ToApiResponseMerchantPolicies(pbResponse *pb.ApiResponseMerchantPolicies) *response.ApiResponseMerchantPolicies {
	return &response.ApiResponseMerchantPolicies{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantPolicy(pbResponse.Data),
	}
}

func (m *merchantPolicyResponseMapper) ToApiResponseMerchantPoliciesDeleteAt(pbResponse *pb.ApiResponseMerchantPoliciesDeleteAt) *response.ApiResponseMerchantPoliciesDeleteAt {
	return &response.ApiResponseMerchantPoliciesDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantPolicyDeleteAt(pbResponse.Data),
	}
}

func (m *merchantPolicyResponseMapper) ToApiResponsesMerchantPolicies(pbResponse *pb.ApiResponsesMerchantPolicies) *response.ApiResponsesMerchantPolicies {
	return &response.ApiResponsesMerchantPolicies{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponsesMerchantPolicy(pbResponse.Data),
	}
}

func (m *merchantPolicyResponseMapper) ToApiResponsePaginationMerchantPoliciesDeleteAt(pbResponse *pb.ApiResponsePaginationMerchantPoliciesDeleteAt) *response.ApiResponsePaginationMerchantPoliciesDeleteAt {
	return &response.ApiResponsePaginationMerchantPoliciesDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesMerchantPolicyDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *merchantPolicyResponseMapper) ToApiResponsePaginationMerchantPolicies(pbResponse *pb.ApiResponsePaginationMerchantPolicies) *response.ApiResponsePaginationMerchantPolicies {
	return &response.ApiResponsePaginationMerchantPolicies{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesMerchantPolicy(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}
