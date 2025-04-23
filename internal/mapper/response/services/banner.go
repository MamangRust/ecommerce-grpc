package response_service

import (
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/response"
)

type bannerResponseMapper struct {
}

func NewBannerResponseMapper() *bannerResponseMapper {
	return &bannerResponseMapper{}
}

func (s *bannerResponseMapper) ToBannerResponse(banner *record.BannerRecord) *response.BannerResponse {
	return &response.BannerResponse{
		ID:        int32(banner.BannerID),
		Name:      banner.Name,
		StartDate: banner.StartDate,
		EndDate:   banner.EndDate,
		StartTime: banner.StartTime,
		EndTime:   banner.EndTime,
		IsActive:  banner.IsActive,
		CreatedAt: banner.CreatedAt,
		UpdatedAt: banner.UpdatedAt,
	}
}

func (s *bannerResponseMapper) ToBannersResponse(banners []*record.BannerRecord) []*response.BannerResponse {
	var responses []*response.BannerResponse

	for _, merchant := range banners {
		responses = append(responses, s.ToBannerResponse(merchant))
	}

	return responses
}

func (s *bannerResponseMapper) ToBannerResponseDeleteAt(banner *record.BannerRecord) *response.BannerResponseDeleteAt {
	return &response.BannerResponseDeleteAt{
		ID:        int32(banner.BannerID),
		Name:      banner.Name,
		StartDate: banner.StartDate,
		EndDate:   banner.EndDate,
		StartTime: banner.StartTime,
		EndTime:   banner.EndTime,
		IsActive:  banner.IsActive,
		CreatedAt: banner.CreatedAt,
		UpdatedAt: banner.UpdatedAt,
		DeletedAt: banner.DeletedAt,
	}
}

func (s *bannerResponseMapper) ToBannersResponseDeleteAt(banners []*record.BannerRecord) []*response.BannerResponseDeleteAt {
	var responses []*response.BannerResponseDeleteAt

	for _, merchant := range banners {
		responses = append(responses, s.ToBannerResponseDeleteAt(merchant))
	}

	return responses
}
