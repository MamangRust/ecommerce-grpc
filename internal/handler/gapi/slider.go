package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/errors/slider_errors"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type sliderHandleGrpc struct {
	pb.UnimplementedSliderServiceServer
	sliderService service.SliderService
}

func NewSliderHandleGrpc(
	sliderService service.SliderService,
) *sliderHandleGrpc {
	return &sliderHandleGrpc{
		sliderService: sliderService,
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

	reqService := requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	sliders, totalRecords, err := s.sliderService.FindAllSlider(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSliders := make([]*pb.SliderResponse, len(sliders))
	for i, slider := range sliders {
		protoSliders[i] = &pb.SliderResponse{
			Id:        int32(slider.SliderID),
			Name:      slider.Name,
			Image:     slider.Image,
			CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationSlider{
		Status:     "success",
		Message:    "Successfully fetched slider records",
		Data:       protoSliders,
		Pagination: paginationMeta,
	}, nil
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

	reqService := requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	sliders, totalRecords, err := s.sliderService.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSliders := make([]*pb.SliderResponseDeleteAt, len(sliders))
	for i, slider := range sliders {
		var deletedAt *wrapperspb.StringValue
		if slider.DeletedAt.Valid {
			deletedAt = &wrapperspb.StringValue{Value: slider.DeletedAt.Time.Format("2006-01-02")}
		}

		protoSliders[i] = &pb.SliderResponseDeleteAt{
			Id:        int32(slider.SliderID),
			Name:      slider.Name,
			Image:     slider.Image,
			CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt: deletedAt,
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationSliderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active slider records",
		Data:       protoSliders,
		Pagination: paginationMeta,
	}, nil
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

	reqService := requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	sliders, totalRecords, err := s.sliderService.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSliders := make([]*pb.SliderResponseDeleteAt, len(sliders))
	for i, slider := range sliders {
		var deletedAt *wrapperspb.StringValue
		if slider.DeletedAt.Valid {
			deletedAt = &wrapperspb.StringValue{Value: slider.DeletedAt.Time.Format("2006-01-02")}
		}

		protoSliders[i] = &pb.SliderResponseDeleteAt{
			Id:        int32(slider.SliderID),
			Name:      slider.Name,
			Image:     slider.Image,
			CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt: deletedAt,
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationSliderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed slider records",
		Data:       protoSliders,
		Pagination: paginationMeta,
	}, nil
}

func (s *sliderHandleGrpc) Create(ctx context.Context, request *pb.CreateSliderRequest) (*pb.ApiResponseSlider, error) {
	req := &requests.CreateSliderRequest{
		Nama:     request.GetName(),
		FilePath: request.GetImage(),
	}

	if err := req.Validate(); err != nil {
		return nil, slider_errors.ErrGrpcValidateCreateSlider
	}

	slider, err := s.sliderService.CreateSlider(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSlider := &pb.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseSlider{
		Status:  "success",
		Message: "Successfully created slider",
		Data:    protoSlider,
	}, nil
}

func (s *sliderHandleGrpc) Update(ctx context.Context, request *pb.UpdateSliderRequest) (*pb.ApiResponseSlider, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, slider_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateSliderRequest{
		ID:       &id,
		Nama:     request.GetName(),
		FilePath: request.GetImage(),
	}

	if err := req.Validate(); err != nil {
		return nil, slider_errors.ErrGrpcValidateUpdateSlider
	}

	slider, err := s.sliderService.UpdateSlider(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSlider := &pb.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseSlider{
		Status:  "success",
		Message: "Successfully updated slider",
		Data:    protoSlider,
	}, nil
}

func (s *sliderHandleGrpc) TrashedSlider(ctx context.Context, request *pb.FindByIdSliderRequest) (*pb.ApiResponseSliderDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, slider_errors.ErrGrpcInvalidID
	}

	slider, err := s.sliderService.TrashSlider(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var deletedAt *wrapperspb.StringValue
	if slider.DeletedAt.Valid {
		deletedAt = &wrapperspb.StringValue{Value: slider.DeletedAt.Time.Format("2006-01-02")}
	}

	protoSlider := &pb.SliderResponseDeleteAt{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: deletedAt,
	}

	return &pb.ApiResponseSliderDeleteAt{
		Status:  "success",
		Message: "Successfully trashed slider",
		Data:    protoSlider,
	}, nil
}

func (s *sliderHandleGrpc) RestoreSlider(ctx context.Context, request *pb.FindByIdSliderRequest) (*pb.ApiResponseSliderDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, slider_errors.ErrGrpcInvalidID
	}

	slider, err := s.sliderService.RestoreSlider(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var deletedAtRestore *wrapperspb.StringValue
	if slider.DeletedAt.Valid {
		deletedAtRestore = &wrapperspb.StringValue{Value: slider.DeletedAt.Time.Format("2006-01-02")}
	}

	protoSliderRestore := &pb.SliderResponseDeleteAt{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: deletedAtRestore,
	}

	return &pb.ApiResponseSliderDeleteAt{
		Status:  "success",
		Message: "Successfully restored slider",
		Data:    protoSliderRestore,
	}, nil
}

func (s *sliderHandleGrpc) DeleteSliderPermanent(ctx context.Context, request *pb.FindByIdSliderRequest) (*pb.ApiResponseSliderDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, slider_errors.ErrGrpcInvalidID
	}

	_, err := s.sliderService.DeleteSliderPermanently(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseSliderDelete{
		Status:  "success",
		Message: "Successfully deleted slider permanently",
	}, nil
}

func (s *sliderHandleGrpc) RestoreAllSlider(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseSliderAll, error) {
	_, err := s.sliderService.RestoreAllSliders(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseSliderAll{
		Status:  "success",
		Message: "Successfully restored all sliders",
	}, nil
}

func (s *sliderHandleGrpc) DeleteAllSliderPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseSliderAll, error) {
	_, err := s.sliderService.DeleteAllPermanentSlider(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseSliderAll{
		Status:  "success",
		Message: "Successfully deleted all sliders permanently",
	}, nil
}
