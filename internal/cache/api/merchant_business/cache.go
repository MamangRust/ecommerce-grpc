package merchantbusiness_cache

import "ecommerce/internal/cache"

type MerchantBusinessMencache interface {
	MerchantBusinessQueryCache
	MerchantBusinessCommandCache
}

type merchantBussinessMencache struct {
	MerchantBusinessQueryCache
	MerchantBusinessCommandCache
}

func NewMerchantBusinessMencache(cacheStore *cache.CacheStore) MerchantBusinessMencache {
	return &merchantBussinessMencache{
		MerchantBusinessQueryCache:   NewMerchantBusinessQueryCache(cacheStore),
		MerchantBusinessCommandCache: NewMerchantBusinessCommandCache(cacheStore),
	}
}
