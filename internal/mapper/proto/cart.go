package protomapper

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"
)

type cartProtoMapper struct {
}

func NewCartProtoMapper() *cartProtoMapper {
	return &cartProtoMapper{}
}

func (c *cartProtoMapper) ToProtoResponseCartDelete(status string, message string) *pb.ApiResponseCartDelete {
	return &pb.ApiResponseCartDelete{
		Status:  status,
		Message: message,
	}
}

func (c *cartProtoMapper) ToProtoResponseCartAll(status string, message string) *pb.ApiResponseCartAll {
	return &pb.ApiResponseCartAll{
		Status:  status,
		Message: message,
	}
}

func (c *cartProtoMapper) mapResponseCart(cart *response.CartResponse) *pb.CartResponse {
	return &pb.CartResponse{
		Id:        int32(cart.ID),
		UserId:    int32(cart.UserID),
		ProductId: int32(cart.ProductID),
		Name:      cart.Name,
		Price:     int32(cart.Price),
		Image:     cart.Image,
		Quantity:  int32(cart.Quantity),
		Weight:    int32(cart.Weight),
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
}

func (c *cartProtoMapper) mapResponsesCart(carts []*response.CartResponse) []*pb.CartResponse {
	var mappedCarts []*pb.CartResponse

	for _, cart := range carts {
		mappedCarts = append(mappedCarts, c.mapResponseCart(cart))
	}

	return mappedCarts
}

func (c *cartProtoMapper) ToProtoResponseCart(status string, message string, categories *response.CartResponse) *pb.ApiResponseCart {
	return &pb.ApiResponseCart{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCart(categories),
	}
}

func (c *cartProtoMapper) ToProtoResponsePaginationCart(pagination *pb.PaginationMeta, status string, message string, categories []*response.CartResponse) *pb.ApiResponsePaginationCart {
	return &pb.ApiResponsePaginationCart{
		Status:     status,
		Message:    message,
		Data:       c.mapResponsesCart(categories),
		Pagination: mapPaginationMeta(pagination),
	}
}
