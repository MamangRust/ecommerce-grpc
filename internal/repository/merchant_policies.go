package repository

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
	merchantpolicy_errors "ecommerce/pkg/errors/merchant_policy_errors"
)

type merchantPolicyRepository struct {
	db *db.Queries
}

func NewMerchantPoliciesRepository(db *db.Queries) *merchantPolicyRepository {
	return &merchantPolicyRepository{
		db: db,
	}
}

func (r *merchantPolicyRepository) FindAllMerchantPolicy(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantPoliciesParams{
		Column1: &req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantPolicies(ctx, reqDb)

	if err != nil {
		return nil, merchantpolicy_errors.ErrFindAllMerchantPolicy
	}

	return res, nil
}

func (r *merchantPolicyRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantPoliciesActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantPoliciesActive(ctx, reqDb)

	if err != nil {
		return nil, merchantpolicy_errors.ErrFindByActiveMerchantPolicy
	}

	return res, nil
}

func (r *merchantPolicyRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantPoliciesTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantPoliciesTrashed(ctx, reqDb)

	if err != nil {
		return nil, merchantpolicy_errors.ErrFindByTrashedMerchantPolicy
	}

	return res, nil
}

func (r *merchantPolicyRepository) FindById(ctx context.Context, user_id int) (*db.GetMerchantPolicyRow, error) {
	res, err := r.db.GetMerchantPolicy(ctx, int32(user_id))

	if err != nil {
		return nil, merchantpolicy_errors.ErrFindByIdMerchantPolicy
	}

	return res, nil
}

func (r *merchantPolicyRepository) CreateMerchantPolicy(ctx context.Context, request *requests.CreateMerchantPolicyRequest) (*db.CreateMerchantPolicyRow, error) {
	req := db.CreateMerchantPolicyParams{
		MerchantID:  int32(request.MerchantID),
		PolicyType:  request.PolicyType,
		Title:       request.Title,
		Description: request.Description,
	}

	policy, err := r.db.CreateMerchantPolicy(ctx, req)
	if err != nil {
		return nil, merchantpolicy_errors.ErrCreateMerchantPolicy
	}

	return policy, nil
}

func (r *merchantPolicyRepository) UpdateMerchantPolicy(ctx context.Context, request *requests.UpdateMerchantPolicyRequest) (*db.UpdateMerchantPolicyRow, error) {
	req := db.UpdateMerchantPolicyParams{
		MerchantPolicyID: int32(*request.MerchantPolicyID),
		PolicyType:       request.PolicyType,
		Title:            request.Title,
		Description:      request.Description,
	}

	res, err := r.db.UpdateMerchantPolicy(ctx, req)
	if err != nil {
		return nil, merchantpolicy_errors.ErrUpdateMerchantPolicy
	}

	return res, nil
}

func (r *merchantPolicyRepository) TrashedMerchantPolicy(ctx context.Context, merchant_id int) (*db.MerchantPolicy, error) {
	res, err := r.db.TrashMerchantPolicy(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantpolicy_errors.ErrTrashedMerchantPolicy
	}

	return res, nil
}

func (r *merchantPolicyRepository) RestoreMerchantPolicy(ctx context.Context, merchant_id int) (*db.MerchantPolicy, error) {
	res, err := r.db.RestoreMerchantPolicy(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantpolicy_errors.ErrRestoreMerchantPolicy
	}

	return res, nil
}

func (r *merchantPolicyRepository) DeleteMerchantPolicyPermanent(ctx context.Context, Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantPermanently(ctx, int32(Merchant_id))

	if err != nil {
		return false, merchantpolicy_errors.ErrDeleteMerchantPolicyPermanent
	}

	return true, nil
}

func (r *merchantPolicyRepository) RestoreAllMerchantPolicy(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchants(ctx)

	if err != nil {
		return false, merchantpolicy_errors.ErrRestoreAllMerchantPolicy
	}
	return true, nil
}

func (r *merchantPolicyRepository) DeleteAllMerchantPolicyPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(ctx)

	if err != nil {
		return false, merchantpolicy_errors.ErrDeleteAllMerchantPolicyPermanent
	}
	return true, nil
}
