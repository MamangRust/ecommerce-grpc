package merchantawards_cache

import "ecommerce/internal/cache"

type MerchantAwardMencache interface {
	MerchantAwardQueryCache
	MerchantAwardCommandCache
}

type merchantAwardMencache struct {
	MerchantAwardQueryCache
	MerchantAwardCommandCache
}

func NewMerchantAward(cacheStore *cache.CacheStore) MerchantAwardMencache {
	return &merchantAwardMencache{
		MerchantAwardQueryCache:   NewMerchantAwardQueryCache(cacheStore),
		MerchantAwardCommandCache: NewMerchantAwardCommandCache(cacheStore),
	}
}
