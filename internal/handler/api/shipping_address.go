package api

import (
	shippingaddress_cache "ecommerce/internal/cache/api/shipping_address"
	"ecommerce/internal/domain/requests"
	response_api "ecommerce/internal/mapper"
	"ecommerce/internal/pb"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/logger"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type shippingAddressHandleApi struct {
	client     pb.ShippingServiceClient
	logger     logger.LoggerInterface
	mapping    response_api.ShippingAddressResponseMapper
	apiHandler errors.ApiHandler
	cache      shippingaddress_cache.ShippingAddressMencache
}

func NewHandlerShippingAddress(
	router *echo.Echo,
	client pb.ShippingServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.ShippingAddressResponseMapper,
	apiHandler errors.ApiHandler,
	cache shippingaddress_cache.ShippingAddressMencache,
) *shippingAddressHandleApi {
	shippingHandler := &shippingAddressHandleApi{
		client:     client,
		logger:     logger,
		mapping:    mapping,
		apiHandler: apiHandler,
		cache:      cache,
	}

	routerShipping := router.Group("/api/shipping-address")

	routerShipping.GET(
		"",
		apiHandler.Handle("findAll", shippingHandler.FindAllShipping),
	)
	routerShipping.GET(
		"/:id",
		apiHandler.Handle("findById", shippingHandler.FindById),
	)
	routerShipping.GET(
		"/order/:id",
		apiHandler.Handle("findByOrder", shippingHandler.FindByOrder),
	)
	routerShipping.GET(
		"/active",
		apiHandler.Handle("findByActive", shippingHandler.FindByActive),
	)
	routerShipping.GET(
		"/trashed",
		apiHandler.Handle("findByTrashed", shippingHandler.FindByTrashed),
	)

	routerShipping.POST(
		"/trashed/:id",
		apiHandler.Handle("trashed", shippingHandler.TrashedShippingAddress),
	)
	routerShipping.POST(
		"/restore/:id",
		apiHandler.Handle("restore", shippingHandler.RestoreShippingAddress),
	)
	routerShipping.DELETE(
		"/permanent/:id",
		apiHandler.Handle("deletePermanent", shippingHandler.DeleteShippingAddressPermanent),
	)

	routerShipping.POST(
		"/restore/all",
		apiHandler.Handle("restoreAll", shippingHandler.RestoreAllShippingAddress),
	)
	routerShipping.POST(
		"/permanent/all",
		apiHandler.Handle("deleteAllPermanent", shippingHandler.DeleteAllShippingAddressPermanent),
	)

	return shippingHandler

}

// @Security Bearer
// @Summary Find all shipping-address
// @Tags shipping address
// @Description Retrieve a list of all shipping-address
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationShippingAddress "List of shipping-address"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve shipping-address data"
// @Router /api/shipping-address [get]
func (h *shippingAddressHandleApi) FindAllShipping(c echo.Context) error {
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

	req := &requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetShippingAddressAllCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllShippingRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindAllShipping")
	}

	apiResponse := h.mapping.ToApiResponsePaginationShippingAddress(res)

	h.cache.SetShippingAddressAllCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find shipping address by ID
// @Tags ShippingAddress
// @Description Retrieve a shipping address by ID
// @Accept json
// @Produce json
// @Param id path int true "Shipping Address ID"
// @Success 200 {object} response.ApiResponseShippingAddress "Shipping address data"
// @Failure 400 {object} response.ErrorResponse "Invalid shipping address ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve shipping address data"
// @Router /api/shipping-address/{id} [get]
func (h *shippingAddressHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedShippingAddressCache(ctx, id)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindByIdShippingRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "FindById")
	}

	apiResponse := h.mapping.ToApiResponseShippingAddress(res)

	h.cache.SetCachedShippingAddressCache(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find shipping address by order ID
// @Tags ShippingAddress
// @Description Retrieve a shipping address by order ID
// @Accept json
// @Produce json
// @Param order_id path int true "Order ID"
// @Success 200 {object} response.ApiResponseShippingAddress "Shipping address data"
// @Failure 400 {object} response.ErrorResponse "Invalid order ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve shipping address data"
// @Router /api/shipping-address/order/{order_id} [get]
func (h *shippingAddressHandleApi) FindByOrder(c echo.Context) error {
	orderID, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		return errors.NewBadRequestError("order_id is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedShippingAddressByOrderCache(ctx, orderID)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindByIdShippingRequest{
		Id: int32(orderID),
	}

	res, err := h.client.FindByOrder(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "FindByOrder")
	}

	apiResponse := h.mapping.ToApiResponseShippingAddress(res)

	h.cache.SetCachedShippingAddressByOrderCache(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active shipping-address
// @Tags ShippingAddress
// @Description Retrieve a list of active shipping-address
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationShippingAddressDeleteAt "List of active shipping-address"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve shipping data"
// @Router /api/shipping-address/active [get]
func (h *shippingAddressHandleApi) FindByActive(c echo.Context) error {
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

	req := &requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetShippingAddressActiveCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllShippingRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByActive")
	}

	apiResponse := h.mapping.ToApiResponsePaginationShippingAddressDeleteAt(res)

	h.cache.SetShippingAddressActiveCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed shipping-address records.
// @Summary Retrieve trashed shipping-address
// @Tags ShippingAddress
// @Description Retrieve a list of trashed shipping-address records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationShippingAddressDeleteAt "List of trashed shipping-address data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve shipping-address data"
// @Router /api/shipping-address/trashed [get]
func (h *shippingAddressHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetShippingAddressTrashedCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllShippingRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByTrashed")
	}

	apiResponse := h.mapping.ToApiResponsePaginationShippingAddressDeleteAt(res)

	h.cache.SetShippingAddressTrashedCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// TrashedShippingAddress retrieves a trashed shipping address record by its ID.
// @Summary Retrieve a trashed shipping address
// @Tags ShippingAddress
// @Description Retrieve a trashed shipping address record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Shipping Address ID"
// @Success 200 {object} response.ApiResponseShippingAddressDeleteAt "Successfully retrieved trashed shipping address"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed shipping address"
// @Router /api/shipping-address/trashed/{id} [get]
func (h *shippingAddressHandleApi) TrashedShippingAddress(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdShippingRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedShipping(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "TrashedShipping")
	}

	so := h.mapping.ToApiResponseShippingAddressDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreShippingAddress restores a shipping address record from the trash by its ID.
// @Summary Restore a trashed shipping address
// @Tags ShippingAddress
// @Description Restore a trashed shipping address record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Shipping Address ID"
// @Success 200 {object} response.ApiResponseShippingAddressDeleteAt "Successfully restored shipping address"
// @Failure 400 {object} response.ErrorResponse "Invalid shipping address ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore shipping address"
// @Router /api/shipping-address/restore/{id} [post]
func (h *shippingAddressHandleApi) RestoreShippingAddress(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdShippingRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreShipping(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "RestoreShipping")
	}

	so := h.mapping.ToApiResponseShippingAddressDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteShippingAddressPermanent permanently deletes a shipping address record by its ID.
// @Summary Permanently delete a shipping address
// @Tags ShippingAddress
// @Description Permanently delete a shipping address record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Shipping Address ID"
// @Success 200 {object} response.ApiResponseShippingAddressDelete "Successfully deleted shipping address record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete shipping address:"
// @Router /api/shipping-address/delete/{id} [delete]
func (h *shippingAddressHandleApi) DeleteShippingAddressPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdShippingRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteShippingPermanent(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "DeleteShipping")
	}

	so := h.mapping.ToApiResponseShippingAddressDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllShippingAddress restores all trashed shipping address records.
// @Summary Restore all trashed shipping addresses
// @Tags ShippingAddress
// @Description Restore all trashed shipping address records.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseShippingAddressAll "Successfully restored all shipping addresses"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all shipping addresses"
// @Router /api/shipping-address/restore/all [post]
func (h *shippingAddressHandleApi) RestoreAllShippingAddress(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllShipping(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "RestoreAll")
	}

	so := h.mapping.ToApiResponseShippingAddressAll(res)

	h.logger.Debug("Successfully restored all shipping addresses")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllShippingAddressPermanent permanently deletes all trashed shipping address records.
// @Summary Permanently delete all trashed shipping addresses
// @Tags ShippingAddress
// @Description Permanently delete all trashed shipping address records.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseShippingAddressAll "Successfully deleted all shipping addresses permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to delete all shipping addresses"
// @Router /api/shipping-address/delete/all [post]
func (h *shippingAddressHandleApi) DeleteAllShippingAddressPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllShippingPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "DeleteAllShipping")
	}

	so := h.mapping.ToApiResponseShippingAddressAll(res)

	h.logger.Debug("Successfully deleted all shipping addresses permanently")

	return c.JSON(http.StatusOK, so)
}

func (h *shippingAddressHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Shipping Address").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Shipping Address already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Shipping Address service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *shippingAddressHandleApi) parseValidationErrors(err error) []errors.ValidationError {
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

func (h *shippingAddressHandleApi) getValidationMessage(fe validator.FieldError) string {
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
