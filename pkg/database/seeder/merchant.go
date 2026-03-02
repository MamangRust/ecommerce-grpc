package seeder

import (
	"context"
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
			Description:  strPtr("Toko elektronik terpercaya dengan berbagai produk gadget dan aksesoris."),
			Address:      strPtr("Jl. Teknologi No.1, Jakarta"),
			ContactEmail: strPtr("support@elektronikstore.com"),
			ContactPhone: strPtr("081234567890"),
			Status:       "active",
		},
		{
			UserID:       2,
			Name:         "Kecantikan Sehat",
			Description:  strPtr("Produk skincare dan kesehatan pilihan."),
			Address:      strPtr("Jl. Kesehatan No.5, Bandung"),
			ContactEmail: strPtr("cs@kecantikansehat.com"),
			ContactPhone: strPtr("082345678901"),
			Status:       "active",
		},
		{
			UserID:       3,
			Name:         "Rumah Indah",
			Description:  strPtr("Peralatan rumah tangga berkualitas dan estetik."),
			Address:      strPtr("Jl. Rumah No.12, Surabaya"),
			ContactEmail: strPtr("info@rumahindah.com"),
			ContactPhone: strPtr("083456789012"),
			Status:       "active",
		},
		{
			UserID:       4,
			Name:         "Mom & Baby Care",
			Description:  strPtr("Semua kebutuhan ibu dan bayi ada di sini."),
			Address:      strPtr("Jl. Keluarga No.7, Depok"),
			ContactEmail: strPtr("support@momandbaby.com"),
			ContactPhone: strPtr("084567890123"),
			Status:       "active",
		},
		{
			UserID:       5,
			Name:         "Sport Zone",
			Description:  strPtr("Perlengkapan olahraga dan outdoor terlengkap."),
			Address:      strPtr("Jl. Atletik No.3, Yogyakarta"),
			ContactEmail: strPtr("halo@sportzone.com"),
			ContactPhone: strPtr("085678901234"),
			Status:       "active",
		},
		{
			UserID:       6,
			Name:         "Fresh Mart",
			Description:  strPtr("Toko makanan dan minuman segar dan kemasan."),
			Address:      strPtr("Jl. Pasar No.10, Semarang"),
			ContactEmail: strPtr("fresh@mart.com"),
			ContactPhone: strPtr("086789012345"),
			Status:       "active",
		},
		{
			UserID:       7,
			Name:         "Gamer Heaven",
			Description:  strPtr("Game, console, dan aksesori lengkap untuk gamers."),
			Address:      strPtr("Jl. Game No.8, Bekasi"),
			ContactEmail: strPtr("gamer@heaven.com"),
			ContactPhone: strPtr("087890123456"),
			Status:       "active",
		},
		{
			UserID:       8,
			Name:         "AutoParts Store",
			Description:  strPtr("Toko perlengkapan otomotif terpercaya."),
			Address:      strPtr("Jl. Otomotif No.6, Medan"),
			ContactEmail: strPtr("service@autoparts.com"),
			ContactPhone: strPtr("088901234567"),
			Status:       "active",
		},
	}

	for _, merchant := range merchants {
		if _, err := r.db.CreateMerchant(r.ctx, merchant); err != nil {
			r.logger.Error(
				"failed to seed merchant",
				zap.Int32("user_id", merchant.UserID),
				zap.String("name", merchant.Name),
				zap.Error(err),
			)
			return err
		}
	}

	r.logger.Info("merchant successfully seeded")
	return nil
}
