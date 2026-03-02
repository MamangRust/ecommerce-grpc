package auth_cache

import (
	"context"
	"ecommerce/internal/cache"
	"ecommerce/internal/domain/response"
	"fmt"
	"time"
)

type loginCache struct {
	store *cache.CacheStore
}

func NewLoginCache(store *cache.CacheStore) *loginCache {
	return &loginCache{store: store}
}

func (s *loginCache) GetCachedLogin(ctx context.Context, email string) (*response.TokenResponse, bool) {
	key := fmt.Sprintf(keylogin, email)

	result, found := cache.GetFromCache[*response.TokenResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *loginCache) SetCachedLogin(ctx context.Context, email string, data *response.TokenResponse, expiration time.Duration) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(keylogin, email)

	cache.SetToCache(ctx, s.store, key, data, expiration)
}
