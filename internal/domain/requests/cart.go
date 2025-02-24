package requests

import "github.com/go-playground/validator/v10"

type CartCreateRecord struct {
	ProductID    int    `json:"product_id"`
	UserID       int    `json:"user_id"`
	Name         string `json:"name"`
	Price        int    `json:"price"`
	ImageProduct string `json:"image_product"`
	Quantity     int    `json:"quantity"`
	Weight       int    `json:"weight"`
}

type CreateCartRequest struct {
	Quantity  int `json:"quantity" validate:"required,gt=0"`
	ProductID int `json:"product_id" validate:"required"`
	UserID    int `json:"user_id,omitempty"`
}

type DeleteCartRequest struct {
	CartIds []int `json:"cart_ids"`
}

func (l *CartCreateRecord) Validate() error {
	validate := validator.New()

	err := validate.Struct(l)

	if err != nil {
		return err
	}

	return nil
}
func (l *CreateCartRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(l)

	if err != nil {
		return err
	}

	return nil
}

func (l *DeleteCartRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(l)

	if err != nil {
		return err
	}

	return nil
}
