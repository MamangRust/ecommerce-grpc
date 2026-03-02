package service

import (
	"context"
	merchantpolicies_cache "ecommerce/internal/cache/merchant_policies"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	merchantpolicy_errors "ecommerce/pkg/errors/merchant_policy_errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantPoliciesService struct {
	merchantPoliciesRepository repository.MerchantPoliciesRepository
	logger                     logger.LoggerInterface
	observability              observability.TraceLoggerObservability
	cache                      merchantpolicies_cache.MerchantPoliciesMencache
}

type MerchantPoliciesServiceDeps struct {
	MerchantPoliciesRepository repository.MerchantPoliciesRepository
	Logger                     logger.LoggerInterface
	Observability              observability.TraceLoggerObservability
	Cache                      merchantpolicies_cache.MerchantPoliciesMencache
}

func NewMerchantPoliciesService(deps MerchantPoliciesServiceDeps) *merchantPoliciesService {
	return &merchantPoliciesService{
		merchantPoliciesRepository: deps.MerchantPoliciesRepository,
		logger:                     deps.Logger,
		observability:              deps.Observability,
		cache:                      deps.Cache,
	}
}

func (s *merchantPoliciesService) FindAllMerchantPolicy(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesRow, *int, error) {
	const method = "FindAllMerchantPolicy"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchantPolicyAll(ctx, req); found {
		logSuccess("Successfully retrieved all merchant policy records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantPoliciesRepository.FindAllMerchantPolicy(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantPoliciesRow](
			s.logger,
			merchantpolicy_errors.ErrFailedFindAllMerchantPolicies,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedMerchantPolicyAll(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched all merchant policies",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantPoliciesService) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesActiveRow, *int, error) {
	const method = "FindByActiveMerchantPolicy"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchantPolicyActive(ctx, req); found {
		logSuccess("Successfully retrieved active merchant policy records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantPoliciesRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantPoliciesActiveRow](
			s.logger,
			merchantpolicy_errors.ErrFailedFindActiveMerchantPolicies,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedMerchantPolicyActive(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched active merchant policies",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantPoliciesService) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesTrashedRow, *int, error) {
	const method = "FindByTrashedMerchantPolicy"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchantPolicyTrashed(ctx, req); found {
		logSuccess("Successfully retrieved trashed merchant policy records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantPoliciesRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantPoliciesTrashedRow](
			s.logger,
			merchantpolicy_errors.ErrFailedFindTrashedMerchantPolicies,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedMerchantPolicyTrashed(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched trashed merchant policies",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantPoliciesService) FindById(ctx context.Context, merchantID int) (*db.GetMerchantPolicyRow, error) {
	const method = "FindMerchantPolicyById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMerchantPolicy(ctx, merchantID); found {
		logSuccess("Successfully retrieved merchant policy by ID from cache",
			zap.Int("merchantID", merchantID))
		return data, nil
	}

	merchant, err := s.merchantPoliciesRepository.FindById(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetMerchantPolicyRow](
			s.logger,
			merchantpolicy_errors.ErrFailedFindMerchantPolicyById,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.SetCachedMerchantPolicy(ctx, merchant)

	logSuccess("Successfully fetched merchant policy by ID",
		zap.Int("merchantID", merchantID))

	return merchant, nil
}

func (s *merchantPoliciesService) CreateMerchantPolicy(ctx context.Context, req *requests.CreateMerchantPolicyRequest) (*db.CreateMerchantPolicyRow, error) {
	const method = "CreateMerchantPolicy"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantPoliciesRepository.CreateMerchantPolicy(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateMerchantPolicyRow](
			s.logger,
			merchantpolicy_errors.ErrFailedCreateMerchantPolicy,
			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteMerchantPolicyCache(ctx, int(merchant.MerchantPolicyID))

	logSuccess("Successfully created merchant policy",
		zap.Int("merchantPolicyID", int(merchant.MerchantPolicyID)))

	return merchant, nil
}

func (s *merchantPoliciesService) UpdateMerchantPolicy(ctx context.Context, req *requests.UpdateMerchantPolicyRequest) (*db.UpdateMerchantPolicyRow, error) {
	const method = "UpdateMerchantPolicy"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantPolicyID", *req.MerchantPolicyID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantPoliciesRepository.UpdateMerchantPolicy(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantPolicyRow](
			s.logger,
			merchantpolicy_errors.ErrFailedUpdateMerchantPolicy,
			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteMerchantPolicyCache(ctx, int(merchant.MerchantPolicyID))

	logSuccess("Successfully updated merchant policy",
		zap.Int("merchantPolicyID", int(merchant.MerchantPolicyID)))

	return merchant, nil
}

func (s *merchantPoliciesService) TrashedMerchantPolicy(ctx context.Context, merchantID int) (*db.MerchantPolicy, error) {
	const method = "TrashedMerchantPolicy"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantPoliciesRepository.TrashedMerchantPolicy(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantPolicy](
			s.logger,
			merchantpolicy_errors.ErrFailedTrashedMerchantPolicy,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	// Invalidate cache
	s.cache.DeleteMerchantPolicyCache(ctx, merchantID)

	logSuccess("Successfully trashed merchant policy",
		zap.Int("merchantID", merchantID))

	return merchant, nil
}

func (s *merchantPoliciesService) RestoreMerchantPolicy(ctx context.Context, merchantID int) (*db.MerchantPolicy, error) {
	const method = "RestoreMerchantPolicy"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantPoliciesRepository.RestoreMerchantPolicy(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantPolicy](
			s.logger,
			merchantpolicy_errors.ErrFailedRestoreMerchantPolicy,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.DeleteMerchantPolicyCache(ctx, merchantID)

	logSuccess("Successfully restored merchant policy",
		zap.Int("merchantID", merchantID))

	return merchant, nil
}

func (s *merchantPoliciesService) DeleteMerchantPolicyPermanent(ctx context.Context, merchantID int) (bool, error) {
	const method = "DeleteMerchantPolicyPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantPoliciesRepository.DeleteMerchantPolicyPermanent(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantpolicy_errors.ErrFailedDeleteMerchantPolicyPermanent,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.DeleteMerchantPolicyCache(ctx, merchantID)

	logSuccess("Successfully permanently deleted merchant policy",
		zap.Int("merchantID", merchantID))

	return success, nil
}

func (s *merchantPoliciesService) RestoreAllMerchantPolicy(ctx context.Context) (bool, error) {
	const method = "RestoreAllMerchantPolicy"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantPoliciesRepository.RestoreAllMerchantPolicy(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantpolicy_errors.ErrFailedRestoreAllMerchantPolicies,
			method,
			span,
		)
	}

	// Note: We can't selectively invalidate cache for all merchant policies,
	// so we would need to implement a cache flush method if needed
	// For now, we'll rely on cache expiration

	logSuccess("Successfully restored all trashed merchant policies")

	return success, nil
}

func (s *merchantPoliciesService) DeleteAllMerchantPolicyPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllMerchantPolicyPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantPoliciesRepository.DeleteAllMerchantPolicyPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantpolicy_errors.ErrFailedDeleteAllMerchantPoliciesPermanent,
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all trashed merchant policies")

	return success, nil
}
