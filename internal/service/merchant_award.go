package service

import (
	"context"
	merchantawards_cache "ecommerce/internal/cache/merchant_awards"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	merchantaward_errors "ecommerce/pkg/errors/merchant_award"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantAwardService struct {
	merchantAwardRepository repository.MerchantAwardRepository
	logger                  logger.LoggerInterface
	observability           observability.TraceLoggerObservability
	cache                   merchantawards_cache.MerchantAwardMencache
}

type MerchantAwardServiceDeps struct {
	MerchantAwardRepository repository.MerchantAwardRepository
	Logger                  logger.LoggerInterface
	Observability           observability.TraceLoggerObservability
	Cache                   merchantawards_cache.MerchantAwardMencache
}

func NewMerchantAwardService(deps MerchantAwardServiceDeps) *merchantAwardService {
	return &merchantAwardService{
		merchantAwardRepository: deps.MerchantAwardRepository,
		logger:                  deps.Logger,
		observability:           deps.Observability,
		cache:                   deps.Cache,
	}
}

func (s *merchantAwardService) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsRow, *int, error) {
	const method = "FindAll"

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

	merchants, err := s.merchantAwardRepository.FindAllMerchants(ctx, req)

	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantCertificationsAndAwardsRow](
			s.logger,
			merchantaward_errors.ErrFailedFindAllMerchantAwards,
			method,
			span,

			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int

	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	} else {
		totalCount = 0
	}

	logSuccess("Successfully fetched merchants",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantAwardService) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsActiveRow, *int, error) {
	const method = "FindByActive"

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

	merchants, err := s.merchantAwardRepository.FindByActive(ctx, req)

	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantCertificationsAndAwardsActiveRow](
			s.logger,
			merchantaward_errors.ErrFailedFindActiveMerchantAwards,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	} else {
		totalCount = 0
	}

	logSuccess("Successfully fetched active merchants",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantAwardService) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsTrashedRow, *int, error) {
	const method = "FindByTrashed"

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

	merchants, err := s.merchantAwardRepository.FindByTrashed(ctx, req)

	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantCertificationsAndAwardsTrashedRow](
			s.logger,
			merchantaward_errors.ErrFailedFindTrashedMerchantAwards,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	} else {
		totalCount = 0
	}

	logSuccess("Successfully fetched trashed merchants",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantAwardService) FindById(ctx context.Context, merchantID int) (*db.GetMerchantCertificationOrAwardRow, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardRepository.FindById(ctx, merchantID)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetMerchantCertificationOrAwardRow](
			s.logger,
			merchantaward_errors.ErrFailedFindMerchantAwardById,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	logSuccess("Successfully fetched merchant by ID", zap.Int("merchantID", merchantID))

	return merchant, nil
}

const method = "CreateMerchant"

func (s *merchantAwardService) CreateMerchantAward(ctx context.Context, req *requests.CreateMerchantCertificationOrAwardRequest) (*db.CreateMerchantCertificationOrAwardRow, error) {

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardRepository.CreateMerchantAward(ctx, req)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateMerchantCertificationOrAwardRow](
			s.logger,
			merchantaward_errors.ErrFailedCreateMerchantAward,
			method,
			span,

			zap.Any("request", req),
		)
	}

	logSuccess("Successfully created new merchant")

	return merchant, nil
}

func (s *merchantAwardService) UpdateMerchantAward(ctx context.Context, req *requests.UpdateMerchantCertificationOrAwardRequest) (*db.UpdateMerchantCertificationOrAwardRow, error) {
	const method = "UpdateMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", *req.MerchantCertificationID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardRepository.UpdateMerchantAward(ctx, req)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantCertificationOrAwardRow](
			s.logger,
			merchantaward_errors.ErrFailedUpdateMerchantAward,
			method,
			span,

			zap.Any("request", req),
		)
	}

	logSuccess("Successfully updated merchant", zap.Int("merchantID", *req.MerchantCertificationID))

	return merchant, nil
}

func (s *merchantAwardService) TrashedMerchantAward(ctx context.Context, merchantID int) (*db.MerchantCertificationsAndAward, error) {
	const method = "TrashedMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardRepository.TrashedMerchantAward(ctx, merchantID)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantCertificationsAndAward](
			s.logger,
			merchantaward_errors.ErrFailedTrashedMerchantAward,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	logSuccess("Successfully moved merchant to trash", zap.Int("merchantID", merchantID))

	return merchant, nil
}

func (s *merchantAwardService) RestoreMerchantAward(ctx context.Context, merchantID int) (*db.MerchantCertificationsAndAward, error) {
	const method = "RestoreMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardRepository.RestoreMerchantAward(ctx, merchantID)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantCertificationsAndAward](
			s.logger,
			merchantaward_errors.ErrFailedRestoreMerchantAward,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	logSuccess("Successfully restored merchant from trash", zap.Int("merchantID", merchantID))

	return merchant, nil
}

func (s *merchantAwardService) DeleteMerchantPermanent(ctx context.Context, merchantID int) (bool, error) {
	const method = "DeleteMerchantPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	success, err := s.merchantAwardRepository.DeleteMerchantPermanent(ctx, merchantID)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantaward_errors.ErrFailedDeleteMerchantAwardPermanent,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	logSuccess("Successfully permanently deleted merchant", zap.Int("merchantID", merchantID))

	return success, nil
}

func (s *merchantAwardService) RestoreAllMerchantAward(ctx context.Context) (bool, error) {
	const method = "RestoreAllMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantAwardRepository.RestoreAllMerchantAward(ctx)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantaward_errors.ErrFailedRestoreAllMerchantAwards,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all trashed merchants")

	return success, nil
}

func (s *merchantAwardService) DeleteAllMerchantAwardPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllMerchantPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantAwardRepository.DeleteAllMerchantAwardPermanent(ctx)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantaward_errors.ErrFailedDeleteAllMerchantAwardsPermanent,
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all merchants")

	return success, nil
}
