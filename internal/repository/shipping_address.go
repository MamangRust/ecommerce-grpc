package repository

import (
	"context"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	shippingaddress_errors "ecommerce/pkg/errors/shipping_address_errors"
)

type shippingAdddressRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.ShippingAddressMapping
}

func NewShippingAddressRepository(
	db *db.Queries,
	ctx context.Context,
	mapping recordmapper.ShippingAddressMapping,
) *shippingAdddressRepository {
	return &shippingAdddressRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *shippingAdddressRepository) FindAllShippingAddress(req *requests.FindAllShippingAddress) ([]*record.ShippingAddressRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetShippingAddressParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetShippingAddress(r.ctx, reqDb)

	if err != nil {
		return nil, nil, shippingaddress_errors.ErrFindAllShippingAddress
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToShippingAddresssRecordPagination(res), &totalCount, nil
}

func (r *shippingAdddressRepository) FindByActive(req *requests.FindAllShippingAddress) ([]*record.ShippingAddressRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetShippingAddressActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetShippingAddressActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, shippingaddress_errors.ErrFindActiveShippingAddress
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToShippingAddresssRecordActivePagination(res), &totalCount, nil
}

func (r *shippingAdddressRepository) FindByTrashed(req *requests.FindAllShippingAddress) ([]*record.ShippingAddressRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetShippingAddressTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetShippingAddressTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, shippingaddress_errors.ErrFindTrashedShippingAddress
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToShippingAddresssRecordTrashedPagination(res), &totalCount, nil
}

func (r *shippingAdddressRepository) FindById(shipping_id int) (*record.ShippingAddressRecord, error) {
	res, err := r.db.GetShippingByID(r.ctx, int32(shipping_id))

	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByID
	}

	return r.mapping.ToShippingAddressRecord(res), nil
}

func (r *shippingAdddressRepository) FindByOrder(order_id int) (*record.ShippingAddressRecord, error) {
	res, err := r.db.GetShippingAddressByOrderID(r.ctx, int32(order_id))

	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByOrder
	}

	return r.mapping.ToShippingAddressRecord(res), nil
}

func (r *shippingAdddressRepository) CreateShippingAddress(request *requests.CreateShippingAddressRequest) (*record.ShippingAddressRecord, error) {
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

	address, err := r.db.CreateShippingAddress(r.ctx, req)

	if err != nil {
		return nil, shippingaddress_errors.ErrCreateShippingAddress
	}

	return r.mapping.ToShippingAddressRecord(address), nil
}

func (r *shippingAdddressRepository) UpdateShippingAddress(request *requests.UpdateShippingAddressRequest) (*record.ShippingAddressRecord, error) {
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


	res, err := r.db.UpdateShippingAddress(r.ctx, req)
	if err != nil {
		return nil, shippingaddress_errors.ErrUpdateShippingAddress
	}

	return r.mapping.ToShippingAddressRecord(res), nil
}

func (r *shippingAdddressRepository) TrashShippingAddress(shipping_id int) (*record.ShippingAddressRecord, error) {
	res, err := r.db.TrashShippingAddress(r.ctx, int32(shipping_id))

	if err != nil {
		return nil, shippingaddress_errors.ErrTrashShippingAddress
	}

	return r.mapping.ToShippingAddressRecord(res), nil
}

func (r *shippingAdddressRepository) RestoreShippingAddress(category_id int) (*record.ShippingAddressRecord, error) {
	res, err := r.db.RestoreShippingAddress(r.ctx, int32(category_id))

	if err != nil {
		return nil, shippingaddress_errors.ErrRestoreShippingAddress
	}

	return r.mapping.ToShippingAddressRecord(res), nil
}

func (r *shippingAdddressRepository) DeleteShippingAddressPermanently(category_id int) (bool, error) {
	err := r.db.DeleteShippingAddressPermanently(r.ctx, int32(category_id))

	if err != nil {
		return false, shippingaddress_errors.ErrDeleteShippingAddressPermanent
	}

	return true, nil
}

func (r *shippingAdddressRepository) RestoreAllShippingAddress() (bool, error) {
	err := r.db.RestoreAllShippingAddress(r.ctx)

	if err != nil {
		return false, shippingaddress_errors.ErrRestoreAllShippingAddresses
	}
	return true, nil
}

func (r *shippingAdddressRepository) DeleteAllPermanentShippingAddress() (bool, error) {
	err := r.db.DeleteAllPermanentShippingAddress(r.ctx)

	if err != nil {
		return false, shippingaddress_errors.ErrDeleteAllPermanentShippingAddress
	}
	return true, nil
}
