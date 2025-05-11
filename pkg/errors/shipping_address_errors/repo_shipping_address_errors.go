package shippingaddress_errors

import "errors"

var (
	ErrFindAllShippingAddress            = errors.New("failed to find all shipping addresses")
	ErrFindActiveShippingAddress         = errors.New("failed to find active shipping addresses")
	ErrFindTrashedShippingAddress        = errors.New("failed to find trashed shipping addresses")
	ErrFindShippingAddressByID           = errors.New("failed to find shipping address by ID")
	ErrFindShippingAddressByOrder        = errors.New("failed to find shipping address by order ID")
	ErrCreateShippingAddress             = errors.New("failed to create shipping address")
	ErrUpdateShippingAddress             = errors.New("failed to update shipping address")
	ErrTrashShippingAddress              = errors.New("failed to trash shipping address")
	ErrRestoreShippingAddress            = errors.New("failed to restore shipping address")
	ErrDeleteShippingAddressPermanent    = errors.New("failed to permanently delete shipping address")
	ErrRestoreAllShippingAddresses       = errors.New("failed to restore all shipping addresses")
	ErrDeleteAllPermanentShippingAddress = errors.New("failed to permanently delete all shipping addresses")
)
