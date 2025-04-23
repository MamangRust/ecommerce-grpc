package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	protomapper "ecommerce/internal/mapper/proto"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors_custom"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type bannerHandleGrpc struct {
	pb.UnimplementedBannerServiceServer
	bannerService service.BannerService
	mapping       protomapper.BannerProtoMapper
}

func NewBannerHaandleGrpc(
	bannerService service.BannerService,
	mapping protomapper.BannerProtoMapper,
) *bannerHandleGrpc {
	return &bannerHandleGrpc{
		bannerService: bannerService,
		mapping:       mapping,
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

	Banner, totalRecords, err := s.bannerService.FindAll(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationBanner(paginationMeta, "success", "Successfully fetched banner", Banner)
	return so, nil
}

func (s *bannerHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdBannerRequest) (*pb.ApiResponseBanner, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_request",
				Message: "Valid banner ID is required",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	Banner, err := s.bannerService.FindById(int(request.GetId()))

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseBanner("success", "Successfully fetched banner", Banner)

	return so, nil

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

	Banner, totalRecords, err := s.bannerService.FindByActive(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationBannerDeleteAt(paginationMeta, "success", "Successfully fetched active banner", Banner)

	return so, nil
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

	users, totalRecords, err := s.bannerService.FindByTrashed(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationBannerDeleteAt(paginationMeta, "success", "Successfully fetched trashed Banner", users)

	return so, nil
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
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Unable to create new Banner. Please check your input.",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	Banner, err := s.bannerService.CreateBanner(req)
	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseBanner("success", "Successfully created banner", Banner)
	return so, nil
}

func (s *bannerHandleGrpc) Update(ctx context.Context, request *pb.UpdateBannerRequest) (*pb.ApiResponseBanner, error) {
	id := int(request.GetBannerId())

	if id == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Banner ID parameter cannot be empty and must be a positive number",
				Code:    int32(codes.InvalidArgument),
			}),
		)
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
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Unable to process Banner update. Please review your data.",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	Banner, err := s.bannerService.UpdateBanner(req)
	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseBanner("success", "Successfully updated banner", Banner)
	return so, nil
}

func (s *bannerHandleGrpc) TrashedBanner(ctx context.Context, request *pb.FindByIdBannerRequest) (*pb.ApiResponseBannerDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Banner ID parameter cannot be empty and must be a positive number",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	Banner, err := s.bannerService.TrashedBanner(id)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseBannerDeleteAt("success", "Successfully trashed Banner", Banner)

	return so, nil
}

func (s *bannerHandleGrpc) RestoreBanner(ctx context.Context, request *pb.FindByIdBannerRequest) (*pb.ApiResponseBannerDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Banner ID parameter cannot be empty and must be a positive number",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	Banner, err := s.bannerService.RestoreBanner(id)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseBannerDeleteAt("success", "Successfully restored Banner", Banner)

	return so, nil
}

func (s *bannerHandleGrpc) DeleteBannerPermanent(ctx context.Context, request *pb.FindByIdBannerRequest) (*pb.ApiResponseBannerDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Banner ID parameter cannot be empty and must be a positive number",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	_, err := s.bannerService.DeleteBannerPermanent(id)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseBannerDelete("success", "Successfully deleted Banner permanently")

	return so, nil
}

func (s *bannerHandleGrpc) RestoreAllBanner(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseBannerAll, error) {
	_, err := s.bannerService.RestoreAllBanner()

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseBannerAll("success", "Successfully restore all Banner")

	return so, nil
}

func (s *bannerHandleGrpc) DeleteAllBannerPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseBannerAll, error) {
	_, err := s.bannerService.DeleteAllBannerPermanent()

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseBannerAll("success", "Successfully delete Banner permanen")

	return so, nil
}
