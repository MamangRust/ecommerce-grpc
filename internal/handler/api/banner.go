package api

import (
	banner_cache "ecommerce/internal/cache/api/banner"
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
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type bannerHandleApi struct {
	client     pb.BannerServiceClient
	logger     logger.LoggerInterface
	mapping    response_api.BannerResponseMapper
	apiHandler errors.ApiHandler
	cache      banner_cache.BannerMencache
}

func NewHandleBanner(
	router *echo.Echo,
	client pb.BannerServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.BannerResponseMapper,
	apiHandler errors.ApiHandler,
	cache banner_cache.BannerMencache,
) *bannerHandleApi {
	bannerHandler := &bannerHandleApi{
		client:     client,
		logger:     logger,
		mapping:    mapping,
		apiHandler: apiHandler,
		cache:      cache,
	}

	routerBanner := router.Group("/api/banner")

	routerBanner.GET("", apiHandler.Handle("findAll", bannerHandler.FindAllBanner))
	routerBanner.GET("/:id", apiHandler.Handle("findById", bannerHandler.FindById))
	routerBanner.GET("/active", apiHandler.Handle("findByActive", bannerHandler.FindByActive))
	routerBanner.GET("/trashed", apiHandler.Handle("findByTrashed", bannerHandler.FindByTrashed))

	routerBanner.POST("/create", apiHandler.Handle("create", bannerHandler.Create))
	routerBanner.POST("/update/:id", apiHandler.Handle("update", bannerHandler.Update))

	routerBanner.POST("/trashed/:id", apiHandler.Handle("trashed", bannerHandler.TrashedBanner))
	routerBanner.POST("/restore/:id", apiHandler.Handle("restore", bannerHandler.RestoreBanner))
	routerBanner.DELETE("/permanent/:id", apiHandler.Handle("deletePermanent", bannerHandler.DeleteBannerPermanent))

	routerBanner.POST("/restore/all", apiHandler.Handle("restoreAll", bannerHandler.RestoreAllBanner))
	routerBanner.POST("/permanent/all", apiHandler.Handle("deleteAllPermanent", bannerHandler.DeleteAllBannerPermanent))

	return bannerHandler
}

// @Security Bearer
// @Summary Find all banners
// @Tags Banner
// @Description Retrieve a list of all banners
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationBanner "List of banners"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve banner data"
// @Router /api/banner [get]
func (h *bannerHandleApi) FindAllBanner(c echo.Context) error {
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

	req := &requests.FindAllBanner{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedBanners(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllBannerRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindAll")
	}

	apiResponse := h.mapping.ToApiResponsePaginationBanner(res)

	h.cache.SetCachedBanners(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find banner by ID
// @Tags Banner
// @Description Retrieve a banner by ID
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBanner "Banner data"
// @Failure 400 {object} response.ErrorResponse "Invalid banner ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve banner data"
// @Router /api/banner/{id} [get]
func (h *bannerHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdBannerRequest{
		Id: int32(id),
	}

	cachedData, found := h.cache.GetCachedBanner(ctx, id)
	if found {
		h.logger.Debug("Returning banner from cache", zap.Int("id", id))
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "FindById")
	}

	so := h.mapping.ToApiResponseBanner(res)

	h.cache.SetCachedBanner(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active banners
// @Tags Banner
// @Description Retrieve a list of active banners
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationBannerDeleteAt "List of active banners"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve banner data"
// @Router /api/banner/active [get]
func (h *bannerHandleApi) FindByActive(c echo.Context) error {
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
	req := &requests.FindAllBanner{Page: page, PageSize: pageSize, Search: search}

	cachedData, found := h.cache.GetCachedActiveBanners(ctx, req)
	if found {
		h.logger.Debug("Returning active banners from cache")
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllBannerRequest{Page: int32(page), PageSize: int32(pageSize), Search: search}
	res, err := h.client.FindByActive(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByActive")
	}

	apiResponse := h.mapping.ToApiResponsePaginationBannerDeleteAt(res)
	h.cache.SetCachedActiveBanners(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve trashed banners
// @Tags Banner
// @Description Retrieve a list of trashed banners
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationBannerDeleteAt "List of active banners"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve banner data"
// @Router /api/banner/trashed [get]
func (h *bannerHandleApi) FindByTrashed(c echo.Context) error {
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
	req := &requests.FindAllBanner{Page: page, PageSize: pageSize, Search: search}

	cachedData, found := h.cache.GetCachedTrashedBanners(ctx, req)
	if found {
		h.logger.Debug("Returning trashed banners from cache")
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllBannerRequest{Page: int32(page), PageSize: int32(pageSize), Search: search}
	res, err := h.client.FindByTrashed(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByTrashed")
	}

	apiResponse := h.mapping.ToApiResponsePaginationBannerDeleteAt(res)
	h.cache.SetCachedTrashedBanners(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// Create handles the creation of a new banner.
// @Summary Create a new banner
// @Tags Banner
// @Description Create a new banner with the provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateBannerRequest true "Create banner request"
// @Success 200 {object} response.ApiResponseBanner "Successfully created banner"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create banner"
// @Router /api/banner/create [post]
func (h *bannerHandleApi) Create(c echo.Context) error {
	var body requests.CreateBannerRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	req := &pb.CreateBannerRequest{
		Name:      body.Name,
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
		StartTime: body.StartTime,
		EndTime:   body.EndTime,
		IsActive:  body.IsActive,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		h.logger.Error("Banner creation failed", zap.Error(err))
		return err
	}

	so := h.mapping.ToApiResponseBanner(res)

	h.cache.SetCachedBanner(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Update handles the update of an existing banner record.
// @Summary Update an existing banner
// @Tags Banner
// @Description Update an existing banner record with the provided details
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Param request body requests.UpdateBannerRequest true "Update banner request"
// @Success 200 {object} response.ApiResponseBanner "Successfully updated banner"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update banner"
// @Router /api/banner/update/{id} [post]
func (h *bannerHandleApi) Update(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	var body requests.UpdateBannerRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	req := &pb.UpdateBannerRequest{
		BannerId:  int32(idInt),
		Name:      body.Name,
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
		StartTime: body.StartTime,
		EndTime:   body.EndTime,
		IsActive:  body.IsActive,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		h.logger.Error("Banner update failed", zap.Error(err))
		return err
	}

	so := h.mapping.ToApiResponseBanner(res)

	h.cache.DeleteBannerCache(ctx, idInt)

	h.cache.SetCachedBanner(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedBanner retrieves a trashed Banner record by its ID.
// @Summary Retrieve a trashed Banner
// @Tags Banner
// @Description Retrieve a trashed Banner record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBannerDeleteAt "Successfully retrieved trashed Banner"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed Banner"
// @Router /api/banner/trashed/{id} [get]
func (h *bannerHandleApi) TrashedBanner(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdBannerRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedBanner(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Trashed")
	}

	so := h.mapping.ToApiResponseBannerDeleteAt(res)

	h.cache.DeleteBannerCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreBanner restores a Banner record from the trash by its ID.
// @Summary Restore a trashed Banner
// @Tags Banner
// @Description Restore a trashed Banner record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBannerDeleteAt "Successfully restored Banner"
// @Failure 400 {object} response.ErrorResponse "Invalid Banner ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore Banner"
// @Router /api/banner/restore/{id} [post]
func (h *bannerHandleApi) RestoreBanner(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdBannerRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreBanner(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Restore")
	}

	so := h.mapping.ToApiResponseBannerDeleteAt(res)

	h.cache.DeleteBannerCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteBannerPermanent permanently deletes a Banner record by its ID.
// @Summary Permanently delete a Banner
// @Tags Banner
// @Description Permanently delete a Banner record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBannerDelete "Successfully deleted Banner record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete Banner:"
// @Router /api/banner/delete/{id} [delete]
func (h *bannerHandleApi) DeleteBannerPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdBannerRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteBannerPermanent(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "DeletePermanent")
	}

	so := h.mapping.ToApiResponseBannerDelete(res)

	h.cache.DeleteBannerCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllBanner restores a Banner record from the trash by its ID.
// @Summary Restore a trashed Banner
// @Tags Banner
// @Description Restore a trashed Banner record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBannerAll "Successfully restored Banner all"
// @Failure 400 {object} response.ErrorResponse "Invalid Banner ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore Banner"
// @Router /api/banner/restore/all [post]
func (h *bannerHandleApi) RestoreAllBanner(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllBanner(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "RestoreAll")
	}

	so := h.mapping.ToApiResponseBannerAll(res)

	h.logger.Debug("Successfully restored all Banner")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllBannerPermanent permanently deletes a banner record by its ID.
// @Summary Permanently delete a banner
// @Tags Banner
// @Description Permanently delete a banner record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "banner ID"
// @Success 200 {object} response.ApiResponseBannerAll "Successfully deleted banner record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete banner:"
// @Router /api/banner/delete/all [post]
func (h *bannerHandleApi) DeleteAllBannerPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllBannerPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "DeleteAll")
	}

	so := h.mapping.ToApiResponseBannerAll(res)

	h.logger.Debug("Successfully deleted all banner permanently")

	return c.JSON(http.StatusOK, so)
}

func (h *bannerHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Banner").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Banner already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Banner Address service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *bannerHandleApi) parseValidationErrors(err error) []errors.ValidationError {
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

func (h *bannerHandleApi) getValidationMessage(fe validator.FieldError) string {
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
