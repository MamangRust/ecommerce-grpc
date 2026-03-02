package shippingaddress_cache

import "ecommerce/internal/cache"

type ShippingAddressMencache interface {
	ShippingAddressQueryCache
	ShippingAddressCommandCache
}

type shippingAddressMencache struct {
	ShippingAddressQueryCache
	ShippingAddressCommandCache
}

func NewShippingAddressMencache(cacheStore *cache.CacheStore) ShippingAddressMencache {
	return &shippingAddressMencache{
		ShippingAddressQueryCache:   NewShippingAddressQueryCache(cacheStore),
		ShippingAddressCommandCache: NewShippingAddressCommandCache(cacheStore),
	}
}
