package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	cart_cache "ecommerce/internal/cache/cart"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/handler/gapi"
	"ecommerce/internal/pb"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"ecommerce/tests"
	"net"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CartGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.CartServiceClient
	conn        *grpc.ClientConn
	userID      int
	productID   int
}

func (s *CartGapiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)

	queries := db.New(pool)
	repos := repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test-cart-gapi", lp)
	obs, _ := observability.NewObservability("test-cart-gapi", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-cart-gapi")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	cartCache := cart_cache.NewCartMencache(cacheStore)

	cartService := service.NewCartService(service.CartServiceDeps{
		ProductRepository: repos.Product,
		UserRepository:    repos.User,
		CartRepository:    repos.Cart,
		Logger:            log,
		Cache:              cartCache,
		Observability:     obs,
	})

	// Start gRPC Server
	cartHandler := gapi.NewCartHandleGrpc(cartService)
	server := grpc.NewServer()
	pb.RegisterCartServiceServer(server, cartHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	// Create Client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewCartServiceClient(conn)

	ctx := context.Background()
	// Setup User and Product
	user, _ := repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Gapi",
		LastName:  "User",
		Email:     "gapi.user@example.com",
		Password:  "password123",
	})
	s.userID = int(user.UserID)

	merchant, _ := repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: s.userID,
		Name:   "Gapi Merchant",
	})

	slugCat := "gapi-category"
	category, _ := repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "Gapi Category",
		SlugCategory: &slugCat,
	})

	slugProd := "gapi-product"
	product, _ := repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:   int(merchant.MerchantID),
		CategoryID:   int(category.CategoryID),
		Name:         "Gapi Product",
		Description:  "Gapi Product Description",
		Price:        3000,
		CountInStock: 30,
		Brand:        "Gapi Brand",
		Weight:       100,
		Rating:       &[]int{5}[0],
		SlugProduct:  &slugProd,
		ImageProduct: "gapi-product.jpg",
	})
	s.productID = int(product.ProductID)
}

func (s *CartGapiTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *CartGapiTestSuite) TestCartLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pb.CreateCartRequest{
		Quantity:  5,
		ProductId: int32(s.productID),
		UserId:    int32(s.userID),
	}

	res, err := s.client.Create(ctx, createReq)
	s.NoError(err)
	s.Equal(createReq.Quantity, res.Data.Quantity)
	cartID := res.Data.Id

	// 2. Find All
	findAllRes, err := s.client.FindAll(ctx, &pb.FindAllCartRequest{
		UserId:   int32(s.userID),
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(findAllRes.Data)

	// 3. Delete
	delRes, err := s.client.Delete(ctx, &pb.FindByIdCartRequest{Id: cartID})
	s.NoError(err)
	s.Equal("success", delRes.Status)

	// 4. Create and Delete All
	c2, _ := s.client.Create(ctx, createReq)
	c3, _ := s.client.Create(ctx, createReq)

	delAllRes, err := s.client.DeleteAll(ctx, &pb.DeleteCartRequest{
		CartIds: []int32{c2.Data.Id, c3.Data.Id},
	})
	s.NoError(err)
	s.Equal("success", delAllRes.Status)
}

func TestCartGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CartGapiTestSuite))
}
