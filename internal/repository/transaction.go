package repository

import (
	"context"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"fmt"
	"time"
)

type transactionRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TransactionRecordMapping
}

func NewTransactionRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TransactionRecordMapping) *transactionRepository {
	return &transactionRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *transactionRepository) FindAllTransactions(req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactions(r.ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to find transactions: %w", err)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordPagination(res), &totalCount, nil
}

func (r *transactionRepository) FindByActive(req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactionsActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to find active transactions: %w", err)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordActivePagination(res), &totalCount, nil
}

func (r *transactionRepository) FindByTrashed(req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactionsTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to find trashed transactions: %w", err)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordTrashedPagination(res), &totalCount, nil
}

func (r *transactionRepository) FindByMerchant(
	req *requests.FindAllTransactionByMerchant,
) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionByMerchantParams{
		Column1: req.Search,
		Column2: int32(req.MerchantID),
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactionByMerchant(r.ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to find merchant transactions: %w", err)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionMerchantsRecordPagination(res), &totalCount, nil
}

func (r *transactionRepository) GetMonthlyAmountSuccess(req *requests.MonthAmountTransaction) ([]*record.TransactionMonthlyAmountSuccessRecord, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyAmountTransactionSuccess(r.ctx, db.GetMonthlyAmountTransactionSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get successful transactions amount: no data available for %d-%02d", req.Year, req.Month)
	}

	return r.mapping.ToTransactionMonthlyAmountSuccess(res), nil
}

func (r *transactionRepository) GetYearlyAmountSuccess(year int) ([]*record.TransactionYearlyAmountSuccessRecord, error) {
	res, err := r.db.GetYearlyAmountTransactionSuccess(r.ctx, int32(year))

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly successful transactions amount: no data for year %d", year)
	}

	return r.mapping.ToTransactionYearlyAmountSuccess(res), nil
}

func (r *transactionRepository) GetMonthlyAmountFailed(req *requests.MonthAmountTransaction) ([]*record.TransactionMonthlyAmountFailedRecord, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyAmountTransactionFailed(r.ctx, db.GetMonthlyAmountTransactionFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get failed transactions amount: no data available for %d-%02d", req.Year, req.Month)
	}

	return r.mapping.ToTransactionMonthlyAmountFailed(res), nil
}

func (r *transactionRepository) GetYearlyAmountFailed(year int) ([]*record.TransactionYearlyAmountFailedRecord, error) {
	res, err := r.db.GetYearlyAmountTransactionFailed(r.ctx, int32(year))

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly failed transactions amount: no data for year %d", year)
	}

	return r.mapping.ToTransactionYearlyAmountFailed(res), nil
}

func (r *transactionRepository) GetMonthlyAmountSuccessByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*record.TransactionMonthlyAmountSuccessRecord, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyAmountTransactionSuccessByMerchant(r.ctx, db.GetMonthlyAmountTransactionSuccessByMerchantParams{
		Column1:    currentDate,
		Column2:    lastDayCurrentMonth,
		Column3:    prevDate,
		Column4:    lastDayPrevMonth,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get merchant %d successful transactions: no data for %d-%02d", req.MerchantID, req.Year, req.Month)
	}

	return r.mapping.ToTransactionMonthlyAmountSuccessByMerchant(res), nil
}

func (r *transactionRepository) GetYearlyAmountSuccessByMerchant(req *requests.YearAmountTransactionMerchant) ([]*record.TransactionYearlyAmountSuccessRecord, error) {
	res, err := r.db.GetYearlyAmountTransactionSuccessByMerchant(r.ctx, db.GetYearlyAmountTransactionSuccessByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get merchant %d yearly successful transactions: no data for year %d", req.MerchantID, req.Year)
	}

	return r.mapping.ToTransactionYearlyAmountSuccessByMerchant(res), nil
}

func (r *transactionRepository) GetMonthlyAmountFailedByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*record.TransactionMonthlyAmountFailedRecord, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyAmountTransactionFailedByMerchant(r.ctx, db.GetMonthlyAmountTransactionFailedByMerchantParams{
		Column1:    currentDate,
		Column2:    lastDayCurrentMonth,
		Column3:    prevDate,
		Column4:    lastDayPrevMonth,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get merchant %d failed transactions: no data for %d-%02d", req.MerchantID, req.Year, req.Month)
	}

	return r.mapping.ToTransactionMonthlyAmountFailedByMerchant(res), nil
}

func (r *transactionRepository) GetYearlyAmountFailedByMerchant(req *requests.YearAmountTransactionMerchant) ([]*record.TransactionYearlyAmountFailedRecord, error) {
	res, err := r.db.GetYearlyAmountTransactionFailedByMerchant(r.ctx, db.GetYearlyAmountTransactionFailedByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get merchant %d yearly failed transactions: no data for year %d", req.MerchantID, req.Year)
	}

	return r.mapping.ToTransactionYearlyAmountFailedByMerchant(res), nil
}

func (r *transactionRepository) GetMonthlyTransactionMethod(year int) ([]*record.TransactionMonthlyMethodRecord, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTransactionMethods(r.ctx, yearStart)

	if err != nil {
		return nil, fmt.Errorf("failed to get monthly payment methods: no transaction data for year %d", year)
	}

	return r.mapping.ToTransactionMonthlyMethod(res), nil
}

func (r *transactionRepository) GetYearlyTransactionMethod(year int) ([]*record.TransactionYearlyMethodRecord, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyTransactionMethods(r.ctx, yearStart)

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly payment methods: no transaction data for year %d", year)
	}

	return r.mapping.ToTransactionYearlyMethod(res), nil
}

func (r *transactionRepository) GetMonthlyTransactionMethodByMerchant(req *requests.MonthlyYearTransactionMethodMerchant) ([]*record.TransactionMonthlyMethodRecord, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTransactionMethodsByMerchant(r.ctx, db.GetMonthlyTransactionMethodsByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get merchant %d payment methods: no monthly data for year %d", req.MerchantID, req.Year)
	}

	return r.mapping.ToTransactionMonthlyByMerchantMethod(res), nil
}

func (r *transactionRepository) GetYearlyTransactionMethodByMerchant(req *requests.MonthlyYearTransactionMethodMerchant) ([]*record.TransactionYearlyMethodRecord, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyTransactionMethodsByMerchant(r.ctx, db.GetYearlyTransactionMethodsByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get merchant %d yearly payment methods: no data for year %d", req.MerchantID, req.Year)
	}

	return r.mapping.ToTransactionYearlyMethodByMerchant(res), nil
}

func (r *transactionRepository) FindById(transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.GetTransactionByID(r.ctx, int32(transaction_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) FindByOrderId(order_id int) (*record.TransactionRecord, error) {
	res, err := r.db.GetTransactionByOrderID(r.ctx, int32(order_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) CreateTransaction(request *requests.CreateTransactionRequest) (*record.TransactionRecord, error) {
	req := db.CreateTransactionParams{
		OrderID:       int32(request.OrderID),
		MerchantID:    int32(request.MerchantID),
		PaymentMethod: request.PaymentMethod,
		Amount:        int32(request.Amount),
		PaymentStatus: *request.PaymentStatus,
	}

	transaction, err := r.db.CreateTransaction(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return r.mapping.ToTransactionRecord(transaction), nil
}

func (r *transactionRepository) UpdateTransaction(request *requests.UpdateTransactionRequest) (*record.TransactionRecord, error) {
	req := db.UpdateTransactionParams{
		TransactionID: int32(*request.TransactionID),
		MerchantID:    int32(request.MerchantID),
		PaymentMethod: request.PaymentMethod,
		Amount:        int32(request.Amount),
		OrderID:       int32(request.OrderID),
		PaymentStatus: *request.PaymentStatus,
	}

	res, err := r.db.UpdateTransaction(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) TrashTransaction(transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.TrashTransaction(r.ctx, int32(transaction_id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash user: %w", err)
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) RestoreTransaction(transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.RestoreTransaction(r.ctx, int32(transaction_id))

	if err != nil {
		return nil, fmt.Errorf("failed to restore topup: %w", err)
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) DeleteTransactionPermanently(transaction_id int) (bool, error) {
	err := r.db.DeleteTransactionPermanently(r.ctx, int32(transaction_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete transactions: %w", err)
	}

	return true, nil
}

func (r *transactionRepository) RestoreAllTransactions() (bool, error) {
	err := r.db.RestoreAllTransactions(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all transactions: %w", err)
	}
	return true, nil
}

func (r *transactionRepository) DeleteAllTransactionPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentTransactions(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all transactions permanently: %w", err)
	}
	return true, nil
}
