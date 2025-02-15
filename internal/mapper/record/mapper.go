package recordmapper

type RecordMapper struct {
	UserRecordMapper         UserRecordMapping
	RoleRecordMapper         RoleRecordMapping
	UserRoleRecordMapper     UserRoleRecordMapping
	RefreshTokenRecordMapper RefreshTokenRecordMapping
	CategoryRecordMapper     CategoryRecordMapper
	MerchantRecordMapper     MerchantRecordMapping
	OrderItemRecordMapper    OrderItemRecordMapping
	OrderRecordMapper        OrderRecordMapping
	ProductRecordMapper      ProductRecordMapping
	TransactionRecordMapper  TransactionRecordMapping
	CartRecordMapping        CartRecordMapping
	ReviewRecordMapping      ReviewRecordMapping
	ShippingAddressMapping   ShippingAddressMapping
	SliderMapping            SliderMapping
}

func NewRecordMapper() *RecordMapper {
	return &RecordMapper{
		UserRecordMapper:         NewUserRecordMapper(),
		RoleRecordMapper:         NewRoleRecordMapper(),
		UserRoleRecordMapper:     NewUserRoleRecordMapper(),
		RefreshTokenRecordMapper: NewRefreshTokenRecordMapper(),
		CategoryRecordMapper:     NewCategoryRecordMapper(),
		MerchantRecordMapper:     NewMerchantRecordMapper(),
		OrderItemRecordMapper:    NewOrderItemRecordMapper(),
		OrderRecordMapper:        NewOrderRecordMapper(),
		ProductRecordMapper:      NewProductRecordMapper(),
		CartRecordMapping:        NewCartRecordMapper(),
		TransactionRecordMapper:  NewTransactionRecordMapper(),
		ReviewRecordMapping:      NewReviewRecordMapper(),
		ShippingAddressMapping:   NewShippingAddressRecordMapper(),
		SliderMapping:            NewSliderRecordMapper(),
	}
}
