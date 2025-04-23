package recordmapper

import (
	"ecommerce/internal/domain/record"
	db "ecommerce/pkg/database/schema"
)

type bannerRecordMapper struct {
}

func NewBannerRecordMapper() *bannerRecordMapper {
	return &bannerRecordMapper{}
}

func (s *bannerRecordMapper) ToBannerRecord(banner *db.Banner) *record.BannerRecord {
	var deletedAt *string
	if banner.DeletedAt.Valid {
		str := banner.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	return &record.BannerRecord{
		BannerID:  int(banner.BannerID),
		Name:      banner.Name,
		StartDate: banner.StartDate.Format("2006-01-02"),
		EndDate:   banner.EndDate.Format("2006-01-02"),
		StartTime: banner.StartTime.Format("15:04:05"),
		EndTime:   banner.EndTime.Format("15:04:05"),
		IsActive:  banner.IsActive.Valid && banner.IsActive.Bool,
		CreatedAt: banner.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: banner.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *bannerRecordMapper) ToBannerRecordPagination(banner *db.GetBannersRow) *record.BannerRecord {
	var deletedAt *string
	if banner.DeletedAt.Valid {
		str := banner.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	return &record.BannerRecord{
		BannerID:  int(banner.BannerID),
		Name:      banner.Name,
		StartDate: banner.StartDate.Format("2006-01-02"),
		EndDate:   banner.EndDate.Format("2006-01-02"),
		StartTime: banner.StartTime.Format("15:04:05"),
		EndTime:   banner.EndTime.Format("15:04:05"),
		IsActive:  banner.IsActive.Valid && banner.IsActive.Bool,
		CreatedAt: banner.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: banner.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *bannerRecordMapper) ToBannersRecordPagination(Merchants []*db.GetBannersRow) []*record.BannerRecord {
	var result []*record.BannerRecord

	for _, Merchant := range Merchants {
		result = append(result, s.ToBannerRecordPagination(Merchant))
	}

	return result
}

func (s *bannerRecordMapper) ToBannerRecordActivePagination(banner *db.GetBannersActiveRow) *record.BannerRecord {
	var deletedAt *string
	if banner.DeletedAt.Valid {
		str := banner.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	return &record.BannerRecord{
		BannerID:  int(banner.BannerID),
		Name:      banner.Name,
		StartDate: banner.StartDate.Format("2006-01-02"),
		EndDate:   banner.EndDate.Format("2006-01-02"),
		StartTime: banner.StartTime.Format("15:04:05"),
		EndTime:   banner.EndTime.Format("15:04:05"),
		IsActive:  banner.IsActive.Valid && banner.IsActive.Bool,
		CreatedAt: banner.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: banner.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *bannerRecordMapper) ToBannersRecordActivePagination(Merchants []*db.GetBannersActiveRow) []*record.BannerRecord {
	var result []*record.BannerRecord

	for _, Merchant := range Merchants {
		result = append(result, s.ToBannerRecordActivePagination(Merchant))
	}

	return result
}

func (s *bannerRecordMapper) ToBannerRecordTrashedPagination(banner *db.GetBannersTrashedRow) *record.BannerRecord {
	var deletedAt *string
	if banner.DeletedAt.Valid {
		str := banner.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &str
	}

	return &record.BannerRecord{
		BannerID:  int(banner.BannerID),
		Name:      banner.Name,
		StartDate: banner.StartDate.Format("2006-01-02"),
		EndDate:   banner.EndDate.Format("2006-01-02"),
		StartTime: banner.StartTime.Format("15:04:05"),
		EndTime:   banner.EndTime.Format("15:04:05"),
		IsActive:  banner.IsActive.Valid && banner.IsActive.Bool,
		CreatedAt: banner.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: banner.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *bannerRecordMapper) ToBannersRecordTrashedPagination(banner []*db.GetBannersTrashedRow) []*record.BannerRecord {
	var result []*record.BannerRecord

	for _, Merchant := range banner {
		result = append(result, s.ToBannerRecordTrashedPagination(Merchant))
	}

	return result
}
