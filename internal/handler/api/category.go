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

type categoryHandleApi struct {
	client  pb.CategoryServiceClient
	logger  logger.LoggerInterface
	mapping response_api.CategoryResponseMapper
}

func NewHandlerCategory(
	router *echo.Echo,
	client pb.CategoryServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.CategoryResponseMapper,
) *categoryHandleApi {
	categoryHandler := &categoryHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
	}

	routercategory := router.Group("/api/category")

	routercategory.GET("", categoryHandler.FindAllCategory)
	routercategory.GET("/:id", categoryHandler.FindById)
	routercategory.GET("/active", categoryHandler.FindByActive)
	routercategory.GET("/trashed", categoryHandler.FindByTrashed)

	routercategory.POST("/create", categoryHandler.Create)
	routercategory.POST("/update/:id", categoryHandler.Update)

	routercategory.POST("/trashed/:id", categoryHandler.TrashedCategory)
	routercategory.POST("/restore/:id", categoryHandler.RestoreCategory)
	routercategory.DELETE("/permanent/:id", categoryHandler.DeleteCategoryPermanent)

	routercategory.POST("/restore/all", categoryHandler.RestoreAllCategory)
	routercategory.POST("/permanent/all", categoryHandler.DeleteAllCategoryPermanent)

	return categoryHandler
}

// @Security Bearer
// @Summary Find all category
// @Tags Category
// @Description Retrieve a list of all category
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationCategory "List of category"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve category data"
// @Router /api/category [get]
func (h *categoryHandleApi) FindAllCategory(c echo.Context) error {
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

	req := &pb.FindAllCategoryRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve category data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve category data: ",
		})
	}

	so := h.mapping.ToApiResponsePaginationCategory(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find category by ID
// @Tags Category
// @Description Retrieve a category by ID
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategory "Category data"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve category data"
// @Router /api/category/{id} [get]
func (h *categoryHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid category ID", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid category ID",
		})
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCategoryRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve category data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve category data: ",
		})
	}

	so := h.mapping.ToApiResponseCategory(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active category
// @Tags Category
// @Description Retrieve a list of active category
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationCategoryDeleteAt "List of active category"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve category data"
// @Router /api/category/active [get]
func (h *categoryHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllCategoryRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve category data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve category data: ",
		})
	}

	so := h.mapping.ToApiResponsePaginationCategoryDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed category records.
// @Summary Retrieve trashed category
// @Tags Category
// @Description Retrieve a list of trashed category records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationCategoryDeleteAt "List of trashed category data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve category data"
// @Router /api/category/trashed [get]
func (h *categoryHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllCategoryRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve category data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve category data: ",
		})
	}

	so := h.mapping.ToApiResponsePaginationCategoryDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Create handles the creation of a new category with image upload.
// @Summary Create a new category
// @Tags Category
// @Description Create a new category with the provided details and an image file
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Category name"
// @Param description formData string true "Category description"
// @Param slug_category formData string true "Category slug"
// @Param image_category formData file true "Category image file"
// @Success 200 {object} response.ApiResponseCategory "Successfully created category"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create category"
// @Router /api/category/create [post]
func (h *categoryHandleApi) Create(c echo.Context) error {
	name := c.FormValue("name")
	description := c.FormValue("description")
	slugCategory := c.FormValue("slug_category")

	file, err := c.FormFile("image_category")
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

	imagePath := "uploads/category/" + file.Filename
	dst, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	ctx := c.Request().Context()

	req := &pb.CreateCategoryRequest{
		Name:          name,
		Description:   description,
		SlugCategory:  slugCategory,
		ImageCategory: imagePath,
	}

	res, err := h.client.Create(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to create category", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create category: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

// @Security Bearer
// Update handles the update of an existing category with image upload.
// @Summary Update an existing category
// @Tags Category
// @Description Update an existing category record with the provided details and an optional image file
// @Accept multipart/form-data
// @Produce json
// @Param category_id formData int true "Category ID"
// @Param name formData string true "Category name"
// @Param description formData string true "Category description"
// @Param slug_category formData string true "Category slug"
// @Param image_category formData file false "New category image file"
// @Success 200 {object} response.ApiResponseCategory "Successfully updated category"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update category"
// @Router /api/category/update [post]
func (h *categoryHandleApi) Update(c echo.Context) error {
	categoryID, err := strconv.Atoi(c.FormValue("category_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid category ID",
		})
	}

	name := c.FormValue("name")
	description := c.FormValue("description")
	slugCategory := c.FormValue("slug_category")

	imagePath := ""
	file, err := c.FormFile("image_category")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		imagePath = "uploads/category/" + file.Filename
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

	req := &pb.UpdateCategoryRequest{
		CategoryId:    int32(categoryID),
		Name:          name,
		Description:   description,
		SlugCategory:  slugCategory,
		ImageCategory: imagePath,
	}

	res, err := h.client.Update(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to update category", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update category: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

// @Security Bearer
// TrashedCategory retrieves a trashed category record by its ID.
// @Summary Retrieve a trashed category
// @Tags Category
// @Description Retrieve a trashed category record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryDeleteAt "Successfully retrieved trashed category"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed category"
// @Router /api/category/trashed/{id} [get]
func (h *categoryHandleApi) TrashedCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid category ID", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid category ID",
		})
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCategoryRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedCategory(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to trashed category", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to trashed category: ",
		})
	}

	so := h.mapping.ToApiResponseCategoryDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreCategory restores a category record from the trash by its ID.
// @Summary Restore a trashed category
// @Tags Category
// @Description Restore a trashed category record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryDeleteAt "Successfully restored category"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore category"
// @Router /api/category/restore/{id} [post]
func (h *categoryHandleApi) RestoreCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid category ID", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid category ID",
		})
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCategoryRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreCategory(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to restore category", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore category: ",
		})
	}

	so := h.mapping.ToApiResponseCategoryDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteCategoryPermanent permanently deletes a category record by its ID.
// @Summary Permanently delete a category
// @Tags Category
// @Description Permanently delete a category record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "category ID"
// @Success 200 {object} response.ApiResponseCategoryDelete "Successfully deleted category record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete category:"
// @Router /api/category/delete/{id} [delete]
func (h *categoryHandleApi) DeleteCategoryPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid category ID", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid category ID",
		})
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCategoryRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteCategoryPermanent(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to delete category", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete category: ",
		})
	}

	so := h.mapping.ToApiResponseCategoryDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllCategory restores a category record from the trash by its ID.
// @Summary Restore a trashed category
// @Tags Category
// @Description Restore a trashed category record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "category ID"
// @Success 200 {object} response.ApiResponseCategoryAll "Successfully restored category all"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore category"
// @Router /api/category/restore/all [post]
func (h *categoryHandleApi) RestoreAllCategory(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllCategory(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to restore all category", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently restore all category",
		})
	}

	so := h.mapping.ToApiResponseCategoryAll(res)

	h.logger.Debug("Successfully restored all category")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllCategoryPermanent permanently deletes a category record by its ID.
// @Summary Permanently delete a category
// @Tags Category
// @Description Permanently delete a category record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "category ID"
// @Success 200 {object} response.ApiResponseCategoryAll "Successfully deleted category record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete category:"
// @Router /api/category/delete/all [post]
func (h *categoryHandleApi) DeleteAllCategoryPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllCategoryPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to permanently delete all category", zap.Error(err))

		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all category",
		})
	}

	so := h.mapping.ToApiResponseCategoryAll(res)

	h.logger.Debug("Successfully deleted all category permanently")

	return c.JSON(http.StatusOK, so)
}
