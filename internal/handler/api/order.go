package api

import (
	order_cache "ecommerce/internal/cache/api/order"
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
	"google.golang.org/protobuf/types/known/emptypb"
)

type orderHandleApi struct {
	client     pb.OrderServiceClient
	logger     logger.LoggerInterface
	mapping    response_api.OrderResponseMapper
	apiHandler errors.ApiHandler
	cache      order_cache.OrderMencache
}

func NewHandlerOrder(
	router *echo.Echo,
	client pb.OrderServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.OrderResponseMapper,
	apiHandler errors.ApiHandler,
	cache order_cache.OrderMencache,
) *orderHandleApi {
	orderHandler := &orderHandleApi{
		client:     client,
		logger:     logger,
		mapping:    mapping,
		apiHandler: apiHandler,
		cache:      cache,
	}

	routerOrder := router.Group("/api/order")

	routerOrder.GET("", apiHandler.Handle("findAll", orderHandler.FindAllOrders))
	routerOrder.GET("/:id", apiHandler.Handle("findById", orderHandler.FindById))
	routerOrder.GET("/active", apiHandler.Handle("findByActive", orderHandler.FindByActive))
	routerOrder.GET("/trashed", apiHandler.Handle("findByTrashed", orderHandler.FindByTrashed))

	routerOrder.GET("/monthly-total-revenue", apiHandler.Handle("findMonthlyTotalRevenue", orderHandler.FindMonthlyTotalRevenue))
	routerOrder.GET("/yearly-total-revenue", apiHandler.Handle("findYearlyTotalRevenue", orderHandler.FindYearlyTotalRevenue))
	routerOrder.GET("/merchant/monthly-total-revenue", apiHandler.Handle("findMonthlyTotalRevenueByMerchant", orderHandler.FindMonthlyTotalRevenueByMerchant))
	routerOrder.GET("/merchant/yearly-total-revenue", apiHandler.Handle("findYearlyTotalRevenueByMerchant", orderHandler.FindYearlyTotalRevenueByMerchant))

	routerOrder.GET("/monthly-revenue", apiHandler.Handle("findMonthlyRevenue", orderHandler.FindMonthlyRevenue))
	routerOrder.GET("/yearly-revenue", apiHandler.Handle("findYearlyRevenue", orderHandler.FindYearlyRevenue))
	routerOrder.GET("/merchant/monthly-revenue", apiHandler.Handle("findMonthlyRevenueByMerchant", orderHandler.FindMonthlyRevenueByMerchant))
	routerOrder.GET("/merchant/yearly-revenue", apiHandler.Handle("findYearlyRevenueByMerchant", orderHandler.FindYearlyRevenueByMerchant))

	routerOrder.POST("/create", apiHandler.Handle("create", orderHandler.Create))
	routerOrder.POST("/update/:id", apiHandler.Handle("update", orderHandler.Update))

	routerOrder.POST("/trashed/:id", apiHandler.Handle("trashed", orderHandler.TrashedOrder))
	routerOrder.POST("/restore/:id", apiHandler.Handle("restore", orderHandler.RestoreOrder))
	routerOrder.DELETE("/permanent/:id", apiHandler.Handle("deletePermanent", orderHandler.DeleteOrderPermanent))

	routerOrder.POST("/restore/all", apiHandler.Handle("restoreAll", orderHandler.RestoreAllOrder))
	routerOrder.POST("/permanent/all", apiHandler.Handle("deleteAllPermanent", orderHandler.DeleteAllOrderPermanent))

	return orderHandler
}

// @Security Bearer
// @Summary Find all orders
// @Tags Order
// @Description Retrieve a list of all orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrder "List of orders"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order data"
// @Router /api/order [get]
func (h *orderHandleApi) FindAllOrders(c echo.Context) error {
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

	req := &requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetOrderAllCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllOrderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindAll")
	}

	apiResponse := h.mapping.ToApiResponsePaginationOrder(res)

	h.cache.SetOrderAllCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find order by ID
// @Tags Order
// @Description Retrieve an order by ID
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.ApiResponseOrder "Order data"
// @Failure 400 {object} response.ErrorResponse "Invalid order ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order data"
// @Router /api/order/{id} [get]
func (h *orderHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedOrderCache(ctx, id)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindByIdOrderRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "FindById")
	}

	apiResponse := h.mapping.ToApiResponseOrder(res)

	h.cache.SetCachedOrderCache(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active orders
// @Tags Order
// @Description Retrieve a list of active orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderDeleteAt "List of active orders"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order data"
// @Router /api/order/active [get]
func (h *orderHandleApi) FindByActive(c echo.Context) error {
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

	req := &requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetOrderActiveCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllOrderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByActive")
	}

	apiResponse := h.mapping.ToApiResponsePaginationOrderDeleteAt(res)

	h.cache.SetOrderActiveCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve trashed orders
// @Tags Order
// @Description Retrieve a list of trashed orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderDeleteAt "List of trashed orders"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order data"
// @Router /api/order/trashed [get]
func (h *orderHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetOrderTrashedCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllOrderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByTrashed")
	}

	apiResponse := h.mapping.ToApiResponsePaginationOrderDeleteAt(res)

	h.cache.SetOrderTrashedCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTotalRevenue retrieves monthly revenue statistics
// @Summary Get monthly revenue report
// @Tags Order
// @Security Bearer
// @Description Retrieve monthly revenue statistics for all orders
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseOrderMonthly "Monthly revenue data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/order/monthly-total-revenue [get]
func (h *orderHandleApi) FindMonthlyTotalRevenue(c echo.Context) error {
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

	req := &requests.MonthTotalRevenue{
		Month: month,
		Year:  year,
	}

	cachedData, found := h.cache.GetMonthlyTotalRevenueCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTotalRevenue(ctx, &pb.FindYearMonthTotalRevenue{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyTotalRevenue")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyTotalRevenue(res)

	h.cache.SetMonthlyTotalRevenueCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTotalRevenue retrieves yearly revenue statistics
// @Summary Get yearly revenue report
// @Tags Order
// @Security Bearer
// @Description Retrieve yearly revenue statistics for all orders
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseOrderYearly "Yearly revenue data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/order/yearly-total-revenue [get]
func (h *orderHandleApi) FindYearlyTotalRevenue(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyTotalRevenueCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTotalRevenue(ctx, &pb.FindYearTotalRevenue{
		Year: int32(year),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindYearlyTotalRevenue")
	}

	apiResponse := h.mapping.ToApiResponseYearlyTotalRevenue(res)

	h.cache.SetYearlyTotalRevenueCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTotalRevenueByMerchant retrieves monthly revenue statistics
// @Summary Get monthly revenue report
// @Tags Order
// @Security Bearer
// @Description Retrieve monthly revenue statistics for all orders
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseOrderMonthly "Monthly revenue data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/order/merchant/monthly-total-revenue [get]
func (h *orderHandleApi) FindMonthlyTotalRevenueByMerchant(c echo.Context) error {
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

	req := &requests.MonthTotalRevenueMerchant{
		MerchantID: merchant,
		Month:      month,
		Year:       year,
	}

	cachedData, found := h.cache.GetMonthlyTotalRevenueByMerchantCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTotalRevenueByMerchant(ctx, &pb.FindYearMonthTotalRevenueByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyTotalRevenueByMerchant")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyTotalRevenue(res)

	h.cache.SetMonthlyTotalRevenueByMerchantCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTotalRevenueByMerchant retrieves yearly revenue statistics
// @Summary Get yearly revenue report
// @Tags Order
// @Security Bearer
// @Description Retrieve yearly revenue statistics for all orders
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseOrderYearly "Yearly revenue data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/order/merchant/yearly-total-revenue [get]
func (h *orderHandleApi) FindYearlyTotalRevenueByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantStr := c.QueryParam("merchant_id")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	merchant, err := strconv.Atoi(merchantStr)
	if err != nil {
		return errors.NewBadRequestError("merchant_id is required")
	}

	ctx := c.Request().Context()

	req := &requests.YearTotalRevenueMerchant{
		MerchantID: merchant,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearlyTotalRevenueByMerchantCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTotalRevenueByMerchant(ctx, &pb.FindYearTotalRevenueByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindYearlyTotalRevenueByMerchant")
	}

	apiResponse := h.mapping.ToApiResponseYearlyTotalRevenue(res)

	h.cache.SetYearlyTotalRevenueByMerchantCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyRevenue retrieves monthly revenue statistics
// @Summary Get monthly revenue report
// @Tags Order
// @Security Bearer
// @Description Retrieve monthly revenue statistics for all orders
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseOrderMonthly "Monthly revenue data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/order/monthly-revenue [get]
func (h *orderHandleApi) FindMonthlyRevenue(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetMonthlyOrderCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyRevenue(ctx, &pb.FindYearOrder{
		Year: int32(year),
	})
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyRevenue")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyOrder(res)

	h.cache.SetMonthlyOrderCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyRevenue retrieves yearly revenue statistics
// @Summary Get yearly revenue report
// @Tags Order
// @Security Bearer
// @Description Retrieve yearly revenue statistics for all orders
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseOrderYearly "Yearly revenue data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/order/yearly-revenue [get]
func (h *orderHandleApi) FindYearlyRevenue(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyOrderCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyRevenue(ctx, &pb.FindYearOrder{
		Year: int32(year),
	})
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyRevenue")
	}

	apiResponse := h.mapping.ToApiResponseYearlyOrder(res)

	h.cache.SetYearlyOrderCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyRevenueByMerchant retrieves monthly revenue by merchant
// @Summary Get monthly revenue by merchant
// @Tags Order
// @Security Bearer
// @Description Retrieve monthly revenue statistics for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseOrderMonthly "Monthly revenue by merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/order/merchant/monthly-revenue [get]
func (h *orderHandleApi) FindMonthlyRevenueByMerchant(c echo.Context) error {
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

	req := &requests.MonthOrderMerchant{
		MerchantID: merchant_id,
		Year:       year,
	}

	cachedData, found := h.cache.GetMonthlyOrderByMerchantCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyRevenueByMerchant(ctx, &pb.FindYearOrderByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyRevenueByMerchant")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyOrder(res)

	h.cache.SetMonthlyOrderByMerchantCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyRevenueByMerchant retrieves yearly revenue by merchant
// @Summary Get yearly revenue by merchant
// @Tags Order
// @Security Bearer
// @Description Retrieve yearly revenue statistics for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseOrderYearly "Yearly revenue by merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/order/merchant/yearly-revenue [get]
func (h *orderHandleApi) FindYearlyRevenueByMerchant(c echo.Context) error {
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

	req := &requests.YearOrderMerchant{
		MerchantID: merchant_id,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearlyOrderByMerchantCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyRevenueByMerchant(ctx, &pb.FindYearOrderByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindYearlyRevenueByMerchant")
	}

	apiResponse := h.mapping.ToApiResponseYearlyOrder(res)

	h.cache.SetYearlyOrderByMerchantCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Create a new order
// @Tags Order
// @Description Create a new order with provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateOrderRequest true "Order details"
// @Success 200 {object} response.ApiResponseOrder "Successfully created order"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create order"
// @Router /api/order/create [post]
func (h *orderHandleApi) Create(c echo.Context) error {
	var body requests.CreateOrderRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	grpcReq := &pb.CreateOrderRequest{
		MerchantId: int32(body.MerchantID),
		UserId:     int32(body.UserID),
		TotalPrice: int32(body.TotalPrice),
		Items:      []*pb.CreateOrderItemRequest{},
		Shipping: &pb.CreateShippingAddressRequest{
			Alamat:         body.ShippingAddress.Alamat,
			Provinsi:       body.ShippingAddress.Provinsi,
			Kota:           body.ShippingAddress.Kota,
			Courier:        body.ShippingAddress.Courier,
			ShippingMethod: body.ShippingAddress.ShippingMethod,
			ShippingCost:   int32(body.ShippingAddress.ShippingCost),
			Negara:         body.ShippingAddress.Negara,
		},
	}

	for _, item := range body.Items {
		grpcReq.Items = append(grpcReq.Items, &pb.CreateOrderItemRequest{
			ProductId: int32(item.ProductID),
			Quantity:  int32(item.Quantity),
			Price:     int32(item.Price),
		})
	}

	h.logger.Debug("Creating new order", zap.Any("request", grpcReq))

	res, err := h.client.Create(ctx, grpcReq)

	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	so := h.mapping.ToApiResponseOrder(res)

	h.cache.SetCachedOrderCache(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Update an existing order
// @Tags Order
// @Description Update an existing order with provided details
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param request body requests.UpdateOrderRequest true "Order update details"
// @Success 200 {object} response.ApiResponseOrder "Successfully updated order"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update order"
// @Router /api/order/update [put]
func (h *orderHandleApi) Update(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	var body requests.UpdateOrderRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	grpcReq := &pb.UpdateOrderRequest{
		OrderId:    int32(idInt),
		UserId:     int32(body.UserID),
		TotalPrice: int32(body.TotalPrice),
		Items:      []*pb.UpdateOrderItemRequest{},
		Shipping: &pb.UpdateShippingAddressRequest{
			ShippingId:     int32(*body.ShippingAddress.ShippingID),
			Alamat:         body.ShippingAddress.Alamat,
			Provinsi:       body.ShippingAddress.Provinsi,
			Kota:           body.ShippingAddress.Kota,
			Courier:        body.ShippingAddress.Courier,
			ShippingMethod: body.ShippingAddress.ShippingMethod,
			ShippingCost:   int32(body.ShippingAddress.ShippingCost),
			Negara:         body.ShippingAddress.Negara,
		},
	}

	for _, item := range body.Items {
		grpcReq.Items = append(grpcReq.Items, &pb.UpdateOrderItemRequest{
			OrderItemId: int32(item.OrderItemID),
			ProductId:   int32(item.ProductID),
			Quantity:    int32(item.Quantity),
			Price:       int32(item.Price),
		})
	}

	res, err := h.client.Update(ctx, grpcReq)

	if err != nil {
		return h.handleGrpcError(err, "Update")
	}

	so := h.mapping.ToApiResponseOrder(res)

	h.cache.DeleteOrderCache(ctx, idInt)
	h.cache.SetCachedOrderCache(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedOrder retrieves a trashed order record by its ID.
// @Summary Retrieve a trashed order
// @Tags Order
// @Description Retrieve a trashed order record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.ApiResponseOrderDeleteAt "Successfully retrieved trashed order"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed order"
// @Router /api/order/trashed/{id} [post]
func (h *orderHandleApi) TrashedOrder(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdOrderRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedOrder(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "TrashedOrder")
	}

	so := h.mapping.ToApiResponseOrderDeleteAt(res)

	h.cache.DeleteOrderCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreOrder restores an order record from the trash by its ID.
// @Summary Restore a trashed order
// @Tags Order
// @Description Restore a trashed order record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.ApiResponseOrderDeleteAt "Successfully restored order"
// @Failure 400 {object} response.ErrorResponse "Invalid order ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore order"
// @Router /api/order/restore/{id} [post]
func (h *orderHandleApi) RestoreOrder(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdOrderRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreOrder(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Restore")
	}

	so := h.mapping.ToApiResponseOrderDeleteAt(res)

	h.cache.DeleteOrderCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteOrderPermanent permanently deletes an order record by its ID.
// @Summary Permanently delete an order
// @Tags Order
// @Description Permanently delete an order record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.ApiResponseOrderDelete "Successfully deleted order record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete order:"
// @Router /api/order/delete/{id} [delete]
func (h *orderHandleApi) DeleteOrderPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdOrderRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteOrderPermanent(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "DeletePermanent")
	}

	so := h.mapping.ToApiResponseOrderDelete(res)

	h.cache.DeleteOrderCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllOrder restores all trashed orders.
// @Summary Restore all trashed orders
// @Tags Order
// @Description Restore all trashed order records.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseOrderAll "Successfully restored all orders"
// @Failure 500 {object} response.ErrorResponse "Failed to restore orders"
// @Router /api/order/restore/all [post]
func (h *orderHandleApi) RestoreAllOrder(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllOrder(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "RestoreAll")
	}

	so := h.mapping.ToApiResponseOrderAll(res)

	h.logger.Debug("Successfully restored all orders")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllOrderPermanent permanently deletes all orders.
// @Summary Permanently delete all orders
// @Tags Order
// @Description Permanently delete all order records.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseOrderAll "Successfully deleted all orders permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to delete orders"
// @Router /api/order/delete/all [post]
func (h *orderHandleApi) DeleteAllOrderPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllOrderPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "DeleteAll")
	}

	so := h.mapping.ToApiResponseOrderAll(res)

	h.logger.Debug("Successfully deleted all orders permanently")

	return c.JSON(http.StatusOK, so)
}

func (h *orderHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Order").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Order already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Order service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *orderHandleApi) parseValidationErrors(err error) []errors.ValidationError {
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

func (h *orderHandleApi) getValidationMessage(fe validator.FieldError) string {
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
