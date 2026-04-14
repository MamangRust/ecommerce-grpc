package repository_test

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	"ecommerce/tests"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	dbPool *pgxpool.Pool
	repo   repository.UserRepository
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	queries := db.New(pool)
	s.repo = repository.NewUserRepository(queries)
}

func (s *UserRepositoryTestSuite) TearDownSuite() {
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *UserRepositoryTestSuite) TestCreateUser() {
	req := &requests.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
	}

	user, err := s.repo.CreateUser(context.Background(), req)
	s.NoError(err)
	s.NotNil(user)
	s.Equal(req.FirstName, user.Firstname)
	s.Equal(req.Email, user.Email)
}

func (s *UserRepositoryTestSuite) TestFindById() {
	// Create user first
	req := &requests.CreateUserRequest{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane.smith@example.com",
		Password:  "password123",
	}
	user, _ := s.repo.CreateUser(context.Background(), req)

	found, err := s.repo.FindById(context.Background(), int(user.UserID))
	s.NoError(err)
	s.NotNil(found)
	s.Equal(user.Email, found.Email)
}

func TestUserRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserRepositoryTestSuite))
}
