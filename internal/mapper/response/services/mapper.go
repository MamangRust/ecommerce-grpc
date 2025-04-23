package response_service

type ResponseServiceMapper struct {
	RoleResponseMapper             RoleResponseMapper
	RefreshTokenResponseMapper     RefreshTokenResponseMapper
	UserResponseMapper             UserResponseMapper
	CategoryResponseMapper         CategoryResponseMapper
	MerchantResponseMapper         MerchantResponseMapper
	OrderResponseMapper            OrderResponseMapper
	OrderItemResponseMapper        OrderItemResponseMapper
	ProductResponseMapper          ProductResponseMapper
	TransactionResponseMapper      TransactionResponseMapper
	CartResponseMapper             CartResponseMapper
	ReviewResponseMapper           ReviewResponseMapper
	ShippingAddressResponseMapper  ShippingAddressResponseMapper
	SliderResponseMapper           SliderResponseMapper
	BannerResponseMapper           BannerResponseMapper
	MerchantAwardResponseMapper    MerchantAwardResponseMapper
	MerchantBusinessResponseMapper MerchantBusinessResponseMapper
	MerchantDetailResponseMapper   MerchantDetailResponseMapper
	MerchantPolicyResponseMapper   MerchantPolicyResponseMapper
	ReviewDetailResponeMapper      ReviewDetailResponeMapper
}

func NewResponseServiceMapper() *ResponseServiceMapper {
	return &ResponseServiceMapper{
		UserResponseMapper:             NewUserResponseMapper(),
		RefreshTokenResponseMapper:     NewRefreshTokenResponseMapper(),
		RoleResponseMapper:             NewRoleResponseMapper(),
		CategoryResponseMapper:         NewCategoryResponseMapper(),
		MerchantResponseMapper:         NewMerchantResponseMapper(),
		OrderResponseMapper:            NewOrderResponseMapper(),
		OrderItemResponseMapper:        NewOrderItemResponseMapper(),
		ProductResponseMapper:          NewProductResponseMapper(),
		TransactionResponseMapper:      NewTransactionResponseMapper(),
		CartResponseMapper:             NewCartResponseMapper(),
		ShippingAddressResponseMapper:  NewShippingAddressResponseMapper(),
		SliderResponseMapper:           NewSliderResponseMapper(),
		ReviewResponseMapper:           NewReviewResponseMapper(),
		BannerResponseMapper:           NewBannerResponseMapper(),
		MerchantAwardResponseMapper:    NewMerchantAwardResponseMapper(),
		MerchantBusinessResponseMapper: NewMerchantBusinessResponseMapper(),
		MerchantDetailResponseMapper:   NewMerchantDetailResponseMapper(),
		MerchantPolicyResponseMapper:   NewMerchantPolicyResponseMapper(),
		ReviewDetailResponeMapper:      NewReviewDetailResponseMapper(),
	}
}
