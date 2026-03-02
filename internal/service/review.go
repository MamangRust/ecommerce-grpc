package service

import (
	"context"
	review_cache "ecommerce/internal/cache/review"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/product_errors"
	review_errors "ecommerce/pkg/errors/review"
	"ecommerce/pkg/errors/user_errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type reviewService struct {
	reviewRepository  repository.ReviewRepository
	productRepository repository.ProductRepository
	userRepository    repository.UserRepository
	logger            logger.LoggerInterface
	observability     observability.TraceLoggerObservability
	cache             review_cache.ReviewMencache
}

type ReviewServiceDeps struct {
	ReviewRepository  repository.ReviewRepository
	ProductRepository repository.ProductRepository
	UserRepository    repository.UserRepository
	Logger            logger.LoggerInterface
	Observability     observability.TraceLoggerObservability
	Cache             review_cache.ReviewMencache
}

func NewReviewService(deps ReviewServiceDeps) ReviewService {
	return &reviewService{
		reviewRepository:  deps.ReviewRepository,
		productRepository: deps.ProductRepository,
		userRepository:    deps.UserRepository,
		logger:            deps.Logger,
		observability:     deps.Observability,
		cache:             deps.Cache,
	}
}

func (s *reviewService) FindAllReview(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsRow, *int, error) {
	const method = "FindAllReview"

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

	if data, total, found := s.cache.GetReviewAllCache(ctx, req); found {
		logSuccess("Successfully retrieved all review records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	reviews, err := s.reviewRepository.FindAllReview(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetReviewsRow](
			s.logger,
			review_errors.ErrFailedFindAllReviews,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(reviews) > 0 {
		totalCount = int(reviews[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetReviewAllCache(ctx, req, reviews, &totalCount)

	logSuccess("Successfully fetched all reviews",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return reviews, &totalCount, nil
}

func (s *reviewService) FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsActiveRow, *int, error) {
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

	if data, total, found := s.cache.GetReviewActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active review records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	reviews, err := s.reviewRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetReviewsActiveRow](
			s.logger,
			review_errors.ErrFailedFindActiveReviews,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(reviews) > 0 {
		totalCount = int(reviews[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetReviewActiveCache(ctx, req, reviews, &totalCount)

	logSuccess("Successfully fetched active reviews",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return reviews, &totalCount, nil
}

func (s *reviewService) FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsTrashedRow, *int, error) {
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

	if data, total, found := s.cache.GetReviewTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved trashed review records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	reviews, err := s.reviewRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetReviewsTrashedRow](
			s.logger,
			review_errors.ErrFailedFindTrashedReviews,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(reviews) > 0 {
		totalCount = int(reviews[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetReviewTrashedCache(ctx, req, reviews, &totalCount)

	logSuccess("Successfully fetched trashed reviews",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return reviews, &totalCount, nil
}

func (s *reviewService) FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*db.GetReviewByProductIdRow, *int, error) {
	const method = "FindByProductReviews"

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
		attribute.String("search", search),
		attribute.Int("productID", req.ProductID))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetReviewByProductCache(ctx, req); found {
		logSuccess("Successfully retrieved product review records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.Int("productID", req.ProductID))
		return data, total, nil
	}

	reviews, err := s.reviewRepository.FindByProduct(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetReviewByProductIdRow](
			s.logger,
			review_errors.ErrFailedFindByProductReviews,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.Int("productID", req.ProductID),
		)
	}

	var totalCount int

	if len(reviews) > 0 {
		totalCount = int(reviews[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetReviewByProductCache(ctx, req, reviews, &totalCount)

	logSuccess("Successfully fetched product reviews",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.Int("productID", req.ProductID))

	return reviews, &totalCount, nil
}

func (s *reviewService) FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*db.GetReviewByMerchantIdRow, *int, error) {
	const method = "FindByMerchantReviews"

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
		attribute.String("search", search),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetReviewByMerchantCache(ctx, req); found {
		logSuccess("Successfully retrieved merchant review records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.Int("merchantID", req.MerchantID))
		return data, total, nil
	}

	reviews, err := s.reviewRepository.FindByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetReviewByMerchantIdRow](
			s.logger,
			review_errors.ErrFailedFindByMerchantReviews,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	var totalCount int

	if len(reviews) > 0 {
		totalCount = int(reviews[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetReviewByMerchantCache(ctx, req, reviews, &totalCount)

	logSuccess("Successfully fetched merchant reviews",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.Int("merchantID", req.MerchantID))

	return reviews, &totalCount, nil
}

func (s *reviewService) FindById(ctx context.Context, id int) (*db.GetReviewByIDRow, error) {
	const method = "FindReviewById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("id", id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetReviewByIdCache(ctx, id); found {
		logSuccess("Successfully retrieved review by ID from cache",
			zap.Int("id", id))
		return data, nil
	}

	review, err := s.reviewRepository.FindById(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetReviewByIDRow](
			s.logger,
			review_errors.ErrFailedReviewNotFound,
			method,
			span,

			zap.Int("id", id),
		)
	}

	s.cache.SetReviewByIdCache(ctx, review)

	logSuccess("Successfully fetched review by ID",
		zap.Int("id", id))

	return review, nil
}

func (s *reviewService) CreateReview(ctx context.Context, req *requests.CreateReviewRequest) (*db.CreateReviewRow, error) {
	const method = "CreateReview"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", req.UserID),
		attribute.Int("product_id", req.ProductID))

	defer func() {
		end(status)
	}()

	_, err := s.userRepository.FindById(ctx, req.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateReviewRow](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,

			zap.Int("user_id", req.UserID),
		)
	}

	_, err = s.productRepository.FindById(ctx, req.ProductID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateReviewRow](
			s.logger,
			product_errors.ErrFailedFindProductById,
			method,
			span,

			zap.Int("product_id", req.ProductID),
		)
	}

	review, err := s.reviewRepository.CreateReview(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateReviewRow](
			s.logger,
			review_errors.ErrFailedCreateReview,
			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteReviewCache(ctx, int(review.ReviewID))

	logSuccess("Successfully created review",
		zap.Int("review_id", int(review.ReviewID)),
		zap.Int("user_id", req.UserID),
		zap.Int("product_id", req.ProductID))

	return review, nil
}

func (s *reviewService) UpdateReview(ctx context.Context, req *requests.UpdateReviewRequest) (*db.UpdateReviewRow, error) {
	const method = "UpdateReview"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("review_id", *req.ReviewID))

	defer func() {
		end(status)
	}()

	_, err := s.reviewRepository.FindById(ctx, *req.ReviewID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateReviewRow](
			s.logger,
			review_errors.ErrFailedReviewNotFound,
			method,
			span,

			zap.Int("review_id", *req.ReviewID),
		)
	}

	review, err := s.reviewRepository.UpdateReview(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateReviewRow](
			s.logger,
			review_errors.ErrFailedUpdateReview,
			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteReviewCache(ctx, int(review.ReviewID))

	logSuccess("Successfully updated review",
		zap.Int("review_id", int(review.ReviewID)))

	return review, nil
}

func (s *reviewService) TrashReview(ctx context.Context, reviewID int) (*db.Review, error) {
	const method = "TrashReview"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("reviewID", reviewID))

	defer func() {
		end(status)
	}()

	review, err := s.reviewRepository.TrashReview(ctx, reviewID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Review](
			s.logger,
			review_errors.ErrFailedTrashedReview,
			method,
			span,

			zap.Int("reviewID", reviewID),
		)
	}

	s.cache.DeleteReviewCache(ctx, int(review.ReviewID))

	logSuccess("Successfully trashed review",
		zap.Int("review_id", int(review.ReviewID)))

	return review, nil
}

func (s *reviewService) RestoreReview(ctx context.Context, reviewID int) (*db.Review, error) {
	const method = "RestoreReview"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("reviewID", reviewID))

	defer func() {
		end(status)
	}()

	review, err := s.reviewRepository.RestoreReview(ctx, reviewID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Review](
			s.logger,
			review_errors.ErrFailedRestoreReview,
			method,
			span,

			zap.Int("reviewID", reviewID),
		)
	}

	s.cache.DeleteReviewCache(ctx, int(review.ReviewID))

	logSuccess("Successfully restored review",
		zap.Int("review_id", int(review.ReviewID)))

	return review, nil
}

func (s *reviewService) DeleteReviewPermanently(ctx context.Context, reviewID int) (bool, error) {
	const method = "DeleteReviewPermanently"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("reviewID", reviewID))

	defer func() {
		end(status)
	}()

	success, err := s.reviewRepository.DeleteReviewPermanently(ctx, reviewID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			review_errors.ErrFailedDeletePermanentReview,
			method,
			span,

			zap.Int("reviewID", reviewID),
		)
	}

	s.cache.DeleteReviewCache(ctx, reviewID)

	logSuccess("Successfully permanently deleted review",
		zap.Int("review_id", reviewID))

	return success, nil
}

func (s *reviewService) RestoreAllReview(ctx context.Context) (bool, error) {
	const method = "RestoreAllReview"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.reviewRepository.RestoreAllReview(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			review_errors.ErrFailedRestoreAllReviews,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all trashed reviews")

	return success, nil
}

func (s *reviewService) DeleteAllPermanentReview(ctx context.Context) (bool, error) {
	const method = "DeleteAllPermanentReview"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.reviewRepository.DeleteAllPermanentReview(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			review_errors.ErrFailedDeleteAllPermanentReviews,
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all reviews")

	return success, nil
}
