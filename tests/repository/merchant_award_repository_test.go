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

type MerchantAwardRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   repository.MerchantAwardRepository
	merchantRepo repository.MerchantRepository
	userRepo repository.UserRepository
}

func (s *MerchantAwardRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = repository.NewMerchantAwardRepository(queries)
	s.merchantRepo = repository.NewMerchantRepository(queries)
	s.userRepo = repository.NewUserRepository(queries)
}

func (s *MerchantAwardRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *MerchantAwardRepositoryTestSuite) TestMerchantAwardLifecycle() {
	ctx := context.Background()

	// 0. Setup: Create User and Merchant
	userReq := &requests.CreateUserRequest{
		FirstName: "Award",
		LastName:  "Winner",
		Email:     "winner@example.com",
		Password:  "password123",
	}
	user, _ := s.userRepo.CreateUser(ctx, userReq)
	
	merchantReq := &requests.CreateMerchantRequest{
		UserID:       int(user.UserID),
		Name:         "Award Merchant",
		Description:  "Award description",
		Address:      "Address",
		ContactEmail: "winner@email.com",
		ContactPhone: "08123456789",
		Status:       "active",
	}
	merchant, _ := s.merchantRepo.CreateMerchant(ctx, merchantReq)
	merchantID := int(merchant.MerchantID)

	// 1. Create Merchant Award
	createReq := &requests.CreateMerchantCertificationOrAwardRequest{
		MerchantID:     merchantID,
		Title:          "Best Merchant 2024",
		Description:    "Best in class",
		IssuedBy:       "E-commerce Assocation",
		IssueDate:      "2024-01-01",
		ExpiryDate:     "2025-01-01",
		CertificateUrl: "https://award.com/cert.pdf",
	}

	award, err := s.repo.CreateMerchantAward(ctx, createReq)
	s.NoError(err)
	s.NotNil(award)
	s.Equal(createReq.Title, award.Title)

	awardID := int(award.MerchantCertificationID)

	// 2. Find By ID
	found, err := s.repo.FindById(ctx, awardID)
	s.NoError(err)
	s.NotNil(found)

	// 3. Update Merchant Award
	updateReq := &requests.UpdateMerchantCertificationOrAwardRequest{
		MerchantCertificationID: &awardID,
		Title:                   "Best Merchant 2024 Updated",
		Description:             "Best in class updated",
		IssuedBy:                "E-commerce Assocation",
		IssueDate:               "2024-02-01",
		ExpiryDate:              "2025-02-01",
		CertificateUrl:          "https://award.com/cert-updated.pdf",
	}

	updated, err := s.repo.UpdateMerchantAward(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(updateReq.Title, updated.Title)

	// 4. Trash
	trashed, err := s.repo.TrashedMerchantAward(ctx, awardID)
	s.NoError(err)
	s.NotNil(trashed)

	// 5. Restore
	restored, err := s.repo.RestoreMerchantAward(ctx, awardID)
	s.NoError(err)
	s.NotNil(restored)

	// 6. Delete Permanent
	_, _ = s.repo.TrashedMerchantAward(ctx, awardID)
	success, err := s.repo.DeleteMerchantPermanent(ctx, awardID)
	s.NoError(err)
	s.True(success)
}

func TestMerchantAwardRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantAwardRepositoryTestSuite))
}
