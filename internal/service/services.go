package service

import (
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/auth"
	"ecommerce/pkg/hash"
	"ecommerce/pkg/logger"
)

type Service struct {
	Auth             AuthService
	User             UserService
	Role             RoleService
	Category         CategoryService
	Merchant         MerchantService
	OrderItem        OrderItemService
	Order            OrderService
	Product          ProductService
	Transaction      TransactionService
	Review           ReviewService
	Cart             CartService
	Shipping         ShippingAddressService
	Slider           SliderService
	Banner           BannerService
	MerchantAward    MerchantAwardService
	MerchantBusiness MerchantBusinessService
	MerchantDetail   MerchantDetailService
	MerchantPolicies MerchantPoliciesService
	ReviewDetail     ReviewDetailService
}

type Deps struct {
	Repositories *repository.Repositories
	Token        auth.TokenManager
	Hash         hash.HashPassword
	Logger       logger.LoggerInterface
	Mapper       response_service.ResponseServiceMapper
}

func NewService(deps Deps) *Service {
	return &Service{
		Auth:             NewAuthService(deps.Repositories.User, deps.Repositories.RefreshToken, deps.Repositories.Role, deps.Repositories.UserRole, deps.Hash, deps.Token, deps.Logger, deps.Mapper.UserResponseMapper),
		User:             NewUserService(deps.Repositories.User, deps.Logger, deps.Mapper.UserResponseMapper, deps.Hash),
		Category:         NewCategoryService(deps.Repositories.Category, deps.Logger, deps.Mapper.CategoryResponseMapper),
		Merchant:         NewMerchantService(deps.Repositories.Merchant, deps.Logger, deps.Mapper.MerchantResponseMapper),
		OrderItem:        NewOrderItemService(deps.Repositories.OrderItem, deps.Logger, deps.Mapper.OrderItemResponseMapper),
		Order:            NewOrderService(deps.Repositories.Order, deps.Repositories.OrderItem, deps.Repositories.User, deps.Repositories.Merchant, deps.Repositories.Product, deps.Repositories.Shipping, deps.Logger, deps.Mapper.OrderResponseMapper),
		Product:          NewProductService(deps.Repositories.Category, deps.Repositories.Merchant, deps.Repositories.Product, deps.Logger, deps.Mapper.ProductResponseMapper),
		Transaction:      NewTransactionService(deps.Repositories.Merchant, deps.Repositories.Transaction, deps.Repositories.Order, deps.Repositories.OrderItem, deps.Logger, deps.Mapper.TransactionResponseMapper),
		Cart:             NewCartService(deps.Repositories.Product, deps.Repositories.User, deps.Repositories.Cart, deps.Logger, deps.Mapper.CartResponseMapper),
		Shipping:         NewShippingAddressService(deps.Repositories.Shipping, deps.Logger, deps.Mapper.ShippingAddressResponseMapper),
		Slider:           NewSliderService(deps.Repositories.Slider, deps.Logger, deps.Mapper.SliderResponseMapper),
		Review:           NewReviewService(deps.Repositories.Review, deps.Repositories.Product, deps.Repositories.User),
		Banner:           NewBannerService(deps.Repositories.Banner, deps.Logger, deps.Mapper.BannerResponseMapper),
		MerchantAward:    NewMerchantAwardService(deps.Repositories.MerchantAward, deps.Logger, deps.Mapper.MerchantAwardResponseMapper),
		MerchantBusiness: NewMerchantBusinessService(deps.Repositories.MerchantBusiness, deps.Logger, deps.Mapper.MerchantBusinessResponseMapper),
		MerchantDetail:   NewMerchantDetailService(deps.Repositories.MerchantDetail, deps.Repositories.MerchantSocialLink, deps.Logger, deps.Mapper.MerchantDetailResponseMapper),
		MerchantPolicies: NewMerchantPoliciesService(deps.Repositories.MerchantPolicies, deps.Logger, deps.Mapper.MerchantPolicyResponseMapper),
		ReviewDetail:     NewReviewDetailService(deps.Repositories.ReviewDetail, deps.Logger, deps.Mapper.ReviewDetailResponeMapper),
	}
}
