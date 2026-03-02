package seeder

import (
	"context"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/logger"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type bannerSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewBannerSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *bannerSeeder {
	return &bannerSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *bannerSeeder) Seed() error {
	banners := []db.CreateBannerParams{
		{
			Name:      "Banner 1",
			StartDate: mustDate(2023, 1, 1),
			EndDate:   mustDate(2023, 1, 31),
			StartTime: mustTime("08:00"),
			EndTime:   mustTime("16:00"),
			IsActive:  boolPtr(true),
		},
		{
			Name:      "Banner 2",
			StartDate: mustDate(2023, 2, 1),
			EndDate:   mustDate(2023, 2, 28),
			StartTime: mustTime("09:00"),
			EndTime:   mustTime("17:00"),
			IsActive:  boolPtr(true),
		},
		{
			Name:      "Banner 3",
			StartDate: mustDate(2023, 3, 1),
			EndDate:   mustDate(2023, 3, 31),
			StartTime: mustTime("10:00"),
			EndTime:   mustTime("18:00"),
			IsActive:  boolPtr(false),
		},
		{
			Name:      "Banner 4",
			StartDate: mustDate(2023, 4, 1),
			EndDate:   mustDate(2023, 4, 30),
			StartTime: mustTime("07:00"),
			EndTime:   mustTime("15:00"),
			IsActive:  boolPtr(true),
		},
		{
			Name:      "Banner 5",
			StartDate: mustDate(2023, 5, 1),
			EndDate:   mustDate(2023, 5, 31),
			StartTime: mustTime("06:00"),
			EndTime:   mustTime("14:00"),
			IsActive:  boolPtr(false),
		},
		{
			Name:      "Banner 6",
			StartDate: mustDate(2023, 6, 1),
			EndDate:   mustDate(2023, 6, 30),
			StartTime: mustTime("12:00"),
			EndTime:   mustTime("20:00"),
			IsActive:  boolPtr(true),
		},
		{
			Name:      "Banner 7",
			StartDate: mustDate(2023, 7, 1),
			EndDate:   mustDate(2023, 7, 31),
			StartTime: mustTime("08:30"),
			EndTime:   mustTime("16:30"),
			IsActive:  boolPtr(true),
		},
		{
			Name:      "Banner 8",
			StartDate: mustDate(2023, 8, 1),
			EndDate:   mustDate(2023, 8, 31),
			StartTime: mustTime("09:30"),
			EndTime:   mustTime("17:30"),
			IsActive:  boolPtr(false),
		},
	}

	for _, banner := range banners {
		if _, err := r.db.CreateBanner(r.ctx, banner); err != nil {
			r.logger.Error("failed to insert banner", zap.Error(err))
			return err
		}
	}

	r.logger.Info("banner successfully seeded")
	return nil
}

func mustDate(y int, m time.Month, d int) pgtype.Date {
	return pgtype.Date{
		Time:  time.Date(y, m, d, 0, 0, 0, 0, time.UTC),
		Valid: true,
	}
}

func mustTime(v string) pgtype.Time {
	t, err := time.Parse("15:04", v)
	if err != nil {
		panic(err)
	}

	return pgtype.Time{
		Microseconds: int64(
			t.Hour()*60*60*1_000_000 +
				t.Minute()*60*1_000_000,
		),
		Valid: true,
	}
}

func boolPtr(v bool) *bool {
	return &v
}
