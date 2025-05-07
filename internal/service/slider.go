package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/errors/slider_errors"
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

func (s *sliderService) FindAll(req *requests.FindAllSlider) ([]*response.SliderResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

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

	sliders, totalRecords, err := s.sliderRepository.FindAllSlider(req)

	if err != nil {
		s.logger.Error("Failed to retrieve sliders",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search))

		return nil, nil, slider_errors.ErrFailedFindAllSliders
	}

	slidersResponse := s.mapping.ToSlidersResponse(sliders)

	s.logger.Debug("Successfully fetched sliders",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return slidersResponse, totalRecords, nil
}

func (s *sliderService) FindByActive(req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

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

	sliders, totalRecords, err := s.sliderRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active sliders",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search))

		return nil, nil, slider_errors.ErrFailedFindActiveSliders
	}

	s.logger.Debug("Successfully fetched sliders",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToSlidersResponseDeleteAt(sliders), totalRecords, nil
}

func (s *sliderService) FindByTrashed(req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

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

	sliders, totalRecords, err := s.sliderRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed slider",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search))

		return nil, nil, slider_errors.ErrFailedFindTrashedSliders
	}

	s.logger.Debug("Successfully fetched slider",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToSlidersResponseDeleteAt(sliders), totalRecords, nil
}

func (s *sliderService) CreateSlider(req *requests.CreateSliderRequest) (*response.SliderResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new slider")

	slider, err := s.sliderRepository.CreateSlider(req)

	if err != nil {
		s.logger.Error("Failed to create new slider record",
			zap.String("slider", req.Nama),
			zap.Error(err))

		return nil, slider_errors.ErrFailedCreateSlider
	}

	return s.mapping.ToSliderResponse(slider), nil
}

func (s *sliderService) UpdateSlider(req *requests.UpdateSliderRequest) (*response.SliderResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating slider", zap.Int("slider_id", *req.ID))

	slider, err := s.sliderRepository.UpdateSlider(req)

	if err != nil {
		s.logger.Error("Failed to update slider record",
			zap.Int("role_id", *req.ID),
			zap.String("new_name", req.Nama),
			zap.Error(err))

		return nil, slider_errors.ErrFailedUpdateSlider
	}

	return s.mapping.ToSliderResponse(slider), nil
}

func (s *sliderService) TrashedSlider(slider_id int) (*response.SliderResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing slider", zap.Int("slider", slider_id))

	slider, err := s.sliderRepository.TrashSlider(slider_id)

	if err != nil {
		s.logger.Error("Failed to move slider to trash",
			zap.Int("slider_id", slider_id),
			zap.Error(err))

		return nil, slider_errors.ErrFailedTrashSlider
	}

	return s.mapping.ToSliderResponseDeleteAt(slider), nil
}

func (s *sliderService) RestoreSlider(sliderID int) (*response.SliderResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring slider", zap.Int("sliderID", sliderID))

	slider, err := s.sliderRepository.RestoreSlider(sliderID)

	if err != nil {
		s.logger.Error("Failed to restore slider from trash",
			zap.Int("sliderID", sliderID),
			zap.Error(err))

		return nil, slider_errors.ErrFailedRestoreSlider
	}

	return s.mapping.ToSliderResponseDeleteAt(slider), nil
}

func (s *sliderService) DeleteSliderPermanent(sliderID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting slider", zap.Int("sliderID", sliderID))

	success, err := s.sliderRepository.DeleteSliderPermanently(sliderID)

	if err != nil {
		s.logger.Error("Failed to permanently delete sliders",
			zap.Int("sliderID", sliderID),
			zap.Error(err))

		return false, slider_errors.ErrFailedDeletePermanentSlider
	}

	return success, nil
}

func (s *sliderService) RestoreAllSliders() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed sliders")

	success, err := s.sliderRepository.RestoreAllSlider()

	if err != nil {
		s.logger.Error("Failed to restore all trashed sliders",
			zap.Error(err))

		return false, slider_errors.ErrFailedRestoreAllSliders
	}

	return success, nil
}

func (s *sliderService) DeleteAllSlidersPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all sliders")

	success, err := s.sliderRepository.DeleteAllPermanentSlider()

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed sliders",
			zap.Error(err))

		return false, slider_errors.ErrFailedDeleteAllPermanentSliders
	}

	return success, nil
}
