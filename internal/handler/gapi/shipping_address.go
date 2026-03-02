package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors"
	shippingaddress_errors "ecommerce/pkg/errors/shipping_address_errors"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type shippingAddressHandleGrpc struct {
	pb.UnimplementedShippingServiceServer
	shippingService service.ShippingAddressService
}

func NewShippingAddressHandleGrpc(
	shipping service.ShippingAddressService,
) *shippingAddressHandleGrpc {
	return &shippingAddressHandleGrpc{
		shippingService: shipping,
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

	shippingAddresses, totalRecords, err := s.shippingService.FindAllShippingAddress(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoShippingAddresses := make([]*pb.ShippingResponse, len(shippingAddresses))
	for i, shipping := range shippingAddresses {
		protoShippingAddresses[i] = &pb.ShippingResponse{
			Id:             int32(shipping.ShippingAddressID),
			OrderId:        int32(shipping.OrderID),
			Alamat:         shipping.Alamat,
			Provinsi:       shipping.Provinsi,
			Negara:         shipping.Negara,
			Kota:           shipping.Kota,
			ShippingMethod: shipping.ShippingMethod,
			ShippingCost:   int32(shipping.ShippingCost),
			CreatedAt:      shipping.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:      shipping.UpdatedAt.Time.Format("2006-01-02"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationShipping{
		Status:     "success",
		Message:    "Successfully fetched shipping addresses",
		Data:       protoShippingAddresses,
		Pagination: paginationMeta,
	}, nil
}

func (s *shippingAddressHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShipping, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	shipping, err := s.shippingService.FindById(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoShipping := &pb.ShippingResponse{
		Id:             int32(shipping.ShippingAddressID),
		OrderId:        int32(shipping.OrderID),
		Alamat:         shipping.Alamat,
		Provinsi:       shipping.Provinsi,
		Negara:         shipping.Negara,
		Kota:           shipping.Kota,
		ShippingMethod: shipping.ShippingMethod,
		ShippingCost:   int32(shipping.ShippingCost),
		CreatedAt:      shipping.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:      shipping.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseShipping{
		Status:  "success",
		Message: "Successfully fetched shipping address",
		Data:    protoShipping,
	}, nil
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

	shippingAddresses, totalRecords, err := s.shippingService.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoShippingAddresses := make([]*pb.ShippingResponseDeleteAt, len(shippingAddresses))
	for i, shipping := range shippingAddresses {
		var deletedAt string
		if shipping.DeletedAt.Valid {
			deletedAt = shipping.DeletedAt.Time.Format("2006-01-02")
		}

		protoShippingAddresses[i] = &pb.ShippingResponseDeleteAt{
			Id:             int32(shipping.ShippingAddressID),
			OrderId:        int32(shipping.OrderID),
			Alamat:         shipping.Alamat,
			Provinsi:       shipping.Provinsi,
			Negara:         shipping.Negara,
			Kota:           shipping.Kota,
			ShippingMethod: shipping.ShippingMethod,
			ShippingCost:   int32(shipping.ShippingCost),
			CreatedAt:      shipping.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:      shipping.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt:      &wrapperspb.StringValue{Value: deletedAt},
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationShippingDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active shipping addresses",
		Data:       protoShippingAddresses,
		Pagination: paginationMeta,
	}, nil
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

	shippingAddresses, totalRecords, err := s.shippingService.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoShippingAddresses := make([]*pb.ShippingResponseDeleteAt, len(shippingAddresses))
	for i, shipping := range shippingAddresses {
		var deletedAt string
		if shipping.DeletedAt.Valid {
			deletedAt = shipping.DeletedAt.Time.Format("2006-01-02")
		}

		protoShippingAddresses[i] = &pb.ShippingResponseDeleteAt{
			Id:             int32(shipping.ShippingAddressID),
			OrderId:        int32(shipping.OrderID),
			Alamat:         shipping.Alamat,
			Provinsi:       shipping.Provinsi,
			Negara:         shipping.Negara,
			Kota:           shipping.Kota,
			ShippingMethod: shipping.ShippingMethod,
			ShippingCost:   int32(shipping.ShippingCost),
			CreatedAt:      shipping.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:      shipping.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt:      &wrapperspb.StringValue{Value: deletedAt},
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationShippingDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed shipping addresses",
		Data:       protoShippingAddresses,
		Pagination: paginationMeta,
	}, nil
}

func (s *shippingAddressHandleGrpc) TrashedShipping(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	shipping, err := s.shippingService.TrashShippingAddress(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var deletedAt string
	if shipping.DeletedAt.Valid {
		deletedAt = shipping.DeletedAt.Time.Format("2006-01-02")
	}

	protoShipping := &pb.ShippingResponseDeleteAt{
		Id:             int32(shipping.ShippingAddressID),
		OrderId:        int32(shipping.OrderID),
		Alamat:         shipping.Alamat,
		Provinsi:       shipping.Provinsi,
		Negara:         shipping.Negara,
		Kota:           shipping.Kota,
		ShippingMethod: shipping.ShippingMethod,
		ShippingCost:   int32(shipping.ShippingCost),
		CreatedAt:      shipping.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:      shipping.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:      &wrapperspb.StringValue{Value: deletedAt},
	}

	return &pb.ApiResponseShippingDeleteAt{
		Status:  "success",
		Message: "Successfully trashed shipping address",
		Data:    protoShipping,
	}, nil
}

func (s *shippingAddressHandleGrpc) RestoreShipping(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	shipping, err := s.shippingService.RestoreShippingAddress(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var deletedAt string
	if shipping.DeletedAt.Valid {
		deletedAt = shipping.DeletedAt.Time.Format("2006-01-02")
	}

	protoShipping := &pb.ShippingResponseDeleteAt{
		Id:             int32(shipping.ShippingAddressID),
		OrderId:        int32(shipping.OrderID),
		Alamat:         shipping.Alamat,
		Provinsi:       shipping.Provinsi,
		Negara:         shipping.Negara,
		Kota:           shipping.Kota,
		ShippingMethod: shipping.ShippingMethod,
		ShippingCost:   int32(shipping.ShippingCost),
		CreatedAt:      shipping.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:      shipping.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:      &wrapperspb.StringValue{Value: deletedAt},
	}

	return &pb.ApiResponseShippingDeleteAt{
		Status:  "success",
		Message: "Successfully restored shipping address",
		Data:    protoShipping,
	}, nil
}

func (s *shippingAddressHandleGrpc) DeleteShippingPermanent(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	_, err := s.shippingService.DeleteShippingAddressPermanently(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseShippingDelete{
		Status:  "success",
		Message: "Successfully deleted shipping address permanently",
	}, nil
}

func (s *shippingAddressHandleGrpc) RestoreAllShipping(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseShippingAll, error) {
	_, err := s.shippingService.RestoreAllShippingAddress(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseShippingAll{
		Status:  "success",
		Message: "Successfully restored all shipping addresses",
	}, nil
}

func (s *shippingAddressHandleGrpc) DeleteShippingAddressPermanently(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseShippingAll, error) {
	_, err := s.shippingService.DeleteAllPermanentShippingAddress(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseShippingAll{
		Status:  "success",
		Message: "Successfully deleted all shipping addresses permanently",
	}, nil
}
