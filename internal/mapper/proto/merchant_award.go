package protomapper

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type merchantAwardProtoMapper struct {
}

func NewMerchantAwardProtoMapper() *merchantAwardProtoMapper {
	return &merchantAwardProtoMapper{}
}

func (m *merchantAwardProtoMapper) ToProtoResponseMerchantAward(status string, message string, pbResponse *response.MerchantAwardResponse) *pb.ApiResponseMerchantAward {
	return &pb.ApiResponseMerchantAward{
		Status:  status,
		Message: message,
		Data:    m.mapResponseMerchantAward(pbResponse),
	}
}

func (m *merchantAwardProtoMapper) ToProtoResponseMerchantAwardDeleteAt(status string, message string, pbResponse *response.MerchantAwardResponseDeleteAt) *pb.ApiResponseMerchantAwardDeleteAt {
	return &pb.ApiResponseMerchantAwardDeleteAt{
		Status:  status,
		Message: message,
		Data:    m.mapResponseMerchantAwardDeleteAt(pbResponse),
	}
}

func (m *merchantAwardProtoMapper) ToProtoResponsesMerchantAward(status string, message string, pbResponse []*response.MerchantAwardResponse) *pb.ApiResponsesMerchantAward {
	return &pb.ApiResponsesMerchantAward{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesMerchantAward(pbResponse),
	}
}

func (m *merchantAwardProtoMapper) ToProtoResponsePaginationMerchantAwardDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.MerchantAwardResponseDeleteAt) *pb.ApiResponsePaginationMerchantAwardDeleteAt {
	return &pb.ApiResponsePaginationMerchantAwardDeleteAt{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesMerchantAwardDeleteAt(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantAwardProtoMapper) ToProtoResponsePaginationMerchantAward(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.MerchantAwardResponse) *pb.ApiResponsePaginationMerchantAward {
	return &pb.ApiResponsePaginationMerchantAward{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesMerchantAward(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantAwardProtoMapper) mapResponseMerchantAward(MerchantAward *response.MerchantAwardResponse) *pb.MerchantAwardResponse {
	return &pb.MerchantAwardResponse{
		Id:             int32(MerchantAward.ID),
		MerchantId:     int32(MerchantAward.MerchantID),
		Title:          MerchantAward.Title,
		Description:    MerchantAward.Description,
		IssuedBy:       MerchantAward.IssuedBy,
		IssueDate:      MerchantAward.IssueDate,
		ExpiryDate:     MerchantAward.ExpiryDate,
		CertificateUrl: MerchantAward.CertificateUrl,
		CreatedAt:      MerchantAward.CreatedAt,
		UpdatedAt:      MerchantAward.UpdatedAt,
	}
}

func (m *merchantAwardProtoMapper) mapResponsesMerchantAward(MerchantAwards []*response.MerchantAwardResponse) []*pb.MerchantAwardResponse {
	var mappedMerchantAwards []*pb.MerchantAwardResponse

	for _, MerchantAward := range MerchantAwards {
		mappedMerchantAwards = append(mappedMerchantAwards, m.mapResponseMerchantAward(MerchantAward))
	}

	return mappedMerchantAwards
}

func (m *merchantAwardProtoMapper) mapResponseMerchantAwardDeleteAt(MerchantAward *response.MerchantAwardResponseDeleteAt) *pb.MerchantAwardResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if MerchantAward.DeletedAt != nil {
		deletedAt = wrapperspb.String(*MerchantAward.DeletedAt)
	}

	return &pb.MerchantAwardResponseDeleteAt{
		Id:             int32(MerchantAward.ID),
		MerchantId:     int32(MerchantAward.MerchantID),
		Title:          MerchantAward.Title,
		Description:    MerchantAward.Description,
		IssuedBy:       MerchantAward.IssuedBy,
		IssueDate:      MerchantAward.IssueDate,
		ExpiryDate:     MerchantAward.ExpiryDate,
		CertificateUrl: MerchantAward.CertificateUrl,
		CreatedAt:      MerchantAward.CreatedAt,
		UpdatedAt:      MerchantAward.UpdatedAt,
		DeletedAt:      deletedAt,
	}
}

func (m *merchantAwardProtoMapper) mapResponsesMerchantAwardDeleteAt(MerchantAwards []*response.MerchantAwardResponseDeleteAt) []*pb.MerchantAwardResponseDeleteAt {
	var mappedMerchantAwards []*pb.MerchantAwardResponseDeleteAt

	for _, MerchantAward := range MerchantAwards {
		mappedMerchantAwards = append(mappedMerchantAwards, m.mapResponseMerchantAwardDeleteAt(MerchantAward))
	}

	return mappedMerchantAwards
}
