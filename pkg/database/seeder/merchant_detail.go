package seeder

import (
	"context"
	"database/sql"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/logger"
	"fmt"

	"go.uber.org/zap"
)

type merchantDetailSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantDetailSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *merchantDetailSeeder {
	return &merchantDetailSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *merchantDetailSeeder) Seed() error {
	details := []db.CreateMerchantDetailParams{
		{
			MerchantID:       1,
			DisplayName:      sql.NullString{String: "Techno Store", Valid: true},
			CoverImageUrl:    sql.NullString{String: "cover/techno.jpg", Valid: true},
			LogoUrl:          sql.NullString{String: "logo/techno.png", Valid: true},
			ShortDescription: sql.NullString{String: "Pusat elektronik terpercaya sejak 2010", Valid: true},
			WebsiteUrl:       sql.NullString{String: "https://technostore.com", Valid: true},
		},
		{
			MerchantID:       2,
			DisplayName:      sql.NullString{String: "Glow Beauty", Valid: true},
			CoverImageUrl:    sql.NullString{String: "cover/beauty.jpg", Valid: true},
			LogoUrl:          sql.NullString{String: "logo/beauty.png", Valid: true},
			ShortDescription: sql.NullString{String: "Produk kecantikan alami dan aman", Valid: true},
			WebsiteUrl:       sql.NullString{String: "https://glowbeauty.id", Valid: true},
		},
		{
			MerchantID:       3,
			DisplayName:      sql.NullString{String: "Dapur Sehat", Valid: true},
			CoverImageUrl:    sql.NullString{String: "cover/dapur.jpg", Valid: true},
			LogoUrl:          sql.NullString{String: "logo/dapur.png", Valid: true},
			ShortDescription: sql.NullString{String: "Makanan sehat dan organik", Valid: true},
			WebsiteUrl:       sql.NullString{String: "https://dapsehat.id", Valid: true},
		},
		{
			MerchantID:       4,
			DisplayName:      sql.NullString{String: "Gadget Hub", Valid: true},
			CoverImageUrl:    sql.NullString{String: "cover/gadget.jpg", Valid: true},
			LogoUrl:          sql.NullString{String: "logo/gadget.png", Valid: true},
			ShortDescription: sql.NullString{String: "Semua tentang gadget terbaru", Valid: true},
			WebsiteUrl:       sql.NullString{String: "https://gadgethub.com", Valid: true},
		},
		{
			MerchantID:       5,
			DisplayName:      sql.NullString{String: "Bayi Ceria", Valid: true},
			CoverImageUrl:    sql.NullString{String: "cover/bayi.jpg", Valid: true},
			LogoUrl:          sql.NullString{String: "logo/bayi.png", Valid: true},
			ShortDescription: sql.NullString{String: "Produk terbaik untuk si kecil", Valid: true},
			WebsiteUrl:       sql.NullString{String: "https://bayiceria.id", Valid: true},
		},
		{
			MerchantID:       6,
			DisplayName:      sql.NullString{String: "Toko Sehat", Valid: true},
			CoverImageUrl:    sql.NullString{String: "cover/sehat.jpg", Valid: true},
			LogoUrl:          sql.NullString{String: "logo/sehat.png", Valid: true},
			ShortDescription: sql.NullString{String: "Peralatan olahraga lengkap", Valid: true},
			WebsiteUrl:       sql.NullString{String: "https://tokosehat.id", Valid: true},
		},
		{
			MerchantID:       7,
			DisplayName:      sql.NullString{String: "Game World", Valid: true},
			CoverImageUrl:    sql.NullString{String: "cover/game.jpg", Valid: true},
			LogoUrl:          sql.NullString{String: "logo/game.png", Valid: true},
			ShortDescription: sql.NullString{String: "Konsol dan game terbaik", Valid: true},
			WebsiteUrl:       sql.NullString{String: "https://gameworld.com", Valid: true},
		},
		{
			MerchantID:       8,
			DisplayName:      sql.NullString{String: "Otomotif Mart", Valid: true},
			CoverImageUrl:    sql.NullString{String: "cover/otomotif.jpg", Valid: true},
			LogoUrl:          sql.NullString{String: "logo/otomotif.png", Valid: true},
			ShortDescription: sql.NullString{String: "Aksesori kendaraan terpercaya", Valid: true},
			WebsiteUrl:       sql.NullString{String: "https://otomotifmart.com", Valid: true},
		},
	}

	for i, detail := range details {
		_, err := r.db.CreateMerchantDetail(r.ctx, detail)
		if err != nil {
			r.logger.Error("failed to seed merchant detail", zap.Error(err))
			return err
		}

		merchantDetailID := int32(i + 1)
		socialMedia := []db.CreateMerchantSocialMediaLinkParams{
			{
				MerchantDetailID: merchantDetailID,
				Platform:         "Facebook",
				Url:              "https://www.facebook.com/merchant" + fmt.Sprint(merchantDetailID),
			},
			{
				MerchantDetailID: merchantDetailID,
				Platform:         "Instagram",
				Url:              "https://www.instagram.com/merchant" + fmt.Sprint(merchantDetailID),
			},
			{
				MerchantDetailID: merchantDetailID,
				Platform:         "Twitter",
				Url:              "https://www.twitter.com/merchant" + fmt.Sprint(merchantDetailID),
			},
		}

		for _, sm := range socialMedia {
			if _, err := r.db.CreateMerchantSocialMediaLink(r.ctx, sm); err != nil {
				r.logger.Error("failed to seed merchant social media link", zap.Error(err))
				return err
			}
		}
	}

	r.logger.Info("merchant detail & merchant social link successfully seeded")

	return nil
}
