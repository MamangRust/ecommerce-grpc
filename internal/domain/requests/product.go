package requests

import "github.com/go-playground/validator/v10"

type CreateProductRequest struct {
	MerchantID   int    `json:"merchant_id" validate:"required"`
	CategoryID   int    `json:"category_id" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description" validate:"required"`
	Price        int    `json:"price" validate:"required"`
	CountInStock int    `json:"count_in_stock" validate:"required"`
	Brand        string `json:"brand" validate:"required"`
	Weight       int    `json:"weight" validate:"required"`
	Rating       int    `json:"rating" validate:"required"`
	SlugProduct  string `json:"slug_product" validate:"required"`
	ImageProduct string `json:"image_product" validate:"required"`
	Barcode      string `json:"barcode" validate:"required"`
}

type UpdateProductRequest struct {
	ProductID    int    `json:"product_id" validate:"required"`
	MerchantID   int    `json:"merchant_id" validate:"required"`
	CategoryID   int    `json:"category_id" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description" validate:"required"`
	Price        int    `json:"price" validate:"required"`
	CountInStock int    `json:"count_in_stock" validate:"required"`
	Brand        string `json:"brand" validate:"required"`
	Weight       int    `json:"weight" validate:"required"`
	Rating       int    `json:"rating" validate:"required"`
	SlugProduct  string `json:"slug_product" validate:"required"`
	ImageProduct string `json:"image_product" validate:"required"`
	Barcode      string `json:"barcode" validate:"required"`
}

func (r *CreateProductRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateProductRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
