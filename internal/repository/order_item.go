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

type orderItemRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.OrderItemRecordMapping
}

func NewOrderItemRepository(db *db.Queries, ctx context.Context, mapping recordmapper.OrderItemRecordMapping) *orderItemRepository {
	return &orderItemRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *orderItemRepository) FindAllOrderItems(search string, page, pageSize int) ([]*record.OrderItemRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetOrderItemsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItems(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find order items : %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrderItemsRecordPagination(res), totalCount, nil
}

func (r *orderItemRepository) FindByActive(search string, page, pageSize int) ([]*record.OrderItemRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetOrderItemsActiveParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItemsActive(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find users: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrderItemsRecordActivePagination(res), totalCount, nil
}

func (r *orderItemRepository) FindByTrashed(search string, page, pageSize int) ([]*record.OrderItemRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetOrderItemsTrashedParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItemsTrashed(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find users: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrderItemsRecordTrashedPagination(res), totalCount, nil
}

func (r *orderItemRepository) FindOrderItemByOrder(order_id int) ([]*record.OrderItemRecord, error) {
	res, err := r.db.GetOrderItemsByOrder(r.ctx, int32(order_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find order by order_id: %w", err)
	}

	return r.mapping.ToOrderItemsRecord(res), nil
}

func (r *orderItemRepository) CalculateTotalPrice(order_id int) (*int32, error) {
	res, err := r.db.CalculateTotalPrice(r.ctx, int32(order_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find order by order_id: %w", err)
	}

	return &res, nil

}

func (r *orderItemRepository) CreateOrderItem(req *requests.CreateOrderItemRecordRequest) (*record.OrderItemRecord, error) {
	res, err := r.db.CreateOrderItem(r.ctx, db.CreateOrderItemParams{
		OrderID:   int32(req.OrderID),
		ProductID: int32(req.ProductID),
		Quantity:  int32(req.Quantity),
		Price:     int32(req.Price),
	})

	if err != nil {
		return nil, errors.New("failed create order item")
	}

	return r.mapping.ToOrderItemRecord(res), nil
}

func (r *orderItemRepository) UpdateOrderItem(req *requests.UpdateOrderItemRecordRequest) (*record.OrderItemRecord, error) {
	res, err := r.db.UpdateOrderItem(r.ctx, db.UpdateOrderItemParams{
		OrderItemID: int32(req.OrderItemID),
		Quantity:    int32(req.Quantity),
		Price:       int32(req.Price),
	})

	if err != nil {
		return nil, errors.New("failed update order item")
	}

	return r.mapping.ToOrderItemRecord(res), nil
}

func (r *orderItemRepository) TrashedOrderItem(order_id int) (*record.OrderItemRecord, error) {
	res, err := r.db.TrashOrderItem(r.ctx, int32(order_id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash order item: %w", err)
	}

	return r.mapping.ToOrderItemRecord(res), nil
}

func (r *orderItemRepository) RestoreOrderItem(order_id int) (*record.OrderItemRecord, error) {
	res, err := r.db.RestoreOrderItem(r.ctx, int32(order_id))

	if err != nil {
		return nil, fmt.Errorf("failed to restore order item: %w", err)
	}

	return r.mapping.ToOrderItemRecord(res), nil
}

func (r *orderItemRepository) DeleteOrderItemPermanent(order_id int) (bool, error) {
	err := r.db.DeleteOrderItemPermanently(r.ctx, int32(order_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete order_item: %w", err)
	}

	return true, nil
}

func (r *orderItemRepository) RestoreAllOrderItem() (bool, error) {
	err := r.db.RestoreAllUsers(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all order_item: %w", err)
	}
	return true, nil
}

func (r *orderItemRepository) DeleteAllOrderPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentOrders(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all order_item permanently: %w", err)
	}
	return true, nil
}
