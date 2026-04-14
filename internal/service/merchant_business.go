package service

import (
	"context"
	merchantbusiness_cache "ecommerce/internal/cache/merchant_business"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	merchantbusiness_errors "ecommerce/pkg/errors/merchant_business"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantBusinessService struct {
	merchantBusinessRepository repository.MerchantBusinessRepository
	logger                     logger.LoggerInterface
	observability              observability.TraceLoggerObservability
	cache                      merchantbusiness_cache.MerchantBusinessMencache
}

type MerchantBusinessServiceDeps struct {
	MerchantBusinessRepository repository.MerchantBusinessRepository
	Logger                     logger.LoggerInterface
	Observability              observability.TraceLoggerObservability
	Cache                      merchantbusiness_cache.MerchantBusinessMencache
}

func NewMerchantBusinessService(deps MerchantBusinessServiceDeps) *merchantBusinessService {
	return &merchantBusinessService{
		merchantBusinessRepository: deps.MerchantBusinessRepository,
		logger:                     deps.Logger,
		observability:              deps.Observability,
		cache:                      deps.Cache,
	}
}

func (s *merchantBusinessService) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedMerchantBusinessAll(ctx, req); found {
		logSuccess("Successfully retrieved all merchant business records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantBusinessRepository.FindAllMerchants(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantsBusinessInformationRow](
			s.logger,
			merchantbusiness_errors.ErrFailedFindAllMerchantBusiness,
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

	s.cache.SetCachedMerchantBusinessAll(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched all merchant businesses",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantBusinessService) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationActiveRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedMerchantBusinessActive(ctx, req); found {
		logSuccess("Successfully retrieved active merchant business records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantBusinessRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantsBusinessInformationActiveRow](
			s.logger,
			merchantbusiness_errors.ErrFailedFindActiveMerchantBusiness,
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

	s.cache.SetCachedMerchantBusinessActive(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched active merchant businesses",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantBusinessService) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationTrashedRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedMerchantBusinessTrashed(ctx, req); found {
		logSuccess("Successfully retrieved trashed merchant business records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantBusinessRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantsBusinessInformationTrashedRow](
			s.logger,
			merchantbusiness_errors.ErrFailedFindTrashedMerchantBusiness,
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

	s.cache.SetCachedMerchantBusinessTrashed(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched trashed merchant businesses",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantBusinessService) FindById(ctx context.Context, merchantID int) (*db.GetMerchantBusinessInformationRow, error) {
	const method = "FindMerchantBusinessById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMerchantBusiness(ctx, merchantID); found {
		logSuccess("Successfully retrieved merchant business by ID from cache",
			zap.Int("merchantID", merchantID))
		return data, nil
	}

	merchant, err := s.merchantBusinessRepository.FindById(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetMerchantBusinessInformationRow](
			s.logger,
			merchantbusiness_errors.ErrFailedFindMerchantBusinessById,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.SetCachedMerchantBusiness(ctx, merchant)

	logSuccess("Successfully fetched merchant business by ID",
		zap.Int("merchantID", merchantID))

	return merchant, nil
}

func (s *merchantBusinessService) CreateMerchantBusiness(ctx context.Context, req *requests.CreateMerchantBusinessInformationRequest) (*db.CreateMerchantBusinessInformationRow, error) {
	const method = "CreateMerchantBusiness"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantBusinessRepository.CreateMerchantBusiness(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateMerchantBusinessInformationRow](
			s.logger,
			merchantbusiness_errors.ErrFailedCreateMerchantBusiness,
			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteMerchantBusinessCache(ctx, int(merchant.MerchantBusinessInfoID))

	logSuccess("Successfully created merchant business",
		zap.Int("merchantID", int(merchant.MerchantBusinessInfoID)))

	return merchant, nil
}

func (s *merchantBusinessService) UpdateMerchantBusiness(ctx context.Context, req *requests.UpdateMerchantBusinessInformationRequest) (*db.UpdateMerchantBusinessInformationRow, error) {
	const method = "UpdateMerchantBusiness"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", *req.MerchantBusinessInfoID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantBusinessRepository.UpdateMerchantBusiness(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantBusinessInformationRow](
			s.logger,
			merchantbusiness_errors.ErrFailedUpdateMerchantBusiness,
			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteMerchantBusinessCache(ctx, int(merchant.MerchantBusinessInfoID))

	logSuccess("Successfully updated merchant business",
		zap.Int("merchantID", int(merchant.MerchantBusinessInfoID)))

	return merchant, nil
}

func (s *merchantBusinessService) TrashedMerchantBusiness(ctx context.Context, merchantID int) (*db.MerchantBusinessInformation, error) {
	const method = "TrashedMerchantBusiness"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantBusinessRepository.TrashedMerchantBusiness(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantBusinessInformation](
			s.logger,
			merchantbusiness_errors.ErrFailedTrashedMerchantBusiness,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.DeleteMerchantBusinessCache(ctx, merchantID)

	logSuccess("Successfully trashed merchant business",
		zap.Int("merchantID", merchantID))

	return merchant, nil
}

func (s *merchantBusinessService) RestoreMerchantBusiness(ctx context.Context, merchantID int) (*db.MerchantBusinessInformation, error) {
	const method = "RestoreMerchantBusiness"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantBusinessRepository.RestoreMerchantBusiness(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantBusinessInformation](
			s.logger,
			merchantbusiness_errors.ErrFailedRestoreMerchantBusiness,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.DeleteMerchantBusinessCache(ctx, merchantID)

	logSuccess("Successfully restored merchant business",
		zap.Int("merchantID", merchantID))

	return merchant, nil
}

func (s *merchantBusinessService) DeleteMerchantBusinessPermanent(ctx context.Context, merchantID int) (bool, error) {
	const method = "DeleteMerchantBusinessPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantBusinessRepository.DeleteMerchantBusinessPermanent(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantbusiness_errors.ErrFailedDeleteMerchantBusinessPermanent,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.DeleteMerchantBusinessCache(ctx, merchantID)

	logSuccess("Successfully permanently deleted merchant business",
		zap.Int("merchantID", merchantID))

	return success, nil
}

func (s *merchantBusinessService) RestoreAllMerchantBusiness(ctx context.Context) (bool, error) {
	const method = "RestoreAllMerchantBusiness"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantBusinessRepository.RestoreAllMerchantBusiness(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantbusiness_errors.ErrFailedRestoreAllMerchantBusiness,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all trashed merchant businesses")

	return success, nil
}

func (s *merchantBusinessService) DeleteAllMerchantBusinessPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllMerchantBusinessPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantBusinessRepository.DeleteAllMerchantBusinessPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantbusiness_errors.ErrFailedDeleteAllMerchantBusinessPermanent,
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all trashed merchant businesses")

	return success, nil
}
