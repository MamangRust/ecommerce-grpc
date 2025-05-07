package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	protomapper "ecommerce/internal/mapper/proto"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	shippingaddress_errors "ecommerce/pkg/errors/shipping_address_errors"
	"math"

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

	reqService := requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	Shipping, totalRecords, err := s.shippingService.FindAll(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationShippingAddress(paginationMeta, "success", "Successfully fetched categories", Shipping)
	return so, nil
}

func (s *shippingAddressHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShipping, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	shipping, err := s.shippingService.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
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

	reqService := requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.shippingService.FindByActive(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
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

	reqService := requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.shippingService.FindByTrashed(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

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
	id := int(request.GetId())

	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	Shipping, err := s.shippingService.TrashShippingAddress(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseShippingAddressDeleteAt("success", "Successfully trashed Shipping", Shipping)

	return so, nil
}

func (s *shippingAddressHandleGrpc) RestoreShipping(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	Shipping, err := s.shippingService.RestoreShippingAddress(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseShippingAddressDeleteAt("success", "Successfully restored Shipping", Shipping)

	return so, nil
}

func (s *shippingAddressHandleGrpc) DeleteShippingPermanent(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	_, err := s.shippingService.DeleteShippingAddressPermanently(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseShippingAddressDelete("success", "Successfully deleted Shipping permanently")

	return so, nil
}

func (s *shippingAddressHandleGrpc) RestoreAllShipping(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseShippingAll, error) {
	_, err := s.shippingService.RestoreAllShippingAddress()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseShippingAddressAll("success", "Successfully restore all Shipping")

	return so, nil
}

func (s *shippingAddressHandleGrpc) DeleteShippingAddressPermanently(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseShippingAll, error) {
	_, err := s.shippingService.DeleteAllPermanentShippingAddress()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseShippingAddressAll("success", "Successfully delete Shipping permanen")

	return so, nil
}
