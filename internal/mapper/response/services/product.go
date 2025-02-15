package response_service

import (
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/response"
)

type productResponseMapper struct {
}

func NewProductResponseMapper() *productResponseMapper {
	return &productResponseMapper{}
}

func (s *productResponseMapper) ToProductResponse(product *record.ProductRecord) *response.ProductResponse {
	return &response.ProductResponse{
		ID:           product.ID,
		MerchantID:   product.MerchantID,
		CategoryID:   product.CategoryID,
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		CountInStock: product.CountInStock,
		Brand:        product.Brand,
		Weight:       product.Weight,
		Rating:       product.Rating,
		SlugProduct:  product.SlugProduct,
		ImageProduct: product.ImageProduct,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
}

func (s *productResponseMapper) ToProductsResponse(products []*record.ProductRecord) []*response.ProductResponse {
	var responses []*response.ProductResponse

	for _, product := range products {
		responses = append(responses, s.ToProductResponse(product))
	}

	return responses
}

func (s *productResponseMapper) ToProductResponseDeleteAt(product *record.ProductRecord) *response.ProductResponseDeleteAt {
	return &response.ProductResponseDeleteAt{
		ID:           product.ID,
		MerchantID:   product.MerchantID,
		CategoryID:   product.CategoryID,
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		CountInStock: product.CountInStock,
		Brand:        product.Brand,
		Weight:       product.Weight,
		Rating:       product.Rating,
		SlugProduct:  product.SlugProduct,
		ImageProduct: product.ImageProduct,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
		DeleteAt:     *product.DeletedAt,
	}
}

func (s *productResponseMapper) ToProductsResponseDeleteAt(products []*record.ProductRecord) []*response.ProductResponseDeleteAt {
	var responses []*response.ProductResponseDeleteAt

	for _, product := range products {
		responses = append(responses, s.ToProductResponseDeleteAt(product))
	}

	return responses
}
