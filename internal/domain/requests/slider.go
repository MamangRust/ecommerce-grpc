package requests

import "github.com/go-playground/validator/v10"

type FindAllSlider struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type CreateSliderRequest struct {
	Nama     string `json:"nama"`
	FilePath string `json:"file_path"`
}

type UpdateSliderRequest struct {
	ID       *int   `json:"id"`
	Nama     string `json:"nama"`
	FilePath string `json:"file_path"`
}

func (l *CreateSliderRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(l)

	if err != nil {
		return err
	}

	return nil
}

func (l *UpdateSliderRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(l)

	if err != nil {
		return err
	}

	return nil
}
