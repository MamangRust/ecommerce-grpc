package repository

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
)

// go generate mockgen -source =interfaces.go -destination=mocks/mock.go
type UserRepository interface {
	FindAllUsers(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersActiveRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUserTrashedRow, error)
	FindById(ctx context.Context, user_id int) (*db.GetUserByIDRow, error)
	FindByIdWithPassword(ctx context.Context, user_id int) (*db.GetUserByIDWithPasswordRow, error)
	FindByEmail(ctx context.Context, email string) (*db.GetUserByEmailRow, error)
	FindByEmailWithPassword(ctx context.Context, email string) (*db.GetUserByEmailWithPasswordRow, error)
	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*db.CreateUserRow, error)
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*db.UpdateUserRow, error)
	TrashedUser(ctx context.Context, user_id int) (*db.TrashUserRow, error)
	RestoreUser(ctx context.Context, user_id int) (*db.RestoreUserRow, error)
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, error)
	RestoreAllUser(ctx context.Context) (bool, error)
	DeleteAllUserPermanent(ctx context.Context) (bool, error)
}

type RoleRepository interface {
	FindAllRoles(ctx context.Context, req *requests.FindAllRole) ([]*db.GetRolesRow, error)
	FindByActiveRole(ctx context.Context, req *requests.FindAllRole) ([]*db.GetActiveRolesRow, error)
	FindByTrashedRole(ctx context.Context, req *requests.FindAllRole) ([]*db.GetTrashedRolesRow, error)
	FindById(ctx context.Context, role_id int) (*db.Role, error)
	FindByName(ctx context.Context, name string) (*db.Role, error)
	FindByUserId(ctx context.Context, user_id int) ([]*db.Role, error)
	CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*db.Role, error)
	UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*db.Role, error)
	TrashedRole(ctx context.Context, role_id int) (*db.Role, error)
	RestoreRole(ctx context.Context, role_id int) (*db.Role, error)
	DeleteRolePermanent(ctx context.Context, role_id int) (bool, error)
	RestoreAllRole(ctx context.Context) (bool, error)
	DeleteAllRolePermanent(ctx context.Context) (bool, error)
}

type RefreshTokenRepository interface {
	FindByToken(ctx context.Context, token string) (*db.RefreshToken, error)
	FindByUserId(ctx context.Context, user_id int) (*db.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*db.RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*db.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error
}

type UserRoleRepository interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*db.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

type BannerRepository interface {
	FindAllBanners(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersActiveRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersTrashedRow, error)
	FindById(ctx context.Context, user_id int) (*db.GetBannerRow, error)

	CreateBanner(ctx context.Context, request *requests.CreateBannerRequest) (*db.CreateBannerRow, error)
	UpdateBanner(ctx context.Context, request *requests.UpdateBannerRequest) (*db.UpdateBannerRow, error)

	TrashedBanner(ctx context.Context, Banner_id int) (*db.Banner, error)
	RestoreBanner(ctx context.Context, Banner_id int) (*db.Banner, error)
	DeleteBannerPermanent(ctx context.Context, banner_id int) (bool, error)

	RestoreAllBanner(ctx context.Context) (bool, error)
	DeleteAllBannerPermanent(ctx context.Context) (bool, error)
}

type CartRepository interface {
	FindCarts(
		ctx context.Context,
		req *requests.FindAllCarts,
	) ([]*db.GetCartsRow, error)

	CreateCart(
		ctx context.Context,
		req *requests.CartCreateRecord,
	) (*db.Cart, error)

	DeletePermanent(
		ctx context.Context,
		cart_id int,
	) (bool, error)

	DeleteAllPermanently(
		ctx context.Context,
		req *requests.DeleteCartRequest,
	) (bool, error)
}

type CategoryRepository interface {
	GetMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*db.GetMonthlyTotalPriceRow, error)
	GetYearlyTotalPrices(ctx context.Context, year int) ([]*db.GetYearlyTotalPriceRow, error)
	GetMonthlyTotalPriceById(
		ctx context.Context,
		req *requests.MonthTotalPriceCategory,
	) ([]*db.GetMonthlyTotalPriceByIdRow, error)

	GetYearlyTotalPricesById(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*db.GetYearlyTotalPriceByIdRow, error)
	GetMonthlyTotalPriceByMerchant(
		ctx context.Context,
		req *requests.MonthTotalPriceMerchant,
	) ([]*db.GetMonthlyTotalPriceByMerchantRow, error)
	GetYearlyTotalPricesByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*db.GetYearlyTotalPriceByMerchantRow, error)

	GetMonthPrice(ctx context.Context, year int) ([]*db.GetMonthlyCategoryRow, error)
	GetYearPrice(ctx context.Context, year int) ([]*db.GetYearlyCategoryRow, error)

	GetMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*db.GetMonthlyCategoryByMerchantRow, error)
	GetYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*db.GetYearlyCategoryByMerchantRow, error)

	GetMonthPriceById(ctx context.Context, req *requests.MonthPriceId) ([]*db.GetMonthlyCategoryByIdRow, error)
	GetYearPriceById(ctx context.Context, req *requests.YearPriceId) ([]*db.GetYearlyCategoryByIdRow, error)

	FindAllCategory(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesRow, error)

	FindByActive(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesActiveRow, error)

	FindByTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesTrashedRow, error)

	FindById(ctx context.Context, category_id int) (*db.GetCategoryByIDRow, error)

	FindByIdTrashed(ctx context.Context, category_id int) (*db.Category, error)

	CreateCategory(
		ctx context.Context,
		request *requests.CreateCategoryRequest,
	) (*db.CreateCategoryRow, error)

	UpdateCategory(
		ctx context.Context,
		request *requests.UpdateCategoryRequest,
	) (*db.UpdateCategoryRow, error)

	TrashedCategory(
		ctx context.Context,
		category_id int,
	) (*db.Category, error)

	RestoreCategory(
		ctx context.Context,
		category_id int,
	) (*db.Category, error)

	DeleteCategoryPermanently(
		ctx context.Context,
		category_id int,
	) (bool, error)

	RestoreAllCategories(ctx context.Context) (bool, error)
	DeleteAllPermanentCategories(ctx context.Context) (bool, error)
}

type MerchantRepository interface {
	FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsRow, error)

	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsActiveRow, error)

	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsTrashedRow, error)

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

type MerchantPoliciesRepository interface {
	FindAllMerchantPolicy(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesActiveRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesTrashedRow, error)
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

type MerchantAwardRepository interface {
	FindAllMerchants(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantCertificationsAndAwardsRow, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantCertificationsAndAwardsActiveRow, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantCertificationsAndAwardsTrashedRow, error)

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

type MerchantBusinessRepository interface {
	FindAllMerchants(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantsBusinessInformationRow, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantsBusinessInformationActiveRow, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantsBusinessInformationTrashedRow, error)

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

type MerchantDetailRepository interface {
	FindAllMerchants(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantDetailsRow, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantDetailsActiveRow, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantDetailsTrashedRow, error)

	FindById(
		ctx context.Context,
		user_id int,
	) (*db.GetMerchantDetailRow, error)
	FindByIdTrashed(ctx context.Context, user_id int) (*db.MerchantDetail, error)

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

type MerchantSocialLinkRepository interface {
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

type OrderRepository interface {
	GetMonthlyTotalRevenue(
		ctx context.Context,
		req *requests.MonthTotalRevenue,
	) ([]*db.GetMonthlyTotalRevenueRow, error)

	GetYearlyTotalRevenue(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyTotalRevenueRow, error)

	GetMonthlyTotalRevenueById(
		ctx context.Context,
		req *requests.MonthTotalRevenueOrder,
	) ([]*db.GetMonthlyTotalRevenueByIdRow, error)

	GetYearlyTotalRevenueById(
		ctx context.Context,
		req *requests.YearTotalRevenueOrder,
	) ([]*db.GetYearlyTotalRevenueByIdRow, error)

	GetMonthlyTotalRevenueByMerchant(
		ctx context.Context,
		req *requests.MonthTotalRevenueMerchant,
	) ([]*db.GetMonthlyTotalRevenueByMerchantRow, error)

	GetYearlyTotalRevenueByMerchant(
		ctx context.Context,
		req *requests.YearTotalRevenueMerchant,
	) ([]*db.GetYearlyTotalRevenueByMerchantRow, error)

	GetMonthlyOrder(
		ctx context.Context,
		year int,
	) ([]*db.GetMonthlyOrderRow, error)

	GetYearlyOrder(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyOrderRow, error)

	GetMonthlyOrderByMerchant(
		ctx context.Context,
		req *requests.MonthOrderMerchant,
	) ([]*db.GetMonthlyOrderByMerchantRow, error)

	GetYearlyOrderByMerchant(
		ctx context.Context,
		req *requests.YearOrderMerchant,
	) ([]*db.GetYearlyOrderByMerchantRow, error)

	FindAllOrders(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*db.GetOrdersRow, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*db.GetOrdersActiveRow, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*db.GetOrdersTrashedRow, error)

	FindByMerchant(
		ctx context.Context,
		req *requests.FindAllOrderByMerchant,
	) ([]*db.GetOrdersByMerchantRow, error)

	FindById(
		ctx context.Context,
		order_id int,
	) (*db.GetOrderByIDRow, error)

	FindByIdTrashed(ctx context.Context, user_id int) (*db.Order, error)

	CreateOrder(
		ctx context.Context,
		request *requests.CreateOrderRecordRequest,
	) (*db.CreateOrderRow, error)

	UpdateOrder(
		ctx context.Context,
		request *requests.UpdateOrderRecordRequest,
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

type OrderItemRepository interface {
	FindAllOrderItems(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*db.GetOrderItemsRow, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*db.GetOrderItemsActiveRow, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*db.GetOrderItemsTrashedRow, error)

	FindOrderItemByOrder(
		ctx context.Context,
		order_id int,
	) ([]*db.GetOrderItemsByOrderRow, error)
	FindOrderItemByOrderTrashed(
		ctx context.Context,
		order_id int,
	) ([]*db.OrderItem, error)

	CalculateTotalPrice(
		ctx context.Context,
		order_id int,
	) (*int32, error)

	CreateOrderItem(
		ctx context.Context,
		req *requests.CreateOrderItemRecordRequest,
	) (*db.CreateOrderItemRow, error)

	UpdateOrderItem(
		ctx context.Context,
		req *requests.UpdateOrderItemRecordRequest,
	) (*db.UpdateOrderItemRow, error)

	TrashedOrderItem(
		ctx context.Context,
		order_id int,
	) (*db.OrderItem, error)

	RestoreOrderItem(
		ctx context.Context,
		order_id int,
	) (*db.OrderItem, error)

	DeleteOrderItemPermanent(
		ctx context.Context,
		order_id int,
	) (bool, error)

	RestoreAllOrderItem(ctx context.Context) (bool, error)
	DeleteAllOrderPermanent(ctx context.Context) (bool, error)
}

type ProductRepository interface {
	FindAllProducts(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsActiveRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsTrashedRow, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*db.GetProductsByMerchantRow, error)
	FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*db.GetProductsByCategoryNameRow, error)
	FindById(ctx context.Context, product_id int) (*db.GetProductByIDRow, error)
	FindByIdTrashed(ctx context.Context, product_id int) (*db.GetProductByIdTrashedRow, error)
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

type TransactionRepository interface {
	GetMonthlyAmountSuccess(
		ctx context.Context,
		req *requests.MonthAmountTransaction,
	) ([]*db.GetMonthlyAmountTransactionSuccessRow, error)

	GetYearlyAmountSuccess(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyAmountTransactionSuccessRow, error)

	GetMonthlyAmountFailed(
		ctx context.Context,
		req *requests.MonthAmountTransaction,
	) ([]*db.GetMonthlyAmountTransactionFailedRow, error)

	GetYearlyAmountFailed(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyAmountTransactionFailedRow, error)

	GetMonthlyAmountSuccessByMerchant(
		ctx context.Context,
		req *requests.MonthAmountTransactionMerchant,
	) ([]*db.GetMonthlyAmountTransactionSuccessByMerchantRow, error)

	GetYearlyAmountSuccessByMerchant(
		ctx context.Context,
		req *requests.YearAmountTransactionMerchant,
	) ([]*db.GetYearlyAmountTransactionSuccessByMerchantRow, error)

	GetMonthlyAmountFailedByMerchant(
		ctx context.Context,
		req *requests.MonthAmountTransactionMerchant,
	) ([]*db.GetMonthlyAmountTransactionFailedByMerchantRow, error)

	GetYearlyAmountFailedByMerchant(
		ctx context.Context,
		req *requests.YearAmountTransactionMerchant,
	) ([]*db.GetYearlyAmountTransactionFailedByMerchantRow, error)

	GetMonthlyTransactionMethodSuccess(
		ctx context.Context,
		req *requests.MonthMethodTransaction,
	) ([]*db.GetMonthlyTransactionMethodsSuccessRow, error)

	GetYearlyTransactionMethodSuccess(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyTransactionMethodsSuccessRow, error)

	GetMonthlyTransactionMethodByMerchantSuccess(
		ctx context.Context,
		req *requests.MonthMethodTransactionMerchant,
	) ([]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow, error)

	GetYearlyTransactionMethodByMerchantSuccess(
		ctx context.Context,
		req *requests.YearMethodTransactionMerchant,
	) ([]*db.GetYearlyTransactionMethodsByMerchantSuccessRow, error)

	GetMonthlyTransactionMethodFailed(
		ctx context.Context,
		req *requests.MonthMethodTransaction,
	) ([]*db.GetMonthlyTransactionMethodsFailedRow, error)

	GetYearlyTransactionMethodFailed(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyTransactionMethodsFailedRow, error)

	GetMonthlyTransactionMethodByMerchantFailed(
		ctx context.Context,
		req *requests.MonthMethodTransactionMerchant,
	) ([]*db.GetMonthlyTransactionMethodsByMerchantFailedRow, error)

	GetYearlyTransactionMethodByMerchantFailed(
		ctx context.Context,
		req *requests.YearMethodTransactionMerchant,
	) ([]*db.GetYearlyTransactionMethodsByMerchantFailedRow, error)

	FindAllTransactions(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*db.GetTransactionsRow, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*db.GetTransactionsActiveRow, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*db.GetTransactionsTrashedRow, error)

	FindByMerchant(
		ctx context.Context,
		req *requests.FindAllTransactionByMerchant,
	) ([]*db.GetTransactionByMerchantRow, error)

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

	TrashTransaction(
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

type ShippingAddressRepository interface {
	FindAllShippingAddress(
		ctx context.Context,
		req *requests.FindAllShippingAddress,
	) ([]*db.GetShippingAddressRow, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllShippingAddress,
	) ([]*db.GetShippingAddressActiveRow, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllShippingAddress,
	) ([]*db.GetShippingAddressTrashedRow, error)

	FindByOrder(
		ctx context.Context,
		shipping_id int,
	) (*db.GetShippingAddressByOrderIDRow, error)

	FindById(
		ctx context.Context,
		shipping_id int,
	) (*db.GetShippingByIDRow, error)

	CreateShippingAddress(
		ctx context.Context,
		request *requests.CreateShippingAddressRequest,
	) (*db.CreateShippingAddressRow, error)

	UpdateShippingAddress(
		ctx context.Context,
		request *requests.UpdateShippingAddressRequest,
	) (*db.UpdateShippingAddressRow, error)

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

type SliderRepository interface {
	FindAllSlider(
		ctx context.Context,
		req *requests.FindAllSlider,
	) ([]*db.GetSlidersRow, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllSlider,
	) ([]*db.GetSlidersActiveRow, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllSlider,
	) ([]*db.GetSlidersTrashedRow, error)

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

	RestoreAllSlider(ctx context.Context) (bool, error)
	DeleteAllPermanentSlider(ctx context.Context) (bool, error)
}

type ReviewRepository interface {
	FindAllReview(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsRow, error)
	FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*db.GetReviewByProductIdRow, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*db.GetReviewByMerchantIdRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsActiveRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsTrashedRow, error)
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

type ReviewDetailRepository interface {
	FindAllReviews(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsActiveRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsTrashedRow, error)
	FindById(ctx context.Context, user_id int) (*db.GetReviewDetailRow, error)
	FindByIdTrashed(ctx context.Context, user_id int) (*db.ReviewDetail, error)
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
