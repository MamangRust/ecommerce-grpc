package api

import (
	slider_cache "ecommerce/internal/cache/api/slider"
	"ecommerce/internal/domain/requests"
	response_api "ecommerce/internal/mapper"
	"ecommerce/internal/pb"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/upload_image"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type sliderHandleApi struct {
	client       pb.SliderServiceClient
	logger       logger.LoggerInterface
	mapping      response_api.SliderResponseMapper
	upload_image upload_image.ImageUploads
	apiHandler   errors.ApiHandler
	cache        slider_cache.SliderMencache
}

func NewHandlerSlider(
	router *echo.Echo,
	client pb.SliderServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.SliderResponseMapper,
	upload_image upload_image.ImageUploads,
	apiHandler errors.ApiHandler,
	cache slider_cache.SliderMencache,
) *sliderHandleApi {
	sliderHandler := &sliderHandleApi{
		client:       client,
		logger:       logger,
		mapping:      mapping,
		upload_image: upload_image,
		apiHandler:   apiHandler,
		cache:        cache,
	}

	routerSlider := router.Group("/api/slider")

	routerSlider.GET(
		"",
		apiHandler.Handle("findAll", sliderHandler.FindAllSlider),
	)
	routerSlider.GET(
		"/active",
		apiHandler.Handle("findByActive", sliderHandler.FindByActive),
	)
	routerSlider.GET(
		"/trashed",
		apiHandler.Handle("findByTrashed", sliderHandler.FindByTrashed),
	)

	routerSlider.POST(
		"/create",
		apiHandler.Handle("create", sliderHandler.Create),
	)
	routerSlider.POST(
		"/update/:id",
		apiHandler.Handle("update", sliderHandler.Update),
	)

	routerSlider.POST(
		"/trashed/:id",
		apiHandler.Handle("trashed", sliderHandler.TrashedSlider),
	)
	routerSlider.POST(
		"/restore/:id",
		apiHandler.Handle("restore", sliderHandler.RestoreSlider),
	)
	routerSlider.DELETE(
		"/permanent/:id",
		apiHandler.Handle("deletePermanent", sliderHandler.DeleteSliderPermanent),
	)

	routerSlider.POST(
		"/restore/all",
		apiHandler.Handle("restoreAll", sliderHandler.RestoreAllSlider),
	)
	routerSlider.POST(
		"/permanent/all",
		apiHandler.Handle("deleteAllPermanent", sliderHandler.DeleteAllSliderPermanent),
	)

	return sliderHandler

}

// @Security Bearer
// @Summary Find all slider
// @Tags Slider
// @Description Retrieve a list of all slider
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationSlider "List of slider"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve slider data"
// @Router /api/slider [get]
func (h *sliderHandleApi) FindAllSlider(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetSliderAllCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllSliderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindAllSlider")
	}

	apiResponse := h.mapping.ToApiResponsePaginationSlider(res)

	h.cache.SetSliderAllCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active slider
// @Tags Slider
// @Description Retrieve a list of active slider
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationSliderDeleteAt "List of active slider"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve slider data"
// @Router /api/slider/active [get]
func (h *sliderHandleApi) FindByActive(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetSliderActiveCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllSliderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByActive")
	}

	apiResponse := h.mapping.ToApiResponsePaginationSliderDeleteAt(res)

	h.cache.SetSliderActiveCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed slider records.
// @Summary Retrieve trashed slider
// @Tags Slider
// @Description Retrieve a list of trashed slider records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationSliderDeleteAt "List of trashed slider data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve slider data"
// @Router /api/slider/trashed [get]
func (h *sliderHandleApi) FindByTrashed(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetSliderTrashedCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllSliderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByTrashed")
	}

	apiResponse := h.mapping.ToApiResponsePaginationSliderDeleteAt(res)

	h.cache.SetSliderTrashedCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// Create handles the creation of a new slider with image upload.
// @Summary Create a new slider
// @Tags Slider
// @Description Create a new slider with the provided details and an image file
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Slider name"
// @Param image_slider formData file true "Slider image file"
// @Success 200 {object} response.ApiResponseSlider "Successfully created slider"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create slider"
// @Router /api/slider/create [post]
func (h *sliderHandleApi) Create(c echo.Context) error {
	formData, err := h.parseSliderForm(c, true)
	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	ctx := c.Request().Context()

	req := &pb.CreateSliderRequest{
		Name:  formData.Nama,
		Image: formData.FilePath,
	}

	res, err := h.client.Create(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	return c.JSON(http.StatusOK, res)
}

// @Security Bearer
// Update handles the update of an existing slider with image upload.
// @Summary Update an existing slider
// @Tags Slider
// @Description Update an existing slider record with the provided details and an optional image file
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Slider ID"
// @Param name formData string true "Slider name"
// @Param image_slider formData file false "New slider image file"
// @Success 200 {object} response.ApiResponseSlider "Successfully updated slider"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update slider"
// @Router /api/slider/update [post]
func (h *sliderHandleApi) Update(c echo.Context) error {
	sliderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	formData, err := h.parseSliderForm(c, true)
	if err != nil {
		return h.handleGrpcError(err, "Updates")
	}

	ctx := c.Request().Context()

	req := &pb.UpdateSliderRequest{
		Id:    int32(sliderID),
		Name:  formData.Nama,
		Image: formData.FilePath,
	}

	res, err := h.client.Update(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Update")
	}

	return c.JSON(http.StatusOK, res)
}

// @Security Bearer
// TrashedSlider retrieves a trashed slider record by its ID.
// @Summary Retrieve a trashed slider
// @Tags Slider
// @Description Retrieve a trashed slider record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "slider ID"
// @Success 200 {object} response.ApiResponseSliderDeleteAt "Successfully retrieved trashed slider"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed slider"
// @Router /api/slider/trashed/{id} [get]
func (h *sliderHandleApi) TrashedSlider(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdSliderRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedSlider(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Trashed")
	}

	so := h.mapping.ToApiResponseSliderDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreSlider restores a slider record from the trash by its ID.
// @Summary Restore a trashed slider
// @Tags Slider
// @Description Restore a trashed slider record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "slider ID"
// @Success 200 {object} response.ApiResponseSliderDeleteAt "Successfully restored slider"
// @Failure 400 {object} response.ErrorResponse "Invalid slider ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore slider"
// @Router /api/slider/restore/{id} [post]
func (h *sliderHandleApi) RestoreSlider(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdSliderRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreSlider(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Restore")
	}

	so := h.mapping.ToApiResponseSliderDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteSliderPermanent permanently deletes a slider record by its ID.
// @Summary Permanently delete a slider
// @Tags Slider
// @Description Permanently delete a slider record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "slider ID"
// @Success 200 {object} response.ApiResponseSliderDelete "Successfully deleted slider record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete slider:"
// @Router /api/slider/delete/{id} [delete]
func (h *sliderHandleApi) DeleteSliderPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdSliderRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteSliderPermanent(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "DeleteSlider")
	}

	so := h.mapping.ToApiResponseSliderDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllSlider restores a slider record from the trash by its ID.
// @Summary Restore a trashed slider
// @Tags Slider
// @Description Restore a trashed slider record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "slider ID"
// @Success 200 {object} response.ApiResponseSliderAll "Successfully restored slider all"
// @Failure 400 {object} response.ErrorResponse "Invalid slider ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore slider"
// @Router /api/slider/restore/all [post]
func (h *sliderHandleApi) RestoreAllSlider(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllSlider(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "RestoreAll")
	}

	so := h.mapping.ToApiResponseSliderAll(res)

	h.logger.Debug("Successfully restored all slider")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllSliderPermanent permanently deletes a slider record by its ID.
// @Summary Permanently delete a slider
// @Tags Slider
// @Description Permanently delete a slider record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "slider ID"
// @Success 200 {object} response.ApiResponseSliderAll "Successfully deleted slider record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete slider:"
// @Router /api/slider/delete/all [post]
func (h *sliderHandleApi) DeleteAllSliderPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllSliderPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "DeleteAll")
	}

	so := h.mapping.ToApiResponseSliderAll(res)

	h.logger.Debug("Successfully deleted all slider permanently")

	return c.JSON(http.StatusOK, so)
}

func (h *sliderHandleApi) parseSliderForm(c echo.Context, requireImage bool) (requests.SliderFormData, error) {
	var formData requests.SliderFormData

	formData.Nama = strings.TrimSpace(c.FormValue("name"))
	if formData.Nama == "" {
		return formData, errors.NewBadRequestError("name is required")
	}

	file, err := c.FormFile("image_slider")
	if err != nil {
		if requireImage {
			h.logger.Debug("Image upload error", zap.Error(err))
			return formData, errors.NewBadRequestError("image_slider is required")
		}
		return formData, nil
	}

	imagePath, err := h.upload_image.ProcessImageUpload(c, file)
	if err != nil {
		return formData, err
	}

	formData.FilePath = imagePath
	return formData, nil
}

func (h *sliderHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Slider").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Slider already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Slider service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *sliderHandleApi) parseValidationErrors(err error) []errors.ValidationError {
	var validationErrs []errors.ValidationError

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			validationErrs = append(validationErrs, errors.ValidationError{
				Field:   fe.Field(),
				Message: h.getValidationMessage(fe),
			})
		}
		return validationErrs
	}

	return []errors.ValidationError{
		{
			Field:   "general",
			Message: err.Error(),
		},
	}
}

func (h *sliderHandleApi) getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s", fe.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fe.Param())
	default:
		return fmt.Sprintf("Validation failed on '%s' tag", fe.Tag())
	}
}
