package response_api

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"
)

type bannerResponseMapper struct {
}

func NewBannerResponseMapper() *bannerResponseMapper {
	return &bannerResponseMapper{}
}

func (m *bannerResponseMapper) ToResponseBanner(banner *pb.BannerResponse) *response.BannerResponse {
	return &response.BannerResponse{
		ID:        banner.BannerId,
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

func (m *bannerResponseMapper) ToResponsesBanner(Banners []*pb.BannerResponse) []*response.BannerResponse {
	var mappedBanners []*response.BannerResponse

	for _, Banner := range Banners {
		mappedBanners = append(mappedBanners, m.ToResponseBanner(Banner))
	}

	return mappedBanners
}

func (m *bannerResponseMapper) ToResponseBannerDeleteAt(banner *pb.BannerResponseDeleteAt) *response.BannerResponseDeleteAt {
	var deletedAt string
	if banner.DeletedAt != nil {
		deletedAt = banner.DeletedAt.Value
	}

	return &response.BannerResponseDeleteAt{
		ID:        banner.BannerId,
		Name:      banner.Name,
		StartDate: banner.StartDate,
		EndDate:   banner.EndDate,
		StartTime: banner.StartTime,
		EndTime:   banner.EndTime,
		IsActive:  banner.IsActive,
		CreatedAt: banner.CreatedAt,
		UpdatedAt: banner.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

func (m *bannerResponseMapper) ToResponsesBannerDeleteAt(Banners []*pb.BannerResponseDeleteAt) []*response.BannerResponseDeleteAt {
	var mappedBanners []*response.BannerResponseDeleteAt

	for _, Banner := range Banners {
		mappedBanners = append(mappedBanners, m.ToResponseBannerDeleteAt(Banner))
	}

	return mappedBanners
}

func (m *bannerResponseMapper) ToApiResponseBanner(pbResponse *pb.ApiResponseBanner) *response.ApiResponseBanner {
	return &response.ApiResponseBanner{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseBanner(pbResponse.Data),
	}
}

func (m *bannerResponseMapper) ToApiResponseBannerDeleteAt(pbResponse *pb.ApiResponseBannerDeleteAt) *response.ApiResponseBannerDeleteAt {
	return &response.ApiResponseBannerDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseBannerDeleteAt(pbResponse.Data),
	}
}

func (m *bannerResponseMapper) ToApiResponsesBanner(pbResponse *pb.ApiResponsesBanner) *response.ApiResponsesBanner {
	return &response.ApiResponsesBanner{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponsesBanner(pbResponse.Data),
	}
}

func (m *bannerResponseMapper) ToApiResponseBannerDelete(pbResponse *pb.ApiResponseBannerDelete) *response.ApiResponseBannerDelete {
	return &response.ApiResponseBannerDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (m *bannerResponseMapper) ToApiResponseBannerAll(pbResponse *pb.ApiResponseBannerAll) *response.ApiResponseBannerAll {
	return &response.ApiResponseBannerAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (m *bannerResponseMapper) ToApiResponsePaginationBannerDeleteAt(pbResponse *pb.ApiResponsePaginationBannerDeleteAt) *response.ApiResponsePaginationBannerDeleteAt {
	return &response.ApiResponsePaginationBannerDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesBannerDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *bannerResponseMapper) ToApiResponsePaginationBanner(pbResponse *pb.ApiResponsePaginationBanner) *response.ApiResponsePaginationBanner {
	return &response.ApiResponsePaginationBanner{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesBanner(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}
