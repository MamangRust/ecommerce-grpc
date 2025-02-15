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

type cartRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.CartRecordMapping
}

func NewCartRepository(
	db *db.Queries,
	ctx context.Context,
	mapping recordmapper.CartRecordMapping,
) *cartRepository {
	return &cartRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *cartRepository) FindCarts(user_id int, search string, page, pageSize int) ([]*record.CartRecord, int, error) {

	offset := (page - 1) * pageSize

	req := db.GetCartsParams{
		UserID:  int32(user_id),
		Column2: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCarts(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find categories: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCartsRecordPagination(res), totalCount, nil
}

func (r *cartRepository) CreateCart(req *requests.CartCreateRecord) (*record.CartRecord, error) {
	res, err := r.db.CreateCart(r.ctx, db.CreateCartParams{
		UserID:    int32(req.UserID),
		ProductID: int32(req.ProductID),
		Name:      req.Name,
		Price:     int32(req.Price),
		Image:     req.ImageProduct,
		Quantity:  int32(req.Quantity),
		Weight:    int32(req.Weight),
	})

	if err != nil {
		return nil, errors.New("failed to create cart")
	}

	return r.mapping.ToCartRecord(res), nil
}

func (r *cartRepository) DeletePermanent(cart_id int) (bool, error) {
	err := r.db.DeleteCart(r.ctx, int32(cart_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete cart: %w", err)
	}

	return true, nil
}

func (r *cartRepository) DeleteAllPermanently(req *requests.DeleteCartRequest) (bool, error) {
	cartIDs := make([]int32, len(req.CartIds))
	for i, id := range req.CartIds {
		cartIDs[i] = int32(id)
	}

	err := r.db.DeleteAllCart(r.ctx, cartIDs)

	if err != nil {
		return false, fmt.Errorf("failed to delete carts: %w", err)
	}

	return true, nil
}
