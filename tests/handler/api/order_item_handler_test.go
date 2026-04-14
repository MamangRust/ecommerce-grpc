package api_test

import (
	"context"
	"ecommerce/internal/cache"
	api_orderitem_cache "ecommerce/internal/cache/api/order_item"
	order_cache "ecommerce/internal/cache/order"
	orderitem_cache "ecommerce/internal/cache/order_item"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/handler/api"
	"ecommerce/internal/handler/gapi"
	response_api "ecommerce/internal/mapper"
	"ecommerce/internal/pb"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"ecommerce/tests"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderItemApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	echo        *echo.Echo
	client      pb.OrderItemServiceClient
	orderSvc    service.OrderService
	conn        *grpc.ClientConn
	merchantID  int
	userID      int
	categoryID  int
	productID   int
	orderID     int
}

func (s *OrderItemApiTestSuite) SetupSuite() {
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
	log, _ := logger.NewLogger("test-api-oi", lp)
	obs, _ := observability.NewObservability("test-api-oi", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-api-oi")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	
	// Service layer cache
	orderItemCacheSrv := orderitem_cache.NewOrderItemMencache(cacheStore)
	orderCacheSrv := order_cache.OrderNewMencache(cacheStore)
	// API layer cache
	orderItemCacheApi := api_orderitem_cache.NewOrderItemMencache(cacheStore)

	orderItemService := service.NewOrderItemService(service.OrderItemServiceDeps{
		OrderItemRepository: repos.OrderItem,
		Logger:              log,
		Observability:       obs,
		Cache:               orderItemCacheSrv,
	})
	
	// Need OrderService to create an order
	s.orderSvc = service.NewOrderService(service.OrderServiceDeps{
		OrderRepository:     repos.Order,
		OrderItemRepository: repos.OrderItem,
		ProductRepository:   repos.Product,
		UserRepository:      repos.User,
		MerchantRepository:  repos.Merchant,
		ShippingRepository:  repos.Shipping,
		Logger:              log,
		Observability:       obs,
		Cache:               orderCacheSrv,
	})

	// Start gRPC Server
	orderItemGapi := gapi.NewOrderItemHandleGrpc(orderItemService)
	server := grpc.NewServer()
	pb.RegisterOrderItemServiceServer(server, orderItemGapi)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	// gRPC Client for the API Handler
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewOrderItemServiceClient(conn)

	// Echo Setup
	s.echo = echo.New()
	mapping := response_api.NewOrderItemResponseMapper()
	apiErrorHandler := errors.NewApiHandler(obs, log)

	api.NewHandlerOrderItem(s.echo, s.client, log, mapping, apiErrorHandler, orderItemCacheApi)

	ctx := context.Background()
	// Setup Dependencies
	user, _ := repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "ApiOI", LastName: "User", Email: "apioi.user@example.com", Password: "password123",
	})
	s.userID = int(user.UserID)

	merchant, _ := repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: s.userID, Name: "ApiOI Merchant",
	})
	s.merchantID = int(merchant.MerchantID)

	slugCat := "apioi-cat"
	category, _ := repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name: "ApiOI Cat", SlugCategory: &slugCat,
	})
	s.categoryID = int(category.CategoryID)

	slugProd := "apioi-prod"
	product, _ := repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID: s.merchantID, CategoryID: s.categoryID, Name: "ApiOI Prod", Price: 100, CountInStock: 100, SlugProduct: &slugProd,
	})
	s.productID = int(product.ProductID)

	// Create Order
	order, _ := s.orderSvc.CreateOrder(ctx, &requests.CreateOrderRequest{
		MerchantID: s.merchantID,
		UserID:     s.userID,
		TotalPrice: 100,
		Items: []requests.CreateOrderItemRequest{
			{
				ProductID: s.productID,
				Quantity:  1,
				Price:     100,
			},
		},
		ShippingAddress: requests.CreateShippingAddressRequest{
			Alamat: "ApiOI Addr",
		},
	})
	s.orderID = int(order.OrderID)
}

func (s *OrderItemApiTestSuite) TearDownSuite() {
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

func (s *OrderItemApiTestSuite) TestOrderItemApiLifecycle() {
	// 1. Find By Order
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/order-item/%d", s.orderID), nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	var res map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &res)
	s.Equal("success", res["status"])

	// 2. Find All
	req = httptest.NewRequest(http.MethodGet, "/api/order-item", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Find By Active
	req = httptest.NewRequest(http.MethodGet, "/api/order-item/active", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Find By Trashed
	req = httptest.NewRequest(http.MethodGet, "/api/order-item/trashed", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestOrderItemApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemApiTestSuite))
}
