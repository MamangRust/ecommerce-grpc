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

type MerchantPolicyRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   repository.MerchantPoliciesRepository
	merchantRepo repository.MerchantRepository
	userRepo repository.UserRepository
}

func (s *MerchantPolicyRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = repository.NewMerchantPoliciesRepository(queries)
	s.merchantRepo = repository.NewMerchantRepository(queries)
	s.userRepo = repository.NewUserRepository(queries)
}

func (s *MerchantPolicyRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *MerchantPolicyRepositoryTestSuite) TestMerchantPolicyLifecycle() {
	ctx := context.Background()

	// 0. Setup: Create User and Merchant
	userReq := &requests.CreateUserRequest{
		FirstName: "Policy",
		LastName:  "Owner",
		Email:     "policy@example.com",
		Password:  "password123",
	}
	user, _ := s.userRepo.CreateUser(ctx, userReq)
	
	merchantReq := &requests.CreateMerchantRequest{
		UserID:       int(user.UserID),
		Name:         "Policy Merchant",
		Description:  "Policy description",
		Address:      "Address",
		ContactEmail: "policy@email.com",
		ContactPhone: "08123456789",
		Status:       "active",
	}
	merchant, _ := s.merchantRepo.CreateMerchant(ctx, merchantReq)
	merchantID := int(merchant.MerchantID)

	// 1. Create Merchant Policy
	createReq := &requests.CreateMerchantPolicyRequest{
		MerchantID:  merchantID,
		PolicyType:  "Shipping",
		Title:       "Free Shipping",
		Description: "Free shipping over $100",
	}

	policy, err := s.repo.CreateMerchantPolicy(ctx, createReq)
	s.NoError(err)
	s.NotNil(policy)
	s.Equal(createReq.Title, policy.Title)

	policyID := int(policy.MerchantPolicyID)

	// 2. Find By ID
	found, err := s.repo.FindById(ctx, policyID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(createReq.Title, found.Title)

	// 3. Update Merchant Policy
	updateReq := &requests.UpdateMerchantPolicyRequest{
		MerchantPolicyID: &policyID,
		PolicyType:       "Shipping",
		Title:            "Standard Shipping",
		Description:      "Flat rate shipping",
	}

	updated, err := s.repo.UpdateMerchantPolicy(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(updateReq.Title, updated.Title)

	// 4. Trash
	trashed, err := s.repo.TrashedMerchantPolicy(ctx, policyID)
	s.NoError(err)
	s.NotNil(trashed)

	// 5. Restore
	restored, err := s.repo.RestoreMerchantPolicy(ctx, policyID)
	s.NoError(err)
	s.NotNil(restored)

	// 6. Delete Permanent
	_, _ = s.repo.TrashedMerchantPolicy(ctx, policyID)
	success, err := s.repo.DeleteMerchantPolicyPermanent(ctx, policyID)
	s.NoError(err)
	s.True(success)
}

func TestMerchantPolicyRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantPolicyRepositoryTestSuite))
}
