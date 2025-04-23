package protomapper

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type merchantDetailProtoMapper struct {
}

func NewMerchantDetailProtoMapper() *merchantDetailProtoMapper {
	return &merchantDetailProtoMapper{}
}

func (m *merchantDetailProtoMapper) ToProtoResponseMerchantDetail(status string, message string, pbResponse *response.MerchantDetailResponse) *pb.ApiResponseMerchantDetail {
	return &pb.ApiResponseMerchantDetail{
		Status:  status,
		Message: message,
		Data:    m.mapResponseMerchantDetail(pbResponse),
	}
}

func (m *merchantDetailProtoMapper) ToProtoResponseMerchantDetailRelation(status string, message string, pbResponse *response.MerchantDetailResponse) *pb.ApiResponseMerchantDetail {
	return &pb.ApiResponseMerchantDetail{
		Status:  status,
		Message: message,
		Data:    m.mapResponseMerchantDetailRelation(pbResponse),
	}
}

func (m *merchantDetailProtoMapper) ToProtoResponseMerchantDetailDeleteAt(status string, message string, pbResponse *response.MerchantDetailResponseDeleteAt) *pb.ApiResponseMerchantDetailDeleteAt {
	return &pb.ApiResponseMerchantDetailDeleteAt{
		Status:  status,
		Message: message,
		Data:    m.mapResponseMerchantDetailDeleteAt(pbResponse),
	}
}

func (m *merchantDetailProtoMapper) ToProtoResponsesMerchantDetail(status string, message string, pbResponse []*response.MerchantDetailResponse) *pb.ApiResponsesMerchantDetail {
	return &pb.ApiResponsesMerchantDetail{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesMerchant(pbResponse),
	}
}

func (m *merchantDetailProtoMapper) ToProtoResponsePaginationMerchantDetailDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.MerchantDetailResponseDeleteAt) *pb.ApiResponsePaginationMerchantDetailDeleteAt {
	return &pb.ApiResponsePaginationMerchantDetailDeleteAt{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesMerchantDetailDeleteAt(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantDetailProtoMapper) ToProtoResponsePaginationMerchantDetail(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.MerchantDetailResponse) *pb.ApiResponsePaginationMerchantDetail {
	return &pb.ApiResponsePaginationMerchantDetail{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesMerchant(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantDetailProtoMapper) mapResponseMerchantDetail(merchant *response.MerchantDetailResponse) *pb.MerchantDetailResponse {
	return &pb.MerchantDetailResponse{
		Id:               int32(merchant.ID),
		MerchantId:       int32(merchant.MerchantID),
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

func (m *merchantDetailProtoMapper) mapResponseMerchantDetailRelation(merchant *response.MerchantDetailResponse) *pb.MerchantDetailResponse {
	var socialMediaLinks []*pb.MerchantSocialMediaLinkResponse
	for _, sm := range merchant.SocialMediaLinks {
		socialMediaLinks = append(socialMediaLinks, &pb.MerchantSocialMediaLinkResponse{
			Id:       int32(sm.ID),
			Platform: sm.Platform,
			Url:      sm.Url,
		})
	}

	return &pb.MerchantDetailResponse{
		Id:               int32(merchant.ID),
		MerchantId:       int32(merchant.MerchantID),
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

func (m *merchantDetailProtoMapper) mapResponsesMerchant(merchants []*response.MerchantDetailResponse) []*pb.MerchantDetailResponse {
	var mappedMerchants []*pb.MerchantDetailResponse

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.mapResponseMerchantDetailRelation(merchant))
	}

	return mappedMerchants
}

func (m *merchantDetailProtoMapper) mapResponseMerchantDetailDeleteAt(merchant *response.MerchantDetailResponseDeleteAt) *pb.MerchantDetailResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if merchant.DeletedAt != nil {
		deletedAt = wrapperspb.String(*merchant.DeletedAt)
	}

	var socialMediaLinks []*pb.MerchantSocialMediaLinkResponse
	for _, sm := range merchant.SocialMediaLinks {
		socialMediaLinks = append(socialMediaLinks, &pb.MerchantSocialMediaLinkResponse{
			Id:       int32(sm.ID),
			Platform: sm.Platform,
			Url:      sm.Url,
		})
	}

	return &pb.MerchantDetailResponseDeleteAt{
		Id:               int32(merchant.ID),
		MerchantId:       int32(merchant.MerchantID),
		DisplayName:      merchant.DisplayName,
		CoverImageUrl:    merchant.CoverImageUrl,
		LogoUrl:          merchant.LogoUrl,
		ShortDescription: merchant.ShortDescription,
		WebsiteUrl:       merchant.WebsiteUrl,
		SocialMediaLinks: socialMediaLinks,
		CreatedAt:        merchant.CreatedAt,
		UpdatedAt:        merchant.UpdatedAt,
		DeletedAt:        deletedAt,
	}
}

func (m *merchantDetailProtoMapper) mapResponsesMerchantDetailDeleteAt(merchants []*response.MerchantDetailResponseDeleteAt) []*pb.MerchantDetailResponseDeleteAt {
	var mappedMerchants []*pb.MerchantDetailResponseDeleteAt

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.mapResponseMerchantDetailDeleteAt(merchant))
	}

	return mappedMerchants
}
