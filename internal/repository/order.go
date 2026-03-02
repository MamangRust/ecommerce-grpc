package repository

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/order_errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type orderRepository struct {
	db *db.Queries
}

func NewOrderRepository(db *db.Queries) *orderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) FindAllOrders(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrders(ctx, reqDb)

	if err != nil {
		return nil, order_errors.ErrFindAllOrders
	}

	return res, nil
}

func (r *orderRepository) FindByActive(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersActive(ctx, reqDb)

	if err != nil {
		return nil, order_errors.ErrFindByActive
	}

	return res, nil
}

func (r *orderRepository) FindByTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersTrashed(ctx, reqDb)

	if err != nil {
		return nil, order_errors.ErrFindByTrashed
	}

	return res, nil
}

func (r *orderRepository) FindByMerchant(ctx context.Context, req *requests.FindAllOrderByMerchant) ([]*db.GetOrdersByMerchantRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersByMerchantParams{
		Column1: req.Search,
		Column4: int32(req.MerchantID),
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersByMerchant(ctx, reqDb)

	if err != nil {
		return nil, order_errors.ErrFindByMerchant
	}

	return res, nil
}

func (r *orderRepository) GetMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*db.GetMonthlyTotalRevenueRow, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalRevenue(ctx, db.GetMonthlyTotalRevenueParams{
		Extract: pgtype.Date{
			Time:  currentMonthStart,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  currentMonthEnd,
			Valid: true,
		},
		CreatedAt_2: pgtype.Timestamp{
			Time:  prevMonthStart,
			Valid: true,
		},
		CreatedAt_3: pgtype.Timestamp{
			Time:  prevMonthEnd,
			Valid: true,
		},
	})

	if err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenue
	}

	return res, nil
}

func (r *orderRepository) GetYearlyTotalRevenue(ctx context.Context, year int) ([]*db.GetYearlyTotalRevenueRow, error) {
	res, err := r.db.GetYearlyTotalRevenue(ctx, int32(year))

	if err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenue
	}

	return res, nil
}

func (r *orderRepository) GetMonthlyTotalRevenueById(ctx context.Context, req *requests.MonthTotalRevenueOrder) ([]*db.GetMonthlyTotalRevenueByIdRow, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalRevenueById(ctx, db.GetMonthlyTotalRevenueByIdParams{
		Extract: pgtype.Date{
			Time:  currentMonthStart,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  currentMonthEnd,
			Valid: true,
		},
		CreatedAt_2: pgtype.Timestamp{
			Time:  prevMonthStart,
			Valid: true,
		},
		CreatedAt_3: pgtype.Timestamp{
			Time:  prevMonthEnd,
			Valid: true,
		},
		OrderID: int32(req.OrderID),
	})

	if err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenueById
	}

	return res, nil
}

func (r *orderRepository) GetYearlyTotalRevenueById(ctx context.Context, req *requests.YearTotalRevenueOrder) ([]*db.GetYearlyTotalRevenueByIdRow, error) {
	res, err := r.db.GetYearlyTotalRevenueById(ctx, db.GetYearlyTotalRevenueByIdParams{
		Column1: int32(req.Year),
		OrderID: int32(req.OrderID),
	})

	if err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenueById
	}

	return res, nil
}

func (r *orderRepository) GetMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*db.GetMonthlyTotalRevenueByMerchantRow, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalRevenueByMerchant(ctx, db.GetMonthlyTotalRevenueByMerchantParams{
		Extract: pgtype.Date{
			Time:  currentMonthStart,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  currentMonthEnd,
			Valid: true,
		},
		CreatedAt_2: pgtype.Timestamp{
			Time:  prevMonthStart,
			Valid: true,
		},
		CreatedAt_3: pgtype.Timestamp{
			Time:  prevMonthEnd,
			Valid: true,
		},
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenueByMerchant
	}

	return res, nil
}

func (r *orderRepository) GetYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*db.GetYearlyTotalRevenueByMerchantRow, error) {
	res, err := r.db.GetYearlyTotalRevenueByMerchant(ctx, db.GetYearlyTotalRevenueByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenueByMerchant
	}

	return res, nil
}

func (r *orderRepository) GetMonthlyOrder(ctx context.Context, year int) ([]*db.GetMonthlyOrderRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyOrder(ctx, yearStart)

	if err != nil {
		return nil, order_errors.ErrGetMonthlyOrder
	}

	return res, nil
}

func (r *orderRepository) GetYearlyOrder(ctx context.Context, year int) ([]*db.GetYearlyOrderRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyOrder(ctx, yearStart)
	if err != nil {
		return nil, order_errors.ErrGetYearlyOrder
	}

	return res, nil
}

func (r *orderRepository) GetMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*db.GetMonthlyOrderByMerchantRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyOrderByMerchant(ctx, db.GetMonthlyOrderByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})
	if err != nil {
		return nil, order_errors.ErrGetMonthlyOrderByMerchant
	}

	return res, nil
}

func (r *orderRepository) GetYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*db.GetYearlyOrderByMerchantRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyOrderByMerchant(ctx, db.GetYearlyOrderByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, order_errors.ErrGetYearlyOrderByMerchant
	}

	return res, nil
}

func (r *orderRepository) FindById(ctx context.Context, user_id int) (*db.GetOrderByIDRow, error) {
	res, err := r.db.GetOrderByID(ctx, int32(user_id))

	if err != nil {
		return nil, order_errors.ErrFindById
	}

	return res, nil
}

func (r *orderRepository) FindByIdTrashed(ctx context.Context, user_id int) (*db.Order, error) {
	res, err := r.db.GetOrderByIDTrashed(ctx, int32(user_id))

	if err != nil {
		return nil, order_errors.ErrFindById
	}

	return res, nil
}

func (r *orderRepository) CreateOrder(ctx context.Context, request *requests.CreateOrderRecordRequest) (*db.CreateOrderRow, error) {
	req := db.CreateOrderParams{
		MerchantID: int32(request.MerchantID),
		UserID:     int32(request.UserID),
		TotalPrice: int32(request.TotalPrice),
	}

	user, err := r.db.CreateOrder(ctx, req)

	if err != nil {
		return nil, order_errors.ErrCreateOrder
	}

	return user, nil
}

func (r *orderRepository) UpdateOrder(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*db.UpdateOrderRow, error) {
	req := db.UpdateOrderParams{
		OrderID:    int32(request.OrderID),
		TotalPrice: int32(request.TotalPrice),
	}

	res, err := r.db.UpdateOrder(ctx, req)

	if err != nil {
		return nil, order_errors.ErrUpdateOrder
	}

	return res, nil
}

func (r *orderRepository) TrashedOrder(ctx context.Context, user_id int) (*db.Order, error) {
	res, err := r.db.TrashedOrder(ctx, int32(user_id))

	if err != nil {
		return nil, order_errors.ErrTrashedOrder
	}

	return res, nil
}

func (r *orderRepository) RestoreOrder(ctx context.Context, user_id int) (*db.Order, error) {
	res, err := r.db.RestoreOrder(ctx, int32(user_id))

	if err != nil {
		return nil, order_errors.ErrRestoreOrder
	}

	return res, nil
}

func (r *orderRepository) DeleteOrderPermanent(ctx context.Context, user_id int) (bool, error) {
	err := r.db.DeleteOrderPermanently(ctx, int32(user_id))

	if err != nil {
		return false, order_errors.ErrDeleteOrderPermanent
	}

	return true, nil
}

func (r *orderRepository) RestoreAllOrder(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllOrders(ctx)

	if err != nil {
		return false, order_errors.ErrRestoreAllOrder
	}
	return true, nil
}

func (r *orderRepository) DeleteAllOrderPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentOrders(ctx)

	if err != nil {
		return false, order_errors.ErrDeleteAllOrderPermanent
	}
	return true, nil
}
