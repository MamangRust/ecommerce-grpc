package recordmapper

import (
	"ecommerce/internal/domain/record"
	db "ecommerce/pkg/database/schema"
)

type UserRecordMapping interface {
	ToUserRecord(user *db.User) *record.UserRecord
	ToUserRecordPagination(user *db.GetUsersRow) *record.UserRecord
	ToUsersRecordPagination(users []*db.GetUsersRow) []*record.UserRecord

	ToUserRecordActivePagination(user *db.GetUsersActiveRow) *record.UserRecord
	ToUsersRecordActivePagination(users []*db.GetUsersActiveRow) []*record.UserRecord
	ToUserRecordTrashedPagination(user *db.GetUserTrashedRow) *record.UserRecord
	ToUsersRecordTrashedPagination(users []*db.GetUserTrashedRow) []*record.UserRecord
}

type RoleRecordMapping interface {
	ToRoleRecord(role *db.Role) *record.RoleRecord
	ToRolesRecord(roles []*db.Role) []*record.RoleRecord

	ToRoleRecordAll(role *db.GetRolesRow) *record.RoleRecord
	ToRolesRecordAll(roles []*db.GetRolesRow) []*record.RoleRecord

	ToRoleRecordActive(role *db.GetActiveRolesRow) *record.RoleRecord
	ToRolesRecordActive(roles []*db.GetActiveRolesRow) []*record.RoleRecord
	ToRoleRecordTrashed(role *db.GetTrashedRolesRow) *record.RoleRecord
	ToRolesRecordTrashed(roles []*db.GetTrashedRolesRow) []*record.RoleRecord
}

type BannerRecordMapping interface {
	ToBannerRecord(banner *db.Banner) *record.BannerRecord
	ToBannerRecordPagination(banner *db.GetBannersRow) *record.BannerRecord
	ToBannersRecordPagination(Merchants []*db.GetBannersRow) []*record.BannerRecord
	ToBannerRecordActivePagination(banner *db.GetBannersActiveRow) *record.BannerRecord
	ToBannersRecordActivePagination(Merchants []*db.GetBannersActiveRow) []*record.BannerRecord
	ToBannerRecordTrashedPagination(banner *db.GetBannersTrashedRow) *record.BannerRecord
	ToBannersRecordTrashedPagination(banner []*db.GetBannersTrashedRow) []*record.BannerRecord
}

type UserRoleRecordMapping interface {
	ToUserRoleRecord(userRole *db.UserRole) *record.UserRoleRecord
}

type RefreshTokenRecordMapping interface {
	ToRefreshTokenRecord(refreshToken *db.RefreshToken) *record.RefreshTokenRecord
	ToRefreshTokensRecord(refreshTokens []*db.RefreshToken) []*record.RefreshTokenRecord
}

type CategoryRecordMapper interface {
	ToCategoryMonthlyTotalPrice(c *db.GetMonthlyTotalPriceRow) *record.CategoriesMonthlyTotalPriceRecord
	ToCategoryMonthlyTotalPrices(c []*db.GetMonthlyTotalPriceRow) []*record.CategoriesMonthlyTotalPriceRecord
	ToCategoryYearlyTotalPrice(c *db.GetYearlyTotalPriceRow) *record.CategoriesYearlyTotalPriceRecord
	ToCategoryYearlyTotalPrices(c []*db.GetYearlyTotalPriceRow) []*record.CategoriesYearlyTotalPriceRecord
	ToCategoryMonthlyTotalPriceById(c *db.GetMonthlyTotalPriceByIdRow) *record.CategoriesMonthlyTotalPriceRecord
	ToCategoryMonthlyTotalPricesById(c []*db.GetMonthlyTotalPriceByIdRow) []*record.CategoriesMonthlyTotalPriceRecord
	ToCategoryYearlyTotalPriceById(c *db.GetYearlyTotalPriceByIdRow) *record.CategoriesYearlyTotalPriceRecord
	ToCategoryYearlyTotalPricesById(c []*db.GetYearlyTotalPriceByIdRow) []*record.CategoriesYearlyTotalPriceRecord

	ToCategoryMonthlyTotalPriceByMerchant(c *db.GetMonthlyTotalPriceByMerchantRow) *record.CategoriesMonthlyTotalPriceRecord
	ToCategoryMonthlyTotalPricesByMerchant(c []*db.GetMonthlyTotalPriceByMerchantRow) []*record.CategoriesMonthlyTotalPriceRecord
	ToCategoryYearlyTotalPriceByMerchant(c *db.GetYearlyTotalPriceByMerchantRow) *record.CategoriesYearlyTotalPriceRecord
	ToCategoryYearlyTotalPricesByMerchant(c []*db.GetYearlyTotalPriceByMerchantRow) []*record.CategoriesYearlyTotalPriceRecord

	ToCategoryMonthlyPriceById(category *db.GetMonthlyCategoryByIdRow) *record.CategoriesMonthPriceRecord
	ToCategoryMonthlyPricesById(c []*db.GetMonthlyCategoryByIdRow) []*record.CategoriesMonthPriceRecord
	ToCategoryYearlyPriceById(category *db.GetYearlyCategoryByIdRow) *record.CategoriesYearPriceRecord
	ToCategoryYearlyPricesById(c []*db.GetYearlyCategoryByIdRow) []*record.CategoriesYearPriceRecord

	ToCategoryMonthlyPrice(category *db.GetMonthlyCategoryRow) *record.CategoriesMonthPriceRecord
	ToCategoryMonthlyPrices(c []*db.GetMonthlyCategoryRow) []*record.CategoriesMonthPriceRecord
	ToCategoryYearlyPrice(category *db.GetYearlyCategoryRow) *record.CategoriesYearPriceRecord
	ToCategoryYearlyPrices(c []*db.GetYearlyCategoryRow) []*record.CategoriesYearPriceRecord
	ToCategoryMonthlyPriceByMerchant(category *db.GetMonthlyCategoryByMerchantRow) *record.CategoriesMonthPriceRecord
	ToCategoryMonthlyPricesByMerchant(c []*db.GetMonthlyCategoryByMerchantRow) []*record.CategoriesMonthPriceRecord
	ToCategoryYearlyPriceByMerchant(category *db.GetYearlyCategoryByMerchantRow) *record.CategoriesYearPriceRecord
	ToCategoryYearlyPricesByMerchant(c []*db.GetYearlyCategoryByMerchantRow) []*record.CategoriesYearPriceRecord

	ToCategoryRecord(category *db.Category) *record.CategoriesRecord
	ToCategoryRecordPagination(category *db.GetCategoriesRow) *record.CategoriesRecord
	ToCategoriesRecordPagination(categories []*db.GetCategoriesRow) []*record.CategoriesRecord
	ToCategoryRecordActivePagination(category *db.GetCategoriesActiveRow) *record.CategoriesRecord
	ToCategoriesRecordActivePagination(categories []*db.GetCategoriesActiveRow) []*record.CategoriesRecord
	ToCategoryRecordTrashedPagination(category *db.GetCategoriesTrashedRow) *record.CategoriesRecord
	ToCategoriesRecordTrashedPagination(categories []*db.GetCategoriesTrashedRow) []*record.CategoriesRecord
}

type MerchantRecordMapping interface {
	ToMerchantRecord(Merchant *db.Merchant) *record.MerchantRecord
	ToMerchantRecordPagination(Merchant *db.GetMerchantsRow) *record.MerchantRecord
	ToMerchantsRecordPagination(Merchants []*db.GetMerchantsRow) []*record.MerchantRecord

	ToMerchantRecordActivePagination(Merchant *db.GetMerchantsActiveRow) *record.MerchantRecord
	ToMerchantsRecordActivePagination(Merchants []*db.GetMerchantsActiveRow) []*record.MerchantRecord
	ToMerchantRecordTrashedPagination(Merchant *db.GetMerchantsTrashedRow) *record.MerchantRecord
	ToMerchantsRecordTrashedPagination(Merchants []*db.GetMerchantsTrashedRow) []*record.MerchantRecord
}

type MerchantPolicyMapping interface {
	ToMerchantPolicyRecord(m *db.MerchantPolicy) *record.MerchantPoliciesRecord
	ToMerchantPolicyRecordPagination(m *db.GetMerchantPoliciesRow) *record.MerchantPoliciesRecord
	ToMerchantPolicysRecordPagination(MerchantPolicys []*db.GetMerchantPoliciesRow) []*record.MerchantPoliciesRecord
	ToMerchantPolicyRecordActivePagination(m *db.GetMerchantPoliciesActiveRow) *record.MerchantPoliciesRecord
	ToMerchantPolicysRecordActivePagination(MerchantPolicys []*db.GetMerchantPoliciesActiveRow) []*record.MerchantPoliciesRecord
	ToMerchantPolicyRecordTrashedPagination(m *db.GetMerchantPoliciesTrashedRow) *record.MerchantPoliciesRecord
	ToMerchantPolicysRecordTrashedPagination(MerchantPolicys []*db.GetMerchantPoliciesTrashedRow) []*record.MerchantPoliciesRecord
}

type MerchantAwardMapping interface {
	ToMerchantAwardRecord(m *db.MerchantCertificationsAndAward) *record.MerchantAwardRecord
	ToMerchantAwardRecordPagination(m *db.GetMerchantCertificationsAndAwardsRow) *record.MerchantAwardRecord
	ToMerchantAwardsRecordPagination(Merchants []*db.GetMerchantCertificationsAndAwardsRow) []*record.MerchantAwardRecord
	ToMerchantAwardRecordActivePagination(m *db.GetMerchantCertificationsAndAwardsActiveRow) *record.MerchantAwardRecord
	ToMerchantAwardsRecordActivePagination(Merchants []*db.GetMerchantCertificationsAndAwardsActiveRow) []*record.MerchantAwardRecord
	ToMerchantAwardRecordTrashedPagination(m *db.GetMerchantCertificationsAndAwardsTrashedRow) *record.MerchantAwardRecord
	ToMerchantAwardsRecordTrashedPagination(m []*db.GetMerchantCertificationsAndAwardsTrashedRow) []*record.MerchantAwardRecord
}

type MerchantBusinessMapping interface {
	ToMerchantBusinessRecord(m *db.MerchantBusinessInformation) *record.MerchantBusinessRecord
	ToMerchantBusinessRecordPagination(m *db.GetMerchantsBusinessInformationRow) *record.MerchantBusinessRecord
	ToMerchantBusinesssRecordPagination(Merchants []*db.GetMerchantsBusinessInformationRow) []*record.MerchantBusinessRecord
	ToMerchantBusinessRecordActivePagination(m *db.GetMerchantsBusinessInformationActiveRow) *record.MerchantBusinessRecord
	ToMerchantBusinesssRecordActivePagination(Merchants []*db.GetMerchantsBusinessInformationActiveRow) []*record.MerchantBusinessRecord
	ToMerchantBusinessRecordTrashedPagination(m *db.GetMerchantsBusinessInformationTrashedRow) *record.MerchantBusinessRecord
	ToMerchantBusinesssRecordTrashedPagination(Merchants []*db.GetMerchantsBusinessInformationTrashedRow) []*record.MerchantBusinessRecord
}

type MerchantDetailMapping interface {
	ToMerchantDetailRecord(detail *db.MerchantDetail) *record.MerchantDetailRecord
	ToMerchantDetailRelationRecord(detail *db.GetMerchantDetailRow) *record.MerchantDetailRecord
	ToMerchantDetailRecordPagination(detail *db.GetMerchantDetailsRow) *record.MerchantDetailRecord
	ToMerchantDetailsRecordPagination(MerchantDetails []*db.GetMerchantDetailsRow) []*record.MerchantDetailRecord
	ToMerchantDetailRecordActivePagination(detail *db.GetMerchantDetailsActiveRow) *record.MerchantDetailRecord
	ToMerchantDetailsRecordActivePagination(MerchantDetails []*db.GetMerchantDetailsActiveRow) []*record.MerchantDetailRecord
	ToMerchantDetailRecordTrashedPagination(detail *db.GetMerchantDetailsTrashedRow) *record.MerchantDetailRecord
	ToMerchantDetailsRecordTrashedPagination(MerchantDetails []*db.GetMerchantDetailsTrashedRow) []*record.MerchantDetailRecord
}

type OrderItemRecordMapping interface {
	ToOrderItemRecord(orderItems *db.OrderItem) *record.OrderItemRecord
	ToOrderItemsRecord(orders []*db.OrderItem) []*record.OrderItemRecord

	ToOrderItemRecordPagination(OrderItem *db.GetOrderItemsRow) *record.OrderItemRecord
	ToOrderItemsRecordPagination(OrderItem []*db.GetOrderItemsRow) []*record.OrderItemRecord

	ToOrderItemRecordActivePagination(OrderItem *db.GetOrderItemsActiveRow) *record.OrderItemRecord
	ToOrderItemsRecordActivePagination(OrderItem []*db.GetOrderItemsActiveRow) []*record.OrderItemRecord
	ToOrderItemRecordTrashedPagination(OrderItem *db.GetOrderItemsTrashedRow) *record.OrderItemRecord
	ToOrderItemsRecordTrashedPagination(OrderItem []*db.GetOrderItemsTrashedRow) []*record.OrderItemRecord
}

type OrderRecordMapping interface {
	ToOrderMonthlyTotalRevenue(c *db.GetMonthlyTotalRevenueRow) *record.OrderMonthlyTotalRevenueRecord
	ToOrderMonthlyTotalRevenues(c []*db.GetMonthlyTotalRevenueRow) []*record.OrderMonthlyTotalRevenueRecord
	ToOrderYearlyTotalRevenue(c *db.GetYearlyTotalRevenueRow) *record.OrderYearlyTotalRevenueRecord
	ToOrderYearlyTotalRevenues(c []*db.GetYearlyTotalRevenueRow) []*record.OrderYearlyTotalRevenueRecord
	ToOrderMonthlyTotalRevenueById(c *db.GetMonthlyTotalRevenueByIdRow) *record.OrderMonthlyTotalRevenueRecord
	ToOrderMonthlyTotalRevenuesById(c []*db.GetMonthlyTotalRevenueByIdRow) []*record.OrderMonthlyTotalRevenueRecord
	ToOrderYearlyTotalRevenueById(c *db.GetYearlyTotalRevenueByIdRow) *record.OrderYearlyTotalRevenueRecord
	ToOrderYearlyTotalRevenuesById(c []*db.GetYearlyTotalRevenueByIdRow) []*record.OrderYearlyTotalRevenueRecord
	ToOrderMonthlyTotalRevenueByMerchant(c *db.GetMonthlyTotalRevenueByMerchantRow) *record.OrderMonthlyTotalRevenueRecord
	ToOrderMonthlyTotalRevenuesByMerchant(c []*db.GetMonthlyTotalRevenueByMerchantRow) []*record.OrderMonthlyTotalRevenueRecord
	ToOrderYearlyTotalRevenueByMerchant(c *db.GetYearlyTotalRevenueByMerchantRow) *record.OrderYearlyTotalRevenueRecord
	ToOrderYearlyTotalRevenuesByMerchant(c []*db.GetYearlyTotalRevenueByMerchantRow) []*record.OrderYearlyTotalRevenueRecord

	ToOrderMonthlyPrice(category *db.GetMonthlyOrderRow) *record.OrderMonthlyRecord
	ToOrderMonthlyPrices(c []*db.GetMonthlyOrderRow) []*record.OrderMonthlyRecord
	ToOrderYearlyPrice(category *db.GetYearlyOrderRow) *record.OrderYearlyRecord
	ToOrderYearlyPrices(c []*db.GetYearlyOrderRow) []*record.OrderYearlyRecord
	ToOrderMonthlyPriceByMerchant(category *db.GetMonthlyOrderByMerchantRow) *record.OrderMonthlyRecord
	ToOrderMonthlyPricesByMerchant(c []*db.GetMonthlyOrderByMerchantRow) []*record.OrderMonthlyRecord
	ToOrderYearlyPriceByMerchant(category *db.GetYearlyOrderByMerchantRow) *record.OrderYearlyRecord
	ToOrderYearlyPricesByMerchant(c []*db.GetYearlyOrderByMerchantRow) []*record.OrderYearlyRecord

	ToOrderRecord(order *db.Order) *record.OrderRecord
	ToOrdersRecord(orders []*db.Order) []*record.OrderRecord
	ToOrderRecordPagination(order *db.GetOrdersRow) *record.OrderRecord
	ToOrdersRecordPagination(orders []*db.GetOrdersRow) []*record.OrderRecord
	ToOrderRecordActivePagination(order *db.GetOrdersActiveRow) *record.OrderRecord
	ToOrdersRecordActivePagination(orders []*db.GetOrdersActiveRow) []*record.OrderRecord
	ToOrderRecordTrashedPagination(order *db.GetOrdersTrashedRow) *record.OrderRecord
	ToOrdersRecordTrashedPagination(orders []*db.GetOrdersTrashedRow) []*record.OrderRecord

	ToOrderRecordByMerchantPagination(order *db.GetOrdersByMerchantRow) *record.OrderRecord
	ToOrdersRecordByMerchantPagination(orders []*db.GetOrdersByMerchantRow) []*record.OrderRecord
}

type ProductRecordMapping interface {
	ToProductRecord(product *db.Product) *record.ProductRecord
	ToProductsRecord(products []*db.Product) []*record.ProductRecord
	ToProductRecordPagination(product *db.GetProductsRow) *record.ProductRecord
	ToProductsRecordPagination(products []*db.GetProductsRow) []*record.ProductRecord
	ToProductRecordActivePagination(product *db.GetProductsActiveRow) *record.ProductRecord
	ToProductsRecordActivePagination(products []*db.GetProductsActiveRow) []*record.ProductRecord
	ToProductRecordTrashedPagination(product *db.GetProductsTrashedRow) *record.ProductRecord
	ToProductsRecordTrashedPagination(products []*db.GetProductsTrashedRow) []*record.ProductRecord

	ToProductRecordMerchantPagination(product *db.GetProductsByMerchantRow) *record.ProductRecord
	ToProductsRecordMerchantPagination(products []*db.GetProductsByMerchantRow) []*record.ProductRecord

	ToProductRecordCategoryPagination(product *db.GetProductsByCategoryNameRow) *record.ProductRecord
	ToProductsRecordCategoryPagination(products []*db.GetProductsByCategoryNameRow) []*record.ProductRecord
}

type TransactionRecordMapping interface {
	ToTransactionMonthAmountSuccess(row *db.GetMonthlyAmountTransactionSuccessRow) *record.TransactionMonthlyAmountSuccessRecord
	ToTransactionMonthlyAmountSuccess(rows []*db.GetMonthlyAmountTransactionSuccessRow) []*record.TransactionMonthlyAmountSuccessRecord
	ToTransactionYearAmountSuccess(row *db.GetYearlyAmountTransactionSuccessRow) *record.TransactionYearlyAmountSuccessRecord
	ToTransactionYearlyAmountSuccess(rows []*db.GetYearlyAmountTransactionSuccessRow) []*record.TransactionYearlyAmountSuccessRecord
	ToTransactionMonthAmountFailed(row *db.GetMonthlyAmountTransactionFailedRow) *record.TransactionMonthlyAmountFailedRecord
	ToTransactionMonthlyAmountFailed(rows []*db.GetMonthlyAmountTransactionFailedRow) []*record.TransactionMonthlyAmountFailedRecord
	ToTransactionYearAmountFailed(row *db.GetYearlyAmountTransactionFailedRow) *record.TransactionYearlyAmountFailedRecord
	ToTransactionYearlyAmountFailed(rows []*db.GetYearlyAmountTransactionFailedRow) []*record.TransactionYearlyAmountFailedRecord

	ToTransactionMonthAmountSuccessByMerchant(row *db.GetMonthlyAmountTransactionSuccessByMerchantRow) *record.TransactionMonthlyAmountSuccessRecord
	ToTransactionMonthlyAmountSuccessByMerchant(rows []*db.GetMonthlyAmountTransactionSuccessByMerchantRow) []*record.TransactionMonthlyAmountSuccessRecord
	ToTransactionYearAmountSuccessByMerchant(row *db.GetYearlyAmountTransactionSuccessByMerchantRow) *record.TransactionYearlyAmountSuccessRecord
	ToTransactionYearlyAmountSuccessByMerchant(rows []*db.GetYearlyAmountTransactionSuccessByMerchantRow) []*record.TransactionYearlyAmountSuccessRecord
	ToTransactionMonthAmountFailedByMerchant(row *db.GetMonthlyAmountTransactionFailedByMerchantRow) *record.TransactionMonthlyAmountFailedRecord
	ToTransactionMonthlyAmountFailedByMerchant(rows []*db.GetMonthlyAmountTransactionFailedByMerchantRow) []*record.TransactionMonthlyAmountFailedRecord
	ToTransactionYearAmountFailedByMerchant(row *db.GetYearlyAmountTransactionFailedByMerchantRow) *record.TransactionYearlyAmountFailedRecord
	ToTransactionYearlyAmountFailedByMerchant(rows []*db.GetYearlyAmountTransactionFailedByMerchantRow) []*record.TransactionYearlyAmountFailedRecord

	ToTransactionMonthMethodSuccess(row *db.GetMonthlyTransactionMethodsSuccessRow) *record.TransactionMonthlyMethodRecord
	ToTransactionMonthlyMethodSuccess(rows []*db.GetMonthlyTransactionMethodsSuccessRow) []*record.TransactionMonthlyMethodRecord

	ToTransactionMonthMethodFailed(row *db.GetMonthlyTransactionMethodsFailedRow) *record.TransactionMonthlyMethodRecord
	ToTransactionMonthlyMethodFailed(rows []*db.GetMonthlyTransactionMethodsFailedRow) []*record.TransactionMonthlyMethodRecord

	ToTransactionYearMethodSuccess(row *db.GetYearlyTransactionMethodsSuccessRow) *record.TransactionYearlyMethodRecord
	ToTransactionYearlyMethodSuccess(rows []*db.GetYearlyTransactionMethodsSuccessRow) []*record.TransactionYearlyMethodRecord

	ToTransactionYearMethodFailed(row *db.GetYearlyTransactionMethodsFailedRow) *record.TransactionYearlyMethodRecord
	ToTransactionYearlyMethodFailed(rows []*db.GetYearlyTransactionMethodsFailedRow) []*record.TransactionYearlyMethodRecord

	ToTransactionMonthMethodByMerchantSuccess(row *db.GetMonthlyTransactionMethodsByMerchantSuccessRow) *record.TransactionMonthlyMethodRecord
	ToTransactionMonthlyByMerchantMethodSuccess(rows []*db.GetMonthlyTransactionMethodsByMerchantSuccessRow) []*record.TransactionMonthlyMethodRecord

	ToTransactionMonthMethodByMerchantFailed(row *db.GetMonthlyTransactionMethodsByMerchantFailedRow) *record.TransactionMonthlyMethodRecord
	ToTransactionMonthlyByMerchantMethodFailed(rows []*db.GetMonthlyTransactionMethodsByMerchantFailedRow) []*record.TransactionMonthlyMethodRecord

	ToTransactionYearMethodByMerchantSuccess(row *db.GetYearlyTransactionMethodsByMerchantSuccessRow) *record.TransactionYearlyMethodRecord
	ToTransactionYearlyMethodByMerchantSuccess(rows []*db.GetYearlyTransactionMethodsByMerchantSuccessRow) []*record.TransactionYearlyMethodRecord

	ToTransactionYearMethodByMerchantFailed(row *db.GetYearlyTransactionMethodsByMerchantFailedRow) *record.TransactionYearlyMethodRecord
	ToTransactionYearlyMethodByMerchantFailed(rows []*db.GetYearlyTransactionMethodsByMerchantFailedRow) []*record.TransactionYearlyMethodRecord

	ToTransactionRecord(transaction *db.Transaction) *record.TransactionRecord
	ToTransactionsRecord(transactions []*db.Transaction) []*record.TransactionRecord
	ToTransactionRecordPagination(transaction *db.GetTransactionsRow) *record.TransactionRecord
	ToTransactionsRecordPagination(transaction []*db.GetTransactionsRow) []*record.TransactionRecord
	ToTransactionRecordActivePagination(transaction *db.GetTransactionsActiveRow) *record.TransactionRecord
	ToTransactionsRecordActivePagination(transactions []*db.GetTransactionsActiveRow) []*record.TransactionRecord
	ToTransactionRecordTrashedPagination(transaction *db.GetTransactionsTrashedRow) *record.TransactionRecord
	ToTransactionsRecordTrashedPagination(products []*db.GetTransactionsTrashedRow) []*record.TransactionRecord

	ToTransactionMerchantRecordPagination(transaction *db.GetTransactionByMerchantRow) *record.TransactionRecord
	ToTransactionMerchantsRecordPagination(products []*db.GetTransactionByMerchantRow) []*record.TransactionRecord
}

type CartRecordMapping interface {
	ToCartRecord(cart *db.Cart) *record.CartRecord
	ToCartRecordPagination(cart *db.GetCartsRow) *record.CartRecord
	ToCartsRecordPagination(carts []*db.GetCartsRow) []*record.CartRecord
}

type ReviewRecordMapping interface {
	ToReviewRecord(review *db.Review) *record.ReviewRecord
	ToReviewRecordPagination(review *db.GetReviewsRow) *record.ReviewRecord
	ToReviewsRecordPagination(reviews []*db.GetReviewsRow) []*record.ReviewRecord

	ToReviewProductRecordPagination(review *db.GetReviewByProductIdRow) *record.ReviewsDetailRecord
	ToReviewsProductRecordPagination(reviews []*db.GetReviewByProductIdRow) []*record.ReviewsDetailRecord

	ToReviewMerchantRecordPagination(review *db.GetReviewByMerchantIdRow) *record.ReviewsDetailRecord
	ToReviewsMerchantRecordPagination(reviews []*db.GetReviewByMerchantIdRow) []*record.ReviewsDetailRecord

	ToReviewRecordActivePagination(review *db.GetReviewsActiveRow) *record.ReviewRecord
	ToReviewsRecordActivePagination(reviews []*db.GetReviewsActiveRow) []*record.ReviewRecord
	ToReviewRecordTrashedPagination(review *db.GetReviewsTrashedRow) *record.ReviewRecord
	ToReviewsRecordTrashedPagination(reviews []*db.GetReviewsTrashedRow) []*record.ReviewRecord
}

type ReviewDetailRecordMapping interface {
	ToReviewDetailRecord(detail *db.ReviewDetail) *record.ReviewDetailRecord
	ToReviewDetailRecordPagination(detail *db.GetReviewDetailsRow) *record.ReviewDetailRecord
	ToReviewDetailsRecordPagination(Merchants []*db.GetReviewDetailsRow) []*record.ReviewDetailRecord
	ToReviewDetailRecordActivePagination(detail *db.GetReviewDetailsActiveRow) *record.ReviewDetailRecord
	ToReviewDetailsRecordActivePagination(m []*db.GetReviewDetailsActiveRow) []*record.ReviewDetailRecord
	ToReviewDetailRecordTrashedPagination(detail *db.GetReviewDetailsTrashedRow) *record.ReviewDetailRecord
	ToReviewDetailsRecordTrashedPagination(Merchants []*db.GetReviewDetailsTrashedRow) []*record.ReviewDetailRecord
}

type ShippingAddressMapping interface {
	ToShippingAddressRecord(address *db.ShippingAddress) *record.ShippingAddressRecord
	ToShippingAddressRecordPagination(address *db.GetShippingAddressRow) *record.ShippingAddressRecord
	ToShippingAddresssRecordPagination(addresses []*db.GetShippingAddressRow) []*record.ShippingAddressRecord
	ToShippingAddressRecordActivePagination(address *db.GetShippingAddressActiveRow) *record.ShippingAddressRecord
	ToShippingAddresssRecordActivePagination(sliders []*db.GetShippingAddressActiveRow) []*record.ShippingAddressRecord
	ToShippingAddressRecordTrashedPagination(address *db.GetShippingAddressTrashedRow) *record.ShippingAddressRecord
	ToShippingAddresssRecordTrashedPagination(sliders []*db.GetShippingAddressTrashedRow) []*record.ShippingAddressRecord
}

type SliderMapping interface {
	ToSliderRecord(slider *db.Slider) *record.SliderRecord
	ToSliderRecordPagination(slider *db.GetSlidersRow) *record.SliderRecord
	ToSlidersRecordPagination(sliders []*db.GetSlidersRow) []*record.SliderRecord
	ToSliderRecordActivePagination(slider *db.GetSlidersActiveRow) *record.SliderRecord
	ToSlidersRecordActivePagination(sliders []*db.GetSlidersActiveRow) []*record.SliderRecord
	ToSliderRecordTrashedPagination(slider *db.GetSlidersTrashedRow) *record.SliderRecord
	ToSlidersRecordTrashedPagination(sliders []*db.GetSlidersTrashedRow) []*record.SliderRecord
}
