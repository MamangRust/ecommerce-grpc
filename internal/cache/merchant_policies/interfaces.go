package merchantpolicies_cache

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
)

type MerchantPolicyQueryCache interface {
	GetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesRow, *int, bool)
	SetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantPoliciesRow, total *int)

	GetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesActiveRow, *int, bool)
	SetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantPoliciesActiveRow, total *int)

	GetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesTrashedRow, *int, bool)
	SetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantPoliciesTrashedRow, total *int)

	GetCachedMerchantPolicy(ctx context.Context, id int) (*db.GetMerchantPolicyRow, bool)
	SetCachedMerchantPolicy(ctx context.Context, data *db.GetMerchantPolicyRow)
}

type MerchantPolicyCommandCache interface {
	DeleteMerchantPolicyCache(ctx context.Context, merchantID int)
}
