package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	transaction_cache "ecommerce/internal/cache/transaction"
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

type TransactionGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	repos       *repository.Repositories
	client      pb.TransactionServiceClient
	conn        *grpc.ClientConn
	listener    *bufconn.Listener
	merchantID  int
	userID      int
	orderID     int
}

func (s *TransactionGapiTestSuite) SetupSuite() {
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
	l, err := logger.NewLogger("test-transaction-gapi", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-transaction-gapi", l)
	s.Require().NoError(err)

	// Cache Setup
	cacheMetrics, err := observability.NewCacheMetrics("test-transaction-gapi")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.redisClient, l, cacheMetrics)
	transCache := transaction_cache.NewTransactionMencache(cacheStore)

	// Service
	transService := service.NewTransactionService(service.TransactionServiceDeps{
		MerchantRepository:    s.repos.Merchant,
		TransactionRepository: s.repos.Transaction,
		OrderRepository:       s.repos.Order,
		OrderItemRepository:   s.repos.OrderItem,
		ShippingRepository:    s.repos.Shipping,
		Logger:                l,
		Cache:                 transCache,
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

	// gRPC Client Setup
	conn, err := grpc.DialContext(context.Background(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return s.listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewTransactionServiceClient(conn)

	ctx := context.Background()

	// Setup data dependencies
	user, err := s.repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "TransGapi",
		LastName:  "User",
		Email:     "trans.gapi@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)

	merchant, err := s.repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:      s.userID,
		Name:        "TransGapi Merchant",
		Description: "Merchant for gAPI testing",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	slugCat := "gapi-test-cat"
	cat, err := s.repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:          "Gapi Test Category",
		Description:   "Description",
		SlugCategory:  &slugCat,
		ImageCategory: "cat.jpg",
	})
	s.Require().NoError(err)

	slugProd := "gapi-test-prod"
	rating := 5
	prod, err := s.repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:   s.merchantID,
		CategoryID:   int(cat.CategoryID),
		Name:         "Gapi Test Product",
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
		Alamat:         "Gapi Alamat Lengkap",
		Provinsi:       "Provinsi",
		Kota:           "Kota",
		Courier:        "JNE",
		ShippingMethod: "REG",
		ShippingCost:   100,
		Negara:         "Indonesia",
	})
	s.Require().NoError(err)
}

func (s *TransactionGapiTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
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

func (s *TransactionGapiTestSuite) TestTransactionGapiLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pb.CreateTransactionRequest{
		OrderId:       int32(s.orderID),
		MerchantId:    int32(s.merchantID),
		PaymentMethod: "credit_card",
		Amount:        2000,
	}

	resCreate, err := s.client.Create(ctx, createReq)
	s.NoError(err)
	s.NotNil(resCreate)
	s.Equal("success", resCreate.Data.PaymentStatus)
	transID := resCreate.Data.Id

	// 2. Find All
	resAll, err := s.client.FindAll(ctx, &pb.FindAllTransactionRequest{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(resAll)
	s.GreaterOrEqual(resAll.Pagination.TotalRecords, int32(1))

	// 3. Find By ID
	resFind, err := s.client.FindById(ctx, &pb.FindByIdTransactionRequest{Id: transID})
	s.NoError(err)
	s.NotNil(resFind)
	s.Equal(transID, resFind.Data.Id)

	// 4. Update (should fail because status is success)
	updateReq := &pb.UpdateTransactionRequest{
		TransactionId: transID,
		OrderId:       int32(s.orderID),
		MerchantId:    int32(s.merchantID),
		PaymentMethod: "debit_card",
		Amount:        2000,
	}
	_, err = s.client.Update(ctx, updateReq)
	s.Error(err)

	// 5. Trash
	resTrash, err := s.client.TrashedTransaction(ctx, &pb.FindByIdTransactionRequest{Id: transID})
	s.NoError(err)
	s.NotNil(resTrash)
	s.NotNil(resTrash.Data.DeletedAt)

	// 6. Restore
	resRestore, err := s.client.RestoreTransaction(ctx, &pb.FindByIdTransactionRequest{Id: transID})
	s.NoError(err)
	s.NotNil(resRestore)
	s.Equal("", resRestore.Data.DeletedAt.GetValue()) // wrapperspb.StringValue empty if nil in some cases or we check GetValue

	// 7. Delete Permanent
	// Trash again
	_, err = s.client.TrashedTransaction(ctx, &pb.FindByIdTransactionRequest{Id: transID})
	s.NoError(err)

	resDelete, err := s.client.DeleteTransactionPermanent(ctx, &pb.FindByIdTransactionRequest{Id: transID})
	s.NoError(err)
	s.NotNil(resDelete)
	s.Equal("success", resDelete.Status)
}

func TestTransactionGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionGapiTestSuite))
}
