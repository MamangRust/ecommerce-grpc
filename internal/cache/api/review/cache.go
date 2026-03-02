package review_cache

import "ecommerce/internal/cache"

type reviewMencache struct {
	ReviewQueryCache
	ReviewCommandCache
}

type ReviewMencache interface {
	ReviewQueryCache
	ReviewCommandCache
}

func NewReviewMencache(cacheStore *cache.CacheStore) ReviewMencache {
	return &reviewMencache{
		ReviewQueryCache:   NewReviewQueryCache(cacheStore),
		ReviewCommandCache: NewReviewCommandCache(cacheStore),
	}
}
