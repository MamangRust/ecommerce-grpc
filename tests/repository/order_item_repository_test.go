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

type OrderItemRepositoryTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	dbPool     *pgxpool.Pool
	repos      *repository.Repositories
	merchantID int
	userID     int
	categoryID int
	productID  int
	orderID    int
}

func (s *OrderItemRepositoryTestSuite) SetupSuite() {
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
		FirstName: "OrderItem",
		LastName:  "User",
		Email:     "orderitem.user@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	// 2. Create Merchant
	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      s.userID,
		Name:        "OrderItem Merchant",
		Description: "A merchant for testing order items",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	// 3. Create Category
	slugCat := "test-cat-oi"
	category, err := s.repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "Test Cat OI",
		Description:  "Test Category for OI",
		SlugCategory: &slugCat,
	})
	s.Require().NoError(err)
	s.categoryID = int(category.CategoryID)

	// 4. Create Product
	slugProd := "test-prod-oi"
	product, err := s.repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:   s.merchantID,
		CategoryID:   s.categoryID,
		Name:         "Test Prod OI",
		Description:  "Test Product for OI",
		Price:        100,
		CountInStock: 100,
		Brand:        "Test Brand",
		Weight:       1000,
		SlugProduct:  &slugProd,
	})
	s.Require().NoError(err)
	s.productID = int(product.ProductID)

	// 5. Create Order
	order, err := s.repos.Order.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: s.merchantID,
		UserID:     s.userID,
		TotalPrice: 100,
	})
	s.Require().NoError(err)
	s.orderID = int(order.OrderID)
}

func (s *OrderItemRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *OrderItemRepositoryTestSuite) TestOrderItemLifecycle() {
	ctx := context.Background()

	// 1. Create Order Item
	createReq := &requests.CreateOrderItemRecordRequest{
		OrderID:   s.orderID,
		ProductID: s.productID,
		Quantity:  2,
		Price:     100,
	}

	orderItem, err := s.repos.OrderItem.CreateOrderItem(ctx, createReq)
	s.NoError(err)
	s.NotNil(orderItem)
	s.Equal(int32(createReq.Quantity), orderItem.Quantity)
	orderItemID := int(orderItem.OrderItemID)

	// 2. Find Item By Order
	items, err := s.repos.OrderItem.FindOrderItemByOrder(ctx, s.orderID)
	s.NoError(err)
	s.NotEmpty(items)
	s.Equal(int32(orderItemID), items[0].OrderItemID)

	// 3. Update Order Item
	updateReq := &requests.UpdateOrderItemRecordRequest{
		OrderItemID: orderItemID,
		OrderID:     s.orderID,
		ProductID:   s.productID,
		Quantity:    3,
		Price:       100,
	}

	updated, err := s.repos.OrderItem.UpdateOrderItem(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(updateReq.Quantity), updated.Quantity)

	// 4. Trash Order Item
	trashed, err := s.repos.OrderItem.TrashedOrderItem(ctx, orderItemID)
	s.NoError(err)
	s.NotNil(trashed)
	s.True(trashed.DeletedAt.Valid)

	// 5. Restore Order Item
	restored, err := s.repos.OrderItem.RestoreOrderItem(ctx, orderItemID)
	s.NoError(err)
	s.NotNil(restored)
	s.False(restored.DeletedAt.Valid)

	// 6. Delete Permanent
	// Trash again first
	_, err = s.repos.OrderItem.TrashedOrderItem(ctx, orderItemID)
	s.NoError(err)

	success, err := s.repos.OrderItem.DeleteOrderItemPermanent(ctx, orderItemID)
	s.NoError(err)
	s.True(success)

	// 7. Verify it's gone
	trashedItems, err := s.repos.OrderItem.FindOrderItemByOrder(ctx, s.orderID)
	s.NoError(err)
	s.Empty(trashedItems)
}

func TestOrderItemRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemRepositoryTestSuite))
}
