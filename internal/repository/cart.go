package repository

import (
	"context"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/cart_errors"
	"log"
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

func (r *cartRepository) FindCarts(req *requests.FindAllCarts) ([]*record.CartRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCartsParams{
		UserID:  int32(req.UserID),
		Column2: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCarts(r.ctx, reqDb)

	if err != nil {
		log.Fatal(err)
		return nil, nil, cart_errors.ErrFindAllCarts
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCartsRecordPagination(res), &totalCount, nil
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
		return nil, cart_errors.ErrCreateCart
	}

	return r.mapping.ToCartRecord(res), nil
}

func (r *cartRepository) DeletePermanent(cart_id int) (bool, error) {
	err := r.db.DeleteCart(r.ctx, int32(cart_id))

	if err != nil {
		return false, cart_errors.ErrDeleteCartPermanent
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
		return false, cart_errors.ErrDeleteAllCarts
	}

	return true, nil
}
