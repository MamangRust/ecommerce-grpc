package response

type MerchantDetailResponse struct {
	ID               int                                `json:"id"`
	MerchantID       int                                `json:"merchant_id"`
	DisplayName      string                             `json:"display_name"`
	CoverImageUrl    string                             `json:"cover_image_url"`
	LogoUrl          string                             `json:"logo_url"`
	ShortDescription string                             `json:"short_description"`
	WebsiteUrl       string                             `json:"website_url"`
	SocialMediaLinks []*MerchantSocialMediaLinkResponse `json:"social_media_links"`
	CreatedAt        string                             `json:"created_at"`
	UpdatedAt        string                             `json:"updated_at"`
}

type MerchantDetailResponseDeleteAt struct {
	ID               int                                `json:"id"`
	MerchantID       int                                `json:"merchant_id"`
	DisplayName      string                             `json:"display_name"`
	CoverImageUrl    string                             `json:"cover_image_url"`
	LogoUrl          string                             `json:"logo_url"`
	ShortDescription string                             `json:"short_description"`
	WebsiteUrl       string                             `json:"website_url"`
	SocialMediaLinks []*MerchantSocialMediaLinkResponse `json:"social_media_links"`
	CreatedAt        string                             `json:"created_at"`
	UpdatedAt        string                             `json:"updated_at"`
	DeletedAt        *string                            `json:"deleted_at"`
}

type MerchantDetailCoreResponse struct {
	ID               int    `json:"id"`
	MerchantID       int    `json:"merchant_id"`
	DisplayName      string `json:"display_name"`
	CoverImageUrl    string `json:"cover_image_url"`
	LogoUrl          string `json:"logo_url"`
	ShortDescription string `json:"short_description"`
	WebsiteUrl       string `json:"website_url"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type MerchantDetailCoreResponseDeleteAt struct {
	ID               int     `json:"id"`
	MerchantID       int     `json:"merchant_id"`
	DisplayName      string  `json:"display_name"`
	CoverImageUrl    string  `json:"cover_image_url"`
	LogoUrl          string  `json:"logo_url"`
	ShortDescription string  `json:"short_description"`
	WebsiteUrl       string  `json:"website_url"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	DeletedAt        *string `json:"deleted_at"`
}

type MerchantSocialMediaLinkResponse struct {
	ID       int    `json:"id"`
	Platform string `json:"platform"`
	Url      string `json:"url"`
}

type ApiResponseMerchantDetail struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Data    *MerchantDetailResponse `json:"data"`
}

type ApiResponseMerchantDetailDeleteAt struct {
	Status  string                          `json:"status"`
	Message string                          `json:"message"`
	Data    *MerchantDetailResponseDeleteAt `json:"data"`
}

type ApiResponseMerchantDetailCore struct {
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Data    *MerchantDetailCoreResponse `json:"data"`
}

type ApiResponseMerchantDetailDeleteAtCore struct {
	Status  string                              `json:"status"`
	Message string                              `json:"message"`
	Data    *MerchantDetailCoreResponseDeleteAt `json:"data"`
}

type ApiResponseMerchantDetailRelation struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Data    *MerchantDetailResponse `json:"data"`
}

type ApiResponsesMerchantDetail struct {
	Status  string                    `json:"status"`
	Message string                    `json:"message"`
	Data    []*MerchantDetailResponse `json:"data"`
}

type ApiResponsePaginationMerchantDetailDeleteAt struct {
	Status     string                            `json:"status"`
	Message    string                            `json:"message"`
	Data       []*MerchantDetailResponseDeleteAt `json:"data"`
	Pagination PaginationMeta                    `json:"pagination"`
}

type ApiResponsePaginationMerchantDetail struct {
	Status     string                    `json:"status"`
	Message    string                    `json:"message"`
	Data       []*MerchantDetailResponse `json:"data"`
	Pagination PaginationMeta            `json:"pagination"`
}
