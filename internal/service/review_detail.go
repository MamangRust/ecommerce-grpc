package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	reviewdetail_errors "ecommerce/pkg/errors/review_detail"
	"ecommerce/pkg/logger"
	"os"

	"go.uber.org/zap"
)

type reviewDetailService struct {
	reviewDetailRepository repository.ReviewDetailRepository
	logger                 logger.LoggerInterface
	mapping                response_service.ReviewDetailResponeMapper
}

func NewReviewDetailService(
	reviewDetailRepository repository.ReviewDetailRepository,
	logger logger.LoggerInterface,
	mapping response_service.ReviewDetailResponeMapper,
) *reviewDetailService {
	return &reviewDetailService{
		reviewDetailRepository: reviewDetailRepository,
		logger:                 logger,
		mapping:                mapping,
	}
}

func (s *reviewDetailService) FindAll(req *requests.FindAllReview) ([]*response.ReviewDetailsResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all Review Details",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.reviewDetailRepository.FindAllReviews(req)

	if err != nil {
		s.logger.Error("Failed to fetch Review Details",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		return nil, nil, reviewdetail_errors.ErrFailedFindAllReview
	}

	s.logger.Debug("Successfully fetched Review Details",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToReviewsDetailsResponse(res), totalRecords, nil
}

func (s *reviewDetailService) FindByActive(req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all Review Details active",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.reviewDetailRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active Review Details",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, reviewdetail_errors.ErrFailedFindActiveReview
	}

	s.logger.Debug("Successfully fetched active Review Detail",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewDetailsResponseDeleteAt(res), totalRecords, nil
}

func (s *reviewDetailService) FindByTrashed(req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all Review Details trashed",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.reviewDetailRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed Review Details",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, reviewdetail_errors.ErrFailedFindTrashedReview
	}

	s.logger.Debug("Successfully fetched trashed Review Detail",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewDetailsResponseDeleteAt(res), totalRecords, nil
}

func (s *reviewDetailService) FindById(review_id int) (*response.ReviewDetailsResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching Review Detail by ID", zap.Int("Review DetailID", review_id))

	res, err := s.reviewDetailRepository.FindById(review_id)

	if err != nil {
		s.logger.Error("Failed to retrieve Review Detail details",
			zap.Error(err),
			zap.Int("review_id", review_id))

		return nil, reviewdetail_errors.ErrReviewDetailNotFoundRes
	}

	return s.mapping.ToReviewDetailsResponse(res), nil
}

func (s *reviewDetailService) CreateReviewDetail(req *requests.CreateReviewDetailRequest) (*response.ReviewDetailsResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new Review Detail")

	res, err := s.reviewDetailRepository.CreateReviewDetail(req)

	if err != nil {
		s.logger.Error("Failed to create new Review Detail",
			zap.Error(err),
			zap.Any("request", req))

		return nil, reviewdetail_errors.ErrFailedCreateReviewDetail
	}

	return s.mapping.ToReviewDetailsResponse(res), nil
}

func (s *reviewDetailService) UpdateReviewDetail(req *requests.UpdateReviewDetailRequest) (*response.ReviewDetailsResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating Review Detail", zap.Int("Review DetailID", *req.ReviewDetailID))

	res, err := s.reviewDetailRepository.UpdateReviewDetail(req)

	if err != nil {
		s.logger.Error("Failed to update Review Detail",
			zap.Error(err),
			zap.Any("request", req))

		return nil, reviewdetail_errors.ErrFailedUpdateReviewDetail
	}

	return s.mapping.ToReviewDetailsResponse(res), nil
}

func (s *reviewDetailService) TrashedReviewDetail(review_id int) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing Review Detail", zap.Int("Review DetailID", review_id))

	res, err := s.reviewDetailRepository.TrashedReviewDetail(review_id)

	if err != nil {
		s.logger.Error("Failed to move Review Detail to trash",
			zap.Error(err),
			zap.Int("Review Detail_id", review_id))

		return nil, reviewdetail_errors.ErrFailedTrashedReviewDetail
	}

	return s.mapping.ToReviewDetailResponseDeleteAt(res), nil
}

func (s *reviewDetailService) RestoreReviewDetail(review_id int) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring Review Detail", zap.Int("Review DetailID", review_id))

	res, err := s.reviewDetailRepository.RestoreReviewDetail(review_id)

	if err != nil {
		s.logger.Error("Failed to restore Review Detail from trash",
			zap.Error(err),
			zap.Int("Review Detail_id", review_id))

		return nil, reviewdetail_errors.ErrFailedRestoreReviewDetail
	}

	return s.mapping.ToReviewDetailResponseDeleteAt(res), nil
}

func (s *reviewDetailService) DeleteReviewDetailPermanent(review_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Deleting Review Detail permanently", zap.Int("Review DetailID", review_id))

	res, err := s.reviewDetailRepository.FindByIdTrashed(review_id)

	if err != nil {
		s.logger.Error("Failed to find review detail",
			zap.Int("review_id", review_id),
			zap.Error(err))

		return false, reviewdetail_errors.ErrFailedDeletePermanentReview
	}

	if res.Url != "" {
		err := os.Remove(res.Url)
		if err != nil {
			if os.IsNotExist(err) {
				s.logger.Debug("review detail upload path file not found, continuing with review detail deletion",
					zap.String("upload path", res.Url))

				return false, reviewdetail_errors.ErrFailedImageNotFound
			} else {
				s.logger.Debug("Failed to delete review detail upload path",
					zap.String("upload path", res.Url),
					zap.Error(err))

				return false, reviewdetail_errors.ErrFailedRemoveImage
			}
		} else {
			s.logger.Debug("Successfully deleted review detail upload path",
				zap.String("upload path", res.Url))
		}
	}

	success, err := s.reviewDetailRepository.DeleteReviewDetailPermanent(review_id)

	if err != nil {
		s.logger.Error("Failed to permanently delete Review Detail",
			zap.Error(err),
			zap.Int("Review Detail_id", review_id))

		return false, reviewdetail_errors.ErrFailedDeletePermanentReview
	}

	return success, nil
}

func (s *reviewDetailService) RestoreAllReviewDetail() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed Review Details")

	success, err := s.reviewDetailRepository.RestoreAllReviewDetail()

	if err != nil {
		s.logger.Error("Failed to restore all trashed Review Details",
			zap.Error(err))

		return false, reviewdetail_errors.ErrFailedRestoreAllReviewDetail
	}

	return success, nil
}

func (s *reviewDetailService) DeleteAllReviewDetailPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all Review Details")

	success, err := s.reviewDetailRepository.DeleteAllReviewDetailPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed Review Details",
			zap.Error(err))

		return false, reviewdetail_errors.ErrFailedDeleteAllReviewDetail
	}

	return success, nil
}
