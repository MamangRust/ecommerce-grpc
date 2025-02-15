package requests

import "github.com/go-playground/validator/v10"

type CreateCategoryRequest struct {
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description" validate:"required"`
	SlugCategory  string `json:"slug_category" validate:"required"`
	ImageCategory string `json:"image_category" validate:"required"`
}

type UpdateCategoryRequest struct {
	CategoryID    int    `json:"category_id" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description" validate:"required"`
	SlugCategory  string `json:"slug_category" validate:"required"`
	ImageCategory string `json:"image_category" validate:"required"`
}

func (r *CreateCategoryRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateCategoryRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
