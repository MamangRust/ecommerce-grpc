package requests

import "github.com/go-playground/validator/v10"

type CreateMerchantSocialFormData struct {
	MerchantDetailID *int
	Platform         string
	Url              string
}

type CreateMerchantDetailFormData struct {
	MerchantID       int
	DisplayName      string
	CoverImageUrl    string
	LogoUrl          string
	ShortDescription string
	WebsiteUrl       string
	SocialLinks      []CreateMerchantSocialFormData
}

type UpdateMerchantSocialFormData struct {
	ID               int
	MerchantDetailID *int
	Platform         string
	Url              string
}

type UpdateMerchantDetailFormData struct {
	DisplayName      string
	CoverImageUrl    string
	LogoUrl          string
	ShortDescription string
	WebsiteUrl       string
	SocialLinks      []UpdateMerchantSocialFormData
}

type CreateMerchantDetailRequest struct {
	MerchantID       int                            `json:"merchant_id" validate:"required"`
	DisplayName      string                         `json:"display_name" validate:"required"`
	CoverImageUrl    string                         `json:"cover_image_url" validate:"required,url"`
	LogoUrl          string                         `json:"logo_url" validate:"required,url"`
	ShortDescription string                         `json:"short_description" validate:"required"`
	WebsiteUrl       string                         `json:"website_url" validate:"omitempty,url"`
	SocialLink       []*CreateMerchantSocialRequest `json:"social_links" validate:"required,dive"`
}

type UpdateMerchantDetailRequest struct {
	MerchantDetailID *int                           `json:"merchant_detail_id" validate:"required"`
	DisplayName      string                         `json:"display_name" validate:"required"`
	CoverImageUrl    string                         `json:"cover_image_url" validate:"required,url"`
	LogoUrl          string                         `json:"logo_url" validate:"required,url"`
	ShortDescription string                         `json:"short_description" validate:"required"`
	WebsiteUrl       string                         `json:"website_url" validate:"omitempty,url"`
	SocialLink       []*UpdateMerchantSocialRequest `json:"social_links" validate:"required,dive"`
}

func (r *CreateMerchantDetailRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateMerchantDetailRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
