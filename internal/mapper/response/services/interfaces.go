package response_service

import (
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/response"
)

type UserResponseMapper interface {
	ToUserResponse(user *record.UserRecord) *response.UserResponse
	ToUsersResponse(users []*record.UserRecord) []*response.UserResponse

	ToUserResponseDeleteAt(user *record.UserRecord) *response.UserResponseDeleteAt
	ToUsersResponseDeleteAt(users []*record.UserRecord) []*response.UserResponseDeleteAt
}

type RoleResponseMapper interface {
	ToRoleResponse(role *record.RoleRecord) *response.RoleResponse
	ToRolesResponse(roles []*record.RoleRecord) []*response.RoleResponse

	ToRoleResponseDeleteAt(role *record.RoleRecord) *response.RoleResponseDeleteAt
	ToRolesResponseDeleteAt(roles []*record.RoleRecord) []*response.RoleResponseDeleteAt
}

type RefreshTokenResponseMapper interface {
	ToRefreshTokenResponse(refresh *record.RefreshTokenRecord) *response.RefreshTokenResponse
	ToRefreshTokenResponses(refreshs []*record.RefreshTokenRecord) []*response.RefreshTokenResponse
}

type BannerResponseMapper interface {
	ToBannerResponse(banner *record.BannerRecord) *response.BannerResponse
	ToBannersResponse(merchants []*record.BannerRecord) []*response.BannerResponse
	ToBannerResponseDeleteAt(banner *record.BannerRecord) *response.BannerResponseDeleteAt
	ToBannersResponseDeleteAt(banners []*record.BannerRecord) []*response.BannerResponseDeleteAt
}

type CategoryResponseMapper interface {
	ToCategoryMonthlyTotalPrice(c *record.CategoriesMonthlyTotalPriceRecord) *response.CategoriesMonthlyTotalPriceResponse
	ToCategoryMonthlyTotalPrices(c []*record.CategoriesMonthlyTotalPriceRecord) []*response.CategoriesMonthlyTotalPriceResponse
	ToCategoryYearlyTotalPrice(c *record.CategoriesYearlyTotalPriceRecord) *response.CategoriesYearlyTotalPriceResponse
	ToCategoryYearlyTotalPrices(c []*record.CategoriesYearlyTotalPriceRecord) []*response.CategoriesYearlyTotalPriceResponse

	ToCategoryMonthlyPrice(category *record.CategoriesMonthPriceRecord) *response.CategoryMonthPriceResponse
	ToCategoryMonthlyPrices(c []*record.CategoriesMonthPriceRecord) []*response.CategoryMonthPriceResponse
	ToCategoryYearlyPrice(category *record.CategoriesYearPriceRecord) *response.CategoryYearPriceResponse
	ToCategoryYearlyPrices(c []*record.CategoriesYearPriceRecord) []*response.CategoryYearPriceResponse

	ToCategoryResponse(category *record.CategoriesRecord) *response.CategoryResponse
	ToCategorysResponse(categories []*record.CategoriesRecord) []*response.CategoryResponse
	ToCategoryResponseDeleteAt(category *record.CategoriesRecord) *response.CategoryResponseDeleteAt
	ToCategorysResponseDeleteAt(categories []*record.CategoriesRecord) []*response.CategoryResponseDeleteAt
}

type MerchantResponseMapper interface {
	ToMerchantResponse(merchant *record.MerchantRecord) *response.MerchantResponse
	ToMerchantsResponse(merchants []*record.MerchantRecord) []*response.MerchantResponse
	ToMerchantResponseDeleteAt(merchant *record.MerchantRecord) *response.MerchantResponseDeleteAt
	ToMerchantsResponseDeleteAt(merchants []*record.MerchantRecord) []*response.MerchantResponseDeleteAt
}

type MerchantAwardResponseMapper interface {
	ToMerchantAwardResponse(merchant *record.MerchantAwardRecord) *response.MerchantAwardResponse
	ToMerchantsAwardResponse(merchants []*record.MerchantAwardRecord) []*response.MerchantAwardResponse
	ToMerchantAwardResponseDeleteAt(merchant *record.MerchantAwardRecord) *response.MerchantAwardResponseDeleteAt
	ToMerchantsAwardResponseDeleteAt(merchants []*record.MerchantAwardRecord) []*response.MerchantAwardResponseDeleteAt
}

type MerchantBusinessResponseMapper interface {
	ToMerchantBusinessResponse(merchant *record.MerchantBusinessRecord) *response.MerchantBusinessResponse
	ToMerchantBusinessResponseRelation(merchant *record.MerchantBusinessRecord) *response.MerchantBusinessResponse
	ToMerchantsBusinessResponse(merchants []*record.MerchantBusinessRecord) []*response.MerchantBusinessResponse
	ToMerchantBusinessResponseDeleteAt(merchant *record.MerchantBusinessRecord) *response.MerchantBusinessResponseDeleteAt
	ToMerchantsBusinessResponseDeleteAt(merchants []*record.MerchantBusinessRecord) []*response.MerchantBusinessResponseDeleteAt
}

type MerchantDetailResponseMapper interface {
	ToMerchantDetailResponse(merchant *record.MerchantDetailRecord) *response.MerchantDetailResponse
	ToMerchantDetailRelationResponse(merchant *record.MerchantDetailRecord) *response.MerchantDetailResponse
	ToMerchantsDetailResponse(merchants []*record.MerchantDetailRecord) []*response.MerchantDetailResponse
	ToMerchantDetailResponseDeleteAt(merchant *record.MerchantDetailRecord) *response.MerchantDetailResponseDeleteAt
	ToMerchantsDetailResponseDeleteAt(merchants []*record.MerchantDetailRecord) []*response.MerchantDetailResponseDeleteAt
}

type MerchantPolicyResponseMapper interface {
	ToMerchantPolicyResponse(merchant *record.MerchantPoliciesRecord) *response.MerchantPoliciesResponse
	ToMerchantsPolicyResponse(merchants []*record.MerchantPoliciesRecord) []*response.MerchantPoliciesResponse
	ToMerchantPolicyResponseDeleteAt(merchant *record.MerchantPoliciesRecord) *response.MerchantPoliciesResponseDeleteAt
	ToMerchantsPolicyResponseDeleteAt(merchants []*record.MerchantPoliciesRecord) []*response.MerchantPoliciesResponseDeleteAt
}

type OrderResponseMapper interface {
	ToOrderMonthlyTotalRevenue(c *record.OrderMonthlyTotalRevenueRecord) *response.OrderMonthlyTotalRevenueResponse
	ToOrderMonthlyTotalRevenues(c []*record.OrderMonthlyTotalRevenueRecord) []*response.OrderMonthlyTotalRevenueResponse
	ToOrderYearlyTotalRevenue(c *record.OrderYearlyTotalRevenueRecord) *response.OrderYearlyTotalRevenueResponse
	ToOrderYearlyTotalRevenues(c []*record.OrderYearlyTotalRevenueRecord) []*response.OrderYearlyTotalRevenueResponse

	ToOrderMonthlyPrice(category *record.OrderMonthlyRecord) *response.OrderMonthlyResponse
	ToOrderMonthlyPrices(c []*record.OrderMonthlyRecord) []*response.OrderMonthlyResponse
	ToOrderYearlyPrice(category *record.OrderYearlyRecord) *response.OrderYearlyResponse
	ToOrderYearlyPrices(c []*record.OrderYearlyRecord) []*response.OrderYearlyResponse

	ToOrderResponse(order *record.OrderRecord) *response.OrderResponse
	ToOrdersResponse(orders []*record.OrderRecord) []*response.OrderResponse
	ToOrderResponseDeleteAt(order *record.OrderRecord) *response.OrderResponseDeleteAt
	ToOrdersResponseDeleteAt(orders []*record.OrderRecord) []*response.OrderResponseDeleteAt
}

type OrderItemResponseMapper interface {
	ToOrderItemResponse(order *record.OrderItemRecord) *response.OrderItemResponse
	ToOrderItemsResponse(orders []*record.OrderItemRecord) []*response.OrderItemResponse
	ToOrderItemResponseDeleteAt(order *record.OrderItemRecord) *response.OrderItemResponseDeleteAt
	ToOrderItemsResponseDeleteAt(orders []*record.OrderItemRecord) []*response.OrderItemResponseDeleteAt
}

type ProductResponseMapper interface {
	ToProductResponse(product *record.ProductRecord) *response.ProductResponse
	ToProductsResponse(products []*record.ProductRecord) []*response.ProductResponse
	ToProductResponseDeleteAt(product *record.ProductRecord) *response.ProductResponseDeleteAt
	ToProductsResponseDeleteAt(products []*record.ProductRecord) []*response.ProductResponseDeleteAt
}

type TransactionResponseMapper interface {
	ToTransactionMonthAmountSuccess(row *record.TransactionMonthlyAmountSuccessRecord) *response.TransactionMonthlyAmountSuccessResponse
	ToTransactionMonthlyAmountSuccess(rows []*record.TransactionMonthlyAmountSuccessRecord) []*response.TransactionMonthlyAmountSuccessResponse
	ToTransactionYearAmountSuccess(row *record.TransactionYearlyAmountSuccessRecord) *response.TransactionYearlyAmountSuccessResponse
	ToTransactionYearlyAmountSuccess(rows []*record.TransactionYearlyAmountSuccessRecord) []*response.TransactionYearlyAmountSuccessResponse
	ToTransactionMonthAmountFailed(row *record.TransactionMonthlyAmountFailedRecord) *response.TransactionMonthlyAmountFailedResponse
	ToTransactionMonthlyAmountFailed(rows []*record.TransactionMonthlyAmountFailedRecord) []*response.TransactionMonthlyAmountFailedResponse
	ToTransactionYearAmountFailed(row *record.TransactionYearlyAmountFailedRecord) *response.TransactionYearlyAmountFailedResponse
	ToTransactionYearlyAmountFailed(rows []*record.TransactionYearlyAmountFailedRecord) []*response.TransactionYearlyAmountFailedResponse
	ToTransactionMonthMethod(row *record.TransactionMonthlyMethodRecord) *response.TransactionMonthlyMethodResponse
	ToTransactionMonthlyMethod(rows []*record.TransactionMonthlyMethodRecord) []*response.TransactionMonthlyMethodResponse
	ToTransactionYearMethod(row *record.TransactionYearlyMethodRecord) *response.TransactionYearlyMethodResponse
	ToTransactionYearlyMethod(rows []*record.TransactionYearlyMethodRecord) []*response.TransactionYearlyMethodResponse

	ToTransactionResponse(transaction *record.TransactionRecord) *response.TransactionResponse
	ToTransactionsResponse(transactions []*record.TransactionRecord) []*response.TransactionResponse
	ToTransactionResponseDeleteAt(transaction *record.TransactionRecord) *response.TransactionResponseDeleteAt
	ToTransactionsResponseDeleteAt(transactions []*record.TransactionRecord) []*response.TransactionResponseDeleteAt
}

type CartResponseMapper interface {
	ToCartResponse(cart *record.CartRecord) *response.CartResponse
	ToCartsResponse(users []*record.CartRecord) []*response.CartResponse
}

type ReviewResponseMapper interface {
	ToReviewDetailResponse(review *record.ReviewsDetailRecord) *response.ReviewsDetailResponse
	ToReviewsDetailResponse(reviews []*record.ReviewsDetailRecord) []*response.ReviewsDetailResponse

	ToReviewResponse(review *record.ReviewRecord) *response.ReviewResponse
	ToReviewsResponse(reviews []*record.ReviewRecord) []*response.ReviewResponse
	ToReviewResponseDeleteAt(review *record.ReviewRecord) *response.ReviewResponseDeleteAt
	ToReviewsResponseDeleteAt(reviews []*record.ReviewRecord) []*response.ReviewResponseDeleteAt
}

type ReviewDetailResponeMapper interface {
	ToReviewDetailsResponse(record *record.ReviewDetailRecord) *response.ReviewDetailsResponse
	ToReviewsDetailsResponse(records []*record.ReviewDetailRecord) []*response.ReviewDetailsResponse
	ToReviewDetailResponseDeleteAt(record *record.ReviewDetailRecord) *response.ReviewDetailsResponseDeleteAt
	ToReviewDetailsResponseDeleteAt(merchants []*record.ReviewDetailRecord) []*response.ReviewDetailsResponseDeleteAt
}

type ShippingAddressResponseMapper interface {
	ToShippingAddressResponse(address *record.ShippingAddressRecord) *response.ShippingAddressResponse
	ToShippingAddressesResponse(addresses []*record.ShippingAddressRecord) []*response.ShippingAddressResponse
	ToShippingAddressResponseDeleteAt(address *record.ShippingAddressRecord) *response.ShippingAddressResponseDeleteAt
	ToShippingAddressesResponseDeleteAt(addresses []*record.ShippingAddressRecord) []*response.ShippingAddressResponseDeleteAt
}

type SliderResponseMapper interface {
	ToSliderResponse(slider *record.SliderRecord) *response.SliderResponse
	ToSlidersResponse(sliders []*record.SliderRecord) []*response.SliderResponse
	ToSliderResponseDeleteAt(slider *record.SliderRecord) *response.SliderResponseDeleteAt
	ToSlidersResponseDeleteAt(sliders []*record.SliderRecord) []*response.SliderResponseDeleteAt
}
