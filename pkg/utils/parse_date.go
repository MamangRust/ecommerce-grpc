package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ParseDate(value string) (pgtype.Date, error) {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return pgtype.Date{}, err
	}

	return pgtype.Date{
		Time:  t,
		Valid: true,
	}, nil
}

func ParseTime(value string) (pgtype.Time, error) {
	t, err := time.Parse("15:04", value)
	if err != nil {
		return pgtype.Time{}, err
	}

	return pgtype.Time{
		Microseconds: int64(t.Hour()*3600+t.Minute()*60) * 1_000_000,
		Valid:        true,
	}, nil
}
