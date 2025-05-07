package repository

import (
	"context"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	orderitem_errors "ecommerce/pkg/errors/order_item_errors"
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

func (r *orderItemRepository) FindAllOrderItems(req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItems(r.ctx, reqDb)

	if err != nil {
		return nil, nil, orderitem_errors.ErrFindAllOrderItems
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrderItemsRecordPagination(res), &totalCount, nil
}

func (r *orderItemRepository) FindByActive(req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItemsActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, orderitem_errors.ErrFindByActive
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrderItemsRecordActivePagination(res), &totalCount, nil
}

func (r *orderItemRepository) FindByTrashed(req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItemsTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, orderitem_errors.ErrFindByTrashed
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrderItemsRecordTrashedPagination(res), &totalCount, nil
}

func (r *orderItemRepository) FindOrderItemByOrder(order_id int) ([]*record.OrderItemRecord, error) {
	res, err := r.db.GetOrderItemsByOrder(r.ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	return r.mapping.ToOrderItemsRecord(res), nil
}

func (r *orderItemRepository) CalculateTotalPrice(order_id int) (*int32, error) {
	res, err := r.db.CalculateTotalPrice(r.ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrCalculateTotalPrice
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
		return nil, orderitem_errors.ErrCreateOrderItem
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
		return nil, orderitem_errors.ErrUpdateOrderItem
	}

	return r.mapping.ToOrderItemRecord(res), nil
}

func (r *orderItemRepository) TrashedOrderItem(order_id int) (*record.OrderItemRecord, error) {
	res, err := r.db.TrashOrderItem(r.ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrTrashedOrderItem
	}

	return r.mapping.ToOrderItemRecord(res), nil
}

func (r *orderItemRepository) RestoreOrderItem(order_id int) (*record.OrderItemRecord, error) {
	res, err := r.db.RestoreOrderItem(r.ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrRestoreOrderItem
	}

	return r.mapping.ToOrderItemRecord(res), nil
}

func (r *orderItemRepository) DeleteOrderItemPermanent(order_id int) (bool, error) {
	err := r.db.DeleteOrderItemPermanently(r.ctx, int32(order_id))

	if err != nil {
		return false, orderitem_errors.ErrDeleteOrderItemPermanent
	}

	return true, nil
}

func (r *orderItemRepository) RestoreAllOrderItem() (bool, error) {
	err := r.db.RestoreAllUsers(r.ctx)

	if err != nil {
		return false, orderitem_errors.ErrRestoreAllOrderItem
	}
	return true, nil
}

func (r *orderItemRepository) DeleteAllOrderPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentOrders(r.ctx)

	if err != nil {
		return false, orderitem_errors.ErrDeleteAllOrderPermanent
	}

	return true, nil
}
