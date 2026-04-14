package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	shippingaddress_cache "ecommerce/internal/cache/shipping_address"
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
	"google.golang.org/grpc/test/bufconn"
)

type ShippingGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	repos       *repository.Repositories
	client      pb.ShippingServiceClient
	server      *grpc.Server
	listener    *bufconn.Listener
	userID      int
	merchantID  int
	orderID     int
}

func (s *ShippingGapiTestSuite) SetupSuite() {
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
	l, err := logger.NewLogger("test-shipping-gapi", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-shipping-gapi", l)
	s.Require().NoError(err)

	// Cache Setup
	cacheMetrics, err := observability.NewCacheMetrics("test-shipping-gapi")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.redisClient, l, cacheMetrics)
	shipCache := shippingaddress_cache.NewShippingAddressMencache(cacheStore)

	// Service
	shipService := service.NewShippingAddressService(service.ShippingAddressServiceDeps{
		ShippingRepository: s.repos.Shipping,
		Logger:             l,
		Observability:      obs,
		Cache:              shipCache,
	})

	// gRPC Server Setup
	s.listener = bufconn.Listen(1024 * 1024)
	s.server = grpc.NewServer()
	h := gapi.NewShippingAddressHandleGrpc(shipService)
	pb.RegisterShippingServiceServer(s.server, h)

	go func() {
		if err := s.server.Serve(s.listener); err != nil {
			panic(err)
		}
	}()

	// gRPC Client Setup
	conn, err := grpc.DialContext(context.Background(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return s.listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err)
	s.client = pb.NewShippingServiceClient(conn)

	// Create Prerequisites
	ctx := context.Background()
	user, err := s.repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Gapi",
		LastName:  "User",
		Email:     "gapi.shipping@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      s.userID,
		Name:        "Gapi Merchant",
		Description: "A test merchant for gapi tests",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	order, err := s.repos.Order.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		UserID:     s.userID,
		MerchantID: s.merchantID,
		TotalPrice: 3000,
	})
	s.Require().NoError(err)
	s.orderID = int(order.OrderID)
}

func (s *ShippingGapiTestSuite) TearDownSuite() {
	s.server.Stop()
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

func (s *ShippingGapiTestSuite) TestShippingAddressGapiLifecycle() {
	ctx := context.Background()

	// 1. Create Shipping Address via Repository
	createReq := &requests.CreateShippingAddressRequest{
		OrderID:        &s.orderID,
		Alamat:         "Jl. Gapi No. 1",
		Provinsi:       "Central Java",
		Kota:           "Semarang",
		Negara:         "Indonesia",
		Courier:        "J&T",
		ShippingMethod: "Economy",
		ShippingCost:   8000,
	}

	address, err := s.repos.Shipping.CreateShippingAddress(ctx, createReq)
	s.NoError(err)
	s.NotNil(address)
	shippingID := int32(address.ShippingAddressID)

	// 2. Find By ID
	res, err := s.client.FindById(ctx, &pb.FindByIdShippingRequest{Id: shippingID})
	s.NoError(err)
	s.NotNil(res)
	s.Equal("success", res.Status)
	s.Equal(createReq.Alamat, res.Data.Alamat)

	// 3. Find By Order
	resOrder, err := s.client.FindByOrder(ctx, &pb.FindByIdShippingRequest{Id: int32(s.orderID)})
	s.NoError(err)
	s.NotNil(resOrder)
	s.Equal(createReq.Alamat, resOrder.Data.Alamat)

	// 4. Find All
	resAll, err := s.client.FindAll(ctx, &pb.FindAllShippingRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotNil(resAll)
	s.GreaterOrEqual(resAll.Pagination.TotalRecords, int32(1))

	// 5. Trash
	resTrash, err := s.client.TrashedShipping(ctx, &pb.FindByIdShippingRequest{Id: shippingID})
	s.NoError(err)
	s.NotNil(resTrash)
	s.NotEmpty(resTrash.Data.DeletedAt.Value)

	// 6. Restore
	resRestore, err := s.client.RestoreShipping(ctx, &pb.FindByIdShippingRequest{Id: shippingID})
	s.NoError(err)
	s.NotNil(resRestore)
	s.Empty(resRestore.Data.DeletedAt.Value)

	// 7. Delete Permanent
	// Trash again before permanent delete
	_, err = s.client.TrashedShipping(ctx, &pb.FindByIdShippingRequest{Id: shippingID})
	s.NoError(err)

	resDelete, err := s.client.DeleteShippingPermanent(ctx, &pb.FindByIdShippingRequest{Id: shippingID})
	s.NoError(err)
	s.NotNil(resDelete)
	s.Equal("success", resDelete.Status)

	// 8. Verify
	_, err = s.client.FindById(ctx, &pb.FindByIdShippingRequest{Id: shippingID})
	s.Error(err)
}

func TestShippingGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ShippingGapiTestSuite))
}
