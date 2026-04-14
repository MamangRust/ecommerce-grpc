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

type OrderRepositoryTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	dbPool     *pgxpool.Pool
	repos      *repository.Repositories
	merchantID int
	userID     int
}

func (s *OrderRepositoryTestSuite) SetupSuite() {
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
		FirstName: "Order",
		LastName:  "User",
		Email:     "order.user@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	// 2. Create Merchant
	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      s.userID,
		Name:        "Order Merchant",
		Description: "A merchant for testing orders",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)
}

func (s *OrderRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *OrderRepositoryTestSuite) TestOrderLifecycle() {
	ctx := context.Background()

	// 1. Create Order
	createReq := &requests.CreateOrderRecordRequest{
		MerchantID: s.merchantID,
		UserID:     s.userID,
		TotalPrice: 500,
	}

	order, err := s.repos.Order.CreateOrder(ctx, createReq)
	s.NoError(err)
	s.NotNil(order)
	s.Equal(int32(createReq.TotalPrice), order.TotalPrice)
	orderID := int(order.OrderID)

	// 2. Find By ID
	found, err := s.repos.Order.FindById(ctx, orderID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(order.TotalPrice, found.TotalPrice)

	// 3. Find All
	orders, err := s.repos.Order.FindAllOrders(ctx, &requests.FindAllOrder{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(orders)

	// 4. Update Order
	updateReq := &requests.UpdateOrderRecordRequest{
		OrderID:    orderID,
		MerchantID: s.merchantID,
		UserID:     s.userID,
		TotalPrice: 750,
	}

	updated, err := s.repos.Order.UpdateOrder(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(updateReq.TotalPrice), updated.TotalPrice)

	// 5. Trash Order
	trashed, err := s.repos.Order.TrashedOrder(ctx, orderID)
	s.NoError(err)
	s.NotNil(trashed)
	s.True(trashed.DeletedAt.Valid)

	// 6. Find By Trashed
	trashedList, err := s.repos.Order.FindByTrashed(ctx, &requests.FindAllOrder{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(trashedList)

	// 7. Restore Order
	restored, err := s.repos.Order.RestoreOrder(ctx, orderID)
	s.NoError(err)
	s.NotNil(restored)
	s.False(restored.DeletedAt.Valid)

	// 8. Delete Permanent
	// Trash again first
	_, err = s.repos.Order.TrashedOrder(ctx, orderID)
	s.NoError(err)

	success, err := s.repos.Order.DeleteOrderPermanent(ctx, orderID)
	s.NoError(err)
	s.True(success)

	// 9. Verify it's gone
	_, err = s.repos.Order.FindById(ctx, orderID)
	s.Error(err)
}

func TestOrderRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderRepositoryTestSuite))
}
