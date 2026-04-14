package api_test

import (
	"bytes"
	"ecommerce/internal/cache"
	api_cart_cache "ecommerce/internal/cache/api/cart"
	cart_cache "ecommerce/internal/cache/cart"
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

type CartApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	echo        *echo.Echo
	client      pb.CartServiceClient
	conn        *grpc.ClientConn
	userID      int
	productID   int
}

func (s *CartApiTestSuite) SetupSuite() {
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
	log, _ := logger.NewLogger("test-cart-api", lp)
	obs, _ := observability.NewObservability("test-cart-api", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-cart-api")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	
	cartCacheSrv := cart_cache.NewCartMencache(cacheStore)
	cartCacheApi := api_cart_cache.NewCartMencache(cacheStore)

	cartService := service.NewCartService(service.CartServiceDeps{
		ProductRepository: repos.Product,
		UserRepository:    repos.User,
		CartRepository:    repos.Cart,
		Logger:            log,
		Cache:              cartCacheSrv,
		Observability:     obs,
	})

	// Start gRPC Server
	cartGapi := gapi.NewCartHandleGrpc(cartService)
	server := grpc.NewServer()
	pb.RegisterCartServiceServer(server, cartGapi)
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
	s.client = pb.NewCartServiceClient(conn)

	// Echo Setup
	s.echo = echo.New()
	mapping := response_api.NewCartResponseMapper()
	apiHandler := errors.NewApiHandler(obs, log)

	api.NewHandlerCart(s.echo, s.client, log, mapping, apiHandler, cartCacheApi)

	// Setup User and Product
	ctx := s.ts.Ctx
	user, _ := repos.User.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Api",
		LastName:  "User",
		Email:     "api.user@example.com",
		Password:  "password123",
	})
	s.userID = int(user.UserID)

	merchant, _ := repos.Merchant.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: s.userID,
		Name:   "Api Merchant",
	})

	slugCat := "api-category"
	category, _ := repos.Category.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:         "Api Category",
		SlugCategory: &slugCat,
	})

	slugProd := "api-product"
	product, _ := repos.Product.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID:   int(merchant.MerchantID),
		CategoryID:   int(category.CategoryID),
		Name:         "Api Product",
		Description:  "Api Product Description",
		Price:        4000,
		CountInStock: 40,
		Brand:        "Api Brand",
		Weight:       100,
		Rating:       &[]int{5}[0],
		SlugProduct:  &slugProd,
		ImageProduct: "api-product.jpg",
	})
	s.productID = int(product.ProductID)
}

func (s *CartApiTestSuite) TearDownSuite() {
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

func (s *CartApiTestSuite) TestCartApiLifecycle() {
	// 1. Create
	createReq := map[string]interface{}{
		"quantity":   2,
		"product_id": s.productID,
		"user_id":    s.userID,
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/cart/create", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusCreated, rec.Code)
	var createRes map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &createRes)
	
	data := createRes["data"].(map[string]interface{})
	cartID := int(data["id"].(float64))

	// 2. Find All
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/cart?user_id=%d", s.userID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Delete
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/cart/%d", cartID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Delete All
	c1Req := map[string]interface{}{"quantity": 1, "product_id": s.productID, "user_id": s.userID}
	b1, _ := json.Marshal(c1Req)
	r1 := httptest.NewRequest(http.MethodPost, "/api/cart/create", bytes.NewBuffer(b1))
	r1.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	re1 := httptest.NewRecorder()
	s.echo.ServeHTTP(re1, r1)
	var cr1 map[string]interface{}
	json.Unmarshal(re1.Body.Bytes(), &cr1)
	id1 := int(cr1["data"].(map[string]interface{})["id"].(float64))

	c2Req := map[string]interface{}{"quantity": 1, "product_id": s.productID, "user_id": s.userID}
	b2, _ := json.Marshal(c2Req)
	r2 := httptest.NewRequest(http.MethodPost, "/api/cart/create", bytes.NewBuffer(b2))
	r2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	re2 := httptest.NewRecorder()
	s.echo.ServeHTTP(re2, r2)
	var cr2 map[string]interface{}
	json.Unmarshal(re2.Body.Bytes(), &cr2)
	id2 := int(cr2["data"].(map[string]interface{})["id"].(float64))

	deleteAllReq := map[string]interface{}{
		"cart_ids": []int{id1, id2},
	}
	body, _ = json.Marshal(deleteAllReq)
	req = httptest.NewRequest(http.MethodPost, "/api/cart/delete-all", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestCartApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CartApiTestSuite))
}
