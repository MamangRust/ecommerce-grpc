package recordmapper

import (
	"ecommerce/internal/domain/record"
	db "ecommerce/pkg/database/schema"
)

type shippingAddressRecordMapper struct {
}

func NewShippingAddressRecordMapper() *shippingAddressRecordMapper {
	return &shippingAddressRecordMapper{}
}

func (s *shippingAddressRecordMapper) ToShippingAddressRecord(address *db.ShippingAddress) *record.ShippingAddressRecord {
	var deletedAt *string

	if address.DeletedAt.Valid {
		formattedDeletedAt := address.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formattedDeletedAt
	}

	return &record.ShippingAddressRecord{
		ID:             int(address.ShippingAddressID),
		OrderID:        int(address.OrderID),
		Alamat:         address.Alamat,
		Provinsi:       address.Provinsi,
		Negara:         address.Negara,
		Kota:           address.Kota,
		ShippingMethod: address.ShippingMethod,
		ShippingCost:   int(address.ShippingCost),
		CreatedAt:      address.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      address.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:      deletedAt,
	}
}

func (s *shippingAddressRecordMapper) ToShippingAddressRecordPagination(address *db.GetShippingAddressRow) *record.ShippingAddressRecord {
	var deletedAt *string

	if address.DeletedAt.Valid {
		formattedDeletedAt := address.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formattedDeletedAt
	}

	return &record.ShippingAddressRecord{
		ID:             int(address.ShippingAddressID),
		OrderID:        int(address.OrderID),
		Alamat:         address.Alamat,
		Provinsi:       address.Provinsi,
		Negara:         address.Negara,
		Kota:           address.Kota,
		ShippingMethod: address.ShippingMethod,
		ShippingCost:   int(address.ShippingCost),
		CreatedAt:      address.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      address.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:      deletedAt,
	}
}

func (s *shippingAddressRecordMapper) ToShippingAddresssRecordPagination(addresses []*db.GetShippingAddressRow) []*record.ShippingAddressRecord {
	var addressRecords []*record.ShippingAddressRecord

	for _, address := range addresses {
		addressRecords = append(addressRecords, s.ToShippingAddressRecordPagination(address))
	}

	return addressRecords
}

func (s *shippingAddressRecordMapper) ToShippingAddressRecordActivePagination(address *db.GetShippingAddressActiveRow) *record.ShippingAddressRecord {
	var deletedAt *string
	if address.DeletedAt.Valid {
		deletedAtStr := address.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.ShippingAddressRecord{
		ID:             int(address.ShippingAddressID),
		OrderID:        int(address.OrderID),
		Alamat:         address.Alamat,
		Provinsi:       address.Provinsi,
		Negara:         address.Negara,
		Kota:           address.Kota,
		ShippingMethod: address.ShippingMethod,
		ShippingCost:   int(address.ShippingCost),
		CreatedAt:      address.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      address.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:      deletedAt,
	}
}

func (s *shippingAddressRecordMapper) ToShippingAddresssRecordActivePagination(sliders []*db.GetShippingAddressActiveRow) []*record.ShippingAddressRecord {
	var result []*record.ShippingAddressRecord

	for _, slider := range sliders {
		result = append(result, s.ToShippingAddressRecordActivePagination(slider))
	}

	return result
}

func (s *shippingAddressRecordMapper) ToShippingAddressRecordTrashedPagination(address *db.GetShippingAddressTrashedRow) *record.ShippingAddressRecord {
	var deletedAt *string
	if address.DeletedAt.Valid {
		deletedAtStr := address.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.ShippingAddressRecord{
		ID:             int(address.ShippingAddressID),
		OrderID:        int(address.OrderID),
		Alamat:         address.Alamat,
		Provinsi:       address.Provinsi,
		Negara:         address.Negara,
		Kota:           address.Kota,
		ShippingMethod: address.ShippingMethod,
		ShippingCost:   int(address.ShippingCost),
		CreatedAt:      address.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:      address.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:      deletedAt,
	}
}

func (s *shippingAddressRecordMapper) ToShippingAddresssRecordTrashedPagination(sliders []*db.GetShippingAddressTrashedRow) []*record.ShippingAddressRecord {
	var result []*record.ShippingAddressRecord

	for _, slider := range sliders {
		result = append(result, s.ToShippingAddressRecordTrashedPagination(slider))
	}

	return result
}
