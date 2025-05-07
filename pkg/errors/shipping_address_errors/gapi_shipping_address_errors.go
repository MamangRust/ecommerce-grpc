package shippingaddress_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidID = response.NewGrpcError("error", "invalid ID", int(codes.InvalidArgument))

	ErrGrpcFailedFindAllShippingAddresses            = response.NewGrpcError("error", "Failed to fetch all shipping addresses", int(codes.Internal))
	ErrGrpcFailedFindShippingAddressById             = response.NewGrpcError("error", "Failed to find shipping address by ID", int(codes.Internal))
	ErrGrpcFailedFindShippingAddressByOrder          = response.NewGrpcError("error", "Failed to find shipping address by order ID", int(codes.Internal))
	ErrGrpcFailedFindActiveShippingAddresses         = response.NewGrpcError("error", "Failed to fetch active shipping addresses", int(codes.Internal))
	ErrGrpcFailedFindTrashedShippingAddresses        = response.NewGrpcError("error", "Failed to fetch trashed shipping addresses", int(codes.Internal))
	ErrGrpcFailedTrashShippingAddress                = response.NewGrpcError("error", "Failed to trash shipping address", int(codes.Internal))
	ErrGrpcFailedRestoreShippingAddress              = response.NewGrpcError("error", "Failed to restore shipping address", int(codes.Internal))
	ErrGrpcFailedDeleteShippingAddressPermanent      = response.NewGrpcError("error", "Failed to permanently delete shipping address", int(codes.Internal))
	ErrGrpcFailedRestoreAllShippingAddresses         = response.NewGrpcError("error", "Failed to restore all shipping addresses", int(codes.Internal))
	ErrGrpcFailedDeleteAllShippingAddressesPermanent = response.NewGrpcError("error", "Failed to permanently delete all shipping addresses", int(codes.Internal))
)
