package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/errors/cart_errors"
	"math"
)

type cartHandleGrpc struct {
	pb.UnimplementedCartServiceServer
	cartService service.CartService
}

func NewCartHandleGrpc(
	cartService service.CartService,
) *cartHandleGrpc {
	return &cartHandleGrpc{
		cartService: cartService,
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

	cartItems, totalRecords, err := s.cartService.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoCartItems := make([]*pb.CartResponse, len(cartItems))
	for i, cartItem := range cartItems {
		protoCartItems[i] = &pb.CartResponse{
			Id:        int32(cartItem.CartID),
			UserId:    int32(cartItem.UserID),
			ProductId: int32(cartItem.ProductID),
			Name:      cartItem.Name,
			Price:     int32(cartItem.Price),
			Image:     cartItem.Image,
			Quantity:  int32(cartItem.Quantity),
			Weight:    int32(cartItem.Weight),
			CreatedAt: cartItem.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: cartItem.UpdatedAt.Time.Format("2006-01-02"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationCart{
		Status:     "success",
		Message:    "Successfully fetched cart items",
		Data:       protoCartItems,
		Pagination: paginationMeta,
	}, nil
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

	cartItem, err := s.cartService.CreateCart(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoCartItem := &pb.CartResponse{
		Id:        int32(cartItem.CartID),
		UserId:    int32(cartItem.UserID),
		ProductId: int32(cartItem.ProductID),
		Name:      cartItem.Name,
		Price:     int32(cartItem.Price),
		Image:     cartItem.Image,
		Quantity:  int32(cartItem.Quantity),
		Weight:    int32(cartItem.Weight),
		CreatedAt: cartItem.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: cartItem.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseCart{
		Status:  "success",
		Message: "Successfully added item to cart",
		Data:    protoCartItem,
	}, nil
}

func (s *cartHandleGrpc) Delete(ctx context.Context, request *pb.FindByIdCartRequest) (*pb.ApiResponseCartDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, cart_errors.ErrGrpcCartInvalidId
	}

	_, err := s.cartService.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseCartDelete{
		Status:  "success",
		Message: "Successfully removed item from cart",
	}, nil
}

func (s *cartHandleGrpc) DeleteAll(ctx context.Context, req *pb.DeleteCartRequest) (*pb.ApiResponseCartAll, error) {
	cartIDs := make([]int, len(req.GetCartIds()))
	for i, id := range req.GetCartIds() {
		cartIDs[i] = int(id)
	}

	deleteRequest := &requests.DeleteCartRequest{
		CartIds: cartIDs,
	}

	_, err := s.cartService.DeleteAllPermanently(ctx, deleteRequest)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseCartAll{
		Status:  "success",
		Message: "Successfully cleared cart",
	}, nil
}
