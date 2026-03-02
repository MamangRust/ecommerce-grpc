package response_api

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"
)

type merchantBusinessResponseMapper struct {
}

func NewMerchantBusinessResponseMapper() *merchantBusinessResponseMapper {
	return &merchantBusinessResponseMapper{}
}

func (m *merchantBusinessResponseMapper) ToResponseMerchantBusiness(merchant *pb.MerchantBusinessResponse) *response.MerchantBusinessResponse {
	return &response.MerchantBusinessResponse{
		ID:                int(merchant.Id),
		MerchantID:        int(merchant.MerchantId),
		BusinessType:      merchant.BusinessType,
		TaxID:             merchant.TaxId,
		EstablishedYear:   int(merchant.EstablishedYear),
		NumberOfEmployees: int(merchant.NumberOfEmployees),
		WebsiteUrl:        merchant.WebsiteUrl,
		MerchantName:      &merchant.MerchantName,
		CreatedAt:         merchant.CreatedAt,
		UpdatedAt:         merchant.UpdatedAt,
	}
}

func (m *merchantBusinessResponseMapper) ToResponsesMerchantBusiness(merchants []*pb.MerchantBusinessResponse) []*response.MerchantBusinessResponse {
	var mappedMerchants []*response.MerchantBusinessResponse

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.ToResponseMerchantBusiness(merchant))
	}

	return mappedMerchants
}

func (m *merchantBusinessResponseMapper) ToResponseMerchantBusinessDeleteAt(merchant *pb.MerchantBusinessResponseDeleteAt) *response.MerchantBusinessResponseDeleteAt {
	var deletedAt string
	if merchant.DeletedAt != nil {
		deletedAt = merchant.DeletedAt.Value
	}

	return &response.MerchantBusinessResponseDeleteAt{
		ID:                int(merchant.Id),
		MerchantID:        int(merchant.MerchantId),
		BusinessType:      merchant.BusinessType,
		TaxID:             merchant.TaxId,
		EstablishedYear:   int(merchant.EstablishedYear),
		NumberOfEmployees: int(merchant.NumberOfEmployees),
		WebsiteUrl:        merchant.WebsiteUrl,
		MerchantName:      merchant.MerchantName,
		CreatedAt:         merchant.CreatedAt,
		UpdatedAt:         merchant.UpdatedAt,
		DeletedAt:         &deletedAt,
	}
}

func (m *merchantBusinessResponseMapper) ToResponsesMerchantBusinessDeleteAt(merchants []*pb.MerchantBusinessResponseDeleteAt) []*response.MerchantBusinessResponseDeleteAt {
	var mappedMerchants []*response.MerchantBusinessResponseDeleteAt

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.ToResponseMerchantBusinessDeleteAt(merchant))
	}

	return mappedMerchants
}

func (m *merchantBusinessResponseMapper) ToApiResponseMerchantBusiness(pbResponse *pb.ApiResponseMerchantBusiness) *response.ApiResponseMerchantBusiness {
	return &response.ApiResponseMerchantBusiness{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantBusiness(pbResponse.Data),
	}
}

func (m *merchantBusinessResponseMapper) ToApiResponseMerchantBusinessDeleteAt(pbResponse *pb.ApiResponseMerchantBusinessDeleteAt) *response.ApiResponseMerchantBusinessDeleteAt {
	return &response.ApiResponseMerchantBusinessDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantBusinessDeleteAt(pbResponse.Data),
	}
}

func (m *merchantBusinessResponseMapper) ToApiResponsesMerchantBusiness(pbResponse *pb.ApiResponsesMerchantBusiness) *response.ApiResponsesMerchantBusiness {
	return &response.ApiResponsesMerchantBusiness{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponsesMerchantBusiness(pbResponse.Data),
	}
}

func (m *merchantBusinessResponseMapper) ToApiResponsePaginationMerchantBusinessDeleteAt(pbResponse *pb.ApiResponsePaginationMerchantBusinessDeleteAt) *response.ApiResponsePaginationMerchantBusinessDeleteAt {
	return &response.ApiResponsePaginationMerchantBusinessDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesMerchantBusinessDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *merchantBusinessResponseMapper) ToApiResponsePaginationMerchantBusiness(pbResponse *pb.ApiResponsePaginationMerchantBusiness) *response.ApiResponsePaginationMerchantBusiness {
	return &response.ApiResponsePaginationMerchantBusiness{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesMerchantBusiness(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}
