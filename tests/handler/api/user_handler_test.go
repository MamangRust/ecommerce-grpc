package api_test

import (
	"bytes"
	"ecommerce/internal/cache"
	api_cache "ecommerce/internal/cache/api/user"
	user_cache "ecommerce/internal/cache/user"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/handler/api"
	"ecommerce/internal/handler/gapi"
	apimapper "ecommerce/internal/mapper"
	"ecommerce/internal/pb"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	db "ecommerce/pkg/database/schema"
	app_errors "ecommerce/pkg/errors"
	"ecommerce/pkg/hash"
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

type UserHandlerTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.UserServiceClient
	conn        *grpc.ClientConn
	router      *echo.Echo
}

func (s *UserHandlerTestSuite) SetupSuite() {
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
	log, _ := logger.NewLogger("test", lp)
	hasher := hash.NewHashingPassword()
	obs, _ := observability.NewObservability("test", log)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	cUser := user_cache.NewUserMencache(cacheStore)

	userService := service.NewUserService(service.UserServiceDeps{
		UserRepo:      repos.User,
		Logger:        log,
		Observability: obs,
		Hashing:       hasher,
		Cache:         cUser,
	})

	// Start gRPC Server
	userGapiHandler := gapi.NewUserHandleGrpc(userService)
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, userGapiHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	// Create gRPC Client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewUserServiceClient(conn)

	// Setup Echo
	s.router = echo.New()
	userMapper := apimapper.NewUserResponseMapper()
	apiErrorHandler := app_errors.NewApiHandler(obs, log)
	apiCacheUser := api_cache.NewUserMencache(cacheStore)

	api.NewHandlerUser(s.router, s.client, log, userMapper, apiErrorHandler, apiCacheUser)
}

func (s *UserHandlerTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *UserHandlerTestSuite) TestUserLifecycle() {
	// 1. Create User
	req := requests.CreateUserRequest{
		FirstName:       "Handler",
		LastName:        "User",
		Email:           "handler.user@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	body, _ := json.Marshal(req)

	request := httptest.NewRequest(http.MethodPost, "/api/user/create", bytes.NewBuffer(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	var createRes map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &createRes)
	userData := createRes["data"].(map[string]interface{})
	userID := int(userData["id"].(float64))

	// 2. Find By ID
	request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/user/%d", userID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Update
	updateReq := requests.UpdateUserRequest{
		FirstName:       "Updated",
		LastName:        "User",
		Email:           "handler.user@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	updateBody, _ := json.Marshal(updateReq)
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/user/update/%d", userID), bytes.NewBuffer(updateBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Delete Permanent
	request = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/user/permanent/%d", userID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)
	s.Equal(http.StatusOK, rec.Code)
}

func TestUserHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserHandlerTestSuite))
}
