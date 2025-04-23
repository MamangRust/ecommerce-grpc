package gapi

import (
	protomapper "ecommerce/internal/mapper/proto"
	"ecommerce/internal/service"
)

type Deps struct {
	Service service.Service
	Mapper  protomapper.ProtoMapper
}

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

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Auth:             NewAuthHandleGrpc(deps.Service.Auth, deps.Mapper.AuthProtoMapper),
		Role:             NewRoleHandleGrpc(deps.Service.Role, deps.Mapper.RoleProtoMapper),
		User:             NewUserHandleGrpc(deps.Service.User, deps.Mapper.UserProtoMapper),
		Category:         NewCategoryHandleGrpc(deps.Service.Category, deps.Mapper.CategoryProtoMapper),
		Merchant:         NewMerchantHandleGrpc(deps.Service.Merchant, deps.Mapper.MerchantProtoMapper),
		OrderItem:        NewOrderItemHandleGrpc(deps.Service.OrderItem, deps.Mapper.OrderItemProtoMapper),
		Order:            NewOrderHandleGrpc(deps.Service.Order, deps.Mapper.OrderProtoMapper),
		Product:          NewProductHandleGrpc(deps.Service.Product, deps.Mapper.ProductProtoMapper),
		Transaction:      NewTransactionHandleGrpc(deps.Service.Transaction, deps.Mapper.TransactionProtoMapper),
		Review:           NewReviewHandleGrpc(deps.Service.Review, deps.Mapper.ReviewProtoMapper),
		Shipping:         NewShippingAddressHandleGrpc(deps.Service.Shipping, deps.Mapper.ShippingProtoMapper),
		Slider:           NewSliderHandleGrpc(deps.Service.Slider, deps.Mapper.SliderProtoMapper),
		Cart:             NewCartHandleGrpc(deps.Service.Cart, deps.Mapper.CartProtoMapper),
		Banner:           NewBannerHaandleGrpc(deps.Service.Banner, deps.Mapper.BannerProtoMapper),
		MerchantAward:    NewMerchantAwardHandleGrpc(deps.Service.MerchantAward, deps.Mapper.MerchantAwardProtoMapper, deps.Mapper.MerchantProtoMapper),
		MerchantBusiness: NewMerchantBusinessHandleGrpc(deps.Service.MerchantBusiness, deps.Mapper.MerchantBusinessProtoMapper, deps.Mapper.MerchantProtoMapper),
		MerchantDetail:   NewMerchantDetailHandleGrpc(deps.Service.MerchantDetail, deps.Mapper.MerchantDetailProtoMapper, deps.Mapper.MerchantProtoMapper),
		MerchantPolicies: NewMerchantPolicyHandleGrpc(deps.Service.MerchantPolicies, deps.Mapper.MerchantPolicyProtoMapper, deps.Mapper.MerchantProtoMapper),
		ReviewDetail:     NewReviewDetailHandleGrpc(deps.Service.ReviewDetail, deps.Mapper.ReviewDetailProtoMapper, deps.Mapper.ReviewProtoMapper),
	}
}
