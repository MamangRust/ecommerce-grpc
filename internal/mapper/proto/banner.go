package protomapper

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type bannerProtoMapper struct{}

func NewBannerProtoMaper() *bannerProtoMapper {
	return &bannerProtoMapper{}
}

func (m *bannerProtoMapper) ToProtoResponseBanner(status string, message string, pbResponse *response.BannerResponse) *pb.ApiResponseBanner {
	return &pb.ApiResponseBanner{
		Status:  status,
		Message: message,
		Data:    m.mapResponseBanner(pbResponse),
	}
}

func (m *bannerProtoMapper) ToProtoResponseBannerDeleteAt(status string, message string, pbResponse *response.BannerResponseDeleteAt) *pb.ApiResponseBannerDeleteAt {
	return &pb.ApiResponseBannerDeleteAt{
		Status:  status,
		Message: message,
		Data:    m.mapResponseBannerDeleteAt(pbResponse),
	}
}

func (m *bannerProtoMapper) ToProtoResponsesBanner(status string, message string, pbResponse []*response.BannerResponse) *pb.ApiResponsesBanner {
	return &pb.ApiResponsesBanner{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesBanner(pbResponse),
	}
}

func (m *bannerProtoMapper) ToProtoResponseBannerDelete(status string, message string) *pb.ApiResponseBannerDelete {
	return &pb.ApiResponseBannerDelete{
		Status:  status,
		Message: message,
	}
}

func (m *bannerProtoMapper) ToProtoResponseBannerAll(status string, message string) *pb.ApiResponseBannerAll {
	return &pb.ApiResponseBannerAll{
		Status:  status,
		Message: message,
	}
}

func (m *bannerProtoMapper) ToProtoResponsePaginationBannerDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.BannerResponseDeleteAt) *pb.ApiResponsePaginationBannerDeleteAt {
	return &pb.ApiResponsePaginationBannerDeleteAt{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesBannerDeleteAt(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *bannerProtoMapper) ToProtoResponsePaginationBanner(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.BannerResponse) *pb.ApiResponsePaginationBanner {
	return &pb.ApiResponsePaginationBanner{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesBanner(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *bannerProtoMapper) mapResponseBanner(banner *response.BannerResponse) *pb.BannerResponse {
	return &pb.BannerResponse{
		BannerId:  banner.ID,
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

func (m *bannerProtoMapper) mapResponsesBanner(Banners []*response.BannerResponse) []*pb.BannerResponse {
	var mappedBanners []*pb.BannerResponse

	for _, Banner := range Banners {
		mappedBanners = append(mappedBanners, m.mapResponseBanner(Banner))
	}

	return mappedBanners
}

func (m *bannerProtoMapper) mapResponseBannerDeleteAt(banner *response.BannerResponseDeleteAt) *pb.BannerResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue

	if banner.DeletedAt != nil {
		deletedAt = wrapperspb.String(*banner.DeletedAt)
	}

	return &pb.BannerResponseDeleteAt{
		BannerId:  banner.ID,
		Name:      banner.Name,
		StartDate: banner.StartDate,
		EndDate:   banner.EndDate,
		StartTime: banner.StartTime,
		EndTime:   banner.EndTime,
		IsActive:  banner.IsActive,
		CreatedAt: banner.CreatedAt,
		UpdatedAt: banner.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (m *bannerProtoMapper) mapResponsesBannerDeleteAt(Banners []*response.BannerResponseDeleteAt) []*pb.BannerResponseDeleteAt {
	var mappedBanners []*pb.BannerResponseDeleteAt

	for _, Banner := range Banners {
		mappedBanners = append(mappedBanners, m.mapResponseBannerDeleteAt(Banner))
	}

	return mappedBanners
}
