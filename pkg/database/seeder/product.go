package seeder

import (
	"context"
	"database/sql"
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
			Description:  sql.NullString{String: "Smartphone dengan performa tinggi dan kamera canggih.", Valid: true},
			Price:        4500000,
			CountInStock: 20,
			Brand:        sql.NullString{String: "Samsung", Valid: true},
			Weight:       sql.NullInt32{Int32: 300, Valid: true},
			Rating:       sql.NullFloat64{Float64: 4.5, Valid: true},
			SlugProduct:  sql.NullString{String: "smartphone-galaxy-x", Valid: true},
			ImageProduct: sql.NullString{String: "galaxy-x.jpg", Valid: true},
		},
		{
			MerchantID:   2,
			CategoryID:   2,
			Name:         "Facial Cleanser Glow",
			Description:  sql.NullString{String: "Pembersih wajah dengan formula ringan untuk semua jenis kulit.", Valid: true},
			Price:        75000,
			CountInStock: 100,
			Brand:        sql.NullString{String: "GlowCare", Valid: true},
			Weight:       sql.NullInt32{Int32: 150, Valid: true},
			Rating:       sql.NullFloat64{Float64: 4.2, Valid: true},
			SlugProduct:  sql.NullString{String: "facial-cleanser-glow", Valid: true},
			ImageProduct: sql.NullString{String: "cleanser.jpg", Valid: true},
		},
		{
			MerchantID:   3,
			CategoryID:   3,
			Name:         "Blender Serbaguna",
			Description:  sql.NullString{String: "Blender 3-in-1 untuk keperluan dapur sehari-hari.", Valid: true},
			Price:        350000,
			CountInStock: 50,
			Brand:        sql.NullString{String: "Maspion", Valid: true},
			Weight:       sql.NullInt32{Int32: 2000, Valid: true},
			Rating:       sql.NullFloat64{Float64: 4.0, Valid: true},
			SlugProduct:  sql.NullString{String: "blender-serbaguna", Valid: true},
			ImageProduct: sql.NullString{String: "blender.jpg", Valid: true},
		},
		{
			MerchantID:   4,
			CategoryID:   4,
			Name:         "Paket Popok Bayi Premium",
			Description:  sql.NullString{String: "Popok bayi dengan teknologi anti bocor dan lembut di kulit.", Valid: true},
			Price:        120000,
			CountInStock: 70,
			Brand:        sql.NullString{String: "BabySoft", Valid: true},
			Weight:       sql.NullInt32{Int32: 1000, Valid: true},
			Rating:       sql.NullFloat64{Float64: 4.6, Valid: true},
			SlugProduct:  sql.NullString{String: "popok-premium", Valid: true},
			ImageProduct: sql.NullString{String: "popok.jpg", Valid: true},
		},
		{
			MerchantID:   5,
			CategoryID:   5,
			Name:         "Matras Yoga Premium",
			Description:  sql.NullString{String: "Matras anti slip dengan ketebalan ideal untuk yoga dan fitness.", Valid: true},
			Price:        220000,
			CountInStock: 40,
			Brand:        sql.NullString{String: "FitZone", Valid: true},
			Weight:       sql.NullInt32{Int32: 700, Valid: true},
			Rating:       sql.NullFloat64{Float64: 4.4, Valid: true},
			SlugProduct:  sql.NullString{String: "matras-yoga-premium", Valid: true},
			ImageProduct: sql.NullString{String: "matras.jpg", Valid: true},
		},
		{
			MerchantID:   6,
			CategoryID:   6,
			Name:         "Snack Kentang Balado",
			Description:  sql.NullString{String: "Cemilan kentang renyah dengan rasa balado khas.", Valid: true},
			Price:        18000,
			CountInStock: 200,
			Brand:        sql.NullString{String: "Snacky", Valid: true},
			Weight:       sql.NullInt32{Int32: 100, Valid: true},
			Rating:       sql.NullFloat64{Float64: 4.1, Valid: true},
			SlugProduct:  sql.NullString{String: "kentang-balado", Valid: true},
			ImageProduct: sql.NullString{String: "snack.jpg", Valid: true},
		},
		{
			MerchantID:   7,
			CategoryID:   7,
			Name:         "Controller PS5 DualSense",
			Description:  sql.NullString{String: "Stik PS5 dengan fitur haptic feedback dan adaptive triggers.", Valid: true},
			Price:        999000,
			CountInStock: 30,
			Brand:        sql.NullString{String: "Sony", Valid: true},
			Weight:       sql.NullInt32{Int32: 450, Valid: true},
			Rating:       sql.NullFloat64{Float64: 4.8, Valid: true},
			SlugProduct:  sql.NullString{String: "controller-ps5", Valid: true},
			ImageProduct: sql.NullString{String: "ps5-controller.jpg", Valid: true},
		},
		{
			MerchantID:   8,
			CategoryID:   8,
			Name:         "Oli Motor Full Synthetic",
			Description:  sql.NullString{String: "Oli mesin motor dengan perlindungan maksimal dan efisiensi tinggi.", Valid: true},
			Price:        95000,
			CountInStock: 60,
			Brand:        sql.NullString{String: "Motul", Valid: true},
			Weight:       sql.NullInt32{Int32: 1000, Valid: true},
			Rating:       sql.NullFloat64{Float64: 4.3, Valid: true},
			SlugProduct:  sql.NullString{String: "oli-motor-synthetic", Valid: true},
			ImageProduct: sql.NullString{String: "oli.jpg", Valid: true},
		},
	}

	for _, product := range products {
		if _, err := r.db.CreateProduct(r.ctx, product); err != nil {
			r.logger.Error("failed to seed product", zap.Error(err))
			return err
		}
	}

	r.logger.Info("product successfully seeded")

	return nil
}
