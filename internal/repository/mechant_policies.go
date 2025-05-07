package repository

import (
	"context"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	merchantpolicy_errors "ecommerce/pkg/errors/merchant_policy_errors"
)

type merchantPolicyRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantPolicyMapping
}

func NewMerchantPolicyRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantPolicyMapping) *merchantPolicyRepository {
	return &merchantPolicyRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantPolicyRepository) FindAllMerchantPolicy(req *requests.FindAllMerchant) ([]*record.MerchantPoliciesRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantPoliciesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantPolicies(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchantpolicy_errors.ErrFindAllMerchantPolicy
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantPolicysRecordPagination(res), &totalCount, nil
}

func (r *merchantPolicyRepository) FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantPoliciesRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantPoliciesActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantPoliciesActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchantpolicy_errors.ErrFindByActiveMerchantPolicy
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantPolicysRecordActivePagination(res), &totalCount, nil
}

func (r *merchantPolicyRepository) FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantPoliciesRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantPoliciesTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantPoliciesTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchantpolicy_errors.ErrFindByTrashedMerchantPolicy
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantPolicysRecordTrashedPagination(res), &totalCount, nil
}

func (r *merchantPolicyRepository) FindById(user_id int) (*record.MerchantPoliciesRecord, error) {
	res, err := r.db.GetMerchantPolicy(r.ctx, int32(user_id))

	if err != nil {
		return nil, merchantpolicy_errors.ErrFindByIdMerchantPolicy
	}

	return r.mapping.ToMerchantPolicyRecord(res), nil
}

func (r *merchantPolicyRepository) CreateMerchantPolicy(request *requests.CreateMerchantPolicyRequest) (*record.MerchantPoliciesRecord, error) {
	req := db.CreateMerchantPolicyParams{
		MerchantID:  int32(request.MerchantID),
		PolicyType:  request.PolicyType,
		Title:       request.Title,
		Description: request.Description,
	}

	policy, err := r.db.CreateMerchantPolicy(r.ctx, req)
	if err != nil {
		return nil, merchantpolicy_errors.ErrCreateMerchantPolicy
	}

	return r.mapping.ToMerchantPolicyRecord(policy), nil
}

func (r *merchantPolicyRepository) UpdateMerchantPolicy(request *requests.UpdateMerchantPolicyRequest) (*record.MerchantPoliciesRecord, error) {
	req := db.UpdateMerchantPolicyParams{
		MerchantPolicyID: int32(*request.MerchantPolicyID),
		PolicyType:       request.PolicyType,
		Title:            request.Title,
		Description:      request.Description,
	}

	res, err := r.db.UpdateMerchantPolicy(r.ctx, req)
	if err != nil {
		return nil, merchantpolicy_errors.ErrUpdateMerchantPolicy
	}

	return r.mapping.ToMerchantPolicyRecord(res), nil
}

func (r *merchantPolicyRepository) TrashedMerchantPolicy(merchant_id int) (*record.MerchantPoliciesRecord, error) {
	res, err := r.db.TrashMerchantPolicy(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantpolicy_errors.ErrTrashedMerchantPolicy
	}

	return r.mapping.ToMerchantPolicyRecord(res), nil
}

func (r *merchantPolicyRepository) RestoreMerchantPolicy(merchant_id int) (*record.MerchantPoliciesRecord, error) {
	res, err := r.db.RestoreMerchantPolicy(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantpolicy_errors.ErrRestoreMerchantPolicy
	}

	return r.mapping.ToMerchantPolicyRecord(res), nil
}

func (r *merchantPolicyRepository) DeleteMerchantPolicyPermanent(Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantPermanently(r.ctx, int32(Merchant_id))

	if err != nil {
		return false, merchantpolicy_errors.ErrDeleteMerchantPolicyPermanent
	}

	return true, nil
}

func (r *merchantPolicyRepository) RestoreAllMerchantPolicy() (bool, error) {
	err := r.db.RestoreAllMerchants(r.ctx)

	if err != nil {
		return false, merchantpolicy_errors.ErrRestoreAllMerchantPolicy
	}
	return true, nil
}

func (r *merchantPolicyRepository) DeleteAllMerchantPolicyPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(r.ctx)

	if err != nil {
		return false, merchantpolicy_errors.ErrDeleteAllMerchantPolicyPermanent
	}
	return true, nil
}
