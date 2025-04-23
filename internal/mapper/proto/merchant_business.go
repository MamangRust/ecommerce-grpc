package protomapper

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type merchantBusinessProtoMapper struct {
}

func NewMerchantBusinessProtoMapper() *merchantBusinessProtoMapper {
	return &merchantBusinessProtoMapper{}
}

func (m *merchantBusinessProtoMapper) ToProtoResponseMerchantBusiness(status string, message string, pbResponse *response.MerchantBusinessResponse) *pb.ApiResponseMerchantBusiness {
	return &pb.ApiResponseMerchantBusiness{
		Status:  status,
		Message: message,
		Data:    m.mapResponseMerchantBusiness(pbResponse),
	}
}

func (m *merchantBusinessProtoMapper) ToProtoResponseMerchantBusinessDeleteAt(status string, message string, pbResponse *response.MerchantBusinessResponseDeleteAt) *pb.ApiResponseMerchantBusinessDeleteAt {
	return &pb.ApiResponseMerchantBusinessDeleteAt{
		Status:  status,
		Message: message,
		Data:    m.mapResponseMerchantBusinessDeleteAt(pbResponse),
	}
}

func (m *merchantBusinessProtoMapper) ToProtoResponsesMerchantBusiness(status string, message string, pbResponse []*response.MerchantBusinessResponse) *pb.ApiResponsesMerchantBusiness {
	return &pb.ApiResponsesMerchantBusiness{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesMerchantBusiness(pbResponse),
	}
}

func (m *merchantBusinessProtoMapper) ToProtoResponsePaginationMerchantBusinessDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.MerchantBusinessResponseDeleteAt) *pb.ApiResponsePaginationMerchantBusinessDeleteAt {
	return &pb.ApiResponsePaginationMerchantBusinessDeleteAt{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesMerchantBusinessDeleteAt(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantBusinessProtoMapper) ToProtoResponsePaginationMerchantBusiness(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.MerchantBusinessResponse) *pb.ApiResponsePaginationMerchantBusiness {
	return &pb.ApiResponsePaginationMerchantBusiness{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesMerchantBusiness(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantBusinessProtoMapper) mapResponseMerchantBusiness(merchant *response.MerchantBusinessResponse) *pb.MerchantBusinessResponse {
	return &pb.MerchantBusinessResponse{
		Id:                int32(merchant.ID),
		MerchantId:        int32(merchant.MerchantID),
		BusinessType:      merchant.BusinessType,
		TaxId:             merchant.TaxID,
		EstablishedYear:   int32(merchant.EstablishedYear),
		NumberOfEmployees: int32(merchant.NumberOfEmployees),
		WebsiteUrl:        merchant.WebsiteUrl,
		MerchantName:      *merchant.MerchantName,
		CreatedAt:         merchant.CreatedAt,
		UpdatedAt:         merchant.UpdatedAt,
	}
}

func (m *merchantBusinessProtoMapper) mapResponsesMerchantBusiness(merchants []*response.MerchantBusinessResponse) []*pb.MerchantBusinessResponse {
	var mappedMerchants []*pb.MerchantBusinessResponse

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.mapResponseMerchantBusiness(merchant))
	}

	return mappedMerchants
}

func (m *merchantBusinessProtoMapper) mapResponseMerchantBusinessDeleteAt(merchant *response.MerchantBusinessResponseDeleteAt) *pb.MerchantBusinessResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if merchant.DeletedAt != nil {
		deletedAt = wrapperspb.String(*merchant.DeletedAt)
	}

	return &pb.MerchantBusinessResponseDeleteAt{
		Id:                int32(merchant.ID),
		MerchantId:        int32(merchant.MerchantID),
		BusinessType:      merchant.BusinessType,
		TaxId:             merchant.TaxID,
		EstablishedYear:   int32(merchant.EstablishedYear),
		NumberOfEmployees: int32(merchant.NumberOfEmployees),
		WebsiteUrl:        merchant.WebsiteUrl,
		MerchantName:      merchant.MerchantName,
		CreatedAt:         merchant.CreatedAt,
		UpdatedAt:         merchant.UpdatedAt,
		DeletedAt:         deletedAt,
	}
}

func (m *merchantBusinessProtoMapper) mapResponsesMerchantBusinessDeleteAt(merchants []*response.MerchantBusinessResponseDeleteAt) []*pb.MerchantBusinessResponseDeleteAt {
	var mappedMerchants []*pb.MerchantBusinessResponseDeleteAt

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.mapResponseMerchantBusinessDeleteAt(merchant))
	}

	return mappedMerchants
}
