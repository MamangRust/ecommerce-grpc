package auth_cache

import "ecommerce/internal/cache"

type AuthMencache struct {
	IdentityCache IdentityCache
	LoginCache    LoginCache
}

func NewMencache(cacheStore *cache.CacheStore) *AuthMencache {
	return &AuthMencache{
		IdentityCache: NewidentityCache(cacheStore),
		LoginCache:    NewLoginCache(cacheStore),
	}
}
