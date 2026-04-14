package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	user_cache "ecommerce/internal/cache/user"
	"ecommerce/internal/handler/gapi"
	"ecommerce/internal/pb"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/hash"
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

type UserGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.UserServiceClient
	conn        *grpc.ClientConn
}

func (s *UserGapiTestSuite) SetupSuite() {
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
	userHandler := gapi.NewUserHandleGrpc(userService)
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, userHandler)
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
	s.client = pb.NewUserServiceClient(conn)
}

func (s *UserGapiTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *UserGapiTestSuite) TestUserLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pb.CreateUserRequest{
		Firstname:       "Gapi",
		Lastname:        "User",
		Email:           "gapi.user@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	res, err := s.client.Create(ctx, createReq)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().NotNil(res.Data)
	s.Equal(createReq.Email, res.Data.Email)

	// 2. Find By ID
	findReq := &pb.FindByIdUserRequest{
		Id: res.Data.Id,
	}
	found, err := s.client.FindById(ctx, findReq)
	s.Require().NoError(err)
	s.Require().NotNil(found)
	s.Require().NotNil(found.Data)
	s.Equal(res.Data.Id, found.Data.Id)
}

func TestUserGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserGapiTestSuite))
}
