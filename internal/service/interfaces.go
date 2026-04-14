package service

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	db "ecommerce/pkg/database/schema"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type AuthService interface {
	Register(ctx context.Context, request *requests.CreateUserRequest) (*db.CreateUserRow, error)
	Login(ctx context.Context, request *requests.AuthRequest) (*response.TokenResponse, error)
	RefreshToken(ctx context.Context, token string) (*response.TokenResponse, error)
	GetMe(ctx context.Context, token string) (*db.GetUserByIDRow, error)
}

type UserService interface {
	FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersRow, *int, error)
	FindByID(ctx context.Context, id int) (*db.GetUserByIDRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUserTrashedRow, *int, error)

	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*db.CreateUserRow, error)
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*db.UpdateUserRow, error)
	TrashedUser(ctx context.Context, user_id int) (*db.TrashUserRow, error)
	RestoreUser(ctx context.Context, user_id int) (*db.RestoreUserRow, error)
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, error)

	RestoreAllUser(ctx context.Context) (bool, error)
	DeleteAllUserPermanent(ctx context.Context) (bool, error)
}

type RoleService interface {
	FindAll(ctx context.Context, req *requests.FindAllRole) ([]*db.GetRolesRow, *int, error)
	FindByActiveRole(ctx context.Context, req *requests.FindAllRole) ([]*db.GetActiveRolesRow, *int, error)
	FindByTrashedRole(ctx context.Context, req *requests.FindAllRole) ([]*db.GetTrashedRolesRow, *int, error)
	FindById(ctx context.Context, role_id int) (*db.Role, error)
	FindByUserId(ctx context.Context, id int) ([]*db.Role, error)

	CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*db.Role, error)
	UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*db.Role, error)
	TrashedRole(ctx context.Context, role_id int) (*db.Role, error)
	RestoreRole(ctx context.Context, role_id int) (*db.Role, error)
	DeleteRolePermanent(ctx context.Context, role_id int) (bool, error)

	RestoreAllRole(ctx context.Context) (bool, error)
	DeleteAllRolePermanent(ctx context.Context) (bool, error)
}

type BannerService interface {
	FindAll(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersTrashedRow, *int, error)

	FindById(ctx context.Context, bannerID int) (*db.GetBannerRow, error)

	CreateBanner(ctx context.Context, req *requests.CreateBannerRequest) (*db.CreateBannerRow, error)
	UpdateBanner(ctx context.Context, req *requests.UpdateBannerRequest) (*db.UpdateBannerRow, error)

	TrashedBanner(ctx context.Context, bannerID int) (*db.Banner, error)
	RestoreBanner(ctx context.Context, bannerID int) (*db.Banner, error)
	DeleteBannerPermanent(ctx context.Context, bannerID int) (bool, error)

	RestoreAllBanner(ctx context.Context) (bool, error)
	DeleteAllBannerPermanent(ctx context.Context) (bool, error)
}

type CategoryService interface {
	FindMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*db.GetMonthlyTotalPriceRow, error)
	FindYearlyTotalPrice(ctx context.Context, year int) ([]*db.GetYearlyTotalPriceRow, error)
	FindMonthPrice(ctx context.Context, year int) ([]*db.GetMonthlyCategoryRow, error)
	FindYearPrice(ctx context.Context, year int) ([]*db.GetYearlyCategoryRow, error)

	FindMonthlyTotalPriceById(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*db.GetMonthlyTotalPriceByIdRow, error)
	FindYearlyTotalPriceById(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*db.GetYearlyTotalPriceByIdRow, error)

	FindMonthPriceById(ctx context.Context, req *requests.MonthPriceId) ([]*db.GetMonthlyCategoryByIdRow, error)
	FindYearPriceById(ctx context.Context, req *requests.YearPriceId) ([]*db.GetYearlyCategoryByIdRow, error)

	FindMonthlyTotalPriceByMerchant(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*db.GetMonthlyTotalPriceByMerchantRow, error)
	FindYearlyTotalPriceByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*db.GetYearlyTotalPriceByMerchantRow, error)
	FindMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*db.GetMonthlyCategoryByMerchantRow, error)
	FindYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*db.GetYearlyCategoryByMerchantRow, error)

	FindAll(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesTrashedRow, *int, error)

	FindById(ctx context.Context, categoryID int) (*db.GetCategoryByIDRow, error)
	CreateCategory(ctx context.Context, req *requests.CreateCategoryRequest) (*db.CreateCategoryRow, error)
	UpdateCategory(ctx context.Context, req *requests.UpdateCategoryRequest) (*db.UpdateCategoryRow, error)

	TrashedCategory(ctx context.Context, categoryID int) (*db.Category, error)
	RestoreCategory(ctx context.Context, categoryID int) (*db.Category, error)
	DeleteCategoryPermanent(ctx context.Context, categoryID int) (bool, error)

	RestoreAllCategories(ctx context.Context) (bool, error)
	DeleteAllCategoriesPermanent(ctx context.Context) (bool, error)
}

type CartService interface {
	FindAll(ctx context.Context, req *requests.FindAllCarts) ([]*db.GetCartsRow, *int, error)
	CreateCart(ctx context.Context, req *requests.CreateCartRequest) (*db.Cart, error)
	DeletePermanent(ctx context.Context, cartID int) (bool, error)
	DeleteAllPermanently(ctx context.Context, req *requests.DeleteCartRequest) (bool, error)
}

type MerchantService interface {
	FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsRow, *int, error)

	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsActiveRow, *int, error)

	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsTrashedRow, *int, error)

	FindById(ctx context.Context, user_id int) (*db.GetMerchantByIDRow, error)

	CreateMerchant(
		ctx context.Context,
		request *requests.CreateMerchantRequest,
	) (*db.CreateMerchantRow, error)

	UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*db.UpdateMerchantRow, error)

	TrashedMerchant(
		ctx context.Context,
		merchant_id int,
	) (*db.Merchant, error)

	RestoreMerchant(
		ctx context.Context,
		merchant_id int,
	) (*db.Merchant, error)

	DeleteMerchantPermanent(
		ctx context.Context,
		merchant_id int,
	) (bool, error)

	RestoreAllMerchant(ctx context.Context) (bool, error)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, error)
}

type MerchantPoliciesService interface {
	FindAllMerchantPolicy(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesTrashedRow, *int, error)
	FindById(ctx context.Context, user_id int) (*db.GetMerchantPolicyRow, error)

	CreateMerchantPolicy(
		ctx context.Context,
		request *requests.CreateMerchantPolicyRequest,
	) (*db.CreateMerchantPolicyRow, error)

	UpdateMerchantPolicy(
		ctx context.Context,
		request *requests.UpdateMerchantPolicyRequest,
	) (*db.UpdateMerchantPolicyRow, error)

	TrashedMerchantPolicy(
		ctx context.Context,
		merchant_id int,
	) (*db.MerchantPolicy, error)

	RestoreMerchantPolicy(
		ctx context.Context,
		merchant_id int,
	) (*db.MerchantPolicy, error)

	DeleteMerchantPolicyPermanent(
		ctx context.Context,
		merchant_id int,
	) (bool, error)

	RestoreAllMerchantPolicy(ctx context.Context) (bool, error)
	DeleteAllMerchantPolicyPermanent(ctx context.Context) (bool, error)
}

type MerchantAwardService interface {
	FindAllMerchants(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantCertificationsAndAwardsRow, *int, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantCertificationsAndAwardsActiveRow, *int, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantCertificationsAndAwardsTrashedRow, *int, error)

	FindById(
		ctx context.Context,
		user_id int,
	) (*db.GetMerchantCertificationOrAwardRow, error)

	CreateMerchantAward(
		ctx context.Context,
		request *requests.CreateMerchantCertificationOrAwardRequest,
	) (*db.CreateMerchantCertificationOrAwardRow, error)

	UpdateMerchantAward(ctx context.Context, request *requests.UpdateMerchantCertificationOrAwardRequest) (*db.UpdateMerchantCertificationOrAwardRow, error)

	TrashedMerchantAward(
		ctx context.Context,
		merchant_id int,
	) (*db.MerchantCertificationsAndAward, error)

	RestoreMerchantAward(
		ctx context.Context,
		merchant_id int,
	) (*db.MerchantCertificationsAndAward, error)

	DeleteMerchantPermanent(
		ctx context.Context,
		merchant_id int,
	) (bool, error)

	RestoreAllMerchantAward(ctx context.Context) (bool, error)
	DeleteAllMerchantAwardPermanent(ctx context.Context) (bool, error)
}

type MerchantBusinessService interface {
	FindAllMerchants(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantsBusinessInformationRow, *int, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantsBusinessInformationActiveRow, *int, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantsBusinessInformationTrashedRow, *int, error)

	FindById(
		ctx context.Context,
		user_id int,
	) (*db.GetMerchantBusinessInformationRow, error)

	CreateMerchantBusiness(
		ctx context.Context,
		request *requests.CreateMerchantBusinessInformationRequest,
	) (*db.CreateMerchantBusinessInformationRow, error)

	UpdateMerchantBusiness(
		ctx context.Context,
		request *requests.UpdateMerchantBusinessInformationRequest,
	) (*db.UpdateMerchantBusinessInformationRow, error)

	TrashedMerchantBusiness(
		ctx context.Context,
		merchant_id int,
	) (*db.MerchantBusinessInformation, error)

	RestoreMerchantBusiness(
		ctx context.Context,
		merchant_id int,
	) (*db.MerchantBusinessInformation, error)

	DeleteMerchantBusinessPermanent(
		ctx context.Context,
		merchant_id int,
	) (bool, error)

	RestoreAllMerchantBusiness(ctx context.Context) (bool, error)
	DeleteAllMerchantBusinessPermanent(ctx context.Context) (bool, error)
}

type MerchantDetailService interface {
	FindAllMerchants(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantDetailsRow, *int, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantDetailsActiveRow, *int, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantDetailsTrashedRow, *int, error)

	FindById(
		ctx context.Context,
		user_id int,
	) (*db.GetMerchantDetailRow, error)

	CreateMerchantDetail(
		ctx context.Context,
		request *requests.CreateMerchantDetailRequest,
	) (*db.CreateMerchantDetailRow, error)

	UpdateMerchantDetail(
		ctx context.Context,
		request *requests.UpdateMerchantDetailRequest,
	) (*db.UpdateMerchantDetailRow, error)

	TrashedMerchantDetail(
		ctx context.Context,
		merchant_id int,
	) (*db.MerchantDetail, error)

	RestoreMerchantDetail(
		ctx context.Context,
		merchant_id int,
	) (*db.MerchantDetail, error)

	DeleteMerchantDetailPermanent(
		ctx context.Context,
		merchant_id int,
	) (bool, error)

	RestoreAllMerchantDetail(ctx context.Context) (bool, error)
	DeleteAllMerchantDetailPermanent(ctx context.Context) (bool, error)
}

type MerchantSocialLinkService interface {
	CreateSocialLink(
		ctx context.Context,
		req *requests.CreateMerchantSocialRequest,
	) (bool, error)

	UpdateSocialLink(
		ctx context.Context,
		req *requests.UpdateMerchantSocialRequest,
	) (bool, error)

	TrashSocialLink(
		ctx context.Context,
		socialID int,
	) (bool, error)

	RestoreSocialLink(
		ctx context.Context,
		socialID int,
	) (bool, error)

	DeletePermanentSocialLink(
		ctx context.Context,
		socialID int,
	) (bool, error)

	RestoreAllSocialLink(ctx context.Context) (bool, error)
	DeleteAllPermanentSocialLink(ctx context.Context) (bool, error)
}

type OrderService interface {
	FindMonthlyTotalRevenue(
		ctx context.Context,
		req *requests.MonthTotalRevenue,
	) ([]*db.GetMonthlyTotalRevenueRow, error)

	FindYearlyTotalRevenue(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyTotalRevenueRow, error)
	FindMonthlyOrder(
		ctx context.Context,
		year int,
	) ([]*db.GetMonthlyOrderRow, error)

	FindYearlyOrder(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyOrderRow, error)

	FindMonthlyTotalRevenueById(
		ctx context.Context,
		req *requests.MonthTotalRevenueOrder,
	) ([]*db.GetMonthlyTotalRevenueByIdRow, error)

	FindYearlyTotalRevenueById(
		ctx context.Context,
		req *requests.YearTotalRevenueOrder,
	) ([]*db.GetYearlyTotalRevenueByIdRow, error)

	FindMonthlyTotalRevenueByMerchant(
		ctx context.Context,
		req *requests.MonthTotalRevenueMerchant,
	) ([]*db.GetMonthlyTotalRevenueByMerchantRow, error)

	FindYearlyTotalRevenueByMerchant(
		ctx context.Context,
		req *requests.YearTotalRevenueMerchant,
	) ([]*db.GetYearlyTotalRevenueByMerchantRow, error)

	FindMonthlyOrderByMerchant(
		ctx context.Context,
		req *requests.MonthOrderMerchant,
	) ([]*db.GetMonthlyOrderByMerchantRow, error)

	FindYearlyOrderByMerchant(
		ctx context.Context,
		req *requests.YearOrderMerchant,
	) ([]*db.GetYearlyOrderByMerchantRow, error)

	FindAllOrders(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*db.GetOrdersRow, *int, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*db.GetOrdersActiveRow, *int, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*db.GetOrdersTrashedRow, *int, error)

	FindById(
		ctx context.Context,
		order_id int,
	) (*db.GetOrderByIDRow, error)

	CreateOrder(
		ctx context.Context,
		request *requests.CreateOrderRequest,
	) (*db.CreateOrderRow, error)

	UpdateOrder(
		ctx context.Context,
		request *requests.UpdateOrderRequest,
	) (*db.UpdateOrderRow, error)

	TrashedOrder(
		ctx context.Context,
		order_id int,
	) (*db.Order, error)

	RestoreOrder(
		ctx context.Context,
		order_id int,
	) (*db.Order, error)

	DeleteOrderPermanent(
		ctx context.Context,
		order_id int,
	) (bool, error)

	RestoreAllOrder(ctx context.Context) (bool, error)
	DeleteAllOrderPermanent(ctx context.Context) (bool, error)
}

type OrderItemService interface {
	FindAllOrderItems(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*db.GetOrderItemsRow, *int, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*db.GetOrderItemsActiveRow, *int, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*db.GetOrderItemsTrashedRow, *int, error)

	FindOrderItemByOrder(
		ctx context.Context,
		order_id int,
	) ([]*db.GetOrderItemsByOrderRow, error)
}

type ProductService interface {
	FindAllProducts(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsTrashedRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*db.GetProductsByMerchantRow, *int, error)
	FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*db.GetProductsByCategoryNameRow, *int, error)
	FindById(ctx context.Context, product_id int) (*db.GetProductByIDRow, error)
	CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*db.CreateProductRow, error)
	UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*db.UpdateProductRow, error)
	UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*db.UpdateProductCountStockRow, error)
	TrashedProduct(ctx context.Context, product_id int) (*db.Product, error)
	RestoreProduct(ctx context.Context, product_id int) (*db.Product, error)

	DeleteProductPermanent(
		ctx context.Context,
		product_id int,
	) (bool, error)

	RestoreAllProducts(ctx context.Context) (bool, error)
	DeleteAllProductPermanent(ctx context.Context) (bool, error)
}

type TransactionService interface {
	FindMonthlyAmountSuccess(
		ctx context.Context,
		req *requests.MonthAmountTransaction,
	) ([]*db.GetMonthlyAmountTransactionSuccessRow, error)

	FindYearlyAmountSuccess(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyAmountTransactionSuccessRow, error)

	FindMonthlyAmountFailed(
		ctx context.Context,
		req *requests.MonthAmountTransaction,
	) ([]*db.GetMonthlyAmountTransactionFailedRow, error)

	FindYearlyAmountFailed(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyAmountTransactionFailedRow, error)

	FindMonthlyAmountSuccessByMerchant(
		ctx context.Context,
		req *requests.MonthAmountTransactionMerchant,
	) ([]*db.GetMonthlyAmountTransactionSuccessByMerchantRow, error)

	FindYearlyAmountSuccessByMerchant(
		ctx context.Context,
		req *requests.YearAmountTransactionMerchant,
	) ([]*db.GetYearlyAmountTransactionSuccessByMerchantRow, error)

	FindMonthlyAmountFailedByMerchant(
		ctx context.Context,
		req *requests.MonthAmountTransactionMerchant,
	) ([]*db.GetMonthlyAmountTransactionFailedByMerchantRow, error)

	FindYearlyAmountFailedByMerchant(
		ctx context.Context,
		req *requests.YearAmountTransactionMerchant,
	) ([]*db.GetYearlyAmountTransactionFailedByMerchantRow, error)

	FindMonthlyTransactionMethodSuccess(
		ctx context.Context,
		req *requests.MonthMethodTransaction,
	) ([]*db.GetMonthlyTransactionMethodsSuccessRow, error)

	FindYearlyTransactionMethodSuccess(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyTransactionMethodsSuccessRow, error)

	FindMonthlyTransactionMethodFailed(
		ctx context.Context,
		req *requests.MonthMethodTransaction,
	) ([]*db.GetMonthlyTransactionMethodsFailedRow, error)

	FindYearlyTransactionMethodFailed(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyTransactionMethodsFailedRow, error)

	FindMonthlyTransactionMethodByMerchantSuccess(
		ctx context.Context,
		req *requests.MonthMethodTransactionMerchant,
	) ([]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow, error)

	FindYearlyTransactionMethodByMerchantSuccess(
		ctx context.Context,
		req *requests.YearMethodTransactionMerchant,
	) ([]*db.GetYearlyTransactionMethodsByMerchantSuccessRow, error)

	FindMonthlyTransactionMethodByMerchantFailed(
		ctx context.Context,
		req *requests.MonthMethodTransactionMerchant,
	) ([]*db.GetMonthlyTransactionMethodsByMerchantFailedRow, error)

	FindYearlyTransactionMethodByMerchantFailed(
		ctx context.Context,
		req *requests.YearMethodTransactionMerchant,
	) ([]*db.GetYearlyTransactionMethodsByMerchantFailedRow, error)

	FindAllTransactions(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*db.GetTransactionsRow, *int, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*db.GetTransactionsActiveRow, *int, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*db.GetTransactionsTrashedRow, *int, error)

	FindByMerchant(
		ctx context.Context,
		req *requests.FindAllTransactionByMerchant,
	) ([]*db.GetTransactionByMerchantRow, *int, error)

	FindById(
		ctx context.Context,
		transaction_id int,
	) (*db.GetTransactionByIDRow, error)

	FindByOrderId(
		ctx context.Context,
		order_id int,
	) (*db.GetTransactionByOrderIDRow, error)

	CreateTransaction(
		ctx context.Context,
		request *requests.CreateTransactionRequest,
	) (*db.CreateTransactionRow, error)

	UpdateTransaction(
		ctx context.Context,
		request *requests.UpdateTransactionRequest,
	) (*db.UpdateTransactionRow, error)

	TrashedTransaction(
		ctx context.Context,
		transaction_id int,
	) (*db.Transaction, error)

	RestoreTransaction(
		ctx context.Context,
		transaction_id int,
	) (*db.Transaction, error)

	DeleteTransactionPermanently(
		ctx context.Context,
		transaction_id int,
	) (bool, error)

	RestoreAllTransactions(ctx context.Context) (bool, error)
	DeleteAllTransactionPermanent(ctx context.Context) (bool, error)
}

type ShippingAddressService interface {
	FindAllShippingAddress(
		ctx context.Context,
		req *requests.FindAllShippingAddress,
	) ([]*db.GetShippingAddressRow, *int, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllShippingAddress,
	) ([]*db.GetShippingAddressActiveRow, *int, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllShippingAddress,
	) ([]*db.GetShippingAddressTrashedRow, *int, error)

	FindByOrder(
		ctx context.Context,
		shipping_id int,
	) (*db.GetShippingAddressByOrderIDRow, error)

	FindById(
		ctx context.Context,
		shipping_id int,
	) (*db.GetShippingByIDRow, error)

	TrashShippingAddress(
		ctx context.Context,
		shipping_id int,
	) (*db.ShippingAddress, error)

	RestoreShippingAddress(
		ctx context.Context,
		shipping_id int,
	) (*db.ShippingAddress, error)

	DeleteShippingAddressPermanently(
		ctx context.Context,
		shipping_id int,
	) (bool, error)

	RestoreAllShippingAddress(ctx context.Context) (bool, error)
	DeleteAllPermanentShippingAddress(ctx context.Context) (bool, error)
}

type SliderService interface {
	FindAllSlider(
		ctx context.Context,
		req *requests.FindAllSlider,
	) ([]*db.GetSlidersRow, *int, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllSlider,
	) ([]*db.GetSlidersActiveRow, *int, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllSlider,
	) ([]*db.GetSlidersTrashedRow, *int, error)

	FindById(
		ctx context.Context,
		slider_id int,
	) (*db.GetSliderByIDRow, error)

	CreateSlider(
		ctx context.Context,
		request *requests.CreateSliderRequest,
	) (*db.CreateSliderRow, error)

	UpdateSlider(
		ctx context.Context,
		request *requests.UpdateSliderRequest,
	) (*db.UpdateSliderRow, error)

	TrashSlider(
		ctx context.Context,
		slider_id int,
	) (*db.Slider, error)

	RestoreSlider(
		ctx context.Context,
		slider_id int,
	) (*db.Slider, error)

	DeleteSliderPermanently(
		ctx context.Context,
		slider_id int,
	) (bool, error)

	RestoreAllSliders(ctx context.Context) (bool, error)
	DeleteAllPermanentSlider(ctx context.Context) (bool, error)
}

type ReviewService interface {
	FindAllReview(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsRow, *int, error)
	FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*db.GetReviewByProductIdRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*db.GetReviewByMerchantIdRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsTrashedRow, *int, error)
	FindById(ctx context.Context, id int) (*db.GetReviewByIDRow, error)
	CreateReview(ctx context.Context, request *requests.CreateReviewRequest) (*db.CreateReviewRow, error)
	UpdateReview(ctx context.Context, request *requests.UpdateReviewRequest) (*db.UpdateReviewRow, error)
	TrashReview(ctx context.Context, shipping_id int) (*db.Review, error)
	RestoreReview(ctx context.Context, category_id int) (*db.Review, error)

	DeleteReviewPermanently(
		ctx context.Context,
		id int,
	) (bool, error)

	RestoreAllReview(ctx context.Context) (bool, error)
	DeleteAllPermanentReview(ctx context.Context) (bool, error)
}

type ReviewDetailService interface {
	FindAllReviews(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsTrashedRow, *int, error)
	FindById(ctx context.Context, user_id int) (*db.GetReviewDetailRow, error)

	CreateReviewDetail(ctx context.Context, request *requests.CreateReviewDetailRequest) (*db.CreateReviewDetailRow, error)
	UpdateReviewDetail(ctx context.Context, request *requests.UpdateReviewDetailRequest) (*db.UpdateReviewDetailRow, error)
	TrashedReviewDetail(ctx context.Context, ReviewDetail_id int) (*db.ReviewDetail, error)
	RestoreReviewDetail(ctx context.Context, ReviewDetail_id int) (*db.ReviewDetail, error)

	DeleteReviewDetailPermanent(
		ctx context.Context,
		review_detail_id int,
	) (bool, error)

	RestoreAllReviewDetail(ctx context.Context) (bool, error)
	DeleteAllReviewDetailPermanent(ctx context.Context) (bool, error)
}
