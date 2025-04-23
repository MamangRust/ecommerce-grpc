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
)

type cartHandleGrpc struct {
	pb.UnimplementedCartServiceServer
	cartService service.CartService
	mapping     protomapper.CartProtoMapper
}

func NewCartHandleGrpc(
	cartService service.CartService,
	mapping protomapper.CartProtoMapper,
) *cartHandleGrpc {
	return &cartHandleGrpc{
		cartService: cartService,
		mapping:     mapping,
	}
}

func (s *cartHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllCartRequest) (*pb.ApiResponsePaginationCart, error) {
	user_id := int(request.GetUserId())
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCarts{
		UserID:   user_id,
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cartItems, totalRecords, err := s.cartService.FindAll(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationCart(paginationMeta, "success", "Successfully fetched cart items", cartItems)
	return so, nil
}

func (s *cartHandleGrpc) Create(ctx context.Context, request *pb.CreateCartRequest) (*pb.ApiResponseCart, error) {
	req := &requests.CreateCartRequest{
		Quantity:  int(request.GetQuantity()),
		ProductID: int(request.GetProductId()),
		UserID:    int(request.GetUserId()),
	}

	if err := req.Validate(); err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Unable to create new cart. Please check your input.",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	cartItem, err := s.cartService.CreateCart(req)

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

	so := s.mapping.ToProtoResponseCart("success", "Successfully added item to cart", cartItem)
	return so, nil
}

func (s *cartHandleGrpc) Delete(ctx context.Context, request *pb.FindByIdCartRequest) (*pb.ApiResponseCartDelete, error) {
	if request.GetId() == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Cart ID parameter cannot be empty and must be a positive number",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	_, err := s.cartService.DeletePermanent(int(request.GetId()))

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

	so := s.mapping.ToProtoResponseCartDelete("success", "Successfully removed item from cart")

	return so, nil
}

func (s *cartHandleGrpc) DeleteAll(ctx context.Context, req *pb.DeleteCartRequest) (*pb.ApiResponseCartAll, error) {
	cartIDs := make([]int, len(req.GetCartIds()))
	for i, id := range req.GetCartIds() {
		cartIDs[i] = int(id)
	}

	deleteRequest := &requests.DeleteCartRequest{
		CartIds: cartIDs,
	}

	_, err := s.cartService.DeleteAllPermanently(deleteRequest)

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

	so := s.mapping.ToProtoResponseCartAll("success", "Successfully cleared cart")
	return so, nil
}
