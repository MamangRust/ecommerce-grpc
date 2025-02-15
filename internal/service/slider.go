package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"

	"go.uber.org/zap"
)

type sliderService struct {
	sliderRepository repository.SliderRepository
	logger           logger.LoggerInterface
	mapping          response_service.SliderResponseMapper
}

func NewSliderService(
	sliderRepository repository.SliderRepository,
	logger logger.LoggerInterface,
	mapping response_service.SliderResponseMapper,
) *sliderService {
	return &sliderService{
		sliderRepository: sliderRepository,
		logger:           logger,
		mapping:          mapping,
	}
}

func (s *sliderService) FindAll(page int, pageSize int, search string) ([]*response.SliderResponse, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching sliders",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	sliders, totalRecords, err := s.sliderRepository.FindAllSlider(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch sliders",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch sliders",
		}
	}

	slidersResponse := s.mapping.ToSlidersResponse(sliders)

	s.logger.Debug("Successfully fetched sliders",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return slidersResponse, int(totalRecords), nil
}

func (s *sliderService) FindByActive(search string, page, pageSize int) ([]*response.SliderResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching sliders",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	sliders, totalRecords, err := s.sliderRepository.FindByActive(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch sliders",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch active sliders"}
	}

	s.logger.Debug("Successfully fetched sliders",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToSlidersResponseDeleteAt(sliders), totalRecords, nil
}

func (s *sliderService) FindByTrashed(search string, page, pageSize int) ([]*response.SliderResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching sliders",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	sliders, totalRecords, err := s.sliderRepository.FindByTrashed(search, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch sliders",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch trashed sliders"}
	}

	s.logger.Debug("Successfully fetched sliders",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToSlidersResponseDeleteAt(sliders), totalRecords, nil
}

func (s *sliderService) CreateSlider(req *requests.CreateSliderRequest) (*response.SliderResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new slider")

	slider, err := s.sliderRepository.CreateSlider(req)
	if err != nil {
		s.logger.Error("Failed to create slider", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to create slider"}
	}

	return s.mapping.ToSliderResponse(slider), nil
}

func (s *sliderService) UpdateSlider(req *requests.UpdateSliderRequest) (*response.SliderResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating slider", zap.Int("slider_id", req.ID))

	slider, err := s.sliderRepository.UpdateSlider(req)
	if err != nil {
		s.logger.Error("Failed to update slider", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to update slider"}
	}

	return s.mapping.ToSliderResponse(slider), nil
}

func (s *sliderService) TrashedSlider(slider_id int) (*response.SliderResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing slider", zap.Int("slider", slider_id))

	slider, err := s.sliderRepository.TrashSlider(slider_id)
	if err != nil {
		s.logger.Error("Failed to trash slider", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to trash slider"}
	}

	return s.mapping.ToSliderResponseDeleteAt(slider), nil
}

func (s *sliderService) RestoreSlider(sliderID int) (*response.SliderResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring slider", zap.Int("sliderID", sliderID))

	slider, err := s.sliderRepository.RestoreSlider(sliderID)
	if err != nil {
		s.logger.Error("Failed to restore slider", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to restore slider"}
	}

	return s.mapping.ToSliderResponseDeleteAt(slider), nil
}

func (s *sliderService) DeleteSliderPermanent(sliderID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting slider", zap.Int("sliderID", sliderID))

	success, err := s.sliderRepository.DeleteSliderPermanently(sliderID)
	if err != nil {
		s.logger.Error("Failed to permanently delete slider", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete slider"}
	}

	return success, nil
}

func (s *sliderService) RestoreAllSliders() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed sliders")

	success, err := s.sliderRepository.RestoreAllSlider()
	if err != nil {
		s.logger.Error("Failed to restore all sliders", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to restore all sliders"}
	}

	return success, nil
}

func (s *sliderService) DeleteAllSlidersPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all sliders")

	success, err := s.sliderRepository.DeleteAllPermanentSlider()
	if err != nil {
		s.logger.Error("Failed to permanently delete all sliders", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete all sliders"}
	}

	return success, nil
}
