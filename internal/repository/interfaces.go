package repository

import (
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
)

// go generate mockgen -source =interfaces.go -destination=mocks/mock.go
type UserRepository interface {
	FindAllUsers(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)
	FindById(user_id int) (*record.UserRecord, error)
	FindByEmail(email string) (*record.UserRecord, error)
	FindByActive(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)
	FindByTrashed(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)
	CreateUser(request *requests.CreateUserRequest) (*record.UserRecord, error)
	UpdateUser(request *requests.UpdateUserRequest) (*record.UserRecord, error)
	TrashedUser(user_id int) (*record.UserRecord, error)
	RestoreUser(user_id int) (*record.UserRecord, error)
	DeleteUserPermanent(user_id int) (bool, error)
	RestoreAllUser() (bool, error)
	DeleteAllUserPermanent() (bool, error)
}

type RoleRepository interface {
	FindAllRoles(req *requests.FindAllRole) ([]*record.RoleRecord, *int, error)
	FindByActiveRole(req *requests.FindAllRole) ([]*record.RoleRecord, *int, error)
	FindByTrashedRole(req *requests.FindAllRole) ([]*record.RoleRecord, *int, error)
	FindById(role_id int) (*record.RoleRecord, error)
	FindByName(name string) (*record.RoleRecord, error)
	FindByUserId(user_id int) ([]*record.RoleRecord, error)
	CreateRole(request *requests.CreateRoleRequest) (*record.RoleRecord, error)
	UpdateRole(request *requests.UpdateRoleRequest) (*record.RoleRecord, error)
	TrashedRole(role_id int) (*record.RoleRecord, error)

	RestoreRole(role_id int) (*record.RoleRecord, error)
	DeleteRolePermanent(role_id int) (bool, error)
	RestoreAllRole() (bool, error)
	DeleteAllRolePermanent() (bool, error)
}

type RefreshTokenRepository interface {
	FindByToken(token string) (*record.RefreshTokenRecord, error)
	FindByUserId(user_id int) (*record.RefreshTokenRecord, error)
	CreateRefreshToken(req *requests.CreateRefreshToken) (*record.RefreshTokenRecord, error)
	UpdateRefreshToken(req *requests.UpdateRefreshToken) (*record.RefreshTokenRecord, error)
	DeleteRefreshToken(token string) error
	DeleteRefreshTokenByUserId(user_id int) error
}

type UserRoleRepository interface {
	AssignRoleToUser(req *requests.CreateUserRoleRequest) (*record.UserRoleRecord, error)
	RemoveRoleFromUser(req *requests.RemoveUserRoleRequest) error
}

type BannerRepository interface {
	FindAllBanners(req *requests.FindAllBanner) ([]*record.BannerRecord, *int, error)
	FindByActive(req *requests.FindAllBanner) ([]*record.BannerRecord, *int, error)
	FindByTrashed(req *requests.FindAllBanner) ([]*record.BannerRecord, *int, error)
	FindById(user_id int) (*record.BannerRecord, error)
	CreateBanner(request *requests.CreateBannerRequest) (*record.BannerRecord, error)
	UpdateBanner(request *requests.UpdateBannerRequest) (*record.BannerRecord, error)
	TrashedBanner(Banner_id int) (*record.BannerRecord, error)
	RestoreBanner(Banner_id int) (*record.BannerRecord, error)
	DeleteBannerPermanent(Banner_id int) (bool, error)
	RestoreAllBanner() (bool, error)
	DeleteAllBannerPermanent() (bool, error)
}

type CategoryRepository interface {
	GetMonthlyTotalPrice(req *requests.MonthTotalPrice) ([]*record.CategoriesMonthlyTotalPriceRecord, error)
	GetYearlyTotalPrices(year int) ([]*record.CategoriesYearlyTotalPriceRecord, error)
	GetMonthlyTotalPriceById(req *requests.MonthTotalPriceCategory) ([]*record.CategoriesMonthlyTotalPriceRecord, error)
	GetYearlyTotalPricesById(req *requests.YearTotalPriceCategory) ([]*record.CategoriesYearlyTotalPriceRecord, error)
	GetMonthlyTotalPriceByMerchant(req *requests.MonthTotalPriceMerchant) ([]*record.CategoriesMonthlyTotalPriceRecord, error)
	GetYearlyTotalPricesByMerchant(req *requests.YearTotalPriceMerchant) ([]*record.CategoriesYearlyTotalPriceRecord, error)

	GetMonthPrice(year int) ([]*record.CategoriesMonthPriceRecord, error)
	GetYearPrice(year int) ([]*record.CategoriesYearPriceRecord, error)
	GetMonthPriceByMerchant(req *requests.MonthPriceMerchant) ([]*record.CategoriesMonthPriceRecord, error)
	GetYearPriceByMerchant(req *requests.YearPriceMerchant) ([]*record.CategoriesYearPriceRecord, error)
	GetMonthPriceById(req *requests.MonthPriceId) ([]*record.CategoriesMonthPriceRecord, error)
	GetYearPriceById(req *requests.YearPriceId) ([]*record.CategoriesYearPriceRecord, error)

	FindAllCategory(req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error)
	FindByActive(req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error)
	FindByTrashed(req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error)
	FindById(category_id int) (*record.CategoriesRecord, error)
	FindByIdTrashed(category_id int) (*record.CategoriesRecord, error)
	CreateCategory(request *requests.CreateCategoryRequest) (*record.CategoriesRecord, error)
	UpdateCategory(request *requests.UpdateCategoryRequest) (*record.CategoriesRecord, error)
	TrashedCategory(category_id int) (*record.CategoriesRecord, error)
	RestoreCategory(category_id int) (*record.CategoriesRecord, error)
	DeleteCategoryPermanently(Category_id int) (bool, error)
	RestoreAllCategories() (bool, error)
	DeleteAllPermanentCategories() (bool, error)
}

type MerchantRepository interface {
	FindAllMerchants(req *requests.FindAllMerchant) ([]*record.MerchantRecord, *int, error)
	FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantRecord, *int, error)
	FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantRecord, *int, error)
	FindById(user_id int) (*record.MerchantRecord, error)
	CreateMerchant(request *requests.CreateMerchantRequest) (*record.MerchantRecord, error)
	UpdateMerchant(request *requests.UpdateMerchantRequest) (*record.MerchantRecord, error)
	TrashedMerchant(merchant_id int) (*record.MerchantRecord, error)
	RestoreMerchant(merchant_id int) (*record.MerchantRecord, error)
	DeleteMerchantPermanent(Merchant_id int) (bool, error)
	RestoreAllMerchant() (bool, error)
	DeleteAllMerchantPermanent() (bool, error)
}

type MerchantPoliciesRepository interface {
	FindAllMerchantPolicy(req *requests.FindAllMerchant) ([]*record.MerchantPoliciesRecord, *int, error)
	FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantPoliciesRecord, *int, error)
	FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantPoliciesRecord, *int, error)
	FindById(user_id int) (*record.MerchantPoliciesRecord, error)
	CreateMerchantPolicy(request *requests.CreateMerchantPolicyRequest) (*record.MerchantPoliciesRecord, error)
	UpdateMerchantPolicy(request *requests.UpdateMerchantPolicyRequest) (*record.MerchantPoliciesRecord, error)
	TrashedMerchantPolicy(merchant_id int) (*record.MerchantPoliciesRecord, error)
	RestoreMerchantPolicy(merchant_id int) (*record.MerchantPoliciesRecord, error)
	DeleteMerchantPolicyPermanent(Merchant_id int) (bool, error)
	RestoreAllMerchantPolicy() (bool, error)
	DeleteAllMerchantPolicyPermanent() (bool, error)
}

type MerchantAwardRepository interface {
	FindAllMerchants(req *requests.FindAllMerchant) ([]*record.MerchantAwardRecord, *int, error)
	FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantAwardRecord, *int, error)
	FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantAwardRecord, *int, error)
	FindById(user_id int) (*record.MerchantAwardRecord, error)
	CreateMerchantAward(request *requests.CreateMerchantCertificationOrAwardRequest) (*record.MerchantAwardRecord, error)
	UpdateMerchantAward(request *requests.UpdateMerchantCertificationOrAwardRequest) (*record.MerchantAwardRecord, error)
	TrashedMerchantAward(merchant_id int) (*record.MerchantAwardRecord, error)
	RestoreMerchantAward(merchant_id int) (*record.MerchantAwardRecord, error)
	DeleteMerchantPermanent(Merchant_id int) (bool, error)
	RestoreAllMerchantAward() (bool, error)
	DeleteAllMerchantAwardPermanent() (bool, error)
}

type MerchantBusinessRepository interface {
	FindAllMerchants(req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error)
	FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error)
	FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error)
	FindById(user_id int) (*record.MerchantBusinessRecord, error)
	CreateMerchantBusiness(request *requests.CreateMerchantBusinessInformationRequest) (*record.MerchantBusinessRecord, error)
	UpdateMerchantBusiness(request *requests.UpdateMerchantBusinessInformationRequest) (*record.MerchantBusinessRecord, error)
	TrashedMerchantBusiness(merchant_id int) (*record.MerchantBusinessRecord, error)
	RestoreMerchantBusiness(merchant_id int) (*record.MerchantBusinessRecord, error)
	DeleteMerchantBusinessPermanent(Merchant_id int) (bool, error)
	RestoreAllMerchantBusiness() (bool, error)
	DeleteAllMerchantBusinessPermanent() (bool, error)
}

type MerchantDetailRepository interface {
	FindAllMerchants(req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error)
	FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error)
	FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error)
	FindById(user_id int) (*record.MerchantDetailRecord, error)
	FindByIdTrashed(user_id int) (*record.MerchantDetailRecord, error)
	CreateMerchantDetail(request *requests.CreateMerchantDetailRequest) (*record.MerchantDetailRecord, error)
	UpdateMerchantDetail(request *requests.UpdateMerchantDetailRequest) (*record.MerchantDetailRecord, error)
	TrashedMerchantDetail(merchant_id int) (*record.MerchantDetailRecord, error)
	RestoreMerchantDetail(merchant_id int) (*record.MerchantDetailRecord, error)
	DeleteMerchantDetailPermanent(Merchant_id int) (bool, error)
	RestoreAllMerchantDetail() (bool, error)
	DeleteAllMerchantDetailPermanent() (bool, error)
}

type MerchantSocialLinkRepository interface {
	CreateSocialLink(req *requests.CreateMerchantSocialRequest) (bool, error)
	UpdateSocialLink(req *requests.UpdateMerchantSocialRequest) (bool, error)
	TrashSocialLink(socialID int) (bool, error)
	RestoreSocialLink(socialID int) (bool, error)
	DeletePermanentSocialLink(socialID int) (bool, error)
	RestoreAllSocialLink() (bool, error)
	DeleteAllPermanentSocialLink() (bool, error)
}

type OrderRepository interface {
	GetMonthlyTotalRevenue(req *requests.MonthTotalRevenue) ([]*record.OrderMonthlyTotalRevenueRecord, error)
	GetYearlyTotalRevenue(year int) ([]*record.OrderYearlyTotalRevenueRecord, error)
	GetMonthlyTotalRevenueById(req *requests.MonthTotalRevenueOrder) ([]*record.OrderMonthlyTotalRevenueRecord, error)
	GetYearlyTotalRevenueById(req *requests.YearTotalRevenueOrder) ([]*record.OrderYearlyTotalRevenueRecord, error)
	GetMonthlyTotalRevenueByMerchant(req *requests.MonthTotalRevenueMerchant) ([]*record.OrderMonthlyTotalRevenueRecord, error)
	GetYearlyTotalRevenueByMerchant(req *requests.YearTotalRevenueMerchant) ([]*record.OrderYearlyTotalRevenueRecord, error)

	GetMonthlyOrder(year int) ([]*record.OrderMonthlyRecord, error)
	GetYearlyOrder(year int) ([]*record.OrderYearlyRecord, error)
	GetMonthlyOrderByMerchant(req *requests.MonthOrderMerchant) ([]*record.OrderMonthlyRecord, error)
	GetYearlyOrderByMerchant(req *requests.YearOrderMerchant) ([]*record.OrderYearlyRecord, error)

	FindAllOrders(req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error)
	FindByActive(req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error)
	FindByTrashed(req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error)
	FindByMerchant(req *requests.FindAllOrderByMerchant) ([]*record.OrderRecord, *int, error)
	FindById(order_id int) (*record.OrderRecord, error)
	CreateOrder(request *requests.CreateOrderRecordRequest) (*record.OrderRecord, error)
	UpdateOrder(request *requests.UpdateOrderRecordRequest) (*record.OrderRecord, error)
	TrashedOrder(order_id int) (*record.OrderRecord, error)
	RestoreOrder(order_id int) (*record.OrderRecord, error)
	DeleteOrderPermanent(order_id int) (bool, error)
	RestoreAllOrder() (bool, error)
	DeleteAllOrderPermanent() (bool, error)
}

type OrderItemRepository interface {
	FindAllOrderItems(req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error)
	FindByActive(req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error)
	FindByTrashed(req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error)
	FindOrderItemByOrder(order_id int) ([]*record.OrderItemRecord, error)
	CalculateTotalPrice(order_id int) (*int32, error)
	CreateOrderItem(req *requests.CreateOrderItemRecordRequest) (*record.OrderItemRecord, error)
	UpdateOrderItem(req *requests.UpdateOrderItemRecordRequest) (*record.OrderItemRecord, error)
	TrashedOrderItem(order_id int) (*record.OrderItemRecord, error)
	RestoreOrderItem(order_id int) (*record.OrderItemRecord, error)
	DeleteOrderItemPermanent(order_id int) (bool, error)
	RestoreAllOrderItem() (bool, error)
	DeleteAllOrderPermanent() (bool, error)
}

type ProductRepository interface {
	FindAllProducts(req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error)
	FindByActive(req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error)
	FindByTrashed(req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error)
	FindByMerchant(req *requests.FindAllProductByMerchant) ([]*record.ProductRecord, *int, error)
	FindByCategory(req *requests.FindAllProductByCategory) ([]*record.ProductRecord, *int, error)

	FindByIdTrashed(product_id int) (*record.ProductRecord, error)
	FindById(product_id int) (*record.ProductRecord, error)
	CreateProduct(request *requests.CreateProductRequest) (*record.ProductRecord, error)
	UpdateProduct(request *requests.UpdateProductRequest) (*record.ProductRecord, error)
	UpdateProductCountStock(product_id int, stock int) (*record.ProductRecord, error)
	TrashedProduct(product_id int) (*record.ProductRecord, error)
	RestoreProduct(product_id int) (*record.ProductRecord, error)
	DeleteProductPermanent(product_id int) (bool, error)
	RestoreAllProducts() (bool, error)
	DeleteAllProductPermanent() (bool, error)
}

type TransactionRepository interface {
	GetMonthlyAmountSuccess(req *requests.MonthAmountTransaction) ([]*record.TransactionMonthlyAmountSuccessRecord, error)
	GetYearlyAmountSuccess(year int) ([]*record.TransactionYearlyAmountSuccessRecord, error)
	GetMonthlyAmountFailed(req *requests.MonthAmountTransaction) ([]*record.TransactionMonthlyAmountFailedRecord, error)
	GetYearlyAmountFailed(year int) ([]*record.TransactionYearlyAmountFailedRecord, error)

	GetMonthlyAmountSuccessByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*record.TransactionMonthlyAmountSuccessRecord, error)
	GetYearlyAmountSuccessByMerchant(req *requests.YearAmountTransactionMerchant) ([]*record.TransactionYearlyAmountSuccessRecord, error)
	GetMonthlyAmountFailedByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*record.TransactionMonthlyAmountFailedRecord, error)
	GetYearlyAmountFailedByMerchant(req *requests.YearAmountTransactionMerchant) ([]*record.TransactionYearlyAmountFailedRecord, error)

	GetMonthlyTransactionMethodSuccess(req *requests.MonthMethodTransaction) ([]*record.TransactionMonthlyMethodRecord, error)
	GetYearlyTransactionMethodSuccess(year int) ([]*record.TransactionYearlyMethodRecord, error)
	GetMonthlyTransactionMethodByMerchantSuccess(req *requests.MonthMethodTransactionMerchant) ([]*record.TransactionMonthlyMethodRecord, error)
	GetYearlyTransactionMethodByMerchantSuccess(req *requests.YearMethodTransactionMerchant) ([]*record.TransactionYearlyMethodRecord, error)

	GetMonthlyTransactionMethodFailed(req *requests.MonthMethodTransaction) ([]*record.TransactionMonthlyMethodRecord, error)
	GetYearlyTransactionMethodFailed(year int) ([]*record.TransactionYearlyMethodRecord, error)
	GetMonthlyTransactionMethodByMerchantFailed(req *requests.MonthMethodTransactionMerchant) ([]*record.TransactionMonthlyMethodRecord, error)
	GetYearlyTransactionMethodByMerchantFailed(req *requests.YearMethodTransactionMerchant) ([]*record.TransactionYearlyMethodRecord, error)

	FindAllTransactions(req *requests.FindAllTransaction) ([]*record.TransactionRecord, *int, error)
	FindByActive(req *requests.FindAllTransaction) ([]*record.TransactionRecord, *int, error)
	FindByTrashed(req *requests.FindAllTransaction) ([]*record.TransactionRecord, *int, error)
	FindByMerchant(req *requests.FindAllTransactionByMerchant) ([]*record.TransactionRecord, *int, error)
	FindById(transaction_id int) (*record.TransactionRecord, error)
	FindByOrderId(order_id int) (*record.TransactionRecord, error)
	CreateTransaction(request *requests.CreateTransactionRequest) (*record.TransactionRecord, error)
	UpdateTransaction(request *requests.UpdateTransactionRequest) (*record.TransactionRecord, error)
	TrashTransaction(transaction_id int) (*record.TransactionRecord, error)
	RestoreTransaction(transaction_id int) (*record.TransactionRecord, error)
	DeleteTransactionPermanently(transaction_id int) (bool, error)
	RestoreAllTransactions() (bool, error)
	DeleteAllTransactionPermanent() (bool, error)
}

type CartRepository interface {
	FindCarts(req *requests.FindAllCarts) ([]*record.CartRecord, *int, error)
	CreateCart(req *requests.CartCreateRecord) (*record.CartRecord, error)
	DeletePermanent(cart_id int) (bool, error)
	DeleteAllPermanently(req *requests.DeleteCartRequest) (bool, error)
}

type ShippingAddressRepository interface {
	FindAllShippingAddress(req *requests.FindAllShippingAddress) ([]*record.ShippingAddressRecord, *int, error)
	FindByActive(req *requests.FindAllShippingAddress) ([]*record.ShippingAddressRecord, *int, error)
	FindByTrashed(req *requests.FindAllShippingAddress) ([]*record.ShippingAddressRecord, *int, error)
	FindByOrder(shipping_id int) (*record.ShippingAddressRecord, error)
	FindById(shipping_id int) (*record.ShippingAddressRecord, error)
	CreateShippingAddress(request *requests.CreateShippingAddressRequest) (*record.ShippingAddressRecord, error)
	UpdateShippingAddress(request *requests.UpdateShippingAddressRequest) (*record.ShippingAddressRecord, error)
	TrashShippingAddress(category_id int) (*record.ShippingAddressRecord, error)
	RestoreShippingAddress(category_id int) (*record.ShippingAddressRecord, error)
	DeleteShippingAddressPermanently(category_id int) (bool, error)
	RestoreAllShippingAddress() (bool, error)
	DeleteAllPermanentShippingAddress() (bool, error)
}

type SliderRepository interface {
	FindAllSlider(req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error)
	FindByActive(req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error)
	FindByTrashed(req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error)
	CreateSlider(request *requests.CreateSliderRequest) (*record.SliderRecord, error)
	UpdateSlider(request *requests.UpdateSliderRequest) (*record.SliderRecord, error)
	TrashSlider(slider_id int) (*record.SliderRecord, error)
	RestoreSlider(slider_id int) (*record.SliderRecord, error)
	DeleteSliderPermanently(slider_id int) (bool, error)
	RestoreAllSlider() (bool, error)
	DeleteAllPermanentSlider() (bool, error)
}

type ReviewRepository interface {
	FindAllReview(req *requests.FindAllReview) ([]*record.ReviewRecord, *int, error)
	FindByProduct(req *requests.FindAllReviewByProduct) ([]*record.ReviewsDetailRecord, *int, error)
	FindByMerchant(req *requests.FindAllReviewByMerchant) ([]*record.ReviewsDetailRecord, *int, error)
	FindByActive(req *requests.FindAllReview) ([]*record.ReviewRecord, *int, error)
	FindByTrashed(req *requests.FindAllReview) ([]*record.ReviewRecord, *int, error)
	FindById(id int) (*record.ReviewRecord, error)
	CreateReview(request *requests.CreateReviewRequest) (*record.ReviewRecord, error)
	UpdateReview(request *requests.UpdateReviewRequest) (*record.ReviewRecord, error)
	TrashReview(shipping_id int) (*record.ReviewRecord, error)
	RestoreReview(category_id int) (*record.ReviewRecord, error)
	DeleteReviewPermanently(category_id int) (bool, error)
	RestoreAllReview() (bool, error)
	DeleteAllPermanentReview() (bool, error)
}

type ReviewDetailRepository interface {
	FindAllReviews(req *requests.FindAllReview) ([]*record.ReviewDetailRecord, *int, error)
	FindByActive(req *requests.FindAllReview) ([]*record.ReviewDetailRecord, *int, error)
	FindByTrashed(req *requests.FindAllReview) ([]*record.ReviewDetailRecord, *int, error)
	FindById(user_id int) (*record.ReviewDetailRecord, error)
	FindByIdTrashed(user_id int) (*record.ReviewDetailRecord, error)
	CreateReviewDetail(request *requests.CreateReviewDetailRequest) (*record.ReviewDetailRecord, error)
	UpdateReviewDetail(request *requests.UpdateReviewDetailRequest) (*record.ReviewDetailRecord, error)
	TrashedReviewDetail(ReviewDetail_id int) (*record.ReviewDetailRecord, error)
	RestoreReviewDetail(ReviewDetail_id int) (*record.ReviewDetailRecord, error)
	DeleteReviewDetailPermanent(ReviewDetail_id int) (bool, error)
	RestoreAllReviewDetail() (bool, error)
	DeleteAllReviewDetailPermanent() (bool, error)
}
