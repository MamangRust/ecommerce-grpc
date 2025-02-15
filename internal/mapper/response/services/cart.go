package response_service

import (
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/response"
)

type cartResponseMapper struct {
}

func NewCartResponseMapper() *cartResponseMapper {
	return &cartResponseMapper{}
}

func (s *cartResponseMapper) ToCartResponse(cart *record.CartRecord) *response.CartResponse {
	return &response.CartResponse{
		ID:        cart.ID,
		UserID:    int(cart.UserID),
		ProductID: int(cart.ProductID),
		Name:      cart.Name,
		Price:     cart.Price,
		Image:     cart.Image,
		Quantity:  int(cart.Quantity),
		Weight:    int(cart.Weight),
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
}

func (s *cartResponseMapper) ToCartsResponse(users []*record.CartRecord) []*response.CartResponse {
	var responses []*response.CartResponse

	for _, user := range users {
		responses = append(responses, s.ToCartResponse(user))
	}

	return responses
}
