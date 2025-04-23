package repository

import (
	"context"
	recordmapper "ecommerce/internal/mapper/record"
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
	DB           *db.Queries
	Ctx          context.Context
	MapperRecord *recordmapper.RecordMapper
}

func NewRepositories(deps Deps) *Repositories {
	return &Repositories{
		User:               NewUserRepository(deps.DB, deps.Ctx, deps.MapperRecord.UserRecordMapper),
		Role:               NewRoleRepository(deps.DB, deps.Ctx, deps.MapperRecord.RoleRecordMapper),
		UserRole:           NewUserRoleRepository(deps.DB, deps.Ctx, deps.MapperRecord.UserRoleRecordMapper),
		Category:           NewCategoryRepository(deps.DB, deps.Ctx, deps.MapperRecord.CategoryRecordMapper),
		RefreshToken:       NewRefreshTokenRepository(deps.DB, deps.Ctx, deps.MapperRecord.RefreshTokenRecordMapper),
		Product:            NewProductRepository(deps.DB, deps.Ctx, deps.MapperRecord.ProductRecordMapper),
		Merchant:           NewMerchantRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantRecordMapper),
		OrderItem:          NewOrderItemRepository(deps.DB, deps.Ctx, deps.MapperRecord.OrderItemRecordMapper),
		Order:              NewOrderRepository(deps.DB, deps.Ctx, deps.MapperRecord.OrderRecordMapper),
		Transaction:        NewTransactionRepository(deps.DB, deps.Ctx, deps.MapperRecord.TransactionRecordMapper),
		Cart:               NewCartRepository(deps.DB, deps.Ctx, deps.MapperRecord.CartRecordMapping),
		Shipping:           NewShippingAddressRepository(deps.DB, deps.Ctx, deps.MapperRecord.ShippingAddressMapping),
		Review:             NewReviewRepository(deps.DB, deps.Ctx, deps.MapperRecord.ReviewRecordMapping),
		Slider:             NewSliderRepository(deps.DB, deps.Ctx, deps.MapperRecord.SliderMapping),
		Banner:             NewBannerRepository(deps.DB, deps.Ctx, deps.MapperRecord.BannerMapping),
		MerchantAward:      NewMerchantAwardRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantAwardMapping),
		MerchantBusiness:   NewMerchantBusinessRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantBusinessMapping),
		MerchantDetail:     NewMerchantDetailRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantDetailMapping),
		MerchantSocialLink: NewMerchantSocialMediaLinkRepository(deps.DB, deps.Ctx),
		MerchantPolicies:   NewMerchantPolicyRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantPolicyMapping),
		ReviewDetail:       NewReviewDetailRepository(deps.DB, deps.Ctx, deps.MapperRecord.ReviewDetailRecordMapping),
	}
}
