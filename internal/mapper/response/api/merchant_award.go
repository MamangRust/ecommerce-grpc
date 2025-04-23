package response_api

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"
)

type merchantAwardResponseMapper struct {
}

func NewMerchantAwardResponseMapper() *merchantAwardResponseMapper {
	return &merchantAwardResponseMapper{}
}

func (m *merchantAwardResponseMapper) ToResponseMerchantAward(MerchantAward *pb.MerchantAwardResponse) *response.MerchantAwardResponse {
	return &response.MerchantAwardResponse{
		ID:             int(MerchantAward.Id),
		MerchantID:     int(MerchantAward.MerchantId),
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

func (m *merchantAwardResponseMapper) ToResponsesMerchantAward(MerchantAwards []*pb.MerchantAwardResponse) []*response.MerchantAwardResponse {
	var mappedMerchantAwards []*response.MerchantAwardResponse

	for _, MerchantAward := range MerchantAwards {
		mappedMerchantAwards = append(mappedMerchantAwards, m.ToResponseMerchantAward(MerchantAward))
	}

	return mappedMerchantAwards
}

func (m *merchantAwardResponseMapper) ToResponseMerchantAwardDeleteAt(MerchantAward *pb.MerchantAwardResponseDeleteAt) *response.MerchantAwardResponseDeleteAt {
	var deletedAt string
	if MerchantAward.DeletedAt != nil {
		deletedAt = MerchantAward.DeletedAt.Value
	}

	return &response.MerchantAwardResponseDeleteAt{
		ID:             int(MerchantAward.Id),
		MerchantID:     int(MerchantAward.MerchantId),
		Title:          MerchantAward.Title,
		Description:    MerchantAward.Description,
		IssuedBy:       MerchantAward.IssuedBy,
		IssueDate:      MerchantAward.IssueDate,
		ExpiryDate:     MerchantAward.ExpiryDate,
		CertificateUrl: MerchantAward.CertificateUrl,
		CreatedAt:      MerchantAward.CreatedAt,
		UpdatedAt:      MerchantAward.UpdatedAt,
		DeletedAt:      &deletedAt,
	}
}

func (m *merchantAwardResponseMapper) ToResponsesMerchantAwardDeleteAt(MerchantAwards []*pb.MerchantAwardResponseDeleteAt) []*response.MerchantAwardResponseDeleteAt {
	var mappedMerchantAwards []*response.MerchantAwardResponseDeleteAt

	for _, MerchantAward := range MerchantAwards {
		mappedMerchantAwards = append(mappedMerchantAwards, m.ToResponseMerchantAwardDeleteAt(MerchantAward))
	}

	return mappedMerchantAwards
}

func (m *merchantAwardResponseMapper) ToApiResponseMerchantAward(pbResponse *pb.ApiResponseMerchantAward) *response.ApiResponseMerchantAward {
	return &response.ApiResponseMerchantAward{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantAward(pbResponse.Data),
	}
}

func (m *merchantAwardResponseMapper) ToApiResponseMerchantAwardDeleteAt(pbResponse *pb.ApiResponseMerchantAwardDeleteAt) *response.ApiResponseMerchantAwardDeleteAt {
	return &response.ApiResponseMerchantAwardDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantAwardDeleteAt(pbResponse.Data),
	}
}

func (m *merchantAwardResponseMapper) ToApiResponsesMerchantAward(pbResponse *pb.ApiResponsesMerchantAward) *response.ApiResponsesMerchantAward {
	return &response.ApiResponsesMerchantAward{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponsesMerchantAward(pbResponse.Data),
	}
}

func (m *merchantAwardResponseMapper) ToApiResponsePaginationMerchantAwardDeleteAt(pbResponse *pb.ApiResponsePaginationMerchantAwardDeleteAt) *response.ApiResponsePaginationMerchantAwardDeleteAt {
	return &response.ApiResponsePaginationMerchantAwardDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesMerchantAwardDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *merchantAwardResponseMapper) ToApiResponsePaginationMerchantAward(pbResponse *pb.ApiResponsePaginationMerchantAward) *response.ApiResponsePaginationMerchantAward {
	return &response.ApiResponsePaginationMerchantAward{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesMerchantAward(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}
