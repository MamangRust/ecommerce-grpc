package api

import (
	merchantdetail_cache "ecommerce/internal/cache/api/merchant_detail"
	"ecommerce/internal/domain/requests"
	response_api "ecommerce/internal/mapper"
	"ecommerce/internal/pb"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/upload_image"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantDetailHandleApi struct {
	client          pb.MerchantDetailServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.MerchantDetailResponseMapper
	mappingMerchant response_api.MerchantResponseMapper
	upload_image    upload_image.ImageUploads
	apiHandler      errors.ApiHandler
	cache           merchantdetail_cache.MerchantDetailMencache
}

func NewHandlerMerchantDetail(
	router *echo.Echo,
	client pb.MerchantDetailServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.MerchantDetailResponseMapper,
	mappingMerchant response_api.MerchantResponseMapper,
	upload_image upload_image.ImageUploads,
	apiHandler errors.ApiHandler,
	cache merchantdetail_cache.MerchantDetailMencache,
) *merchantDetailHandleApi {
	merchantDetailHandler := &merchantDetailHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		mappingMerchant: mappingMerchant,
		upload_image:    upload_image,
		apiHandler:      apiHandler,
		cache:           cache,
	}

	routerMerchantDetail := router.Group("/api/merchant-detail")

	routerMerchantDetail.GET("", apiHandler.Handle("findAll", merchantDetailHandler.FindAllMerchantDetail))
	routerMerchantDetail.GET("/:id", apiHandler.Handle("findById", merchantDetailHandler.FindById))
	routerMerchantDetail.GET("/active", apiHandler.Handle("findByActive", merchantDetailHandler.FindByActive))
	routerMerchantDetail.GET("/trashed", apiHandler.Handle("findByTrashed", merchantDetailHandler.FindByTrashed))

	routerMerchantDetail.POST("/create", apiHandler.Handle("create", merchantDetailHandler.Create))
	routerMerchantDetail.POST("/update/:id", apiHandler.Handle("update", merchantDetailHandler.Update))

	routerMerchantDetail.POST("/trashed/:id", apiHandler.Handle("trashed", merchantDetailHandler.TrashedMerchant))
	routerMerchantDetail.POST("/restore/:id", apiHandler.Handle("restore", merchantDetailHandler.RestoreMerchant))
	routerMerchantDetail.DELETE("/permanent/:id", apiHandler.Handle("deletePermanent", merchantDetailHandler.DeleteMerchantPermanent))

	routerMerchantDetail.POST("/restore/all", apiHandler.Handle("restoreAll", merchantDetailHandler.RestoreAllMerchant))
	routerMerchantDetail.POST("/permanent/all", apiHandler.Handle("deleteAllPermanent", merchantDetailHandler.DeleteAllMerchantPermanent))

	return merchantDetailHandler
}

// @Security Bearer
// @Summary Find all merchant
// @Tags MerchantDetail
// @Description Retrieve a list of all merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDetail "List of merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-detail [get]
func (h *merchantDetailHandleApi) FindAllMerchantDetail(c echo.Context) error {
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

	cachedData, found := h.cache.GetCachedMerchantDetailAll(ctx, req)
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

	apiResponse := h.mapping.ToApiResponsePaginationMerchantDetail(res)

	h.cache.SetCachedMerchantDetailAll(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find merchant by ID
// @Tags MerchantDetail
// @Description Retrieve a merchant by ID
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDetail "merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-detail/{id} [get]
func (h *merchantDetailHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedMerchantDetail(ctx, id)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindByIdMerchantDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "FindById")
	}

	apiResponse := h.mapping.ToApiResponseMerchantDetailRelation(res)

	h.cache.SetCachedMerchantDetailRelation(ctx, id, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active merchant
// @Tags MerchantDetail
// @Description Retrieve a list of active merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDetailDeleteAt "List of active merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-detail/active [get]
func (h *merchantDetailHandleApi) FindByActive(c echo.Context) error {
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

	cachedData, found := h.cache.GetCachedMerchantDetailActive(ctx, req)
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

	apiResponse := h.mapping.ToApiResponsePaginationMerchantDetailDeleteAt(res)

	h.cache.SetCachedMerchantDetailActive(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed merchant records.
// @Summary Retrieve trashed merchant
// @Tags MerchantDetail
// @Description Retrieve a list of trashed merchant records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationMerchantDetailDeleteAt "List of trashed merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-detail/trashed [get]
func (h *merchantDetailHandleApi) FindByTrashed(c echo.Context) error {
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

	cachedData, found := h.cache.GetCachedMerchantDetailTrashed(ctx, req)
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

	apiResponse := h.mapping.ToApiResponsePaginationMerchantDetailDeleteAt(res)

	h.cache.SetCachedMerchantDetailTrashed(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Create a new merchant detail
// @Tags MerchantDetail
// @Description Create a new merchant detail with display name, cover image, logo, etc.
// @Accept multipart/form-data
// @Produce json
// @Param merchant_id formData int true "Merchant ID"
// @Param display_name formData string true "Display name"
// @Param short_description formData string true "Short description"
// @Param website_url formData string false "Website URL"
// @Param cover_image_url formData file true "Cover image file"
// @Param logo_url formData file true "Logo file"
// @Param social_links formData string true "Social links in JSON format (e.g., [{\"platform\": \"instagram\", \"url\": \"https://insta...\", \"merchant_detail_id\": 1}])"
// @Success 200 {object} response.ApiResponseMerchantDetail "Successfully created merchant detail"
// @Failure 400 {object} response.ErrorResponse "Invalid request or validation error"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/merchant-detail/create [post]
func (h *merchantDetailHandleApi) Create(c echo.Context) error {
	formData, err := h.parseMerchantDetailCreate(c)
	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	ctx := c.Request().Context()

	var pbSocialLinks []*pb.CreateMerchantSocialRequest
	for _, link := range formData.SocialLinks {
		var detailID int32
		if link.MerchantDetailID != nil {
			detailID = int32(*link.MerchantDetailID)
		}
		pbSocialLinks = append(pbSocialLinks, &pb.CreateMerchantSocialRequest{
			MerchantDetailId: detailID,
			Platform:         link.Platform,
			Url:              link.Url,
		})
	}

	req := &pb.CreateMerchantDetailRequest{
		MerchantId:       int32(formData.MerchantID),
		DisplayName:      formData.DisplayName,
		CoverImageUrl:    formData.CoverImageUrl,
		LogoUrl:          formData.LogoUrl,
		ShortDescription: formData.ShortDescription,
		WebsiteUrl:       formData.WebsiteUrl,
		SocialLinks:      pbSocialLinks,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	so := h.mapping.ToApiResponseMerchantDetail(res)
	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Update existing merchant detail
// @Tags MerchantDetail
// @Description Update an existing merchant detail by ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant Detail ID"
// @Param merchant_id formData int true "Merchant ID"
// @Param display_name formData string true "Display name"
// @Param short_description formData string true "Short description"
// @Param website_url formData string false "Website URL"
// @Param cover_image_url formData file true "Cover image file"
// @Param logo_url formData file true "Logo file"
// @Param social_links formData string true "Social links in JSON format (e.g., [{\"platform\": \"instagram\", \"url\": \"https://insta...\", \"merchant_detail_id\": 1}])"
// @Success 200 {object} response.ApiResponseMerchantDetail "Successfully updated merchant detail"
// @Failure 400 {object} response.ErrorResponse "Invalid request or validation error"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/merchant-detail/update/{id} [post]
func (h *merchantDetailHandleApi) Update(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	formData, err := h.parseMerchantDetailUpdate(c)
	if err != nil {
		return h.handleGrpcError(err, "Update")
	}

	var pbSocialLinks []*pb.UpdateMerchantSocialRequest
	for _, link := range formData.SocialLinks {
		pbSocialLinks = append(pbSocialLinks, &pb.UpdateMerchantSocialRequest{
			Id:               int32(link.ID),
			MerchantDetailId: int32(*link.MerchantDetailID),
			Platform:         link.Platform,
			Url:              link.Url,
		})
	}

	ctx := c.Request().Context()

	req := &pb.UpdateMerchantDetailRequest{
		MerchantDetailId: int32(idInt),
		DisplayName:      formData.DisplayName,
		CoverImageUrl:    formData.CoverImageUrl,
		LogoUrl:          formData.LogoUrl,
		ShortDescription: formData.ShortDescription,
		WebsiteUrl:       formData.WebsiteUrl,
		SocialLinks:      pbSocialLinks,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "Update")
	}

	so := h.mapping.ToApiResponseMerchantDetail(res)
	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedMerchant retrieves a trashed merchant record by its ID.
// @Summary Retrieve a trashed merchant
// @Tags MerchantDetail
// @Description Retrieve a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDetailDeleteAt "Successfully retrieved trashed merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed merchant"
// @Router /api/merchant-detail/trashed/{id} [get]
func (h *merchantDetailHandleApi) TrashedMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedMerchantDetail(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Trashed")
	}

	so := h.mapping.ToApiResponseMerchantDetailDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags MerchantDetail
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDeleteAt "Successfully restored merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant-detail/restore/{id} [post]
func (h *merchantDetailHandleApi) RestoreMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreMerchantDetail(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Restore")
	}

	so := h.mapping.ToApiResponseMerchantDetailDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags MerchantDetail
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/merchant-detail/delete/{id} [delete]
func (h *merchantDetailHandleApi) DeleteMerchantPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteMerchantDetailPermanent(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "DeletePermanent")
	}

	so := h.mappingMerchant.ToApiResponseMerchantDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags MerchantDetail
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored merchant all"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant-detail/restore/all [post]
func (h *merchantDetailHandleApi) RestoreAllMerchant(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllMerchantDetail(ctx, &emptypb.Empty{})

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
// @Tags MerchantDetail
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/merchant-detail/delete/all [post]
func (h *merchantDetailHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllMerchantDetailPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "DeleteAll")
	}

	so := h.mappingMerchant.ToApiResponseMerchantAll(res)

	h.logger.Debug("Successfully deleted all merchant permanently")

	return c.JSON(http.StatusOK, so)
}

func (h *merchantDetailHandleApi) parseMerchantDetailCreate(
	c echo.Context,
) (requests.CreateMerchantDetailFormData, error) {

	var formData requests.CreateMerchantDetailFormData
	var err error

	formData.MerchantID, err = strconv.Atoi(c.FormValue("merchant_id"))
	if err != nil || formData.MerchantID <= 0 {
		return formData, errors.NewBadRequestError("merchant_id must be a valid positive integer")
	}

	formData.DisplayName = strings.TrimSpace(c.FormValue("display_name"))
	if formData.DisplayName == "" {
		return formData, errors.NewBadRequestError("display_name is required")
	}

	formData.ShortDescription = strings.TrimSpace(c.FormValue("short_description"))
	if formData.ShortDescription == "" {
		return formData, errors.NewBadRequestError("short_description is required")
	}

	formData.WebsiteUrl = strings.TrimSpace(c.FormValue("website_url"))

	coverFile, err := c.FormFile("cover_image_url")
	if err != nil {
		return formData, errors.NewBadRequestError("cover_image_url is required")
	}

	coverPath, err := h.upload_image.ProcessImageUpload(c, coverFile)
	if err != nil {
		return formData, err
	}
	formData.CoverImageUrl = coverPath

	logoFile, err := c.FormFile("logo_url")
	if err != nil {
		return formData, errors.NewBadRequestError("logo_url is required")
	}

	logoPath, err := h.upload_image.ProcessImageUpload(c, logoFile)
	if err != nil {
		return formData, err
	}
	formData.LogoUrl = logoPath

	socialLinksJson := strings.TrimSpace(c.FormValue("social_links"))
	if socialLinksJson == "" {
		return formData, errors.NewBadRequestError("social_links is required")
	}

	var parsedSocialLinks []requests.CreateMerchantSocialFormData
	if err := json.Unmarshal([]byte(socialLinksJson), &parsedSocialLinks); err != nil {
		return formData, errors.NewBadRequestError("social_links must be a valid JSON array")
	}
	formData.SocialLinks = parsedSocialLinks

	return formData, nil
}

func (h *merchantDetailHandleApi) parseMerchantDetailUpdate(
	c echo.Context,
) (requests.UpdateMerchantDetailFormData, error) {
	var formData requests.UpdateMerchantDetailFormData

	formData.DisplayName = strings.TrimSpace(c.FormValue("display_name"))
	formData.ShortDescription = strings.TrimSpace(c.FormValue("short_description"))
	formData.WebsiteUrl = strings.TrimSpace(c.FormValue("website_url"))

	coverFile, err := c.FormFile("cover_image_url")
	if err == nil {
		coverPath, err := h.upload_image.ProcessImageUpload(c, coverFile)
		if err != nil {
			return formData, err
		}
		formData.CoverImageUrl = coverPath
	}

	logoFile, err := c.FormFile("logo_url")
	if err == nil {
		logoPath, err := h.upload_image.ProcessImageUpload(c, logoFile)
		if err != nil {
			return formData, err
		}
		formData.LogoUrl = logoPath
	}

	socialLinksRaw := strings.TrimSpace(c.FormValue("social_links"))
	if socialLinksRaw != "" {
		var links []requests.UpdateMerchantSocialFormData
		if err := json.Unmarshal([]byte(socialLinksRaw), &links); err != nil {
			return formData, errors.NewBadRequestError("social_links must be a valid JSON array")
		}
		formData.SocialLinks = links
	}

	return formData, nil
}

func (h *merchantDetailHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Merchant Detaik").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Merchant Detaik already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Merchant Detaik service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *merchantDetailHandleApi) parseValidationErrors(err error) []errors.ValidationError {
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

func (h *merchantDetailHandleApi) getValidationMessage(fe validator.FieldError) string {
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
