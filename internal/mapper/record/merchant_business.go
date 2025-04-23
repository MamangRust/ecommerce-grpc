package recordmapper

import (
	"ecommerce/internal/domain/record"
	db "ecommerce/pkg/database/schema"
)

type merchantBusinessRecordMapper struct {
}

func NewMerchantBusinessRecordMapper() *merchantBusinessRecordMapper {
	return &merchantBusinessRecordMapper{}
}

func (s *merchantBusinessRecordMapper) ToMerchantBusinessRecord(m *db.MerchantBusinessInformation) *record.MerchantBusinessRecord {
	var deletedAt *string
	if m.DeletedAt.Valid {
		str := m.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	return &record.MerchantBusinessRecord{
		ID:                int(m.MerchantBusinessInfoID),
		MerchantID:        int(m.MerchantID),
		BusinessType:      m.BusinessType.String,
		TaxID:             m.TaxID.String,
		EstablishedYear:   int(m.EstablishedYear.Int32),
		NumberOfEmployees: int(m.NumberOfEmployees.Int32),
		WebsiteUrl:        m.WebsiteUrl.String,
		CreatedAt:         m.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:         m.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:         deletedAt,
	}
}

func (s *merchantBusinessRecordMapper) ToMerchantBusinessRecordPagination(m *db.GetMerchantsBusinessInformationRow) *record.MerchantBusinessRecord {
	var deletedAt *string
	if m.DeletedAt.Valid {
		str := m.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	return &record.MerchantBusinessRecord{
		ID:                int(m.MerchantBusinessInfoID),
		MerchantID:        int(m.MerchantID),
		BusinessType:      m.BusinessType.String,
		TaxID:             m.TaxID.String,
		EstablishedYear:   int(m.EstablishedYear.Int32),
		NumberOfEmployees: int(m.NumberOfEmployees.Int32),
		WebsiteUrl:        m.WebsiteUrl.String,
		MerchantName:      &m.MerchantName,
		CreatedAt:         m.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:         m.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:         deletedAt,
	}
}

func (s *merchantBusinessRecordMapper) ToMerchantBusinesssRecordPagination(Merchants []*db.GetMerchantsBusinessInformationRow) []*record.MerchantBusinessRecord {
	var result []*record.MerchantBusinessRecord

	for _, Merchant := range Merchants {
		result = append(result, s.ToMerchantBusinessRecordPagination(Merchant))
	}

	return result
}

func (s *merchantBusinessRecordMapper) ToMerchantBusinessRecordActivePagination(m *db.GetMerchantsBusinessInformationActiveRow) *record.MerchantBusinessRecord {
	var deletedAt *string
	if m.DeletedAt.Valid {
		str := m.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	return &record.MerchantBusinessRecord{
		ID:                int(m.MerchantBusinessInfoID),
		MerchantID:        int(m.MerchantID),
		BusinessType:      m.BusinessType.String,
		TaxID:             m.TaxID.String,
		EstablishedYear:   int(m.EstablishedYear.Int32),
		NumberOfEmployees: int(m.NumberOfEmployees.Int32),
		WebsiteUrl:        m.WebsiteUrl.String,
		MerchantName:      &m.MerchantName,
		CreatedAt:         m.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:         m.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:         deletedAt,
	}
}

func (s *merchantBusinessRecordMapper) ToMerchantBusinesssRecordActivePagination(Merchants []*db.GetMerchantsBusinessInformationActiveRow) []*record.MerchantBusinessRecord {
	var result []*record.MerchantBusinessRecord

	for _, Merchant := range Merchants {
		result = append(result, s.ToMerchantBusinessRecordActivePagination(Merchant))
	}

	return result
}

func (s *merchantBusinessRecordMapper) ToMerchantBusinessRecordTrashedPagination(m *db.GetMerchantsBusinessInformationTrashedRow) *record.MerchantBusinessRecord {
	var deletedAt *string
	if m.DeletedAt.Valid {
		str := m.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	return &record.MerchantBusinessRecord{
		ID:                int(m.MerchantBusinessInfoID),
		MerchantID:        int(m.MerchantID),
		BusinessType:      m.BusinessType.String,
		TaxID:             m.TaxID.String,
		EstablishedYear:   int(m.EstablishedYear.Int32),
		NumberOfEmployees: int(m.NumberOfEmployees.Int32),
		WebsiteUrl:        m.WebsiteUrl.String,
		MerchantName:      &m.MerchantName,
		CreatedAt:         m.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:         m.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:         deletedAt,
	}
}

func (s *merchantBusinessRecordMapper) ToMerchantBusinesssRecordTrashedPagination(Merchants []*db.GetMerchantsBusinessInformationTrashedRow) []*record.MerchantBusinessRecord {
	var result []*record.MerchantBusinessRecord

	for _, Merchant := range Merchants {
		result = append(result, s.ToMerchantBusinessRecordTrashedPagination(Merchant))
	}

	return result
}
