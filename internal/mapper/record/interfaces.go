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

type UserRoleRecordMapping interface {
	ToUserRoleRecord(userRole *db.UserRole) *record.UserRoleRecord
}

type RefreshTokenRecordMapping interface {
	ToRefreshTokenRecord(refreshToken *db.RefreshToken) *record.RefreshTokenRecord
	ToRefreshTokensRecord(refreshTokens []*db.RefreshToken) []*record.RefreshTokenRecord
}

type CategoryRecordMapper interface {
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

	ToReviewProductRecordPagination(review *db.GetReviewsByProductIDRow) *record.ReviewRecord
	ToReviewsProductRecordPagination(reviews []*db.GetReviewsByProductIDRow) []*record.ReviewRecord

	ToReviewRecordActivePagination(review *db.GetReviewsActiveRow) *record.ReviewRecord
	ToReviewsRecordActivePagination(reviews []*db.GetReviewsActiveRow) []*record.ReviewRecord
	ToReviewRecordTrashedPagination(review *db.GetReviewsTrashedRow) *record.ReviewRecord
	ToReviewsRecordTrashedPagination(reviews []*db.GetReviewsTrashedRow) []*record.ReviewRecord
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
