package recordmapper

import (
	"ecommerce/internal/domain/record"
	db "ecommerce/pkg/database/schema"
)

type merchantPolicyRecordMapper struct {
}

func NewMerchantPolicyRecordMapper() *merchantPolicyRecordMapper {
	return &merchantPolicyRecordMapper{}
}

func (s *merchantPolicyRecordMapper) ToMerchantPolicyRecord(m *db.MerchantPolicy) *record.MerchantPoliciesRecord {
	var deletedAt *string
	if m.DeletedAt.Valid {
		str := m.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	return &record.MerchantPoliciesRecord{
		ID:          int(m.MerchantPolicyID),
		MerchantID:  int(m.MerchantPolicyID),
		PolicyType:  m.PolicyType,
		Title:       m.Title,
		Description: m.Description,
		CreatedAt:   m.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:   m.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:   deletedAt,
	}
}

func (s *merchantPolicyRecordMapper) ToMerchantPolicyRecordPagination(m *db.GetMerchantPoliciesRow) *record.MerchantPoliciesRecord {
	var deletedAt *string
	if m.DeletedAt.Valid {
		deletedAtStr := m.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.MerchantPoliciesRecord{
		ID:          int(m.MerchantPolicyID),
		MerchantID:  int(m.MerchantPolicyID),
		PolicyType:  m.PolicyType,
		Title:       m.Title,
		Description: m.Description,
		CreatedAt:   m.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:   m.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:   deletedAt,
	}
}

func (s *merchantPolicyRecordMapper) ToMerchantPolicysRecordPagination(MerchantPolicys []*db.GetMerchantPoliciesRow) []*record.MerchantPoliciesRecord {
	var result []*record.MerchantPoliciesRecord

	for _, MerchantPolicy := range MerchantPolicys {
		result = append(result, s.ToMerchantPolicyRecordPagination(MerchantPolicy))
	}

	return result
}

func (s *merchantPolicyRecordMapper) ToMerchantPolicyRecordActivePagination(m *db.GetMerchantPoliciesActiveRow) *record.MerchantPoliciesRecord {
	var deletedAt *string
	if m.DeletedAt.Valid {
		deletedAtStr := m.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.MerchantPoliciesRecord{
		ID:          int(m.MerchantPolicyID),
		MerchantID:  int(m.MerchantPolicyID),
		PolicyType:  m.PolicyType,
		Title:       m.Title,
		Description: m.Description,
		CreatedAt:   m.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:   m.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:   deletedAt,
	}
}

func (s *merchantPolicyRecordMapper) ToMerchantPolicysRecordActivePagination(MerchantPolicys []*db.GetMerchantPoliciesActiveRow) []*record.MerchantPoliciesRecord {
	var result []*record.MerchantPoliciesRecord

	for _, MerchantPolicy := range MerchantPolicys {
		result = append(result, s.ToMerchantPolicyRecordActivePagination(MerchantPolicy))
	}

	return result
}

func (s *merchantPolicyRecordMapper) ToMerchantPolicyRecordTrashedPagination(m *db.GetMerchantPoliciesTrashedRow) *record.MerchantPoliciesRecord {
	var deletedAt *string
	if m.DeletedAt.Valid {
		deletedAtStr := m.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.MerchantPoliciesRecord{
		ID:          int(m.MerchantPolicyID),
		MerchantID:  int(m.MerchantPolicyID),
		PolicyType:  m.PolicyType,
		Title:       m.Title,
		Description: m.Description,
		CreatedAt:   m.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:   m.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:   deletedAt,
	}
}

func (s *merchantPolicyRecordMapper) ToMerchantPolicysRecordTrashedPagination(MerchantPolicys []*db.GetMerchantPoliciesTrashedRow) []*record.MerchantPoliciesRecord {
	var result []*record.MerchantPoliciesRecord

	for _, MerchantPolicy := range MerchantPolicys {
		result = append(result, s.ToMerchantPolicyRecordTrashedPagination(MerchantPolicy))
	}

	return result
}
