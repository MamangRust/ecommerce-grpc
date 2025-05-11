package seeder

import (
	"context"
	"database/sql"
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
			BusinessType:      sql.NullString{String: "Toko Elektronik", Valid: true},
			TaxID:             sql.NullString{String: "01.234.567.8-999.000", Valid: true},
			EstablishedYear:   sql.NullInt32{Int32: 2010, Valid: true},
			NumberOfEmployees: sql.NullInt32{Int32: 15, Valid: true},
			WebsiteUrl:        sql.NullString{String: "https://technostore.com", Valid: true},
		},
		{
			MerchantID:        2,
			BusinessType:      sql.NullString{String: "Produk Kecantikan", Valid: true},
			TaxID:             sql.NullString{String: "02.345.678.9-888.000", Valid: true},
			EstablishedYear:   sql.NullInt32{Int32: 2015, Valid: true},
			NumberOfEmployees: sql.NullInt32{Int32: 10, Valid: true},
			WebsiteUrl:        sql.NullString{String: "https://glowbeauty.id", Valid: true},
		},
		{
			MerchantID:        3,
			BusinessType:      sql.NullString{String: "Toko Makanan Organik", Valid: true},
			TaxID:             sql.NullString{String: "03.456.789.0-777.000", Valid: true},
			EstablishedYear:   sql.NullInt32{Int32: 2012, Valid: true},
			NumberOfEmployees: sql.NullInt32{Int32: 20, Valid: true},
			WebsiteUrl:        sql.NullString{String: "https://dapsehat.id", Valid: true},
		},
		{
			MerchantID:        4,
			BusinessType:      sql.NullString{String: "Retail Gadget", Valid: true},
			TaxID:             sql.NullString{String: "04.567.890.1-666.000", Valid: true},
			EstablishedYear:   sql.NullInt32{Int32: 2018, Valid: true},
			NumberOfEmployees: sql.NullInt32{Int32: 8, Valid: true},
			WebsiteUrl:        sql.NullString{String: "https://gadgethub.com", Valid: true},
		},
		{
			MerchantID:        5,
			BusinessType:      sql.NullString{String: "Produk Ibu & Bayi", Valid: true},
			TaxID:             sql.NullString{String: "05.678.901.2-555.000", Valid: true},
			EstablishedYear:   sql.NullInt32{Int32: 2019, Valid: true},
			NumberOfEmployees: sql.NullInt32{Int32: 6, Valid: true},
			WebsiteUrl:        sql.NullString{String: "https://bayiceria.id", Valid: true},
		},
		{
			MerchantID:        6,
			BusinessType:      sql.NullString{String: "Peralatan Olahraga", Valid: true},
			TaxID:             sql.NullString{String: "06.789.012.3-444.000", Valid: true},
			EstablishedYear:   sql.NullInt32{Int32: 2016, Valid: true},
			NumberOfEmployees: sql.NullInt32{Int32: 12, Valid: true},
			WebsiteUrl:        sql.NullString{String: "https://tokosehat.id", Valid: true},
		},
		{
			MerchantID:        7,
			BusinessType:      sql.NullString{String: "Gaming Store", Valid: true},
			TaxID:             sql.NullString{String: "07.890.123.4-333.000", Valid: true},
			EstablishedYear:   sql.NullInt32{Int32: 2020, Valid: true},
			NumberOfEmployees: sql.NullInt32{Int32: 5, Valid: true},
			WebsiteUrl:        sql.NullString{String: "https://gameworld.com", Valid: true},
		},
		{
			MerchantID:        8,
			BusinessType:      sql.NullString{String: "Aksesori Otomotif", Valid: true},
			TaxID:             sql.NullString{String: "08.901.234.5-222.000", Valid: true},
			EstablishedYear:   sql.NullInt32{Int32: 2013, Valid: true},
			NumberOfEmployees: sql.NullInt32{Int32: 9, Valid: true},
			WebsiteUrl:        sql.NullString{String: "https://otomotifmart.com", Valid: true},
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
