package service

import (
	"context"
	merchant_cache "ecommerce/internal/cache/merchant"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	merchant_errors "ecommerce/pkg/errors/merchant"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantService struct {
	merchantRepository repository.MerchantRepository
	logger             logger.LoggerInterface
	observability      observability.TraceLoggerObservability
	cache              merchant_cache.MerchantMencache
}

type MerchantServiceDeps struct {
	MerchantRepository repository.MerchantRepository
	Logger             logger.LoggerInterface
	Observability      observability.TraceLoggerObservability
	Cache              merchant_cache.MerchantMencache
}

func NewMerchantService(deps MerchantServiceDeps) *merchantService {
	return &merchantService{
		merchantRepository: deps.MerchantRepository,
		logger:             deps.Logger,
		observability:      deps.Observability,
		cache:              deps.Cache,
	}
}

func (s *merchantService) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsRow, *int, error) {
	const method = "FindAllMerchants"

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

	if data, total, found := s.cache.GetCachedMerchants(ctx, req); found {
		logSuccess("Successfully retrieved all merchant records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantRepository.FindAllMerchants(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantsRow](
			s.logger,
			merchant_errors.ErrFailedFindAllMerchants,
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

	s.cache.SetCachedMerchants(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched all merchants",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantService) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsActiveRow, *int, error) {
	const method = "FindByActiveMerchants"

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

	if data, total, found := s.cache.GetCachedMerchantActive(ctx, req); found {
		logSuccess("Successfully retrieved active merchant records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantsActiveRow](
			s.logger,
			merchant_errors.ErrFailedFindActiveMerchants,
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

	s.cache.SetCachedMerchantActive(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched active merchants",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantService) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsTrashedRow, *int, error) {
	const method = "FindByTrashedMerchants"

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

	if data, total, found := s.cache.GetCachedMerchantTrashed(ctx, req); found {
		logSuccess("Successfully retrieved trashed merchant records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantsTrashedRow](
			s.logger,
			merchant_errors.ErrFailedFindTrashedMerchants,
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

	s.cache.SetCachedMerchantTrashed(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched trashed merchants",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantService) FindById(ctx context.Context, merchantID int) (*db.GetMerchantByIDRow, error) {
	const method = "FindMerchantById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMerchant(ctx, merchantID); found {
		logSuccess("Successfully retrieved merchant by ID from cache",
			zap.Int("merchantID", merchantID))
		return data, nil
	}

	merchant, err := s.merchantRepository.FindById(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetMerchantByIDRow](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.SetCachedMerchant(ctx, merchant)

	logSuccess("Successfully fetched merchant by ID",
		zap.Int("merchantID", merchantID))

	return merchant, nil
}

func (s *merchantService) CreateMerchant(ctx context.Context, req *requests.CreateMerchantRequest) (*db.CreateMerchantRow, error) {
	const method = "CreateMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantRepository.CreateMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateMerchantRow](
			s.logger,
			merchant_errors.ErrFailedCreateMerchant,
			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteCachedMerchant(ctx, int(merchant.MerchantID))

	logSuccess("Successfully created merchant",
		zap.Int("merchantID", int(merchant.MerchantID)))

	return merchant, nil
}

func (s *merchantService) UpdateMerchant(ctx context.Context, req *requests.UpdateMerchantRequest) (*db.UpdateMerchantRow, error) {
	const method = "UpdateMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", *req.MerchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantRepository.UpdateMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantRow](
			s.logger,
			merchant_errors.ErrFailedUpdateMerchant,
			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteCachedMerchant(ctx, int(merchant.MerchantID))

	logSuccess("Successfully updated merchant",
		zap.Int("merchantID", int(merchant.MerchantID)))

	return merchant, nil
}

func (s *merchantService) TrashedMerchant(ctx context.Context, merchantID int) (*db.Merchant, error) {
	const method = "TrashedMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantRepository.TrashedMerchant(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedTrashedMerchant,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.DeleteCachedMerchant(ctx, merchantID)

	logSuccess("Successfully trashed merchant",
		zap.Int("merchantID", merchantID))

	return merchant, nil
}

func (s *merchantService) RestoreMerchant(ctx context.Context, merchantID int) (*db.Merchant, error) {
	const method = "RestoreMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantRepository.RestoreMerchant(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedRestoreMerchant,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.DeleteCachedMerchant(ctx, merchantID)

	logSuccess("Successfully restored merchant",
		zap.Int("merchantID", merchantID))

	return merchant, nil
}

func (s *merchantService) DeleteMerchantPermanent(ctx context.Context, merchantID int) (bool, error) {
	const method = "DeleteMerchantPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	success, err := s.merchantRepository.DeleteMerchantPermanent(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchant_errors.ErrFailedDeleteMerchantPermanent,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.DeleteCachedMerchant(ctx, merchantID)

	logSuccess("Successfully permanently deleted merchant",
		zap.Int("merchantID", merchantID))

	return success, nil
}

func (s *merchantService) RestoreAllMerchant(ctx context.Context) (bool, error) {
	const method = "RestoreAllMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantRepository.RestoreAllMerchant(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchant_errors.ErrFailedRestoreAllMerchants,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all trashed merchants")

	return success, nil
}

func (s *merchantService) DeleteAllMerchantPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllMerchantPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantRepository.DeleteAllMerchantPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchant_errors.ErrFailedDeleteAllMerchantsPermanent,
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all trashed merchants")

	return success, nil
}
