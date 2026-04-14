package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors"
	merchantdetail_errors "ecommerce/pkg/errors/merchant_detail"
	"encoding/json"
	"log"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type merchantDetailHandleGrpc struct {
	pb.UnimplementedMerchantDetailServiceServer
	merchantDetailService service.MerchantDetailService
}

func NewMerchantDetailHandleGrpc(
	merchantDetailService service.MerchantDetailService,
) *merchantDetailHandleGrpc {
	return &merchantDetailHandleGrpc{
		merchantDetailService: merchantDetailService,
	}
}

func (s *merchantDetailHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDetail, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchantDetails, totalRecords, err := s.merchantDetailService.FindAllMerchants(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchants := make([]*pb.MerchantDetailResponse, len(merchantDetails))
	for i, m := range merchantDetails {
		pbMerchant := &pb.MerchantDetailResponse{
			Id:               int32(m.MerchantDetailID),
			MerchantId:       int32(m.MerchantID),
			DisplayName:      *m.DisplayName,
			CoverImageUrl:    *m.CoverImageUrl,
			LogoUrl:          *m.LogoUrl,
			ShortDescription: *m.ShortDescription,
			WebsiteUrl:       *m.WebsiteUrl,
			CreatedAt:        m.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
			UpdatedAt:        m.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		}

		rawJSON, ok := normalizeJSON(m.SocialMediaLinks)
		if ok && len(rawJSON) > 0 {
			var socialLinks []struct {
				ID       int    `json:"id"`
				Platform string `json:"platform"`
				URL      string `json:"url"`
			}

			if err := json.Unmarshal(rawJSON, &socialLinks); err != nil {
				log.Printf("Error unmarshaling social media links: %v", err)
			} else {
				pbSocialLinks := make([]*pb.MerchantSocialMediaLinkResponse, len(socialLinks))
				for j, link := range socialLinks {
					pbSocialLinks[j] = &pb.MerchantSocialMediaLinkResponse{
						Id:       int32(link.ID),
						Platform: link.Platform,
						Url:      link.URL,
					}
				}
				pbMerchant.SocialMediaLinks = pbSocialLinks
			}
		}

		pbMerchants[i] = pbMerchant
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantDetail{
		Status:     "success",
		Message:    "Successfully fetched merchant",
		Data:       pbMerchants,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantDetailHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantDetailRequest) (*pb.ApiResponseMerchantDetail, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	m, err := s.merchantDetailService.FindById(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantDetailResponse{
		Id:               int32(m.MerchantDetailID),
		MerchantId:       int32(m.MerchantID),
		DisplayName:      *m.DisplayName,
		CoverImageUrl:    *m.CoverImageUrl,
		LogoUrl:          *m.LogoUrl,
		ShortDescription: *m.ShortDescription,
		WebsiteUrl:       *m.WebsiteUrl,
		CreatedAt:        m.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:        m.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
	}

	rawJSON, ok := normalizeJSON(m.SocialMediaLinks)
	if ok && len(rawJSON) > 0 {
		var socialLinks []struct {
			ID       int    `json:"id"`
			Platform string `json:"platform"`
			URL      string `json:"url"`
		}

		if err := json.Unmarshal(rawJSON, &socialLinks); err != nil {
			log.Printf("Error unmarshaling social media links: %v", err)
		} else {
			pbSocialLinks := make([]*pb.MerchantSocialMediaLinkResponse, len(socialLinks))
			for j, link := range socialLinks {
				pbSocialLinks[j] = &pb.MerchantSocialMediaLinkResponse{
					Id:       int32(link.ID),
					Platform: link.Platform,
					Url:      link.URL,
				}
			}
			pbMerchant.SocialMediaLinks = pbSocialLinks
		}
	}

	return &pb.ApiResponseMerchantDetail{
		Status:  "success",
		Message: "Successfully fetched merchant",
		Data:    pbMerchant,
	}, nil
}

func (s *merchantDetailHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDetailDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchantDetails, totalRecords, err := s.merchantDetailService.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchants := make([]*pb.MerchantDetailResponseDeleteAt, len(merchantDetails))
	for i, m := range merchantDetails {
		pbMerchant := &pb.MerchantDetailResponseDeleteAt{
			Id:               int32(m.MerchantDetailID),
			MerchantId:       int32(m.MerchantID),
			DisplayName:      *m.DisplayName,
			CoverImageUrl:    *m.CoverImageUrl,
			LogoUrl:          *m.LogoUrl,
			ShortDescription: *m.ShortDescription,
			WebsiteUrl:       *m.WebsiteUrl,
			CreatedAt:        m.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
			UpdatedAt:        m.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		}

		rawJSON, ok := normalizeJSON(m.SocialMediaLinks)
		if ok && len(rawJSON) > 0 {
			var socialLinks []struct {
				ID       int    `json:"id"`
				Platform string `json:"platform"`
				URL      string `json:"url"`
			}

			if err := json.Unmarshal(rawJSON, &socialLinks); err != nil {
				log.Printf("Error unmarshaling social media links: %v", err)
			} else {
				pbSocialLinks := make([]*pb.MerchantSocialMediaLinkResponse, len(socialLinks))
				for j, link := range socialLinks {
					pbSocialLinks[j] = &pb.MerchantSocialMediaLinkResponse{
						Id:       int32(link.ID),
						Platform: link.Platform,
						Url:      link.URL,
					}
				}
				pbMerchant.SocialMediaLinks = pbSocialLinks
			}
		}

		pbMerchants[i] = pbMerchant
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantDetailDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active merchant",
		Data:       pbMerchants,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantDetailHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDetailDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchantDetails, totalRecords, err := s.merchantDetailService.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchants := make([]*pb.MerchantDetailResponseDeleteAt, len(merchantDetails))
	for i, m := range merchantDetails {
		pbMerchant := &pb.MerchantDetailResponseDeleteAt{
			Id:               int32(m.MerchantDetailID),
			MerchantId:       int32(m.MerchantID),
			DisplayName:      *m.DisplayName,
			CoverImageUrl:    *m.CoverImageUrl,
			LogoUrl:          *m.LogoUrl,
			ShortDescription: *m.ShortDescription,
			WebsiteUrl:       *m.WebsiteUrl,
			CreatedAt:        m.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
			UpdatedAt:        m.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		}

		if m.DeletedAt.Valid {
			deletedAt := m.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
			pbMerchant.DeletedAt = &wrapperspb.StringValue{Value: deletedAt}
		}

		rawJSON, ok := normalizeJSON(m.SocialMediaLinks)
		if ok && len(rawJSON) > 0 {
			var socialLinks []struct {
				ID       int    `json:"id"`
				Platform string `json:"platform"`
				URL      string `json:"url"`
			}

			if err := json.Unmarshal(rawJSON, &socialLinks); err != nil {
				log.Printf("Error unmarshaling social media links: %v", err)
			} else {
				pbSocialLinks := make([]*pb.MerchantSocialMediaLinkResponse, len(socialLinks))
				for j, link := range socialLinks {
					pbSocialLinks[j] = &pb.MerchantSocialMediaLinkResponse{
						Id:       int32(link.ID),
						Platform: link.Platform,
						Url:      link.URL,
					}
				}
				pbMerchant.SocialMediaLinks = pbSocialLinks
			}
		}

		pbMerchants[i] = pbMerchant
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantDetailDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed merchant",
		Data:       pbMerchants,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantDetailHandleGrpc) Create(ctx context.Context, request *pb.CreateMerchantDetailRequest) (*pb.ApiResponseMerchantDetailCore, error) {
	socialLinks := make([]*requests.CreateMerchantSocialRequest, 0)
	for _, link := range request.GetSocialLinks() {
		socialLinks = append(socialLinks, &requests.CreateMerchantSocialRequest{
			Platform: link.GetPlatform(),
			Url:      link.GetUrl(),
		})
	}

	req := &requests.CreateMerchantDetailRequest{
		MerchantID:       int(request.GetMerchantId()),
		DisplayName:      request.GetDisplayName(),
		CoverImageUrl:    request.GetCoverImageUrl(),
		LogoUrl:          request.GetLogoUrl(),
		ShortDescription: request.GetShortDescription(),
		WebsiteUrl:       request.GetWebsiteUrl(),
		SocialLink:       socialLinks,
	}

	if err := req.Validate(); err != nil {
		return nil, merchantdetail_errors.ErrGrpcValidateCreateMerchantDetail
	}

	merchant, err := s.merchantDetailService.CreateMerchantDetail(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantDetailCoreResponse{
		Id:               int32(merchant.MerchantDetailID),
		MerchantId:       int32(merchant.MerchantID),
		DisplayName:      getStringValue(merchant.DisplayName),
		CoverImageUrl:    getStringValue(merchant.CoverImageUrl),
		LogoUrl:          getStringValue(merchant.LogoUrl),
		ShortDescription: getStringValue(merchant.ShortDescription),
		WebsiteUrl:       getStringValue(merchant.WebsiteUrl),
		CreatedAt:        merchant.CreatedAt.Time.String(),
		UpdatedAt:        merchant.UpdatedAt.Time.String(),
	}

	return &pb.ApiResponseMerchantDetailCore{
		Status:  "success",
		Message: "Successfully created merchant Detail",
		Data:    pbMerchant,
	}, nil
}

func (s *merchantDetailHandleGrpc) Update(ctx context.Context, request *pb.UpdateMerchantDetailRequest) (*pb.ApiResponseMerchantDetailCore, error) {
	id := int(request.GetMerchantDetailId())

	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	socialLinks := make([]*requests.UpdateMerchantSocialRequest, 0)
	for _, link := range request.GetSocialLinks() {
		socialLinks = append(socialLinks, &requests.UpdateMerchantSocialRequest{
			ID:               int(link.GetId()),
			Platform:         link.GetPlatform(),
			Url:              link.GetUrl(),
			MerchantDetailID: &id,
		})
	}

	req := &requests.UpdateMerchantDetailRequest{
		MerchantDetailID: &id,
		DisplayName:      request.GetDisplayName(),
		CoverImageUrl:    request.GetCoverImageUrl(),
		LogoUrl:          request.GetLogoUrl(),
		ShortDescription: request.GetShortDescription(),
		WebsiteUrl:       request.GetWebsiteUrl(),
		SocialLink:       socialLinks,
	}

	if err := req.Validate(); err != nil {
		return nil, merchantdetail_errors.ErrGrpcValidateUpdateMerchantDetail
	}

	merchant, err := s.merchantDetailService.UpdateMerchantDetail(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantDetailCoreResponse{
		Id:               int32(merchant.MerchantDetailID),
		MerchantId:       int32(merchant.MerchantID),
		DisplayName:      getStringValue(merchant.DisplayName),
		CoverImageUrl:    getStringValue(merchant.CoverImageUrl),
		LogoUrl:          getStringValue(merchant.LogoUrl),
		ShortDescription: getStringValue(merchant.ShortDescription),
		WebsiteUrl:       getStringValue(merchant.WebsiteUrl),
		CreatedAt:        merchant.CreatedAt.Time.String(),
		UpdatedAt:        merchant.UpdatedAt.Time.String(),
	}

	return &pb.ApiResponseMerchantDetailCore{
		Status:  "success",
		Message: "Successfully updated merchant Detail",
		Data:    pbMerchant,
	}, nil
}

func (s *merchantDetailHandleGrpc) TrashedMerchantDetail(ctx context.Context, request *pb.FindByIdMerchantDetailRequest) (*pb.ApiResponseMerchantDetailDeleteAtCore, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	merchant, err := s.merchantDetailService.TrashedMerchantDetail(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var deletedAt string
	if merchant.DeletedAt.Valid {
		deletedAt = merchant.DeletedAt.Time.Format("2006-01-02")
	}

	pbMerchant := &pb.MerchantDetailCoreResponseDeleteAt{
		Id:               int32(merchant.MerchantDetailID),
		MerchantId:       int32(merchant.MerchantID),
		DisplayName:      getStringValue(merchant.DisplayName),
		CoverImageUrl:    getStringValue(merchant.CoverImageUrl),
		LogoUrl:          getStringValue(merchant.LogoUrl),
		ShortDescription: getStringValue(merchant.ShortDescription),
		WebsiteUrl:       getStringValue(merchant.WebsiteUrl),
		CreatedAt:        merchant.CreatedAt.Time.String(),
		UpdatedAt:        merchant.UpdatedAt.Time.String(),
		DeletedAt:        &wrapperspb.StringValue{Value: deletedAt},
	}

	return &pb.ApiResponseMerchantDetailDeleteAtCore{
		Status:  "success",
		Message: "Successfully trashed merchant",
		Data:    pbMerchant,
	}, nil
}

func (s *merchantDetailHandleGrpc) RestoreMerchantDetail(ctx context.Context, request *pb.FindByIdMerchantDetailRequest) (*pb.ApiResponseMerchantDetailDeleteAtCore, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	merchant, err := s.merchantDetailService.RestoreMerchantDetail(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var deletedAt string
	if merchant.DeletedAt.Valid {
		deletedAt = merchant.DeletedAt.Time.Format("2006-01-02")
	}

	pbMerchant := &pb.MerchantDetailCoreResponseDeleteAt{
		Id:               int32(merchant.MerchantDetailID),
		MerchantId:       int32(merchant.MerchantID),
		DisplayName:      getStringValue(merchant.DisplayName),
		CoverImageUrl:    getStringValue(merchant.CoverImageUrl),
		LogoUrl:          getStringValue(merchant.LogoUrl),
		ShortDescription: getStringValue(merchant.ShortDescription),
		WebsiteUrl:       getStringValue(merchant.WebsiteUrl),
		CreatedAt:        merchant.CreatedAt.Time.String(),
		UpdatedAt:        merchant.UpdatedAt.Time.String(),
		DeletedAt:        &wrapperspb.StringValue{Value: deletedAt},
	}

	return &pb.ApiResponseMerchantDetailDeleteAtCore{
		Status:  "success",
		Message: "Successfully restored merchant",
		Data:    pbMerchant,
	}, nil
}

func (s *merchantDetailHandleGrpc) DeleteMerchantDetailPermanent(ctx context.Context, request *pb.FindByIdMerchantDetailRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	_, err := s.merchantDetailService.DeleteMerchantDetailPermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant permanently",
	}, nil
}

func (s *merchantDetailHandleGrpc) RestoreAllMerchantDetail(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantDetailService.RestoreAllMerchantDetail(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restore all merchant",
	}, nil
}

func (s *merchantDetailHandleGrpc) DeleteAllMerchantDetailPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantDetailService.DeleteAllMerchantDetailPermanent(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully delete all merchant permanently",
	}, nil
}

func normalizeJSON(v interface{}) ([]byte, bool) {
	if v == nil {
		return nil, false
	}

	switch t := v.(type) {
	case []byte:
		return t, true
	case string:
		return []byte(t), true
	default:
		b, err := json.Marshal(t)
		if err != nil {
			return nil, false
		}
		return b, true
	}
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
