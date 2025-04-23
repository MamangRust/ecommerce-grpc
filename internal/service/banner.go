package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"
	"net/http"

	"go.uber.org/zap"
)

type bannerService struct {
	bannerRepository repository.BannerRepository
	logger           logger.LoggerInterface
	mapping          response_service.BannerResponseMapper
}

func NewBannerService(
	bannerRepository repository.BannerRepository,
	logger logger.LoggerInterface,
	mapping response_service.BannerResponseMapper,
) *bannerService {
	return &bannerService{
		bannerRepository: bannerRepository,
		logger:           logger,
		mapping:          mapping,
	}
}

func (s *bannerService) FindAll(req *requests.FindAllBanner) ([]*response.BannerResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all Banners",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	Banners, totalRecords, err := s.bannerRepository.FindAllBanners(req)

	if err != nil {
		s.logger.Error("Failed to fetch Banners",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve Banners list",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched Banners",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToBannersResponse(Banners), totalRecords, nil
}

func (s *bannerService) FindByActive(req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all Banners active",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	Banners, totalRecords, err := s.bannerRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active Banners",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve active Banner",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched active Banner",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToBannersResponseDeleteAt(Banners), totalRecords, nil
}

func (s *bannerService) FindByTrashed(req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all Banners trashed",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	Banners, totalRecords, err := s.bannerRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed Banners",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve trashed Banner",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched trashed Banner",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToBannersResponseDeleteAt(Banners), totalRecords, nil
}

func (s *bannerService) FindById(BannerID int) (*response.BannerResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching Banner by ID", zap.Int("BannerID", BannerID))

	Banner, err := s.bannerRepository.FindById(BannerID)

	if err != nil {
		s.logger.Error("Failed to retrieve Banner details",
			zap.Error(err),
			zap.Int("Banner_id", BannerID))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve Banner details",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToBannerResponse(Banner), nil
}

func (s *bannerService) CreateBanner(req *requests.CreateBannerRequest) (*response.BannerResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new Banner")

	Banner, err := s.bannerRepository.CreateBanner(req)

	if err != nil {
		s.logger.Error("Failed to create new Banner",
			zap.Error(err),
			zap.Any("request", req))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create new Banner record",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToBannerResponse(Banner), nil
}

func (s *bannerService) UpdateBanner(req *requests.UpdateBannerRequest) (*response.BannerResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating Banner", zap.Int("BannerID", *req.BannerID))

	Banner, err := s.bannerRepository.UpdateBanner(req)

	if err != nil {
		s.logger.Error("Failed to update Banner",
			zap.Error(err),
			zap.Any("request", req))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update Banner record",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToBannerResponse(Banner), nil
}

func (s *bannerService) TrashedBanner(BannerID int) (*response.BannerResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing Banner", zap.Int("BannerID", BannerID))

	Banner, err := s.bannerRepository.TrashedBanner(BannerID)

	if err != nil {
		s.logger.Error("Failed to move Banner to trash",
			zap.Error(err),
			zap.Int("Banner_id", BannerID))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to move Banner to trash",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToBannerResponseDeleteAt(Banner), nil
}

func (s *bannerService) RestoreBanner(BannerID int) (*response.BannerResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring Banner", zap.Int("BannerID", BannerID))

	Banner, err := s.bannerRepository.RestoreBanner(BannerID)

	if err != nil {
		s.logger.Error("Failed to restore Banner from trash",
			zap.Error(err),
			zap.Int("Banner_id", BannerID))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore Banner from trash",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToBannerResponseDeleteAt(Banner), nil
}

func (s *bannerService) DeleteBannerPermanent(BannerID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Deleting Banner permanently", zap.Int("BannerID", BannerID))

	success, err := s.bannerRepository.DeleteBannerPermanent(BannerID)

	if err != nil {
		s.logger.Error("Failed to permanently delete Banner",
			zap.Error(err),
			zap.Int("Banner_id", BannerID))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete Banner",
			Code:    http.StatusInternalServerError,
		}
	}

	return success, nil
}

func (s *bannerService) RestoreAllBanner() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed Banners")

	success, err := s.bannerRepository.RestoreAllBanner()

	if err != nil {
		s.logger.Error("Failed to restore all trashed Banners",
			zap.Error(err))

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all trashed Banners",
			Code:    http.StatusInternalServerError,
		}
	}

	return success, nil
}

func (s *bannerService) DeleteAllBannerPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all Banners")

	success, err := s.bannerRepository.DeleteAllBannerPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed Banners",
			zap.Error(err))

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all trashed Banners",
			Code:    http.StatusInternalServerError,
		}
	}

	return success, nil
}
