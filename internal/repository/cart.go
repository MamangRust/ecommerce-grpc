package repository

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/cart_errors"
)

type cartRepository struct {
	db *db.Queries
}

func NewCartRepository(
	db *db.Queries,
) *cartRepository {
	return &cartRepository{
		db: db,
	}
}

func (r *cartRepository) FindCarts(ctx context.Context, req *requests.FindAllCarts) ([]*db.GetCartsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCartsParams{
		UserID:  int32(req.UserID),
		Column2: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCarts(ctx, reqDb)

	if err != nil {
		return nil, cart_errors.ErrFindAllCarts
	}

	return res, nil
}

func (r *cartRepository) CreateCart(ctx context.Context, req *requests.CartCreateRecord) (*db.Cart, error) {
	res, err := r.db.CreateCart(ctx, db.CreateCartParams{
		UserID:    int32(req.UserID),
		ProductID: int32(req.ProductID),
		Name:      req.Name,
		Price:     int32(req.Price),
		Image:     req.ImageProduct,
		Quantity:  int32(req.Quantity),
		Weight:    int32(req.Weight),
	})

	if err != nil {
		return nil, cart_errors.ErrCreateCart
	}

	return res, nil
}

func (r *cartRepository) DeletePermanent(ctx context.Context, cart_id int) (bool, error) {
	err := r.db.DeleteCart(ctx, int32(cart_id))

	if err != nil {
		return false, cart_errors.ErrDeleteCartPermanent
	}

	return true, nil
}

func (r *cartRepository) DeleteAllPermanently(ctx context.Context, req *requests.DeleteCartRequest) (bool, error) {
	cartIDs := make([]int32, len(req.CartIds))

	for i, id := range req.CartIds {
		cartIDs[i] = int32(id)
	}

	err := r.db.DeleteAllCart(ctx, cartIDs)

	if err != nil {
		return false, cart_errors.ErrDeleteAllCarts
	}

	return true, nil
}
