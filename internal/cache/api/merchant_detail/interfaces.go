package merchantdetail_cache

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
)

type MerchantDetailQueryCache interface {
	GetCachedMerchantDetailAll(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantDetail, bool)
	SetCachedMerchantDetailAll(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantDetail)

	GetCachedMerchantDetailActive(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantDetailDeleteAt, bool)
	SetCachedMerchantDetailActive(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantDetailDeleteAt)

	GetCachedMerchantDetailTrashed(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantDetailDeleteAt, bool)
	SetCachedMerchantDetailTrashed(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantDetailDeleteAt)

	GetCachedMerchantDetail(ctx context.Context, id int) (*response.ApiResponseMerchantDetailCore, bool)
	SetCachedMerchantDetail(ctx context.Context, data *response.ApiResponseMerchantDetailCore)

	GetCachedMerchantDetailRelation(
		ctx context.Context,
		merchantID int,
	) (*response.ApiResponseMerchantDetailRelation, bool)

	SetCachedMerchantDetailRelation(
		ctx context.Context,
		merchantID int,
		data *response.ApiResponseMerchantDetailRelation,
	)
}

type MerchantDetailCommandCache interface {
	DeleteMerchantDetailCache(ctx context.Context, id int)
}
