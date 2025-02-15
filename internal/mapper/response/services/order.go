package response_service

import (
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/response"
)

type orderResponseMapper struct {
}

func NewOrderResponseMapper() *orderResponseMapper {
	return &orderResponseMapper{}
}

func (s *orderResponseMapper) ToOrderResponse(order *record.OrderRecord) *response.OrderResponse {
	return &response.OrderResponse{
		ID:         order.ID,
		MerchantID: order.MerchantID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (s *orderResponseMapper) ToOrdersResponse(orders []*record.OrderRecord) []*response.OrderResponse {
	var responses []*response.OrderResponse

	for _, order := range orders {
		responses = append(responses, s.ToOrderResponse(order))
	}

	return responses
}

func (s *orderResponseMapper) ToOrderResponseDeleteAt(order *record.OrderRecord) *response.OrderResponseDeleteAt {
	return &response.OrderResponseDeleteAt{
		ID:         order.ID,
		MerchantID: order.MerchantID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		DeleteAt:   *order.DeletedAt,
	}
}

func (s *orderResponseMapper) ToOrdersResponseDeleteAt(orders []*record.OrderRecord) []*response.OrderResponseDeleteAt {
	var responses []*response.OrderResponseDeleteAt

	for _, order := range orders {
		responses = append(responses, s.ToOrderResponseDeleteAt(order))
	}

	return responses
}
