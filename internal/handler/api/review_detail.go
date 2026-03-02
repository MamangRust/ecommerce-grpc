package api

import (
	reviewdetail_cache "ecommerce/internal/cache/api/review_detail"
	"ecommerce/internal/domain/requests"
	response_api "ecommerce/internal/mapper"
	"fmt"

	"ecommerce/internal/pb"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/upload_image"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type reviewDetailHandleApi struct {
	client        pb.ReviewDetailServiceClient
	logger        logger.LoggerInterface
	mapping       response_api.ReviewDetailResponseMapper
	mappingReview response_api.ReviewResponseMapper
	upload_image  upload_image.ImageUploads
	apiHandler    errors.ApiHandler
	cache         reviewdetail_cache.ReviewDetailMencache
}

func NewHandlerReviewDetail(
	router *echo.Echo,
	client pb.ReviewDetailServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.ReviewDetailResponseMapper,
	mappingReview response_api.ReviewResponseMapper,
	upload_image upload_image.ImageUploads,
	apiHandler errors.ApiHandler,
	cache reviewdetail_cache.ReviewDetailMencache,
) *reviewDetailHandleApi {
	reviewDetailHandler := &reviewDetailHandleApi{
		client:        client,
		logger:        logger,
		mapping:       mapping,
		mappingReview: mappingReview,
		upload_image:  upload_image,
		apiHandler:    apiHandler,
		cache:         cache,
	}

	routerReviewDetail := router.Group("/api/review-detail")

	routerReviewDetail.GET(
		"",
		apiHandler.Handle("findAll", reviewDetailHandler.FindAllReviewDetail),
	)
	routerReviewDetail.GET(
		"/:id",
		apiHandler.Handle("findById", reviewDetailHandler.FindById),
	)
	routerReviewDetail.GET(
		"/active",
		apiHandler.Handle("findByActive", reviewDetailHandler.FindByActive),
	)
	routerReviewDetail.GET(
		"/trashed",
		apiHandler.Handle("findByTrashed", reviewDetailHandler.FindByTrashed),
	)

	routerReviewDetail.POST(
		"/create",
		apiHandler.Handle("create", reviewDetailHandler.Create),
	)
	routerReviewDetail.POST(
		"/update/:id",
		apiHandler.Handle("update", reviewDetailHandler.Update),
	)

	routerReviewDetail.POST(
		"/trashed/:id",
		apiHandler.Handle("trashed", reviewDetailHandler.TrashedMerchant),
	)
	routerReviewDetail.POST(
		"/restore/:id",
		apiHandler.Handle("restore", reviewDetailHandler.RestoreMerchant),
	)
	routerReviewDetail.DELETE(
		"/permanent/:id",
		apiHandler.Handle("deletePermanent", reviewDetailHandler.DeleteMerchantPermanent),
	)

	routerReviewDetail.POST(
		"/restore/all",
		apiHandler.Handle("restoreAll", reviewDetailHandler.RestoreAllMerchant),
	)
	routerReviewDetail.POST(
		"/permanent/all",
		apiHandler.Handle("deleteAllPermanent", reviewDetailHandler.DeleteAllMerchantPermanent),
	)

	return reviewDetailHandler
}

// @Security Bearer
// @Summary Find all merchant
// @Tags ReviewDetail
// @Description Retrieve a list of all merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewDetails "List of merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/review-detail [get]
func (h *reviewDetailHandleApi) FindAllReviewDetail(c echo.Context) error {
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

	req := &requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetReviewDetailAllCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindAllReview")
	}

	apiResponse := h.mapping.ToApiResponsePaginationReviewDetail(res)

	h.cache.SetReviewDetailAllCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find merchant by ID
// @Tags ReviewDetail
// @Description Retrieve a merchant by ID
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseReviewDetail "merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/review-detail/{id} [get]
func (h *reviewDetailHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedReviewDetailCache(ctx, id)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindByIdReviewDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "FindById")
	}

	apiResponse := h.mapping.ToApiResponseReviewDetail(res)

	h.cache.SetCachedReviewDetailCache(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active merchant
// @Tags ReviewDetail
// @Description Retrieve a list of active merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewDetailsDeleteAt "List of active merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/review-detail/active [get]
func (h *reviewDetailHandleApi) FindByActive(c echo.Context) error {
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

	req := &requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetReviewDetailActiveCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByActive")
	}

	apiResponse := h.mapping.ToApiResponsePaginationReviewDetailDeleteAt(res)

	h.cache.SetReviewDetailActiveCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed merchant records.
// @Summary Retrieve trashed merchant
// @Tags ReviewDetail
// @Description Retrieve a list of trashed merchant records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationReviewDetailsDeleteAt "List of trashed merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/review-detail/trashed [get]
func (h *reviewDetailHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetReviewDetailTrashedCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByTrashed")
	}

	apiResponse := h.mapping.ToApiResponsePaginationReviewDetailDeleteAt(res)

	h.cache.SetReviewDetailTrashedCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// Create handles the creation of a new merchant review detail.
// @Summary Create a new merchant review detail
// @Tags ReviewDetail
// @Description Create a new merchant review detail with the provided details
// @Accept multipart/form-data
// @Produce json
// @Param type formData string true "Type"
// @Param url formData file true "url"
// @Param caption formData string true "Product name"
// @Success 200 {object} response.ApiResponseReviewDetail "Successfully created review detail"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create review detail"
// @Router /api/review-detail/create [post]
func (h *reviewDetailHandleApi) Create(c echo.Context) error {
	formData, err := h.parseReviewDetailForm(c)
	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	ctx := c.Request().Context()

	req := &pb.CreateReviewDetailRequest{
		ReviewId: int32(formData.ReviewID),
		Type:     formData.Type,
		Url:      formData.Url,
		Caption:  formData.Caption,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	return c.JSON(http.StatusOK, h.mapping.ToApiResponseReviewDetail(res))
}

// @Security Bearer
// Update handles the update of an existing merchant review detail.
// @Summary Update an existing merchant review detail
// @Tags ReviewDetail
// @Description Update an existing merchant review detail with the provided details
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Review Detail ID"
// @Param type formData string true "Type"
// @Param url formData file true "url"
// @Param caption formData string true "Product name"
// @Param request body requests.UpdateReviewDetailRequest true "Update review detail request"
// @Success 200 {object} response.ApiResponseReviewDetail "Successfully updated review detail"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update review detail"
// @Router /api/review-detail/update/{id} [post]
func (h *reviewDetailHandleApi) Update(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	formData, err := h.parseReviewDetailForm(c)
	if err != nil {
		return h.handleGrpcError(err, "Update")
	}

	ctx := c.Request().Context()

	req := &pb.UpdateReviewDetailRequest{
		ReviewDetailId: int32(idInt),
		Type:           formData.Type,
		Url:            formData.Url,
		Caption:        formData.Caption,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "Update")
	}

	return c.JSON(http.StatusOK, h.mapping.ToApiResponseReviewDetail(res))
}

// @Security Bearer
// TrashedMerchant retrieves a trashed merchant record by its ID.
// @Summary Retrieve a trashed merchant
// @Tags ReviewDetail
// @Description Retrieve a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseReviewDetailDeleteAt "Successfully retrieved trashed merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed merchant"
// @Router /api/review-detail/trashed/{id} [get]
func (h *reviewDetailHandleApi) TrashedMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdReviewDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedReviewDetail(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Trashed")
	}

	so := h.mapping.ToApiResponseReviewDetailDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags ReviewDetail
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseReviewDetailDeleteAt "Successfully restored merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/review-detail/restore/{id} [post]
func (h *reviewDetailHandleApi) RestoreMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return h.handleGrpcError(err, "id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdReviewDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreReviewDetail(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Restore")
	}

	so := h.mapping.ToApiResponseReviewDetailDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags ReviewDetail
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/review-detail/delete/{id} [delete]
func (h *reviewDetailHandleApi) DeleteMerchantPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return h.handleGrpcError(err, "id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdReviewDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteReviewDetailPermanent(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Delete Review")
	}

	so := h.mappingReview.ToApiResponseReviewDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags ReviewDetail
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored merchant all"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/review-detail/restore/all [post]
func (h *reviewDetailHandleApi) RestoreAllMerchant(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllReviewDetail(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "RestoreAll")
	}

	so := h.mappingReview.ToApiResponseReviewAll(res)

	h.logger.Debug("Successfully restored all merchant")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags ReviewDetail
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/review-detail/delete/all [post]
func (h *reviewDetailHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllReviewDetailPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "DeleteAll")
	}

	so := h.mappingReview.ToApiResponseReviewAll(res)

	h.logger.Debug("Successfully deleted all merchant permanently")

	return c.JSON(http.StatusOK, so)
}

func (h *reviewDetailHandleApi) parseReviewDetailForm(
	c echo.Context,
) (requests.ReviewDetailFormData, error) {

	var formData requests.ReviewDetailFormData
	var err error

	formData.ReviewID, err = strconv.Atoi(c.FormValue("review_id"))
	if err != nil || formData.ReviewID <= 0 {
		return formData, errors.NewBadRequestError("review_id must be a valid positive integer")
	}

	formData.Type = strings.TrimSpace(c.FormValue("type"))
	if formData.Type == "" {
		return formData, errors.NewBadRequestError("type is required")
	}

	formData.Caption = strings.TrimSpace(c.FormValue("caption"))
	if formData.Caption == "" {
		return formData, errors.NewBadRequestError("caption is required")
	}

	file, err := c.FormFile("url")
	if err != nil {
		return formData, errors.NewBadRequestError("url file is required")
	}

	uploadPath, err := h.upload_image.ProcessImageUpload(c, file)
	if err != nil {
		return formData, err
	}

	formData.Url = uploadPath
	return formData, nil
}

func (h *reviewDetailHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Review Detail").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Review Detail already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Review Detail service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *reviewDetailHandleApi) parseValidationErrors(err error) []errors.ValidationError {
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

func (h *reviewDetailHandleApi) getValidationMessage(fe validator.FieldError) string {
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
