package gapi

import "ecommerce/internal/pb"

type AuthHandleGrpc interface {
	pb.AuthServiceServer
}

type RoleHandleGrpc interface {
	pb.RoleServiceServer
}

type UserHandleGrpc interface {
	pb.UserServiceServer
}

type CategoryHandleGrpc interface {
	pb.CategoryServiceServer
}

type MerchantHandleGrpc interface {
	pb.MerchantServiceServer
}

type OrderItemHandleGrpc interface {
	pb.OrderItemServiceServer
}

type OrderHandleGrpc interface {
	pb.OrderServiceServer
}

type ProductHandleGrpc interface {
	pb.ProductServiceServer
}

type TransactionHandleGrpc interface {
	pb.TransactionServiceServer
}

type CartHandleGrpc interface {
	pb.CartServiceServer
}

type ReviewHandleGrpc interface {
	pb.ReviewServiceServer
}

type ShippingAdddressHandleGrpc interface {
	pb.ShippingServiceServer
}

type SliderHandleGrpc interface {
	pb.SliderServiceServer
}

type BannerHandleGrpc interface {
	pb.BannerServiceServer
}

type MerchantAwardHandleGrpc interface {
	pb.MerchantAwardServiceServer
}

type MerchantBusinessHandleGrpc interface {
	pb.MerchantBusinessServiceServer
}

type MerchantDetailHandleGrpc interface {
	pb.MerchantDetailServiceServer
}

type MerchantPoliciesHandleGrpc interface {
	pb.MerchantPoliciesServiceServer
}

type ReviewDetailHandleGrpc interface {
	pb.ReviewDetailServiceServer
}
