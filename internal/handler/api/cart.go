package api

import (
	cart_cache "ecommerce/internal/cache/api/cart"
	"ecommerce/internal/domain/requests"
	response_api "ecommerce/internal/mapper"
	"fmt"

	"ecommerce/internal/pb"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/logger"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type cartHandleApi struct {
	client     pb.CartServiceClient
	logger     logger.LoggerInterface
	mapping    response_api.CartResponseMapper
	apiHandler errors.ApiHandler
	cache      cart_cache.CartMencache
}

func NewHandlerCart(
	router *echo.Echo,
	client pb.CartServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.CartResponseMapper,
	apiHandler errors.ApiHandler,
	cache cart_cache.CartMencache,
) *cartHandleApi {
	cartHandler := &cartHandleApi{
		client:     client,
		logger:     logger,
		mapping:    mapping,
		apiHandler: apiHandler,
		cache:      cache,
	}

	routerCart := router.Group("/api/cart")

	routerCart.GET("", apiHandler.Handle("findAll", cartHandler.FindAll))
	routerCart.POST("/create", apiHandler.Handle("create", cartHandler.Create))
	routerCart.DELETE("/:id", apiHandler.Handle("delete", cartHandler.Delete))
	routerCart.POST("/delete-all", apiHandler.Handle("deleteAll", cartHandler.DeleteAll))

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
		return errors.NewBadRequestError("user_id is required")
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

	reqCache := &requests.FindAllCarts{
		UserID:   userID,
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedCarts(ctx, reqCache)

	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindAllCartRequest{
		UserId:   int32(userID),
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch cart details", zap.Error(err))
		return err
	}

	so := h.mapping.ToApiResponseCartPagination(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Create a new cart
// @Tags Cart
// @Description Create a new cart item
// @Accept json
// @Produce json
// @Param body body requests.CreateCartRequest true "Cart creation data"
// @Success 200 {object} response.ApiResponseCart "Created cart details"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 502 {object} response.ErrorResponse "Failed to create cart"
// @Router /api/cart [post]
func (h *cartHandleApi) Create(c echo.Context) error {
	var body requests.CreateCartRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()
	req := &pb.CreateCartRequest{
		Quantity:  int32(body.Quantity),
		ProductId: int32(body.ProductID),
		UserId:    int32(body.UserID),
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	so := h.mapping.ToApiResponseCart(res)
	return c.JSON(http.StatusCreated, so)
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
	id := c.Param("id")

	idStr, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCartRequest{
		Id: int32(idStr),
	}

	res, err := h.client.Delete(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Delete")
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
// @Param request body requests.DeleteCartRequest true "Cart IDs"
// @Success 200 {object} response.ApiResponseCartAll "Successfully deleted carts"
// @Failure 500 {object} response.ErrorResponse "Failed to delete carts"
// @Router /api/cart/delete-all [post]
func (h *cartHandleApi) DeleteAll(c echo.Context) error {
	var req requests.DeleteCartRequest

	if err := c.Bind(&req); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := req.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	cartIdsPb := make([]int32, len(req.CartIds))
	for i, id := range req.CartIds {
		cartIdsPb[i] = int32(id)
	}

	reqPb := &pb.DeleteCartRequest{
		CartIds: cartIdsPb,
	}

	res, err := h.client.DeleteAll(ctx, reqPb)

	if err != nil {
		return h.handleGrpcError(err, "DeleteAll")
	}

	so := h.mapping.ToApiResponseCartAll(res)
	return c.JSON(http.StatusOK, so)
}

func (h *cartHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Cart").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Cart already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Cart Address service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *cartHandleApi) parseValidationErrors(err error) []errors.ValidationError {
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

func (h *cartHandleApi) getValidationMessage(fe validator.FieldError) string {
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
