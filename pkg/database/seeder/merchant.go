package seeder

import (
	"context"
	"database/sql"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/logger"

	"go.uber.org/zap"
)

type merchantSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *merchantSeeder {
	return &merchantSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *merchantSeeder) Seed() error {
	merchants := []db.CreateMerchantParams{
		{
			UserID:       1,
			Name:         "Elektronik Store",
			Description:  sql.NullString{String: "Toko elektronik terpercaya dengan berbagai produk gadget dan aksesoris.", Valid: true},
			Address:      sql.NullString{String: "Jl. Teknologi No.1, Jakarta", Valid: true},
			ContactEmail: sql.NullString{String: "support@elektronikstore.com", Valid: true},
			ContactPhone: sql.NullString{String: "081234567890", Valid: true},
			Status:       "active",
		},
		{
			UserID:       2,
			Name:         "Kecantikan Sehat",
			Description:  sql.NullString{String: "Produk skincare dan kesehatan pilihan.", Valid: true},
			Address:      sql.NullString{String: "Jl. Kesehatan No.5, Bandung", Valid: true},
			ContactEmail: sql.NullString{String: "cs@kecantikansehat.com", Valid: true},
			ContactPhone: sql.NullString{String: "082345678901", Valid: true},
			Status:       "active",
		},
		{
			UserID:       3,
			Name:         "Rumah Indah",
			Description:  sql.NullString{String: "Peralatan rumah tangga berkualitas dan estetik.", Valid: true},
			Address:      sql.NullString{String: "Jl. Rumah No.12, Surabaya", Valid: true},
			ContactEmail: sql.NullString{String: "info@rumahindah.com", Valid: true},
			ContactPhone: sql.NullString{String: "083456789012", Valid: true},
			Status:       "active",
		},
		{
			UserID:       4,
			Name:         "Mom & Baby Care",
			Description:  sql.NullString{String: "Semua kebutuhan ibu dan bayi ada di sini.", Valid: true},
			Address:      sql.NullString{String: "Jl. Keluarga No.7, Depok", Valid: true},
			ContactEmail: sql.NullString{String: "support@momandbaby.com", Valid: true},
			ContactPhone: sql.NullString{String: "084567890123", Valid: true},
			Status:       "active",
		},
		{
			UserID:       5,
			Name:         "Sport Zone",
			Description:  sql.NullString{String: "Perlengkapan olahraga dan outdoor terlengkap.", Valid: true},
			Address:      sql.NullString{String: "Jl. Atletik No.3, Yogyakarta", Valid: true},
			ContactEmail: sql.NullString{String: "halo@sportzone.com", Valid: true},
			ContactPhone: sql.NullString{String: "085678901234", Valid: true},
			Status:       "active",
		},
		{
			UserID:       6,
			Name:         "Fresh Mart",
			Description:  sql.NullString{String: "Toko makanan dan minuman segar dan kemasan.", Valid: true},
			Address:      sql.NullString{String: "Jl. Pasar No.10, Semarang", Valid: true},
			ContactEmail: sql.NullString{String: "fresh@mart.com", Valid: true},
			ContactPhone: sql.NullString{String: "086789012345", Valid: true},
			Status:       "active",
		},
		{
			UserID:       7,
			Name:         "Gamer Heaven",
			Description:  sql.NullString{String: "Game, console, dan aksesori lengkap untuk gamers.", Valid: true},
			Address:      sql.NullString{String: "Jl. Game No.8, Bekasi", Valid: true},
			ContactEmail: sql.NullString{String: "gamer@heaven.com", Valid: true},
			ContactPhone: sql.NullString{String: "087890123456", Valid: true},
			Status:       "active",
		},
		{
			UserID:       8,
			Name:         "AutoParts Store",
			Description:  sql.NullString{String: "Toko perlengkapan otomotif terpercaya.", Valid: true},
			Address:      sql.NullString{String: "Jl. Otomotif No.6, Medan", Valid: true},
			ContactEmail: sql.NullString{String: "service@autoparts.com", Valid: true},
			ContactPhone: sql.NullString{String: "088901234567", Valid: true},
			Status:       "active",
		},
	}

	for _, merchant := range merchants {
		if _, err := r.db.CreateMerchant(r.ctx, merchant); err != nil {
			r.logger.Error("failed to seed merchant", zap.Error(err))
			return err
		}
	}

	r.logger.Info("merchant succesfully seeded")

	return nil
}
