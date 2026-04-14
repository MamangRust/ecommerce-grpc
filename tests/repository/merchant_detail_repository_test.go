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

type MerchantDetailRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   repository.MerchantDetailRepository
	merchantRepo repository.MerchantRepository
	userRepo repository.UserRepository
}

func (s *MerchantDetailRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = repository.NewMerchantDetailRepository(queries)
	s.merchantRepo = repository.NewMerchantRepository(queries)
	s.userRepo = repository.NewUserRepository(queries)
}

func (s *MerchantDetailRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *MerchantDetailRepositoryTestSuite) TestMerchantDetailLifecycle() {
	ctx := context.Background()

	// 0. Setup: Create User and Merchant
	userReq := &requests.CreateUserRequest{
		FirstName: "Detail",
		LastName:  "Owner",
		Email:     "detail@example.com",
		Password:  "password123",
	}
	user, _ := s.userRepo.CreateUser(ctx, userReq)
	
	merchantReq := &requests.CreateMerchantRequest{
		UserID:       int(user.UserID),
		Name:         "Detail Merchant",
		Description:  "Detail description",
		Address:      "Address",
		ContactEmail: "detail@email.com",
		ContactPhone: "08123456789",
		Status:       "active",
	}
	merchant, _ := s.merchantRepo.CreateMerchant(ctx, merchantReq)
	merchantID := int(merchant.MerchantID)

	// 1. Create Merchant Detail
	createReq := &requests.CreateMerchantDetailRequest{
		MerchantID:       merchantID,
		DisplayName:      "Merchant Display Name",
		CoverImageUrl:    "https://example.com/cover.jpg",
		LogoUrl:          "https://example.com/logo.jpg",
		ShortDescription: "Short desc",
		WebsiteUrl:       "https://merchant.com",
	}

	detail, err := s.repo.CreateMerchantDetail(ctx, createReq)
	s.NoError(err)
	s.NotNil(detail)
	s.Equal(createReq.DisplayName, *detail.DisplayName)

	detailID := int(detail.MerchantDetailID)

	// 2. Find By ID
	found, err := s.repo.FindById(ctx, detailID)
	s.NoError(err)
	s.NotNil(found)

	// 3. Update Merchant Detail
	updateReq := &requests.UpdateMerchantDetailRequest{
		MerchantDetailID: &detailID,
		DisplayName:      "Merchant Display Name Updated",
		CoverImageUrl:    "https://example.com/cover-updated.jpg",
		LogoUrl:          "https://example.com/logo-updated.jpg",
		ShortDescription: "Short desc updated",
		WebsiteUrl:       "https://merchant-updated.com",
	}

	updated, err := s.repo.UpdateMerchantDetail(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(updateReq.DisplayName, *updated.DisplayName)

	// 4. Trash
	trashed, err := s.repo.TrashedMerchantDetail(ctx, detailID)
	s.NoError(err)
	s.NotNil(trashed)

	// 5. Restore
	restored, err := s.repo.RestoreMerchantDetail(ctx, detailID)
	s.NoError(err)
	s.NotNil(restored)

	// 6. Delete Permanent
	_, _ = s.repo.TrashedMerchantDetail(ctx, detailID)
	success, err := s.repo.DeleteMerchantDetailPermanent(ctx, detailID)
	s.NoError(err)
	s.True(success)
}

func TestMerchantDetailRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantDetailRepositoryTestSuite))
}
