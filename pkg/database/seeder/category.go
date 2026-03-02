package seeder

import (
	"context"
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
			Description:   strPtr("Produk elektronik seperti smartphone, laptop, dan aksesori elektronik lainnya."),
			SlugCategory:  strPtr("elektronik"),
			ImageCategory: strPtr("elektronik.jpg"),
		},
		{
			Name:          "Kesehatan & Kecantikan",
			Description:   strPtr("Produk perawatan tubuh, skincare, dan suplemen kesehatan."),
			SlugCategory:  strPtr("kesehatan-kecantikan"),
			ImageCategory: strPtr("kesehatan.jpg"),
		},
		{
			Name:          "Peralatan Rumah Tangga",
			Description:   strPtr("Peralatan dapur, perlengkapan rumah, dan furnitur."),
			SlugCategory:  strPtr("peralatan-rumah"),
			ImageCategory: strPtr("rumah.jpg"),
		},
		{
			Name:          "Ibu & Bayi",
			Description:   strPtr("Produk khusus untuk ibu hamil, menyusui, dan bayi."),
			SlugCategory:  strPtr("ibu-bayi"),
			ImageCategory: strPtr("ibu-bayi.jpg"),
		},
		{
			Name:          "Olahraga & Outdoor",
			Description:   strPtr("Perlengkapan olahraga, fitness, dan kegiatan luar ruangan."),
			SlugCategory:  strPtr("olahraga-outdoor"),
			ImageCategory: strPtr("olahraga.jpg"),
		},
		{
			Name:          "Makanan & Minuman",
			Description:   strPtr("Makanan ringan, minuman, bahan makanan segar dan kemasan."),
			SlugCategory:  strPtr("makanan-minuman"),
			ImageCategory: strPtr("makanan.jpg"),
		},
		{
			Name:          "Gaming & Console",
			Description:   strPtr("Konsol game, aksesori, dan game terbaru dari berbagai platform."),
			SlugCategory:  strPtr("gaming-console"),
			ImageCategory: strPtr("gaming.jpg"),
		},
		{
			Name:          "Perlengkapan Otomotif",
			Description:   strPtr("Aksesori mobil dan motor, oli, serta sparepart kendaraan."),
			SlugCategory:  strPtr("otomotif"),
			ImageCategory: strPtr("otomotif.jpg"),
		},
	}

	for _, category := range categories {
		if _, err := r.db.CreateCategory(r.ctx, category); err != nil {
			r.logger.Error("failed to insert category", zap.Error(err))
			return err
		}
	}

	r.logger.Info("successfully seeded categories")
	return nil
}

func strPtr(v string) *string {
	return &v
}
