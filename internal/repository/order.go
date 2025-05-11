package repository

import (
	"context"
	"database/sql"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/order_errors"
	"fmt"
	"time"
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

func (r *orderRepository) FindAllOrders(req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrders(r.ctx, reqDb)

	if err != nil {
		return nil, nil, order_errors.ErrFindAllOrders
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrdersRecordPagination(res), &totalCount, nil
}

func (r *orderRepository) FindByActive(req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, order_errors.ErrFindByActive
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrdersRecordActivePagination(res), &totalCount, nil
}

func (r *orderRepository) FindByTrashed(req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, order_errors.ErrFindByTrashed
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrdersRecordTrashedPagination(res), &totalCount, nil
}

func (r *orderRepository) FindByMerchant(req *requests.FindAllOrderByMerchant) ([]*record.OrderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersByMerchantParams{
		Column1: req.Search,
		Column4: int32(req.MerchantID),
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersByMerchant(r.ctx, reqDb)

	if err != nil {
		return nil, nil, order_errors.ErrFindByMerchant
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrdersRecordByMerchantPagination(res), &totalCount, nil
}

func (r *orderRepository) GetMonthlyTotalRevenue(req *requests.MonthTotalRevenue) ([]*record.OrderMonthlyTotalRevenueRecord, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalRevenue(r.ctx, db.GetMonthlyTotalRevenueParams{
		Extract:     currentMonthStart,
		CreatedAt:   sql.NullTime{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: sql.NullTime{Time: prevMonthStart, Valid: true},
		CreatedAt_3: sql.NullTime{Time: prevMonthEnd, Valid: true},
	})

	if err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenue
	}

	so := r.mapping.ToOrderMonthlyTotalRevenues(res)

	return so, nil
}

func (r *orderRepository) GetYearlyTotalRevenue(year int) ([]*record.OrderYearlyTotalRevenueRecord, error) {
	res, err := r.db.GetYearlyTotalRevenue(r.ctx, int32(year))

	if err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenue
	}

	fmt.Println("err", err)

	so := r.mapping.ToOrderYearlyTotalRevenues(res)

	return so, nil
}

func (r *orderRepository) GetMonthlyTotalRevenueById(req *requests.MonthTotalRevenueOrder) ([]*record.OrderMonthlyTotalRevenueRecord, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalRevenueById(r.ctx, db.GetMonthlyTotalRevenueByIdParams{
		Extract:     currentMonthStart,
		CreatedAt:   sql.NullTime{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: sql.NullTime{Time: prevMonthStart, Valid: true},
		CreatedAt_3: sql.NullTime{Time: prevMonthEnd, Valid: true},
		OrderID:     int32(req.OrderID),
	})

	if err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenueById
	}

	so := r.mapping.ToOrderMonthlyTotalRevenuesById(res)

	return so, nil
}

func (r *orderRepository) GetYearlyTotalRevenueById(req *requests.YearTotalRevenueOrder) ([]*record.OrderYearlyTotalRevenueRecord, error) {
	res, err := r.db.GetYearlyTotalRevenueById(r.ctx, db.GetYearlyTotalRevenueByIdParams{
		Column1: int32(req.Year),
		OrderID: int32(req.OrderID),
	})

	if err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenueById
	}

	so := r.mapping.ToOrderYearlyTotalRevenuesById(res)

	return so, nil
}

func (r *orderRepository) GetMonthlyTotalRevenueByMerchant(req *requests.MonthTotalRevenueMerchant) ([]*record.OrderMonthlyTotalRevenueRecord, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalRevenueByMerchant(r.ctx, db.GetMonthlyTotalRevenueByMerchantParams{
		Extract:     currentMonthStart,
		CreatedAt:   sql.NullTime{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: sql.NullTime{Time: prevMonthStart, Valid: true},
		CreatedAt_3: sql.NullTime{Time: prevMonthEnd, Valid: true},
		MerchantID:  int32(req.MerchantID),
	})

	if err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenueByMerchant
	}

	so := r.mapping.ToOrderMonthlyTotalRevenuesByMerchant(res)

	return so, nil
}

func (r *orderRepository) GetYearlyTotalRevenueByMerchant(req *requests.YearTotalRevenueMerchant) ([]*record.OrderYearlyTotalRevenueRecord, error) {
	res, err := r.db.GetYearlyTotalRevenueByMerchant(r.ctx, db.GetYearlyTotalRevenueByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenueByMerchant
	}

	so := r.mapping.ToOrderYearlyTotalRevenuesByMerchant(res)

	return so, nil
}

func (r *orderRepository) GetMonthlyOrder(year int) ([]*record.OrderMonthlyRecord, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyOrder(r.ctx, yearStart)

	if err != nil {
		return nil, order_errors.ErrGetMonthlyOrder
	}

	return r.mapping.ToOrderMonthlyPrices(res), nil
}

func (r *orderRepository) GetYearlyOrder(year int) ([]*record.OrderYearlyRecord, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyOrder(r.ctx, yearStart)
	if err != nil {
		return nil, order_errors.ErrGetYearlyOrder
	}

	return r.mapping.ToOrderYearlyPrices(res), nil
}

func (r *orderRepository) GetMonthlyOrderByMerchant(req *requests.MonthOrderMerchant) ([]*record.OrderMonthlyRecord, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyOrderByMerchant(r.ctx, db.GetMonthlyOrderByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})
	if err != nil {
		return nil, order_errors.ErrGetMonthlyOrderByMerchant
	}

	return r.mapping.ToOrderMonthlyPricesByMerchant(res), nil
}

func (r *orderRepository) GetYearlyOrderByMerchant(req *requests.YearOrderMerchant) ([]*record.OrderYearlyRecord, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyOrderByMerchant(r.ctx, db.GetYearlyOrderByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, order_errors.ErrGetYearlyOrderByMerchant
	}

	return r.mapping.ToOrderYearlyPricesByMerchant(res), nil
}

func (r *orderRepository) FindById(user_id int) (*record.OrderRecord, error) {
	res, err := r.db.GetOrderByID(r.ctx, int32(user_id))

	if err != nil {
		return nil, order_errors.ErrFindById
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
		return nil, order_errors.ErrCreateOrder
	}

	return r.mapping.ToOrderRecord(user), nil
}

func (r *orderRepository) UpdateOrder(request *requests.UpdateOrderRecordRequest) (*record.OrderRecord, error) {
	req := db.UpdateOrderParams{
		OrderID:    int32(request.OrderID),
		TotalPrice: int32(request.TotalPrice),
	}

	res, err := r.db.UpdateOrder(r.ctx, req)

	if err != nil {
		return nil, order_errors.ErrUpdateOrder
	}

	return r.mapping.ToOrderRecord(res), nil
}

func (r *orderRepository) TrashedOrder(user_id int) (*record.OrderRecord, error) {
	res, err := r.db.TrashedOrder(r.ctx, int32(user_id))

	if err != nil {
		return nil, order_errors.ErrTrashedOrder
	}

	return r.mapping.ToOrderRecord(res), nil
}

func (r *orderRepository) RestoreOrder(user_id int) (*record.OrderRecord, error) {
	res, err := r.db.RestoreOrder(r.ctx, int32(user_id))

	if err != nil {
		return nil, order_errors.ErrRestoreOrder
	}

	return r.mapping.ToOrderRecord(res), nil
}

func (r *orderRepository) DeleteOrderPermanent(user_id int) (bool, error) {
	err := r.db.DeleteOrderPermanently(r.ctx, int32(user_id))

	if err != nil {
		return false, order_errors.ErrDeleteOrderPermanent
	}

	return true, nil
}

func (r *orderRepository) RestoreAllOrder() (bool, error) {
	err := r.db.RestoreAllOrders(r.ctx)

	if err != nil {
		return false, order_errors.ErrRestoreAllOrder
	}
	return true, nil
}

func (r *orderRepository) DeleteAllOrderPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentOrders(r.ctx)

	if err != nil {
		return false, order_errors.ErrDeleteAllOrderPermanent
	}
	return true, nil
}
