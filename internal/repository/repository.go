package repository

import (
	db "ecommerce/pkg/database/schema"
)

type Repositories struct {
	User               UserRepository
	Role               RoleRepository
	UserRole           UserRoleRepository
	Category           CategoryRepository
	RefreshToken       RefreshTokenRepository
	Product            ProductRepository
	Merchant           MerchantRepository
	OrderItem          OrderItemRepository
	Order              OrderRepository
	Transaction        TransactionRepository
	Cart               CartRepository
	Shipping           ShippingAddressRepository
	Review             ReviewRepository
	Slider             SliderRepository
	Banner             BannerRepository
	MerchantAward      MerchantAwardRepository
	MerchantBusiness   MerchantBusinessRepository
	MerchantDetail     MerchantDetailRepository
	MerchantSocialLink MerchantSocialLinkRepository
	MerchantPolicies   MerchantPoliciesRepository
	ReviewDetail       ReviewDetailRepository
}

type Deps struct {
	DB *db.Queries
}

func NewRepositories(db *db.Queries) *Repositories {
	return &Repositories{
		User:               NewUserRepository(db),
		Role:               NewRoleRepository(db),
		UserRole:           NewUserRoleRepository(db),
		Cart:               NewCartRepository(db),
		RefreshToken:       NewRefreshTokenRepository(db),
		Category:           NewCategoryRepository(db),
		Product:            NewProductRepository(db),
		Merchant:           NewMerchantRepository(db),
		OrderItem:          NewOrderItemRepository(db),
		Order:              NewOrderRepository(db),
		Transaction:        NewTransactionRepository(db),
		Shipping:           NewShippingAddressRepository(db),
		Review:             NewReviewRepository(db),
		Slider:             NewSliderRepository(db),
		Banner:             NewBannerRepository(db),
		MerchantAward:      NewMerchantAwardRepository(db),
		MerchantBusiness:   NewMerchantBusinessRepository(db),
		MerchantDetail:     NewMerchantDetailRepository(db),
		MerchantSocialLink: NewMerchantSocialMediaLinkRepository(db),
		MerchantPolicies:   NewMerchantPolicyRepository(db),
		ReviewDetail:       NewReviewDetailRepository(db),
	}
}
