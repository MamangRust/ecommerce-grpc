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

type ShippingAddressRepositoryTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	dbPool     *pgxpool.Pool
	repos      *repository.Repositories
	userID     int
	merchantID int
	orderID    int
}

func (s *ShippingAddressRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repos = repository.NewRepositories(queries)

	ctx := context.Background()

	// 1. Create User
	user, err := s.repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Shipping",
		LastName:  "User",
		Email:     "shipping.user@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	// 2. Create Merchant
	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      s.userID,
		Name:        "Shipping Merchant",
		Description: "A merchant for shipping tests",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	// 3. Create Order
	order, err := s.repos.Order.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		UserID:     s.userID,
		MerchantID: s.merchantID,
		TotalPrice: 1000,
	})
	s.Require().NoError(err)
	s.orderID = int(order.OrderID)
}

func (s *ShippingAddressRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *ShippingAddressRepositoryTestSuite) TestShippingAddressLifecycle() {
	ctx := context.Background()

	// 1. Create Shipping Address
	createReq := &requests.CreateShippingAddressRequest{
		OrderID:        &s.orderID,
		Alamat:         "Jl. Merdeka No. 1",
		Provinsi:       "DKI Jakarta",
		Kota:           "Jakarta Pusat",
		Negara:         "Indonesia",
		Courier:        "JNE",
		ShippingMethod: "Reguler",
		ShippingCost:   15000,
	}

	address, err := s.repos.Shipping.CreateShippingAddress(ctx, createReq)
	s.NoError(err)
	s.NotNil(address)
	s.Equal(createReq.Alamat, address.Alamat)
	shippingID := int(address.ShippingAddressID)

	// 2. Find By ID
	found, err := s.repos.Shipping.FindById(ctx, shippingID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(address.Alamat, found.Alamat)

	// 3. Find By Order
	foundByOrder, err := s.repos.Shipping.FindByOrder(ctx, s.orderID)
	s.NoError(err)
	s.NotNil(foundByOrder)
	s.Equal(address.Alamat, foundByOrder.Alamat)

	// 4. Update Shipping Address
	updateReq := &requests.UpdateShippingAddressRequest{
		ShippingID:     &shippingID,
		Alamat:         "Jl. Merdeka No. 2 (Updated)",
		Provinsi:       "DKI Jakarta",
		Kota:           "Jakarta Selatan",
		Negara:         "Indonesia",
		Courier:        "TIKI",
		ShippingMethod: "Overnight",
		ShippingCost:   25000,
	}

	updated, err := s.repos.Shipping.UpdateShippingAddress(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(updateReq.Alamat, updated.Alamat)

	// 5. Trash Shipping Address
	trashed, err := s.repos.Shipping.TrashShippingAddress(ctx, shippingID)
	s.NoError(err)
	s.NotNil(trashed)
	s.True(trashed.DeletedAt.Valid)

	// 6. Find By Trashed
	trashedList, err := s.repos.Shipping.FindByTrashed(ctx, &requests.FindAllShippingAddress{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(trashedList)

	// 7. Restore Shipping Address
	restored, err := s.repos.Shipping.RestoreShippingAddress(ctx, shippingID)
	s.NoError(err)
	s.NotNil(restored)
	s.False(restored.DeletedAt.Valid)

	// 8. Delete Permanent
	// Trash again first
	_, err = s.repos.Shipping.TrashShippingAddress(ctx, shippingID)
	s.NoError(err)

	success, err := s.repos.Shipping.DeleteShippingAddressPermanently(ctx, shippingID)
	s.NoError(err)
	s.True(success)

	// 9. Verify it's gone
	_, err = s.repos.Shipping.FindById(ctx, shippingID)
	s.Error(err)
}

func TestShippingAddressRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ShippingAddressRepositoryTestSuite))
}
