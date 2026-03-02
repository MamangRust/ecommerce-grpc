package auth_cache

import (
	"context"
	"ecommerce/internal/domain/response"
	db "ecommerce/pkg/database/schema"
	"time"
)

type IdentityCache interface {
	SetRefreshToken(ctx context.Context, token string, expiration time.Duration)
	GetRefreshToken(ctx context.Context, token string) (string, bool)
	DeleteRefreshToken(ctx context.Context, token string)
	SetCachedUserInfo(ctx context.Context, user *db.GetUserByIDRow, expiration time.Duration)
	GetCachedUserInfo(ctx context.Context, userId string) (*db.GetUserByIDRow, bool)
	DeleteCachedUserInfo(ctx context.Context, userId string)
}

type LoginCache interface {
	SetCachedLogin(ctx context.Context, email string, data *response.TokenResponse, expiration time.Duration)
	GetCachedLogin(ctx context.Context, email string) (*response.TokenResponse, bool)
}
