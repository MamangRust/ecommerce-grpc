package api_test

import (
	"bytes"
	"context"
	"ecommerce/internal/cache"
	api_order_cache "ecommerce/internal/cache/api/order"
	order_cache "ecommerce/internal/cache/order"
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

type OrderApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	echo        *echo.Echo
	client      pb.OrderServiceClient
	repos       *repository.Repositories
	conn        *grpc.ClientConn
	merchantID  int
	userID      int
	categoryID  int
	productID   int
}

func (s *OrderApiTestSuite) SetupSuite() {
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
	log, _ := logger.NewLogger("test-api", lp)
	obs, _ := observability.NewObservability("test-api", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-api")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	
	// Service layer cache
	orderCacheSrv := order_cache.OrderNewMencache(cacheStore)
	// API layer cache
	orderCacheApi := api_order_cache.OrderNewMencache(cacheStore)

	orderService := service.NewOrderService(service.OrderServiceDeps{
		OrderRepository:     s.repos.Order,
		OrderItemRepository: s.repos.OrderItem,
		ProductRepository:   s.repos.Product,
		UserRepository:      s.repos.User,
		MerchantRepository:  s.repos.Merchant,
		ShippingRepository:  s.repos.Shipping,
		Logger:              log,
		Observability:       obs,
		Cache:               orderCacheSrv,
	})

	// Start gRPC Server
	orderGapi := gapi.NewOrderHandleGrpc(orderService)
	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, orderGapi)
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
	s.client = pb.NewOrderServiceClient(conn)

	// Echo Setup
	s.echo = echo.New()
	mapping := response_api.NewOrderResponseMapper()
	apiErrorHandler := errors.NewApiHandler(obs, log)

	api.NewHandlerOrder(s.echo, s.client, log, mapping, apiErrorHandler, orderCacheApi)

	ctx := context.Background()
	// Setup Dependencies
	user, _ := s.repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Api", LastName: "User", Email: "api.user@example.com", Password: "password123",
	})
	s.userID = int(user.UserID)

	merchant, _ := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: s.userID, Name: "Api Merchant",
	})
	s.merchantID = int(merchant.MerchantID)

	slugCat := "api-cat"
	category, _ := s.repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name: "Api Cat", SlugCategory: &slugCat,
	})
	s.categoryID = int(category.CategoryID)

	slugProd := "api-prod"
	product, _ := s.repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID: s.merchantID, CategoryID: s.categoryID, Name: "Api Prod", Price: 100, CountInStock: 100, SlugProduct: &slugProd,
	})
	s.productID = int(product.ProductID)
}

func (s *OrderApiTestSuite) TearDownSuite() {
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

func (s *OrderApiTestSuite) TestOrderApiLifecycle() {
	ctx := context.Background()
	// 1. Create
	createReq := map[string]interface{}{
		"merchant_id": s.merchantID,
		"user_id":     s.userID,
		"total_price": 200,
		"items": []map[string]interface{}{
			{
				"product_id": s.productID,
				"quantity":   2,
				"price":      100,
			},
		},
		"shipping_address": map[string]interface{}{
			"alamat": "API Addr", "provinsi": "P", "kota": "K", "negara": "N", "courier": "C", "shipping_method": "M", "shipping_cost": 1000,
		},
	}
	body, _ := json.Marshal(createReq)

	req := httptest.NewRequest(http.MethodPost, "/api/order/create", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var createRes map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &createRes)
	s.Equal("success", createRes["status"])
	
	data := createRes["data"].(map[string]interface{})
	orderID := int(data["id"].(float64))

	// Fetch ShippingID
	shipping, err := s.repos.Shipping.FindByOrder(ctx, orderID)
	s.NoError(err)
	shippingID := int(shipping.ShippingAddressID)

	// 2. Find By ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/order/%d", orderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Find All
	req = httptest.NewRequest(http.MethodGet, "/api/order", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// Fetch Order Items to get IDs
	orderItems, err := s.repos.OrderItem.FindOrderItemByOrder(ctx, orderID)
	s.NoError(err)
	s.NotEmpty(orderItems)
	orderItemID := int(orderItems[0].OrderItemID)

	// 4. Update
	updateReq := map[string]interface{}{
		"user_id":     s.userID,
		"total_price": 300,
		"items": []map[string]interface{}{
			{
				"order_item_id": orderItemID,
				"product_id":    s.productID,
				"quantity":      3,
				"price":         100,
			},
		},
		"shipping_address": map[string]interface{}{
			"shipping_id":     shippingID,
			"alamat":          "API Addr Updated",
			"courier":         "C",
			"shipping_method": "M",
			"shipping_cost":   1000,
		},
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order/update/%d", orderID), bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order/trashed/%d", orderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order/restore/%d", orderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. Delete Permanent
	// Trash again
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order/trashed/%d", orderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/order/permanent/%d", orderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestOrderApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderApiTestSuite))
}
