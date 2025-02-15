package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	protomapper "ecommerce/internal/mapper/proto"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type sliderHandleGrpc struct {
	pb.UnimplementedSliderServiceServer
	sliderService service.SliderService
	mapping       protomapper.SliderProtoMapper
}

func NewSliderHandleGrpc(
	sliderService service.SliderService,
	mapping protomapper.SliderProtoMapper,
) *sliderHandleGrpc {
	return &sliderHandleGrpc{
		sliderService: sliderService,
		mapping:       mapping,
	}
}

func (s *sliderHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllSliderRequest) (*pb.ApiResponsePaginationSlider, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	category, totalRecords, err := s.sliderService.FindAll(page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch sliders: ",
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationSlider(paginationMeta, "success", "Successfully fetched slider", category)
	return so, nil
}

func (s *sliderHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllSliderRequest) (*pb.ApiResponsePaginationSliderDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	users, totalRecords, err := s.sliderService.FindByActive(search, page, pageSize)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch active categories: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(0),
	}
	so := s.mapping.ToProtoResponsePaginationSliderDeleteAt(paginationMeta, "success", "Successfully fetched active slider", users)

	return so, nil
}

func (s *sliderHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllSliderRequest) (*pb.ApiResponsePaginationSliderDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	users, totalRecords, err := s.sliderService.FindByTrashed(search, page, pageSize)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch trashed slider: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(0),
	}

	so := s.mapping.ToProtoResponsePaginationSliderDeleteAt(paginationMeta, "success", "Successfully fetched trashed slider", users)

	return so, nil
}

func (s *sliderHandleGrpc) Create(ctx context.Context, request *pb.CreateSliderRequest) (*pb.ApiResponseSlider, error) {
	req := &requests.CreateSliderRequest{
		Nama:     request.GetName(),
		FilePath: request.GetImage(),
	}

	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err.Error())
	}

	slider, err := s.sliderService.CreateSlider(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create slider:")
	}

	return s.mapping.ToProtoResponseSlider("success", "Successfully created slider", slider), nil
}

func (s *sliderHandleGrpc) Update(ctx context.Context, request *pb.UpdateSliderRequest) (*pb.ApiResponseSlider, error) {
	if request.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid slider ID")
	}

	req := &requests.UpdateSliderRequest{
		ID:       int(request.GetId()),
		Nama:     request.GetName(),
		FilePath: request.GetImage(),
	}

	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err.Error())
	}

	slider, err := s.sliderService.UpdateSlider(req)
	if err != nil {

		return nil, status.Errorf(codes.Internal, "Failed to update slider: ")
	}

	return s.mapping.ToProtoResponseSlider("success", "Successfully updated slider", slider), nil
}

func (s *sliderHandleGrpc) TrashedSlider(ctx context.Context, request *pb.FindByIdSliderRequest) (*pb.ApiResponseSliderDeleteAt, error) {
	if request.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid slider id",
		})
	}

	slider, err := s.sliderService.TrashedSlider(int(request.GetId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to trashed slider: ",
		})
	}

	so := s.mapping.ToProtoResponseSliderDeleteAt("success", "Successfully trashed slider", slider)

	return so, nil
}

func (s *sliderHandleGrpc) RestoreSlider(ctx context.Context, request *pb.FindByIdSliderRequest) (*pb.ApiResponseSliderDeleteAt, error) {
	if request.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid slider id",
		})
	}

	slider, err := s.sliderService.RestoreSlider(int(request.GetId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore slider: ",
		})
	}

	so := s.mapping.ToProtoResponseSliderDeleteAt("success", "Successfully restored slider", slider)

	return so, nil
}

func (s *sliderHandleGrpc) DeleteSliderPermanent(ctx context.Context, request *pb.FindByIdSliderRequest) (*pb.ApiResponseSliderDelete, error) {
	if request.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid slider id",
		})
	}

	_, err := s.sliderService.DeleteSliderPermanent(int(request.GetId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete slider permanently: ",
		})
	}

	so := s.mapping.ToProtoResponseSliderDelete("success", "Successfully deleted slider permanently")

	return so, nil
}

func (s *sliderHandleGrpc) RestoreAllSlider(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseSliderAll, error) {
	_, err := s.sliderService.RestoreAllSliders()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all sliders",
		})
	}

	so := s.mapping.ToProtoResponseSliderAll("success", "Successfully restored all sliders")

	return so, nil
}

func (s *sliderHandleGrpc) DeleteAllSliderPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseSliderAll, error) {
	_, err := s.sliderService.DeleteAllSlidersPermanent()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete all sliders permanently",
		})
	}

	so := s.mapping.ToProtoResponseSliderAll("success", "Successfully deleted all sliders permanently")

	return so, nil
}
