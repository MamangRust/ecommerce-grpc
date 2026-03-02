package api

import (
	merchantbusiness_cache "ecommerce/internal/cache/api/merchant_business"
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

type merchantBusinessHandleApi struct {
	client          pb.MerchantBusinessServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.MerchantBusinessResponseMapper
	mappingMerchant response_api.MerchantResponseMapper
	apiHandler      errors.ApiHandler
	cache           merchantbusiness_cache.MerchantBusinessMencache
}

func NewHandlerMerchantBusiness(
	router *echo.Echo,
	client pb.MerchantBusinessServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.MerchantBusinessResponseMapper,
	mappingMerchant response_api.MerchantResponseMapper,
	apiHandler errors.ApiHandler,
	cache merchantbusiness_cache.MerchantBusinessMencache,
) *merchantBusinessHandleApi {
	merchantBusinessHandler := &merchantBusinessHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		mappingMerchant: mappingMerchant,
		apiHandler:      apiHandler,
		cache:           cache,
	}

	routerMerchantBusiness := router.Group("/api/merchant-business")

	routerMerchantBusiness.GET("", apiHandler.Handle("findAll", merchantBusinessHandler.FindAllMerchantBusiness))
	routerMerchantBusiness.GET("/:id", apiHandler.Handle("findById", merchantBusinessHandler.FindById))
	routerMerchantBusiness.GET("/active", apiHandler.Handle("findByActive", merchantBusinessHandler.FindByActive))
	routerMerchantBusiness.GET("/trashed", apiHandler.Handle("findByTrashed", merchantBusinessHandler.FindByTrashed))

	routerMerchantBusiness.POST("/create", apiHandler.Handle("create", merchantBusinessHandler.Create))
	routerMerchantBusiness.POST("/update/:id", apiHandler.Handle("update", merchantBusinessHandler.Update))

	routerMerchantBusiness.POST("/trashed/:id", apiHandler.Handle("trashed", merchantBusinessHandler.TrashedMerchant))
	routerMerchantBusiness.POST("/restore/:id", apiHandler.Handle("restore", merchantBusinessHandler.RestoreMerchant))
	routerMerchantBusiness.DELETE("/permanent/:id", apiHandler.Handle("deletePermanent", merchantBusinessHandler.DeleteMerchantPermanent))

	routerMerchantBusiness.POST("/restore/all", apiHandler.Handle("restoreAll", merchantBusinessHandler.RestoreAllMerchant))
	routerMerchantBusiness.POST("/permanent/all", apiHandler.Handle("deleteAllPermanent", merchantBusinessHandler.DeleteAllMerchantPermanent))

	return merchantBusinessHandler
}

// @Security Bearer
// @Summary Find all merchant
// @Tags MerchantCertification
// @Description Retrieve a list of all merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantBusiness "List of merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-business [get]
func (h *merchantBusinessHandleApi) FindAllMerchantBusiness(c echo.Context) error {
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

	req := &requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedMerchantBusinessAll(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindAll")
	}

	apiResponse := h.mapping.ToApiResponsePaginationMerchantBusiness(res)

	h.cache.SetCachedMerchantBusinessAll(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find merchant by ID
// @Tags MerchantCertification
// @Description Retrieve a merchant by ID
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantBusiness "merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-business/{id} [get]
func (h *merchantBusinessHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedMerchantBusiness(ctx, id)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindByIdMerchantBusinessRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "FindById")
	}

	apiResponse := h.mapping.ToApiResponseMerchantBusiness(res)

	h.cache.SetCachedMerchantBusiness(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active merchant
// @Tags MerchantCertification
// @Description Retrieve a list of active merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantBusinessDeleteAt "List of active merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-business/active [get]
func (h *merchantBusinessHandleApi) FindByActive(c echo.Context) error {
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

	req := &requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedMerchantBusinessActive(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByActive")
	}

	apiResponse := h.mapping.ToApiResponsePaginationMerchantBusinessDeleteAt(res)

	h.cache.SetCachedMerchantBusinessActive(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed merchant records.
// @Summary Retrieve trashed merchant
// @Tags MerchantCertification
// @Description Retrieve a list of trashed merchant records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationMerchantBusinessDeleteAt "List of trashed merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-business/trashed [get]
func (h *merchantBusinessHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedMerchantBusinessTrashed(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByTrashed")
	}

	apiResponse := h.mapping.ToApiResponsePaginationMerchantBusinessDeleteAt(res)

	h.cache.SetCachedMerchantBusinessTrashed(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Create a new merchant business information
// @Tags MerchantBusiness
// @Description Create merchant business info (e.g., type, tax ID, website, etc.)
// @Accept json
// @Produce json
// @Param request body requests.CreateMerchantBusinessInformationRequest true "Create merchant business request"
// @Success 200 {object} response.ApiResponseMerchantBusiness "Successfully created merchant business info"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create merchant business info"
// @Router /api/merchant-business/create [post]
func (h *merchantBusinessHandleApi) Create(c echo.Context) error {
	var body requests.CreateMerchantBusinessInformationRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	req := &pb.CreateMerchantBusinessRequest{
		MerchantId:        int32(body.MerchantID),
		BusinessType:      body.BusinessType,
		TaxId:             body.TaxID,
		EstablishedYear:   int32(body.EstablishedYear),
		NumberOfEmployees: int32(body.NumberOfEmployees),
		WebsiteUrl:        body.WebsiteUrl,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	so := h.mapping.ToApiResponseMerchantBusiness(res)

	h.cache.SetCachedMerchantBusiness(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Update existing merchant business information
// @Tags MerchantBusiness
// @Description Update merchant business info by ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant Business Info ID"
// @Param request body requests.UpdateMerchantBusinessInformationRequest true "Update merchant business request"
// @Success 200 {object} response.ApiResponseMerchantBusiness "Successfully updated merchant business info"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update merchant business info"
// @Router /api/merchant-business/update/{id} [post]
func (h *merchantBusinessHandleApi) Update(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	var body requests.UpdateMerchantBusinessInformationRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	req := &pb.UpdateMerchantBusinessRequest{
		MerchantBusinessInfoId: int32(idInt),
		BusinessType:           body.BusinessType,
		TaxId:                  body.TaxID,
		EstablishedYear:        int32(body.EstablishedYear),
		NumberOfEmployees:      int32(body.NumberOfEmployees),
		WebsiteUrl:             body.WebsiteUrl,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "Update")
	}

	so := h.mapping.ToApiResponseMerchantBusiness(res)

	h.cache.DeleteMerchantBusinessCache(ctx, idInt)
	h.cache.SetCachedMerchantBusiness(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedMerchant retrieves a trashed merchant record by its ID.
// @Summary Retrieve a trashed merchant
// @Tags MerchantCertification
// @Description Retrieve a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantBusinessDeleteAt "Successfully retrieved trashed merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed merchant"
// @Router /api/merchant-business/trashed/{id} [get]
func (h *merchantBusinessHandleApi) TrashedMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return h.handleGrpcError(err, "FindById")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantBusinessRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedMerchantBusiness(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Trashed")
	}

	so := h.mapping.ToApiResponseMerchantBusinessDeleteAt(res)

	h.cache.DeleteMerchantBusinessCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags MerchantCertification
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantBusinessDeleteAt "Successfully restored merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant-business/restore/{id} [post]
func (h *merchantBusinessHandleApi) RestoreMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantBusinessRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreMerchantBusiness(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Restore")
	}

	so := h.mapping.ToApiResponseMerchantBusinessDeleteAt(res)

	h.cache.DeleteMerchantBusinessCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags MerchantCertification
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/merchant-business/delete/{id} [delete]
func (h *merchantBusinessHandleApi) DeleteMerchantPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantBusinessRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteMerchantBusinessPermanent(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "DeleteMerchant")
	}

	so := h.mappingMerchant.ToApiResponseMerchantDelete(res)

	h.cache.DeleteMerchantBusinessCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags MerchantCertification
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored merchant all"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant-business/restore/all [post]
func (h *merchantBusinessHandleApi) RestoreAllMerchant(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllMerchantBusiness(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "RestoreAll")
	}

	so := h.mappingMerchant.ToApiResponseMerchantAll(res)

	h.logger.Debug("Successfully restored all merchant")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags MerchantCertification
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/merchant-business/delete/all [post]
func (h *merchantBusinessHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllMerchantBusinessPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "DeleteAll")
	}

	so := h.mappingMerchant.ToApiResponseMerchantAll(res)

	h.logger.Debug("Successfully deleted all merchant permanently")

	return c.JSON(http.StatusOK, so)
}

func (h *merchantBusinessHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Merchant Business").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Merchant Business already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Merchant Business service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *merchantBusinessHandleApi) parseValidationErrors(err error) []errors.ValidationError {
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

func (h *merchantBusinessHandleApi) getValidationMessage(fe validator.FieldError) string {
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
