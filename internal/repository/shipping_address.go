package repository

import (
	"context"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"fmt"
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

func (r *shippingAdddressRepository) FindAllShippingAddress(search string, page, pageSize int) ([]*record.ShippingAddressRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetShippingAddressParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetShippingAddress(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find shipping address: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToShippingAddresssRecordPagination(res), totalCount, nil
}

func (r *shippingAdddressRepository) FindById(shipping_id int) (*record.ShippingAddressRecord, error) {
	res, err := r.db.GetShippingByID(r.ctx, int32(shipping_id))

	if err != nil {

		return nil, fmt.Errorf("failed to find shipping address: %w", err)
	}

	return r.mapping.ToShippingAddressRecord(res), nil
}

func (r *shippingAdddressRepository) FindByOrder(order_id int) (*record.ShippingAddressRecord, error) {
	res, err := r.db.GetShippingAddressByOrderID(r.ctx, int32(order_id))

	if err != nil {

		return nil, fmt.Errorf("failed to find shipping address: %w", err)
	}

	return r.mapping.ToShippingAddressRecord(res), nil
}

func (r *shippingAdddressRepository) FindByActive(search string, page, pageSize int) ([]*record.ShippingAddressRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetShippingAddressActiveParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetShippingAddressActive(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find shipping address: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToShippingAddresssRecordActivePagination(res), totalCount, nil
}

func (r *shippingAdddressRepository) FindByTrashed(search string, page, pageSize int) ([]*record.ShippingAddressRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetShippingAddressTrashedParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetShippingAddressTrashed(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find shipping address: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToShippingAddresssRecordTrashedPagination(res), totalCount, nil
}

func (r *shippingAdddressRepository) CreateShippingAddress(request *requests.CreateShippingAddressRequest) (*record.ShippingAddressRecord, error) {
	req := db.CreateShippingAddressParams{
		OrderID:        int32(request.OrderID),
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
		return nil, fmt.Errorf("failed to create shipping address: %w", err)
	}

	return r.mapping.ToShippingAddressRecord(address), nil
}

func (r *shippingAdddressRepository) UpdateShippingAddress(request *requests.UpdateShippingAddressRequest) (*record.ShippingAddressRecord, error) {
	req := db.UpdateShippingAddressParams{
		ShippingAddressID: int32(request.ShippingID),
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
		return nil, fmt.Errorf("failed to update shipping address: %w", err)
	}

	return r.mapping.ToShippingAddressRecord(res), nil
}

func (r *shippingAdddressRepository) TrashShippingAddress(shipping_id int) (*record.ShippingAddressRecord, error) {
	res, err := r.db.TrashShippingAddress(r.ctx, int32(shipping_id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash shipping address: %w", err)
	}

	return r.mapping.ToShippingAddressRecord(res), nil
}

func (r *shippingAdddressRepository) RestoreShippingAddress(category_id int) (*record.ShippingAddressRecord, error) {
	res, err := r.db.RestoreShippingAddress(r.ctx, int32(category_id))

	if err != nil {
		return nil, fmt.Errorf("failed to shipping address: %w", err)
	}

	return r.mapping.ToShippingAddressRecord(res), nil
}

func (r *shippingAdddressRepository) DeleteShippingAddressPermanently(category_id int) (bool, error) {
	err := r.db.DeleteShippingAddressPermanently(r.ctx, int32(category_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete shipping address: %w", err)
	}

	return true, nil
}

func (r *shippingAdddressRepository) RestoreAllShippingAddress() (bool, error) {
	err := r.db.RestoreAllShippingAddress(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all shipping address: %w", err)
	}
	return true, nil
}

func (r *shippingAdddressRepository) DeleteAllPermanentShippingAddress() (bool, error) {
	err := r.db.DeleteAllPermanentShippingAddress(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all shipping address permanently: %w", err)
	}
	return true, nil
}
