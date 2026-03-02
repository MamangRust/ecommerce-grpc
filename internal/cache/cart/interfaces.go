package cart_cache

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
)

type CartQueryCache interface {
	GetCachedCartsCache(ctx context.Context, request *requests.FindAllCarts) ([]*db.GetCartsRow, *int, bool)
	SetCartsCache(ctx context.Context, request *requests.FindAllCarts, response []*db.GetCartsRow, total *int)
}
