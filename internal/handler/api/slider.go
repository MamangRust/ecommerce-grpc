package api

import (
	"ecommerce/internal/domain/response"
	response_api "ecommerce/internal/mapper/response/api"
	"ecommerce/internal/pb"
	"ecommerce/pkg/logger"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type sliderHandleApi struct {
	client  pb.SliderServiceClient
	logger  logger.LoggerInterface
	mapping response_api.SliderResponseMapper
}

func NewHandlerSlider(
	router *echo.Echo,
	client pb.SliderServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.SliderResponseMapper,
) *sliderHandleApi {
	sliderHandler := &sliderHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
	}

	routerSlider := router.Group("/api/slider")

	routerSlider.GET("", sliderHandler.FindAllSlider)
	routerSlider.GET("/active", sliderHandler.FindByActive)
	routerSlider.GET("/trashed", sliderHandler.FindByTrashed)

	routerSlider.POST("/create", sliderHandler.Create)
	routerSlider.POST("/update/:id", sliderHandler.Update)

	routerSlider.POST("/trashed/:id", sliderHandler.TrashedSlider)
	routerSlider.POST("/restore/:id", sliderHandler.RestoreSlider)
	routerSlider.DELETE("/permanent/:id", sliderHandler.DeleteSliderPermanent)

	routerSlider.POST("/restore/all", sliderHandler.RestoreAllSlider)
	routerSlider.POST("/permanent/all", sliderHandler.DeleteAllSliderPermanent)

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

	req := &pb.FindAllSliderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve slider data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve slider data: ",
		})
	}

	so := h.mapping.ToApiResponsePaginationSlider(res)

	return c.JSON(http.StatusOK, so)
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

	req := &pb.FindAllSliderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve slider data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve slider data: ",
		})
	}

	so := h.mapping.ToApiResponsePaginationSliderDeleteAt(res)

	return c.JSON(http.StatusOK, so)
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

	req := &pb.FindAllSliderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve slider data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve slider data: ",
		})
	}

	so := h.mapping.ToApiResponsePaginationSliderDeleteAt(res)

	return c.JSON(http.StatusOK, so)
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
	name := c.FormValue("name")

	file, err := c.FormFile("image_slider")
	if err != nil {
		h.logger.Debug("Invalid image file", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid image file",
		})
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	imagePath := "uploads/slider/" + file.Filename
	dst, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	ctx := c.Request().Context()

	req := &pb.CreateSliderRequest{
		Name:  name,
		Image: imagePath,
	}

	res, err := h.client.Create(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to create slider", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create slider: " + err.Error(),
		})
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
// @Param slider_id formData int true "Slider ID"
// @Param name formData string true "Slider name"
// @Param image_slider formData file false "New slider image file"
// @Success 200 {object} response.ApiResponseSlider "Successfully updated slider"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update slider"
// @Router /api/slider/update [post]
func (h *sliderHandleApi) Update(c echo.Context) error {
	sliderID, err := strconv.Atoi(c.FormValue("slider_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid slider ID",
		})
	}

	name := c.FormValue("name")

	imagePath := ""
	file, err := c.FormFile("image_slider")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		imagePath = "uploads/slider/" + file.Filename
		dst, err := os.Create(imagePath)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
	}

	ctx := c.Request().Context()

	req := &pb.UpdateSliderRequest{
		Id:    int32(sliderID),
		Name:  name,
		Image: imagePath,
	}

	res, err := h.client.Update(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to update slider", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update slider: " + err.Error(),
		})
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
		h.logger.Debug("Invalid slider ID", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid slider ID",
		})
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdSliderRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedSlider(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to trashed slider", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to trashed slider: ",
		})
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
		h.logger.Debug("Invalid slider ID", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid slider ID",
		})
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdSliderRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreSlider(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to restore slider", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore slider: ",
		})
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
		h.logger.Debug("Invalid slider ID", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid slider ID",
		})
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdSliderRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteSliderPermanent(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to delete slider", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete slider: ",
		})
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
		h.logger.Error("Failed to restore all slider", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently restore all slider",
		})
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
		h.logger.Error("Failed to permanently delete all slider", zap.Error(err))

		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all slider",
		})
	}

	so := h.mapping.ToApiResponseSliderAll(res)

	h.logger.Debug("Successfully deleted all slider permanently")

	return c.JSON(http.StatusOK, so)
}
