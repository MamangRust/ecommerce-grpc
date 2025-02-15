package recordmapper

import (
	"ecommerce/internal/domain/record"
	db "ecommerce/pkg/database/schema"
)

type categoryRecordMapper struct {
}

func NewCategoryRecordMapper() *categoryRecordMapper {
	return &categoryRecordMapper{}
}

func (s *categoryRecordMapper) ToCategoryRecord(category *db.Category) *record.CategoriesRecord {
	var deletedAt *string
	if category.DeletedAt.Valid {
		deletedAtStr := category.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.CategoriesRecord{
		ID:            int(category.CategoryID),
		Name:          category.Name,
		Description:   category.Description.String,
		SlugCategory:  category.SlugCategory.String,
		ImageCategory: category.ImageCategory.String,
		CreatedAt:     category.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:     category.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:     deletedAt,
	}
}

func (s *categoryRecordMapper) ToCategoryRecordPagination(category *db.GetCategoriesRow) *record.CategoriesRecord {
	var deletedAt *string
	if category.DeletedAt.Valid {
		deletedAtStr := category.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.CategoriesRecord{
		ID:            int(category.CategoryID),
		Name:          category.Name,
		Description:   category.Description.String,
		SlugCategory:  category.SlugCategory.String,
		ImageCategory: category.ImageCategory.String,
		CreatedAt:     category.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:     category.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:     deletedAt,
	}
}

func (s *categoryRecordMapper) ToCategoriesRecordPagination(categories []*db.GetCategoriesRow) []*record.CategoriesRecord {
	var result []*record.CategoriesRecord

	for _, category := range categories {
		result = append(result, s.ToCategoryRecordPagination(category))
	}

	return result
}

func (s *categoryRecordMapper) ToCategoryRecordActivePagination(category *db.GetCategoriesActiveRow) *record.CategoriesRecord {
	var deletedAt *string
	if category.DeletedAt.Valid {
		deletedAtStr := category.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.CategoriesRecord{
		ID:            int(category.CategoryID),
		Name:          category.Name,
		Description:   category.Description.String,
		SlugCategory:  category.SlugCategory.String,
		ImageCategory: category.ImageCategory.String,
		CreatedAt:     category.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:     category.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:     deletedAt,
	}
}

func (s *categoryRecordMapper) ToCategoriesRecordActivePagination(categories []*db.GetCategoriesActiveRow) []*record.CategoriesRecord {
	var result []*record.CategoriesRecord

	for _, category := range categories {
		result = append(result, s.ToCategoryRecordActivePagination(category))
	}

	return result
}

func (s *categoryRecordMapper) ToCategoryRecordTrashedPagination(category *db.GetCategoriesTrashedRow) *record.CategoriesRecord {
	var deletedAt *string
	if category.DeletedAt.Valid {
		deletedAtStr := category.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.CategoriesRecord{
		ID:            int(category.CategoryID),
		Name:          category.Name,
		Description:   category.Description.String,
		SlugCategory:  category.SlugCategory.String,
		ImageCategory: category.ImageCategory.String,
		CreatedAt:     category.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:     category.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:     deletedAt,
	}
}

func (s *categoryRecordMapper) ToCategoriesRecordTrashedPagination(categories []*db.GetCategoriesTrashedRow) []*record.CategoriesRecord {
	var result []*record.CategoriesRecord

	for _, category := range categories {
		result = append(result, s.ToCategoryRecordTrashedPagination(category))
	}

	return result
}
