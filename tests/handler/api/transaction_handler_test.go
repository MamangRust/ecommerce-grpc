package api_test

import (
	"bytes"
	"context"
	"ecommerce/internal/cache"
	api_cache "ecommerce/internal/cache/api/transaction"
	service_cache "ecommerce/internal/cache/transaction"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/handler/api"
	"ecommerce/internal/handler/gapi"
	mapper "ecommerce/internal/mapper"
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
	"google.golang.org/grpc/test/bufconn"
)

type TransactionApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	repos       *repository.Repositories
	echo        *echo.Echo
	gRpcConn    *grpc.ClientConn
	listener    *bufconn.Listener
	merchantID  int
	userID      int
	orderID     int
}

func (s *TransactionApiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	// DB Setup
	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	// Redis Setup
	opt, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opt)

	// Repositories
	queries := db.New(pool)
	s.repos = repository.NewRepositories(queries)

	// Logging & Observability
	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	l, err := logger.NewLogger("test-transaction-api", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-transaction-api", l)
	s.Require().NoError(err)

	// Cache Setup
	cacheMetrics, err := observability.NewCacheMetrics("test-transaction-api")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.redisClient, l, cacheMetrics)
	transServiceCache := service_cache.NewTransactionMencache(cacheStore)
	transApiCache := api_cache.NewTransactionMencache(cacheStore)

	// Service
	transService := service.NewTransactionService(service.TransactionServiceDeps{
		MerchantRepository:    s.repos.Merchant,
		TransactionRepository: s.repos.Transaction,
		OrderRepository:       s.repos.Order,
		OrderItemRepository:   s.repos.OrderItem,
		ShippingRepository:    s.repos.Shipping,
		Logger:                l,
		Cache:                 transServiceCache,
		Observability:         obs,
	})

	// gRPC Server Setup
	s.listener = bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	hGrpc := gapi.NewTransactionHandleGrpc(transService)
	pb.RegisterTransactionServiceServer(server, hGrpc)

	go func() {
		if err := server.Serve(s.listener); err != nil {
			panic(err)
		}
	}()

	// gRPC Client Setup for API Handler
	conn, err := grpc.DialContext(context.Background(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return s.listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err)
	s.gRpcConn = conn
	transClient := pb.NewTransactionServiceClient(conn)

	// API Handler Setup
	s.echo = echo.New()
	transMapper := mapper.NewTransactionResponseMapper()
	apiHandler := errors.NewApiHandler(obs, l)
	
	api.NewHandlerTransaction(s.echo, transClient, l, transMapper, apiHandler, transApiCache)

	ctx := context.Background()

	// Setup data dependencies
	user, err := s.repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "TransApi",
		LastName:  "User",
		Email:     "trans.api@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      s.userID,
		Name:        "TransApi Merchant",
		Description: "Merchant for API testing",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	slugCat := "api-test-cat"
	cat, err := s.repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:          "Api Test Category",
		Description:   "Description",
		SlugCategory:  &slugCat,
		ImageCategory: "cat.jpg",
	})
	s.Require().NoError(err)

	slugProd := "api-test-prod"
	rating := 5
	prod, err := s.repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:   s.merchantID,
		CategoryID:   int(cat.CategoryID),
		Name:         "Api Test Product",
		Description:  "Product Description",
		Price:        1000,
		CountInStock: 100,
		Brand:        "Brand",
		Weight:       1,
		Rating:       &rating,
		SlugProduct:  &slugProd,
		ImageProduct: "prod.jpg",
	})
	s.Require().NoError(err)

	order, err := s.repos.Order.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: s.merchantID,
		UserID:     s.userID,
		TotalPrice: 1100,
	})
	s.Require().NoError(err)
	s.orderID = int(order.OrderID)

	_, err = s.repos.OrderItem.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{
		OrderID:   s.orderID,
		ProductID: int(prod.ProductID),
		Quantity:  1,
		Price:     1000,
	})
	s.Require().NoError(err)

	_, err = s.repos.Shipping.CreateShippingAddress(ctx, &requests.CreateShippingAddressRequest{
		OrderID:        &s.orderID,
		Alamat:         "Api Alamat Lengkap",
		Provinsi:       "Provinsi",
		Kota:           "Kota",
		Courier:        "JNE",
		ShippingMethod: "REG",
		ShippingCost:   100,
		Negara:         "Indonesia",
	})
	s.Require().NoError(err)
}

func (s *TransactionApiTestSuite) TearDownSuite() {
	if s.gRpcConn != nil {
		s.gRpcConn.Close()
	}
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *TransactionApiTestSuite) TestTransactionApiLifecycle() {
	// 1. Create
	createReq := &requests.CreateTransactionRequest{
		OrderID:       s.orderID,
		MerchantID:    s.merchantID,
		PaymentMethod: "credit_card",
		Amount:        2000,
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/transaction/create", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	var resCreate map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resCreate)
	s.NoError(err)
	s.Equal("success", resCreate["status"])
	
	data := resCreate["data"].(map[string]interface{})
	transID := int(data["id"].(float64))

	// 2. Find All
	req = httptest.NewRequest(http.MethodGet, "/api/transaction?page=1&page_size=10", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Find By ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/transaction/%d", transID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Update (should fail if status is success)
	updateReq := &requests.UpdateTransactionRequest{
		OrderID:       s.orderID,
		MerchantID:    s.merchantID,
		PaymentMethod: "debit_card",
		Amount:        2000,
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transaction/update/%d", transID), bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	// Service layer returns an error, which is mapped to 500 or 400
	s.NotEqual(http.StatusOK, rec.Code)

	// 5. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transaction/trashed/%d", transID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transaction/restore/%d", transID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. Delete Permanent
	// Trash again
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transaction/trashed/%d", transID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/transaction/permanent/%d", transID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *TransactionApiTestSuite) TestTransactionReports() {
	// Setup data for reporting
	// (Already have one order setup, let's create a transaction for it)
	createReq := &requests.CreateTransactionRequest{
		OrderID:       s.orderID,
		MerchantID:    s.merchantID,
		PaymentMethod: "credit_card",
		Amount:        2000,
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/transaction/create", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 1. Monthly Success
	req = httptest.NewRequest(http.MethodGet, "/api/transaction/monthly-success?year=2024&month=4", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	// We might not get data if current date is different, but the endpoint should be reachable
	s.Equal(http.StatusOK, rec.Code)

	// 2. Yearly Success
	req = httptest.NewRequest(http.MethodGet, "/api/transaction/yearly-success?year=2024", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestTransactionApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionApiTestSuite))
}
