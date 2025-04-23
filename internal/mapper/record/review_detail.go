package recordmapper

import (
	"ecommerce/internal/domain/record"
	db "ecommerce/pkg/database/schema"
)

type reviewDetailMapper struct {
}

func NewReviewDetailRecordMapper() *reviewDetailMapper {
	return &reviewDetailMapper{}
}

func (s *reviewDetailMapper) ToReviewDetailRecord(detail *db.ReviewDetail) *record.ReviewDetailRecord {
	var deletedAt *string
	if detail.DeletedAt.Valid {
		str := detail.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	reviewId := int(detail.ReviewID)

	return &record.ReviewDetailRecord{
		ID:        int(detail.ReviewDetailID),
		ReviewID:  &reviewId,
		Type:      detail.Type,
		Url:       detail.Url,
		Caption:   detail.Caption.String,
		CreatedAt: detail.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: detail.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *reviewDetailMapper) ToReviewDetailRecordPagination(detail *db.GetReviewDetailsRow) *record.ReviewDetailRecord {
	var deletedAt *string
	if detail.DeletedAt.Valid {
		str := detail.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	reviewId := int(detail.ReviewID)

	return &record.ReviewDetailRecord{
		ID:        int(detail.ReviewDetailID),
		ReviewID:  &reviewId,
		Type:      detail.Type,
		Url:       detail.Url,
		Caption:   detail.Caption.String,
		CreatedAt: detail.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: detail.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *reviewDetailMapper) ToReviewDetailsRecordPagination(Merchants []*db.GetReviewDetailsRow) []*record.ReviewDetailRecord {
	var result []*record.ReviewDetailRecord

	for _, Merchant := range Merchants {
		result = append(result, s.ToReviewDetailRecordPagination(Merchant))
	}

	return result
}

func (s *reviewDetailMapper) ToReviewDetailRecordActivePagination(detail *db.GetReviewDetailsActiveRow) *record.ReviewDetailRecord {
	var deletedAt *string
	if detail.DeletedAt.Valid {
		str := detail.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	reviewId := int(detail.ReviewID)

	return &record.ReviewDetailRecord{
		ID:        int(detail.ReviewDetailID),
		ReviewID:  &reviewId,
		Type:      detail.Type,
		Url:       detail.Url,
		Caption:   detail.Caption.String,
		CreatedAt: detail.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: detail.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *reviewDetailMapper) ToReviewDetailsRecordActivePagination(m []*db.GetReviewDetailsActiveRow) []*record.ReviewDetailRecord {
	var result []*record.ReviewDetailRecord

	for _, Merchant := range m {
		result = append(result, s.ToReviewDetailRecordActivePagination(Merchant))
	}

	return result
}

func (s *reviewDetailMapper) ToReviewDetailRecordTrashedPagination(detail *db.GetReviewDetailsTrashedRow) *record.ReviewDetailRecord {
	var deletedAt *string
	if detail.DeletedAt.Valid {
		str := detail.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	reviewId := int(detail.ReviewID)

	return &record.ReviewDetailRecord{
		ID:        int(detail.ReviewDetailID),
		ReviewID:  &reviewId,
		Type:      detail.Type,
		Url:       detail.Url,
		Caption:   detail.Caption.String,
		CreatedAt: detail.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: detail.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *reviewDetailMapper) ToReviewDetailsRecordTrashedPagination(Merchants []*db.GetReviewDetailsTrashedRow) []*record.ReviewDetailRecord {
	var result []*record.ReviewDetailRecord

	for _, Merchant := range Merchants {
		result = append(result, s.ToReviewDetailRecordTrashedPagination(Merchant))
	}

	return result
}
