package seeder

import (
	"context"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/logger"

	"go.uber.org/zap"
)

type productSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewProductSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *productSeeder {
	return &productSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *productSeeder) Seed() error {
	products := []db.CreateProductParams{
		{
			MerchantID:   1,
			CategoryID:   1,
			Name:         "Smartphone Galaxy X",
			Description:  strPtr("Smartphone dengan performa tinggi dan kamera canggih."),
			Price:        4_500_000,
			CountInStock: 20,
			Brand:        strPtr("Samsung"),
			Weight:       int32Ptr(300),
			Rating:       float64Ptr(4.5),
			SlugProduct:  strPtr("smartphone-galaxy-x"),
			ImageProduct: strPtr("galaxy-x.jpg"),
		},
		{
			MerchantID:   2,
			CategoryID:   2,
			Name:         "Facial Cleanser Glow",
			Description:  strPtr("Pembersih wajah dengan formula ringan untuk semua jenis kulit."),
			Price:        75_000,
			CountInStock: 100,
			Brand:        strPtr("GlowCare"),
			Weight:       int32Ptr(150),
			Rating:       float64Ptr(4.2),
			SlugProduct:  strPtr("facial-cleanser-glow"),
			ImageProduct: strPtr("cleanser.jpg"),
		},
		{
			MerchantID:   3,
			CategoryID:   3,
			Name:         "Blender Serbaguna",
			Description:  strPtr("Blender 3-in-1 untuk keperluan dapur sehari-hari."),
			Price:        350_000,
			CountInStock: 50,
			Brand:        strPtr("Maspion"),
			Weight:       int32Ptr(2000),
			Rating:       float64Ptr(4.0),
			SlugProduct:  strPtr("blender-serbaguna"),
			ImageProduct: strPtr("blender.jpg"),
		},
		{
			MerchantID:   4,
			CategoryID:   4,
			Name:         "Paket Popok Bayi Premium",
			Description:  strPtr("Popok bayi dengan teknologi anti bocor dan lembut di kulit."),
			Price:        120_000,
			CountInStock: 70,
			Brand:        strPtr("BabySoft"),
			Weight:       int32Ptr(1000),
			Rating:       float64Ptr(4.6),
			SlugProduct:  strPtr("popok-premium"),
			ImageProduct: strPtr("popok.jpg"),
		},
		{
			MerchantID:   5,
			CategoryID:   5,
			Name:         "Matras Yoga Premium",
			Description:  strPtr("Matras anti slip dengan ketebalan ideal untuk yoga dan fitness."),
			Price:        220_000,
			CountInStock: 40,
			Brand:        strPtr("FitZone"),
			Weight:       int32Ptr(700),
			Rating:       float64Ptr(4.4),
			SlugProduct:  strPtr("matras-yoga-premium"),
			ImageProduct: strPtr("matras.jpg"),
		},
		{
			MerchantID:   6,
			CategoryID:   6,
			Name:         "Snack Kentang Balado",
			Description:  strPtr("Cemilan kentang renyah dengan rasa balado khas."),
			Price:        18_000,
			CountInStock: 200,
			Brand:        strPtr("Snacky"),
			Weight:       int32Ptr(100),
			Rating:       float64Ptr(4.1),
			SlugProduct:  strPtr("kentang-balado"),
			ImageProduct: strPtr("snack.jpg"),
		},
		{
			MerchantID:   7,
			CategoryID:   7,
			Name:         "Controller PS5 DualSense",
			Description:  strPtr("Stik PS5 dengan fitur haptic feedback dan adaptive triggers."),
			Price:        999_000,
			CountInStock: 30,
			Brand:        strPtr("Sony"),
			Weight:       int32Ptr(450),
			Rating:       float64Ptr(4.8),
			SlugProduct:  strPtr("controller-ps5"),
			ImageProduct: strPtr("ps5-controller.jpg"),
		},
		{
			MerchantID:   8,
			CategoryID:   8,
			Name:         "Oli Motor Full Synthetic",
			Description:  strPtr("Oli mesin motor dengan perlindungan maksimal dan efisiensi tinggi."),
			Price:        95_000,
			CountInStock: 60,
			Brand:        strPtr("Motul"),
			Weight:       int32Ptr(1000),
			Rating:       float64Ptr(4.3),
			SlugProduct:  strPtr("oli-motor-synthetic"),
			ImageProduct: strPtr("oli.jpg"),
		},
	}

	for _, product := range products {
		if _, err := r.db.CreateProduct(r.ctx, product); err != nil {
			r.logger.Error(
				"failed to seed product",
				zap.String("name", product.Name),
				zap.Int32("merchant_id", product.MerchantID),
				zap.Error(err),
			)
			return err
		}
	}

	r.logger.Info("product successfully seeded")
	return nil
}

func float64Ptr(f float64) *float64 { return &f }
