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

type CategoryResponseMapper interface {
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

type OrderResponseMapper interface {
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
	ToReviewResponse(review *record.ReviewRecord) *response.ReviewResponse
	ToReviewsResponse(reviews []*record.ReviewRecord) []*response.ReviewResponse
	ToReviewResponseDeleteAt(review *record.ReviewRecord) *response.ReviewResponseDeleteAt
	ToReviewsResponseDeleteAt(reviews []*record.ReviewRecord) []*response.ReviewResponseDeleteAt
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
