package banner_cache

import (
	"ecommerce/internal/cache"
)

type bannerMencache struct {
	BannerQueryCache
	BannerCommandCache
}

type BannerMencache interface {
	BannerQueryCache
	BannerCommandCache
}

func NewBannerMencache(cacheStore *cache.CacheStore) BannerMencache {
	return &bannerMencache{
		BannerQueryCache:   NewBannerQueryCache(cacheStore),
		BannerCommandCache: NewBannerCommandCache(cacheStore),
	}
}
