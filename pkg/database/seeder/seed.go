package seeder

import (
	"context"

	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/hash"
	"ecommerce/pkg/logger"
	"fmt"
	"time"
)

type Deps struct {
	Db     *db.Queries
	Ctx    context.Context
	Logger logger.LoggerInterface
	Hash   hash.HashPassword
}

type Seeder struct {
	User     *userSeeder
	Role     *roleSeeder
	UserRole *userRoleSeeder
}

func NewSeeder(deps Deps) *Seeder {
	return &Seeder{
		User:     NewUserSeeder(deps.Db, deps.Hash, deps.Ctx, deps.Logger),
		Role:     NewRoleSeeder(deps.Db, deps.Ctx, deps.Logger),
		UserRole: NewUserRoleSeeder(deps.Db, deps.Ctx, deps.Logger),
	}
}

func (s *Seeder) Run() error {
	if err := s.seedWithDelay("users", s.User.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("roles", s.Role.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("user_roles", s.UserRole.Seed); err != nil {
		return err
	}

	return nil
}

func (s *Seeder) seedWithDelay(entityName string, seedFunc func() error) error {
	if err := seedFunc(); err != nil {
		return fmt.Errorf("failed to seed %s: %w", entityName, err)
	}

	time.Sleep(30 * time.Second)
	return nil
}
