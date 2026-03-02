package service

import (
	"ecommerce/internal/cache"
	auth_cache "ecommerce/internal/cache/auth"
	banner_cache "ecommerce/internal/cache/banner"
	cart_cache "ecommerce/internal/cache/cart"
	category_cache "ecommerce/internal/cache/category"
	merchant_cache "ecommerce/internal/cache/merchant"
	merchantawards_cache "ecommerce/internal/cache/merchant_awards"
	merchantbusiness_cache "ecommerce/internal/cache/merchant_business"
	merchantdetail_cache "ecommerce/internal/cache/merchant_detail"
	merchantpolicies_cache "ecommerce/internal/cache/merchant_policies"
	order_cache "ecommerce/internal/cache/order"
	orderitem_cache "ecommerce/internal/cache/order_item"
	product_cache "ecommerce/internal/cache/product"
	review_cache "ecommerce/internal/cache/review"
	reviewdetail_cache "ecommerce/internal/cache/review_detail"
	role_cache "ecommerce/internal/cache/role"
	shippingaddress_cache "ecommerce/internal/cache/shipping_address"
	slider_cache "ecommerce/internal/cache/slider"
	transaction_cache "ecommerce/internal/cache/transaction"
	user_cache "ecommerce/internal/cache/user"
	"ecommerce/internal/repository"
	"ecommerce/pkg/auth"
	"ecommerce/pkg/hash"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
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
	Cache        *cache.CacheStore
}

func NewService(deps Deps) *Service {
	observability, _ := observability.NewObservability("grpc-server", deps.Logger)

	cacheAuth := auth_cache.NewMencache(deps.Cache)
	cacheUser := user_cache.NewUserMencache(deps.Cache)
	cacheRole := role_cache.NewRoleMencache(deps.Cache)
	cacheCategory := category_cache.NewCategoryMencache(deps.Cache)
	cacheMerchant := merchant_cache.NewMerchantMencache(deps.Cache)
	cacheOrder := order_cache.OrderNewMencache(deps.Cache)
	cacheOrderItem := orderitem_cache.NewOrderItemMencache(deps.Cache)
	cacheProduct := product_cache.NewProductMencache(deps.Cache)
	cacheCart := cart_cache.NewCartMencache(deps.Cache)
	cacheShipping := shippingaddress_cache.NewShippingAddressMencache(deps.Cache)
	cacheSlider := slider_cache.NewSliderMencache(deps.Cache)
	cacheReview := review_cache.NewReviewMencache(deps.Cache)
	cacheReviewDetail := reviewdetail_cache.NewReviewDetailMencache(deps.Cache)
	cacheMerchantAward := merchantawards_cache.NewMerchantAward(deps.Cache)
	cacheMerchantPolicies := merchantpolicies_cache.NewMerchantPoliciesMencache(deps.Cache)
	cacheMerchantBusiness := merchantbusiness_cache.NewMerchantBusinessMencache(deps.Cache)
	cacheMerchantDetail := merchantdetail_cache.NewMerchantDetailMencache(deps.Cache)
	cacheBanner := banner_cache.NewBannerMencache(deps.Cache)
	cacheTransaction := transaction_cache.NewTransactionMencache(deps.Cache)

	return &Service{
		Auth: NewAuthService(AuthServiceDeps{
			AuthRepo:         deps.Repositories.User,
			RefreshTokenRepo: deps.Repositories.RefreshToken,
			RoleRepo:         deps.Repositories.Role,
			UserRoleRepo:     deps.Repositories.UserRole,
			Hash:             deps.Hash,
			Token:            deps.Token,
			Logger:           deps.Logger,
			Tracer:           observability,
			CacheIdentity:    cacheAuth.IdentityCache,
			CacheLogin:       cacheAuth.LoginCache,
		}),
		User: NewUserService(UserServiceDeps{
			UserRepo:      deps.Repositories.User,
			Logger:        deps.Logger,
			Cache:         cacheUser,
			Observability: observability,
			Hashing:       deps.Hash,
		}),

		Role: NewRoleService(RoleServiceDeps{
			RoleRepo:      deps.Repositories.Role,
			Logger:        deps.Logger,
			Cache:         cacheRole,
			Observability: observability,
		}),
		Category: NewCategoryService(CategoryServiceDeps{
			CategoryRepository: deps.Repositories.Category,
			Logger:             deps.Logger,
			Cache:              cacheCategory,
			Observability:      observability,
		}),
		Product: NewProductService(ProductServiceDeps{
			CategoryRepository: deps.Repositories.Category,
			MerchantRepository: deps.Repositories.Merchant,
			ProductRepository:  deps.Repositories.Product,
			Logger:             deps.Logger,
			Cache:              cacheProduct,
			Observability:      observability,
		}),
		Merchant: NewMerchantService(MerchantServiceDeps{
			MerchantRepository: deps.Repositories.Merchant,
			Logger:             deps.Logger,
			Cache:              cacheMerchant,
			Observability:      observability,
		}),

		MerchantAward: NewMerchantAwardService(MerchantAwardServiceDeps{
			MerchantAwardRepository: deps.Repositories.MerchantAward,
			Logger:                  deps.Logger,
			Cache:                   cacheMerchantAward,
			Observability:           observability,
		}),

		MerchantBusiness: NewMerchantBusinessService(MerchantBusinessServiceDeps{
			MerchantBusinessRepository: deps.Repositories.MerchantBusiness,
			Logger:                     deps.Logger,
			Cache:                      cacheMerchantBusiness,
			Observability:              observability,
		}),

		MerchantDetail: NewMerchantDetailService(MerchantDetailServiceDeps{
			MerchantDetailRepository:     deps.Repositories.MerchantDetail,
			MerchantSocialLinkRepository: deps.Repositories.MerchantSocialLink,
			Logger:                       deps.Logger,
			Cache:                        cacheMerchantDetail,
			Observability:                observability,
		}),

		MerchantPolicies: NewMerchantPoliciesService(MerchantPoliciesServiceDeps{
			MerchantPoliciesRepository: deps.Repositories.MerchantPolicies,
			Logger:                     deps.Logger,
			Cache:                      cacheMerchantPolicies,
			Observability:              observability,
		}),

		OrderItem: NewOrderItemService(OrderItemServiceDeps{
			OrderItemRepository: deps.Repositories.OrderItem,
			Logger:              deps.Logger,
			Cache:               cacheOrderItem,
			Observability:       observability,
		}),

		Order: NewOrderService(OrderServiceDeps{
			OrderRepository:     deps.Repositories.Order,
			OrderItemRepository: deps.Repositories.OrderItem,
			ProductRepository:   deps.Repositories.Product,
			UserRepository:      deps.Repositories.User,
			MerchantRepository:  deps.Repositories.Merchant,
			ShippingRepository:  deps.Repositories.Shipping,
			Logger:              deps.Logger,
			Cache:               cacheOrder,
			Observability:       observability,
		}),

		Transaction: NewTransactionService(TransactionServiceDeps{
			MerchantRepository:    deps.Repositories.Merchant,
			TransactionRepository: deps.Repositories.Transaction,
			OrderRepository:       deps.Repositories.Order,
			OrderItemRepository:   deps.Repositories.OrderItem,
			ShippingRepository:    deps.Repositories.Shipping,
			Logger:                deps.Logger,
			Cache:                 cacheTransaction,
			Observability:         observability,
		}),

		Cart: NewCartService(CartServiceDeps{
			ProductRepository: deps.Repositories.Product,
			UserRepository:    deps.Repositories.User,
			CartRepository:    deps.Repositories.Cart,
			Logger:            deps.Logger,
			Cache:             cacheCart,
			Observability:     observability,
		}),
		Shipping: NewShippingAddressService(ShippingAddressServiceDeps{
			ShippingRepository: deps.Repositories.Shipping,
			Logger:             deps.Logger,
			Cache:              cacheShipping,
			Observability:      observability,
		}),

		Slider: NewSliderService(SliderServiceDeps{
			SliderRepository: deps.Repositories.Slider,
			Logger:           deps.Logger,
			Cache:            cacheSlider,
			Observability:    observability,
		}),

		Banner: NewBannerService(BannerServiceDeps{
			BannerRepository: deps.Repositories.Banner,
			Logger:           deps.Logger,
			Cache:            cacheBanner,
			Observability:    observability,
		}),

		Review: NewReviewService(ReviewServiceDeps{
			ReviewRepository:  deps.Repositories.Review,
			ProductRepository: deps.Repositories.Product,
			UserRepository:    deps.Repositories.User,
			Logger:            deps.Logger,
			Cache:             cacheReview,
			Observability:     observability,
		}),

		ReviewDetail: NewReviewDetailService(ReviewDetailServiceDeps{
			ReviewDetailRepository: deps.Repositories.ReviewDetail,
			Logger:                 deps.Logger,
			Cache:                  cacheReviewDetail,
			Observability:          observability,
		}),
	}
}
