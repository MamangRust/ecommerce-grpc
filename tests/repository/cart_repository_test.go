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

type CartRepositoryTestSuite struct {
	suite.Suite
	ts        *tests.TestSuite
	dbPool    *pgxpool.Pool
	repos     *repository.Repositories
	userID    int
	productID int
}

func (s *CartRepositoryTestSuite) SetupSuite() {
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
		FirstName: "Cart",
		LastName:  "User",
		Email:     "cart.user@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	// 2. Create Merchant (needed for product)
	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      s.userID,
		Name:        "Cart Merchant",
		Description: "Merchant for cart tests",
	})
	s.Require().NoError(err)

	// 3. Create Category (needed for product)
	slugCat := "cart-category"
	category, err := s.repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "Cart Category",
		Description:  "Category for cart tests",
		SlugCategory: &slugCat,
	})
	s.Require().NoError(err)

	// 4. Create Product
	slugProd := "cart-product"
	product, err := s.repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:   int(merchant.MerchantID),
		CategoryID:   int(category.CategoryID),
		Name:         "Cart Product",
		Description:  "Product for cart tests",
		Price:        1000,
		CountInStock: 10,
		Brand:        "Cart Brand",
		Weight:       500,
		SlugProduct:  &slugProd,
		ImageProduct: "cart-product.jpg",
	})
	s.Require().NoError(err)
	s.productID = int(product.ProductID)
}

func (s *CartRepositoryTestSuite) TearDownSuite() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *CartRepositoryTestSuite) TestCartLifecycle() {
	ctx := context.Background()

	// 1. Create Cart
	createReq := &requests.CartCreateRecord{
		ProductID:    s.productID,
		UserID:       s.userID,
		Name:         "Cart Product",
		Price:        1000,
		ImageProduct: "cart-product.jpg",
		Quantity:     2,
		Weight:       500,
	}

	cart, err := s.repos.Cart.CreateCart(ctx, createReq)
	s.NoError(err)
	s.NotNil(cart)
	s.Equal(int32(createReq.Quantity), cart.Quantity)

	cartID := int(cart.CartID)

	// 2. Find Carts
	carts, err := s.repos.Cart.FindCarts(ctx, &requests.FindAllCarts{
		UserID:   s.userID,
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.NotEmpty(carts)

	// 3. Delete Permanent
	success, err := s.repos.Cart.DeletePermanent(ctx, cartID)
	s.NoError(err)
	s.True(success)

	// 4. Create and Delete All Permanently
	cart2, _ := s.repos.Cart.CreateCart(ctx, createReq)
	cart3, _ := s.repos.Cart.CreateCart(ctx, createReq)

	success, err = s.repos.Cart.DeleteAllPermanently(ctx, &requests.DeleteCartRequest{
		CartIds: []int{int(cart2.CartID), int(cart3.CartID)},
	})
	s.NoError(err)
	s.True(success)
}

func TestCartRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CartRepositoryTestSuite))
}
