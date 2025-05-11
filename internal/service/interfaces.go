package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type AuthService interface {
	Register(request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse)
	Login(request *requests.AuthRequest) (*response.TokenResponse, *response.ErrorResponse)
	RefreshToken(token string) (*response.TokenResponse, *response.ErrorResponse)
	GetMe(token string) (*response.UserResponse, *response.ErrorResponse)
}

type RoleService interface {
	FindAll(req *requests.FindAllRole) ([]*response.RoleResponse, *int, *response.ErrorResponse)
	FindByActiveRole(req *requests.FindAllRole) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashedRole(req *requests.FindAllRole) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(role_id int) (*response.RoleResponse, *response.ErrorResponse)
	FindByUserId(id int) ([]*response.RoleResponse, *response.ErrorResponse)
	CreateRole(request *requests.CreateRoleRequest) (*response.RoleResponse, *response.ErrorResponse)
	UpdateRole(request *requests.UpdateRoleRequest) (*response.RoleResponse, *response.ErrorResponse)
	TrashedRole(role_id int) (*response.RoleResponse, *response.ErrorResponse)
	RestoreRole(role_id int) (*response.RoleResponse, *response.ErrorResponse)
	DeleteRolePermanent(role_id int) (bool, *response.ErrorResponse)

	RestoreAllRole() (bool, *response.ErrorResponse)
	DeleteAllRolePermanent() (bool, *response.ErrorResponse)
}

type UserService interface {
	FindAll(req *requests.FindAllUsers) ([]*response.UserResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse)
	FindByID(id int) (*response.UserResponse, *response.ErrorResponse)
	Create(request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse)
	Update(request *requests.UpdateUserRequest) (*response.UserResponse, *response.ErrorResponse)
	Trashed(user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse)
	Restore(user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse)
	DeletePermanent(user_id int) (bool, *response.ErrorResponse)

	RestoreAll() (bool, *response.ErrorResponse)
	DeleteAllPermanent() (bool, *response.ErrorResponse)
}

type BannerService interface {
	FindAll(req *requests.FindAllBanner) ([]*response.BannerResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(BannerID int) (*response.BannerResponse, *response.ErrorResponse)
	CreateBanner(req *requests.CreateBannerRequest) (*response.BannerResponse, *response.ErrorResponse)
	UpdateBanner(req *requests.UpdateBannerRequest) (*response.BannerResponse, *response.ErrorResponse)
	TrashedBanner(BannerID int) (*response.BannerResponseDeleteAt, *response.ErrorResponse)
	RestoreBanner(BannerID int) (*response.BannerResponseDeleteAt, *response.ErrorResponse)
	DeleteBannerPermanent(BannerID int) (bool, *response.ErrorResponse)
	RestoreAllBanner() (bool, *response.ErrorResponse)
	DeleteAllBannerPermanent() (bool, *response.ErrorResponse)
}

type CategoryService interface {
	FindMonthlyTotalPrice(req *requests.MonthTotalPrice) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse)
	FindYearlyTotalPrice(year int) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse)
	FindMonthlyTotalPriceById(req *requests.MonthTotalPriceCategory) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse)
	FindYearlyTotalPriceById(req *requests.YearTotalPriceCategory) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse)
	FindMonthlyTotalPriceByMerchant(req *requests.MonthTotalPriceMerchant) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse)
	FindYearlyTotalPriceByMerchant(req *requests.YearTotalPriceMerchant) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse)

	FindMonthPrice(year int) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse)
	FindYearPrice(year int) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse)
	FindMonthPriceByMerchant(req *requests.MonthPriceMerchant) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse)
	FindYearPriceByMerchant(req *requests.YearPriceMerchant) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse)
	FindMonthPriceById(req *requests.MonthPriceId) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse)
	FindYearPriceById(req *requests.YearPriceId) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse)

	FindAll(req *requests.FindAllCategory) ([]*response.CategoryResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(category_id int) (*response.CategoryResponse, *response.ErrorResponse)
	CreateCategory(req *requests.CreateCategoryRequest) (*response.CategoryResponse, *response.ErrorResponse)
	UpdateCategory(req *requests.UpdateCategoryRequest) (*response.CategoryResponse, *response.ErrorResponse)
	TrashedCategory(category_id int) (*response.CategoryResponseDeleteAt, *response.ErrorResponse)
	RestoreCategory(categoryID int) (*response.CategoryResponseDeleteAt, *response.ErrorResponse)
	DeleteCategoryPermanent(categoryID int) (bool, *response.ErrorResponse)
	RestoreAllCategories() (bool, *response.ErrorResponse)
	DeleteAllCategoriesPermanent() (bool, *response.ErrorResponse)
}

type MerchantService interface {
	FindAll(req *requests.FindAllMerchant) ([]*response.MerchantResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(merchantID int) (*response.MerchantResponse, *response.ErrorResponse)
	CreateMerchant(req *requests.CreateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse)
	UpdateMerchant(req *requests.UpdateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse)
	TrashedMerchant(merchantID int) (*response.MerchantResponseDeleteAt, *response.ErrorResponse)
	RestoreMerchant(merchantID int) (*response.MerchantResponseDeleteAt, *response.ErrorResponse)
	DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse)
	RestoreAllMerchant() (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent() (bool, *response.ErrorResponse)
}

type MerchantAwardService interface {
	FindAll(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(merchantID int) (*response.MerchantAwardResponse, *response.ErrorResponse)
	CreateMerchant(req *requests.CreateMerchantCertificationOrAwardRequest) (*response.MerchantAwardResponse, *response.ErrorResponse)
	UpdateMerchant(req *requests.UpdateMerchantCertificationOrAwardRequest) (*response.MerchantAwardResponse, *response.ErrorResponse)
	TrashedMerchant(merchantID int) (*response.MerchantAwardResponseDeleteAt, *response.ErrorResponse)
	RestoreMerchant(merchantID int) (*response.MerchantAwardResponseDeleteAt, *response.ErrorResponse)
	DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse)
	RestoreAllMerchant() (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent() (bool, *response.ErrorResponse)
}

type MerchantBusinessService interface {
	FindAll(req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(merchantID int) (*response.MerchantBusinessResponse, *response.ErrorResponse)
	CreateMerchant(req *requests.CreateMerchantBusinessInformationRequest) (*response.MerchantBusinessResponse, *response.ErrorResponse)
	UpdateMerchant(req *requests.UpdateMerchantBusinessInformationRequest) (*response.MerchantBusinessResponse, *response.ErrorResponse)
	TrashedMerchant(merchantID int) (*response.MerchantBusinessResponseDeleteAt, *response.ErrorResponse)
	RestoreMerchant(merchantID int) (*response.MerchantBusinessResponseDeleteAt, *response.ErrorResponse)
	DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse)
	RestoreAllMerchant() (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent() (bool, *response.ErrorResponse)
}

type MerchantDetailService interface {
	FindAll(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(merchantID int) (*response.MerchantDetailResponse, *response.ErrorResponse)
	CreateMerchant(req *requests.CreateMerchantDetailRequest) (*response.MerchantDetailResponse, *response.ErrorResponse)
	UpdateMerchant(req *requests.UpdateMerchantDetailRequest) (*response.MerchantDetailResponse, *response.ErrorResponse)
	TrashedMerchant(merchantID int) (*response.MerchantDetailResponseDeleteAt, *response.ErrorResponse)
	RestoreMerchant(merchantID int) (
		*response.MerchantDetailResponseDeleteAt, *response.ErrorResponse)
	DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse)
	RestoreAllMerchant() (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent() (bool, *response.ErrorResponse)
}

type MerchantPoliciesService interface {
	FindAll(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(merchantID int) (*response.MerchantPoliciesResponse, *response.ErrorResponse)
	CreateMerchant(req *requests.CreateMerchantPolicyRequest) (*response.MerchantPoliciesResponse, *response.ErrorResponse)
	UpdateMerchant(req *requests.UpdateMerchantPolicyRequest) (*response.MerchantPoliciesResponse, *response.ErrorResponse)
	TrashedMerchant(merchantID int) (*response.MerchantPoliciesResponseDeleteAt, *response.ErrorResponse)
	RestoreMerchant(merchantID int) (*response.MerchantPoliciesResponseDeleteAt, *response.ErrorResponse)
	DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse)
	RestoreAllMerchant() (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent() (bool, *response.ErrorResponse)
}

type OrderItemService interface {
	FindAllOrderItems(req *requests.FindAllOrderItems) ([]*response.OrderItemResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse)
	FindOrderItemByOrder(orderID int) ([]*response.OrderItemResponse, *response.ErrorResponse)
}

type OrderService interface {
	FindMonthlyTotalRevenue(req *requests.MonthTotalRevenue) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse)
	FindYearlyTotalRevenue(year int) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse)
	FindMonthlyTotalRevenueById(req *requests.MonthTotalRevenueOrder) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse)
	FindYearlyTotalRevenueById(req *requests.YearTotalRevenueOrder) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse)
	FindMonthlyTotalRevenueByMerchant(req *requests.MonthTotalRevenueMerchant) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse)
	FindYearlyTotalRevenueByMerchant(req *requests.YearTotalRevenueMerchant) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse)

	FindMonthlyOrder(year int) ([]*response.OrderMonthlyResponse, *response.ErrorResponse)
	FindYearlyOrder(year int) ([]*response.OrderYearlyResponse, *response.ErrorResponse)
	FindMonthlyOrderByMerchant(req *requests.MonthOrderMerchant) ([]*response.OrderMonthlyResponse, *response.ErrorResponse)
	FindYearlyOrderByMerchant(req *requests.YearOrderMerchant) ([]*response.OrderYearlyResponse, *response.ErrorResponse)

	FindAll(req *requests.FindAllOrder) ([]*response.OrderResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(order_id int) (*response.OrderResponse, *response.ErrorResponse)
	CreateOrder(req *requests.CreateOrderRequest) (*response.OrderResponse, *response.ErrorResponse)
	UpdateOrder(req *requests.UpdateOrderRequest) (*response.OrderResponse, *response.ErrorResponse)
	TrashedOrder(order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse)
	RestoreOrder(order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse)
	DeleteOrderPermanent(order_id int) (bool, *response.ErrorResponse)
	RestoreAllOrder() (bool, *response.ErrorResponse)
	DeleteAllOrderPermanent() (bool, *response.ErrorResponse)
}

type ProductService interface {
	FindAll(req *requests.FindAllProduct) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	FindByMerchant(req *requests.FindAllProductByMerchant) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	FindByCategory(req *requests.FindAllProductByCategory) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(productID int) (*response.ProductResponse, *response.ErrorResponse)
	CreateProduct(req *requests.CreateProductRequest) (*response.ProductResponse, *response.ErrorResponse)
	UpdateProduct(req *requests.UpdateProductRequest) (*response.ProductResponse, *response.ErrorResponse)
	TrashProduct(productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse)
	RestoreProduct(productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse)
	DeleteProductPermanent(productID int) (bool, *response.ErrorResponse)
	RestoreAllProducts() (bool, *response.ErrorResponse)
	DeleteAllProductsPermanent() (bool, *response.ErrorResponse)
}

type TransactionService interface {
	FindMonthlyAmountSuccess(req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse)
	FindYearlyAmountSuccess(year int) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse)
	FindMonthlyAmountFailed(req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse)
	FindYearlyAmountFailed(year int) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse)

	FindMonthlyAmountSuccessByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse)
	FindYearlyAmountSuccessByMerchant(req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse)
	FindMonthlyAmountFailedByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse)
	FindYearlyAmountFailedByMerchant(req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse)

	FindMonthlyMethodSuccess(req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse)
	FindYearlyMethodSuccess(year int) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse)
	FindMonthlyMethodByMerchantSuccess(req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse)
	FindYearlyMethodByMerchantSuccess(req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse)

	FindMonthlyMethodFailed(req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse)
	FindYearlyMethodFailed(year int) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse)
	FindMonthlyMethodByMerchantFailed(req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse)
	FindYearlyMethodByMerchantFailed(req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse)

	FindAllTransactions(req *requests.FindAllTransaction) ([]*response.TransactionResponse, *int, *response.ErrorResponse)
	FindByMerchant(req *requests.FindAllTransactionByMerchant) ([]*response.TransactionResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(transactionID int) (*response.TransactionResponse, *response.ErrorResponse)
	FindByOrderId(orderID int) (*response.TransactionResponse, *response.ErrorResponse)
	CreateTransaction(req *requests.CreateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse)
	UpdateTransaction(req *requests.UpdateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse)
	TrashedTransaction(transaction_id int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse)
	RestoreTransaction(transaction_id int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse)
	DeleteTransactionPermanently(transactionID int) (bool, *response.ErrorResponse)
	RestoreAllTransactions() (bool, *response.ErrorResponse)
	DeleteAllTransactionPermanent() (bool, *response.ErrorResponse)
}

type CartService interface {
	FindAll(req *requests.FindAllCarts) ([]*response.CartResponse, *int, *response.ErrorResponse)
	CreateCart(req *requests.CreateCartRequest) (*response.CartResponse, *response.ErrorResponse)
	DeletePermanent(cart_id int) (bool, *response.ErrorResponse)
	DeleteAllPermanently(req *requests.DeleteCartRequest) (bool, *response.ErrorResponse)
}

type ReviewService interface {
	FindAllReviews(req *requests.FindAllReview) ([]*response.ReviewResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, *response.ErrorResponse)
	FindByProduct(req *requests.FindAllReviewByProduct) ([]*response.ReviewsDetailResponse, *int, *response.ErrorResponse)
	FindByMerchant(req *requests.FindAllReviewByMerchant) ([]*response.ReviewsDetailResponse, *int, *response.ErrorResponse)
	CreateReview(req *requests.CreateReviewRequest) (*response.ReviewResponse, *response.ErrorResponse)
	UpdateReview(req *requests.UpdateReviewRequest) (*response.ReviewResponse, *response.ErrorResponse)
	TrashedReview(reviewID int) (*response.ReviewResponseDeleteAt, *response.ErrorResponse)
	RestoreReview(reviewID int) (*response.ReviewResponseDeleteAt, *response.ErrorResponse)
	DeleteReviewPermanent(reviewID int) (bool, *response.ErrorResponse)
	RestoreAllReviews() (bool, *response.ErrorResponse)
	DeleteAllReviewsPermanent() (bool, *response.ErrorResponse)
}

type ReviewDetailService interface {
	FindAll(req *requests.FindAllReview) ([]*response.ReviewDetailsResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(review_id int) (*response.ReviewDetailsResponse, *response.ErrorResponse)
	CreateReviewDetail(req *requests.CreateReviewDetailRequest) (*response.ReviewDetailsResponse, *response.ErrorResponse)
	UpdateReviewDetail(req *requests.UpdateReviewDetailRequest) (*response.ReviewDetailsResponse, *response.ErrorResponse)
	TrashedReviewDetail(review_id int) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse)
	RestoreReviewDetail(review_id int) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse)
	DeleteReviewDetailPermanent(review_id int) (bool, *response.ErrorResponse)
	RestoreAllReviewDetail() (bool, *response.ErrorResponse)
	DeleteAllReviewDetailPermanent() (bool, *response.ErrorResponse)
}

type ShippingAddressService interface {
	FindAll(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(shipping_id int) (*response.ShippingAddressResponse, *response.ErrorResponse)
	FindByOrder(order_id int) (*response.ShippingAddressResponse, *response.ErrorResponse)
	TrashShippingAddress(shipping_id int) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse)
	RestoreShippingAddress(shipping_id int) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse)
	DeleteShippingAddressPermanently(categoryID int) (bool, *response.ErrorResponse)
	RestoreAllShippingAddress() (bool, *response.ErrorResponse)
	DeleteAllPermanentShippingAddress() (bool, *response.ErrorResponse)
}

type SliderService interface {
	FindAll(req *requests.FindAllSlider) ([]*response.SliderResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, *response.ErrorResponse)
	CreateSlider(req *requests.CreateSliderRequest) (*response.SliderResponse, *response.ErrorResponse)
	UpdateSlider(req *requests.UpdateSliderRequest) (*response.SliderResponse, *response.ErrorResponse)
	TrashedSlider(slider_id int) (*response.SliderResponseDeleteAt, *response.ErrorResponse)
	RestoreSlider(sliderID int) (*response.SliderResponseDeleteAt, *response.ErrorResponse)
	DeleteSliderPermanent(sliderID int) (bool, *response.ErrorResponse)
	RestoreAllSliders() (bool, *response.ErrorResponse)
	DeleteAllSlidersPermanent() (bool, *response.ErrorResponse)
}
