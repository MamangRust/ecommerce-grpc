package api

import (
	"ecommerce/internal/cache"
	auth_cache "ecommerce/internal/cache/api/auth"
	banner_cache "ecommerce/internal/cache/api/banner"
	cart_cache "ecommerce/internal/cache/api/cart"
	category_cache "ecommerce/internal/cache/api/category"
	merchant_cache "ecommerce/internal/cache/api/merchant"
	merchantawards_cache "ecommerce/internal/cache/api/merchant_awards"
	merchantbusiness_cache "ecommerce/internal/cache/api/merchant_business"
	merchantdetail_cache "ecommerce/internal/cache/api/merchant_detail"
	merchantpolicies_cache "ecommerce/internal/cache/api/merchant_policies"
	order_cache "ecommerce/internal/cache/api/order"
	orderitem_cache "ecommerce/internal/cache/api/order_item"
	product_cache "ecommerce/internal/cache/api/product"
	review_cache "ecommerce/internal/cache/api/review"
	reviewdetail_cache "ecommerce/internal/cache/api/review_detail"
	role_cache "ecommerce/internal/cache/api/role"
	shippingaddress_cache "ecommerce/internal/cache/api/shipping_address"
	slider_cache "ecommerce/internal/cache/api/slider"
	transaction_cache "ecommerce/internal/cache/api/transaction"
	user_cache "ecommerce/internal/cache/api/user"
	response_api "ecommerce/internal/mapper"
	"ecommerce/internal/pb"
	"ecommerce/pkg/auth"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"ecommerce/pkg/upload_image"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type Deps struct {
	Conn    *grpc.ClientConn
	Token   auth.TokenManager
	E       *echo.Echo
	Logger  logger.LoggerInterface
	Mapping *response_api.ResponseApiMapper
	Image   upload_image.ImageUploads
	Cache   *cache.CacheStore
}

func NewHandler(deps Deps) {
	observability, _ := observability.NewObservability("client", deps.Logger)

	apiHandler := errors.NewApiHandler(observability, deps.Logger)

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

	clientAuth := pb.NewAuthServiceClient(deps.Conn)
	clientRole := pb.NewRoleServiceClient(deps.Conn)
	clientUser := pb.NewUserServiceClient(deps.Conn)
	clientCategory := pb.NewCategoryServiceClient(deps.Conn)
	clientMerchant := pb.NewMerchantServiceClient(deps.Conn)
	clientOrderItem := pb.NewOrderItemServiceClient(deps.Conn)
	clientOrder := pb.NewOrderServiceClient(deps.Conn)
	clientProduct := pb.NewProductServiceClient(deps.Conn)
	clientTransaction := pb.NewTransactionServiceClient(deps.Conn)
	clientCart := pb.NewCartServiceClient(deps.Conn)
	clientReview := pb.NewReviewServiceClient(deps.Conn)
	clientSlider := pb.NewSliderServiceClient(deps.Conn)
	clientShipping := pb.NewShippingServiceClient(deps.Conn)
	clientBanner := pb.NewBannerServiceClient(deps.Conn)
	clientMerchantAward := pb.NewMerchantAwardServiceClient(deps.Conn)
	clientMerchantBusiness := pb.NewMerchantBusinessServiceClient(deps.Conn)
	clientMerchantDetail := pb.NewMerchantDetailServiceClient(deps.Conn)
	clientMerchantPolicy := pb.NewMerchantPoliciesServiceClient(deps.Conn)
	clientReviewDetail := pb.NewReviewDetailServiceClient(deps.Conn)

	NewHandlerAuth(
		deps.E,
		clientAuth,
		deps.Logger,
		deps.Mapping.AuthResponseMapper,
		apiHandler,
		cacheAuth,
	)

	NewHandlerRole(
		deps.E,
		clientRole,
		deps.Logger,
		deps.Mapping.RoleResponseMapper,
		apiHandler,
		cacheRole,
	)

	NewHandlerUser(
		deps.E,
		clientUser,
		deps.Logger,
		deps.Mapping.UserResponseMapper,
		apiHandler,
		cacheUser,
	)

	NewHandlerCategory(
		deps.E,
		clientCategory,
		deps.Logger,
		deps.Mapping.CategoryResponseMapper,
		deps.Image,
		apiHandler,
		cacheCategory,
	)

	NewHandlerMerchant(
		deps.E,
		clientMerchant,
		deps.Logger,
		deps.Mapping.MerchantResponseMapper,
		apiHandler,
		cacheMerchant,
	)

	NewHandlerOrder(
		deps.E,
		clientOrder,
		deps.Logger,
		deps.Mapping.OrderResponseMapper,
		apiHandler,
		cacheOrder,
	)

	NewHandlerOrderItem(
		deps.E,
		clientOrderItem,
		deps.Logger,
		deps.Mapping.OrderItemResponseMapper,
		apiHandler,
		cacheOrderItem,
	)

	NewHandlerProduct(
		deps.E,
		clientProduct,
		deps.Logger,
		deps.Mapping.ProductResponseMapper,
		deps.Image,
		apiHandler,
		cacheProduct,
	)

	NewHandlerCart(
		deps.E,
		clientCart,
		deps.Logger,
		deps.Mapping.CartResponseMapper,
		apiHandler,
		cacheCart,
	)

	NewHandlerShippingAddress(
		deps.E,
		clientShipping,
		deps.Logger,
		deps.Mapping.ShippingAddressResponseMapper,
		apiHandler,
		cacheShipping,
	)

	NewHandlerSlider(
		deps.E,
		clientSlider,
		deps.Logger,
		deps.Mapping.SliderMapper,
		deps.Image,
		apiHandler,
		cacheSlider,
	)

	NewHandlerReview(
		deps.E,
		clientReview,
		deps.Logger,
		deps.Mapping.ReviewMapper,
		apiHandler,
		cacheReview,
	)

	NewHandlerReviewDetail(
		deps.E,
		clientReviewDetail,
		deps.Logger,
		deps.Mapping.ReviewDetailResponseMapper,
		deps.Mapping.ReviewMapper,
		deps.Image,
		apiHandler,
		cacheReviewDetail,
	)

	NewHandleBanner(
		deps.E,
		clientBanner,
		deps.Logger,
		deps.Mapping.BannerResponseMapper,
		apiHandler,
		cacheBanner,
	)

	NewHandlerMerchantAward(
		deps.E,
		clientMerchantAward,
		deps.Logger,
		deps.Mapping.MerchantAwardResponseMapper,
		deps.Mapping.MerchantResponseMapper,
		apiHandler,
		cacheMerchantAward,
	)

	NewHandlerMerchantBusiness(
		deps.E,
		clientMerchantBusiness,
		deps.Logger,
		deps.Mapping.MerchantBusinessMapper,
		deps.Mapping.MerchantResponseMapper,
		apiHandler,
		cacheMerchantBusiness,
	)

	NewHandlerMerchantPolicies(
		deps.E,
		clientMerchantPolicy,
		deps.Logger,
		deps.Mapping.MerchantPolicyResponseMapper,
		deps.Mapping.MerchantResponseMapper,
		apiHandler,
		cacheMerchantPolicies,
	)

	NewHandlerMerchantDetail(
		deps.E,
		clientMerchantDetail,
		deps.Logger,
		deps.Mapping.MerchantDetailResponseMapper,
		deps.Mapping.MerchantResponseMapper,
		deps.Image,
		apiHandler,
		cacheMerchantDetail,
	)

	NewHandlerTransaction(
		deps.E,
		clientTransaction,
		deps.Logger,
		deps.Mapping.TransactionResponseMapper,
		apiHandler,
		cacheTransaction,
	)
}
