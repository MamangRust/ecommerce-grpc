package protomapper

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"
)

type AuthProtoMapper interface {
	ToProtoResponseLogin(status string, message string, response *response.TokenResponse) *pb.ApiResponseLogin
	ToProtoResponseRegister(status string, message string, response *response.UserResponse) *pb.ApiResponseRegister
	ToProtoResponseRefreshToken(status string, message string, response *response.TokenResponse) *pb.ApiResponseRefreshToken
	ToProtoResponseGetMe(status string, message string, response *response.UserResponse) *pb.ApiResponseGetMe
}

type UserProtoMapper interface {
	ToProtoResponseUserDeleteAt(status string, message string, pbResponse *response.UserResponseDeleteAt) *pb.ApiResponseUserDeleteAt
	ToProtoResponsesUser(status string, message string, pbResponse []*response.UserResponse) *pb.ApiResponsesUser
	ToProtoResponseUser(status string, message string, pbResponse *response.UserResponse) *pb.ApiResponseUser
	ToProtoResponseUserDelete(status string, message string) *pb.ApiResponseUserDelete
	ToProtoResponseUserAll(status string, message string) *pb.ApiResponseUserAll
	ToProtoResponsePaginationUserDeleteAt(pagination *pb.PaginationMeta, status string, message string, users []*response.UserResponseDeleteAt) *pb.ApiResponsePaginationUserDeleteAt
	ToProtoResponsePaginationUser(pagination *pb.PaginationMeta, status string, message string, users []*response.UserResponse) *pb.ApiResponsePaginationUser
}

type RoleProtoMapper interface {
	ToProtoResponseRoleAll(status string, message string) *pb.ApiResponseRoleAll
	ToProtoResponseRoleDelete(status string, message string) *pb.ApiResponseRoleDelete
	ToProtoResponseRole(status string, message string, pbResponse *response.RoleResponse) *pb.ApiResponseRole
	ToProtoResponsesRole(status string, message string, pbResponse []*response.RoleResponse) *pb.ApiResponsesRole
	ToProtoResponsePaginationRole(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.RoleResponse) *pb.ApiResponsePaginationRole
	ToProtoResponsePaginationRoleDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.RoleResponseDeleteAt) *pb.ApiResponsePaginationRoleDeleteAt
}

type BannerProtoMapper interface {
	ToProtoResponseBanner(status string, message string, pbResponse *response.BannerResponse) *pb.ApiResponseBanner
	ToProtoResponseBannerDeleteAt(status string, message string, pbResponse *response.BannerResponseDeleteAt) *pb.ApiResponseBannerDeleteAt
	ToProtoResponsesBanner(status string, message string, pbResponse []*response.BannerResponse) *pb.ApiResponsesBanner
	ToProtoResponseBannerDelete(status string, message string) *pb.ApiResponseBannerDelete
	ToProtoResponseBannerAll(status string, message string) *pb.ApiResponseBannerAll
	ToProtoResponsePaginationBannerDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.BannerResponseDeleteAt) *pb.ApiResponsePaginationBannerDeleteAt
	ToProtoResponsePaginationBanner(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.BannerResponse) *pb.ApiResponsePaginationBanner
}

type CategoryProtoMapper interface {
	ToProtoResponseMonthlyTotalPrice(status string, message string, row []*response.CategoriesMonthlyTotalPriceResponse) *pb.ApiResponseCategoryMonthlyTotalPrice
	ToProtoResponseYearlyTotalPrice(status string, message string, row []*response.CategoriesYearlyTotalPriceResponse) *pb.ApiResponseCategoryYearlyTotalPrice
	ToProtoResponseCategoryMonthlyPrice(status string, message string, row []*response.CategoryMonthPriceResponse) *pb.ApiResponseCategoryMonthPrice
	ToProtoResponseCategoryYearlyPrice(status string, message string, row []*response.CategoryYearPriceResponse) *pb.ApiResponseCategoryYearPrice

	ToProtoResponsesCategory(status string, message string, pbResponse []*response.CategoryResponse) *pb.ApiResponsesCategory
	ToProtoResponseCategoryDeleteAt(status string, message string, pbResponse *response.CategoryResponseDeleteAt) *pb.ApiResponseCategoryDeleteAt

	ToProtoResponseCategoryAll(status string, message string) *pb.ApiResponseCategoryAll
	ToProtoResponseCategory(status string, message string, pbResponse *response.CategoryResponse) *pb.ApiResponseCategory
	ToProtoResponseCategoryDelete(status string, message string) *pb.ApiResponseCategoryDelete
	ToProtoResponsePaginationCategoryDeleteAt(pagination *pb.PaginationMeta, status string, message string, categories []*response.CategoryResponseDeleteAt) *pb.ApiResponsePaginationCategoryDeleteAt
	ToProtoResponsePaginationCategory(pagination *pb.PaginationMeta, status string, message string, categories []*response.CategoryResponse) *pb.ApiResponsePaginationCategory
}

type MerchantProtoMapper interface {
	ToProtoResponseMerchant(status string, message string, pbResponse *response.MerchantResponse) *pb.ApiResponseMerchant
	ToProtoResponseMerchantDeleteAt(status string, message string, pbResponse *response.MerchantResponseDeleteAt) *pb.ApiResponseMerchantDeleteAt

	ToProtoResponsesMerchant(status string, message string, pbResponse []*response.MerchantResponse) *pb.ApiResponsesMerchant
	ToProtoResponseMerchantDelete(status string, message string) *pb.ApiResponseMerchantDelete
	ToProtoResponseMerchantAll(status string, message string) *pb.ApiResponseMerchantAll
	ToProtoResponsePaginationMerchantDeleteAt(pagination *pb.PaginationMeta, status string, message string, merchants []*response.MerchantResponseDeleteAt) *pb.ApiResponsePaginationMerchantDeleteAt
	ToProtoResponsePaginationMerchant(pagination *pb.PaginationMeta, status string, message string, merchants []*response.MerchantResponse) *pb.ApiResponsePaginationMerchant
}

type MerchantAwardProtoMapper interface {
	ToProtoResponseMerchantAward(status string, message string, pbResponse *response.MerchantAwardResponse) *pb.ApiResponseMerchantAward
	ToProtoResponseMerchantAwardDeleteAt(status string, message string, pbResponse *response.MerchantAwardResponseDeleteAt) *pb.ApiResponseMerchantAwardDeleteAt
	ToProtoResponsesMerchantAward(status string, message string, pbResponse []*response.MerchantAwardResponse) *pb.ApiResponsesMerchantAward
	ToProtoResponsePaginationMerchantAwardDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.MerchantAwardResponseDeleteAt) *pb.ApiResponsePaginationMerchantAwardDeleteAt
	ToProtoResponsePaginationMerchantAward(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.MerchantAwardResponse) *pb.ApiResponsePaginationMerchantAward
}

type MerchantBusinessProtoMapper interface {
	ToProtoResponseMerchantBusiness(status string, message string, pbResponse *response.MerchantBusinessResponse) *pb.ApiResponseMerchantBusiness
	ToProtoResponseMerchantBusinessDeleteAt(status string, message string, pbResponse *response.MerchantBusinessResponseDeleteAt) *pb.ApiResponseMerchantBusinessDeleteAt
	ToProtoResponsesMerchantBusiness(status string, message string, pbResponse []*response.MerchantBusinessResponse) *pb.ApiResponsesMerchantBusiness
	ToProtoResponsePaginationMerchantBusinessDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.MerchantBusinessResponseDeleteAt) *pb.ApiResponsePaginationMerchantBusinessDeleteAt
	ToProtoResponsePaginationMerchantBusiness(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.MerchantBusinessResponse) *pb.ApiResponsePaginationMerchantBusiness
}

type MerchantDetailProtoMapper interface {
	ToProtoResponseMerchantDetail(status string, message string, pbResponse *response.MerchantDetailResponse) *pb.ApiResponseMerchantDetail
	ToProtoResponseMerchantDetailRelation(status string, message string, pbResponse *response.MerchantDetailResponse) *pb.ApiResponseMerchantDetail
	ToProtoResponseMerchantDetailDeleteAt(status string, message string, pbResponse *response.MerchantDetailResponseDeleteAt) *pb.ApiResponseMerchantDetailDeleteAt
	ToProtoResponsesMerchantDetail(status string, message string, pbResponse []*response.MerchantDetailResponse) *pb.ApiResponsesMerchantDetail
	ToProtoResponsePaginationMerchantDetailDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.MerchantDetailResponseDeleteAt) *pb.ApiResponsePaginationMerchantDetailDeleteAt
	ToProtoResponsePaginationMerchantDetail(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.MerchantDetailResponse) *pb.ApiResponsePaginationMerchantDetail
}

type MerchantPolicyProtoMapper interface {
	ToProtoResponseMerchantPolicy(status string, message string, pbResponse *response.MerchantPoliciesResponse) *pb.ApiResponseMerchantPolicies
	ToProtoResponseMerchantPolicyDeleteAt(status string, message string, pbResponse *response.MerchantPoliciesResponseDeleteAt) *pb.ApiResponseMerchantPoliciesDeleteAt
	ToProtoResponsesMerchantPolicy(status string, message string, pbResponse []*response.MerchantPoliciesResponse) *pb.ApiResponsesMerchantPolicies
	ToProtoResponsePaginationMerchantPolicyDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.MerchantPoliciesResponseDeleteAt) *pb.ApiResponsePaginationMerchantPoliciesDeleteAt
	ToProtoResponsePaginationMerchantPolicy(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.MerchantPoliciesResponse) *pb.ApiResponsePaginationMerchantPolicies
}

type OrderItemProtoMapper interface {
	ToProtoResponseOrderItem(status string, message string, pbResponse *response.OrderItemResponse) *pb.ApiResponseOrderItem
	ToProtoResponsesOrderItem(status string, message string, pbResponse []*response.OrderItemResponse) *pb.ApiResponsesOrderItem
	ToProtoResponseOrderItemDelete(status string, message string) *pb.ApiResponseOrderItemDelete
	ToProtoResponseOrderItemAll(status string, message string) *pb.ApiResponseOrderItemAll
	ToProtoResponsePaginationOrderItemDeleteAt(pagination *pb.PaginationMeta, status string, message string, orderItems []*response.OrderItemResponseDeleteAt) *pb.ApiResponsePaginationOrderItemDeleteAt
	ToProtoResponsePaginationOrderItem(pagination *pb.PaginationMeta, status string, message string, orderItems []*response.OrderItemResponse) *pb.ApiResponsePaginationOrderItem
}

type OrderProtoMapper interface {
	ToProtoResponseMonthlyTotalRevenue(status string, message string, row []*response.OrderMonthlyTotalRevenueResponse) *pb.ApiResponseOrderMonthlyTotalRevenue
	ToProtoResponseYearlyTotalRevenue(status string, message string, row []*response.OrderYearlyTotalRevenueResponse) *pb.ApiResponseOrderYearlyTotalRevenue

	ToProtoResponseMonthlyRevenue(status string, message string, row []*response.OrderMonthlyResponse) *pb.ApiResponseOrderMonthly
	ToProtoResponseYearlyRevenue(status string, message string, row []*response.OrderYearlyResponse) *pb.ApiResponseOrderYearly

	ToProtoResponseOrder(status string, message string, pbResponse *response.OrderResponse) *pb.ApiResponseOrder
	ToProtoResponseOrderDeleteAt(status string, message string, pbResponse *response.OrderResponseDeleteAt) *pb.ApiResponseOrderDeleteAt
	ToProtoResponsesOrder(status string, message string, pbResponse []*response.OrderResponse) *pb.ApiResponsesOrder
	ToProtoResponseOrderDelete(status string, message string) *pb.ApiResponseOrderDelete
	ToProtoResponseOrderAll(status string, message string) *pb.ApiResponseOrderAll
	ToProtoResponsePaginationOrderDeleteAt(pagination *pb.PaginationMeta, status string, message string, orders []*response.OrderResponseDeleteAt) *pb.ApiResponsePaginationOrderDeleteAt
	ToProtoResponsePaginationOrder(pagination *pb.PaginationMeta, status string, message string, orders []*response.OrderResponse) *pb.ApiResponsePaginationOrder
}

type ProductProtoMapper interface {
	ToProtoResponseProduct(status string, message string, pbResponse *response.ProductResponse) *pb.ApiResponseProduct
	ToProtoResponseProductDeleteAt(status string, message string, pbResponse *response.ProductResponseDeleteAt) *pb.ApiResponseProductDeleteAt

	ToProtoResponsesProduct(status string, message string, pbResponse []*response.ProductResponse) *pb.ApiResponsesProduct
	ToProtoResponseProductDelete(status string, message string) *pb.ApiResponseProductDelete
	ToProtoResponseProductAll(status string, message string) *pb.ApiResponseProductAll
	ToProtoResponsePaginationProductDeleteAt(pagination *pb.PaginationMeta, status string, message string, products []*response.ProductResponseDeleteAt) *pb.ApiResponsePaginationProductDeleteAt
	ToProtoResponsePaginationProduct(pagination *pb.PaginationMeta, status string, message string, products []*response.ProductResponse) *pb.ApiResponsePaginationProduct
}

type TransactionProtoMapper interface {
	ToProtoResponseMonthAmountSuccess(status string, message string, row []*response.TransactionMonthlyAmountSuccessResponse) *pb.ApiResponseTransactionMonthAmountSuccess
	ToProtoResponseYearAmountSuccess(status string, message string, row []*response.TransactionYearlyAmountSuccessResponse) *pb.ApiResponseTransactionYearAmountSuccess
	ToProtoResponseMonthAmountFailed(status string, message string, row []*response.TransactionMonthlyAmountFailedResponse) *pb.ApiResponseTransactionMonthAmountFailed
	ToProtoResponseYearAmountFailed(status string, message string, row []*response.TransactionYearlyAmountFailedResponse) *pb.ApiResponseTransactionYearAmountFailed
	ToProtoResponseMonthMethod(status string, message string, row []*response.TransactionMonthlyMethodResponse) *pb.ApiResponseTransactionMonthPaymentMethod
	ToProtoResponseYearMethod(status string, message string, row []*response.TransactionYearlyMethodResponse) *pb.ApiResponseTransactionYearPaymentmethod

	ToProtoResponseTransaction(status string, message string, trans *response.TransactionResponse) *pb.ApiResponseTransaction
	ToProtoResponseTransactionDeleteAt(status string, message string, trans *response.TransactionResponseDeleteAt) *pb.ApiResponseTransactionDeleteAt
	ToProtoResponsesTransaction(status string, message string, transList []*response.TransactionResponse) *pb.ApiResponsesTransaction
	ToProtoResponseTransactionDelete(status string, message string) *pb.ApiResponseTransactionDelete
	ToProtoResponseTransactionAll(status string, message string) *pb.ApiResponseTransactionAll
	ToProtoResponsePaginationTransactionDeleteAt(pagination *pb.PaginationMeta, status string, message string, transactions []*response.TransactionResponseDeleteAt) *pb.ApiResponsePaginationTransactionDeleteAt
	ToProtoResponsePaginationTransaction(pagination *pb.PaginationMeta, status string, message string, transactions []*response.TransactionResponse) *pb.ApiResponsePaginationTransaction
}

type CartProtoMapper interface {
	ToProtoResponseCartDelete(status string, message string) *pb.ApiResponseCartDelete
	ToProtoResponseCartAll(status string, message string) *pb.ApiResponseCartAll
	ToProtoResponseCart(status string, message string, categories *response.CartResponse) *pb.ApiResponseCart
	ToProtoResponsePaginationCart(pagination *pb.PaginationMeta, status string, message string, categories []*response.CartResponse) *pb.ApiResponsePaginationCart
}

type ReviewProtoMapper interface {
	ToProtoResponseReview(status string, message string, pbResponse *response.ReviewResponse) *pb.ApiResponseReview
	ToProtoResponseReviewDeleteAt(status string, message string, pbResponse *response.ReviewResponseDeleteAt) *pb.ApiResponseReviewDeleteAt
	ToProtoResponsesReview(status string, message string, pbResponse []*response.ReviewResponse) *pb.ApiResponsesReview
	ToProtoResponseReviewDelete(status string, message string) *pb.ApiResponseReviewDelete
	ToProtoResponseReviewAll(status string, message string) *pb.ApiResponseReviewAll
	ToProtoResponsePaginationReviewDeleteAt(pagination *pb.PaginationMeta, status string, message string, reviews []*response.ReviewResponseDeleteAt) *pb.ApiResponsePaginationReviewDeleteAt
	ToProtoResponsePaginationReview(pagination *pb.PaginationMeta, status string, message string, reviews []*response.ReviewResponse) *pb.ApiResponsePaginationReview

	ToProtoResponsePaginationReviewsDetail(pagination *pb.PaginationMeta, status string, message string, reviews []*response.ReviewsDetailResponse) *pb.ApiResponsePaginationReviewDetail
}

type ReviewDetailProtoMapper interface {
	ToProtoResponseReviewDetail(status string, message string, pbResponse *response.ReviewDetailsResponse) *pb.ApiResponseReviewDetail
	ToProtoResponseReviewDetailDeleteAt(status string, message string, pbResponse *response.ReviewDetailsResponseDeleteAt) *pb.ApiResponseReviewDetailDeleteAt
	ToProtoResponsesReviewDetail(status string, message string, pbResponse []*response.ReviewDetailsResponse) *pb.ApiResponsesReviewDetails
	ToProtoResponsePaginationReviewDetailDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.ReviewDetailsResponseDeleteAt) *pb.ApiResponsePaginationReviewDetailsDeleteAt
	ToProtoResponsePaginationReviewDetail(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.ReviewDetailsResponse) *pb.ApiResponsePaginationReviewDetails
}

type ShippingAddresProtoMapper interface {
	ToProtoResponseShippingAddress(status string, message string, pbResponse *response.ShippingAddressResponse) *pb.ApiResponseShipping
	ToProtoResponseShippingAddressDeleteAt(status string, message string, pbResponse *response.ShippingAddressResponseDeleteAt) *pb.ApiResponseShippingDeleteAt
	ToProtoResponsesShippingAddress(status string, message string, pbResponse []*response.ShippingAddressResponse) *pb.ApiResponsesShipping
	ToProtoResponseShippingAddressDelete(status string, message string) *pb.ApiResponseShippingDelete
	ToProtoResponseShippingAddressAll(status string, message string) *pb.ApiResponseShippingAll
	ToProtoResponsePaginationShippingAddressDeleteAt(pagination *pb.PaginationMeta, status string, message string, addresses []*response.ShippingAddressResponseDeleteAt) *pb.ApiResponsePaginationShippingDeleteAt
	ToProtoResponsePaginationShippingAddress(pagination *pb.PaginationMeta, status string, message string, addresses []*response.ShippingAddressResponse) *pb.ApiResponsePaginationShipping
}

type SliderProtoMapper interface {
	ToProtoResponseSlider(status string, message string, pbResponse *response.SliderResponse) *pb.ApiResponseSlider
	ToProtoResponseSliderDeleteAt(status string, message string, pbResponse *response.SliderResponseDeleteAt) *pb.ApiResponseSliderDeleteAt
	ToProtoResponsesSlider(status string, message string, pbResponse []*response.SliderResponse) *pb.ApiResponsesSlider
	ToProtoResponseSliderDelete(status string, message string) *pb.ApiResponseSliderDelete
	ToProtoResponseSliderAll(status string, message string) *pb.ApiResponseSliderAll
	ToProtoResponsePaginationSliderDeleteAt(pagination *pb.PaginationMeta, status string, message string, sliders []*response.SliderResponseDeleteAt) *pb.ApiResponsePaginationSliderDeleteAt
	ToProtoResponsePaginationSlider(pagination *pb.PaginationMeta, status string, message string, sliders []*response.SliderResponse) *pb.ApiResponsePaginationSlider
}
