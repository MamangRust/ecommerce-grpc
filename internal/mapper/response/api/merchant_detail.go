package response_api

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"
)

type merchantDetailResponseMapper struct {
}

func NewMerchantDetailResponseMapper() *merchantDetailResponseMapper {
	return &merchantDetailResponseMapper{}
}

func (m *merchantDetailResponseMapper) ToResponseMerchantDetail(merchant *pb.MerchantDetailResponse) *response.MerchantDetailResponse {
	return &response.MerchantDetailResponse{
		ID:               int(merchant.Id),
		MerchantID:       int(merchant.MerchantId),
		DisplayName:      merchant.DisplayName,
		CoverImageUrl:    merchant.CoverImageUrl,
		LogoUrl:          merchant.LogoUrl,
		ShortDescription: merchant.ShortDescription,
		WebsiteUrl:       merchant.WebsiteUrl,
		SocialMediaLinks: nil,
		CreatedAt:        merchant.CreatedAt,
		UpdatedAt:        merchant.UpdatedAt,
	}
}

func (m *merchantDetailResponseMapper) ToResponseMerchantDetailRelation(merchant *pb.MerchantDetailResponse) *response.MerchantDetailResponse {
	var socialMediaLinks []*response.MerchantSocialMediaLinkResponse
	for _, sm := range merchant.SocialMediaLinks {
		socialMediaLinks = append(socialMediaLinks, &response.MerchantSocialMediaLinkResponse{
			ID:       int(sm.Id),
			Platform: sm.Platform,
			Url:      sm.Url,
		})
	}

	return &response.MerchantDetailResponse{
		ID:               int(merchant.Id),
		MerchantID:       int(merchant.MerchantId),
		DisplayName:      merchant.DisplayName,
		CoverImageUrl:    merchant.CoverImageUrl,
		LogoUrl:          merchant.LogoUrl,
		ShortDescription: merchant.ShortDescription,
		WebsiteUrl:       merchant.WebsiteUrl,
		SocialMediaLinks: socialMediaLinks,
		CreatedAt:        merchant.CreatedAt,
		UpdatedAt:        merchant.UpdatedAt,
	}
}

func (m *merchantDetailResponseMapper) ToResponsesMerchant(merchants []*pb.MerchantDetailResponse) []*response.MerchantDetailResponse {
	var mappedMerchants []*response.MerchantDetailResponse

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.ToResponseMerchantDetailRelation(merchant))
	}

	return mappedMerchants
}

func (m *merchantDetailResponseMapper) ToResponseMerchantDetailDeleteAt(merchant *pb.MerchantDetailResponseDeleteAt) *response.MerchantDetailResponseDeleteAt {
	var deletedAt string
	if merchant.DeletedAt != nil {
		deletedAt = merchant.DeletedAt.Value
	}

	var socialMediaLinks []*response.MerchantSocialMediaLinkResponse
	for _, sm := range merchant.SocialMediaLinks {
		socialMediaLinks = append(socialMediaLinks, &response.MerchantSocialMediaLinkResponse{
			ID:       int(sm.Id),
			Platform: sm.Platform,
			Url:      sm.Url,
		})
	}

	return &response.MerchantDetailResponseDeleteAt{
		ID:               int(merchant.Id),
		MerchantID:       int(merchant.MerchantId),
		DisplayName:      merchant.DisplayName,
		CoverImageUrl:    merchant.CoverImageUrl,
		LogoUrl:          merchant.LogoUrl,
		ShortDescription: merchant.ShortDescription,
		WebsiteUrl:       merchant.WebsiteUrl,
		SocialMediaLinks: socialMediaLinks,
		CreatedAt:        merchant.CreatedAt,
		UpdatedAt:        merchant.UpdatedAt,
		DeletedAt:        &deletedAt,
	}
}

func (m *merchantDetailResponseMapper) ToResponsesMerchantDetailDeleteAt(merchants []*pb.MerchantDetailResponseDeleteAt) []*response.MerchantDetailResponseDeleteAt {
	var mappedMerchants []*response.MerchantDetailResponseDeleteAt

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.ToResponseMerchantDetailDeleteAt(merchant))
	}

	return mappedMerchants
}

func (m *merchantDetailResponseMapper) ToApiResponseMerchantDetail(pbResponse *pb.ApiResponseMerchantDetail) *response.ApiResponseMerchantDetail {
	return &response.ApiResponseMerchantDetail{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantDetail(pbResponse.Data),
	}
}

func (m *merchantDetailResponseMapper) ToApiResponseMerchantDetailRelation(pbResponse *pb.ApiResponseMerchantDetail) *response.ApiResponseMerchantDetailRelation {
	return &response.ApiResponseMerchantDetailRelation{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantDetailRelation(pbResponse.Data),
	}
}

func (m *merchantDetailResponseMapper) ToApiResponseMerchantDetailDeleteAt(pbResponse *pb.ApiResponseMerchantDetailDeleteAt) *response.ApiResponseMerchantDetailDeleteAt {
	return &response.ApiResponseMerchantDetailDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantDetailDeleteAt(pbResponse.Data),
	}
}

func (m *merchantDetailResponseMapper) ToApiResponsesMerchantDetail(pbResponse *pb.ApiResponsesMerchantDetail) *response.ApiResponsesMerchantDetail {
	return &response.ApiResponsesMerchantDetail{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponsesMerchant(pbResponse.Data),
	}
}

func (m *merchantDetailResponseMapper) ToApiResponsePaginationMerchantDetailDeleteAt(pbResponse *pb.ApiResponsePaginationMerchantDetailDeleteAt) *response.ApiResponsePaginationMerchantDetailDeleteAt {
	return &response.ApiResponsePaginationMerchantDetailDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesMerchantDetailDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *merchantDetailResponseMapper) ToApiResponsePaginationMerchantDetail(pbResponse *pb.ApiResponsePaginationMerchantDetail) *response.ApiResponsePaginationMerchantDetail {
	return &response.ApiResponsePaginationMerchantDetail{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesMerchant(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}
