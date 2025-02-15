package api

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_api "ecommerce/internal/mapper/response/api"
	"ecommerce/internal/pb"
	"ecommerce/pkg/logger"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type reviewHandleApi struct {
	client  pb.ReviewServiceClient
	logger  logger.LoggerInterface
	mapping response_api.ReviewResponseMapper
}

func NewHandlerReview(
	router *echo.Echo,
	client pb.ReviewServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.ReviewResponseMapper,
) *reviewHandleApi {
	reviewHandler := &reviewHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
	}

	routercategory := router.Group("/api/category")

	routercategory.GET("", reviewHandler.FindAll)
	routercategory.GET("/product/:id", reviewHandler.FindByProduct)
	routercategory.GET("/active", reviewHandler.FindByActive)
	routercategory.GET("/trashed", reviewHandler.FindByTrashed)

	routercategory.POST("/create", reviewHandler.Create)
	routercategory.POST("/update/:id", reviewHandler.Update)

	routercategory.POST("/trashed/:id", reviewHandler.TrashedReview)
	routercategory.POST("/restore/:id", reviewHandler.RestoreReview)
	routercategory.DELETE("/permanent/:id", reviewHandler.DeleteReviewPermanent)

	routercategory.POST("/restore/all", reviewHandler.RestoreAllReview)
	routercategory.POST("/permanent/all", reviewHandler.DeleteAllReviewPermanent)

	return reviewHandler
}

// @Security Bearer
// @Summary Find all review
// @Tags Review
// @Description Retrieve a list of all review
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReview "List of review"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve review data"
// @Router /api/review [get]
func (h *reviewHandleApi) FindAll(c echo.Context) error {
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

	req := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve review data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve review data: ",
		})
	}

	so := h.mapping.ToApiResponsePaginationReview(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find reviews by product ID
// @Tags Review
// @Description Retrieve a list of reviews for a specific product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReview "List of reviews for the product"
// @Failure 400 {object} response.ErrorResponse "Invalid product ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve review data"
// @Router /api/review/product/{id} [get]
func (h *reviewHandleApi) FindByProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid product ID",
		})
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")
	ctx := c.Request().Context()

	req := &pb.FindAllReviewProductRequest{
		ProductId: int32(id),
		Page:      int32(page),
		PageSize:  int32(pageSize),
		Search:    search,
	}

	res, err := h.client.FindByProduct(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to retrieve review data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve review data",
		})
	}

	so := h.mapping.ToApiResponsePaginationReview(res)
	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active review
// @Tags Review
// @Description Retrieve a list of active review
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationReviewDeleteAt "List of active review"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve review data"
// @Router /api/review/active [get]
func (h *reviewHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllReviewRequest{
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

	so := h.mapping.ToApiResponsePaginationReviewDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed review records.
// @Summary Retrieve trashed review
// @Tags Review
// @Description Retrieve a list of trashed review records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationReviewDeleteAt "List of trashed review data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve review data"
// @Router /api/review/trashed [get]
func (h *reviewHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve review data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve review data: ",
		})
	}

	so := h.mapping.ToApiResponsePaginationReviewDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

func (h *reviewHandleApi) Create(c echo.Context) error {
	var req requests.CreateReviewRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	ctx := c.Request().Context()

	grpcReq := &pb.CreateReviewRequest{
		UserId:    int32(req.UserID),
		ProductId: int32(req.ProductID),
		Comment:   req.Comment,
		Rating:    int32(req.Rating),
	}

	res, err := h.client.Create(ctx, grpcReq)
	if err != nil {
		h.logger.Debug("Failed to create review", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create review: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *reviewHandleApi) Update(c echo.Context) error {
	var req requests.UpdateReviewRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	ctx := c.Request().Context()

	grpcReq := &pb.UpdateReviewRequest{
		ReviewId: int32(req.ReviewID),
		Name:     req.Name,
		Comment:  req.Comment,
		Rating:   int32(req.Rating),
	}

	res, err := h.client.Update(ctx, grpcReq)
	if err != nil {
		h.logger.Debug("Failed to update review", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update review: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

// @Security Bearer
// TrashedReview retrieves a trashed review record by its ID.
// @Summary Retrieve a trashed review
// @Tags Review
// @Description Retrieve a trashed review record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "review ID"
// @Success 200 {object} response.ApiResponseReviewDeleteAt "Successfully retrieved trashed review"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed review"
// @Router /api/review/trashed/{id} [get]
func (h *reviewHandleApi) TrashedReview(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid review ID", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid review ID",
		})
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdReviewRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedReview(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to trashed review", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to trashed review: ",
		})
	}

	so := h.mapping.ToApiResponseReviewDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreReview restores a review record from the trash by its ID.
// @Summary Restore a trashed review
// @Tags Review
// @Description Restore a trashed review record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} response.ApiResponseReviewDeleteAt "Successfully restored review"
// @Failure 400 {object} response.ErrorResponse "Invalid review ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore review"
// @Router /api/review/restore/{id} [post]
func (h *reviewHandleApi) RestoreReview(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid review ID", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid review ID",
		})
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdReviewRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreReview(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to restore review", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore review: ",
		})
	}

	so := h.mapping.ToApiResponseReviewDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteReviewPermanent permanently deletes a review record by its ID.
// @Summary Permanently delete a review
// @Tags Category
// @Description Permanently delete a review record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "review ID"
// @Success 200 {object} response.ApiResponseReviewDelete "Successfully deleted review record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete review:"
// @Router /api/review/delete/{id} [delete]
func (h *reviewHandleApi) DeleteReviewPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid review ID", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid review ID",
		})
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdReviewRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteReviewPermanent(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to delete review", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete review: ",
		})
	}

	so := h.mapping.ToApiResponseReviewDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllReview restores a review record from the trash by its ID.
// @Summary Restore a trashed review
// @Tags Category
// @Description Restore a trashed review record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "review ID"
// @Success 200 {object} response.ApiResponseReviewAll "Successfully restored review all"
// @Failure 400 {object} response.ErrorResponse "Invalid review ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore review"
// @Router /api/review/restore/all [post]
func (h *reviewHandleApi) RestoreAllReview(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllReview(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to restore all review", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently restore all review",
		})
	}

	so := h.mapping.ToApiResponseReviewAll(res)

	h.logger.Debug("Successfully restored all review")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllReviewPermanent permanently deletes a review record by its ID.
// @Summary Permanently delete a review
// @Tags Review
// @Description Permanently delete a review record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "review ID"
// @Success 200 {object} response.ApiResponseReviewAll "Successfully deleted review record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete review:"
// @Router /api/review/delete/all [post]
func (h *reviewHandleApi) DeleteAllReviewPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllReviewPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to permanently delete all category", zap.Error(err))

		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all category",
		})
	}

	so := h.mapping.ToApiResponseReviewAll(res)

	h.logger.Debug("Successfully deleted all category permanently")

	return c.JSON(http.StatusOK, so)
}
