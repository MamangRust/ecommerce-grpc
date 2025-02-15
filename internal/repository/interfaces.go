package repository

import (
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
)

type UserRepository interface {
	FindAllUsers(search string, page, pageSize int) ([]*record.UserRecord, int, error)
	FindById(user_id int) (*record.UserRecord, error)
	FindByEmail(email string) (*record.UserRecord, error)
	FindByActive(search string, page, pageSize int) ([]*record.UserRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.UserRecord, int, error)
	CreateUser(request *requests.CreateUserRequest) (*record.UserRecord, error)
	UpdateUser(request *requests.UpdateUserRequest) (*record.UserRecord, error)
	TrashedUser(user_id int) (*record.UserRecord, error)
	RestoreUser(user_id int) (*record.UserRecord, error)
	DeleteUserPermanent(user_id int) (bool, error)
	RestoreAllUser() (bool, error)
	DeleteAllUserPermanent() (bool, error)
}

type RoleRepository interface {
	FindAllRoles(page int, pageSize int, search string) ([]*record.RoleRecord, int, error)
	FindById(role_id int) (*record.RoleRecord, error)
	FindByName(name string) (*record.RoleRecord, error)
	FindByUserId(user_id int) ([]*record.RoleRecord, error)
	FindByActiveRole(page int, pageSize int, search string) ([]*record.RoleRecord, int, error)
	FindByTrashedRole(page int, pageSize int, search string) ([]*record.RoleRecord, int, error)
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

type CategoryRepository interface {
	FindAllCategory(search string, page, pageSize int) ([]*record.CategoriesRecord, int, error)
	FindById(category_id int) (*record.CategoriesRecord, error)
	FindByActive(search string, page, pageSize int) ([]*record.CategoriesRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.CategoriesRecord, int, error)
	CreateCategory(request *requests.CreateCategoryRequest) (*record.CategoriesRecord, error)
	UpdateCategory(request *requests.UpdateCategoryRequest) (*record.CategoriesRecord, error)
	TrashedCategory(category_id int) (*record.CategoriesRecord, error)
	RestoreCategory(category_id int) (*record.CategoriesRecord, error)
	DeleteCategoryPermanently(Category_id int) (bool, error)
	RestoreAllCategories() (bool, error)
	DeleteAllPermanentCategories() (bool, error)
}

type MerchantRepository interface {
	FindAllMerchants(search string, page, pageSize int) ([]*record.MerchantRecord, int, error)
	FindByActive(search string, page, pageSize int) ([]*record.MerchantRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.MerchantRecord, int, error)
	FindById(user_id int) (*record.MerchantRecord, error)
	CreateMerchant(request *requests.CreateMerchantRequest) (*record.MerchantRecord, error)
	UpdateMerchant(request *requests.UpdateMerchantRequest) (*record.MerchantRecord, error)
	TrashedMerchant(merchant_id int) (*record.MerchantRecord, error)
	RestoreMerchant(merchant_id int) (*record.MerchantRecord, error)
	DeleteMerchantPermanent(Merchant_id int) (bool, error)
	RestoreAllMerchant() (bool, error)
	DeleteAllMerchantPermanent() (bool, error)
}

type OrderRepository interface {
	FindAllOrders(search string, page, pageSize int) ([]*record.OrderRecord, int, error)
	FindByActive(search string, page, pageSize int) ([]*record.OrderRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.OrderRecord, int, error)
	FindByMerchant(merchant_id int, search string, page, pageSize int) ([]*record.OrderRecord, int, error)
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
	FindAllOrderItems(search string, page, pageSize int) ([]*record.OrderItemRecord, int, error)
	FindByActive(search string, page, pageSize int) ([]*record.OrderItemRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.OrderItemRecord, int, error)
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
	FindAllProducts(search string, page, pageSize int) ([]*record.ProductRecord, int, error)
	FindByActive(search string, page, pageSize int) ([]*record.ProductRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.ProductRecord, int, error)
	FindByMerchant(merchant_id int, search string, page, pageSize int) ([]*record.ProductRecord, int, error)

	FindByCategory(category_name string, search string, page, pageSize int) ([]*record.ProductRecord, int, error)

	FindById(user_id int) (*record.ProductRecord, error)
	CreateProduct(request *requests.CreateProductRequest) (*record.ProductRecord, error)
	UpdateProduct(request *requests.UpdateProductRequest) (*record.ProductRecord, error)
	UpdateProductCountStock(product_id int, stock int) (*record.ProductRecord, error)
	TrashedProduct(user_id int) (*record.ProductRecord, error)
	RestoreProduct(user_id int) (*record.ProductRecord, error)
	DeleteProductPermanent(user_id int) (bool, error)
	RestoreAllProducts() (bool, error)
	DeleteAllProductPermanent() (bool, error)
}

type TransactionRepository interface {
	FindAllTransactions(search string, page, pageSize int) ([]*record.TransactionRecord, int, error)
	FindByActive(search string, page, pageSize int) ([]*record.TransactionRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.TransactionRecord, int, error)
	FindByMerchant(merchant_id int, search string, page, pageSize int) ([]*record.TransactionRecord, int, error)
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
	FindCarts(user_id int, search string, page, pageSize int) ([]*record.CartRecord, int, error)
	CreateCart(req *requests.CartCreateRecord) (*record.CartRecord, error)
	DeletePermanent(cart_id int) (bool, error)
	DeleteAllPermanently(req *requests.DeleteCartRequest) (bool, error)
}

type ShippingAddressRepository interface {
	FindAllShippingAddress(search string, page, pageSize int) ([]*record.ShippingAddressRecord, int, error)
	FindByOrder(shipping_id int) (*record.ShippingAddressRecord, error)
	FindById(shipping_id int) (*record.ShippingAddressRecord, error)
	FindByActive(search string, page, pageSize int) ([]*record.ShippingAddressRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.ShippingAddressRecord, int, error)
	CreateShippingAddress(request *requests.CreateShippingAddressRequest) (*record.ShippingAddressRecord, error)
	UpdateShippingAddress(request *requests.UpdateShippingAddressRequest) (*record.ShippingAddressRecord, error)
	TrashShippingAddress(category_id int) (*record.ShippingAddressRecord, error)
	RestoreShippingAddress(category_id int) (*record.ShippingAddressRecord, error)
	DeleteShippingAddressPermanently(category_id int) (bool, error)
	RestoreAllShippingAddress() (bool, error)
	DeleteAllPermanentShippingAddress() (bool, error)
}

type SliderRepository interface {
	FindAllSlider(search string, page, pageSize int) ([]*record.SliderRecord, int, error)
	FindByActive(search string, page, pageSize int) ([]*record.SliderRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.SliderRecord, int, error)
	CreateSlider(request *requests.CreateSliderRequest) (*record.SliderRecord, error)
	UpdateSlider(request *requests.UpdateSliderRequest) (*record.SliderRecord, error)
	TrashSlider(slider_id int) (*record.SliderRecord, error)
	RestoreSlider(slider_id int) (*record.SliderRecord, error)
	DeleteSliderPermanently(slider_id int) (bool, error)
	RestoreAllSlider() (bool, error)
	DeleteAllPermanentSlider() (bool, error)
}

type ReviewRepository interface {
	FindAllReview(search string, page, pageSize int) ([]*record.ReviewRecord, int, error)
	FindByProduct(product_id int, search string, page, pageSize int) ([]*record.ReviewRecord, int, error)
	FindByActive(search string, page, pageSize int) ([]*record.ReviewRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.ReviewRecord, int, error)
	FindById(id int) (*record.ReviewRecord, error)
	CreateReview(request *requests.CreateReviewRequest) (*record.ReviewRecord, error)
	UpdateReview(request *requests.UpdateReviewRequest) (*record.ReviewRecord, error)
	TrashReview(shipping_id int) (*record.ReviewRecord, error)
	RestoreReview(category_id int) (*record.ReviewRecord, error)
	DeleteReviewPermanently(category_id int) (bool, error)
	RestoreAllReview() (bool, error)
	DeleteAllPermanentReview() (bool, error)
}
