package requests

import "github.com/go-playground/validator/v10"

type CreateTransactionRequest struct {
	OrderID       int    `json:"order_id" validate:"required"`
	MerchantID    int    `json:"merchant_id" validate:"required"`
	PaymentMethod string `json:"payment_method" validate:"required"`
	Amount        int    `json:"amount" validate:"required"`
	ChangeAmount  int    `json:"change_amount" validate:"required"`
	PaymentStatus string `json:"payment_status" validate:"required"`
}

type UpdateTransactionRequest struct {
	TransactionID int    `json:"transaction_id" validate:"required"`
	OrderID       int    `json:"order_id" validate:"required"`
	MerchantID    int    `json:"merchant_id" validate:"required"`
	PaymentMethod string `json:"payment_method" validate:"required"`
	Amount        int    `json:"amount" validate:"required"`
	ChangeAmount  int    `json:"change_amount" validate:"required"`
	PaymentStatus string `json:"payment_status" validate:"required"`
}

func (r *CreateTransactionRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateTransactionRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
