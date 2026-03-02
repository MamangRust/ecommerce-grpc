package service

import (
	"context"
	banner_cache "ecommerce/internal/cache/banner"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/banner_errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type bannerService struct {
	bannerRepository repository.BannerRepository
	observability    observability.TraceLoggerObservability
	logger           logger.LoggerInterface
	cache            banner_cache.BannerMencache
}

type BannerServiceDeps struct {
	BannerRepository repository.BannerRepository
	Observability    observability.TraceLoggerObservability
	Logger           logger.LoggerInterface
	Cache            banner_cache.BannerMencache
}

func NewBannerService(deps BannerServiceDeps) *bannerService {
	return &bannerService{
		bannerRepository: deps.BannerRepository,
		observability:    deps.Observability,
		logger:           deps.Logger,
		cache:            deps.Cache,
	}
}

func (s *bannerService) FindAll(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersRow, *int, error) {
	const method = "FindAllBanners"

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

	if data, total, found := s.cache.GetCachedBannersCache(ctx, req); found {
		logSuccess("Successfully retrieved all banner records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	banners, err := s.bannerRepository.FindAllBanners(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetBannersRow](
			s.logger,
			banner_errors.ErrFailedFindAllBanners,
			method,
			span,

			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int

	if len(banners) > 0 {
		totalCount = int(banners[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedBannersCache(ctx, req, banners, &totalCount)

	logSuccess("Successfully fetched all banners",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return banners, &totalCount, nil
}

func (s *bannerService) FindByActive(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersActiveRow, *int, error) {
	const method = "FindActiveBanners"

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

	if data, total, found := s.cache.GetCachedBannerActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active banner records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	banners, err := s.bannerRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetBannersActiveRow](
			s.logger,
			banner_errors.ErrFailedFindActiveBanners,
			method,
			span,

			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int

	if len(banners) > 0 {
		totalCount = int(banners[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedBannerActiveCache(ctx, req, banners, &totalCount)

	logSuccess("Successfully fetched active banners",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return banners, &totalCount, nil
}

func (s *bannerService) FindByTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersTrashedRow, *int, error) {
	const method = "FindTrashedBanners"

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

	if data, total, found := s.cache.GetCachedBannerTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved trashed banner records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	banners, err := s.bannerRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetBannersTrashedRow](
			s.logger,
			banner_errors.ErrFailedFindTrashedBanners,
			method,
			span,

			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int

	if len(banners) > 0 {
		totalCount = int(banners[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedBannerTrashedCache(ctx, req, banners, &totalCount)

	logSuccess("Successfully fetched trashed banners",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return banners, &totalCount, nil
}

func (s *bannerService) FindById(ctx context.Context, bannerID int) (*db.GetBannerRow, error) {
	const method = "FindByIdBanner"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("bannerID", bannerID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedBannerCache(ctx, bannerID); found {
		logSuccess("Successfully retrieved banner from cache", zap.Int("bannerID", bannerID))
		return data, nil
	}

	res, err := s.bannerRepository.FindById(ctx, bannerID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetBannerRow](
			s.logger,
			banner_errors.ErrBannerNotFound,
			method,
			span,

			zap.Int("banner_id", bannerID),
		)
	}

	s.cache.SetCachedBannerCache(ctx, res)

	logSuccess("Successfully fetched banner", zap.Int("bannerID", bannerID))
	return res, nil
}

func (s *bannerService) CreateBanner(ctx context.Context, req *requests.CreateBannerRequest) (*db.CreateBannerRow, error) {
	const method = "CreateBanner"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	banner, err := s.bannerRepository.CreateBanner(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateBannerRow](
			s.logger,
			banner_errors.ErrFailedCreateBanner,
			method,
			span,

			zap.Any("request", req),
		)
	}

	logSuccess("Successfully created banner", zap.Int("bannerID", int(banner.BannerID)))
	return banner, nil
}

func (s *bannerService) UpdateBanner(ctx context.Context, req *requests.UpdateBannerRequest) (*db.UpdateBannerRow, error) {
	const method = "UpdateBanner"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("bannerID", *req.BannerID))

	defer func() {
		end(status)
	}()

	banner, err := s.bannerRepository.UpdateBanner(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateBannerRow](
			s.logger,
			banner_errors.ErrFailedUpdateBanner,
			method,
			span,

			zap.Int("banner_id", *req.BannerID),
		)
	}

	s.cache.DeleteBannerCache(ctx, *req.BannerID)

	logSuccess("Successfully updated banner", zap.Int("bannerID", *req.BannerID))
	return banner, nil
}

func (s *bannerService) TrashedBanner(ctx context.Context, bannerID int) (*db.Banner, error) {
	const method = "TrashedBanner"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("bannerID", bannerID))

	defer func() {
		end(status)
	}()

	banner, err := s.bannerRepository.TrashedBanner(ctx, bannerID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Banner](
			s.logger,
			banner_errors.ErrFailedTrashedBanner,
			method,
			span,

			zap.Int("banner_id", bannerID),
		)
	}

	s.cache.DeleteBannerCache(ctx, bannerID)

	logSuccess("Successfully trashed banner", zap.Int("bannerID", bannerID))
	return banner, nil
}

func (s *bannerService) RestoreBanner(ctx context.Context, bannerID int) (*db.Banner, error) {
	const method = "RestoreBanner"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("bannerID", bannerID))

	defer func() {
		end(status)
	}()

	banner, err := s.bannerRepository.RestoreBanner(ctx, bannerID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Banner](
			s.logger,
			banner_errors.ErrFailedRestoreBanner,
			method,
			span,

			zap.Int("banner_id", bannerID),
		)
	}

	s.cache.DeleteBannerCache(ctx, bannerID)

	logSuccess("Successfully restored banner", zap.Int("bannerID", bannerID))
	return banner, nil
}

func (s *bannerService) DeleteBannerPermanent(ctx context.Context, bannerID int) (bool, error) {
	const method = "DeleteBannerPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("bannerID", bannerID))

	defer func() {
		end(status)
	}()

	success, err := s.bannerRepository.DeleteBannerPermanent(ctx, bannerID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			banner_errors.ErrFailedDeleteBanner,
			method,
			span,

			zap.Int("banner_id", bannerID),
		)
	}

	s.cache.DeleteBannerCache(ctx, bannerID)

	logSuccess("Successfully deleted banner permanently", zap.Int("bannerID", bannerID))
	return success, nil
}

func (s *bannerService) RestoreAllBanner(ctx context.Context) (bool, error) {
	const method = "RestoreAllBanner"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.bannerRepository.RestoreAllBanner(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			banner_errors.ErrFailedRestoreAllBanners,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all trashed banners")
	return success, nil
}

func (s *bannerService) DeleteAllBannerPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllBannerPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.bannerRepository.DeleteAllBannerPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			banner_errors.ErrFailedDeleteAllBanners,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all trashed banners permanently")
	return success, nil
}
