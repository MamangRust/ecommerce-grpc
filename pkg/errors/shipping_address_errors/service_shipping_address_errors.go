package shippingaddress_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrFailedCreateShippingAddress = response.NewErrorResponse("Failed to create shipping address", http.StatusInternalServerError)
	ErrFailedUpdateShippingAddress = response.NewErrorResponse("Failed to update shipping address", http.StatusInternalServerError)

	ErrFailedFindAllShippingAddresses            = response.NewErrorResponse("Failed to fetch all shipping addresses", http.StatusInternalServerError)
	ErrFailedFindActiveShippingAddresses         = response.NewErrorResponse("Failed to fetch active shipping addresses", http.StatusInternalServerError)
	ErrFailedFindTrashedShippingAddresses        = response.NewErrorResponse("Failed to fetch trashed shipping addresses", http.StatusInternalServerError)
	ErrFailedFindShippingAddressByID             = response.NewErrorResponse("Failed to find shipping address by ID", http.StatusInternalServerError)
	ErrFailedFindShippingAddressByOrder          = response.NewErrorResponse("Failed to find shipping address by order ID", http.StatusInternalServerError)
	ErrFailedTrashShippingAddress                = response.NewErrorResponse("Failed to trash shipping address", http.StatusInternalServerError)
	ErrFailedRestoreShippingAddress              = response.NewErrorResponse("Failed to restore shipping address", http.StatusInternalServerError)
	ErrFailedDeleteShippingAddressPermanent      = response.NewErrorResponse("Failed to permanently delete shipping address", http.StatusInternalServerError)
	ErrFailedRestoreAllShippingAddresses         = response.NewErrorResponse("Failed to restore all shipping addresses", http.StatusInternalServerError)
	ErrFailedDeleteAllShippingAddressesPermanent = response.NewErrorResponse("Failed to permanently delete all shipping addresses", http.StatusInternalServerError)
)
