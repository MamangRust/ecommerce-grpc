package repository

import (
	"context"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"errors"
	"fmt"
)

type orderRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.OrderRecordMapping
}

func NewOrderRepository(db *db.Queries, ctx context.Context, mapping recordmapper.OrderRecordMapping) *orderRepository {
	return &orderRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *orderRepository) FindAllOrders(search string, page, pageSize int) ([]*record.OrderRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetOrdersParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrders(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find users: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrdersRecordPagination(res), totalCount, nil
}

func (r *orderRepository) FindByActive(search string, page, pageSize int) ([]*record.OrderRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetOrdersActiveParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersActive(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find users: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrdersRecordActivePagination(res), totalCount, nil
}

func (r *orderRepository) FindByTrashed(search string, page, pageSize int) ([]*record.OrderRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetOrdersTrashedParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersTrashed(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find users: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrdersRecordTrashedPagination(res), totalCount, nil
}

func (r *orderRepository) FindByMerchant(merchant_id int, search string, page, pageSize int) ([]*record.OrderRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetOrdersByMerchantParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersByMerchant(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find users: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrdersRecordByMerchantPagination(res), totalCount, nil
}

func (r *orderRepository) FindById(user_id int) (*record.OrderRecord, error) {
	res, err := r.db.GetOrderByID(r.ctx, int32(user_id))

	if err != nil {
		fmt.Printf("Error fetching user: %v\n", err)

		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	return r.mapping.ToOrderRecord(res), nil
}

func (r *orderRepository) CreateOrder(request *requests.CreateOrderRecordRequest) (*record.OrderRecord, error) {
	req := db.CreateOrderParams{
		MerchantID: int32(request.MerchantID),
		UserID:     int32(request.UserID),
		TotalPrice: int32(request.TotalPrice),
	}

	user, err := r.db.CreateOrder(r.ctx, req)

	if err != nil {
		return nil, errors.New("failed create order")
	}

	return r.mapping.ToOrderRecord(user), nil
}

func (r *orderRepository) UpdateOrder(request *requests.UpdateOrderRecordRequest) (*record.OrderRecord, error) {
	req := db.UpdateOrderParams{
		OrderID:    int32(request.OrderID),
		UserID:     int32(request.UserID),
		TotalPrice: int32(request.TotalPrice),
	}

	res, err := r.db.UpdateOrder(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return r.mapping.ToOrderRecord(res), nil
}

func (r *orderRepository) TrashedOrder(user_id int) (*record.OrderRecord, error) {
	res, err := r.db.TrashOrder(r.ctx, int32(user_id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash order: %w", err)
	}

	return r.mapping.ToOrderRecord(res), nil
}

func (r *orderRepository) RestoreOrder(user_id int) (*record.OrderRecord, error) {
	res, err := r.db.RestoreOrder(r.ctx, int32(user_id))

	if err != nil {
		return nil, fmt.Errorf("failed to restore order: %w", err)
	}

	return r.mapping.ToOrderRecord(res), nil
}

func (r *orderRepository) DeleteOrderPermanent(user_id int) (bool, error) {
	err := r.db.DeleteOrderPermanently(r.ctx, int32(user_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete order: %w", err)
	}

	return true, nil
}

func (r *orderRepository) RestoreAllOrder() (bool, error) {
	err := r.db.RestoreAllOrders(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all orders: %w", err)
	}
	return true, nil
}

func (r *orderRepository) DeleteAllOrderPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentOrders(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all orders permanently: %w", err)
	}
	return true, nil
}
