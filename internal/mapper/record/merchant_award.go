package recordmapper

import (
	"ecommerce/internal/domain/record"
	db "ecommerce/pkg/database/schema"
)

type merchantAwardRecordMapper struct {
}

func NewMerchantAwardRecordMapper() *merchantAwardRecordMapper {
	return &merchantAwardRecordMapper{}
}

func (s *merchantAwardRecordMapper) ToMerchantAwardRecord(m *db.MerchantCertificationsAndAward) *record.MerchantAwardRecord {
	var deletedAt *string
	if m.DeletedAt.Valid {
		str := m.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	issueDate := ""
	if m.IssueDate.Valid {
		issueDate = m.IssueDate.Time.Format("2006-01-02")
	}

	expiryDate := ""
	if m.ExpiryDate.Valid {
		expiryDate = m.ExpiryDate.Time.Format("2006-01-02")
	}

	return &record.MerchantAwardRecord{
		ID:             int(m.MerchantCertificationID),
		MerchantID:     int(m.MerchantID),
		Title:          m.Title,
		Description:    m.Description.String,
		IssuedBy:       m.IssuedBy.String,
		IssueDate:      issueDate,
		ExpiryDate:     expiryDate,
		CertificateUrl: m.CertificateUrl.String,
		CreatedAt:      m.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:      m.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:      deletedAt,
	}
}

func (s *merchantAwardRecordMapper) ToMerchantAwardRecordPagination(m *db.GetMerchantCertificationsAndAwardsRow) *record.MerchantAwardRecord {
	var deletedAt *string
	if m.DeletedAt.Valid {
		str := m.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	issueDate := ""
	if m.IssueDate.Valid {
		issueDate = m.IssueDate.Time.Format("2006-01-02")
	}

	expiryDate := ""
	if m.ExpiryDate.Valid {
		expiryDate = m.ExpiryDate.Time.Format("2006-01-02")
	}

	return &record.MerchantAwardRecord{
		ID:             int(m.MerchantCertificationID),
		MerchantID:     int(m.MerchantID),
		Title:          m.Title,
		Description:    m.Description.String,
		IssuedBy:       m.IssuedBy.String,
		IssueDate:      issueDate,
		ExpiryDate:     expiryDate,
		CertificateUrl: m.CertificateUrl.String,
		CreatedAt:      m.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:      m.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:      deletedAt,
	}
}

func (s *merchantAwardRecordMapper) ToMerchantAwardsRecordPagination(Merchants []*db.GetMerchantCertificationsAndAwardsRow) []*record.MerchantAwardRecord {
	var result []*record.MerchantAwardRecord

	for _, Merchant := range Merchants {
		result = append(result, s.ToMerchantAwardRecordPagination(Merchant))
	}

	return result
}

func (s *merchantAwardRecordMapper) ToMerchantAwardRecordActivePagination(m *db.GetMerchantCertificationsAndAwardsActiveRow) *record.MerchantAwardRecord {
	var deletedAt *string
	if m.DeletedAt.Valid {
		str := m.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	issueDate := ""
	if m.IssueDate.Valid {
		issueDate = m.IssueDate.Time.Format("2006-01-02")
	}

	expiryDate := ""
	if m.ExpiryDate.Valid {
		expiryDate = m.ExpiryDate.Time.Format("2006-01-02")
	}

	return &record.MerchantAwardRecord{
		ID:             int(m.MerchantCertificationID),
		MerchantID:     int(m.MerchantID),
		Title:          m.Title,
		Description:    m.Description.String,
		IssuedBy:       m.IssuedBy.String,
		IssueDate:      issueDate,
		ExpiryDate:     expiryDate,
		CertificateUrl: m.CertificateUrl.String,
		CreatedAt:      m.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:      m.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:      deletedAt,
	}
}

func (s *merchantAwardRecordMapper) ToMerchantAwardsRecordActivePagination(Merchants []*db.GetMerchantCertificationsAndAwardsActiveRow) []*record.MerchantAwardRecord {
	var result []*record.MerchantAwardRecord

	for _, Merchant := range Merchants {
		result = append(result, s.ToMerchantAwardRecordActivePagination(Merchant))
	}

	return result
}

func (s *merchantAwardRecordMapper) ToMerchantAwardRecordTrashedPagination(m *db.GetMerchantCertificationsAndAwardsTrashedRow) *record.MerchantAwardRecord {
	var deletedAt *string
	if m.DeletedAt.Valid {
		str := m.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	issueDate := ""
	if m.IssueDate.Valid {
		issueDate = m.IssueDate.Time.Format("2006-01-02")
	}

	expiryDate := ""
	if m.ExpiryDate.Valid {
		expiryDate = m.ExpiryDate.Time.Format("2006-01-02")
	}

	return &record.MerchantAwardRecord{
		ID:             int(m.MerchantCertificationID),
		MerchantID:     int(m.MerchantID),
		Title:          m.Title,
		Description:    m.Description.String,
		IssuedBy:       m.IssuedBy.String,
		IssueDate:      issueDate,
		ExpiryDate:     expiryDate,
		CertificateUrl: m.CertificateUrl.String,
		CreatedAt:      m.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:      m.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:      deletedAt,
	}
}

func (s *merchantAwardRecordMapper) ToMerchantAwardsRecordTrashedPagination(m []*db.GetMerchantCertificationsAndAwardsTrashedRow) []*record.MerchantAwardRecord {
	var result []*record.MerchantAwardRecord

	for _, Merchant := range m {
		result = append(result, s.ToMerchantAwardRecordTrashedPagination(Merchant))
	}

	return result
}
