package api

import (
	"ecommerce/internal/domain/response"
	response_api "ecommerce/internal/mapper/response/api"
	"ecommerce/internal/pb"
	"ecommerce/pkg/logger"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type cartHandleApi struct {
	client  pb.CartServiceClient
	logger  logger.LoggerInterface
	mapping response_api.CartResponseMapper
}

func NewHandlerCart(
	router *echo.Echo,
	client pb.CartServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.CartResponseMapper,
) *cartHandleApi {
	cartHandler := &cartHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
	}

	routerCart := router.Group("/api/cart")
	routerCart.GET("", cartHandler.FindAll)
	routerCart.DELETE("/:id", cartHandler.Delete)
	routerCart.POST("/delete-all", cartHandler.DeleteAll)

	return cartHandler
}

// @Security Bearer
// @Summary Find all carts
// @Tags Cart
// @Description Retrieve a list of all carts
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponseCartPagination "List of carts"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve cart data"
// @Router /api/cart [get]
func (h *cartHandleApi) FindAll(c echo.Context) error {
	userID, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil || userID <= 0 {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid user ID",
		})
	}

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

	req := &pb.FindAllCartRequest{
		UserId:   int32(userID),
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to retrieve cart data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve cart data",
		})
	}

	so := h.mapping.ToApiResponseCartPagination(res)
	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Delete a cart
// @Tags Cart
// @Description Delete a cart by ID
// @Accept json
// @Produce json
// @Param id path int true "Cart ID"
// @Success 200 {object} response.ApiResponseCartDelete "Successfully deleted cart"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete cart"
// @Router /api/cart/{id} [delete]
func (h *cartHandleApi) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Debug("Invalid cart ID", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid cart ID",
		})
	}

	ctx := c.Request().Context()
	req := &pb.FindByIdCartRequest{
		Id: int32(id),
	}

	res, err := h.client.Delete(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to delete cart", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete cart",
		})
	}

	so := h.mapping.ToApiResponseCartDelete(res)
	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Delete multiple carts
// @Tags Cart
// @Description Delete multiple carts by IDs
// @Accept json
// @Produce json
// @Param request body pb.DeleteCartRequest true "Cart IDs"
// @Success 200 {object} response.ApiResponseCartAll "Successfully deleted carts"
// @Failure 500 {object} response.ErrorResponse "Failed to delete carts"
// @Router /api/cart/delete-all [post]
func (h *cartHandleApi) DeleteAll(c echo.Context) error {
	var req pb.DeleteCartRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Debug("Invalid request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	ctx := c.Request().Context()
	res, err := h.client.DeleteAll(ctx, &req)
	if err != nil {
		h.logger.Debug("Failed to delete carts", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete carts",
		})
	}

	so := h.mapping.ToApiResponseCartAll(res)
	return c.JSON(http.StatusOK, so)
}
