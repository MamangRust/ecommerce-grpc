package seeder

import (
	"context"
	"database/sql"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/logger"
	"time"

	"go.uber.org/zap"
)

type merchantAwardSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantAwardSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *merchantAwardSeeder {
	return &merchantAwardSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *merchantAwardSeeder) Seed() error {
	awards := []db.CreateMerchantCertificationOrAwardParams{
		{
			MerchantID:     1,
			Title:          "ISO 9001 Certified",
			Description:    sql.NullString{String: "Manajemen mutu bersertifikat", Valid: true},
			IssuedBy:       sql.NullString{String: "ISO Organization", Valid: true},
			IssueDate:      sql.NullTime{Time: time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			ExpiryDate:     sql.NullTime{Time: time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			CertificateUrl: sql.NullString{String: "https://example.com/iso9001-cert.pdf", Valid: true},
		},
		{
			MerchantID:     2,
			Title:          "Top UMKM 2023",
			Description:    sql.NullString{String: "Penghargaan untuk UMKM terbaik tahun 2023", Valid: true},
			IssuedBy:       sql.NullString{String: "Kementerian Koperasi", Valid: true},
			IssueDate:      sql.NullTime{Time: time.Date(2023, time.July, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			ExpiryDate:     sql.NullTime{Valid: false},
			CertificateUrl: sql.NullString{String: "https://example.com/umkm-award-2023.pdf", Valid: true},
		},
		{
			MerchantID:     3,
			Title:          "Halal Certified",
			Description:    sql.NullString{String: "Sertifikasi halal dari MUI", Valid: true},
			IssuedBy:       sql.NullString{String: "Majelis Ulama Indonesia", Valid: true},
			IssueDate:      sql.NullTime{Time: time.Date(2021, time.March, 12, 0, 0, 0, 0, time.UTC), Valid: true},
			ExpiryDate:     sql.NullTime{Time: time.Date(2024, time.March, 12, 0, 0, 0, 0, time.UTC), Valid: true},
			CertificateUrl: sql.NullString{String: "https://example.com/halal-cert.pdf", Valid: true},
		},
		{
			MerchantID:     4,
			Title:          "Best Food Product 2022",
			Description:    sql.NullString{String: "Penghargaan untuk produk makanan terbaik tahun 2022", Valid: true},
			IssuedBy:       sql.NullString{String: "Asosiasi Kuliner Indonesia", Valid: true},
			IssueDate:      sql.NullTime{Time: time.Date(2022, time.November, 5, 0, 0, 0, 0, time.UTC), Valid: true},
			ExpiryDate:     sql.NullTime{Valid: false},
			CertificateUrl: sql.NullString{String: "https://example.com/best-food-2022.pdf", Valid: true},
		},
		{
			MerchantID:     5,
			Title:          "Eco-Friendly Business",
			Description:    sql.NullString{String: "Sertifikasi bisnis ramah lingkungan", Valid: true},
			IssuedBy:       sql.NullString{String: "Green Business Council", Valid: true},
			IssueDate:      sql.NullTime{Time: time.Date(2023, time.April, 22, 0, 0, 0, 0, time.UTC), Valid: true},
			ExpiryDate:     sql.NullTime{Time: time.Date(2026, time.April, 22, 0, 0, 0, 0, time.UTC), Valid: true},
			CertificateUrl: sql.NullString{String: "https://example.com/eco-friendly-cert.pdf", Valid: true},
		},
		{
			MerchantID:     6,
			Title:          "Top Seller 2023",
			Description:    sql.NullString{String: "Penjual terbaik platform e-commerce tahun 2023", Valid: true},
			IssuedBy:       sql.NullString{String: "Tokopedia", Valid: true},
			IssueDate:      sql.NullTime{Time: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC), Valid: true},
			ExpiryDate:     sql.NullTime{Valid: false},
			CertificateUrl: sql.NullString{String: "https://example.com/top-seller-2023.pdf", Valid: true},
		},
		{
			MerchantID:     7,
			Title:          "BPOM Certified",
			Description:    sql.NullString{String: "Sertifikasi produk dari Badan Pengawas Obat dan Makanan", Valid: true},
			IssuedBy:       sql.NullString{String: "Badan POM RI", Valid: true},
			IssueDate:      sql.NullTime{Time: time.Date(2022, time.August, 3, 0, 0, 0, 0, time.UTC), Valid: true},
			ExpiryDate:     sql.NullTime{Time: time.Date(2025, time.August, 3, 0, 0, 0, 0, time.UTC), Valid: true},
			CertificateUrl: sql.NullString{String: "https://example.com/bpom-cert.pdf", Valid: true},
		},
		{
			MerchantID:     8,
			Title:          "Creativepreneur Award",
			Description:    sql.NullString{String: "Penghargaan untuk wirausaha kreatif", Valid: true},
			IssuedBy:       sql.NullString{String: "Kementerian Pariwisata dan Ekonomi Kreatif", Valid: true},
			IssueDate:      sql.NullTime{Time: time.Date(2023, time.December, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			ExpiryDate:     sql.NullTime{Valid: false},
			CertificateUrl: sql.NullString{String: "https://example.com/creativepreneur-award.pdf", Valid: true},
		},
	}

	for _, award := range awards {
		if _, err := r.db.CreateMerchantCertificationOrAward(r.ctx, award); err != nil {
			r.logger.Error("failed to seed merchant award", zap.Error(err))
			return err
		}
	}

	r.logger.Info("merchant awards seeded successfully")

	return nil
}
