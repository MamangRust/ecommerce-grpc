package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/errors/banner_errors"
	"fmt"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type bannerHandleGrpc struct {
	pb.UnimplementedBannerServiceServer
	bannerService service.BannerService
}

func NewBannerHandleGrpc(
	bannerService service.BannerService,
) *bannerHandleGrpc {
	return &bannerHandleGrpc{
		bannerService: bannerService,
	}
}

func (s *bannerHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllBannerRequest) (*pb.ApiResponsePaginationBanner, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllBanner{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	banners, totalRecords, err := s.bannerService.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoBanners := make([]*pb.BannerResponse, len(banners))
	for i, banner := range banners {
		protoBanner := &pb.BannerResponse{
			BannerId:  int32(banner.BannerID),
			Name:      banner.Name,
			IsActive:  *banner.IsActive,
			CreatedAt: banner.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
			UpdatedAt: banner.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		}

		if banner.StartDate.Valid {
			protoBanner.StartDate = banner.StartDate.Time.Format("2006-01-02")
		}

		if banner.EndDate.Valid {
			protoBanner.EndDate = banner.EndDate.Time.Format("2006-01-02")
		}

		if banner.StartTime.Valid {
			hours := banner.StartTime.Microseconds / (1000000 * 60 * 60)
			minutes := (banner.StartTime.Microseconds / (1000000 * 60)) % 60
			seconds := (banner.StartTime.Microseconds / 1000000) % 60
			protoBanner.StartTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
		}

		if banner.EndTime.Valid {
			hours := banner.EndTime.Microseconds / (1000000 * 60 * 60)
			minutes := (banner.EndTime.Microseconds / (1000000 * 60)) % 60
			seconds := (banner.EndTime.Microseconds / 1000000) % 60
			protoBanner.EndTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
		}

		protoBanners[i] = protoBanner
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationBanner{
		Status:     "success",
		Message:    "Successfully fetched banners",
		Data:       protoBanners,
		Pagination: paginationMeta,
	}, nil
}

func (s *bannerHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdBannerRequest) (*pb.ApiResponseBanner, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, banner_errors.ErrGrpcBannerInvalidId
	}

	banner, err := s.bannerService.FindById(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoBanner := &pb.BannerResponse{
		BannerId:  int32(banner.BannerID),
		Name:      banner.Name,
		IsActive:  *banner.IsActive,
		CreatedAt: banner.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: banner.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
	}

	if banner.StartDate.Valid {
		protoBanner.StartDate = banner.StartDate.Time.Format("2006-01-02")
	}

	if banner.EndDate.Valid {
		protoBanner.EndDate = banner.EndDate.Time.Format("2006-01-02")
	}

	if banner.StartTime.Valid {
		hours := banner.StartTime.Microseconds / (1000000 * 60 * 60)
		minutes := (banner.StartTime.Microseconds / (1000000 * 60)) % 60
		seconds := (banner.StartTime.Microseconds / 1000000) % 60
		protoBanner.StartTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	if banner.EndTime.Valid {
		hours := banner.EndTime.Microseconds / (1000000 * 60 * 60)
		minutes := (banner.EndTime.Microseconds / (1000000 * 60)) % 60
		seconds := (banner.EndTime.Microseconds / 1000000) % 60
		protoBanner.EndTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	return &pb.ApiResponseBanner{
		Status:  "success",
		Message: "Successfully fetched banner",
		Data:    protoBanner,
	}, nil
}

func (s *bannerHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllBannerRequest) (*pb.ApiResponsePaginationBannerDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllBanner{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	banners, totalRecords, err := s.bannerService.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoBanners := make([]*pb.BannerResponseDeleteAt, len(banners))
	for i, banner := range banners {
		protoBanner := &pb.BannerResponseDeleteAt{
			BannerId:  int32(banner.BannerID),
			Name:      banner.Name,
			IsActive:  *banner.IsActive,
			CreatedAt: banner.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
			UpdatedAt: banner.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		}

		if banner.StartDate.Valid {
			protoBanner.StartDate = banner.StartDate.Time.Format("2006-01-02")
		}

		if banner.EndDate.Valid {
			protoBanner.EndDate = banner.EndDate.Time.Format("2006-01-02")
		}

		if banner.StartTime.Valid {
			hours := banner.StartTime.Microseconds / (1000000 * 60 * 60)
			minutes := (banner.StartTime.Microseconds / (1000000 * 60)) % 60
			seconds := (banner.StartTime.Microseconds / 1000000) % 60
			protoBanner.StartTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
		}

		if banner.EndTime.Valid {
			hours := banner.EndTime.Microseconds / (1000000 * 60 * 60)
			minutes := (banner.EndTime.Microseconds / (1000000 * 60)) % 60
			seconds := (banner.EndTime.Microseconds / 1000000) % 60
			protoBanner.EndTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
		}

		if banner.DeletedAt.Valid {
			deletedAt := banner.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
			protoBanner.DeletedAt = &wrapperspb.StringValue{Value: deletedAt}
		}

		protoBanners[i] = protoBanner
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationBannerDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active banners",
		Data:       protoBanners,
		Pagination: paginationMeta,
	}, nil
}

func (s *bannerHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllBannerRequest) (*pb.ApiResponsePaginationBannerDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllBanner{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	banners, totalRecords, err := s.bannerService.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoBanners := make([]*pb.BannerResponseDeleteAt, len(banners))
	for i, banner := range banners {
		protoBanner := &pb.BannerResponseDeleteAt{
			BannerId:  int32(banner.BannerID),
			Name:      banner.Name,
			IsActive:  *banner.IsActive,
			CreatedAt: banner.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
			UpdatedAt: banner.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		}

		if banner.StartDate.Valid {
			protoBanner.StartDate = banner.StartDate.Time.Format("2006-01-02")
		}

		if banner.EndDate.Valid {
			protoBanner.EndDate = banner.EndDate.Time.Format("2006-01-02")
		}

		if banner.StartTime.Valid {
			hours := banner.StartTime.Microseconds / (1000000 * 60 * 60)
			minutes := (banner.StartTime.Microseconds / (1000000 * 60)) % 60
			seconds := (banner.StartTime.Microseconds / 1000000) % 60
			protoBanner.StartTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
		}

		if banner.EndTime.Valid {
			hours := banner.EndTime.Microseconds / (1000000 * 60 * 60)
			minutes := (banner.EndTime.Microseconds / (1000000 * 60)) % 60
			seconds := (banner.EndTime.Microseconds / 1000000) % 60
			protoBanner.EndTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
		}

		if banner.DeletedAt.Valid {
			deletedAt := banner.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
			protoBanner.DeletedAt = &wrapperspb.StringValue{Value: deletedAt}
		}

		protoBanners[i] = protoBanner
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationBannerDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed banners",
		Data:       protoBanners,
		Pagination: paginationMeta,
	}, nil
}

func (s *bannerHandleGrpc) Create(ctx context.Context, request *pb.CreateBannerRequest) (*pb.ApiResponseBanner, error) {
	req := &requests.CreateBannerRequest{
		Name:      request.GetName(),
		StartDate: request.GetStartDate(),
		EndDate:   request.GetEndDate(),
		StartTime: request.GetStartTime(),
		EndTime:   request.GetEndTime(),
		IsActive:  request.GetIsActive(),
	}

	if err := req.Validate(); err != nil {
		return nil, banner_errors.ErrGrpcValidateCreateBanner
	}

	banner, err := s.bannerService.CreateBanner(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoBanner := &pb.BannerResponse{
		BannerId:  int32(banner.BannerID),
		Name:      banner.Name,
		IsActive:  *banner.IsActive,
		CreatedAt: banner.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: banner.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
	}

	if banner.StartDate.Valid {
		protoBanner.StartDate = banner.StartDate.Time.Format("2006-01-02")
	}

	if banner.EndDate.Valid {
		protoBanner.EndDate = banner.EndDate.Time.Format("2006-01-02")
	}

	if banner.StartTime.Valid {
		hours := banner.StartTime.Microseconds / (1000000 * 60 * 60)
		minutes := (banner.StartTime.Microseconds / (1000000 * 60)) % 60
		seconds := (banner.StartTime.Microseconds / 1000000) % 60
		protoBanner.StartTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	if banner.EndTime.Valid {
		hours := banner.EndTime.Microseconds / (1000000 * 60 * 60)
		minutes := (banner.EndTime.Microseconds / (1000000 * 60)) % 60
		seconds := (banner.EndTime.Microseconds / 1000000) % 60
		protoBanner.EndTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	return &pb.ApiResponseBanner{
		Status:  "success",
		Message: "Successfully created banner",
		Data:    protoBanner,
	}, nil
}

func (s *bannerHandleGrpc) Update(ctx context.Context, request *pb.UpdateBannerRequest) (*pb.ApiResponseBanner, error) {
	id := int(request.GetBannerId())

	if id == 0 {
		return nil, banner_errors.ErrGrpcBannerInvalidId
	}

	req := &requests.UpdateBannerRequest{
		BannerID:  &id,
		Name:      request.GetName(),
		StartDate: request.GetStartDate(),
		EndDate:   request.GetEndDate(),
		StartTime: request.GetStartTime(),
		EndTime:   request.GetEndTime(),
		IsActive:  request.GetIsActive(),
	}

	if err := req.Validate(); err != nil {
		return nil, banner_errors.ErrGrpcValidateUpdateBanner
	}

	banner, err := s.bannerService.UpdateBanner(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoBanner := &pb.BannerResponse{
		BannerId:  int32(banner.BannerID),
		Name:      banner.Name,
		IsActive:  *banner.IsActive,
		CreatedAt: banner.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: banner.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
	}

	if banner.StartDate.Valid {
		protoBanner.StartDate = banner.StartDate.Time.Format("2006-01-02")
	}

	if banner.EndDate.Valid {
		protoBanner.EndDate = banner.EndDate.Time.Format("2006-01-02")
	}

	if banner.StartTime.Valid {
		hours := banner.StartTime.Microseconds / (1000000 * 60 * 60)
		minutes := (banner.StartTime.Microseconds / (1000000 * 60)) % 60
		seconds := (banner.StartTime.Microseconds / 1000000) % 60
		protoBanner.StartTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	if banner.EndTime.Valid {
		hours := banner.EndTime.Microseconds / (1000000 * 60 * 60)
		minutes := (banner.EndTime.Microseconds / (1000000 * 60)) % 60
		seconds := (banner.EndTime.Microseconds / 1000000) % 60
		protoBanner.EndTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	return &pb.ApiResponseBanner{
		Status:  "success",
		Message: "Successfully updated banner",
		Data:    protoBanner,
	}, nil
}

func (s *bannerHandleGrpc) TrashedBanner(ctx context.Context, request *pb.FindByIdBannerRequest) (*pb.ApiResponseBannerDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, banner_errors.ErrGrpcBannerInvalidId
	}

	banner, err := s.bannerService.TrashedBanner(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoBanner := &pb.BannerResponseDeleteAt{
		BannerId:  int32(banner.BannerID),
		Name:      banner.Name,
		IsActive:  *banner.IsActive,
		CreatedAt: banner.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: banner.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
	}

	if banner.StartDate.Valid {
		protoBanner.StartDate = banner.StartDate.Time.Format("2006-01-02")
	}

	if banner.EndDate.Valid {
		protoBanner.EndDate = banner.EndDate.Time.Format("2006-01-02")
	}

	if banner.StartTime.Valid {
		hours := banner.StartTime.Microseconds / (1000000 * 60 * 60)
		minutes := (banner.StartTime.Microseconds / (1000000 * 60)) % 60
		seconds := (banner.StartTime.Microseconds / 1000000) % 60
		protoBanner.StartTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	if banner.EndTime.Valid {
		hours := banner.EndTime.Microseconds / (1000000 * 60 * 60)
		minutes := (banner.EndTime.Microseconds / (1000000 * 60)) % 60
		seconds := (banner.EndTime.Microseconds / 1000000) % 60
		protoBanner.EndTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	if banner.DeletedAt.Valid {
		deletedAt := banner.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		protoBanner.DeletedAt = &wrapperspb.StringValue{Value: deletedAt}
	}

	return &pb.ApiResponseBannerDeleteAt{
		Status:  "success",
		Message: "Successfully trashed banner",
		Data:    protoBanner,
	}, nil
}

func (s *bannerHandleGrpc) RestoreBanner(ctx context.Context, request *pb.FindByIdBannerRequest) (*pb.ApiResponseBannerDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, banner_errors.ErrGrpcBannerInvalidId
	}

	banner, err := s.bannerService.RestoreBanner(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoBanner := &pb.BannerResponseDeleteAt{
		BannerId:  int32(banner.BannerID),
		Name:      banner.Name,
		IsActive:  *banner.IsActive,
		CreatedAt: banner.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: banner.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
	}

	if banner.StartDate.Valid {
		protoBanner.StartDate = banner.StartDate.Time.Format("2006-01-02")
	}

	if banner.EndDate.Valid {
		protoBanner.EndDate = banner.EndDate.Time.Format("2006-01-02")
	}

	if banner.StartTime.Valid {
		hours := banner.StartTime.Microseconds / (1000000 * 60 * 60)
		minutes := (banner.StartTime.Microseconds / (1000000 * 60)) % 60
		seconds := (banner.StartTime.Microseconds / 1000000) % 60
		protoBanner.StartTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	if banner.EndTime.Valid {
		hours := banner.EndTime.Microseconds / (1000000 * 60 * 60)
		minutes := (banner.EndTime.Microseconds / (1000000 * 60)) % 60
		seconds := (banner.EndTime.Microseconds / 1000000) % 60
		protoBanner.EndTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	if banner.DeletedAt.Valid {
		deletedAt := banner.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		protoBanner.DeletedAt = &wrapperspb.StringValue{Value: deletedAt}
	}

	return &pb.ApiResponseBannerDeleteAt{
		Status:  "success",
		Message: "Successfully restored banner",
		Data:    protoBanner,
	}, nil
}

func (s *bannerHandleGrpc) DeleteBannerPermanent(ctx context.Context, request *pb.FindByIdBannerRequest) (*pb.ApiResponseBannerDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, banner_errors.ErrGrpcBannerInvalidId
	}

	_, err := s.bannerService.DeleteBannerPermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseBannerDelete{
		Status:  "success",
		Message: "Successfully deleted banner permanently",
	}, nil
}

func (s *bannerHandleGrpc) RestoreAllBanner(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseBannerAll, error) {
	_, err := s.bannerService.RestoreAllBanner(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseBannerAll{
		Status:  "success",
		Message: "Successfully restored all banners",
	}, nil
}

func (s *bannerHandleGrpc) DeleteAllBannerPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseBannerAll, error) {
	_, err := s.bannerService.DeleteAllBannerPermanent(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseBannerAll{
		Status: "Successfully deleted all banners permanently",
	}, nil
}
