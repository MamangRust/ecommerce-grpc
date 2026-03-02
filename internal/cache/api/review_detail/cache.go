package reviewdetail_cache

import "ecommerce/internal/cache"

type ReviewDetailMencache interface {
	ReviewDetailQueryCache
	ReviewDetailCommandCache
}

type reviewDetaiMencache struct {
	ReviewDetailQueryCache
	ReviewDetailCommandCache
}

func NewReviewDetailMencache(cacheStore *cache.CacheStore) ReviewDetailMencache {
	return &reviewDetaiMencache{
		ReviewDetailQueryCache:   NewReviewDetailQueryCache(cacheStore),
		ReviewDetailCommandCache: NewReviewDetailCommandCache(cacheStore),
	}
}
