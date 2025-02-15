package gapi

import (
	"context"
	protomapper "ecommerce/internal/mapper/proto"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type shippingAddressHandleGrpc struct {
	pb.UnimplementedShippingServiceServer
	shippingService service.ShippingAddressService
	mapping         protomapper.ShippingAddresProtoMapper
}

func NewShippingAddressHandleGrpc(
	shipping service.ShippingAddressService,
	mapping protomapper.ShippingAddresProtoMapper,
) *shippingAddressHandleGrpc {
	return &shippingAddressHandleGrpc{
		shippingService: shipping,
		mapping:         mapping,
	}
}

func (s *shippingAddressHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllShippingRequest) (*pb.ApiResponsePaginationShipping, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	Shipping, totalRecords, err := s.shippingService.FindAll(page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch Shippings: ",
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationShippingAddress(paginationMeta, "success", "Successfully fetched categories", Shipping)
	return so, nil
}

func (s *shippingAddressHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShipping, error) {
	if request.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid shipping address id",
		})
	}

	shipping, err := s.shippingService.FindById(int(request.GetId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch shipping address: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseShippingAddress("success", "Successfully fetched shipping address", shipping)

	return so, nil

}

func (s *shippingAddressHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllShippingRequest) (*pb.ApiResponsePaginationShippingDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	users, totalRecords, err := s.shippingService.FindByActive(search, page, pageSize)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch active Shippings: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(0),
	}
	so := s.mapping.ToProtoResponsePaginationShippingAddressDeleteAt(paginationMeta, "success", "Successfully fetched active categories", users)

	return so, nil
}

func (s *shippingAddressHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllShippingRequest) (*pb.ApiResponsePaginationShippingDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	users, totalRecords, err := s.shippingService.FindByTrashed(search, page, pageSize)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch trashed Shippings: " + err.Message,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(0),
	}

	so := s.mapping.ToProtoResponsePaginationShippingAddressDeleteAt(paginationMeta, "success", "Successfully fetched trashed categories", users)

	return so, nil
}

func (s *shippingAddressHandleGrpc) TrashedShipping(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDeleteAt, error) {
	if request.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid Shipping id",
		})
	}

	Shipping, err := s.shippingService.TrashShippingAddress(int(request.GetId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to trashed Shipping: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseShippingAddressDeleteAt("success", "Successfully trashed Shipping", Shipping)

	return so, nil
}

func (s *shippingAddressHandleGrpc) RestoreShipping(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDeleteAt, error) {
	if request.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid Shipping id",
		})
	}

	Shipping, err := s.shippingService.RestoreShippingAddress(int(request.GetId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore Shipping: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseShippingAddressDeleteAt("success", "Successfully restored Shipping", Shipping)

	return so, nil
}

func (s *shippingAddressHandleGrpc) DeleteShippingPermanent(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDelete, error) {
	if request.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid Shipping id",
		})
	}

	_, err := s.shippingService.DeleteShippingAddressPermanently(int(request.GetId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete Shipping permanently: " + err.Message,
		})
	}

	so := s.mapping.ToProtoResponseShippingAddressDelete("success", "Successfully deleted Shipping permanently")

	return so, nil
}

func (s *shippingAddressHandleGrpc) RestoreAllShipping(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseShippingAll, error) {
	_, err := s.shippingService.RestoreAllShippingAddress()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all Shipping: ",
		})
	}

	so := s.mapping.ToProtoResponseShippingAddressAll("success", "Successfully restore all Shipping")

	return so, nil
}

func (s *shippingAddressHandleGrpc) DeleteShippingAddressPermanently(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseShippingAll, error) {
	_, err := s.shippingService.DeleteAllPermanentShippingAddress()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete Shipping permanent: ",
		})
	}

	so := s.mapping.ToProtoResponseShippingAddressAll("success", "Successfully delete Shipping permanen")

	return so, nil
}
