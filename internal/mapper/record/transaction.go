package recordmapper

import (
	"ecommerce/internal/domain/record"
	db "ecommerce/pkg/database/schema"
)

type transactionRecordMapper struct {
}

func NewTransactionRecordMapper() *transactionRecordMapper {
	return &transactionRecordMapper{}
}

func (s *transactionRecordMapper) ToTransactionRecord(transaction *db.Transaction) *record.TransactionRecord {
	var deletedAt *string
	if transaction.DeletedAt.Valid {
		deletedAtStr := transaction.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.TransactionRecord{
		ID:            int(transaction.TransactionID),
		OrderID:       int(transaction.OrderID),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int(transaction.Amount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:     transaction.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:     deletedAt,
	}
}

func (s *transactionRecordMapper) ToTransactionsRecord(transactions []*db.Transaction) []*record.TransactionRecord {
	var result []*record.TransactionRecord

	for _, transaction := range transactions {
		result = append(result, s.ToTransactionRecord(transaction))
	}

	return result
}

func (s *transactionRecordMapper) ToTransactionMerchantRecordPagination(transaction *db.GetTransactionByMerchantRow) *record.TransactionRecord {
	var deletedAt *string
	if transaction.DeletedAt.Valid {
		deletedAtStr := transaction.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.TransactionRecord{
		ID:            int(transaction.TransactionID),
		OrderID:       int(transaction.OrderID),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int(transaction.Amount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:     transaction.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:     deletedAt,
	}
}

func (s *transactionRecordMapper) ToTransactionMerchantsRecordPagination(products []*db.GetTransactionByMerchantRow) []*record.TransactionRecord {
	var result []*record.TransactionRecord

	for _, product := range products {
		result = append(result, s.ToTransactionMerchantRecordPagination(product))
	}

	return result
}

func (s *transactionRecordMapper) ToTransactionRecordPagination(transaction *db.GetTransactionsRow) *record.TransactionRecord {
	var deletedAt *string
	if transaction.DeletedAt.Valid {
		deletedAtStr := transaction.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.TransactionRecord{
		ID:            int(transaction.TransactionID),
		OrderID:       int(transaction.OrderID),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int(transaction.Amount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:     transaction.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:     deletedAt,
	}
}

func (s *transactionRecordMapper) ToTransactionsRecordPagination(products []*db.GetTransactionsRow) []*record.TransactionRecord {
	var result []*record.TransactionRecord

	for _, product := range products {
		result = append(result, s.ToTransactionRecordPagination(product))
	}

	return result
}

func (s *transactionRecordMapper) ToTransactionRecordActivePagination(transaction *db.GetTransactionsActiveRow) *record.TransactionRecord {
	var deletedAt *string
	if transaction.DeletedAt.Valid {
		deletedAtStr := transaction.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.TransactionRecord{
		ID:            int(transaction.TransactionID),
		OrderID:       int(transaction.OrderID),
		MerchantID:    int(transaction.MerchantID),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int(transaction.Amount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:     transaction.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:     deletedAt,
	}
}

func (s *transactionRecordMapper) ToTransactionsRecordActivePagination(transactions []*db.GetTransactionsActiveRow) []*record.TransactionRecord {
	var result []*record.TransactionRecord

	for _, transaction := range transactions {
		result = append(result, s.ToTransactionRecordActivePagination(transaction))
	}

	return result
}

func (s *transactionRecordMapper) ToTransactionRecordTrashedPagination(transaction *db.GetTransactionsTrashedRow) *record.TransactionRecord {
	var deletedAt *string
	if transaction.DeletedAt.Valid {
		deletedAtStr := transaction.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.TransactionRecord{
		ID:            int(transaction.TransactionID),
		OrderID:       int(transaction.OrderID),
		MerchantID:    int(transaction.MerchantID),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int(transaction.Amount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:     transaction.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:     deletedAt,
	}
}

func (s *transactionRecordMapper) ToTransactionsRecordTrashedPagination(products []*db.GetTransactionsTrashedRow) []*record.TransactionRecord {
	var result []*record.TransactionRecord

	for _, product := range products {
		result = append(result, s.ToTransactionRecordTrashedPagination(product))
	}

	return result
}
