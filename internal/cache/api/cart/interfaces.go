package cart_cache

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
)

type CartQueryCache interface {
	GetCachedCarts(
		ctx context.Context,
		request *requests.FindAllCarts,
	) (*response.ApiResponseCartPagination, bool)

	SetCachedCarts(
		ctx context.Context,
		request *requests.FindAllCarts,
		response *response.ApiResponseCartPagination,
	)
}
