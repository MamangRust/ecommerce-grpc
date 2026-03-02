package gapi

import (
	"ecommerce/internal/service"
)

type Handler struct {
	Auth             AuthHandleGrpc
	Role             RoleHandleGrpc
	User             UserHandleGrpc
	Category         CategoryHandleGrpc
	Merchant         MerchantHandleGrpc
	OrderItem        OrderItemHandleGrpc
	Order            OrderHandleGrpc
	Product          ProductHandleGrpc
	Transaction      TransactionHandleGrpc
	Cart             CartHandleGrpc
	Review           ReviewHandleGrpc
	Shipping         ShippingAdddressHandleGrpc
	Slider           SliderHandleGrpc
	Banner           BannerHandleGrpc
	MerchantAward    MerchantAwardHandleGrpc
	MerchantBusiness MerchantBusinessHandleGrpc
	MerchantDetail   MerchantDetailHandleGrpc
	MerchantPolicies MerchantPoliciesHandleGrpc
	ReviewDetail     ReviewDetailHandleGrpc
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		Auth:             NewAuthHandleGrpc(service.Auth),
		Role:             NewRoleHandleGrpc(service.Role),
		User:             NewUserHandleGrpc(service.User),
		Category:         NewCategoryHandleGrpc(service.Category),
		Merchant:         NewMerchantHandleGrpc(service.Merchant),
		OrderItem:        NewOrderItemHandleGrpc(service.OrderItem),
		Order:            NewOrderHandleGrpc(service.Order),
		Product:          NewProductHandleGrpc(service.Product),
		Transaction:      NewTransactionHandleGrpc(service.Transaction),
		Review:           NewReviewHandleGrpc(service.Review),
		Shipping:         NewShippingAddressHandleGrpc(service.Shipping),
		Slider:           NewSliderHandleGrpc(service.Slider),
		Cart:             NewCartHandleGrpc(service.Cart),
		Banner:           NewBannerHandleGrpc(service.Banner),
		MerchantAward:    NewMerchantAwardHandleGrpc(service.MerchantAward),
		MerchantBusiness: NewMerchantBusinessHandleGrpc(service.MerchantBusiness),
		MerchantDetail:   NewMerchantDetailHandleGrpc(service.MerchantDetail),
		MerchantPolicies: NewMerchantPolicyHandleGrpc(service.MerchantPolicies),
		ReviewDetail:     NewReviewDetailHandleGrpc(service.ReviewDetail),
	}
}
