package protomapper

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type merchantPolicyProtoMapper struct {
}

func NewMerchantPolicyProtoMapper() *merchantPolicyProtoMapper {
	return &merchantPolicyProtoMapper{}
}

func (m *merchantPolicyProtoMapper) ToProtoResponseMerchantPolicy(status string, message string, pbResponse *response.MerchantPoliciesResponse) *pb.ApiResponseMerchantPolicies {
	return &pb.ApiResponseMerchantPolicies{
		Status:  status,
		Message: message,
		Data:    m.mapResponseMerchantPolicy(pbResponse),
	}
}

func (m *merchantPolicyProtoMapper) ToProtoResponseMerchantPolicyDeleteAt(status string, message string, pbResponse *response.MerchantPoliciesResponseDeleteAt) *pb.ApiResponseMerchantPoliciesDeleteAt {
	return &pb.ApiResponseMerchantPoliciesDeleteAt{
		Status:  status,
		Message: message,
		Data:    m.mapResponseMerchantPolicyDeleteAt(pbResponse),
	}
}

func (m *merchantPolicyProtoMapper) ToProtoResponsesMerchantPolicy(status string, message string, pbResponse []*response.MerchantPoliciesResponse) *pb.ApiResponsesMerchantPolicies {
	return &pb.ApiResponsesMerchantPolicies{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesMerchantPolicy(pbResponse),
	}
}

func (m *merchantPolicyProtoMapper) ToProtoResponsePaginationMerchantPolicyDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.MerchantPoliciesResponseDeleteAt) *pb.ApiResponsePaginationMerchantPoliciesDeleteAt {
	return &pb.ApiResponsePaginationMerchantPoliciesDeleteAt{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesMerchantPolicyDeleteAt(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantPolicyProtoMapper) ToProtoResponsePaginationMerchantPolicy(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.MerchantPoliciesResponse) *pb.ApiResponsePaginationMerchantPolicies {
	return &pb.ApiResponsePaginationMerchantPolicies{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesMerchantPolicy(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantPolicyProtoMapper) mapResponseMerchantPolicy(merchant *response.MerchantPoliciesResponse) *pb.MerchantPoliciesResponse {
	return &pb.MerchantPoliciesResponse{
		Id:           int32(merchant.ID),
		MerchantId:   int32(merchant.MerchantID),
		PolicyType:   merchant.PolicyType,
		Title:        merchant.Title,
		Description:  merchant.Description,
		CreatedAt:    merchant.CreatedAt,
		UpdatedAt:    merchant.UpdatedAt,
		MerchantName: merchant.MerchantName,
	}
}

func (m *merchantPolicyProtoMapper) mapResponsesMerchantPolicy(merchants []*response.MerchantPoliciesResponse) []*pb.MerchantPoliciesResponse {
	var mappedMerchants []*pb.MerchantPoliciesResponse

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.mapResponseMerchantPolicy(merchant))
	}

	return mappedMerchants
}

func (m *merchantPolicyProtoMapper) mapResponseMerchantPolicyDeleteAt(merchant *response.MerchantPoliciesResponseDeleteAt) *pb.MerchantPoliciesResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if merchant.DeletedAt != nil {
		deletedAt = wrapperspb.String(*merchant.DeletedAt)
	}

	return &pb.MerchantPoliciesResponseDeleteAt{
		Id:           int32(merchant.ID),
		MerchantId:   int32(merchant.MerchantID),
		PolicyType:   merchant.PolicyType,
		Title:        merchant.Title,
		Description:  merchant.Description,
		CreatedAt:    merchant.CreatedAt,
		UpdatedAt:    merchant.UpdatedAt,
		MerchantName: merchant.MerchantName,
		DeletedAt:    deletedAt,
	}
}

func (m *merchantPolicyProtoMapper) mapResponsesMerchantPolicyDeleteAt(merchants []*response.MerchantPoliciesResponseDeleteAt) []*pb.MerchantPoliciesResponseDeleteAt {
	var mappedMerchants []*pb.MerchantPoliciesResponseDeleteAt

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.mapResponseMerchantPolicyDeleteAt(merchant))
	}

	return mappedMerchants
}
