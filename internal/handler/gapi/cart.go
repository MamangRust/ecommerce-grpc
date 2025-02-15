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
	cart_id := request.GetUserId()
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	cartItems, totalRecords, err := s.cartService.FindAll(int(cart_id), page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch cart items",
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
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
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to add item to cart: " + err.Error(),
		})
	}

	cartItem, err := s.cartService.CreateCart(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to add item to cart",
		})
	}

	so := s.mapping.ToProtoResponseCart("success", "Successfully added item to cart", cartItem)
	return so, nil
}

func (s *cartHandleGrpc) Delete(ctx context.Context, request *pb.FindByIdCartRequest) (*pb.ApiResponseCartDelete, error) {
	if request.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid cart item id",
		})
	}

	_, err := s.cartService.DeletePermanent(int(request.GetId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to remove item from cart: ",
		})
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
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to clear cart",
		})
	}

	so := s.mapping.ToProtoResponseCartAll("success", "Successfully cleared cart")
	return so, nil
}
