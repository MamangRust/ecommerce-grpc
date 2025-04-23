package recordmapper

import (
	"ecommerce/internal/domain/record"
	db "ecommerce/pkg/database/schema"
	"encoding/json"
	"fmt"
)

type merchantDetailRecordMapper struct {
}

func NewMerchantDetailRecordMapper() *merchantDetailRecordMapper {
	return &merchantDetailRecordMapper{}
}

func (s *merchantDetailRecordMapper) ToMerchantDetailRecord(detail *db.MerchantDetail) *record.MerchantDetailRecord {
	var deletedAt *string
	if detail.DeletedAt.Valid {
		str := detail.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	return &record.MerchantDetailRecord{
		ID:               int(detail.MerchantDetailID),
		MerchantID:       int(detail.MerchantID),
		DisplayName:      detail.DisplayName.String,
		CoverImageUrl:    detail.CoverImageUrl.String,
		LogoUrl:          detail.LogoUrl.String,
		ShortDescription: detail.ShortDescription.String,
		WebsiteUrl:       detail.WebsiteUrl.String,
		CreatedAt:        detail.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:        detail.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:        deletedAt,
	}
}

func (s *merchantDetailRecordMapper) ToMerchantDetailRelationRecord(detail *db.GetMerchantDetailRow) *record.MerchantDetailRecord {
	var deletedAt *string
	if detail.DeletedAt.Valid {
		str := detail.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	var socialMediaLinks []*record.MerchantSocialMediaLinkRecord
	if len(detail.SocialMediaLinks) > 0 {
		err := json.Unmarshal(detail.SocialMediaLinks, &socialMediaLinks)
		if err != nil {
			fmt.Println("Error unmarshaling social media links:", err)
			socialMediaLinks = nil
		}
	}

	return &record.MerchantDetailRecord{
		ID:               int(detail.MerchantDetailID),
		MerchantID:       int(detail.MerchantID),
		DisplayName:      detail.DisplayName.String,
		CoverImageUrl:    detail.CoverImageUrl.String,
		LogoUrl:          detail.LogoUrl.String,
		ShortDescription: detail.ShortDescription.String,
		WebsiteUrl:       detail.WebsiteUrl.String,
		SocialMediaLinks: socialMediaLinks,
		CreatedAt:        detail.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:        detail.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:        deletedAt,
	}
}

func (s *merchantDetailRecordMapper) ToMerchantDetailRecordPagination(detail *db.GetMerchantDetailsRow) *record.MerchantDetailRecord {
	var deletedAt *string
	if detail.DeletedAt.Valid {
		deletedAtStr := detail.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	var socialMediaLinks []*record.MerchantSocialMediaLinkRecord
	if len(detail.SocialMediaLinks) > 0 {
		err := json.Unmarshal(detail.SocialMediaLinks, &socialMediaLinks)
		if err != nil {
			fmt.Println("Error unmarshaling social media links:", err)
			socialMediaLinks = nil
		}
	}

	return &record.MerchantDetailRecord{
		ID:               int(detail.MerchantDetailID),
		MerchantID:       int(detail.MerchantID),
		DisplayName:      detail.DisplayName.String,
		CoverImageUrl:    detail.CoverImageUrl.String,
		LogoUrl:          detail.LogoUrl.String,
		ShortDescription: detail.ShortDescription.String,
		WebsiteUrl:       detail.WebsiteUrl.String,
		SocialMediaLinks: socialMediaLinks,
		CreatedAt:        detail.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:        detail.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:        deletedAt,
	}
}

func (s *merchantDetailRecordMapper) ToMerchantDetailsRecordPagination(MerchantDetails []*db.GetMerchantDetailsRow) []*record.MerchantDetailRecord {
	var result []*record.MerchantDetailRecord

	for _, MerchantDetail := range MerchantDetails {
		result = append(result, s.ToMerchantDetailRecordPagination(MerchantDetail))
	}

	return result
}

func (s *merchantDetailRecordMapper) ToMerchantDetailRecordActivePagination(detail *db.GetMerchantDetailsActiveRow) *record.MerchantDetailRecord {
	var deletedAt *string
	if detail.DeletedAt.Valid {
		deletedAtStr := detail.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	var socialMediaLinks []*record.MerchantSocialMediaLinkRecord
	if len(detail.SocialMediaLinks) > 0 {
		err := json.Unmarshal(detail.SocialMediaLinks, &socialMediaLinks)
		if err != nil {
			fmt.Println("Error unmarshaling social media links:", err)
			socialMediaLinks = nil
		}
	}

	return &record.MerchantDetailRecord{
		ID:               int(detail.MerchantDetailID),
		MerchantID:       int(detail.MerchantID),
		DisplayName:      detail.DisplayName.String,
		CoverImageUrl:    detail.CoverImageUrl.String,
		LogoUrl:          detail.LogoUrl.String,
		ShortDescription: detail.ShortDescription.String,
		WebsiteUrl:       detail.WebsiteUrl.String,
		SocialMediaLinks: socialMediaLinks,
		CreatedAt:        detail.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:        detail.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:        deletedAt,
	}
}

func (s *merchantDetailRecordMapper) ToMerchantDetailsRecordActivePagination(MerchantDetails []*db.GetMerchantDetailsActiveRow) []*record.MerchantDetailRecord {
	var result []*record.MerchantDetailRecord

	for _, MerchantDetail := range MerchantDetails {
		result = append(result, s.ToMerchantDetailRecordActivePagination(MerchantDetail))
	}

	return result
}

func (s *merchantDetailRecordMapper) ToMerchantDetailRecordTrashedPagination(detail *db.GetMerchantDetailsTrashedRow) *record.MerchantDetailRecord {
	var deletedAt *string
	if detail.DeletedAt.Valid {
		deletedAtStr := detail.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	var socialMediaLinks []*record.MerchantSocialMediaLinkRecord
	if len(detail.SocialMediaLinks) > 0 {
		err := json.Unmarshal(detail.SocialMediaLinks, &socialMediaLinks)
		if err != nil {
			fmt.Println("Error unmarshaling social media links:", err)
			socialMediaLinks = nil
		}
	}

	return &record.MerchantDetailRecord{
		ID:               int(detail.MerchantDetailID),
		MerchantID:       int(detail.MerchantID),
		DisplayName:      detail.DisplayName.String,
		CoverImageUrl:    detail.CoverImageUrl.String,
		LogoUrl:          detail.LogoUrl.String,
		ShortDescription: detail.ShortDescription.String,
		WebsiteUrl:       detail.WebsiteUrl.String,
		SocialMediaLinks: socialMediaLinks,
		CreatedAt:        detail.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:        detail.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:        deletedAt,
	}
}

func (s *merchantDetailRecordMapper) ToMerchantDetailsRecordTrashedPagination(MerchantDetails []*db.GetMerchantDetailsTrashedRow) []*record.MerchantDetailRecord {
	var result []*record.MerchantDetailRecord

	for _, MerchantDetail := range MerchantDetails {
		result = append(result, s.ToMerchantDetailRecordTrashedPagination(MerchantDetail))
	}

	return result
}
