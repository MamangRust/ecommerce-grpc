package service

import (
	"context"
	reviewdetail_cache "ecommerce/internal/cache/review_detail"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	reviewdetail_errors "ecommerce/pkg/errors/review_detail"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type reviewDetailService struct {
	reviewDetailRepository repository.ReviewDetailRepository
	logger                 logger.LoggerInterface
	cache                  reviewdetail_cache.ReviewDetailMencache
	observability          observability.TraceLoggerObservability
}

type ReviewDetailServiceDeps struct {
	ReviewDetailRepository repository.ReviewDetailRepository
	Logger                 logger.LoggerInterface
	Cache                  reviewdetail_cache.ReviewDetailMencache
	Observability          observability.TraceLoggerObservability
}

func NewReviewDetailService(deps ReviewDetailServiceDeps) *reviewDetailService {
	return &reviewDetailService{
		reviewDetailRepository: deps.ReviewDetailRepository,
		logger:                 deps.Logger,
		cache:                  deps.Cache,
		observability:          deps.Observability,
	}
}

func (s *reviewDetailService) FindAllReviews(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsRow, *int, error) {
	const method = "FindAllReviews"

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

	if data, total, found := s.cache.GetReviewDetailAllCache(ctx, req); found {
		logSuccess("Successfully retrieved all review detail records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	reviewDetails, err := s.reviewDetailRepository.FindAllReviews(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetReviewDetailsRow](
			s.logger,
			reviewdetail_errors.ErrFailedFindAllReview,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(reviewDetails) > 0 {
		totalCount = int(reviewDetails[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetReviewDetailAllCache(ctx, req, reviewDetails, &totalCount)

	logSuccess("Successfully fetched all review details",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return reviewDetails, &totalCount, nil
}

func (s *reviewDetailService) FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsActiveRow, *int, error) {
	const method = "FindByActiveReviews"

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

	if data, total, found := s.cache.GetReviewDetailActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active review detail records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	reviewDetails, err := s.reviewDetailRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetReviewDetailsActiveRow](
			s.logger,
			reviewdetail_errors.ErrFailedFindActiveReview,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(reviewDetails) > 0 {
		totalCount = int(reviewDetails[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetReviewDetailActiveCache(ctx, req, reviewDetails, &totalCount)

	logSuccess("Successfully fetched active review details",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return reviewDetails, &totalCount, nil
}

func (s *reviewDetailService) FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsTrashedRow, *int, error) {
	const method = "FindByTrashedReviews"

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

	if data, total, found := s.cache.GetReviewDetailTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved trashed review detail records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	reviewDetails, err := s.reviewDetailRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetReviewDetailsTrashedRow](
			s.logger,
			reviewdetail_errors.ErrFailedFindTrashedReview,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(reviewDetails) > 0 {
		totalCount = int(reviewDetails[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetReviewDetailTrashedCache(ctx, req, reviewDetails, &totalCount)

	logSuccess("Successfully fetched trashed review details",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return reviewDetails, &totalCount, nil
}

func (s *reviewDetailService) FindById(ctx context.Context, review_id int) (*db.GetReviewDetailRow, error) {
	const method = "FindReviewDetailById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("review_id", review_id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedReviewDetailCache(ctx, review_id); found {
		logSuccess("Successfully retrieved review detail by ID from cache",
			zap.Int("review_id", review_id))
		return data, nil
	}

	reviewDetail, err := s.reviewDetailRepository.FindById(ctx, review_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetReviewDetailRow](
			s.logger,
			reviewdetail_errors.ErrReviewDetailNotFoundRes,
			method,
			span,

			zap.Int("review_id", review_id),
		)
	}

	s.cache.SetCachedReviewDetailCache(ctx, reviewDetail)

	logSuccess("Successfully fetched review detail by ID",
		zap.Int("review_id", review_id))

	return reviewDetail, nil
}

func (s *reviewDetailService) CreateReviewDetail(ctx context.Context, req *requests.CreateReviewDetailRequest) (*db.CreateReviewDetailRow, error) {
	const method = "CreateReviewDetail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	reviewDetail, err := s.reviewDetailRepository.CreateReviewDetail(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateReviewDetailRow](
			s.logger,
			reviewdetail_errors.ErrFailedCreateReviewDetail,
			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteReviewDetailCache(ctx, int(reviewDetail.ReviewDetailID))

	logSuccess("Successfully created review detail",
		zap.Int("review_detail_id", int(reviewDetail.ReviewDetailID)))

	return reviewDetail, nil
}

func (s *reviewDetailService) UpdateReviewDetail(ctx context.Context, req *requests.UpdateReviewDetailRequest) (*db.UpdateReviewDetailRow, error) {
	const method = "UpdateReviewDetail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("review_detail_id", *req.ReviewDetailID))

	defer func() {
		end(status)
	}()

	reviewDetail, err := s.reviewDetailRepository.UpdateReviewDetail(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateReviewDetailRow](
			s.logger,
			reviewdetail_errors.ErrFailedUpdateReviewDetail,
			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteReviewDetailCache(ctx, int(reviewDetail.ReviewDetailID))

	logSuccess("Successfully updated review detail",
		zap.Int("review_detail_id", int(reviewDetail.ReviewDetailID)))

	return reviewDetail, nil
}

func (s *reviewDetailService) TrashedReviewDetail(ctx context.Context, review_id int) (*db.ReviewDetail, error) {
	const method = "TrashedReviewDetail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("review_id", review_id))

	defer func() {
		end(status)
	}()

	reviewDetail, err := s.reviewDetailRepository.TrashedReviewDetail(ctx, review_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.ReviewDetail](
			s.logger,
			reviewdetail_errors.ErrFailedTrashedReviewDetail,
			method,
			span,

			zap.Int("review_id", review_id),
		)
	}

	s.cache.DeleteReviewDetailCache(ctx, int(reviewDetail.ReviewDetailID))

	logSuccess("Successfully trashed review detail",
		zap.Int("review_id", review_id))

	return reviewDetail, nil
}

func (s *reviewDetailService) RestoreReviewDetail(ctx context.Context, review_id int) (*db.ReviewDetail, error) {
	const method = "RestoreReviewDetail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("review_id", review_id))

	defer func() {
		end(status)
	}()

	reviewDetail, err := s.reviewDetailRepository.RestoreReviewDetail(ctx, review_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.ReviewDetail](
			s.logger,
			reviewdetail_errors.ErrFailedRestoreReviewDetail,
			method,
			span,

			zap.Int("review_id", review_id),
		)
	}

	s.cache.DeleteReviewDetailCache(ctx, int(reviewDetail.ReviewDetailID))

	logSuccess("Successfully restored review detail",
		zap.Int("review_id", review_id))

	return reviewDetail, nil
}

func (s *reviewDetailService) DeleteReviewDetailPermanent(ctx context.Context, review_id int) (bool, error) {
	const method = "DeleteReviewDetailPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("review_id", review_id))

	defer func() {
		end(status)
	}()

	reviewDetail, err := s.reviewDetailRepository.FindByIdTrashed(ctx, review_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			reviewdetail_errors.ErrFailedDeletePermanentReview,
			method,
			span,

			zap.Int("review_id", review_id),
		)
	}

	if reviewDetail.Url != "" {
		err := os.Remove(reviewDetail.Url)
		if err != nil {
			if os.IsNotExist(err) {
				status = "error"
				return errorhandler.HandleError[bool](
					s.logger,
					reviewdetail_errors.ErrFailedImageNotFound,
					method,
					span,
					zap.String("upload_path", reviewDetail.Url),
				)
			} else {
				status = "error"
				return errorhandler.HandleError[bool](
					s.logger,
					reviewdetail_errors.ErrFailedRemoveImage,
					method,
					span,
					zap.String("upload_path", reviewDetail.Url),
				)
			}
		}
	}

	success, err := s.reviewDetailRepository.DeleteReviewDetailPermanent(ctx, review_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			reviewdetail_errors.ErrFailedDeletePermanentReview,
			method,
			span,

			zap.Int("review_id", review_id),
		)
	}

	s.cache.DeleteReviewDetailCache(ctx, review_id)

	logSuccess("Successfully permanently deleted review detail",
		zap.Int("review_id", review_id))

	return success, nil
}

func (s *reviewDetailService) RestoreAllReviewDetail(ctx context.Context) (bool, error) {
	const method = "RestoreAllReviewDetail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.reviewDetailRepository.RestoreAllReviewDetail(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			reviewdetail_errors.ErrFailedRestoreAllReviewDetail,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all trashed review details")

	return success, nil
}

func (s *reviewDetailService) DeleteAllReviewDetailPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllReviewDetailPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.reviewDetailRepository.DeleteAllReviewDetailPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			reviewdetail_errors.ErrFailedDeleteAllReviewDetail,
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all trashed review details")

	return success, nil
}
