package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors"
	orderitem_errors "ecommerce/pkg/errors/order_item_errors"
	"math"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type orderItemHandleGrpc struct {
	pb.UnimplementedOrderItemServiceServer
	orderItemService service.OrderItemService
}

func NewOrderItemHandleGrpc(
	orderItemService service.OrderItemService,
) *orderItemHandleGrpc {
	return &orderItemHandleGrpc{
		orderItemService: orderItemService,
	}
}

func (s *orderItemHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllOrderItemRequest) (*pb.ApiResponsePaginationOrderItem, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindAllOrderItems(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var pbOrderItems []*pb.OrderItemResponse
	for _, item := range orderItems {
		pbOrderItems = append(pbOrderItems, &pb.OrderItemResponse{
			Id:        int32(item.OrderItemID),
			OrderId:   int32(item.OrderID),
			ProductId: int32(item.ProductID),
			Quantity:  int32(item.Quantity),
			Price:     int32(item.Price),
			CreatedAt: item.CreatedAt.Time.String(),
			UpdatedAt: item.UpdatedAt.Time.String(),
		})
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationOrderItem{
		Status:     "success",
		Message:    "Successfully fetched order items",
		Data:       pbOrderItems,
		Pagination: paginationMeta,
	}, nil
}

func (s *orderItemHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllOrderItemRequest) (*pb.ApiResponsePaginationOrderItemDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var pbOrderItems []*pb.OrderItemResponseDeleteAt
	for _, item := range orderItems {
		var deletedAt string

		if item.DeletedAt.Valid {
			deletedAt = item.DeletedAt.Time.Format("2006-01-02")
		}

		pbOrderItems = append(pbOrderItems, &pb.OrderItemResponseDeleteAt{
			Id:        int32(item.OrderItemID),
			OrderId:   int32(item.OrderID),
			ProductId: int32(item.ProductID),
			Quantity:  int32(item.Quantity),
			Price:     int32(item.Price),
			CreatedAt: item.CreatedAt.Time.String(),
			UpdatedAt: item.UpdatedAt.Time.String(),
			DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
		})
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationOrderItemDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active order items",
		Data:       pbOrderItems,
		Pagination: paginationMeta,
	}, nil
}

func (s *orderItemHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllOrderItemRequest) (*pb.ApiResponsePaginationOrderItemDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var pbOrderItems []*pb.OrderItemResponseDeleteAt
	for _, item := range orderItems {
		var deletedAt string

		if item.DeletedAt.Valid {
			deletedAt = item.DeletedAt.Time.Format("2006-01-02")
		}

		pbOrderItems = append(pbOrderItems, &pb.OrderItemResponseDeleteAt{
			Id:        int32(item.OrderItemID),
			OrderId:   int32(item.OrderID),
			ProductId: int32(item.ProductID),
			Quantity:  int32(item.Quantity),
			Price:     int32(item.Price),
			CreatedAt: item.CreatedAt.Time.String(),
			UpdatedAt: item.UpdatedAt.Time.String(),
			DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
		})
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationOrderItemDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed order items",
		Data:       pbOrderItems,
		Pagination: paginationMeta,
	}, nil
}

func (s *orderItemHandleGrpc) FindOrderItemByOrder(ctx context.Context, request *pb.FindByIdOrderItemRequest) (*pb.ApiResponsesOrderItem, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	orderItems, err := s.orderItemService.FindOrderItemByOrder(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var pbOrderItems []*pb.OrderItemResponse
	for _, item := range orderItems {
		pbOrderItems = append(pbOrderItems, &pb.OrderItemResponse{
			Id:        int32(item.OrderItemID),
			OrderId:   int32(item.OrderID),
			ProductId: int32(item.ProductID),
			Quantity:  int32(item.Quantity),
			Price:     int32(item.Price),
			CreatedAt: item.CreatedAt.Time.String(),
			UpdatedAt: item.UpdatedAt.Time.String(),
		})
	}

	return &pb.ApiResponsesOrderItem{
		Status:  "success",
		Message: "Successfully fetched order items by order",
		Data:    pbOrderItems,
	}, nil
}
