package requests

import "github.com/go-playground/validator/v10"

type CreateMerchantRequest struct {
	UserID       int    `json:"user_id" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description" validate:"required"`
	Address      string `json:"address" validate:"required"`
	ContactEmail string `json:"contact_email" validate:"required,email"`
	ContactPhone string `json:"contact_phone" validate:"required"`
	Status       string `json:"status" validate:"required"`
}

type UpdateMerchantRequest struct {
	MerchantID   int    `json:"merchant_id" validate:"required"`
	UserID       int    `json:"user_id" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description" validate:"required"`
	Address      string `json:"address" validate:"required"`
	ContactEmail string `json:"contact_email" validate:"required,email"`
	ContactPhone string `json:"contact_phone" validate:"required"`
	Status       string `json:"status" validate:"required"`
}

func (r *CreateMerchantRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateMerchantRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
