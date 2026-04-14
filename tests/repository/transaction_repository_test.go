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

type TransactionRepositoryTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	dbPool     *pgxpool.Pool
	repos      *repository.Repositories
	merchantID int
	userID     int
	orderID    int
}

func (s *TransactionRepositoryTestSuite) SetupSuite() {
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
		FirstName: "Transaction",
		LastName:  "User",
		Email:     "trans.user@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	// 2. Create Merchant
	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      s.userID,
		Name:        "Transaction Merchant",
		Description: "A merchant for testing transactions",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	// 3. Create Order
	order, err := s.repos.Order.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: s.merchantID,
		UserID:     s.userID,
		TotalPrice: 1000,
	})
	s.Require().NoError(err)
	s.orderID = int(order.OrderID)
}

func (s *TransactionRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *TransactionRepositoryTestSuite) TestTransactionLifecycle() {
	ctx := context.Background()
	statusSuccess := "success"

	// 1. Create Transaction
	createReq := &requests.CreateTransactionRequest{
		OrderID:       s.orderID,
		MerchantID:    s.merchantID,
		PaymentMethod: "credit_card",
		Amount:        1000,
		PaymentStatus: &statusSuccess,
	}

	trans, err := s.repos.Transaction.CreateTransaction(ctx, createReq)
	s.NoError(err)
	s.NotNil(trans)
	s.Equal(int32(createReq.Amount), trans.Amount)
	transID := int(trans.TransactionID)

	// 2. Find By ID
	found, err := s.repos.Transaction.FindById(ctx, transID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(trans.Amount, found.Amount)

	// 3. Find By Order ID
	foundByOrder, err := s.repos.Transaction.FindByOrderId(ctx, s.orderID)
	s.NoError(err)
	s.NotNil(foundByOrder)
	s.Equal(int32(transID), foundByOrder.TransactionID)

	// 4. Update Transaction
	statusFailed := "failed"
	updateReq := &requests.UpdateTransactionRequest{
		TransactionID: &transID,
		OrderID:       s.orderID,
		MerchantID:    s.merchantID,
		PaymentMethod: "debit_card",
		Amount:        1200,
		PaymentStatus: &statusFailed,
	}

	updated, err := s.repos.Transaction.UpdateTransaction(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(updateReq.Amount), updated.Amount)
	s.Equal("debit_card", updated.PaymentMethod)

	// 5. Trash Transaction
	trashed, err := s.repos.Transaction.TrashTransaction(ctx, transID)
	s.NoError(err)
	s.NotNil(trashed)
	s.True(trashed.DeletedAt.Valid)

	// 6. Find By Trashed
	trashedList, err := s.repos.Transaction.FindByTrashed(ctx, &requests.FindAllTransaction{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.NotEmpty(trashedList)

	// 7. Restore Transaction
	restored, err := s.repos.Transaction.RestoreTransaction(ctx, transID)
	s.NoError(err)
	s.NotNil(restored)
	s.False(restored.DeletedAt.Valid)

	// 8. Find All
	all, err := s.repos.Transaction.FindAllTransactions(ctx, &requests.FindAllTransaction{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.NotEmpty(all)

	// 9. Delete Permanent
	// Trash again
	_, err = s.repos.Transaction.TrashTransaction(ctx, transID)
	s.NoError(err)

	success, err := s.repos.Transaction.DeleteTransactionPermanently(ctx, transID)
	s.NoError(err)
	s.True(success)

	// 10. Verify it's gone
	_, err = s.repos.Transaction.FindById(ctx, transID)
	s.Error(err)
}

func (s *TransactionRepositoryTestSuite) TestTransactionReporting() {
	ctx := context.Background()
	statusSuccess := "success"

	// Create a transaction to have some data
	_, err := s.repos.Transaction.CreateTransaction(ctx, &requests.CreateTransactionRequest{
		OrderID:       s.orderID,
		MerchantID:    s.merchantID,
		PaymentMethod: "credit_card",
		Amount:        5000,
		PaymentStatus: &statusSuccess,
	})
	s.Require().NoError(err)

	// 1. Get Monthly Amount Success
	monthly, err := s.repos.Transaction.GetMonthlyAmountSuccess(ctx, &requests.MonthAmountTransaction{
		Year:  2026,
		Month: 4,
	})
	s.NoError(err)
	// Even if it's empty in this month (unless we seed with explicit dates), it shouldn't error.
	s.NotNil(monthly)

	// 2. Get Monthly Amount Success By Merchant
	monthlyMerchant, err := s.repos.Transaction.GetMonthlyAmountSuccessByMerchant(ctx, &requests.MonthAmountTransactionMerchant{
		MerchantID: s.merchantID,
		Year:       2026,
		Month:      4,
	})
	s.NoError(err)
	s.NotNil(monthlyMerchant)
}

func TestTransactionRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionRepositoryTestSuite))
}
