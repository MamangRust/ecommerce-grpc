package api_test

import (
	"context"
	"ecommerce/internal/cache"
	api_cache "ecommerce/internal/cache/api/shipping_address"
	service_cache "ecommerce/internal/cache/shipping_address"
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

type ShippingApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	repos       *repository.Repositories
	echo        *echo.Echo
	gRpcConn    *grpc.ClientConn
	listener    *bufconn.Listener
	userID      int
	merchantID  int
	orderID     int
}

func (s *ShippingApiTestSuite) SetupSuite() {
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
	l, err := logger.NewLogger("test-shipping-api", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-shipping-api", l)
	s.Require().NoError(err)

	// Cache Setup
	cacheMetrics, err := observability.NewCacheMetrics("test-shipping-api")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.redisClient, l, cacheMetrics)
	shipServiceCache := service_cache.NewShippingAddressMencache(cacheStore)
	shipApiCache := api_cache.NewShippingAddressMencache(cacheStore)

	// Service
	shipService := service.NewShippingAddressService(service.ShippingAddressServiceDeps{
		ShippingRepository: s.repos.Shipping,
		Logger:             l,
		Observability:      obs,
		Cache:              shipServiceCache,
	})

	// gRPC Server Setup
	s.listener = bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	hGrpc := gapi.NewShippingAddressHandleGrpc(shipService)
	pb.RegisterShippingServiceServer(server, hGrpc)

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
	shipClient := pb.NewShippingServiceClient(conn)

	// API Handler Setup
	s.echo = echo.New()
	shipMapper := mapper.NewShippingAddressResponseMapper()
	apiHandler := errors.NewApiHandler(obs, l)
	
	api.NewHandlerShippingAddress(s.echo, shipClient, l, shipMapper, apiHandler, shipApiCache)

	// Create Prerequisites
	ctx := context.Background()
	user, err := s.repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Api",
		LastName:  "User",
		Email:     "api.shipping@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      s.userID,
		Name:        "Api Merchant",
		Description: "A test merchant for api tests",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	order, err := s.repos.Order.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		UserID:     s.userID,
		MerchantID: s.merchantID,
		TotalPrice: 4000,
	})
	s.Require().NoError(err)
	s.orderID = int(order.OrderID)
}

func (s *ShippingApiTestSuite) TearDownSuite() {
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

func (s *ShippingApiTestSuite) TestShippingAddressApiLifecycle() {
	// 1. Create via Repository (API usually just reads)
	ctx := context.Background()
	createReq := &requests.CreateShippingAddressRequest{
		OrderID:        &s.orderID,
		Alamat:         "Jl. Api No. 1",
		Provinsi:       "East Java",
		Kota:           "Surabaya",
		Negara:         "Indonesia",
		Courier:        "JNE",
		ShippingMethod: "REG",
		ShippingCost:   10000,
	}

	address, err := s.repos.Shipping.CreateShippingAddress(ctx, createReq)
	s.NoError(err)
	s.NotNil(address)
	shippingID := int(address.ShippingAddressID)

	// 2. Find By ID
	// Note the base path from NewHandlerShippingAddress: /api/shipping-address
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/shipping-address/%d", shippingID), nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	var res map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	s.NoError(err)
	s.Equal("success", res["status"])

	// 3. Find By Order
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/shipping-address/order/%d", s.orderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Find All
	req = httptest.NewRequest(http.MethodGet, "/api/shipping-address?page=1&pageSize=10", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/shipping-address/trashed/%d", shippingID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/shipping-address/restore/%d", shippingID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. Delete Permanent
	// Trash again
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/shipping-address/trashed/%d", shippingID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/shipping-address/permanent/%d", shippingID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 8. Verify
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/shipping-address/%d", shippingID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusInternalServerError, rec.Code) // Service returns error if not found
}

func TestShippingApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ShippingApiTestSuite))
}
