package seeder

import (
	"context"
	"database/sql"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/logger"

	"go.uber.org/zap"
)

type categorySeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewCategorySeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *categorySeeder {
	return &categorySeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *categorySeeder) Seed() error {
	categories := []db.CreateCategoryParams{
		{
			Name:          "Elektronik",
			Description:   sql.NullString{String: "Produk elektronik seperti smartphone, laptop, dan aksesori elektronik lainnya.", Valid: true},
			SlugCategory:  sql.NullString{String: "elektronik", Valid: true},
			ImageCategory: sql.NullString{String: "elektronik.jpg", Valid: true},
		},
		{
			Name:          "Kesehatan & Kecantikan",
			Description:   sql.NullString{String: "Produk perawatan tubuh, skincare, dan suplemen kesehatan.", Valid: true},
			SlugCategory:  sql.NullString{String: "kesehatan-kecantikan", Valid: true},
			ImageCategory: sql.NullString{String: "kesehatan.jpg", Valid: true},
		},
		{
			Name:          "Peralatan Rumah Tangga",
			Description:   sql.NullString{String: "Peralatan dapur, perlengkapan rumah, dan furnitur.", Valid: true},
			SlugCategory:  sql.NullString{String: "peralatan-rumah", Valid: true},
			ImageCategory: sql.NullString{String: "rumah.jpg", Valid: true},
		},
		{
			Name:          "Ibu & Bayi",
			Description:   sql.NullString{String: "Produk khusus untuk ibu hamil, menyusui, dan bayi.", Valid: true},
			SlugCategory:  sql.NullString{String: "ibu-bayi", Valid: true},
			ImageCategory: sql.NullString{String: "ibu-bayi.jpg", Valid: true},
		},
		{
			Name:          "Olahraga & Outdoor",
			Description:   sql.NullString{String: "Perlengkapan olahraga, fitness, dan kegiatan luar ruangan.", Valid: true},
			SlugCategory:  sql.NullString{String: "olahraga-outdoor", Valid: true},
			ImageCategory: sql.NullString{String: "olahraga.jpg", Valid: true},
		},
		{
			Name:          "Makanan & Minuman",
			Description:   sql.NullString{String: "Makanan ringan, minuman, bahan makanan segar dan kemasan.", Valid: true},
			SlugCategory:  sql.NullString{String: "makanan-minuman", Valid: true},
			ImageCategory: sql.NullString{String: "makanan.jpg", Valid: true},
		},
		{
			Name:          "Gaming & Console",
			Description:   sql.NullString{String: "Konsol game, aksesori, dan game terbaru dari berbagai platform.", Valid: true},
			SlugCategory:  sql.NullString{String: "gaming-console", Valid: true},
			ImageCategory: sql.NullString{String: "gaming.jpg", Valid: true},
		},
		{
			Name:          "Perlengkapan Otomotif",
			Description:   sql.NullString{String: "Aksesori mobil dan motor, oli, serta sparepart kendaraan.", Valid: true},
			SlugCategory:  sql.NullString{String: "otomotif", Valid: true},
			ImageCategory: sql.NullString{String: "otomotif.jpg", Valid: true},
		},
	}

	for _, category := range categories {
		if _, err := r.db.CreateCategory(r.ctx, category); err != nil {
			r.logger.Error("Failed to insert category", zap.Error(err))
			return err
		}
	}

	r.logger.Info("Successfully seeded 10 categories")
	return nil
}
