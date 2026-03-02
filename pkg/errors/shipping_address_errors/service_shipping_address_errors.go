package shippingaddress_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrFailedCreateShippingAddress = errors.NewErrorResponse("Failed to create shipping address", http.StatusInternalServerError)
	ErrFailedUpdateShippingAddress = errors.NewErrorResponse("Failed to update shipping address", http.StatusInternalServerError)

	ErrFailedFindAllShippingAddresses            = errors.NewErrorResponse("Failed to fetch all shipping addresses", http.StatusInternalServerError)
	ErrFailedFindActiveShippingAddresses         = errors.NewErrorResponse("Failed to fetch active shipping addresses", http.StatusInternalServerError)
	ErrFailedFindTrashedShippingAddresses        = errors.NewErrorResponse("Failed to fetch trashed shipping addresses", http.StatusInternalServerError)
	ErrFailedFindShippingAddressByID             = errors.NewErrorResponse("Failed to find shipping address by ID", http.StatusInternalServerError)
	ErrFailedFindShippingAddressByOrder          = errors.NewErrorResponse("Failed to find shipping address by order ID", http.StatusInternalServerError)
	ErrFailedTrashShippingAddress                = errors.NewErrorResponse("Failed to trash shipping address", http.StatusInternalServerError)
	ErrFailedRestoreShippingAddress              = errors.NewErrorResponse("Failed to restore shipping address", http.StatusInternalServerError)
	ErrFailedDeleteShippingAddressPermanent      = errors.NewErrorResponse("Failed to permanently delete shipping address", http.StatusInternalServerError)
	ErrFailedRestoreAllShippingAddresses         = errors.NewErrorResponse("Failed to restore all shipping addresses", http.StatusInternalServerError)
	ErrFailedDeleteAllShippingAddressesPermanent = errors.NewErrorResponse("Failed to permanently delete all shipping addresses", http.StatusInternalServerError)
)
