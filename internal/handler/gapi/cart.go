package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	protomapper "ecommerce/internal/mapper/proto"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors/cart_errors"
	"math"
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
		return nil, response.ToGrpcErrorFromErrorResponse(err)
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
		return nil, cart_errors.ErrGrpcValidateCreateCart
	}

	cartItem, err := s.cartService.CreateCart(req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCart("success", "Successfully added item to cart", cartItem)
	return so, nil
}

func (s *cartHandleGrpc) Delete(ctx context.Context, request *pb.FindByIdCartRequest) (*pb.ApiResponseCartDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, cart_errors.ErrGrpcCartInvalidId
	}

	_, err := s.cartService.DeletePermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
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
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCartAll("success", "Successfully cleared cart")
	return so, nil
}
