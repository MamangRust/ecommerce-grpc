package seeder

import (
	"context"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/logger"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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
			Description:    strPtr("Manajemen mutu bersertifikat"),
			IssuedBy:       strPtr("ISO Organization"),
			IssueDate:      mustDate(2020, time.January, 15),
			ExpiryDate:     mustDate(2025, time.January, 15),
			CertificateUrl: strPtr("https://example.com/iso9001-cert.pdf"),
		},
		{
			MerchantID:     2,
			Title:          "Top UMKM 2023",
			Description:    strPtr("Penghargaan untuk UMKM terbaik tahun 2023"),
			IssuedBy:       strPtr("Kementerian Koperasi"),
			IssueDate:      mustDate(2023, time.July, 1),
			ExpiryDate:     nullDate(),
			CertificateUrl: strPtr("https://example.com/umkm-award-2023.pdf"),
		},
		{
			MerchantID:     3,
			Title:          "Halal Certified",
			Description:    strPtr("Sertifikasi halal dari MUI"),
			IssuedBy:       strPtr("Majelis Ulama Indonesia"),
			IssueDate:      mustDate(2021, time.March, 12),
			ExpiryDate:     mustDate(2024, time.March, 12),
			CertificateUrl: strPtr("https://example.com/halal-cert.pdf"),
		},
		{
			MerchantID:     4,
			Title:          "Best Food Product 2022",
			Description:    strPtr("Penghargaan untuk produk makanan terbaik tahun 2022"),
			IssuedBy:       strPtr("Asosiasi Kuliner Indonesia"),
			IssueDate:      mustDate(2022, time.November, 5),
			ExpiryDate:     nullDate(),
			CertificateUrl: strPtr("https://example.com/best-food-2022.pdf"),
		},
		{
			MerchantID:     5,
			Title:          "Eco-Friendly Business",
			Description:    strPtr("Sertifikasi bisnis ramah lingkungan"),
			IssuedBy:       strPtr("Green Business Council"),
			IssueDate:      mustDate(2023, time.April, 22),
			ExpiryDate:     mustDate(2026, time.April, 22),
			CertificateUrl: strPtr("https://example.com/eco-friendly-cert.pdf"),
		},
		{
			MerchantID:     6,
			Title:          "Top Seller 2023",
			Description:    strPtr("Penjual terbaik platform e-commerce tahun 2023"),
			IssuedBy:       strPtr("Tokopedia"),
			IssueDate:      mustDate(2024, time.January, 10),
			ExpiryDate:     nullDate(),
			CertificateUrl: strPtr("https://example.com/top-seller-2023.pdf"),
		},
		{
			MerchantID:     7,
			Title:          "BPOM Certified",
			Description:    strPtr("Sertifikasi produk dari Badan Pengawas Obat dan Makanan"),
			IssuedBy:       strPtr("Badan POM RI"),
			IssueDate:      mustDate(2022, time.August, 3),
			ExpiryDate:     mustDate(2025, time.August, 3),
			CertificateUrl: strPtr("https://example.com/bpom-cert.pdf"),
		},
		{
			MerchantID:     8,
			Title:          "Creativepreneur Award",
			Description:    strPtr("Penghargaan untuk wirausaha kreatif"),
			IssuedBy:       strPtr("Kementerian Pariwisata dan Ekonomi Kreatif"),
			IssueDate:      mustDate(2023, time.December, 15),
			ExpiryDate:     nullDate(),
			CertificateUrl: strPtr("https://example.com/creativepreneur-award.pdf"),
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

func nullDate() pgtype.Date {
	return pgtype.Date{Valid: false}
}
