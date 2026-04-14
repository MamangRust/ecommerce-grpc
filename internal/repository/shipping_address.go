package repository

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
	shippingaddress_errors "ecommerce/pkg/errors/shipping_address_errors"
)

type shippingAdddressRepository struct {
	db *db.Queries
}

func NewShippingAddressRepository(
	db *db.Queries,

) *shippingAdddressRepository {
	return &shippingAdddressRepository{
		db: db,
	}
}

func (r *shippingAdddressRepository) FindAllShippingAddress(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetShippingAddressParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetShippingAddress(ctx, reqDb)

	if err != nil {
		return nil, shippingaddress_errors.ErrFindAllShippingAddress
	}

	return res, nil
}

func (r *shippingAdddressRepository) FindByActive(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetShippingAddressActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetShippingAddressActive(ctx, reqDb)

	if err != nil {
		return nil, shippingaddress_errors.ErrFindActiveShippingAddress
	}

	return res, nil
}

func (r *shippingAdddressRepository) FindByTrashed(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetShippingAddressTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetShippingAddressTrashed(ctx, reqDb)

	if err != nil {
		return nil, shippingaddress_errors.ErrFindTrashedShippingAddress
	}

	return res, nil
}

func (r *shippingAdddressRepository) FindById(ctx context.Context, shipping_id int) (*db.GetShippingByIDRow, error) {
	res, err := r.db.GetShippingByID(ctx, int32(shipping_id))

	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByID
	}

	return res, nil
}

func (r *shippingAdddressRepository) FindByOrder(ctx context.Context, order_id int) (*db.GetShippingAddressByOrderIDRow, error) {
	res, err := r.db.GetShippingAddressByOrderID(ctx, int32(order_id))

	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByOrder
	}

	return res, nil
}

func (r *shippingAdddressRepository) FindTrashedByOrder(ctx context.Context, order_id int) (*db.ShippingAddress, error) {
	res, err := r.db.GetShippingAddressByOrderIDTrashed(ctx, int32(order_id))

	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByOrder
	}

	return res, nil
}

func (r *shippingAdddressRepository) CreateShippingAddress(ctx context.Context, request *requests.CreateShippingAddressRequest) (*db.CreateShippingAddressRow, error) {
	req := db.CreateShippingAddressParams{
		OrderID:        int32(*request.OrderID),
		Alamat:         request.Alamat,
		Provinsi:       request.Provinsi,
		Kota:           request.Kota,
		Negara:         request.Negara,
		Courier:        request.Courier,
		ShippingMethod: request.ShippingMethod,
		ShippingCost:   float64(request.ShippingCost),
	}

	address, err := r.db.CreateShippingAddress(ctx, req)

	if err != nil {
		return nil, shippingaddress_errors.ErrCreateShippingAddress
	}

	return address, nil
}

func (r *shippingAdddressRepository) UpdateShippingAddress(ctx context.Context, request *requests.UpdateShippingAddressRequest) (*db.UpdateShippingAddressRow, error) {
	req := db.UpdateShippingAddressParams{
		ShippingAddressID: int32(*request.ShippingID),
		Alamat:            request.Alamat,
		Provinsi:          request.Provinsi,
		Kota:              request.Kota,
		Negara:            request.Negara,
		Courier:           request.Courier,
		ShippingMethod:    request.ShippingMethod,
		ShippingCost:      float64(request.ShippingCost),
	}

	res, err := r.db.UpdateShippingAddress(ctx, req)
	if err != nil {
		return nil, shippingaddress_errors.ErrUpdateShippingAddress
	}

	return res, nil
}

func (r *shippingAdddressRepository) TrashShippingAddress(ctx context.Context, shipping_id int) (*db.ShippingAddress, error) {
	res, err := r.db.TrashShippingAddress(ctx, int32(shipping_id))

	if err != nil {
		return nil, shippingaddress_errors.ErrTrashShippingAddress
	}

	return res, nil
}

func (r *shippingAdddressRepository) RestoreShippingAddress(ctx context.Context, category_id int) (*db.ShippingAddress, error) {
	res, err := r.db.RestoreShippingAddress(ctx, int32(category_id))

	if err != nil {
		return nil, shippingaddress_errors.ErrRestoreShippingAddress
	}

	return res, nil
}

func (r *shippingAdddressRepository) DeleteShippingAddressPermanently(ctx context.Context, category_id int) (bool, error) {
	err := r.db.DeleteShippingAddressPermanently(ctx, int32(category_id))

	if err != nil {
		return false, shippingaddress_errors.ErrDeleteShippingAddressPermanent
	}

	return true, nil
}

func (r *shippingAdddressRepository) RestoreAllShippingAddress(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllShippingAddress(ctx)

	if err != nil {
		return false, shippingaddress_errors.ErrRestoreAllShippingAddresses
	}
	return true, nil
}

func (r *shippingAdddressRepository) DeleteAllPermanentShippingAddress(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentShippingAddress(ctx)

	if err != nil {
		return false, shippingaddress_errors.ErrDeleteAllPermanentShippingAddress
	}
	return true, nil
}
