package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
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

func (s *reviewService) FindAllReviews(search string, page, pageSize int) ([]*response.ReviewResponse, int, *response.ErrorResponse) {
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

	Reviews, totalRecords, err := s.reviewRepository.FindAllReview(search, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch Reviews",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch Reviews"}
	}

	s.logger.Debug("Successfully fetched Reviews",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewsResponse(Reviews), totalRecords, nil
}

func (s *reviewService) FindByActive(search string, page, pageSize int) ([]*response.ReviewResponseDeleteAt, int, *response.ErrorResponse) {
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

	Reviews, totalRecords, err := s.reviewRepository.FindByActive(search, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch Reviews",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch active Reviews"}
	}

	s.logger.Debug("Successfully fetched Reviews",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewsResponseDeleteAt(Reviews), totalRecords, nil
}

func (s *reviewService) FindByTrashed(search string, page, pageSize int) ([]*response.ReviewResponseDeleteAt, int, *response.ErrorResponse) {
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

	Reviews, totalRecords, err := s.reviewRepository.FindByTrashed(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch Reviews",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch trashed Reviews"}
	}

	s.logger.Debug("Successfully fetched Reviews",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewsResponseDeleteAt(Reviews), totalRecords, nil
}

func (s *reviewService) FindByProduct(product_id int, search string, page, pageSize int) ([]*response.ReviewResponse, int, *response.ErrorResponse) {
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

	reviews, totalRecords, err := s.reviewRepository.FindByProduct(product_id, search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch Reviews",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch trashed Reviews"}
	}

	s.logger.Debug("Successfully fetched Reviews",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewsResponse(reviews), totalRecords, nil
}

func (s *reviewService) CreateReview(req *requests.CreateReviewRequest) (*response.ReviewResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new cashier")

	_, err := s.userRepository.FindById(req.UserID)

	if err != nil {
		s.logger.Error("Failed to find user", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to find user"}
	}

	_, err = s.productRepository.FindById(req.ProductID)

	if err != nil {
		s.logger.Error("Failed to find product", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to find product"}
	}

	review, err := s.reviewRepository.CreateReview(req)
	if err != nil {
		s.logger.Error("Failed to create review", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to create review"}
	}

	return s.mapping.ToReviewResponse(review), nil
}

func (s *reviewService) UpdateReview(req *requests.UpdateReviewRequest) (*response.ReviewResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating review", zap.Int("review_id", req.ReviewID))

	_, err := s.reviewRepository.FindById(req.ReviewID)

	if err != nil {
		s.logger.Error("Failed to update review", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to update review"}
	}

	review, err := s.reviewRepository.UpdateReview(req)
	if err != nil {
		s.logger.Error("Failed to update review", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to update review"}
	}

	return s.mapping.ToReviewResponse(review), nil
}

func (s *reviewService) TrashedReview(reviewID int) (*response.ReviewResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing review", zap.Int("reviewID", reviewID))

	review, err := s.reviewRepository.TrashReview(reviewID)
	if err != nil {
		s.logger.Error("Failed to trash review", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to trash review"}
	}

	return s.mapping.ToReviewResponseDeleteAt(review), nil
}

func (s *reviewService) RestoreReview(reviewID int) (*response.ReviewResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring review", zap.Int("reviewID", reviewID))

	review, err := s.reviewRepository.RestoreReview(reviewID)
	if err != nil {
		s.logger.Error("Failed to restore review", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to restore review"}
	}

	return s.mapping.ToReviewResponseDeleteAt(review), nil
}

func (s *reviewService) DeleteReviewPermanent(reviewID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting review", zap.Int("reviewID", reviewID))

	success, err := s.reviewRepository.DeleteReviewPermanently(reviewID)
	if err != nil {
		s.logger.Error("Failed to permanently delete review", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete review"}
	}

	return success, nil
}

func (s *reviewService) RestoreAllReviews() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed reviews")

	success, err := s.reviewRepository.RestoreAllReview()
	if err != nil {
		s.logger.Error("Failed to restore all reviews", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to restore all reviews"}
	}

	return success, nil
}

func (s *reviewService) DeleteAllReviewsPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all reviews")

	success, err := s.reviewRepository.DeleteAllPermanentReview()
	if err != nil {
		s.logger.Error("Failed to permanently delete all reviews", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete all reviews"}
	}

	return success, nil
}
