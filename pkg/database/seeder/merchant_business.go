package seeder

import (
	"context"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/logger"

	"go.uber.org/zap"
)

type merchantBusinessSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantBusinessSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *merchantBusinessSeeder {
	return &merchantBusinessSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *merchantBusinessSeeder) Seed() error {
	businessInfos := []db.CreateMerchantBusinessInformationParams{
		{
			MerchantID:        1,
			BusinessType:      strPtr("Toko Elektronik"),
			TaxID:             strPtr("01.234.567.8-999.000"),
			EstablishedYear:   int32Ptr(2010),
			NumberOfEmployees: int32Ptr(15),
			WebsiteUrl:        strPtr("https://technostore.com"),
		},
		{
			MerchantID:        2,
			BusinessType:      strPtr("Produk Kecantikan"),
			TaxID:             strPtr("02.345.678.9-888.000"),
			EstablishedYear:   int32Ptr(2015),
			NumberOfEmployees: int32Ptr(10),
			WebsiteUrl:        strPtr("https://glowbeauty.id"),
		},
		{
			MerchantID:        3,
			BusinessType:      strPtr("Toko Makanan Organik"),
			TaxID:             strPtr("03.456.789.0-777.000"),
			EstablishedYear:   int32Ptr(2012),
			NumberOfEmployees: int32Ptr(20),
			WebsiteUrl:        strPtr("https://dapsehat.id"),
		},
		{
			MerchantID:        4,
			BusinessType:      strPtr("Retail Gadget"),
			TaxID:             strPtr("04.567.890.1-666.000"),
			EstablishedYear:   int32Ptr(2018),
			NumberOfEmployees: int32Ptr(8),
			WebsiteUrl:        strPtr("https://gadgethub.com"),
		},
		{
			MerchantID:        5,
			BusinessType:      strPtr("Produk Ibu & Bayi"),
			TaxID:             strPtr("05.678.901.2-555.000"),
			EstablishedYear:   int32Ptr(2019),
			NumberOfEmployees: int32Ptr(6),
			WebsiteUrl:        strPtr("https://bayiceria.id"),
		},
		{
			MerchantID:        6,
			BusinessType:      strPtr("Peralatan Olahraga"),
			TaxID:             strPtr("06.789.012.3-444.000"),
			EstablishedYear:   int32Ptr(2016),
			NumberOfEmployees: int32Ptr(12),
			WebsiteUrl:        strPtr("https://tokosehat.id"),
		},
		{
			MerchantID:        7,
			BusinessType:      strPtr("Gaming Store"),
			TaxID:             strPtr("07.890.123.4-333.000"),
			EstablishedYear:   int32Ptr(2020),
			NumberOfEmployees: int32Ptr(5),
			WebsiteUrl:        strPtr("https://gameworld.com"),
		},
		{
			MerchantID:        8,
			BusinessType:      strPtr("Aksesori Otomotif"),
			TaxID:             strPtr("08.901.234.5-222.000"),
			EstablishedYear:   int32Ptr(2013),
			NumberOfEmployees: int32Ptr(9),
			WebsiteUrl:        strPtr("https://otomotifmart.com"),
		},
	}

	for _, info := range businessInfos {
		if _, err := r.db.CreateMerchantBusinessInformation(r.ctx, info); err != nil {
			r.logger.Error("failed to seed merchant business info", zap.Error(err))
			return err
		}
	}

	r.logger.Info("merchant business successfully seeded")
	return nil
}

func int32Ptr(v int32) *int32 {
	return &v
}
