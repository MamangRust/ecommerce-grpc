package service

import (
	"context"
	merchantdetail_cache "ecommerce/internal/cache/merchant_detail"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	merchantdetail_errors "ecommerce/pkg/errors/merchant_detail"
	merchantsociallink_errors "ecommerce/pkg/errors/merchant_social_link_errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantDetailService struct {
	merchantDetailRepository     repository.MerchantDetailRepository
	merchantSocialLinkRepository repository.MerchantSocialLinkRepository
	logger                       logger.LoggerInterface
	cache                        merchantdetail_cache.MerchantDetailMencache
	observability                observability.TraceLoggerObservability
}

type MerchantDetailServiceDeps struct {
	MerchantDetailRepository     repository.MerchantDetailRepository
	MerchantSocialLinkRepository repository.MerchantSocialLinkRepository
	Logger                       logger.LoggerInterface
	Cache                        merchantdetail_cache.MerchantDetailMencache
	Observability                observability.TraceLoggerObservability
}

func NewMerchantDetailService(deps MerchantDetailServiceDeps) *merchantDetailService {
	return &merchantDetailService{
		merchantDetailRepository:     deps.MerchantDetailRepository,
		merchantSocialLinkRepository: deps.MerchantSocialLinkRepository,
		logger:                       deps.Logger,
		cache:                        deps.Cache,
		observability:                deps.Observability,
	}
}

func (s *merchantDetailService) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedMerchantDetailAll(ctx, req); found {
		logSuccess("Successfully retrieved all merchant detail records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantDetailRepository.FindAllMerchants(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantDetailsRow](
			s.logger,
			merchantdetail_errors.ErrFailedFindAllMerchantDetail,
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

	s.cache.SetCachedMerchantDetailAll(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched all merchant details",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantDetailService) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsActiveRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedMerchantDetailActive(ctx, req); found {
		logSuccess("Successfully retrieved active merchant detail records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantDetailRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantDetailsActiveRow](
			s.logger,
			merchantdetail_errors.ErrFailedFindActiveMerchantDetail,
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

	s.cache.SetCachedMerchantDetailActive(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched active merchant details",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantDetailService) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsTrashedRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedMerchantDetailTrashed(ctx, req); found {
		logSuccess("Successfully retrieved trashed merchant detail records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantDetailRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantDetailsTrashedRow](
			s.logger,
			merchantdetail_errors.ErrFailedFindTrashedMerchantDetail,
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

	s.cache.SetCachedMerchantDetailTrashed(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched trashed merchant details",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantDetailService) FindById(ctx context.Context, merchantID int) (*db.GetMerchantDetailRow, error) {
	const method = "FindMerchantDetailById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMerchantDetail(ctx, merchantID); found {
		logSuccess("Successfully retrieved merchant detail by ID from cache",
			zap.Int("merchantID", merchantID))
		return data, nil
	}

	merchant, err := s.merchantDetailRepository.FindById(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetMerchantDetailRow](
			s.logger,
			merchantdetail_errors.ErrFailedFindMerchantDetailById,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.SetCachedMerchantDetail(ctx, merchant)

	logSuccess("Successfully fetched merchant detail by ID",
		zap.Int("merchantID", merchantID))

	return merchant, nil
}

func (s *merchantDetailService) CreateMerchantDetail(ctx context.Context, req *requests.CreateMerchantDetailRequest) (*db.CreateMerchantDetailRow, error) {
	const method = "CreateMerchantDetail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantDetailRepository.CreateMerchantDetail(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateMerchantDetailRow](
			s.logger,
			merchantdetail_errors.ErrFailedCreateMerchantDetail,
			method,
			span,

			zap.Any("request", req),
		)
	}

	merchantId := int(merchant.MerchantDetailID)

	for _, social := range req.SocialLink {
		social.MerchantDetailID = &merchantId
		_, err := s.merchantSocialLinkRepository.CreateSocialLink(ctx, social)
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.CreateMerchantDetailRow](
				s.logger,
				merchantsociallink_errors.ErrFailedCreateMerchantSocialLink,
				method,
				span,

				zap.Any("social_link", social),
			)
		}
	}

	s.cache.DeleteMerchantDetailCache(ctx, int(merchant.MerchantDetailID))

	logSuccess("Successfully created merchant detail",
		zap.Int("merchantID", int(merchant.MerchantDetailID)))

	return merchant, nil
}

func (s *merchantDetailService) UpdateMerchantDetail(ctx context.Context, req *requests.UpdateMerchantDetailRequest) (*db.UpdateMerchantDetailRow, error) {
	const method = "UpdateMerchantDetail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", *req.MerchantDetailID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantDetailRepository.UpdateMerchantDetail(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantDetailRow](
			s.logger,
			merchantdetail_errors.ErrFailedUpdateMerchantDetail,
			method,
			span,

			zap.Any("request", req),
		)
	}

	merchantId := int(merchant.MerchantID)

	for _, social := range req.SocialLink {
		social.MerchantDetailID = &merchantId
		_, err := s.merchantSocialLinkRepository.UpdateSocialLink(ctx, social)
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.UpdateMerchantDetailRow](
				s.logger,
				merchantsociallink_errors.ErrFailedUpdateMerchantSocialLink,
				method,
				span,

				zap.Any("social_link", social),
			)
		}
	}

	s.cache.DeleteMerchantDetailCache(ctx, int(merchant.MerchantDetailID))

	logSuccess("Successfully updated merchant detail",
		zap.Int("merchantID", int(merchant.MerchantDetailID)))

	return merchant, nil
}

func (s *merchantDetailService) TrashedMerchantDetail(ctx context.Context, merchantID int) (*db.MerchantDetail, error) {
	const method = "TrashedMerchantDetail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantDetailRepository.TrashedMerchantDetail(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantDetail](
			s.logger,
			merchantdetail_errors.ErrFailedTrashedMerchantDetail,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	_, err = s.merchantSocialLinkRepository.TrashSocialLink(ctx, int(merchant.MerchantID))
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantDetail](
			s.logger,
			merchantsociallink_errors.ErrFailedTrashMerchantSocialLink,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.DeleteMerchantDetailCache(ctx, merchantID)

	logSuccess("Successfully trashed merchant detail",
		zap.Int("merchantID", merchantID))

	return merchant, nil
}

func (s *merchantDetailService) RestoreMerchantDetail(ctx context.Context, merchantID int) (*db.MerchantDetail, error) {
	const method = "RestoreMerchantDetail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantDetailRepository.RestoreMerchantDetail(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantDetail](
			s.logger,
			merchantdetail_errors.ErrFailedRestoreMerchantDetail,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	_, err = s.merchantSocialLinkRepository.RestoreSocialLink(ctx, int(merchant.MerchantDetailID))
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantDetail](
			s.logger,
			merchantsociallink_errors.ErrFailedRestoreMerchantSocialLink,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.DeleteMerchantDetailCache(ctx, merchantID)

	logSuccess("Successfully restored merchant detail",
		zap.Int("merchantID", merchantID))

	return merchant, nil
}

func (s *merchantDetailService) DeleteMerchantDetailPermanent(ctx context.Context, merchantID int) (bool, error) {
	const method = "DeleteMerchantDetailPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantDetailRepository.FindByIdTrashed(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantdetail_errors.ErrFailedFindMerchantDetailById,
			method,
			span,

			zap.Int("merchant_detail_id", merchantID),
		)
	}

	if merchant.CoverImageUrl != nil && *merchant.CoverImageUrl != "" {
		err := os.Remove(*merchant.CoverImageUrl)
		if err != nil {
			status = "error"

			if os.IsNotExist(err) {
				return errorhandler.HandleError[bool](
					s.logger,
					merchantdetail_errors.ErrFailedImageNotFound,
					method,
					span,
					zap.String("cover_image_path", *merchant.CoverImageUrl),
				)
			}

			return errorhandler.HandleError[bool](
				s.logger,
				merchantdetail_errors.ErrFailedRemoveImageMerchantDetail,
				method,
				span,
				zap.String("cover_image_path", *merchant.CoverImageUrl),
			)
		}
	}

	if merchant.LogoUrl != nil && *merchant.LogoUrl != "" {
		err := os.Remove(*merchant.LogoUrl)
		if err != nil {
			status = "error"

			if os.IsNotExist(err) {
				return errorhandler.HandleError[bool](
					s.logger,
					merchantdetail_errors.ErrFailedLogoNotFound,
					method,
					span,
					zap.String("logo_path", *merchant.LogoUrl),
				)
			}

			return errorhandler.HandleError[bool](
				s.logger,
				merchantdetail_errors.ErrFailedRemoveImageMerchantDetail,
				method,
				span,
				zap.String("logo_path", *merchant.LogoUrl),
			)
		}
	}

	success, err := s.merchantDetailRepository.DeleteMerchantDetailPermanent(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantdetail_errors.ErrFailedDeleteMerchantDetailPermanent,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	_, err = s.merchantSocialLinkRepository.DeletePermanentSocialLink(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantsociallink_errors.ErrFailedDeletePermanentMerchantSocialLink,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	// Invalidate cache
	s.cache.DeleteMerchantDetailCache(ctx, merchantID)

	logSuccess("Successfully permanently deleted merchant detail",
		zap.Int("merchantID", merchantID))

	return success, nil
}

func (s *merchantDetailService) RestoreAllMerchantDetail(ctx context.Context) (bool, error) {
	const method = "RestoreAllMerchantDetail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantDetailRepository.RestoreAllMerchantDetail(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantdetail_errors.ErrFailedRestoreAllMerchantDetail,
			method,
			span,
		)
	}

	_, err = s.merchantSocialLinkRepository.RestoreAllSocialLink(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantsociallink_errors.ErrFailedRestoreAllMerchantSocialLinks,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all trashed merchant details")

	return success, nil
}

func (s *merchantDetailService) DeleteAllMerchantDetailPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllMerchantDetailPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantDetailRepository.DeleteAllMerchantDetailPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantdetail_errors.ErrFailedDeleteAllMerchantDetailPermanent,
			method,
			span,
		)
	}

	_, err = s.merchantSocialLinkRepository.DeleteAllPermanentSocialLink(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantsociallink_errors.ErrFailedDeleteAllPermanentMerchantSocialLinks,
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all trashed merchant details")

	return success, nil
}
