package response_service

import (
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/response"
)

type categoryResponseMapper struct {
}

func NewCategoryResponseMapper() *categoryResponseMapper {
	return &categoryResponseMapper{}
}

func (s *categoryResponseMapper) ToCategoryResponse(category *record.CategoriesRecord) *response.CategoryResponse {
	return &response.CategoryResponse{
		ID:            category.ID,
		Name:          category.Name,
		Description:   category.Description,
		SlugCategory:  category.SlugCategory,
		ImageCategory: category.ImageCategory,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
	}
}

func (s *categoryResponseMapper) ToCategorysResponse(categories []*record.CategoriesRecord) []*response.CategoryResponse {
	var responses []*response.CategoryResponse

	for _, category := range categories {
		responses = append(responses, s.ToCategoryResponse(category))
	}

	return responses
}

func (s *categoryResponseMapper) ToCategoryResponseDeleteAt(category *record.CategoriesRecord) *response.CategoryResponseDeleteAt {
	return &response.CategoryResponseDeleteAt{
		ID:            category.ID,
		Name:          category.Name,
		Description:   category.Description,
		SlugCategory:  category.SlugCategory,
		ImageCategory: category.ImageCategory,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
		DeletedAt:     *category.DeletedAt,
	}
}

func (s *categoryResponseMapper) ToCategorysResponseDeleteAt(categories []*record.CategoriesRecord) []*response.CategoryResponseDeleteAt {
	var responses []*response.CategoryResponseDeleteAt

	for _, category := range categories {
		responses = append(responses, s.ToCategoryResponseDeleteAt(category))
	}

	return responses
}
