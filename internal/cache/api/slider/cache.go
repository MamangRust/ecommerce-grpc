package slider_cache

import "ecommerce/internal/cache"

type sliderMencache struct {
	SliderQueryCache
	SliderCommandCache
}

type SliderMencache interface {
	SliderQueryCache
	SliderCommandCache
}

func NewSliderMencache(cacheStore *cache.CacheStore) SliderMencache {
	return sliderMencache{
		SliderQueryCache:   NewSliderQueryCache(cacheStore),
		SliderCommandCache: NewSliderCommandCache(cacheStore),
	}
}
