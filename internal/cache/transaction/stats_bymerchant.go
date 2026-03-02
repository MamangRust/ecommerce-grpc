package transaction_cache

import (
	"context"
	"ecommerce/internal/cache"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
	"fmt"
)

const (
	transactonMonthAmountSuccessByMerchantKey = "transaction:month:amount:success:merchant:%d:month:%d:year:%d"
	transactonMonthAmountFailedByMerchantKey  = "transaction:month:amount:failed:merchant:%d:month:%d:year:%d"
	transactonYearAmountSuccessByMerchantKey  = "transaction:year:amount:success:merchant:%d:year:%d"
	transactonYearAmountFailedByMerchantKey   = "transaction:year:amount:failed:merchant:%d:year:%d"
	transactonMonthMethodSuccessByMerchantKey = "transaction:month:method:success:merchant:%d:month:%d:year:%d"
	transactonMonthMethodFailedByMerchantKey  = "transaction:month:method:failed:merchant:%d:month:%d:year:%d"
	transactonYearMethodSuccessByMerchantKey  = "transaction:year:method:success:merchant:%d:year:%d"
	transactonYearMethodFailedByMerchantKey   = "transaction:year:method:failed:merchant:%d:year:%d"
)

type transactionStatsByMerchantCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsByMerchantCache(store *cache.CacheStore) *transactionStatsByMerchantCache {
	return &transactionStatsByMerchantCache{store: store}
}

func (t *transactionStatsByMerchantCache) GetCachedMonthAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionSuccessByMerchantRow, bool) {
	key := fmt.Sprintf(transactonMonthAmountSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyAmountTransactionSuccessByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*db.GetMonthlyAmountTransactionSuccessByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonMonthAmountSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionSuccessByMerchantRow, bool) {
	key := fmt.Sprintf(transactonYearAmountSuccessByMerchantKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*db.GetYearlyAmountTransactionSuccessByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*db.GetYearlyAmountTransactionSuccessByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonYearAmountSuccessByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionFailedByMerchantRow, bool) {
	key := fmt.Sprintf(transactonMonthAmountFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyAmountTransactionFailedByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*db.GetMonthlyAmountTransactionFailedByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonMonthAmountFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionFailedByMerchantRow, bool) {
	key := fmt.Sprintf(transactonYearAmountFailedByMerchantKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*db.GetYearlyAmountTransactionFailedByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*db.GetYearlyAmountTransactionFailedByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonYearAmountFailedByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthMethodSuccessByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow, bool) {
	key := fmt.Sprintf(transactonMonthMethodSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthMethodSuccessByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*db.GetMonthlyTransactionMethodsByMerchantSuccessRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonMonthMethodSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearMethodSuccessByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantSuccessRow, bool) {
	key := fmt.Sprintf(transactonYearMethodSuccessByMerchantKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTransactionMethodsByMerchantSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearMethodSuccessByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*db.GetYearlyTransactionMethodsByMerchantSuccessRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonYearMethodSuccessByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthMethodFailedByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantFailedRow, bool) {
	key := fmt.Sprintf(transactonMonthMethodFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyTransactionMethodsByMerchantFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthMethodFailedByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*db.GetMonthlyTransactionMethodsByMerchantFailedRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonMonthMethodFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearMethodFailedByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantFailedRow, bool) {
	key := fmt.Sprintf(transactonYearMethodFailedByMerchantKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTransactionMethodsByMerchantFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearMethodFailedByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*db.GetYearlyTransactionMethodsByMerchantFailedRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonYearMethodFailedByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}
