package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/errors/product_errors"
	review_errors "ecommerce/pkg/errors/review"
	"ecommerce/pkg/errors/user_errors"
	"ecommerce/pkg/logger"

	"go.uber.org/zap"
)

type reviewService struct {
	reviewRepository  repository.ReviewRepository
	productRepository repository.ProductRepository
	userRepository    repository.UserRepository
	logger            logger.LoggerInterface
	mapping           response_service.ReviewResponseMapper
}

func NewReviewService(
	reviewRepository repository.ReviewRepository,
	productRepository repository.ProductRepository,
	userRepository repository.UserRepository,
) *reviewService {
	return &reviewService{
		reviewRepository:  reviewRepository,
		productRepository: productRepository,
		userRepository:    userRepository,
	}
}
func (s *reviewService) FindAllReviews(req *requests.FindAllReview) ([]*response.ReviewResponse, *int, *response.ErrorResponse) {

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching Reviews",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	Reviews, totalRecords, err := s.reviewRepository.FindAllReview(req)
	if err != nil {
		s.logger.Error("Failed to retrieve review list",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, review_errors.ErrFailedFindAllReviews
	}

	s.logger.Debug("Successfully fetched Reviews",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewsResponse(Reviews), totalRecords, nil
}

func (s *reviewService) FindByActive(req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching Reviews",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	Reviews, totalRecords, err := s.reviewRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve review active list",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, review_errors.ErrFailedFindActiveReviews
	}

	s.logger.Debug("Successfully fetched Reviews",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewsResponseDeleteAt(Reviews), totalRecords, nil
}

func (s *reviewService) FindByTrashed(req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching Reviews",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	Reviews, totalRecords, err := s.reviewRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve review trashed list",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, review_errors.ErrFailedFindTrashedReviews
	}

	s.logger.Debug("Successfully fetched Reviews",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewsResponseDeleteAt(Reviews), totalRecords, nil
}

func (s *reviewService) FindByProduct(req *requests.FindAllReviewByProduct) ([]*response.ReviewsDetailResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching Reviews",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	reviews, totalRecords, err := s.reviewRepository.FindByProduct(req)

	if err != nil {
		s.logger.Error("Failed to retrieve review product list",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, review_errors.ErrFailedFindByProductReviews
	}

	s.logger.Debug("Successfully fetched Reviews",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewsDetailResponse(reviews), totalRecords, nil
}

func (s *reviewService) FindByMerchant(req *requests.FindAllReviewByMerchant) ([]*response.ReviewsDetailResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching Reviews",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	reviews, totalRecords, err := s.reviewRepository.FindByMerchant(req)

	if err != nil {
		s.logger.Error("Failed to retrieve review product list",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, review_errors.ErrFailedFindByMerchantReviews
	}

	s.logger.Debug("Successfully fetched Reviews",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewsDetailResponse(reviews), totalRecords, nil
}

func (s *reviewService) CreateReview(req *requests.CreateReviewRequest) (*response.ReviewResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new cashier")

	_, err := s.userRepository.FindById(req.UserID)

	if err != nil {
		s.logger.Error("Failed to retrieve user details",
			zap.Error(err),
			zap.Int("user_id", req.UserID))

		return nil, user_errors.ErrUserNotFoundRes
	}

	_, err = s.productRepository.FindById(req.ProductID)

	if err != nil {
		s.logger.Error("Failed to retrieve product details",
			zap.Error(err),
			zap.Int("product_id", req.UserID))

		return nil, product_errors.ErrFailedFindProductById
	}

	review, err := s.reviewRepository.CreateReview(req)

	if err != nil {
		s.logger.Error("Failed to create new review",
			zap.Error(err),
			zap.Any("request", req))

		return nil, review_errors.ErrFailedCreateReview
	}

	return s.mapping.ToReviewResponse(review), nil
}

func (s *reviewService) UpdateReview(req *requests.UpdateReviewRequest) (*response.ReviewResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating review", zap.Int("review_id", *req.ReviewID))

	_, err := s.reviewRepository.FindById(*req.ReviewID)

	if err != nil {
		s.logger.Error("Failed to retrieve review details",
			zap.Error(err),
			zap.Int("review_id", *req.ReviewID))

		return nil, review_errors.ErrFailedReviewNotFound
	}

	review, err := s.reviewRepository.UpdateReview(req)

	if err != nil {
		s.logger.Error("Failed to update category",
			zap.Error(err),
			zap.Any("request", req))

		return nil, review_errors.ErrFailedUpdateReview
	}

	return s.mapping.ToReviewResponse(review), nil
}

func (s *reviewService) TrashedReview(reviewID int) (*response.ReviewResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing review", zap.Int("reviewID", reviewID))

	review, err := s.reviewRepository.TrashReview(reviewID)

	if err != nil {
		s.logger.Error("Failed to move category to trash",
			zap.Error(err),
			zap.Int("reviewID", reviewID))

		return nil, review_errors.ErrFailedTrashedReview
	}

	return s.mapping.ToReviewResponseDeleteAt(review), nil
}

func (s *reviewService) RestoreReview(reviewID int) (*response.ReviewResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring review", zap.Int("reviewID", reviewID))

	review, err := s.reviewRepository.RestoreReview(reviewID)

	if err != nil {
		s.logger.Error("Failed to restore review from trash",
			zap.Error(err),
			zap.Int("reviewID", reviewID))

		return nil, review_errors.ErrFailedRestoreReview
	}

	return s.mapping.ToReviewResponseDeleteAt(review), nil
}

func (s *reviewService) DeleteReviewPermanent(reviewID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting review", zap.Int("reviewID", reviewID))

	success, err := s.reviewRepository.DeleteReviewPermanently(reviewID)

	if err != nil {
		s.logger.Error("Failed to permanently delete review",
			zap.Error(err),
			zap.Int("reviewID", reviewID))

		return false, review_errors.ErrFailedDeletePermanentReview
	}
	return success, nil
}

func (s *reviewService) RestoreAllReviews() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed reviews")

	success, err := s.reviewRepository.RestoreAllReview()
	if err != nil {
		s.logger.Error("Failed to restore all trashed reviews",
			zap.Error(err))

		return false, review_errors.ErrFailedRestoreAllReviews
	}

	return success, nil
}

func (s *reviewService) DeleteAllReviewsPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all reviews")

	success, err := s.reviewRepository.DeleteAllPermanentReview()
	if err != nil {
		s.logger.Error("Failed to permanently delete all reviews", zap.Error(err))

		return false, review_errors.ErrFailedDeleteAllPermanentReviews
	}

	return success, nil
}
