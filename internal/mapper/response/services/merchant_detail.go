package response_service

import (
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/response"
)

type merchantDetailResponseMapper struct {
}

func NewMerchantDetailResponseMapper() *merchantDetailResponseMapper {
	return &merchantDetailResponseMapper{}
}

func (s *merchantDetailResponseMapper) ToMerchantDetailResponse(merchant *record.MerchantDetailRecord) *response.MerchantDetailResponse {
	return &response.MerchantDetailResponse{
		ID:               merchant.ID,
		MerchantID:       merchant.MerchantID,
		DisplayName:      merchant.DisplayName,
		CoverImageUrl:    merchant.CoverImageUrl,
		LogoUrl:          merchant.LogoUrl,
		ShortDescription: merchant.ShortDescription,
		WebsiteUrl:       merchant.WebsiteUrl,
		CreatedAt:        merchant.CreatedAt,
		UpdatedAt:        merchant.UpdatedAt,
	}
}

func (s *merchantDetailResponseMapper) ToMerchantDetailRelationResponse(merchant *record.MerchantDetailRecord) *response.MerchantDetailResponse {
	var socialMediaLinks []*response.MerchantSocialMediaLinkResponse
	for _, link := range merchant.SocialMediaLinks {
		socialMediaLinks = append(socialMediaLinks, &response.MerchantSocialMediaLinkResponse{
			ID:       link.ID,
			Platform: link.Platform,
			Url:      link.Url,
		})
	}

	return &response.MerchantDetailResponse{
		ID:               merchant.ID,
		MerchantID:       merchant.MerchantID,
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

func (s *merchantDetailResponseMapper) ToMerchantsDetailResponse(merchants []*record.MerchantDetailRecord) []*response.MerchantDetailResponse {
	var responses []*response.MerchantDetailResponse

	for _, merchant := range merchants {
		responses = append(responses, s.ToMerchantDetailRelationResponse(merchant))
	}

	return responses
}

func (s *merchantDetailResponseMapper) ToMerchantDetailResponseDeleteAt(merchant *record.MerchantDetailRecord) *response.MerchantDetailResponseDeleteAt {
	var socialMediaLinks []*response.MerchantSocialMediaLinkResponse
	for _, link := range merchant.SocialMediaLinks {
		socialMediaLinks = append(socialMediaLinks, &response.MerchantSocialMediaLinkResponse{
			ID:       link.ID,
			Platform: link.Platform,
			Url:      link.Url,
		})
	}

	return &response.MerchantDetailResponseDeleteAt{
		ID:               merchant.ID,
		MerchantID:       merchant.MerchantID,
		DisplayName:      merchant.DisplayName,
		CoverImageUrl:    merchant.CoverImageUrl,
		LogoUrl:          merchant.LogoUrl,
		ShortDescription: merchant.ShortDescription,
		WebsiteUrl:       merchant.WebsiteUrl,
		SocialMediaLinks: socialMediaLinks,
		CreatedAt:        merchant.CreatedAt,
		UpdatedAt:        merchant.UpdatedAt,
		DeletedAt:        merchant.DeletedAt,
	}
}

func (s *merchantDetailResponseMapper) ToMerchantsDetailResponseDeleteAt(merchants []*record.MerchantDetailRecord) []*response.MerchantDetailResponseDeleteAt {
	var responses []*response.MerchantDetailResponseDeleteAt

	for _, merchant := range merchants {
		responses = append(responses, s.ToMerchantDetailResponseDeleteAt(merchant))
	}

	return responses
}
