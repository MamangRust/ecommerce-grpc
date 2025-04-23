package protomapper

type ProtoMapper struct {
	AuthProtoMapper             AuthProtoMapper
	RoleProtoMapper             RoleProtoMapper
	UserProtoMapper             UserProtoMapper
	CategoryProtoMapper         CategoryProtoMapper
	MerchantProtoMapper         MerchantProtoMapper
	OrderItemProtoMapper        OrderItemProtoMapper
	OrderProtoMapper            OrderProtoMapper
	ProductProtoMapper          ProductProtoMapper
	TransactionProtoMapper      TransactionProtoMapper
	CartProtoMapper             CartProtoMapper
	ReviewProtoMapper           ReviewProtoMapper
	SliderProtoMapper           SliderProtoMapper
	ShippingProtoMapper         ShippingAddresProtoMapper
	BannerProtoMapper           BannerProtoMapper
	MerchantAwardProtoMapper    MerchantAwardProtoMapper
	MerchantBusinessProtoMapper MerchantBusinessProtoMapper
	MerchantDetailProtoMapper   MerchantDetailProtoMapper
	MerchantPolicyProtoMapper   MerchantPolicyProtoMapper
	ReviewDetailProtoMapper     ReviewDetailProtoMapper
}

func NewProtoMapper() *ProtoMapper {
	return &ProtoMapper{
		AuthProtoMapper:             NewAuthProtoMapper(),
		RoleProtoMapper:             NewRoleProtoMapper(),
		UserProtoMapper:             NewUserProtoMapper(),
		CategoryProtoMapper:         NewCategoryProtoMapper(),
		MerchantProtoMapper:         NewMerchantProtoMaper(),
		OrderItemProtoMapper:        NewOrderItemProtoMapper(),
		OrderProtoMapper:            NewOrderProtoMapper(),
		ProductProtoMapper:          NewProductProtoMapper(),
		TransactionProtoMapper:      NewTransactionProtoMapper(),
		CartProtoMapper:             NewCartProtoMapper(),
		ReviewProtoMapper:           NewReviewProtoMapper(),
		SliderProtoMapper:           NewSliderProtoMapper(),
		ShippingProtoMapper:         NewShippingAddressProtoMapper(),
		BannerProtoMapper:           NewBannerProtoMaper(),
		MerchantAwardProtoMapper:    NewMerchantAwardProtoMapper(),
		MerchantBusinessProtoMapper: NewMerchantBusinessProtoMapper(),
		MerchantDetailProtoMapper:   NewMerchantDetailProtoMapper(),
		MerchantPolicyProtoMapper:   NewMerchantPolicyProtoMapper(),
		ReviewDetailProtoMapper:     NewReviewDetailProtoMapper(),
	}
}
