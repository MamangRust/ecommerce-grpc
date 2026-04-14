package service_test

import (
	"context"
	"ecommerce/internal/cache"
	user_cache "ecommerce/internal/cache/user"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
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

type UserServiceTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	userService service.UserService
}

func (s *UserServiceTestSuite) SetupSuite() {
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
	cacheStore := cache.NewCacheStore(s.redisClient, log, nil)
	cacheUser := user_cache.NewUserMencache(cacheStore)

	s.userService = service.NewUserService(service.UserServiceDeps{
		UserRepo:      repos.User,
		Logger:        log,
		Observability: obs,
		Hashing:       hasher,
		Cache:         cacheUser,
	})
}

func (s *UserServiceTestSuite) TearDownSuite() {
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *UserServiceTestSuite) TestUserLifecycle() {
	ctx := context.Background()

	// 1. Create User
	req := &requests.CreateUserRequest{
		FirstName: "User",
		LastName:  "Service",
		Email:     "user.service@example.com",
		Password:  "password123",
	}
	user, err := s.userService.CreateUser(ctx, req)
	s.NoError(err)
	s.NotNil(user)
	s.Equal(req.Email, user.Email)

	// 2. Find By ID
	found, err := s.userService.FindByID(ctx, int(user.UserID))
	s.NoError(err)
	s.NotNil(found)
	s.Equal(user.Email, found.Email)

	// 3. Update User
	updateID := int(user.UserID)
	updateReq := &requests.UpdateUserRequest{
		UserID:    &updateID,
		FirstName: "Updated",
	}
	updated, err := s.userService.UpdateUser(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("Updated", updated.Firstname)

	// 4. Trash and Restore
	_, err = s.userService.TrashedUser(ctx, int(user.UserID))
	s.NoError(err)

	_, err = s.userService.RestoreUser(ctx, int(user.UserID))
	s.NoError(err)

	// 5. Delete Permanent
	success, err := s.userService.DeleteUserPermanent(ctx, int(user.UserID))
	s.NoError(err)
	s.True(success)
}

func TestUserServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserServiceTestSuite))
}
