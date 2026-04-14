package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	auth_cache "ecommerce/internal/cache/auth"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/handler/gapi"
	"ecommerce/internal/pb"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	"ecommerce/pkg/auth"
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

type AuthGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.AuthServiceClient
	conn        *grpc.ClientConn
}

func (s *AuthGapiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	s.redisClient = redis.NewClient(&redis.Options{
		Addr: s.ts.RedisURL,
	})

	queries := db.New(pool)
	repos := repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	tokenManager, _ := auth.NewManager("test-secret-key-that-is-long-enough-32")
	hasher := hash.NewHashingPassword()
	obs, _ := observability.NewObservability("test", log)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	cacheAuth := auth_cache.NewMencache(cacheStore)

	authService := service.NewAuthService(service.AuthServiceDeps{
		AuthRepo:         repos.User,
		RefreshTokenRepo: repos.RefreshToken,
		RoleRepo:         repos.Role,
		UserRoleRepo:     repos.UserRole,
		Hash:             hasher,
		Token:            tokenManager,
		Logger:           log,
		Tracer:           obs,
		CacheIdentity:    cacheAuth.IdentityCache,
		CacheLogin:       cacheAuth.LoginCache,
	})

	// Seed ROLE_ADMIN
	_, _ = repos.Role.CreateRole(context.Background(), &requests.CreateRoleRequest{Name: "ROLE_ADMIN"})

	// Start gRPC Server
	authHandler := gapi.NewAuthHandleGrpc(authService)
	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, authHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		if err := server.Serve(lis); err != nil {
			return
		}
	}()

	// Create Client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewAuthServiceClient(conn)
}

func (s *AuthGapiTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *AuthGapiTestSuite) TestRegisterUser() {
	req := &pb.RegisterRequest{
		Firstname:       "Gapi",
		Lastname:        "Test",
		Email:           "gapi@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	res, err := s.client.RegisterUser(context.Background(), req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.Equal(req.Email, res.Data.Email)
}

func TestAuthGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthGapiTestSuite))
}
