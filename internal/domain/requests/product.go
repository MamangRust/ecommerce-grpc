package requests

import "github.com/go-playground/validator/v10"

type FindAllProduct struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type FindAllProductByMerchant struct {
	CategoryID int    `json:"category_id" validate:"required"`
	MerchantID int    `json:"merchant_id" validate:"required"`
	Search     string `json:"search" validate:"required"`
	MinPrice   *int   `json:"min_price"`
	MaxPrice   *int   `json:"max_price"`
	Page       int    `json:"page" validate:"min=1"`
	PageSize   int    `json:"page_size" validate:"min=1,max=100"`
}

type FindAllProductByCategory struct {
	CategoryName string `json:"category_name" validate:"required"`
	Search       string `json:"search" validate:"required"`
	MinPrice     *int   `json:"min_price"`
	MaxPrice     *int   `json:"max_price"`
	Page         int    `json:"page" validate:"min=1"`
	PageSize     int    `json:"page_size" validate:"min=1,max=100"`
}

type CreateProductRequest struct {
	MerchantID   int     `json:"merchant_id" validate:"required"`
	CategoryID   int     `json:"category_id" validate:"required"`
	Name         string  `json:"name" validate:"required"`
	Description  string  `json:"description" validate:"required"`
	Price        int     `json:"price" validate:"required"`
	CountInStock int     `json:"count_in_stock" validate:"required"`
	Brand        string  `json:"brand" validate:"required"`
	Weight       int     `json:"weight" validate:"required"`
	Rating       *int    `json:"rating" validate:"required"`
	SlugProduct  *string `json:"slug_product" validate:"required"`
	ImageProduct string  `json:"image_product" validate:"required"`
}

type UpdateProductRequest struct {
	ProductID    *int    `json:"product_id" `
	MerchantID   int     `json:"merchant_id" validate:"required"`
	CategoryID   int     `json:"category_id" validate:"required"`
	Name         string  `json:"name" validate:"required"`
	Description  string  `json:"description" validate:"required"`
	Price        int     `json:"price" validate:"required"`
	CountInStock int     `json:"count_in_stock" validate:"required"`
	Brand        string  `json:"brand" validate:"required"`
	Weight       int     `json:"weight" validate:"required"`
	Rating       *int    `json:"rating" validate:"required"`
	SlugProduct  *string `json:"slug_product" validate:"required"`
	ImageProduct string  `json:"image_product" validate:"required"`
}

type ProductFormData struct {
	MerchantID   int
	CategoryID   int
	Name         string
	Description  string
	Price        int
	CountInStock int
	Brand        string
	Weight       int
	ImagePath    string
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
