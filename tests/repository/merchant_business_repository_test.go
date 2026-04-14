package repository_test

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	"ecommerce/tests"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type MerchantBusinessRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   repository.MerchantBusinessRepository
	merchantRepo repository.MerchantRepository
	userRepo repository.UserRepository
}

func (s *MerchantBusinessRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = repository.NewMerchantBusinessRepository(queries)
	s.merchantRepo = repository.NewMerchantRepository(queries)
	s.userRepo = repository.NewUserRepository(queries)
}

func (s *MerchantBusinessRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *MerchantBusinessRepositoryTestSuite) TestMerchantBusinessLifecycle() {
	ctx := context.Background()

	// 0. Setup: Create User and Merchant
	userReq := &requests.CreateUserRequest{
		FirstName: "Business",
		LastName:  "Owner",
		Email:     "business@example.com",
		Password:  "password123",
	}
	user, _ := s.userRepo.CreateUser(ctx, userReq)
	
	merchantReq := &requests.CreateMerchantRequest{
		UserID:       int(user.UserID),
		Name:         "Business Merchant",
		Description:  "Business description",
		Address:      "Address",
		ContactEmail: "business@email.com",
		ContactPhone: "08123456789",
		Status:       "active",
	}
	merchant, _ := s.merchantRepo.CreateMerchant(ctx, merchantReq)
	merchantID := int(merchant.MerchantID)

	// 1. Create Merchant Business Info
	createReq := &requests.CreateMerchantBusinessInformationRequest{
		MerchantID:        merchantID,
		BusinessType:      "PT",
		TaxID:             "123-456-789",
		EstablishedYear:   2020,
		NumberOfEmployees: 50,
		WebsiteUrl:        "https://business.com",
	}

	business, err := s.repo.CreateMerchantBusiness(ctx, createReq)
	s.NoError(err)
	s.NotNil(business)
	s.Equal(createReq.TaxID, *business.TaxID)

	businessID := int(business.MerchantBusinessInfoID)

	// 2. Find By ID
	found, err := s.repo.FindById(ctx, businessID)
	s.NoError(err)
	s.NotNil(found)

	// 3. Update Merchant Business Info
	updateReq := &requests.UpdateMerchantBusinessInformationRequest{
		MerchantBusinessInfoID: &businessID,
		BusinessType:           "CV",
		TaxID:                  "987-654-321",
		EstablishedYear:        2021,
		NumberOfEmployees:      60,
		WebsiteUrl:             "https://business-updated.com",
	}

	updated, err := s.repo.UpdateMerchantBusiness(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(updateReq.TaxID, *updated.TaxID)

	// 4. Trash
	trashed, err := s.repo.TrashedMerchantBusiness(ctx, businessID)
	s.NoError(err)
	s.NotNil(trashed)

	// 5. Restore
	restored, err := s.repo.RestoreMerchantBusiness(ctx, businessID)
	s.NoError(err)
	s.NotNil(restored)

	// 6. Delete Permanent
	_, _ = s.repo.TrashedMerchantBusiness(ctx, businessID)
	success, err := s.repo.DeleteMerchantBusinessPermanent(ctx, businessID)
	s.NoError(err)
	s.True(success)
}

func TestMerchantBusinessRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantBusinessRepositoryTestSuite))
}
