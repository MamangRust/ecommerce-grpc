package requests

import "github.com/go-playground/validator/v10"

type CreateOrderRecordRequest struct {
	MerchantID int `json:"merchant_id" validate:"required"`
	UserID     int `json:"user_id" validate:"required"`
	TotalPrice int `json:"total_price"`
}

type UpdateOrderRecordRequest struct {
	OrderID    int `json:"order_id" validate:"required"`
	MerchantID int `json:"merchant_id" validate:"required"`
	UserID     int `json:"user_id" validate:"required"`
	TotalPrice int `json:"total_price" validate:"required"`
}

type CreateOrderRequest struct {
	MerchantID      int                          `json:"merchant_id" validate:"required"`
	UserID          int                          `json:"user_id" validate:"required"`
	TotalPrice      int                          `json:"total_price" validate:"required"`
	Items           []CreateOrderItemRequest     `json:"items" validate:"required"`
	ShippingAddress CreateShippingAddressRequest `json:"shipping_address"`
}

type UpdateOrderRequest struct {
	OrderID         int                          `json:"order_id" validate:"required"`
	UserID          int                          `json:"user_id" validate:"required"`
	TotalPrice      int                          `json:"total_price" validate:"required"`
	Items           []UpdateOrderItemRequest     `json:"items" validate:"required"`
	ShippingAddress UpdateShippingAddressRequest `json:"shipping_address"`
}

type CreateOrderItemRequest struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required"`
	Price     int `json:"price" validate:"required"`
}

type UpdateOrderItemRequest struct {
	OrderItemID int `json:"order_item_id" validate:"required"`
	ProductID   int `json:"product_id" validate:"required"`
	Quantity    int `json:"quantity" validate:"required"`
	Price       int `json:"price" validate:"required"`
}

func (r *CreateOrderRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateOrderRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
