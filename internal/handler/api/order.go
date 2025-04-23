package api

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_api "ecommerce/internal/mapper/response/api"
	"ecommerce/internal/pb"
	"ecommerce/pkg/logger"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orderHandleApi struct {
	client  pb.OrderServiceClient
	logger  logger.LoggerInterface
	mapping response_api.OrderResponseMapper
}

func NewHandlerOrder(
	router *echo.Echo,
	client pb.OrderServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.OrderResponseMapper,
) *orderHandleApi {
	orderHandler := &orderHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
	}

	routerOrder := router.Group("/api/order")

	routerOrder.GET("", orderHandler.FindAllOrders)
	routerOrder.GET("/:id", orderHandler.FindById)
	routerOrder.GET("/active", orderHandler.FindByActive)
	routerOrder.GET("/trashed", orderHandler.FindByTrashed)

	routerOrder.GET("/monthly-total-revenue", orderHandler.FindMonthlyTotalRevenue)
	routerOrder.GET("/yearly-total-revenue", orderHandler.FindYearlyTotalRevenue)
	routerOrder.GET("/merchant/monthly-total-revenue", orderHandler.FindMonthlyTotalRevenueByMerchant)
	routerOrder.GET("/merchant/yearly-total-revenue", orderHandler.FindYearlyTotalRevenueByMerchant)

	routerOrder.GET("/monthly-revenue", orderHandler.FindMonthlyRevenue)
	routerOrder.GET("/yearly-revenue", orderHandler.FindYearlyRevenue)
	routerOrder.GET("/merchant/monthly-revenue", orderHandler.FindMonthlyRevenueByMerchant)
	routerOrder.GET("/merchant/yearly-revenue", orderHandler.FindYearlyRevenueByMerchant)

	routerOrder.POST("/create", orderHandler.Create)
	routerOrder.POST("/update/:id", orderHandler.Update)

	routerOrder.POST("/trashed/:id", orderHandler.TrashedOrder)
	routerOrder.POST("/restore/:id", orderHandler.RestoreOrder)
	routerOrder.DELETE("/permanent/:id", orderHandler.DeleteOrderPermanent)

	routerOrder.POST("/restore/all", orderHandler.RestoreAllOrder)
	routerOrder.POST("/permanent/all", orderHandler.DeleteAllOrderPermanent)

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

	req := &pb.FindAllOrderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve order data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve order data",
			Code:    http.StatusInternalServerError,
		})
	}

	so := h.mapping.ToApiResponsePaginationOrder(res)

	return c.JSON(http.StatusOK, so)
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
		h.logger.Debug("Invalid order ID", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid order ID",
			Code:    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdOrderRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve order data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve order data",
			Code:    http.StatusInternalServerError,
		})
	}

	so := h.mapping.ToApiResponseOrder(res)

	return c.JSON(http.StatusOK, so)
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

	req := &pb.FindAllOrderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to retrieve active order data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve active order data",
			Code:    http.StatusInternalServerError,
		})
	}

	so := h.mapping.ToApiResponsePaginationOrderDeleteAt(res)

	return c.JSON(http.StatusOK, so)
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

	req := &pb.FindAllOrderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve trashed order data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve trashed order data",
			Code:    http.StatusInternalServerError,
		})
	}

	so := h.mapping.ToApiResponsePaginationOrderDeleteAt(res)

	return c.JSON(http.StatusOK, so)
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

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))

		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid year parameter",
			Code:    http.StatusBadRequest,
		})
	}

	monthStr := c.QueryParam("month")

	month, err := strconv.Atoi(monthStr)

	if err != nil {
		h.logger.Debug("Invalid month parameter", zap.Error(err))

		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid month parameter",
			Code:    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTotalRevenue(ctx, &pb.FindYearMonthTotalRevenue{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly order revenue", zap.Error(err))

		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve monthly order revenue",
			Code:    http.StatusInternalServerError,
		})
	}

	so := h.mapping.ToApiResponseMonthlyTotalRevenue(res)

	return c.JSON(http.StatusOK, so)
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
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid year parameter",
			Code:    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTotalRevenue(ctx, &pb.FindYearTotalRevenue{
		Year: int32(year),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly order revenue", zap.Error(err))

		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve yearly order revenue",
			Code:    http.StatusInternalServerError,
		})
	}

	so := h.mapping.ToApiResponseYearlyTotalRevenue(res)

	return c.JSON(http.StatusOK, so)
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

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid year parameter",
			Code:    http.StatusBadRequest,
		})
	}

	monthStr := c.QueryParam("month")

	month, err := strconv.Atoi(monthStr)

	if err != nil {
		h.logger.Debug("Invalid month parameter", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid month parameter",
			Code:    http.StatusBadRequest,
		})
	}

	merchantStr := c.QueryParam("merchant_id")

	merchant, err := strconv.Atoi(merchantStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid merchant id parameter",
			Code:    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTotalRevenueByMerchant(ctx, &pb.FindYearMonthTotalRevenueByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly order revenue", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve monthly order revenue",
			Code:    http.StatusInternalServerError,
		})
	}

	so := h.mapping.ToApiResponseMonthlyTotalRevenue(res)

	return c.JSON(http.StatusOK, so)
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

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid year parameter",
			Code:    http.StatusBadRequest,
		})
	}

	merchantStr := c.QueryParam("merchant_id")

	merchant, err := strconv.Atoi(merchantStr)
	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid merchant id parameter",
			Code:    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTotalRevenueByMerchant(ctx, &pb.FindYearTotalRevenueByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly order revenue", zap.Error(err))

		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve yearly order revenue",
			Code:    http.StatusInternalServerError,
		})
	}

	so := h.mapping.ToApiResponseYearlyTotalRevenue(res)

	return c.JSON(http.StatusOK, so)
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
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid year parameter",
			Code:    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyRevenue(ctx, &pb.FindYearOrder{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly order revenue", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve monthly order revenue",
			Code:    http.StatusInternalServerError,
		})
	}

	so := h.mapping.ToApiResponseMonthlyOrder(res)

	return c.JSON(http.StatusOK, so)
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
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid year parameter",
			Code:    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyRevenue(ctx, &pb.FindYearOrder{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly order revenue", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve yearly order revenue",
			Code:    http.StatusInternalServerError,
		})
	}

	so := h.mapping.ToApiResponseYearlyOrder(res)

	return c.JSON(http.StatusOK, so)
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
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid year parameter",
			Code:    http.StatusBadRequest,
		})
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid merchant id parameter",
			Code:    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyRevenueByMerchant(ctx, &pb.FindYearOrderByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly order revenue", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve monthly order revenue",
			Code:    http.StatusInternalServerError,
		})
	}

	so := h.mapping.ToApiResponseMonthlyOrder(res)

	return c.JSON(http.StatusOK, so)
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
	year, err := strconv.Atoi(yearStr)

	merchantIdStr := c.QueryParam("merchant_id")

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid year parameter",
			Code:    http.StatusBadRequest,
		})
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid merchant id parameter",
			Code:    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyRevenueByMerchant(ctx, &pb.FindYearOrderByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly order revenue", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve yearly order revenue",
			Code:    http.StatusBadRequest,
		})
	}

	so := h.mapping.ToApiResponseYearlyOrder(res)

	return c.JSON(http.StatusOK, so)
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
	var req requests.CreateOrderRequest

	if err := c.Bind(&req); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))

		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "invalid_request",
			Message: "Invalid request format. Please check your input.",
			Code:    http.StatusBadRequest,
		})
	}

	if err := req.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))

		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "validation_error",
			Message: "Please provide valid order information.",
			Code:    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()

	grpcReq := &pb.CreateOrderRequest{
		MerchantId: int32(req.MerchantID),
		UserId:     int32(req.UserID),
		TotalPrice: int32(req.TotalPrice),
		Items:      []*pb.CreateOrderItemRequest{},
		Shipping: &pb.CreateShippingAddressRequest{
			Alamat:         req.ShippingAddress.Alamat,
			Provinsi:       req.ShippingAddress.Provinsi,
			Kota:           req.ShippingAddress.Kota,
			Courier:        req.ShippingAddress.Courier,
			ShippingMethod: req.ShippingAddress.ShippingMethod,
			ShippingCost:   int32(req.ShippingAddress.ShippingCost),
			Negara:         req.ShippingAddress.Negara,
		},
	}

	for _, item := range req.Items {
		grpcReq.Items = append(grpcReq.Items, &pb.CreateOrderItemRequest{
			ProductId: int32(item.ProductID),
			Quantity:  int32(item.Quantity),
			Price:     int32(item.Price),
		})
	}

	res, err := h.client.Create(ctx, grpcReq)

	if err != nil {
		h.logger.Error("Cashier creation failed", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "creation_failed",
			Message: "We couldn't create the order. Please try again.",
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, res)
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
		h.logger.Debug("Invalid id parameter", zap.Error(err))

		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid id parameter",
			Code:    http.StatusBadRequest,
		})
	}

	var req requests.UpdateOrderRequest

	if err := c.Bind(&req); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "invalid_request",
			Message: "Invalid request format. Please check your input.",
			Code:    http.StatusBadRequest,
		})
	}

	if err := c.Validate(req); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "validation_error",
			Message: "Please provide valid order information.",
			Code:    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()

	grpcReq := &pb.UpdateOrderRequest{
		OrderId:    int32(idInt),
		TotalPrice: int32(req.TotalPrice),
		Items:      []*pb.UpdateOrderItemRequest{},
		Shipping: &pb.UpdateShippingAddressRequest{
			ShippingId:     int32(*req.ShippingAddress.ShippingID),
			OrderId:        int32(req.ShippingAddress.OrderID),
			Alamat:         req.ShippingAddress.Alamat,
			Provinsi:       req.ShippingAddress.Provinsi,
			Kota:           req.ShippingAddress.Kota,
			Courier:        req.ShippingAddress.Courier,
			ShippingMethod: req.ShippingAddress.ShippingMethod,
			ShippingCost:   int32(req.ShippingAddress.ShippingCost),
			Negara:         req.ShippingAddress.Negara,
		},
	}

	for _, item := range req.Items {
		grpcReq.Items = append(grpcReq.Items, &pb.UpdateOrderItemRequest{
			OrderItemId: int32(item.OrderItemID),
			ProductId:   int32(item.ProductID),
			Quantity:    int32(item.Quantity),
			Price:       int32(item.Price),
		})
	}

	res, err := h.client.Update(ctx, grpcReq)

	if err != nil {
		h.logger.Debug("Failed to update order", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update order",
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, res)
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
		h.logger.Debug("Invalid order ID format", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "invalid_input",
			Message: "Please provide a valid order ID.",
			Code:    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdOrderRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedOrder(ctx, req)

	if err != nil {
		h.logger.Error("Failed to archive order", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "archive_failed",
			Message: "We couldn't archive the order. Please try again.",
			Code:    http.StatusInternalServerError,
		})
	}

	so := h.mapping.ToApiResponseOrderDeleteAt(res)

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
		h.logger.Debug("Invalid order ID format", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "invalid_input",
			Message: "Please provide a valid order ID.",
			Code:    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdOrderRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreOrder(ctx, req)

	if err != nil {
		h.logger.Error("Failed to restore order", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "restore_failed",
			Message: "We couldn't restore the order. Please try again.",
			Code:    http.StatusInternalServerError,
		})
	}

	so := h.mapping.ToApiResponseOrderDeleteAt(res)

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
		h.logger.Debug("Invalid order ID format", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "invalid_input",
			Message: "Please provide a valid order ID.",
			Code:    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdOrderRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteOrderPermanent(ctx, req)

	if err != nil {
		h.logger.Error("Failed to delete order", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "deletion_failed",
			Message: "We couldn't permanently delete the order. Please try again.",
			Code:    http.StatusInternalServerError,
		})
	}

	so := h.mapping.ToApiResponseOrderDelete(res)

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
		h.logger.Error("Bulk orders restoration failed", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "restoration_failed",
			Message: "We couldn't restore all orders. Please try again later.",
			Code:    http.StatusInternalServerError,
		})
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
		h.logger.Error("Bulk order deletion failed", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "deletion_failed",
			Message: "We couldn't permanently delete all orders. Please try again later.",
			Code:    http.StatusInternalServerError,
		})
	}

	so := h.mapping.ToApiResponseOrderAll(res)

	h.logger.Debug("Successfully deleted all orders permanently")

	return c.JSON(http.StatusOK, so)
}
