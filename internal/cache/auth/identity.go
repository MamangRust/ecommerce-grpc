package auth_cache

import (
	"context"
	"ecommerce/internal/cache"
	db "ecommerce/pkg/database/schema"
	"fmt"
	"time"
)

type identityCache struct {
	store *cache.CacheStore
}

func NewidentityCache(store *cache.CacheStore) *identityCache {
	return &identityCache{store: store}
}

func (c *identityCache) SetRefreshToken(ctx context.Context, token string, expiration time.Duration) {
	key := fmt.Sprintf(keyIdentityRefreshToken, token)

	cache.SetToCache(ctx, c.store, key, &token, expiration)
}

func (c *identityCache) GetRefreshToken(ctx context.Context, token string) (string, bool) {
	key := fmt.Sprintf(keyIdentityRefreshToken, token)

	result, found := cache.GetFromCache[string](ctx, c.store, key)

	if !found || result == nil {
		return "", false
	}

	return *result, true
}

func (c *identityCache) DeleteRefreshToken(ctx context.Context, token string) {
	key := fmt.Sprintf(keyIdentityRefreshToken, token)
	cache.DeleteFromCache(ctx, c.store, key)
}

func (c *identityCache) SetCachedUserInfo(ctx context.Context, user *db.GetUserByIDRow, expiration time.Duration) {
	if user == nil {
		return
	}

	key := fmt.Sprintf(keyIdentityUserInfo, user.UserID)

	cache.SetToCache(ctx, c.store, key, user, expiration)
}

func (c *identityCache) GetCachedUserInfo(ctx context.Context, userId string) (*db.GetUserByIDRow, bool) {
	key := fmt.Sprintf(keyIdentityUserInfo, userId)

	return cache.GetFromCache[db.GetUserByIDRow](ctx, c.store, key)
}

func (c *identityCache) DeleteCachedUserInfo(ctx context.Context, userId string) {
	key := fmt.Sprintf(keyIdentityUserInfo, userId)

	cache.DeleteFromCache(ctx, c.store, key)
}
