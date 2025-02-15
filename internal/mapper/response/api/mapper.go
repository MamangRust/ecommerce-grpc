package response_api

type ResponseApiMapper struct {
	AuthResponseMapper            AuthResponseMapper
	RoleResponseMapper            RoleResponseMapper
	UserResponseMapper            UserResponseMapper
	CategoryResponseMapper        CategoryResponseMapper
	MerchantResponseMapper        MerchantResponseMapper
	OrderItemResponseMapper       OrderItemResponseMapper
	OrderResponseMapper           OrderResponseMapper
	ProductResponseMapper         ProductResponseMapper
	TransactionResponseMapper     TransactionResponseMapper
	CartResponseMapper            CartResponseMapper
	ReviewMapper                  ReviewResponseMapper
	SliderMapper                  SliderResponseMapper
	ShippingAddressResponseMapper ShippingAddressResponseMapper
}

func NewResponseApiMapper() *ResponseApiMapper {
	return &ResponseApiMapper{
		AuthResponseMapper:            NewAuthResponseMapper(),
		UserResponseMapper:            NewUserResponseMapper(),
		RoleResponseMapper:            NewRoleResponseMapper(),
		CategoryResponseMapper:        NewCategoryResponseMapper(),
		MerchantResponseMapper:        NewMerchantResponseMapper(),
		OrderItemResponseMapper:       NewOrderItemResponseMapper(),
		OrderResponseMapper:           NewOrderResponseMapper(),
		ProductResponseMapper:         NewProductResponseMapper(),
		TransactionResponseMapper:     NewTransactionResponseMapper(),
		CartResponseMapper:            NewCartResponseMapper(),
		ReviewMapper:                  NewReviewResponseMapper(),
		SliderMapper:                  NewSliderResponseMapper(),
		ShippingAddressResponseMapper: NewShippingAddressResponseMapper(),
	}
}
