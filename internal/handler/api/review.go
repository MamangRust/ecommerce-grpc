package api

import (
	review_cache "ecommerce/internal/cache/api/review"
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

type reviewHandleApi struct {
	client     pb.ReviewServiceClient
	logger     logger.LoggerInterface
	mapping    response_api.ReviewResponseMapper
	apiHandler errors.ApiHandler
	cache      review_cache.ReviewMencache
}

func NewHandlerReview(
	router *echo.Echo,
	client pb.ReviewServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.ReviewResponseMapper,
	apiHandler errors.ApiHandler,
	cache review_cache.ReviewMencache,
) *reviewHandleApi {
	reviewHandler := &reviewHandleApi{
		client:     client,
		logger:     logger,
		mapping:    mapping,
		apiHandler: apiHandler,
		cache:      cache,
	}

	routerReview := router.Group("/api/review")

	routerReview.GET(
		"",
		apiHandler.Handle("findAll", reviewHandler.FindAll),
	)
	routerReview.GET(
		"/product/:id",
		apiHandler.Handle("findByProduct", reviewHandler.FindByProduct),
	)
	routerReview.GET(
		"/active",
		apiHandler.Handle("findByActive", reviewHandler.FindByActive),
	)
	routerReview.GET(
		"/trashed",
		apiHandler.Handle("findByTrashed", reviewHandler.FindByTrashed),
	)

	routerReview.POST(
		"/create",
		apiHandler.Handle("create", reviewHandler.Create),
	)
	routerReview.POST(
		"/update/:id",
		apiHandler.Handle("update", reviewHandler.Update),
	)

	routerReview.POST(
		"/trashed/:id",
		apiHandler.Handle("trashed", reviewHandler.TrashedReview),
	)
	routerReview.POST(
		"/restore/:id",
		apiHandler.Handle("restore", reviewHandler.RestoreReview),
	)
	routerReview.DELETE(
		"/permanent/:id",
		apiHandler.Handle("deletePermanent", reviewHandler.DeleteReviewPermanent),
	)

	routerReview.POST(
		"/restore/all",
		apiHandler.Handle("restoreAll", reviewHandler.RestoreAllReview),
	)
	routerReview.POST(
		"/permanent/all",
		apiHandler.Handle("deleteAllPermanent", reviewHandler.DeleteAllReviewPermanent),
	)

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
		return h.handleGrpcError(err, "FindAll")
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
		return errors.NewBadRequestError("id is required")
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
		return h.handleGrpcError(err, "FindByProduct")
	}

	so := h.mapping.ToApiResponsePaginationReviewsDetail(res)
	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find reviews by merchant ID
// @Tags Review
// @Description Retrieve a list of reviews for a specific merchant
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReview "List of reviews for the merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve review data"
// @Router /api/review/merchant/{id} [get]
func (h *reviewHandleApi) FindByMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("id is required")
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

	cacheReq := &requests.FindAllReviewByMerchant{
		MerchantID: id,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	cachedData, found := h.cache.GetReviewByMerchantCache(ctx, cacheReq)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllReviewMerchantRequest{
		MerchantId: int32(id),
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindByMerchant(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByMerchant")
	}

	apiResponse := h.mapping.ToApiResponsePaginationReviewsDetail(res)

	h.cache.SetReviewByMerchantCache(ctx, cacheReq, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active review
// @Tags Review
// @Description Retrieve a list of active review
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
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

	req := &requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetReviewActiveCache(ctx, req)
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

	apiResponse := h.mapping.ToApiResponsePaginationReviewDeleteAt(res)

	h.cache.SetReviewActiveCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed review records.
// @Summary Retrieve trashed review
// @Tags Review
// @Description Retrieve a list of trashed review records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
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

	req := &requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetReviewTrashedCache(ctx, req)
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

	apiResponse := h.mapping.ToApiResponsePaginationReviewDeleteAt(res)

	h.cache.SetReviewTrashedCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// Create handles the creation of a new review without image upload.
// @Summary Create a new review
// @Tags Review
// @Description Create a new review with the provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateReviewRequest true "review details"
// @Success 201 {object} response.ApiResponseReview "Successfully created review"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create review"
// @Router /api/review/create [post]
func (h *reviewHandleApi) Create(c echo.Context) error {
	var body requests.CreateReviewRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	grpcReq := &pb.CreateReviewRequest{
		UserId:    int32(body.UserID),
		ProductId: int32(body.ProductID),
		Comment:   body.Comment,
		Rating:    int32(body.Rating),
	}

	res, err := h.client.Create(ctx, grpcReq)

	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	return c.JSON(http.StatusOK, res)
}

// @Security Bearer
// Update handles the update of an existing review.
// @Summary Update an existing review
// @Tags Review
// @Description Update an existing review record with the provided details
// @Accept json
// @Produce json
// @Param id path int true "review ID"
// @Param request body requests.UpdateReviewRequest true "review update details"
// @Success 200 {object} response.ApiResponseReview "Successfully updated review"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update review"
// @Router /api/review/update/{id} [post]
func (h *reviewHandleApi) Update(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	var body requests.UpdateReviewRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	grpcReq := &pb.UpdateReviewRequest{
		ReviewId: int32(idInt),
		Name:     body.Name,
		Comment:  body.Comment,
		Rating:   int32(body.Rating),
	}

	res, err := h.client.Update(ctx, grpcReq)

	if err != nil {
		return h.handleGrpcError(err, "Update")
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
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdReviewRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedReview(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Trashed")
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
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdReviewRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreReview(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "RestoreReview")
	}

	so := h.mapping.ToApiResponseReviewDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteReviewPermanent permanently deletes a review record by its ID.
// @Summary Permanently delete a review
// @Tags review
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
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdReviewRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteReviewPermanent(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "DeleteReviewPermanent")
	}

	so := h.mapping.ToApiResponseReviewDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllReview restores a review record from the trash by its ID.
// @Summary Restore a trashed review
// @Tags review
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
		return h.handleGrpcError(err, "RestoreAll")
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
		return h.handleGrpcError(err, "DeleteAll")
	}

	so := h.mapping.ToApiResponseReviewAll(res)

	h.logger.Debug("Successfully deleted all review permanently")

	return c.JSON(http.StatusOK, so)
}

func (h *reviewHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Review").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Review already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Review service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *reviewHandleApi) parseValidationErrors(err error) []errors.ValidationError {
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

func (h *reviewHandleApi) getValidationMessage(fe validator.FieldError) string {
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
