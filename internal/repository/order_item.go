package repository

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
	orderitem_errors "ecommerce/pkg/errors/order_item_errors"
)

type orderItemRepository struct {
	db *db.Queries
}

func NewOrderItemRepository(db *db.Queries) *orderItemRepository {
	return &orderItemRepository{
		db: db,
	}
}

func (r *orderItemRepository) FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItems(ctx, reqDb)

	if err != nil {
		return nil, orderitem_errors.ErrFindAllOrderItems
	}

	return res, nil
}

func (r *orderItemRepository) FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItemsActive(ctx, reqDb)

	if err != nil {
		return nil, orderitem_errors.ErrFindByActive
	}

	return res, nil
}

func (r *orderItemRepository) FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItemsTrashed(ctx, reqDb)

	if err != nil {
		return nil, orderitem_errors.ErrFindByTrashed
	}

	return res, nil
}

func (r *orderItemRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*db.GetOrderItemsByOrderRow, error) {
	res, err := r.db.GetOrderItemsByOrder(ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	return res, nil
}

func (r *orderItemRepository) FindOrderItemByOrderTrashed(ctx context.Context, order_id int) ([]*db.OrderItem, error) {
	res, err := r.db.GetOrderItemsByOrderTrashed(ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	return res, nil
}

func (r *orderItemRepository) CalculateTotalPrice(ctx context.Context, order_id int) (*int32, error) {
	res, err := r.db.CalculateTotalPrice(ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrCalculateTotalPrice
	}

	return &res, nil

}

func (r *orderItemRepository) CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*db.CreateOrderItemRow, error) {
	res, err := r.db.CreateOrderItem(ctx, db.CreateOrderItemParams{
		OrderID:   int32(req.OrderID),
		ProductID: int32(req.ProductID),
		Quantity:  int32(req.Quantity),
		Price:     int32(req.Price),
	})

	if err != nil {
		return nil, orderitem_errors.ErrCreateOrderItem
	}

	return res, nil
}

func (r *orderItemRepository) UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*db.UpdateOrderItemRow, error) {
	res, err := r.db.UpdateOrderItem(ctx, db.UpdateOrderItemParams{
		OrderItemID: int32(req.OrderItemID),
		Quantity:    int32(req.Quantity),
		Price:       int32(req.Price),
	})

	if err != nil {
		return nil, orderitem_errors.ErrUpdateOrderItem
	}

	return res, nil
}

func (r *orderItemRepository) TrashedOrderItem(ctx context.Context, order_id int) (*db.OrderItem, error) {
	res, err := r.db.TrashOrderItem(ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrTrashedOrderItem
	}

	return res, nil
}

func (r *orderItemRepository) RestoreOrderItem(ctx context.Context, order_id int) (*db.OrderItem, error) {
	res, err := r.db.RestoreOrderItem(ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrRestoreOrderItem
	}

	return res, nil
}

func (r *orderItemRepository) DeleteOrderItemPermanent(ctx context.Context, order_id int) (bool, error) {
	err := r.db.DeleteOrderItemPermanently(ctx, int32(order_id))

	if err != nil {
		return false, orderitem_errors.ErrDeleteOrderItemPermanent
	}

	return true, nil
}

func (r *orderItemRepository) RestoreAllOrderItem(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllUsers(ctx)

	if err != nil {
		return false, orderitem_errors.ErrRestoreAllOrderItem
	}
	return true, nil
}

func (r *orderItemRepository) DeleteAllOrderPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentOrders(ctx)

	if err != nil {
		return false, orderitem_errors.ErrDeleteAllOrderPermanent
	}

	return true, nil
}
