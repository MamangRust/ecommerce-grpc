package recordmapper

import (
	"ecommerce/internal/domain/record"
	db "ecommerce/pkg/database/schema"
)

type reviewRecordMapper struct {
}

func NewReviewRecordMapper() *reviewRecordMapper {
	return &reviewRecordMapper{}
}

func (s *reviewRecordMapper) ToReviewRecord(review *db.Review) *record.ReviewRecord {
	var deletedAt *string
	if review.DeletedAt.Valid {
		deletedAtStr := review.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.ReviewRecord{
		ID:        int(review.ReviewID),
		UserID:    int(review.UserID),
		ProductID: int(review.ProductID),
		Name:      review.Name,
		Comment:   review.Comment,
		Rating:    int(review.Rating),
		CreatedAt: review.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: review.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *reviewRecordMapper) ToReviewRecordPagination(review *db.GetReviewsRow) *record.ReviewRecord {
	var deletedAt *string
	if review.DeletedAt.Valid {
		deletedAtStr := review.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.ReviewRecord{
		ID:        int(review.ReviewID),
		UserID:    int(review.UserID),
		ProductID: int(review.ProductID),
		Name:      review.Name,
		Comment:   review.Comment,
		Rating:    int(review.Rating),
		CreatedAt: review.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: review.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *reviewRecordMapper) ToReviewsRecordPagination(reviews []*db.GetReviewsRow) []*record.ReviewRecord {
	var result []*record.ReviewRecord

	for _, review := range reviews {
		result = append(result, s.ToReviewRecordPagination(review))
	}

	return result
}

func (s *reviewRecordMapper) ToReviewProductRecordPagination(review *db.GetReviewsByProductIDRow) *record.ReviewRecord {
	var deletedAt *string
	if review.DeletedAt.Valid {
		deletedAtStr := review.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.ReviewRecord{
		ID:        int(review.ReviewID),
		UserID:    int(review.UserID),
		ProductID: int(review.ProductID),
		Name:      review.Name,
		Comment:   review.Comment,
		Rating:    int(review.Rating),
		CreatedAt: review.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: review.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *reviewRecordMapper) ToReviewsProductRecordPagination(reviews []*db.GetReviewsByProductIDRow) []*record.ReviewRecord {
	var result []*record.ReviewRecord

	for _, review := range reviews {
		result = append(result, s.ToReviewProductRecordPagination(review))
	}

	return result
}

func (s *reviewRecordMapper) ToReviewRecordActivePagination(review *db.GetReviewsActiveRow) *record.ReviewRecord {
	var deletedAt *string
	if review.DeletedAt.Valid {
		deletedAtStr := review.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.ReviewRecord{
		ID:        int(review.ReviewID),
		UserID:    int(review.UserID),
		ProductID: int(review.ProductID),
		Name:      review.Name,
		Comment:   review.Comment,
		Rating:    int(review.Rating),
		CreatedAt: review.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: review.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *reviewRecordMapper) ToReviewsRecordActivePagination(reviews []*db.GetReviewsActiveRow) []*record.ReviewRecord {
	var result []*record.ReviewRecord

	for _, review := range reviews {
		result = append(result, s.ToReviewRecordActivePagination(review))
	}

	return result
}

func (s *reviewRecordMapper) ToReviewRecordTrashedPagination(review *db.GetReviewsTrashedRow) *record.ReviewRecord {
	var deletedAt *string
	if review.DeletedAt.Valid {
		deletedAtStr := review.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.ReviewRecord{
		ID:        int(review.ReviewID),
		UserID:    int(review.UserID),
		ProductID: int(review.ProductID),
		Name:      review.Name,
		Comment:   review.Comment,
		Rating:    int(review.Rating),
		CreatedAt: review.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: review.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *reviewRecordMapper) ToReviewsRecordTrashedPagination(reviews []*db.GetReviewsTrashedRow) []*record.ReviewRecord {
	var result []*record.ReviewRecord

	for _, review := range reviews {
		result = append(result, s.ToReviewRecordTrashedPagination(review))
	}

	return result
}
