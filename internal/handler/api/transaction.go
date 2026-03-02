package api

import (
	transaction_cache "ecommerce/internal/cache/api/transaction"
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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transactionHandleApi struct {
	client     pb.TransactionServiceClient
	logger     logger.LoggerInterface
	mapping    response_api.TransactionResponseMapper
	apiHandler errors.ApiHandler
	cache      transaction_cache.TransactionMencache
}

func NewHandlerTransaction(
	router *echo.Echo,
	client pb.TransactionServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.TransactionResponseMapper,
	apiHandler errors.ApiHandler,
	cache transaction_cache.TransactionMencache,
) *transactionHandleApi {
	transactionHandler := &transactionHandleApi{
		client:     client,
		logger:     logger,
		mapping:    mapping,
		apiHandler: apiHandler,
		cache:      cache,
	}

	routerTransaction := router.Group("/api/transaction")

	routerTransaction.GET(
		"",
		apiHandler.Handle("findAll", transactionHandler.FindAllTransaction),
	)
	routerTransaction.GET(
		"/active",
		apiHandler.Handle("findByActive", transactionHandler.FindByActive),
	)
	routerTransaction.GET(
		"/trashed",
		apiHandler.Handle("findByTrashed", transactionHandler.FindByTrashed),
	)

	routerTransaction.GET(
		"/merchant/:merchant_id",
		apiHandler.Handle("findByMerchant", transactionHandler.FindByMerchant),
	)

	routerTransaction.GET(
		"/monthly-success",
		apiHandler.Handle("monthlySuccess", transactionHandler.FindMonthStatusSuccess),
	)
	routerTransaction.GET(
		"/yearly-success",
		apiHandler.Handle("yearlySuccess", transactionHandler.FindYearStatusSuccess),
	)
	routerTransaction.GET(
		"/monthly-failed",
		apiHandler.Handle("monthlyFailed", transactionHandler.FindMonthStatusFailed),
	)
	routerTransaction.GET(
		"/yearly-failed",
		apiHandler.Handle("yearlyFailed", transactionHandler.FindYearStatusFailed),
	)

	routerTransaction.GET(
		"/merchant/monthly-success",
		apiHandler.Handle("merchantMonthlySuccess", transactionHandler.FindMonthStatusSuccessByMerchant),
	)
	routerTransaction.GET(
		"/merchant/yearly-success",
		apiHandler.Handle("merchantYearlySuccess", transactionHandler.FindYearStatusSuccessByMerchant),
	)
	routerTransaction.GET(
		"/merchant/monthly-failed",
		apiHandler.Handle("merchantMonthlyFailed", transactionHandler.FindMonthStatusFailedByMerchant),
	)
	routerTransaction.GET(
		"/merchant/yearly-failed",
		apiHandler.Handle("merchantYearlyFailed", transactionHandler.FindYearStatusFailedByMerchant),
	)

	routerTransaction.GET(
		"/monthly-method-success",
		apiHandler.Handle("monthlyMethodSuccess", transactionHandler.FindMonthMethodSuccess),
	)
	routerTransaction.GET(
		"/yearly-method-success",
		apiHandler.Handle("yearlyMethodSuccess", transactionHandler.FindYearMethodSuccess),
	)
	routerTransaction.GET(
		"/monthly-method-failed",
		apiHandler.Handle("monthlyMethodFailed", transactionHandler.FindMonthMethodFailed),
	)
	routerTransaction.GET(
		"/yearly-method-failed",
		apiHandler.Handle("yearlyMethodFailed", transactionHandler.FindYearMethodFailed),
	)

	routerTransaction.GET(
		"/merchant/monthly-method-success",
		apiHandler.Handle("merchantMonthlyMethodSuccess", transactionHandler.FindMonthMethodByMerchantSuccess),
	)
	routerTransaction.GET(
		"/merchant/yearly-method-success",
		apiHandler.Handle("merchantYearlyMethodSuccess", transactionHandler.FindYearMethodByMerchantSuccess),
	)
	routerTransaction.GET(
		"/merchant/monthly-method-failed",
		apiHandler.Handle("merchantMonthlyMethodFailed", transactionHandler.FindMonthMethodByMerchantFailed),
	)
	routerTransaction.GET(
		"/merchant/yearly-method-failed",
		apiHandler.Handle("merchantYearlyMethodFailed", transactionHandler.FindYearMethodByMerchantFailed),
	)

	routerTransaction.GET(
		"/:id",
		apiHandler.Handle("findById", transactionHandler.FindById),
	)

	routerTransaction.POST(
		"/create",
		apiHandler.Handle("create", transactionHandler.Create),
	)
	routerTransaction.POST(
		"/update/:id",
		apiHandler.Handle("update", transactionHandler.Update),
	)

	routerTransaction.POST(
		"/trashed/:id",
		apiHandler.Handle("trashed", transactionHandler.TrashedTransaction),
	)
	routerTransaction.POST(
		"/restore/:id",
		apiHandler.Handle("restore", transactionHandler.RestoreTransaction),
	)
	routerTransaction.DELETE(
		"/permanent/:id",
		apiHandler.Handle("deletePermanent", transactionHandler.DeleteTransactionPermanent),
	)

	routerTransaction.POST(
		"/restore/all",
		apiHandler.Handle("restoreAll", transactionHandler.RestoreAllTransaction),
	)
	routerTransaction.POST(
		"/permanent/all",
		apiHandler.Handle("deleteAllPermanent", transactionHandler.DeleteAllTransactionPermanent),
	)

	return transactionHandler
}

// @Security Bearer
// @Summary Find all transactions
// @Tags Transaction
// @Description Retrieve a list of all transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction [get]
func (h *transactionHandleApi) FindAllTransaction(c echo.Context) error {
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

	req := &requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedTransactionsCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindAllTransaction")
	}

	apiResponse := h.mapping.ToApiResponsePaginationTransaction(res)

	h.cache.SetCachedTransactionsCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find all transactions by merchant
// @Tags Transaction
// @Description Retrieve a list of all transactions filtered by merchant
// @Accept json
// @Produce json
// @Param merchant_id path int true "Merchant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction/merchant/{merchant_id} [get]
func (h *transactionHandleApi) FindByMerchant(c echo.Context) error {
	merchantID, err := strconv.Atoi(c.Param("merchant_id"))
	if err != nil || merchantID <= 0 {
		return errors.NewBadRequestError("merchant is required")
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

	req := &requests.FindAllTransactionByMerchant{
		MerchantID: merchantID,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	cachedData, found := h.cache.GetCachedTransactionByMerchant(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllTransactionMerchantRequest{
		MerchantId: int32(merchantID),
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindByMerchant(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByMerchant")
	}

	apiResponse := h.mapping.ToApiResponsePaginationTransaction(res)

	h.cache.SetCachedTransactionByMerchant(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find transaction by ID
// @Tags Transaction
// @Description Retrieve a transaction by ID
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransaction "Transaction data"
// @Failure 400 {object} response.ErrorResponse "Invalid transaction ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction/{id} [get]
func (h *transactionHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedTransactionCache(ctx, id)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindByIdTransactionRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "FindById")
	}

	apiResponse := h.mapping.ToApiResponseTransaction(res)

	h.cache.SetCachedTransactionCache(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active transactions
// @Tags Transaction
// @Description Retrieve a list of active transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransactionDeleteAt "List of active transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction/active [get]
func (h *transactionHandleApi) FindByActive(c echo.Context) error {
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

	req := &requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedTransactionActiveCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByActive")
	}

	apiResponse := h.mapping.ToApiResponsePaginationTransactionDeleteAt(res)

	h.cache.SetCachedTransactionActiveCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed transaction records.
// @Summary Retrieve trashed transactions
// @Tags Transaction
// @Description Retrieve a list of trashed transaction records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransactionDeleteAt "List of trashed transaction data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction/trashed [get]
func (h *transactionHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedTransactionTrashedCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByTrashed")
	}

	apiResponse := h.mapping.ToApiResponsePaginationTransactionDeleteAt(res)

	h.cache.SetCachedTransactionTrashedCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthStatusSuccess retrieves monthly successful transactions
// @Summary Get monthly successful transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of successful transactions by month
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month in MM format (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthSuccess
// @Failure 400 {object} response.ErrorResponse "Invalid year or month parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/monthly-success [get]
func (h *transactionHandleApi) FindMonthStatusSuccess(c echo.Context) error {
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

	req := &requests.MonthAmountTransaction{
		Month: month,
		Year:  year,
	}

	cachedData, found := h.cache.GetCachedMonthAmountSuccessCached(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthStatusSuccess(ctx, &pb.FindMonthlyTransactionStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindMonthStatusSuccess")
	}

	apiResponse := h.mapping.ToApiResponseTransactionMonthAmountSuccess(res)

	h.cache.SetCachedMonthAmountSuccessCached(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearStatusSuccess retrieves yearly successful transactions
// @Summary Get yearly successful transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of successful transactions by year
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearSuccess
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/yearly-success [get]
func (h *transactionHandleApi) FindYearStatusSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedYearAmountSuccessCached(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearStatusSuccess(ctx, &pb.FindYearlyTransactionStatus{
		Year: int32(year),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindYearStatusSuccess")
	}

	apiResponse := h.mapping.ToApiResponseTransactionYearAmountSuccess(res)

	h.cache.SetCachedYearAmountSuccessCached(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthStatusFailed retrieves monthly failed transactions
// @Summary Get monthly failed transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of failed transactions by month
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month in MM format (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthFailed
// @Failure 400 {object} response.ErrorResponse "Invalid year or month parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/monthly-failed [get]
func (h *transactionHandleApi) FindMonthStatusFailed(c echo.Context) error {
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

	req := &requests.MonthAmountTransaction{
		Month: month,
		Year:  year,
	}

	cachedData, found := h.cache.GetCachedMonthAmountFailedCached(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthStatusFailed(ctx, &pb.FindMonthlyTransactionStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindMonthStatusFailed")
	}

	apiResponse := h.mapping.ToApiResponseTransactionMonthAmountFailed(res)

	h.cache.SetCachedMonthAmountFailedCached(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearStatusFailed retrieves yearly failed transactions
// @Summary Get yearly failed transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of failed transactions by year
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearFailed
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/yearly-failed [get]
func (h *transactionHandleApi) FindYearStatusFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedYearAmountFailedCached(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearStatusFailed(ctx, &pb.FindYearlyTransactionStatus{
		Year: int32(year),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindYearStatusFailed")
	}

	apiResponse := h.mapping.ToApiResponseTransactionYearAmountFailed(res)

	h.cache.SetCachedYearAmountFailedCached(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthStatusSuccessByMerchant retrieves monthly successful transactions by merchant
// @Summary Get monthly successful transactions by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of successful transactions by month for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month in MM format (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthSuccess
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID, year or month parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/monthly-success [get]
func (h *transactionHandleApi) FindMonthStatusSuccessByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return errors.NewBadRequestError("month is required")
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)
	if err != nil {
		return errors.NewBadRequestError("merchant is required")
	}

	ctx := c.Request().Context()

	req := &requests.MonthAmountTransactionMerchant{
		MerchantID: merchant_id,
		Month:      month,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedMonthAmountSuccessByMerchant(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthStatusSuccessByMerchant(ctx, &pb.FindMonthlyTransactionStatusByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindMonthSuccessByMerchant")
	}

	apiResponse := h.mapping.ToApiResponseTransactionMonthAmountSuccess(res)

	h.cache.SetCachedMonthAmountSuccessByMerchant(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearStatusSuccessByMerchant retrieves yearly successful transactions by merchant
// @Summary Get yearly successful transactions by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of successful transactions by year for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearSuccess
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/yearly-success [get]
func (h *transactionHandleApi) FindYearStatusSuccessByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)
	if err != nil {
		return errors.NewBadRequestError("merchant is required")
	}

	ctx := c.Request().Context()

	req := &requests.YearAmountTransactionMerchant{
		MerchantID: merchant_id,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedYearAmountSuccessByMerchant(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearStatusSuccessByMerchant(ctx, &pb.FindYearlyTransactionStatusByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindYearStatusSuccessByMerchant")
	}

	apiResponse := h.mapping.ToApiResponseTransactionYearAmountSuccess(res)

	h.cache.SetCachedYearAmountSuccessByMerchant(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthStatusFailedByMerchant retrieves monthly failed transactions by merchant
// @Summary Get monthly failed transactions by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of failed transactions by month for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month in MM format (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthFailed
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID, year or month parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/monthly-failed [get]
func (h *transactionHandleApi) FindMonthStatusFailedByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return errors.NewBadRequestError("month is required")
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)
	if err != nil {
		return errors.NewBadRequestError("merchant is required")
	}

	ctx := c.Request().Context()

	req := &requests.MonthAmountTransactionMerchant{
		MerchantID: merchant_id,
		Month:      month,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedMonthAmountFailedByMerchant(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthStatusFailedByMerchant(ctx, &pb.FindMonthlyTransactionStatusByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindMonthStatusFailedByMerchant")
	}

	apiResponse := h.mapping.ToApiResponseTransactionMonthAmountFailed(res)

	h.cache.SetCachedMonthAmountFailedByMerchant(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearStatusFailedByMerchant retrieves yearly failed transactions by merchant
// @Summary Get yearly failed transactions by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of failed transactions by year for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearFailed
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/yearly-failed [get]
func (h *transactionHandleApi) FindYearStatusFailedByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)
	if err != nil {
		return errors.NewBadRequestError("merchant id is required")
	}

	ctx := c.Request().Context()

	req := &requests.YearAmountTransactionMerchant{
		MerchantID: merchant_id,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedYearAmountFailedByMerchant(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearStatusFailedByMerchant(ctx, &pb.FindYearlyTransactionStatusByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindYearStatusFailedByMerchant")
	}

	apiResponse := h.mapping.ToApiResponseTransactionYearAmountFailed(res)

	h.cache.SetCachedYearAmountFailedByMerchant(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthMethod retrieves monthly payment method statistics
// @Summary Get monthly payment method distribution
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by month
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/monthly-method-success [get]
func (h *transactionHandleApi) FindMonthMethodSuccess(c echo.Context) error {
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

	req := &requests.MonthMethodTransaction{
		Month: month,
		Year:  year,
	}

	cachedData, found := h.cache.GetCachedMonthMethodSuccessCached(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthMethodSuccess(ctx, &pb.MonthTransactionMethod{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		return h.handleGrpcError(err, "FindMonthMethodSuccess")
	}

	apiResponse := h.mapping.ToApiResponseTransactionMonthMethod(res)

	h.cache.SetCachedMonthMethodSuccessCached(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearMethod retrieves yearly payment method statistics
// @Summary Get yearly payment method distribution
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by year
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/yearly-method-success [get]
func (h *transactionHandleApi) FindYearMethodSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedYearMethodSuccessCached(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearMethodSuccess(ctx, &pb.YearTransactionMethod{
		Year: int32(year),
	})
	if err != nil {
		return h.handleGrpcError(err, "FindYearMethodSuccess")
	}

	apiResponse := h.mapping.ToApiResponseTransactionYearMethod(res)

	h.cache.SetCachedYearMethodSuccessCached(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthMethodByMerchant retrieves monthly payment method statistics by merchant
// @Summary Get monthly payment method distribution by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by month for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/monthly-method-success/{merchant_id} [get]
func (h *transactionHandleApi) FindMonthMethodByMerchantSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return errors.NewBadRequestError("month is required")
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)
	if err != nil {
		return errors.NewBadRequestError("merchant id is required")
	}

	ctx := c.Request().Context()

	req := &requests.MonthMethodTransactionMerchant{
		MerchantID: merchant_id,
		Month:      month,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedMonthMethodSuccessByMerchant(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthMethodByMerchantSuccess(ctx, &pb.MonthTransactionMethodByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindMonthMethodByMerchantSuccess")
	}

	apiResponse := h.mapping.ToApiResponseTransactionMonthMethod(res)

	h.cache.SetCachedMonthMethodSuccessByMerchant(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearMethodByMerchant retrieves yearly payment method statistics by merchant
// @Summary Get yearly payment method distribution by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by year for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/yearly-method-success/{merchant_id} [get]
func (h *transactionHandleApi) FindYearMethodByMerchantSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)
	if err != nil {
		return errors.NewBadRequestError("merchant id is required")
	}

	ctx := c.Request().Context()

	req := &requests.YearMethodTransactionMerchant{
		MerchantID: merchant_id,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedYearMethodSuccessByMerchant(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearMethodByMerchantSuccess(ctx, &pb.YearTransactionMethodByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})
	if err != nil {
		return h.handleGrpcError(err, "FindYearMethodByMerchantSuccess")
	}

	apiResponse := h.mapping.ToApiResponseTransactionYearMethod(res)

	h.cache.SetCachedYearMethodSuccessByMerchant(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthMethod retrieves monthly payment method statistics
// @Summary Get monthly payment method distribution
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by month
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/monthly-method-failed [get]
func (h *transactionHandleApi) FindMonthMethodFailed(c echo.Context) error {
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

	req := &requests.MonthMethodTransaction{
		Month: month,
		Year:  year,
	}

	cachedData, found := h.cache.GetCachedMonthMethodFailedCached(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthMethodFailed(ctx, &pb.MonthTransactionMethod{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		return h.handleGrpcError(err, "FindMonthMethodFailed")
	}

	apiResponse := h.mapping.ToApiResponseTransactionMonthMethod(res)

	h.cache.SetCachedMonthMethodFailedCached(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearMethod retrieves yearly payment method statistics
// @Summary Get yearly payment method distribution
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by year
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/yearly-method-failed [get]
func (h *transactionHandleApi) FindYearMethodFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedYearMethodFailedCached(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearMethodFailed(ctx, &pb.YearTransactionMethod{
		Year: int32(year),
	})
	if err != nil {
		return h.handleGrpcError(err, "FindYearMethodFailed")
	}

	apiResponse := h.mapping.ToApiResponseTransactionYearMethod(res)

	h.cache.SetCachedYearMethodFailedCached(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthMethodByMerchant retrieves monthly payment method statistics by merchant
// @Summary Get monthly payment method distribution by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by month for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/monthly-method-failed/{merchant_id} [get]
func (h *transactionHandleApi) FindMonthMethodByMerchantFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantIdStr := c.QueryParam("merchant_id")
	monthStr := c.QueryParam("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return errors.NewBadRequestError("month is required")
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)
	if err != nil {
		return errors.NewBadRequestError("merchant id is required")
	}

	ctx := c.Request().Context()

	req := &requests.MonthMethodTransactionMerchant{
		MerchantID: merchant_id,
		Month:      month,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedMonthMethodFailedByMerchant(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthMethodByMerchantFailed(ctx, &pb.MonthTransactionMethodByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		return h.handleGrpcError(err, "FindMonthMethodByMerchantFailed")
	}

	apiResponse := h.mapping.ToApiResponseTransactionMonthMethod(res)

	h.cache.SetCachedMonthMethodFailedByMerchant(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearMethodByMerchant retrieves yearly payment method statistics by merchant
// @Summary Get yearly payment method distribution by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by year for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/yearly-method-failed/{merchant_id} [get]
func (h *transactionHandleApi) FindYearMethodByMerchantFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errors.NewBadRequestError("year is required")
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)
	if err != nil {
		return errors.NewBadRequestError("merchant id is required")
	}

	ctx := c.Request().Context()

	req := &requests.YearMethodTransactionMerchant{
		MerchantID: merchant_id,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedYearMethodFailedByMerchant(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearMethodByMerchantFailed(ctx, &pb.YearTransactionMethodByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})
	if err != nil {
		return h.handleGrpcError(err, "FindYearMethodByMerchantFailed")
	}

	apiResponse := h.mapping.ToApiResponseTransactionYearMethod(res)

	h.cache.SetCachedYearMethodFailedByMerchant(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Create a new transaction
// @Tags Transaction
// @Description Create a new transaction record
// @Accept json
// @Produce json
// @Param request body requests.CreateTransactionRequest true "Transaction details"
// @Success 200 {object} response.ApiResponseTransaction "Successfully created transaction"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create transaction"
// @Router /api/transaction/create [post]
func (h *transactionHandleApi) Create(c echo.Context) error {
	var body requests.CreateTransactionRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()
	grpcReq := &pb.CreateTransactionRequest{
		OrderId:       int32(body.OrderID),
		MerchantId:    int32(body.MerchantID),
		PaymentMethod: body.PaymentMethod,
		Amount:        int32(body.Amount),
	}

	res, err := h.client.Create(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	so := h.mapping.ToApiResponseTransaction(res)

	h.cache.SetCachedTransactionCache(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Update an existing transaction
// @Tags Transaction
// @Description Update an existing transaction record
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param request body requests.UpdateTransactionRequest true "Updated transaction details"
// @Success 200 {object} response.ApiResponseTransaction "Successfully updated transaction"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update transaction"
// @Router /api/transaction/update [post]
func (h *transactionHandleApi) Update(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	var req requests.UpdateTransactionRequest

	if err := c.Bind(&req); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := req.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()
	grpcReq := &pb.UpdateTransactionRequest{
		TransactionId: int32(idInt),
		OrderId:       int32(req.OrderID),
		MerchantId:    int32(req.MerchantID),
		PaymentMethod: req.PaymentMethod,
		Amount:        int32(req.Amount),
	}

	res, err := h.client.Update(ctx, grpcReq)

	if err != nil {
		return h.handleGrpcError(err, "Update")
	}

	so := h.mapping.ToApiResponseTransaction(res)

	h.cache.DeleteTransactionCache(ctx, idInt)
	h.cache.SetCachedTransactionCache(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedTransaction retrieves a trashed transaction record by its ID.
// @Summary Retrieve a trashed transaction
// @Tags Transaction
// @Description Retrieve a trashed transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransactionDeleteAt "Successfully retrieved trashed transaction"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed transaction"
// @Router /api/transaction/trashed/{id} [get]
func (h *transactionHandleApi) TrashedTransaction(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdTransactionRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedTransaction(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Trashed")
	}

	so := h.mapping.ToApiResponseTransactionDeleteAt(res)

	h.cache.DeleteTransactionCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreTransaction restores a transaction record from the trash by its ID.
// @Summary Restore a trashed transaction
// @Tags Transaction
// @Description Restore a trashed transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransactionDeleteAt "Successfully restored transaction"
// @Failure 400 {object} response.ErrorResponse "Invalid transaction ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transaction"
// @Router /api/transaction/restore/{id} [post]
func (h *transactionHandleApi) RestoreTransaction(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdTransactionRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreTransaction(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Restore")
	}

	so := h.mapping.ToApiResponseTransactionDeleteAt(res)

	h.cache.DeleteTransactionCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteTransactionPermanent permanently deletes a transaction record by its ID.
// @Summary Permanently delete a transaction
// @Tags Transaction
// @Description Permanently delete a transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransactionDelete "Successfully deleted transaction record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transaction"
// @Router /api/transaction/delete/{id} [delete]
func (h *transactionHandleApi) DeleteTransactionPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdTransactionRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteTransactionPermanent(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "DeleteTransaction")
	}

	so := h.mapping.ToApiResponseTransactionDelete(res)

	h.cache.DeleteTransactionCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllTransaction restores all trashed transactions.
// @Summary Restore all trashed transactions
// @Tags Transaction
// @Description Restore all trashed transactions.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransactionAll "Successfully restored all transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transactions"
// @Router /api/transaction/restore/all [post]
func (h *transactionHandleApi) RestoreAllTransaction(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllTransaction(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "RestoreAllTransaction")
	}

	so := h.mapping.ToApiResponseTransactionAll(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllTransactionPermanent permanently deletes all transactions.
// @Summary Permanently delete all transactions
// @Tags Transaction
// @Description Permanently delete all transactions.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransactionAll "Successfully deleted all transactions permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transactions"
// @Router /api/transaction/delete/all [post]
func (h *transactionHandleApi) DeleteAllTransactionPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllTransactionPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "DeleteAllTransaction")
	}

	so := h.mapping.ToApiResponseTransactionAll(res)

	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Transaction").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Transaction already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Transaction service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *transactionHandleApi) parseValidationErrors(err error) []errors.ValidationError {
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

func (h *transactionHandleApi) getValidationMessage(fe validator.FieldError) string {
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
