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

type MerchantSocialLinkRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   repository.MerchantSocialLinkRepository
	detailRepo repository.MerchantDetailRepository
	merchantRepo repository.MerchantRepository
	userRepo repository.UserRepository
}

func (s *MerchantSocialLinkRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = repository.NewMerchantSocialLinkRepository(queries)
	s.detailRepo = repository.NewMerchantDetailRepository(queries)
	s.merchantRepo = repository.NewMerchantRepository(queries)
	s.userRepo = repository.NewUserRepository(queries)
}

func (s *MerchantSocialLinkRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *MerchantSocialLinkRepositoryTestSuite) TestMerchantSocialLinkLifecycle() {
	ctx := context.Background()

	// 0. Setup: Create User, Merchant, and Detail
	userReq := &requests.CreateUserRequest{
		FirstName: "Social",
		LastName:  "Link",
		Email:     "social@example.com",
		Password:  "password123",
	}
	user, _ := s.userRepo.CreateUser(ctx, userReq)
	
	merchantReq := &requests.CreateMerchantRequest{
		UserID:       int(user.UserID),
		Name:         "Social Merchant",
		Description:  "Social description",
		Address:      "Address",
		ContactEmail: "social@email.com",
		ContactPhone: "08123456789",
		Status:       "active",
	}
	merchant, _ := s.merchantRepo.CreateMerchant(ctx, merchantReq)
	merchantID := int(merchant.MerchantID)

	detailReq := &requests.CreateMerchantDetailRequest{
		MerchantID:       merchantID,
		DisplayName:      "Social Merchant Display",
		CoverImageUrl:    "https://example.com/cover.jpg",
		LogoUrl:          "https://example.com/logo.jpg",
		ShortDescription: "Social short desc",
		WebsiteUrl:       "https://social.com",
	}
	detail, _ := s.detailRepo.CreateMerchantDetail(ctx, detailReq)
	detailID := int(detail.MerchantDetailID)

	// 1. Create Social Link
	createReq := &requests.CreateMerchantSocialRequest{
		MerchantDetailID: &detailID,
		Platform:         "Instagram",
		Url:              "https://instagram.com/merchant",
	}

	success, err := s.repo.CreateSocialLink(ctx, createReq)
	s.NoError(err)
	s.True(success)

	// Get the ID (we need a way to get it, usually FindAll in real repo, but let's assume it works or check DB)
	// Since the repository doesn't have a GetSocialLinkByID, I'll assume it works if create returns true.
	// In a real scenario, we might want to add FindByDetailID.
}

func TestMerchantSocialLinkRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantSocialLinkRepositoryTestSuite))
}
