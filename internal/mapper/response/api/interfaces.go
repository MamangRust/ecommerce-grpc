package response_api

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"
)

type AuthResponseMapper interface {
	ToResponseLogin(res *pb.ApiResponseLogin) *response.ApiResponseLogin
	ToResponseRegister(res *pb.ApiResponseRegister) *response.ApiResponseRegister
	ToResponseRefreshToken(res *pb.ApiResponseRefreshToken) *response.ApiResponseRefreshToken
	ToResponseGetMe(res *pb.ApiResponseGetMe) *response.ApiResponseGetMe
}

type RoleResponseMapper interface {
	ToApiResponseRoleAll(pbResponse *pb.ApiResponseRoleAll) *response.ApiResponseRoleAll
	ToApiResponseRoleDelete(pbResponse *pb.ApiResponseRoleDelete) *response.ApiResponseRoleDelete
	ToApiResponseRole(pbResponse *pb.ApiResponseRole) *response.ApiResponseRole
	ToApiResponsesRole(pbResponse *pb.ApiResponsesRole) *response.ApiResponsesRole
	ToApiResponsePaginationRole(pbResponse *pb.ApiResponsePaginationRole) *response.ApiResponsePaginationRole
	ToApiResponsePaginationRoleDeleteAt(pbResponse *pb.ApiResponsePaginationRoleDeleteAt) *response.ApiResponsePaginationRoleDeleteAt
}

type UserResponseMapper interface {
	ToApiResponseUserDeleteAt(pbResponse *pb.ApiResponseUserDeleteAt) *response.ApiResponseUserDeleteAt
	ToApiResponseUser(pbResponse *pb.ApiResponseUser) *response.ApiResponseUser
	ToApiResponsesUser(pbResponse *pb.ApiResponsesUser) *response.ApiResponsesUser

	ToApiResponseUserDelete(pbResponse *pb.ApiResponseUserDelete) *response.ApiResponseUserDelete
	ToApiResponseUserAll(pbResponse *pb.ApiResponseUserAll) *response.ApiResponseUserAll
	ToApiResponsePaginationUserDeleteAt(pbResponse *pb.ApiResponsePaginationUserDeleteAt) *response.ApiResponsePaginationUserDeleteAt
	ToApiResponsePaginationUser(pbResponse *pb.ApiResponsePaginationUser) *response.ApiResponsePaginationUser
}

type CategoryResponseMapper interface {
	ToApiResponseCategory(pbResponse *pb.ApiResponseCategory) *response.ApiResponseCategory
	ToApiResponseCategoryDeleteAt(pbResponse *pb.ApiResponseCategoryDeleteAt) *response.ApiResponseCategoryDeleteAt
	ToApiResponsesCategory(pbResponse *pb.ApiResponsesCategory) *response.ApiResponsesCategory
	ToApiResponseCategoryDelete(pbResponse *pb.ApiResponseCategoryDelete) *response.ApiResponseCategoryDelete
	ToApiResponseCategoryAll(pbResponse *pb.ApiResponseCategoryAll) *response.ApiResponseCategoryAll
	ToApiResponsePaginationCategoryDeleteAt(pbResponse *pb.ApiResponsePaginationCategoryDeleteAt) *response.ApiResponsePaginationCategoryDeleteAt
	ToApiResponsePaginationCategory(pbResponse *pb.ApiResponsePaginationCategory) *response.ApiResponsePaginationCategory
}

type MerchantResponseMapper interface {
	ToApiResponseMerchant(pbResponse *pb.ApiResponseMerchant) *response.ApiResponseMerchant

	ToApiResponseMerchantDeleteAt(pbResponse *pb.ApiResponseMerchantDeleteAt) *response.ApiResponseMerchantDeleteAt
	ToApiResponsesMerchant(pbResponse *pb.ApiResponsesMerchant) *response.ApiResponsesMerchant
	ToApiResponseMerchantDelete(pbResponse *pb.ApiResponseMerchantDelete) *response.ApiResponseMerchantDelete
	ToApiResponseMerchantAll(pbResponse *pb.ApiResponseMerchantAll) *response.ApiResponseMerchantAll
	ToApiResponsePaginationMerchantDeleteAt(pbResponse *pb.ApiResponsePaginationMerchantDeleteAt) *response.ApiResponsePaginationMerchantDeleteAt
	ToApiResponsePaginationMerchant(pbResponse *pb.ApiResponsePaginationMerchant) *response.ApiResponsePaginationMerchant
}

type OrderItemResponseMapper interface {
	ToApiResponseOrderItem(pbResponse *pb.ApiResponseOrderItem) *response.ApiResponseOrderItem
	ToApiResponsesOrderItem(pbResponse *pb.ApiResponsesOrderItem) *response.ApiResponsesOrderItem
	ToApiResponseOrderItemDelete(pbResponse *pb.ApiResponseOrderItemDelete) *response.ApiResponseOrderItemDelete
	ToApiResponseOrderItemAll(pbResponse *pb.ApiResponseOrderItemAll) *response.ApiResponseOrderItemAll
	ToApiResponsePaginationOrderItemDeleteAt(pbResponse *pb.ApiResponsePaginationOrderItemDeleteAt) *response.ApiResponsePaginationOrderItemDeleteAt
	ToApiResponsePaginationOrderItem(pbResponse *pb.ApiResponsePaginationOrderItem) *response.ApiResponsePaginationOrderItem
}

type OrderResponseMapper interface {
	ToApiResponseOrder(pbResponse *pb.ApiResponseOrder) *response.ApiResponseOrder
	ToApiResponseOrderDeleteAt(pbResponse *pb.ApiResponseOrderDeleteAt) *response.ApiResponseOrderDeleteAt
	ToApiResponsesOrder(pbResponse *pb.ApiResponsesOrder) *response.ApiResponsesOrder
	ToApiResponseOrderDelete(pbResponse *pb.ApiResponseOrderDelete) *response.ApiResponseOrderDelete
	ToApiResponseOrderAll(pbResponse *pb.ApiResponseOrderAll) *response.ApiResponseOrderAll
	ToApiResponsePaginationOrderDeleteAt(pbResponse *pb.ApiResponsePaginationOrderDeleteAt) *response.ApiResponsePaginationOrderDeleteAt
	ToApiResponsePaginationOrder(pbResponse *pb.ApiResponsePaginationOrder) *response.ApiResponsePaginationOrder
}

type ProductResponseMapper interface {
	ToApiResponseProduct(pbResponse *pb.ApiResponseProduct) *response.ApiResponseProduct
	ToApiResponsesProductDeleteAt(pbResponse *pb.ApiResponseProductDeleteAt) *response.ApiResponseProductDeleteAt
	ToApiResponsesProduct(pbResponse *pb.ApiResponsesProduct) *response.ApiResponsesProduct
	ToApiResponseProductDelete(pbResponse *pb.ApiResponseProductDelete) *response.ApiResponseProductDelete
	ToApiResponseProductAll(pbResponse *pb.ApiResponseProductAll) *response.ApiResponseProductAll
	ToApiResponsePaginationProductDeleteAt(pbResponse *pb.ApiResponsePaginationProductDeleteAt) *response.ApiResponsePaginationProductDeleteAt
	ToApiResponsePaginationProduct(pbResponse *pb.ApiResponsePaginationProduct) *response.ApiResponsePaginationProduct
}

type TransactionResponseMapper interface {
	ToApiResponseTransaction(pbResponse *pb.ApiResponseTransaction) *response.ApiResponseTransaction
	ToApiResponseTransactionDeleteAt(pbResponse *pb.ApiResponseTransactionDeleteAt) *response.ApiResponseTransactionDeleteAt
	ToApiResponsesTransaction(pbResponse *pb.ApiResponsesTransaction) *response.ApiResponsesTransaction
	ToApiResponseTransactionDelete(pbResponse *pb.ApiResponseTransactionDelete) *response.ApiResponseTransactionDelete
	ToApiResponseTransactionAll(pbResponse *pb.ApiResponseTransactionAll) *response.ApiResponseTransactionAll
	ToApiResponsePaginationTransactionDeleteAt(pbResponse *pb.ApiResponsePaginationTransactionDeleteAt) *response.ApiResponsePaginationTransactionDeleteAt
	ToApiResponsePaginationTransaction(pbResponse *pb.ApiResponsePaginationTransaction) *response.ApiResponsePaginationTransaction
}

type CartResponseMapper interface {
	ToApiResponseCartPagination(pbResponse *pb.ApiResponsePaginationCart) *response.ApiResponseCartPagination
	ToApiResponseCartDelete(pbResponse *pb.ApiResponseCartDelete) *response.ApiResponseCartDelete
	ToApiResponseCartAll(pbResponse *pb.ApiResponseCartAll) *response.ApiResponseCartAll
}

type ReviewResponseMapper interface {
	ToApiResponseReview(pbResponse *pb.ApiResponseReview) *response.ApiResponseReview
	ToApiResponseReviewDeleteAt(pbResponse *pb.ApiResponseReviewDeleteAt) *response.ApiResponseReviewDeleteAt
	ToApiResponsesReview(pbResponse *pb.ApiResponsesReview) *response.ApiResponsesReview
	ToApiResponseReviewDelete(pbResponse *pb.ApiResponseReviewDelete) *response.ApiResponseReviewDelete
	ToApiResponseReviewAll(pbResponse *pb.ApiResponseReviewAll) *response.ApiResponseReviewAll
	ToApiResponsePaginationReviewDeleteAt(pbResponse *pb.ApiResponsePaginationReviewDeleteAt) *response.ApiResponsePaginationReviewDeleteAt
	ToApiResponsePaginationReview(pbResponse *pb.ApiResponsePaginationReview) *response.ApiResponsePaginationReview
}

type ShippingAddressResponseMapper interface {
	ToApiResponseShippingAddress(pbResponse *pb.ApiResponseShipping) *response.ApiResponseShippingAddress
	ToApiResponseShippingAddressDeleteAt(pbResponse *pb.ApiResponseShippingDeleteAt) *response.ApiResponseShippingAddressDeleteAt
	ToApiResponsesShippingAddress(pbResponse *pb.ApiResponsesShipping) *response.ApiResponsesShippingAddress
	ToApiResponseShippingAddressDelete(pbResponse *pb.ApiResponseShippingDelete) *response.ApiResponseShippingAddressDelete
	ToApiResponseShippingAddressAll(pbResponse *pb.ApiResponseShippingAll) *response.ApiResponseShippingAddressAll
	ToApiResponsePaginationShippingAddressDeleteAt(pbResponse *pb.ApiResponsePaginationShippingDeleteAt) *response.ApiResponsePaginationShippingAddressDeleteAt
	ToApiResponsePaginationShippingAddress(pbResponse *pb.ApiResponsePaginationShipping) *response.ApiResponsePaginationShippingAddress
}

type SliderResponseMapper interface {
	ToApiResponseSlider(pbResponse *pb.ApiResponseSlider) *response.ApiResponseSlider
	ToApiResponseSliderAll(pbResponse *pb.ApiResponseSliderAll) *response.ApiResponseSliderAll
	ToApiResponseSliderDeleteAt(pbResponse *pb.ApiResponseSliderDeleteAt) *response.ApiResponseSliderDeleteAt
	ToApiResponsesSlider(pbResponse *pb.ApiResponsesSlider) *response.ApiResponsesSlider
	ToApiResponseSliderDelete(pbResponse *pb.ApiResponseSliderDelete) *response.ApiResponseSliderDelete
	ToApiResponsePaginationSliderDeleteAt(pbResponse *pb.ApiResponsePaginationSliderDeleteAt) *response.ApiResponsePaginationSliderDeleteAt
	ToApiResponsePaginationSlider(pbResponse *pb.ApiResponsePaginationSlider) *response.ApiResponsePaginationSlider
}
