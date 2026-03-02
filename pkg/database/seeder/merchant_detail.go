package seeder

import (
	"context"
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
			DisplayName:      strPtr("Techno Store"),
			CoverImageUrl:    strPtr("cover/techno.jpg"),
			LogoUrl:          strPtr("logo/techno.png"),
			ShortDescription: strPtr("Pusat elektronik terpercaya sejak 2010"),
			WebsiteUrl:       strPtr("https://technostore.com"),
		},
		{
			MerchantID:       2,
			DisplayName:      strPtr("Glow Beauty"),
			CoverImageUrl:    strPtr("cover/beauty.jpg"),
			LogoUrl:          strPtr("logo/beauty.png"),
			ShortDescription: strPtr("Produk kecantikan alami dan aman"),
			WebsiteUrl:       strPtr("https://glowbeauty.id"),
		},
		{
			MerchantID:       3,
			DisplayName:      strPtr("Dapur Sehat"),
			CoverImageUrl:    strPtr("cover/dapur.jpg"),
			LogoUrl:          strPtr("logo/dapur.png"),
			ShortDescription: strPtr("Makanan sehat dan organik"),
			WebsiteUrl:       strPtr("https://dapsehat.id"),
		},
		{
			MerchantID:       4,
			DisplayName:      strPtr("Gadget Hub"),
			CoverImageUrl:    strPtr("cover/gadget.jpg"),
			LogoUrl:          strPtr("logo/gadget.png"),
			ShortDescription: strPtr("Semua tentang gadget terbaru"),
			WebsiteUrl:       strPtr("https://gadgethub.com"),
		},
		{
			MerchantID:       5,
			DisplayName:      strPtr("Bayi Ceria"),
			CoverImageUrl:    strPtr("cover/bayi.jpg"),
			LogoUrl:          strPtr("logo/bayi.png"),
			ShortDescription: strPtr("Produk terbaik untuk si kecil"),
			WebsiteUrl:       strPtr("https://bayiceria.id"),
		},
		{
			MerchantID:       6,
			DisplayName:      strPtr("Toko Sehat"),
			CoverImageUrl:    strPtr("cover/sehat.jpg"),
			LogoUrl:          strPtr("logo/sehat.png"),
			ShortDescription: strPtr("Peralatan olahraga lengkap"),
			WebsiteUrl:       strPtr("https://tokosehat.id"),
		},
		{
			MerchantID:       7,
			DisplayName:      strPtr("Game World"),
			CoverImageUrl:    strPtr("cover/game.jpg"),
			LogoUrl:          strPtr("logo/game.png"),
			ShortDescription: strPtr("Konsol dan game terbaik"),
			WebsiteUrl:       strPtr("https://gameworld.com"),
		},
		{
			MerchantID:       8,
			DisplayName:      strPtr("Otomotif Mart"),
			CoverImageUrl:    strPtr("cover/otomotif.jpg"),
			LogoUrl:          strPtr("logo/otomotif.png"),
			ShortDescription: strPtr("Aksesori kendaraan terpercaya"),
			WebsiteUrl:       strPtr("https://otomotifmart.com"),
		},
	}

	for i, detail := range details {
		if _, err := r.db.CreateMerchantDetail(r.ctx, detail); err != nil {
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
