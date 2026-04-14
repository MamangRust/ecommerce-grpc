package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	order_cache "ecommerce/internal/cache/order"
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

type OrderGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.OrderServiceClient
	repos       *repository.Repositories
	conn        *grpc.ClientConn
	merchantID  int
	userID      int
	categoryID  int
	productID   int
}

func (s *OrderGapiTestSuite) SetupSuite() {
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
	s.repos = repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test-gapi", lp)
	obs, _ := observability.NewObservability("test-gapi", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-gapi")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	orderCache := order_cache.OrderNewMencache(cacheStore)

	orderService := service.NewOrderService(service.OrderServiceDeps{
		OrderRepository:     s.repos.Order,
		OrderItemRepository: s.repos.OrderItem,
		ProductRepository:   s.repos.Product,
		UserRepository:      s.repos.User,
		MerchantRepository:  s.repos.Merchant,
		ShippingRepository:  s.repos.Shipping,
		Logger:              log,
		Observability:       obs,
		Cache:               orderCache,
	})

	// Start gRPC Server
	orderHandler := gapi.NewOrderHandleGrpc(orderService)
	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, orderHandler)
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
	s.client = pb.NewOrderServiceClient(conn)

	ctx := context.Background()

	// Setup Dependencies
	user, err := s.repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Gapi", LastName: "User", Email: "gapi.user@example.com", Password: "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: s.userID, Name: "Gapi Merchant",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	slugCat := "gapi-cat"
	category, err := s.repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name: "Gapi Cat", SlugCategory: &slugCat,
	})
	s.Require().NoError(err)
	s.categoryID = int(category.CategoryID)

	slugProd := "gapi-prod"
	product, err := s.repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID: s.merchantID, CategoryID: s.categoryID, Name: "Gapi Prod", Price: 100, CountInStock: 100, SlugProduct: &slugProd,
	})
	s.Require().NoError(err)
	s.productID = int(product.ProductID)
}

func (s *OrderGapiTestSuite) TearDownSuite() {
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

func (s *OrderGapiTestSuite) TestOrderGapiLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pb.CreateOrderRequest{
		MerchantId: int32(s.merchantID),
		UserId:     int32(s.userID),
		TotalPrice: 200,
		Items: []*pb.CreateOrderItemRequest{
			{
				ProductId: int32(s.productID),
				Quantity:  2,
				Price:     100,
			},
		},
		Shipping: &pb.CreateShippingAddressRequest{
			Alamat: "Gapi Addr", Courier: "Gapi Courier", ShippingMethod: "Method", ShippingCost: 1000, Negara: "ID", Provinsi: "Prov", Kota: "Kota",
		},
	}

	res, err := s.client.Create(ctx, createReq)
	s.NoError(err)
	s.NotNil(res)
	orderID := res.Data.Id

	// Fetch ShippingID
	shipping, err := s.repos.Shipping.FindByOrder(ctx, int(orderID))
	s.NoError(err)
	shippingID := int32(shipping.ShippingAddressID)

	// 2. Find By ID
	found, err := s.client.FindById(ctx, &pb.FindByIdOrderRequest{Id: orderID})
	s.NoError(err)
	s.Equal(orderID, found.Data.Id)

	// 3. Find All
	orders, err := s.client.FindAll(ctx, &pb.FindAllOrderRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(orders.Data)

	// Fetch Order Items to get IDs
	orderItems, err := s.repos.OrderItem.FindOrderItemByOrder(ctx, int(orderID))
	s.Require().NoError(err)
	s.NotEmpty(orderItems)
	orderItemID := orderItems[0].OrderItemID

	// 4. Update
	updateReq := &pb.UpdateOrderRequest{
		OrderId:    orderID,
		UserId:     int32(s.userID),
		TotalPrice: 300,
		Items: []*pb.UpdateOrderItemRequest{
			{
				OrderItemId: int32(orderItemID),
				ProductId:   int32(s.productID),
				Quantity:    3,
				Price:       100,
			},
		},
		Shipping: &pb.UpdateShippingAddressRequest{
			ShippingId:     shippingID,
			Alamat:         "Gapi Addr Updated",
			Courier:        "Gapi Courier",
			ShippingMethod: "Method",
			ShippingCost:   1000,
		},
	}
	updated, err := s.client.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)

	// 5. Trash
	_, err = s.client.TrashedOrder(ctx, &pb.FindByIdOrderRequest{Id: orderID})
	s.NoError(err)

	// 6. Restore
	_, err = s.client.RestoreOrder(ctx, &pb.FindByIdOrderRequest{Id: orderID})
	s.NoError(err)

	// 7. Delete Permanent
	_, err = s.client.TrashedOrder(ctx, &pb.FindByIdOrderRequest{Id: orderID})
	s.NoError(err)

	delRes, err := s.client.DeleteOrderPermanent(ctx, &pb.FindByIdOrderRequest{Id: orderID})
	s.NoError(err)
	s.Equal("success", delRes.Status)
}

func TestOrderGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderGapiTestSuite))
}
