package service

import (
	"context"
	slider_cache "ecommerce/internal/cache/slider"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/slider_errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type sliderService struct {
	sliderRepository repository.SliderRepository
	logger           logger.LoggerInterface
	cache            slider_cache.SliderMencache
	observability    observability.TraceLoggerObservability
}

type SliderServiceDeps struct {
	SliderRepository repository.SliderRepository
	Logger           logger.LoggerInterface
	Cache            slider_cache.SliderMencache
	Observability    observability.TraceLoggerObservability
}

func NewSliderService(deps SliderServiceDeps) *sliderService {
	return &sliderService{
		sliderRepository: deps.SliderRepository,
		logger:           deps.Logger,
		cache:            deps.Cache,
		observability:    deps.Observability,
	}
}

func (s *sliderService) FindAllSlider(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersRow, *int, error) {
	const method = "FindAll"

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

	if data, total, found := s.cache.GetSliderAllCache(ctx, req); found {
		logSuccess("Successfully retrieved sliders from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	sliders, err := s.sliderRepository.FindAllSlider(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetSlidersRow](
			s.logger,
			slider_errors.ErrFailedFindAllSliders,
			method,
			span,

			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(sliders) > 0 {
		totalCount = int(sliders[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetSliderAllCache(ctx, req, sliders, &totalCount)

	logSuccess("Successfully fetched sliders from repository",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return sliders, &totalCount, nil
}

func (s *sliderService) FindByActive(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersActiveRow, *int, error) {
	const method = "FindByActive"

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

	if data, total, found := s.cache.GetSliderActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active sliders from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	sliders, err := s.sliderRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetSlidersActiveRow](
			s.logger,
			slider_errors.ErrFailedFindActiveSliders,
			method,
			span,

			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(sliders) > 0 {
		totalCount = int(sliders[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetSliderActiveCache(ctx, req, sliders, &totalCount)

	logSuccess("Successfully fetched active sliders from repository",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return sliders, &totalCount, nil
}

func (s *sliderService) FindById(ctx context.Context, slider_id int) (*db.GetSliderByIDRow, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("slider_id", slider_id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetSliderCache(ctx, slider_id); found {
		logSuccess("Successfully retrieved slider by ID from cache",
			zap.Int("slider_id", slider_id))
		return data, nil
	}

	slider, err := s.sliderRepository.FindById(ctx, slider_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetSliderByIDRow](
			s.logger,
			slider_errors.ErrFailedFindSliderByID,
			method,
			span,
			zap.Int("slider_id", slider_id),
		)
	}

	s.cache.SetSliderCache(ctx, slider)

	logSuccess("Successfully fetched slider by ID from repository",
		zap.Int("slider_id", slider_id))

	return slider, nil
}

func (s *sliderService) FindByTrashed(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersTrashedRow, *int, error) {
	const method = "FindByTrashed"

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

	if data, total, found := s.cache.GetSliderTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved trashed sliders from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	sliders, err := s.sliderRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetSlidersTrashedRow](
			s.logger,
			slider_errors.ErrFailedFindTrashedSliders,
			method,
			span,

			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(sliders) > 0 {
		totalCount = int(sliders[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetSliderTrashedCache(ctx, req, sliders, &totalCount)

	logSuccess("Successfully fetched trashed sliders from repository",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return sliders, &totalCount, nil
}

func (s *sliderService) CreateSlider(ctx context.Context, req *requests.CreateSliderRequest) (*db.CreateSliderRow, error) {
	const method = "CreateSlider"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("slider", req.Nama))

	defer func() {
		end(status)
	}()

	slider, err := s.sliderRepository.CreateSlider(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateSliderRow](
			s.logger,
			slider_errors.ErrFailedCreateSlider,
			method,
			span,
			zap.String("slider", req.Nama),
		)
	}

	s.cache.InvalidateSliderCache(ctx)

	logSuccess("Successfully created slider",
		zap.Int("slider_id", int(slider.SliderID)),
		zap.String("slider_name", slider.Name))

	return slider, nil
}

func (s *sliderService) UpdateSlider(ctx context.Context, req *requests.UpdateSliderRequest) (*db.UpdateSliderRow, error) {
	const method = "UpdateSlider"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("slider_id", *req.ID),
		attribute.String("new_name", req.Nama))

	defer func() {
		end(status)
	}()

	slider, err := s.sliderRepository.UpdateSlider(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateSliderRow](
			s.logger,
			slider_errors.ErrFailedUpdateSlider,
			method,
			span,
			zap.Int("slider_id", *req.ID),
			zap.String("new_name", req.Nama),
		)
	}

	s.cache.InvalidateSliderCache(ctx)

	logSuccess("Successfully updated slider",
		zap.Int("slider_id", int(slider.SliderID)),
		zap.String("slider_name", slider.Name))

	return slider, nil
}

func (s *sliderService) TrashSlider(ctx context.Context, slider_id int) (*db.Slider, error) {
	const method = "TrashSlider"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("slider_id", slider_id))

	defer func() {
		end(status)
	}()

	slider, err := s.sliderRepository.TrashSlider(ctx, slider_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Slider](
			s.logger,
			slider_errors.ErrFailedTrashSlider,
			method,
			span,
			zap.Int("slider_id", slider_id),
		)
	}

	s.cache.InvalidateSliderCache(ctx)

	logSuccess("Successfully trashed slider",
		zap.Int("slider_id", int(slider.SliderID)))

	return slider, nil
}

func (s *sliderService) RestoreSlider(ctx context.Context, sliderID int) (*db.Slider, error) {
	const method = "RestoreSlider"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("sliderID", sliderID))

	defer func() {
		end(status)
	}()

	slider, err := s.sliderRepository.RestoreSlider(ctx, sliderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Slider](
			s.logger,
			slider_errors.ErrFailedRestoreSlider,
			method,
			span,
			zap.Int("sliderID", sliderID),
		)
	}

	s.cache.InvalidateSliderCache(ctx)

	logSuccess("Successfully restored slider",
		zap.Int("slider_id", int(slider.SliderID)))

	return slider, nil
}

func (s *sliderService) DeleteSliderPermanently(ctx context.Context, sliderID int) (bool, error) {
	const method = "DeleteSliderPermanently"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("sliderID", sliderID))

	defer func() {
		end(status)
	}()

	success, err := s.sliderRepository.DeleteSliderPermanently(ctx, sliderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			slider_errors.ErrFailedDeletePermanentSlider,
			method,
			span,
			zap.Int("sliderID", sliderID),
		)
	}

	s.cache.InvalidateSliderCache(ctx)

	logSuccess("Successfully permanently deleted slider",
		zap.Int("sliderID", sliderID))

	return success, nil
}

func (s *sliderService) RestoreAllSliders(ctx context.Context) (bool, error) {
	const method = "RestoreAllSliders"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.sliderRepository.RestoreAllSlider(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			slider_errors.ErrFailedRestoreAllSliders,
			method,
			span,
		)
	}

	s.cache.InvalidateSliderCache(ctx)
	logSuccess("Successfully restored all trashed sliders")

	return success, nil
}

func (s *sliderService) DeleteAllPermanentSlider(ctx context.Context) (bool, error) {
	const method = "DeleteAllSlidersPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.sliderRepository.DeleteAllPermanentSlider(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			slider_errors.ErrFailedDeleteAllPermanentSliders,
			method,
			span,
		)
	}

	s.cache.InvalidateSliderCache(ctx)
	logSuccess("Successfully permanently deleted all trashed sliders")

	return success, nil
}
