package repository

import (
	"context"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/transaction_errors"
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

func (r *transactionRepository) FindAllTransactions(req *requests.FindAllTransaction) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactions(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindAllTransactions
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordPagination(res), &totalCount, nil
}

func (r *transactionRepository) FindByActive(req *requests.FindAllTransaction) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactionsActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindByActive
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordActivePagination(res), &totalCount, nil
}

func (r *transactionRepository) FindByTrashed(req *requests.FindAllTransaction) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactionsTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindByTrashed
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
		return nil, nil, transaction_errors.ErrFindByMerchant
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
		return nil, transaction_errors.ErrGetMonthlyAmountSuccess
	}

	return r.mapping.ToTransactionMonthlyAmountSuccess(res), nil
}

func (r *transactionRepository) GetYearlyAmountSuccess(year int) ([]*record.TransactionYearlyAmountSuccessRecord, error) {
	res, err := r.db.GetYearlyAmountTransactionSuccess(r.ctx, int32(year))

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountSuccess
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
		return nil, transaction_errors.ErrGetMonthlyAmountFailed
	}

	return r.mapping.ToTransactionMonthlyAmountFailed(res), nil
}

func (r *transactionRepository) GetYearlyAmountFailed(year int) ([]*record.TransactionYearlyAmountFailedRecord, error) {
	res, err := r.db.GetYearlyAmountTransactionFailed(r.ctx, int32(year))

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountFailed
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
		return nil, transaction_errors.ErrGetMonthlyAmountSuccessByMerchant
	}

	return r.mapping.ToTransactionMonthlyAmountSuccessByMerchant(res), nil
}

func (r *transactionRepository) GetYearlyAmountSuccessByMerchant(req *requests.YearAmountTransactionMerchant) ([]*record.TransactionYearlyAmountSuccessRecord, error) {
	res, err := r.db.GetYearlyAmountTransactionSuccessByMerchant(r.ctx, db.GetYearlyAmountTransactionSuccessByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountSuccessByMerchant
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
		return nil, transaction_errors.ErrGetMonthlyAmountFailedByMerchant
	}

	return r.mapping.ToTransactionMonthlyAmountFailedByMerchant(res), nil
}

func (r *transactionRepository) GetYearlyAmountFailedByMerchant(req *requests.YearAmountTransactionMerchant) ([]*record.TransactionYearlyAmountFailedRecord, error) {
	res, err := r.db.GetYearlyAmountTransactionFailedByMerchant(r.ctx, db.GetYearlyAmountTransactionFailedByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountFailedByMerchant
	}

	return r.mapping.ToTransactionYearlyAmountFailedByMerchant(res), nil
}

func (r *transactionRepository) GetMonthlyTransactionMethodSuccess(req *requests.MonthMethodTransaction) ([]*record.TransactionMonthlyMethodRecord, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTransactionMethodsSuccess(r.ctx, db.GetMonthlyTransactionMethodsSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethod
	}

	return r.mapping.ToTransactionMonthlyMethodSuccess(res), nil
}

func (r *transactionRepository) GetYearlyTransactionMethodSuccess(year int) ([]*record.TransactionYearlyMethodRecord, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyTransactionMethodsSuccess(r.ctx, yearStart)

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethod
	}

	return r.mapping.ToTransactionYearlyMethodSuccess(res), nil
}

func (r *transactionRepository) GetMonthlyTransactionMethodFailed(req *requests.MonthMethodTransaction) ([]*record.TransactionMonthlyMethodRecord, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTransactionMethodsFailed(r.ctx, db.GetMonthlyTransactionMethodsFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethod
	}

	return r.mapping.ToTransactionMonthlyMethodFailed(res), nil
}

func (r *transactionRepository) GetYearlyTransactionMethodFailed(year int) ([]*record.TransactionYearlyMethodRecord, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyTransactionMethodsFailed(r.ctx, yearStart)

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethod
	}

	return r.mapping.ToTransactionYearlyMethodFailed(res), nil
}

func (r *transactionRepository) GetMonthlyTransactionMethodByMerchantSuccess(req *requests.MonthMethodTransactionMerchant) ([]*record.TransactionMonthlyMethodRecord, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTransactionMethodsByMerchantSuccess(r.ctx, db.GetMonthlyTransactionMethodsByMerchantSuccessParams{
		Column1:    currentDate,
		Column2:    lastDayCurrentMonth,
		Column3:    prevDate,
		Column4:    lastDayPrevMonth,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethodByMerchant
	}

	return r.mapping.ToTransactionMonthlyByMerchantMethodSuccess(res), nil
}

func (r *transactionRepository) GetYearlyTransactionMethodByMerchantSuccess(req *requests.YearMethodTransactionMerchant) ([]*record.TransactionYearlyMethodRecord, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyTransactionMethodsByMerchantSuccess(r.ctx, db.GetYearlyTransactionMethodsByMerchantSuccessParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethodByMerchant
	}

	return r.mapping.ToTransactionYearlyMethodByMerchantSuccess(res), nil
}

func (r *transactionRepository) GetMonthlyTransactionMethodByMerchantFailed(req *requests.MonthMethodTransactionMerchant) ([]*record.TransactionMonthlyMethodRecord, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTransactionMethodsByMerchantFailed(r.ctx, db.GetMonthlyTransactionMethodsByMerchantFailedParams{
		Column1:    currentDate,
		Column2:    lastDayCurrentMonth,
		Column3:    prevDate,
		Column4:    lastDayPrevMonth,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethodByMerchant
	}

	return r.mapping.ToTransactionMonthlyByMerchantMethodFailed(res), nil
}

func (r *transactionRepository) GetYearlyTransactionMethodByMerchantFailed(req *requests.YearMethodTransactionMerchant) ([]*record.TransactionYearlyMethodRecord, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyTransactionMethodsByMerchantFailed(r.ctx, db.GetYearlyTransactionMethodsByMerchantFailedParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethodByMerchant
	}

	return r.mapping.ToTransactionYearlyMethodByMerchantFailed(res), nil
}

func (r *transactionRepository) FindById(transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.GetTransactionByID(r.ctx, int32(transaction_id))

	if err != nil {
		return nil, transaction_errors.ErrFindById
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) FindByOrderId(order_id int) (*record.TransactionRecord, error) {
	res, err := r.db.GetTransactionByOrderID(r.ctx, int32(order_id))

	if err != nil {
		return nil, transaction_errors.ErrFindByOrderId
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
		return nil, transaction_errors.ErrCreateTransaction
	}

	fmt.Println("transaction", transaction)

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
		return nil, transaction_errors.ErrUpdateTransaction
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) TrashTransaction(transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.TrashTransaction(r.ctx, int32(transaction_id))

	if err != nil {
		return nil, transaction_errors.ErrTrashTransaction
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) RestoreTransaction(transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.RestoreTransaction(r.ctx, int32(transaction_id))

	if err != nil {
		return nil, transaction_errors.ErrRestoreTransaction
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionRepository) DeleteTransactionPermanently(transaction_id int) (bool, error) {
	err := r.db.DeleteTransactionPermanently(r.ctx, int32(transaction_id))

	if err != nil {
		return false, transaction_errors.ErrDeleteTransactionPermanently
	}

	return true, nil
}

func (r *transactionRepository) RestoreAllTransactions() (bool, error) {
	err := r.db.RestoreAllTransactions(r.ctx)

	if err != nil {
		return false, transaction_errors.ErrRestoreAllTransactions
	}
	return true, nil
}

func (r *transactionRepository) DeleteAllTransactionPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentTransactions(r.ctx)

	if err != nil {
		return false, transaction_errors.ErrDeleteAllTransactionPermanent
	}
	return true, nil
}
