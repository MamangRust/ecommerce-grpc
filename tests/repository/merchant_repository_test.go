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

type MerchantRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   repository.MerchantRepository
	userRepo repository.UserRepository
}

func (s *MerchantRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = repository.NewMerchantRepository(queries)
	s.userRepo = repository.NewUserRepository(queries)
}

func (s *MerchantRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *MerchantRepositoryTestSuite) TestMerchantLifecycle() {
	ctx := context.Background()

	// 0. Create User first
	userReq := &requests.CreateUserRequest{
		FirstName: "Merchant",
		LastName:  "Owner",
		Email:     "owner@example.com",
		Password:  "password123",
	}
	user, err := s.userRepo.CreateUser(ctx, userReq)
	s.NoError(err)
	s.NotNil(user)
	userID := int(user.UserID)

	// 1. Create Merchant
	createReq := &requests.CreateMerchantRequest{
		UserID:       userID,
		Name:         "Gopay Merchant",
		Description:  "Fast payment merchant",
		Address:      "Jakarta, Indonesia",
		ContactEmail: "gopay@example.com",
		ContactPhone: "08123456789",
		Status:       "active",
	}

	merchant, err := s.repo.CreateMerchant(ctx, createReq)
	s.NoError(err)
	s.NotNil(merchant)
	s.Equal(createReq.Name, merchant.Name)

	merchantID := int(merchant.MerchantID)

	// 2. Find By ID
	found, err := s.repo.FindById(ctx, merchantID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(merchant.Name, found.Name)

	// 3. Update Merchant
	updateReq := &requests.UpdateMerchantRequest{
		MerchantID:   &merchantID,
		UserID:       userID,
		Name:         "Gopay Merchant Updated",
		Description:  "Faster payment merchant",
		Address:      "Bandung, Indonesia",
		ContactEmail: "gopay-updated@example.com",
		ContactPhone: "08987654321",
		Status:       "active",
	}

	updated, err := s.repo.UpdateMerchant(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(updateReq.Name, updated.Name)

	// 4. Find All
	all, err := s.repo.FindAllMerchants(ctx, &requests.FindAllMerchant{
		Search:   "Gopay",
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(all)

	// 5. Trash Merchant
	trashed, err := s.repo.TrashedMerchant(ctx, merchantID)
	s.NoError(err)
	s.NotNil(trashed)
	s.True(trashed.DeletedAt.Valid)

	// 6. Find By Trashed
	trashedList, err := s.repo.FindByTrashed(ctx, &requests.FindAllMerchant{
		Search:   "Gopay",
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(trashedList)

	// 7. Restore Merchant
	restored, err := s.repo.RestoreMerchant(ctx, merchantID)
	s.NoError(err)
	s.NotNil(restored)
	s.False(restored.DeletedAt.Valid)

	// 8. Delete Permanent
	// Trash again first
	_, err = s.repo.TrashedMerchant(ctx, merchantID)
	s.NoError(err)

	success, err := s.repo.DeleteMerchantPermanent(ctx, merchantID)
	s.NoError(err)
	s.True(success)

	// 9. Verify it's gone
	_, err = s.repo.FindById(ctx, merchantID)
	s.Error(err)
}

func TestMerchantRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantRepositoryTestSuite))
}
