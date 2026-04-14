package service_test

import (
	"context"
	"ecommerce/internal/cache"
	auth_cache "ecommerce/internal/cache/auth"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	"ecommerce/pkg/auth"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/hash"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"ecommerce/tests"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type AuthServiceTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	authService service.AuthService
}

func (s *AuthServiceTestSuite) SetupSuite() {
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
	tokenManager, _ := auth.NewManager("test-secret-key-that-is-long-enough-32")
	hasher := hash.NewHashingPassword()
	obs, _ := observability.NewObservability("test", log)
	cacheStore := cache.NewCacheStore(s.redisClient, log, nil)
	cacheAuth := auth_cache.NewMencache(cacheStore)

	s.authService = service.NewAuthService(service.AuthServiceDeps{
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

	// Seed ROLE_ADMIN as it's required for registration
	_, err = repos.Role.CreateRole(context.Background(), &requests.CreateRoleRequest{
		Name: "ROLE_ADMIN",
	})
	s.Require().NoError(err)
}

func (s *AuthServiceTestSuite) TearDownSuite() {
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *AuthServiceTestSuite) TestRegisterAndLogin() {
	// Register
	regReq := &requests.CreateUserRequest{
		FirstName: "Auth",
		LastName:  "Service",
		Email:     "auth.service@example.com",
		Password:  "password123",
	}
	user, err := s.authService.Register(context.Background(), regReq)
	s.NoError(err)
	s.NotNil(user)

	// Login
	loginReq := &requests.AuthRequest{
		Email:    regReq.Email,
		Password: "password123",
	}
	tokenRes, err := s.authService.Login(context.Background(), loginReq)
	s.NoError(err)
	s.NotNil(tokenRes)
	s.NotEmpty(tokenRes.AccessToken)
	s.NotEmpty(tokenRes.RefreshToken)
}

func TestAuthServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthServiceTestSuite))
}
