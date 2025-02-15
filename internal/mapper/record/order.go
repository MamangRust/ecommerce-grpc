package recordmapper

import (
	"ecommerce/internal/domain/record"
	db "ecommerce/pkg/database/schema"
)

type orderRecordMapper struct {
}

func NewOrderRecordMapper() *orderRecordMapper {
	return &orderRecordMapper{}
}

func (s *orderRecordMapper) ToOrderRecord(order *db.Order) *record.OrderRecord {
	var deletedAt *string
	if order.DeletedAt.Valid {
		deletedAtStr := order.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderRecord{
		ID:         int(order.OrderID),
		MerchantID: int(order.MerchantID),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *orderRecordMapper) ToOrdersRecord(orders []*db.Order) []*record.OrderRecord {
	var result []*record.OrderRecord

	for _, order := range orders {
		result = append(result, s.ToOrderRecord(order))
	}

	return result
}

func (s *orderRecordMapper) ToOrderRecordPagination(order *db.GetOrdersRow) *record.OrderRecord {
	var deletedAt *string
	if order.DeletedAt.Valid {
		deletedAtStr := order.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderRecord{
		ID:         int(order.OrderID),
		MerchantID: int(order.MerchantID),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *orderRecordMapper) ToOrdersRecordPagination(orders []*db.GetOrdersRow) []*record.OrderRecord {
	var result []*record.OrderRecord

	for _, order := range orders {
		result = append(result, s.ToOrderRecordPagination(order))
	}

	return result
}

func (s *orderRecordMapper) ToOrderRecordActivePagination(order *db.GetOrdersActiveRow) *record.OrderRecord {
	var deletedAt *string
	if order.DeletedAt.Valid {
		deletedAtStr := order.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderRecord{
		ID:         int(order.OrderID),
		MerchantID: int(order.MerchantID),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *orderRecordMapper) ToOrdersRecordActivePagination(orders []*db.GetOrdersActiveRow) []*record.OrderRecord {
	var result []*record.OrderRecord

	for _, order := range orders {
		result = append(result, s.ToOrderRecordActivePagination(order))
	}

	return result
}

func (s *orderRecordMapper) ToOrderRecordTrashedPagination(order *db.GetOrdersTrashedRow) *record.OrderRecord {
	var deletedAt *string
	if order.DeletedAt.Valid {
		deletedAtStr := order.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderRecord{
		ID:         int(order.OrderID),
		MerchantID: int(order.MerchantID),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *orderRecordMapper) ToOrdersRecordTrashedPagination(orders []*db.GetOrdersTrashedRow) []*record.OrderRecord {
	var result []*record.OrderRecord

	for _, order := range orders {
		result = append(result, s.ToOrderRecordTrashedPagination(order))
	}

	return result
}

func (s *orderRecordMapper) ToOrderRecordByMerchantPagination(order *db.GetOrdersByMerchantRow) *record.OrderRecord {
	var deletedAt *string
	if order.DeletedAt.Valid {
		deletedAtStr := order.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderRecord{
		ID:         int(order.OrderID),
		MerchantID: int(order.MerchantID),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *orderRecordMapper) ToOrdersRecordByMerchantPagination(orders []*db.GetOrdersByMerchantRow) []*record.OrderRecord {
	var result []*record.OrderRecord

	for _, order := range orders {
		result = append(result, s.ToOrderRecordByMerchantPagination(order))
	}

	return result
}
