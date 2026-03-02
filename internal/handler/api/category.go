package api

import (
	category_cache "ecommerce/internal/cache/api/category"
	"ecommerce/internal/domain/requests"
	response_api "ecommerce/internal/mapper"
	"ecommerce/internal/pb"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/upload_image"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type categoryHandleApi struct {
	client       pb.CategoryServiceClient
	logger       logger.LoggerInterface
	mapping      response_api.CategoryResponseMapper
	upload_image upload_image.ImageUploads
	apiHandler   errors.ApiHandler
	cache        category_cache.CategoryMencache
}

func NewHandlerCategory(
	router *echo.Echo,
	client pb.CategoryServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.CategoryResponseMapper,
	upload_image upload_image.ImageUploads,
	apiHandler errors.ApiHandler,
	cache category_cache.CategoryMencache,
) *categoryHandleApi {
	categoryHandler := &categoryHandleApi{
		client:       client,
		logger:       logger,
		mapping:      mapping,
		upload_image: upload_image,
		apiHandler:   apiHandler,
		cache:        cache,
	}

	routerCategory := router.Group("/api/category")

	routerCategory.GET("", apiHandler.Handle("findAll", categoryHandler.FindAllCategory))
	routerCategory.GET("/:id", apiHandler.Handle("findById", categoryHandler.FindById))
	routerCategory.GET("/active", apiHandler.Handle("findByActive", categoryHandler.FindByActive))
	routerCategory.GET("/trashed", apiHandler.Handle("findByTrashed", categoryHandler.FindByTrashed))

	routerCategory.GET("/monthly-total-pricing", apiHandler.Handle("findMonthlyTotalPricing", categoryHandler.FindMonthTotalPrice))
	routerCategory.GET("/yearly-total-pricing", apiHandler.Handle("findYearlyTotalPricing", categoryHandler.FindYearTotalPrice))
	routerCategory.GET("/merchant/monthly-total-pricing", apiHandler.Handle("findMonthlyTotalPricingByMerchant", categoryHandler.FindMonthTotalPriceByMerchant))
	routerCategory.GET("/merchant/yearly-total-pricing", apiHandler.Handle("findYearlyTotalPricingByMerchant", categoryHandler.FindYearTotalPriceByMerchant))
	routerCategory.GET("/mycategory/monthly-total-pricing", apiHandler.Handle("findMonthlyTotalPricingById", categoryHandler.FindMonthTotalPriceById))
	routerCategory.GET("/mycategory/yearly-total-pricing", apiHandler.Handle("findYearlyTotalPricingById", categoryHandler.FindYearTotalPriceById))

	routerCategory.GET("/monthly-pricing", apiHandler.Handle("findMonthlyPricing", categoryHandler.FindMonthPrice))
	routerCategory.GET("/yearly-pricing", apiHandler.Handle("findYearlyPricing", categoryHandler.FindYearPrice))
	routerCategory.GET("/merchant/monthly-pricing", apiHandler.Handle("findMonthlyPricingByMerchant", categoryHandler.FindMonthPriceByMerchant))
	routerCategory.GET("/merchant/yearly-pricing", apiHandler.Handle("findYearlyPricingByMerchant", categoryHandler.FindYearPriceByMerchant))
	routerCategory.GET("/mycategory/monthly-pricing", apiHandler.Handle("findMonthlyPricingById", categoryHandler.FindMonthPriceById))
	routerCategory.GET("/mycategory/yearly-pricing", apiHandler.Handle("findYearlyPricingById", categoryHandler.FindYearPriceById))

	routerCategory.POST("/create", apiHandler.Handle("create", categoryHandler.Create))
	routerCategory.POST("/update/:id", apiHandler.Handle("update", categoryHandler.Update))

	routerCategory.POST("/trashed/:id", apiHandler.Handle("trashed", categoryHandler.TrashedCategory))
	routerCategory.POST("/restore/:id", apiHandler.Handle("restore", categoryHandler.RestoreCategory))
	routerCategory.DELETE("/permanent/:id", apiHandler.Handle("deletePermanent", categoryHandler.DeleteCategoryPermanent))

	routerCategory.POST("/restore/all", apiHandler.Handle("restoreAll", categoryHandler.RestoreAllCategory))
	routerCategory.POST("/permanent/all", apiHandler.Handle("deleteAllPermanent", categoryHandler.DeleteAllCategoryPermanent))

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

	reqCache := &requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedCategoriesCache(ctx, reqCache)

	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindAllCategoryRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "FindAll")
	}

	so := h.mapping.ToApiResponsePaginationCategory(res)

	h.cache.SetCachedCategoriesCache(ctx, reqCache, so)

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
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCategoryRequest{
		Id: int32(id),
	}

	cachedData, found := h.cache.GetCachedCategoryCache(ctx, id)

	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "FindById")
	}

	so := h.mapping.ToApiResponseCategory(res)

	h.cache.SetCachedCategoryCache(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active category
// @Tags Category
// @Description Retrieve a list of active category
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
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

	reqCache := &requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedCategoryActiveCache(ctx, reqCache)

	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindAllCategoryRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "FindByActive")
	}

	so := h.mapping.ToApiResponsePaginationCategoryDeleteAt(res)

	h.cache.SetCachedCategoryActiveCache(ctx, reqCache, so)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed category records.
// @Summary Retrieve trashed category
// @Tags Category
// @Description Retrieve a list of trashed category records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
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

	reqCache := &requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedCategoryTrashedCache(ctx, reqCache)

	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindAllCategoryRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "FindByTrashed")
	}

	so := h.mapping.ToApiResponsePaginationCategoryDeleteAt(res)

	h.cache.SetCachedCategoryTrashedCache(ctx, reqCache, so)

	return c.JSON(http.StatusOK, so)
}

// FindMonthTotalPrice retrieves monthly category pricing statistics
// @Summary Get monthly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/monthly-total-pricing [get]
func (h *categoryHandleApi) FindMonthTotalPrice(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return errors.NewBadRequestError("month is required")
	}

	ctx := c.Request().Context()

	req := &requests.MonthTotalPrice{
		Month: month,
		Year:  year,
	}

	cachedData, found := h.cache.GetCachedMonthTotalPriceCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTotalPrices(ctx, &pb.FindYearMonthTotalPrices{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindMonthTotalPrice")
	}

	apiResponse := h.mapping.ToApiResponseCategoryMonthlyTotalPrice(res)

	h.cache.SetCachedMonthTotalPriceCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearTotalPrice retrieves yearly category pricing statistics
// @Summary Get yearly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/yearly-total-pricing [get]

func (h *categoryHandleApi) FindYearTotalPrice(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedYearTotalPriceCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTotalPrices(ctx, &pb.FindYearTotalPrices{
		Year: int32(year),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindYearTotalPrice")
	}

	apiResponse := h.mapping.ToApiResponseCategoryYearlyTotalPrice(res)

	h.cache.SetCachedYearTotalPriceCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthTotalPriceById retrieves monthly category pricing statistics
// @Summary Get monthly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/mycategory/monthly-total-pricing [get]

func (h *categoryHandleApi) FindMonthTotalPriceById(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	categoryStr := c.QueryParam("category_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return errors.NewBadRequestError("month is required")
	}

	category, err := strconv.Atoi(categoryStr)
	if err != nil {
		return errors.NewBadRequestError("category_id is required")
	}

	ctx := c.Request().Context()

	req := &requests.MonthTotalPriceCategory{
		CategoryID: category,
		Month:      month,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedMonthTotalPriceByIdCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTotalPricesById(ctx, &pb.FindYearMonthTotalPriceById{
		Year:       int32(year),
		Month:      int32(month),
		CategoryId: int32(category),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyTotalPriceById")
	}

	apiResponse := h.mapping.ToApiResponseCategoryMonthlyTotalPrice(res)

	h.cache.SetCachedMonthTotalPriceByIdCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearTotalPriceById retrieves yearly category pricing statistics
// @Summary Get yearly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param category_id path int true "Category ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/yearly-total-pricing/{id} [get]

func (h *categoryHandleApi) FindYearTotalPriceById(c echo.Context) error {
	yearStr := c.QueryParam("year")
	categoryStr := c.QueryParam("category_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	category, err := strconv.Atoi(categoryStr)
	if err != nil {
		return errors.NewBadRequestError("category_id is required")
	}

	ctx := c.Request().Context()

	req := &requests.YearTotalPriceCategory{
		CategoryID: category,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedYearTotalPriceByIdCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTotalPricesById(ctx, &pb.FindYearTotalPriceById{
		Year:       int32(year),
		CategoryId: int32(category),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindYearlyTotalPricesById")
	}

	apiResponse := h.mapping.ToApiResponseCategoryYearlyTotalPrice(res)

	h.cache.SetCachedYearTotalPriceByIdCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthTotalPriceByMerchant retrieves monthly category pricing statistics
// @Summary Get monthly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/merchant/monthly-total-pricing [get]
func (h *categoryHandleApi) FindMonthTotalPriceByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	merchantStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return errors.NewBadRequestError("month is required")
	}

	merchant, err := strconv.Atoi(merchantStr)
	if err != nil {
		return errors.NewBadRequestError("merchant_id is required")
	}

	ctx := c.Request().Context()

	req := &requests.MonthTotalPriceMerchant{
		MerchantID: merchant,
		Month:      month,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedMonthTotalPriceByMerchantCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTotalPricesByMerchant(ctx, &pb.FindYearMonthTotalPriceByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyTotalPriceByMerchant")
	}

	apiResponse := h.mapping.ToApiResponseCategoryMonthlyTotalPrice(res)

	h.cache.SetCachedMonthTotalPriceByMerchantCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearTotalPriceByMerchant retrieves yearly category total pricing statistics
// @Summary Get yearly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/merchant/yearly-total-pricing [get]
func (h *categoryHandleApi) FindYearTotalPriceByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	merchant, err := strconv.Atoi(merchantStr)
	if err != nil {
		return errors.NewBadRequestError("merchant_id required")
	}

	ctx := c.Request().Context()

	req := &requests.YearTotalPriceMerchant{
		MerchantID: merchant,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedYearTotalPriceByMerchantCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTotalPricesByMerchant(ctx, &pb.FindYearTotalPriceByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindYearTotalPriceByMerchant")
	}

	apiResponse := h.mapping.ToApiResponseCategoryYearlyTotalPrice(res)

	h.cache.SetCachedYearTotalPriceByMerchantCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthPrice retrieves monthly category pricing statistics
// @Summary Get monthly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/monthly-pricing [get]
func (h *categoryHandleApi) FindMonthPrice(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedMonthPriceCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthPrice(ctx, &pb.FindYearCategory{
		Year: int32(year),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindMonthPrice")
	}

	apiResponse := h.mapping.ToApiResponseCategoryMonthlyPrice(res)

	h.cache.SetCachedMonthPriceCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearPrice retrieves yearly category pricing statistics
// @Summary Get yearly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/yearly-pricing [get]
func (h *categoryHandleApi) FindYearPrice(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedYearPriceCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearPrice(ctx, &pb.FindYearCategory{
		Year: int32(year),
	})
	if err != nil {
		return h.handleGrpcError(err, "FindYearPrice")
	}

	apiResponse := h.mapping.ToApiResponseCategoryYearlyPrice(res)

	h.cache.SetCachedYearPriceCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthPriceByMerchant retrieves monthly category pricing by merchant
// @Summary Get monthly category pricing by merchant
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for categories by specific merchant
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly category pricing by merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/merchant/monthly-pricing [get]
func (h *categoryHandleApi) FindMonthPriceByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)
	if err != nil {
		return errors.NewBadRequestError("merchant_id is required")
	}

	ctx := c.Request().Context()

	req := &requests.MonthPriceMerchant{
		MerchantID: merchant_id,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedMonthPriceByMerchantCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthPriceByMerchant(ctx, &pb.FindYearCategoryByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})
	if err != nil {
		return h.handleGrpcError(err, "FindMonthPriceByMerchant")
	}

	apiResponse := h.mapping.ToApiResponseCategoryMonthlyPrice(res)

	h.cache.SetCachedMonthPriceByMerchantCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearPriceByMerchant retrieves yearly category pricing by merchant
// @Summary Get yearly category pricing by merchant
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for categories by specific merchant
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly category pricing by merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/merchant/yearly-pricing [get]
func (h *categoryHandleApi) FindYearPriceByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)
	if err != nil {
		return errors.NewBadRequestError("merchant_id is required")
	}

	ctx := c.Request().Context()

	req := &requests.YearPriceMerchant{
		MerchantID: merchant_id,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedYearPriceByMerchantCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearPriceByMerchant(ctx, &pb.FindYearCategoryByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})
	if err != nil {
		return h.handleGrpcError(err, "FindYearPriceByMerchant")
	}

	apiResponse := h.mapping.ToApiResponseCategoryYearlyPrice(res)

	h.cache.SetCachedYearPriceByMerchantCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthPriceById retrieves monthly pricing for specific category
// @Summary Get monthly pricing by category ID
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for specific category
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param category_id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly pricing by category"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Category not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/mycategory/monthly-pricing [get]
func (h *categoryHandleApi) FindMonthPriceById(c echo.Context) error {
	yearStr := c.QueryParam("year")
	categoryIdStr := c.QueryParam("category_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	category_id, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		return errors.NewBadRequestError("category_id is required")
	}

	ctx := c.Request().Context()

	req := &requests.MonthPriceId{
		CategoryID: category_id,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedMonthPriceByIdCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthPriceById(ctx, &pb.FindYearCategoryById{
		Year:       int32(year),
		CategoryId: int32(category_id),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindMonthPriceById")
	}

	apiResponse := h.mapping.ToApiResponseCategoryMonthlyPrice(res)

	h.cache.SetCachedMonthPriceByIdCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearPriceById retrieves yearly pricing for specific category
// @Summary Get yearly pricing by category ID
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for specific category
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param category_id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly pricing by category"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Category not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/mycategory/yearly-pricing [get]
func (h *categoryHandleApi) FindYearPriceById(c echo.Context) error {
	yearStr := c.QueryParam("year")
	categoryIdStr := c.QueryParam("category_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	category_id, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		return errors.NewBadRequestError("category_id is required")
	}

	ctx := c.Request().Context()

	req := &requests.YearPriceId{
		CategoryID: category_id,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedYearPriceByIdCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearPriceById(ctx, &pb.FindYearCategoryById{
		Year:       int32(year),
		CategoryId: int32(category_id),
	})
	if err != nil {
		return h.handleGrpcError(err, "FindYearPriceById")
	}

	apiResponse := h.mapping.ToApiResponseCategoryYearlyPrice(res)

	h.cache.SetCachedYearPriceByIdCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
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
	formData, err := h.parseCategoryForm(c, true)
	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	ctx := c.Request().Context()

	req := &pb.CreateCategoryRequest{
		Name:          formData.Name,
		Description:   formData.Description,
		SlugCategory:  *formData.SlugCategory,
		ImageCategory: formData.ImageCategory,
	}

	res, err := h.client.Create(ctx, req)

	if err != nil {

		if formData.ImageCategory != "" {
			h.upload_image.CleanupImageOnFailure(formData.ImageCategory)
		}

		return h.handleGrpcError(err, "Create")
	}

	so := h.mapping.ToApiResponseCategory(res)

	h.cache.SetCachedCategoryCache(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Update handles the update of an existing category with image upload.
// @Summary Update an existing category
// @Tags Category
// @Description Update an existing category record with the provided details and an optional image file
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Order ID"
// @Param name formData string true "Category name"
// @Param description formData string true "Category description"
// @Param slug_category formData string true "Category slug"
// @Param image_category formData file false "New category image file"
// @Success 200 {object} response.ApiResponseCategory "Successfully updated category"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update category"
// @Router /api/category/update [post]
func (h *categoryHandleApi) Update(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	formData, err := h.parseCategoryForm(c, false)
	if err != nil {
		return h.handleGrpcError(err, "Update")
	}

	ctx := c.Request().Context()

	req := &pb.UpdateCategoryRequest{
		CategoryId:    int32(idInt),
		Name:          formData.Name,
		Description:   formData.Description,
		SlugCategory:  *formData.SlugCategory,
		ImageCategory: formData.ImageCategory,
	}

	res, err := h.client.Update(ctx, req)

	if err != nil {
		if formData.ImageCategory != "" {
			h.upload_image.CleanupImageOnFailure(formData.ImageCategory)
		}

		h.logger.Error("Category update failed",
			zap.Error(err),
			zap.Any("request", req),
		)

		return h.handleGrpcError(err, "Update")
	}

	so := h.mapping.ToApiResponseCategory(res)

	h.cache.DeleteCachedCategoryCache(ctx, idInt)

	h.cache.SetCachedCategoryCache(ctx, so)

	return c.JSON(http.StatusOK, so)
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
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCategoryRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedCategory(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Trashed")
	}

	so := h.mapping.ToApiResponseCategoryDeleteAt(res)

	h.cache.DeleteCachedCategoryCache(ctx, id)

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
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCategoryRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreCategory(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Restore")
	}

	so := h.mapping.ToApiResponseCategoryDeleteAt(res)

	h.cache.DeleteCachedCategoryCache(ctx, id)

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
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCategoryRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteCategoryPermanent(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "DeletePermanent")
	}

	so := h.mapping.ToApiResponseCategoryDelete(res)

	h.cache.DeleteCachedCategoryCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllCategory restores a category record from the trash by its ID.
// @Summary Restore a trashed category
// @Tags Category
// @Description Restore a trashed category record by its ID.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCategoryAll "Successfully restored category all"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore category"
// @Router /api/category/restore/all [post]
func (h *categoryHandleApi) RestoreAllCategory(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllCategory(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "RestoreAll")
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
// @Success 200 {object} response.ApiResponseCategoryAll "Successfully deleted category record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete category:"
// @Router /api/category/delete/all [post]
func (h *categoryHandleApi) DeleteAllCategoryPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllCategoryPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "DeleteAll")
	}

	so := h.mapping.ToApiResponseCategoryAll(res)

	h.logger.Debug("Successfully deleted all category permanently")

	return c.JSON(http.StatusOK, so)
}

func (h *categoryHandleApi) parseCategoryForm(
	c echo.Context,
	requireImage bool,
) (requests.CategoryFormData, error) {

	var formData requests.CategoryFormData

	formData.Name = strings.TrimSpace(c.FormValue("name"))
	if formData.Name == "" {
		return formData, errors.NewBadRequestError("category name is required")
	}

	formData.Description = strings.TrimSpace(c.FormValue("description"))
	if formData.Description == "" {
		return formData, errors.NewBadRequestError("category description is required")
	}

	slug := strings.TrimSpace(c.FormValue("slug_category"))
	if slug == "" {
		return formData, errors.NewBadRequestError("slug_category is required")
	}
	formData.SlugCategory = &slug

	file, err := c.FormFile("image_category")
	if err != nil {
		if requireImage {
			return formData, errors.NewBadRequestError("image_category is required")
		}

		return formData, nil
	}

	imagePath, err := h.upload_image.ProcessImageUpload(c, file)
	if err != nil {
		return formData, err
	}

	formData.ImageCategory = imagePath
	return formData, nil
}

func (h *categoryHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Category").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Category already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Category service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *categoryHandleApi) parseValidationErrors(err error) []errors.ValidationError {
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

func (h *categoryHandleApi) getValidationMessage(fe validator.FieldError) string {
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
