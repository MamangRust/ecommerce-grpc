package requests

import "github.com/go-playground/validator/v10"

type CreateUserRequest struct {
	FirstName       string `json:"firstname" validate:"required,alpha"`
	LastName        string `json:"lastname" validate:"required,alpha"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type UpdateUserRequest struct {
	UserID          int    `json:"user_id" validate:"required,min=1"`
	FirstName       string `json:"firstname" validate:"required,alpha"`
	LastName        string `json:"lastname" validate:"required,alpha"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

func (r *CreateUserRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateUserRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
