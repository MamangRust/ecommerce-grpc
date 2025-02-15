package requests

import "github.com/go-playground/validator/v10"

type CreateOrderItemRecordRequest struct {
	OrderID   int `json:"order_id" validate:"required"`
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required"`
	Price     int `json:"price" validate:"required"`
}

type UpdateOrderItemRecordRequest struct {
	OrderItemID int `json:"order_item_id" validate:"required"`
	OrderID     int `json:"order_id" validate:"required"`
	ProductID   int `json:"product_id" validate:"required"`
	Quantity    int `json:"quantity" validate:"required"`
	Price       int `json:"price" validate:"required"`
}

func (r *CreateOrderItemRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateOrderItemRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
