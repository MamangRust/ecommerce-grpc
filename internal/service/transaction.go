package service

import (
	"context"
	transaction_cache "ecommerce/internal/cache/transaction"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	merchant_errors "ecommerce/pkg/errors/merchant"
	"ecommerce/pkg/errors/order_errors"
	orderitem_errors "ecommerce/pkg/errors/order_item_errors"
	shippingaddress_errors "ecommerce/pkg/errors/shipping_address_errors"
	"ecommerce/pkg/errors/transaction_errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type transactionService struct {
	merchantRepository    repository.MerchantRepository
	transactionRepository repository.TransactionRepository
	orderRepository       repository.OrderRepository
	orderItemRepository   repository.OrderItemRepository
	shippingRepository    repository.ShippingAddressRepository
	logger                logger.LoggerInterface
	cache                 transaction_cache.TransactionMencache
	observability         observability.TraceLoggerObservability
}

type TransactionServiceDeps struct {
	MerchantRepository    repository.MerchantRepository
	TransactionRepository repository.TransactionRepository
	OrderRepository       repository.OrderRepository
	OrderItemRepository   repository.OrderItemRepository
	ShippingRepository    repository.ShippingAddressRepository
	Logger                logger.LoggerInterface
	Cache                 transaction_cache.TransactionMencache
	Observability         observability.TraceLoggerObservability
}

func NewTransactionService(deps TransactionServiceDeps) *transactionService {
	return &transactionService{
		merchantRepository:    deps.MerchantRepository,
		transactionRepository: deps.TransactionRepository,
		orderRepository:       deps.OrderRepository,
		orderItemRepository:   deps.OrderItemRepository,
		shippingRepository:    deps.ShippingRepository,
		logger:                deps.Logger,
		cache:                 deps.Cache,
		observability:         deps.Observability,
	}
}

func (s *transactionService) FindAllTransactions(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsRow, *int, error) {
	const method = "FindAllTransactions"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedTransactionsCache(ctx, req); found {
		logSuccess("Successfully retrieved transactions from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, err := s.transactionRepository.FindAllTransactions(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTransactionsRow](
			s.logger,
			transaction_errors.ErrFailedFindAllTransactions,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
		)
	}

	var totalCount int

	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedTransactionsCache(ctx, req, transactions, &totalCount)

	logSuccess("Successfully fetched transactions from repository",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return transactions, &totalCount, nil
}

func (s *transactionService) FindByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*db.GetTransactionByMerchantRow, *int, error) {
	const method = "FindByMerchant"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search
	merchant_id := req.MerchantID

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
		attribute.Int("merchant_id", merchant_id))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedTransactionByMerchant(ctx, req); found {
		logSuccess("Successfully retrieved merchant transactions from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, err := s.transactionRepository.FindByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTransactionByMerchantRow](
			s.logger,
			transaction_errors.ErrFailedFindTransactionsByMerchant,
			method,
			span,

			zap.Int("merchant_id", merchant_id),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
		)
	}

	var totalCount int

	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedTransactionByMerchant(ctx, req, transactions, &totalCount)

	logSuccess("Successfully fetched merchant transactions from repository",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return transactions, &totalCount, nil
}

func (s *transactionService) FindByActive(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsActiveRow, *int, error) {
	const method = "FindByActive"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedTransactionActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active transactions from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, err := s.transactionRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTransactionsActiveRow](
			s.logger,
			transaction_errors.ErrFailedFindTransactionsByActive,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
		)
	}

	var totalCount int

	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedTransactionActiveCache(ctx, req, transactions, &totalCount)

	logSuccess("Successfully fetched active transactions from repository",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return transactions, &totalCount, nil
}

func (s *transactionService) FindByTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsTrashedRow, *int, error) {
	const method = "FindByTrashed"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedTransactionTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved trashed transactions from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, err := s.transactionRepository.FindByTrashed(ctx, req)

	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTransactionsTrashedRow](
			s.logger,
			transaction_errors.ErrFailedFindTransactionsByTrashed,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
		)
	}

	var totalCount int

	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedTransactionTrashedCache(ctx, req, transactions, &totalCount)

	logSuccess("Successfully fetched trashed transactions from repository",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return transactions, &totalCount, nil
}

func (s *transactionService) FindById(ctx context.Context, transactionID int) (*db.GetTransactionByIDRow, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transactionID", transactionID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedTransactionCache(ctx, transactionID); found {
		logSuccess("Successfully retrieved transaction from cache",
			zap.Int("transactionID", transactionID))
		return data, nil
	}

	transaction, err := s.transactionRepository.FindById(ctx, transactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetTransactionByIDRow](
			s.logger,
			transaction_errors.ErrFailedFindTransactionById,
			method,
			span,

			zap.Int("transaction_id", transactionID),
		)
	}

	s.cache.SetCachedTransactionCache(ctx, transaction)

	logSuccess("Successfully fetched transaction from repository",
		zap.Int("transactionID", transactionID))

	return transaction, nil
}

func (s *transactionService) FindByOrderId(ctx context.Context, orderID int) (*db.GetTransactionByOrderIDRow, error) {
	const method = "FindByOrderId"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderID", orderID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedTransactionByOrderId(ctx, orderID); found {
		logSuccess("Successfully retrieved transaction by order ID from cache",
			zap.Int("orderID", orderID))
		return data, nil
	}

	transaction, err := s.transactionRepository.FindByOrderId(ctx, orderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetTransactionByOrderIDRow](
			s.logger,
			transaction_errors.ErrFailedFindTransactionByOrderId,
			method,
			span,

			zap.Int("order_id", orderID),
		)
	}

	s.cache.SetCachedTransactionByOrderId(ctx, orderID, transaction)

	logSuccess("Successfully fetched transaction by order ID from repository",
		zap.Int("orderID", orderID))

	return transaction, nil
}

func (s *transactionService) FindMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionSuccessRow, error) {
	const method = "FindMonthlyAmountSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthAmountSuccessCached(ctx, req); found {
		logSuccess("Successfully retrieved monthly amount success from cache",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month))
		return data, nil
	}

	res, err := s.transactionRepository.GetMonthlyAmountSuccess(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyAmountTransactionSuccessRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyAmountSuccess,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetCachedMonthAmountSuccessCached(ctx, req, res)

	logSuccess("Successfully fetched monthly amount success from repository",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	return res, nil
}

func (s *transactionService) FindYearlyAmountSuccess(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionSuccessRow, error) {
	const method = "FindYearlyAmountSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearAmountSuccessCached(ctx, year); found {
		logSuccess("Successfully retrieved yearly amount success from cache",
			zap.Int("year", year))
		return data, nil
	}

	res, err := s.transactionRepository.GetYearlyAmountSuccess(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyAmountTransactionSuccessRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyAmountSuccess,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetCachedYearAmountSuccessCached(ctx, year, res)

	logSuccess("Successfully fetched yearly amount success from repository",
		zap.Int("year", year))

	return res, nil
}

func (s *transactionService) FindMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionFailedRow, error) {
	const method = "FindMonthlyAmountFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	// Check cache first
	if data, found := s.cache.GetCachedMonthAmountFailedCached(ctx, req); found {
		logSuccess("Successfully retrieved monthly amount failed from cache",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month))
		return data, nil
	}

	res, err := s.transactionRepository.GetMonthlyAmountFailed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyAmountTransactionFailedRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyAmountFailed,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetCachedMonthAmountFailedCached(ctx, req, res)

	logSuccess("Successfully fetched monthly amount failed from repository",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	return res, nil
}

func (s *transactionService) FindYearlyAmountFailed(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionFailedRow, error) {
	const method = "FindYearlyAmountFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearAmountFailedCached(ctx, year); found {
		logSuccess("Successfully retrieved yearly amount failed from cache",
			zap.Int("year", year))
		return data, nil
	}

	res, err := s.transactionRepository.GetYearlyAmountFailed(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyAmountTransactionFailedRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyAmountFailed,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetCachedYearAmountFailedCached(ctx, year, res)

	logSuccess("Successfully fetched yearly amount failed from repository",
		zap.Int("year", year))

	return res, nil
}

func (s *transactionService) FindMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionSuccessByMerchantRow, error) {
	const method = "FindMonthlyAmountSuccessByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthAmountSuccessByMerchant(ctx, req); found {
		logSuccess("Successfully retrieved monthly amount success by merchant from cache",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.Int("merchantID", req.MerchantID))
		return data, nil
	}

	res, err := s.transactionRepository.GetMonthlyAmountSuccessByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyAmountTransactionSuccessByMerchantRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyAmountSuccessByMerchant,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	s.cache.SetCachedMonthAmountSuccessByMerchant(ctx, req, res)

	logSuccess("Successfully fetched monthly amount success by merchant from repository",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month),
		zap.Int("merchantID", req.MerchantID))

	return res, nil
}

func (s *transactionService) FindYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionSuccessByMerchantRow, error) {
	const method = "FindYearlyAmountSuccessByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearAmountSuccessByMerchant(ctx, req); found {
		logSuccess("Successfully retrieved yearly amount success by merchant from cache",
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID))
		return data, nil
	}

	res, err := s.transactionRepository.GetYearlyAmountSuccessByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyAmountTransactionSuccessByMerchantRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyAmountSuccessByMerchant,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	s.cache.SetCachedYearAmountSuccessByMerchant(ctx, req, res)

	logSuccess("Successfully fetched yearly amount success by merchant from repository",
		zap.Int("year", req.Year),
		zap.Int("merchantID", req.MerchantID))

	return res, nil
}

func (s *transactionService) FindMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionFailedByMerchantRow, error) {
	const method = "FindMonthlyAmountFailedByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthAmountFailedByMerchant(ctx, req); found {
		logSuccess("Successfully retrieved monthly amount failed by merchant from cache",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.Int("merchantID", req.MerchantID))
		return data, nil
	}

	res, err := s.transactionRepository.GetMonthlyAmountFailedByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyAmountTransactionFailedByMerchantRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyAmountFailedByMerchant,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	s.cache.SetCachedMonthAmountFailedByMerchant(ctx, req, res)

	logSuccess("Successfully fetched monthly amount failed by merchant from repository",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month),
		zap.Int("merchantID", req.MerchantID))

	return res, nil
}

func (s *transactionService) FindYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionFailedByMerchantRow, error) {
	const method = "FindYearlyAmountFailedByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearAmountFailedByMerchant(ctx, req); found {
		logSuccess("Successfully retrieved yearly amount failed by merchant from cache",
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID))
		return data, nil
	}

	res, err := s.transactionRepository.GetYearlyAmountFailedByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyAmountTransactionFailedByMerchantRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyAmountFailedByMerchant,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	s.cache.SetCachedYearAmountFailedByMerchant(ctx, req, res)

	logSuccess("Successfully fetched yearly amount failed by merchant from repository",
		zap.Int("year", req.Year),
		zap.Int("merchantID", req.MerchantID))

	return res, nil
}

func (s *transactionService) FindMonthlyTransactionMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsSuccessRow, error) {
	const method = "FindMonthlyTransactionMethodSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthMethodSuccessCached(ctx, req); found {
		logSuccess("Successfully retrieved monthly method success from cache",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month))
		return data, nil
	}

	res, err := s.transactionRepository.GetMonthlyTransactionMethodSuccess(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransactionMethodsSuccessRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyMethod,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetCachedMonthMethodSuccessCached(ctx, req, res)

	logSuccess("Successfully fetched monthly method success from repository",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	return res, nil
}

func (s *transactionService) FindYearlyTransactionMethodSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsSuccessRow, error) {
	const method = "FindYearlyTransactionMethodSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearMethodSuccessCached(ctx, year); found {
		logSuccess("Successfully retrieved yearly method success from cache",
			zap.Int("year", year))
		return data, nil
	}

	res, err := s.transactionRepository.GetYearlyTransactionMethodSuccess(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransactionMethodsSuccessRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyMethod,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetCachedYearMethodSuccessCached(ctx, year, res)

	logSuccess("Successfully fetched yearly method success from repository",
		zap.Int("year", year))

	return res, nil
}

func (s *transactionService) FindMonthlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow, error) {
	const method = "FindMonthlyTransactionMethodByMerchantSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthMethodSuccessByMerchant(ctx, req); found {
		logSuccess("Successfully retrieved monthly method success by merchant from cache",
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID))
		return data, nil
	}

	res, err := s.transactionRepository.GetMonthlyTransactionMethodByMerchantSuccess(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyMethodByMerchant,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	s.cache.SetCachedMonthMethodSuccessByMerchant(ctx, req, res)

	logSuccess("Successfully fetched monthly method success by merchant from repository",
		zap.Int("year", req.Year),
		zap.Int("merchantID", req.MerchantID))

	return res, nil
}

func (s *transactionService) FindYearlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantSuccessRow, error) {
	const method = "FindYearlyTransactionMethodByMerchantSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearMethodSuccessByMerchant(ctx, req); found {
		logSuccess("Successfully retrieved yearly method success by merchant from cache",
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID))
		return data, nil
	}

	res, err := s.transactionRepository.GetYearlyTransactionMethodByMerchantSuccess(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransactionMethodsByMerchantSuccessRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyMethodByMerchant,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	s.cache.SetCachedYearMethodSuccessByMerchant(ctx, req, res)

	logSuccess("Successfully fetched yearly method success by merchant from repository",
		zap.Int("year", req.Year),
		zap.Int("merchantID", req.MerchantID))

	return res, nil
}

func (s *transactionService) FindMonthlyTransactionMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsFailedRow, error) {
	const method = "FindMonthlyTransactionMethodFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthMethodFailedCached(ctx, req); found {
		logSuccess("Successfully retrieved monthly method failed from cache",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month))
		return data, nil
	}

	res, err := s.transactionRepository.GetMonthlyTransactionMethodFailed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransactionMethodsFailedRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyMethod,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetCachedMonthMethodFailedCached(ctx, req, res)

	logSuccess("Successfully fetched monthly method failed from repository",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	return res, nil
}

func (s *transactionService) FindYearlyTransactionMethodFailed(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsFailedRow, error) {
	const method = "FindYearlyTransactionMethodFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearMethodFailedCached(ctx, year); found {
		logSuccess("Successfully retrieved yearly method failed from cache",
			zap.Int("year", year))
		return data, nil
	}

	res, err := s.transactionRepository.GetYearlyTransactionMethodFailed(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransactionMethodsFailedRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyMethod,
			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetCachedYearMethodFailedCached(ctx, year, res)

	logSuccess("Successfully fetched yearly method failed from repository",
		zap.Int("year", year))

	return res, nil
}

func (s *transactionService) FindMonthlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantFailedRow, error) {
	const method = "FindMonthlyTransactionMethodByMerchantFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthMethodFailedByMerchant(ctx, req); found {
		logSuccess("Successfully retrieved monthly method failed by merchant from cache",
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID))
		return data, nil
	}

	res, err := s.transactionRepository.GetMonthlyTransactionMethodByMerchantFailed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransactionMethodsByMerchantFailedRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyMethodByMerchant,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	s.cache.SetCachedMonthMethodFailedByMerchant(ctx, req, res)

	logSuccess("Successfully fetched monthly method failed by merchant from repository",
		zap.Int("year", req.Year),
		zap.Int("merchantID", req.MerchantID))

	return res, nil
}

func (s *transactionService) FindYearlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantFailedRow, error) {
	const method = "FindYearlyTransactionMethodByMerchantFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearMethodFailedByMerchant(ctx, req); found {
		logSuccess("Successfully retrieved yearly method failed by merchant from cache",
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID))
		return data, nil
	}

	res, err := s.transactionRepository.GetYearlyTransactionMethodByMerchantFailed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransactionMethodsByMerchantFailedRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyMethodByMerchant,
			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	s.cache.SetCachedYearMethodFailedByMerchant(ctx, req, res)

	logSuccess("Successfully fetched yearly method failed by merchant from repository",
		zap.Int("year", req.Year),
		zap.Int("merchantID", req.MerchantID))

	return res, nil
}

func (s *transactionService) CreateTransaction(ctx context.Context, req *requests.CreateTransactionRequest) (*db.CreateTransactionRow, error) {
	const method = "CreateTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderID", req.OrderID),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	_, err := s.merchantRepository.FindById(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById,
			method,
			span,
			zap.Int("merchantId", req.MerchantID),
		)
	}

	_, err = s.orderRepository.FindById(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](
			s.logger,
			order_errors.ErrFailedFindOrderById,
			method,
			span,
			zap.Int("orderID", req.OrderID),
		)
	}

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(ctx, req.OrderID)
	if err != nil || len(orderItems) == 0 {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemByOrder,
			method,
			span,
			zap.Int("orderID", req.OrderID),
		)
	}

	for _, item := range orderItems {
		if item.Quantity <= 0 {
			status = "error"
			return errorhandler.HandleError[*db.CreateTransactionRow](
				s.logger,
				orderitem_errors.ErrFailedFindOrderItemByOrder,
				method,
				span,
				zap.Int("orderItemID", int(item.OrderItemID)),
				zap.Int("quantity", int(item.Quantity)),
			)
		}
	}

	shipping, err := s.shippingRepository.FindByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](
			s.logger,
			shippingaddress_errors.ErrFailedFindShippingAddressByOrder,
			method,
			span,
			zap.Int("orderID", req.OrderID),
		)
	}

	var totalAmount int
	for _, item := range orderItems {
		totalAmount += int(item.Price) * int(item.Quantity)
	}

	totalAmount += int(shipping.ShippingCost)

	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn

	var paymentStatus string
	if req.Amount >= totalAmountWithTax {
		paymentStatus = "success"
	} else {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](
			s.logger,
			transaction_errors.ErrFailedPaymentInsufficientBalance,
			method,
			span,
			zap.Int("requiredAmount", totalAmountWithTax),
			zap.Int("providedAmount", req.Amount),
		)
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionRepository.CreateTransaction(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](
			s.logger,
			transaction_errors.ErrFailedCreateTransaction,
			method,
			span,
		)
	}

	s.cache.InvalidateTransactionCache(ctx)

	logSuccess("Successfully created transaction",
		zap.Int("transactionID", int(transaction.TransactionID)),
		zap.Int("orderID", req.OrderID),
		zap.String("paymentStatus", paymentStatus),
		zap.Int("totalAmount", totalAmountWithTax))

	return transaction, nil
}

func (s *transactionService) UpdateTransaction(ctx context.Context, req *requests.UpdateTransactionRequest) (*db.UpdateTransactionRow, error) {
	const method = "UpdateTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transactionID", *req.TransactionID),
		attribute.Int("merchantID", req.MerchantID),
		attribute.Int("orderID", req.OrderID))

	defer func() {
		end(status)
	}()

	existingTx, err := s.transactionRepository.FindById(ctx, *req.TransactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			transaction_errors.ErrFailedFindTransactionById,
			method,
			span,
			zap.Int("transactionID", *req.TransactionID),
		)
	}

	if existingTx.PaymentStatus == "success" || existingTx.PaymentStatus == "refunded" {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			transaction_errors.ErrFailedPaymentStatusCannotBeModified,
			method,
			span,
			zap.Int("transactionID", *req.TransactionID),
			zap.String("paymentStatus", existingTx.PaymentStatus),
		)
	}

	_, err = s.merchantRepository.FindById(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById,
			method,
			span,
			zap.Int("merchantId", req.MerchantID),
		)
	}

	_, err = s.orderRepository.FindById(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			order_errors.ErrFailedFindOrderById,
			method,
			span,
			zap.Int("orderID", req.OrderID),
		)
	}

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemByOrder,
			method,
			span,
			zap.Int("orderID", req.OrderID),
		)
	}

	for _, item := range orderItems {
		if item.Quantity <= 0 {
			status = "error"
			return errorhandler.HandleError[*db.UpdateTransactionRow](
				s.logger,
				orderitem_errors.ErrFailedFindOrderItemByOrder,
				method,
				span,
				zap.Int("orderItemID", int(item.OrderItemID)),
				zap.Int("quantity", int(item.Quantity)),
			)
		}
	}

	shipping, err := s.shippingRepository.FindByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			shippingaddress_errors.ErrFailedFindShippingAddressByOrder,
			method,
			span,
			zap.Int("orderID", req.OrderID),
		)
	}

	var totalAmount int
	for _, item := range orderItems {
		totalAmount += int(item.Price) * int(item.Quantity)
	}

	totalAmount += int(shipping.ShippingCost)

	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn

	var paymentStatus string
	if req.Amount >= totalAmountWithTax {
		paymentStatus = "success"
	} else {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			transaction_errors.ErrFailedPaymentInsufficientBalance,
			method,
			span,
			zap.Int("requiredAmount", totalAmountWithTax),
			zap.Int("providedAmount", req.Amount),
		)
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionRepository.UpdateTransaction(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](
			s.logger,
			transaction_errors.ErrFailedUpdateTransaction,
			method,
			span,
		)
	}

	s.cache.InvalidateTransactionCache(ctx)

	logSuccess("Successfully updated transaction",
		zap.Int("transactionID", *req.TransactionID),
		zap.Int("orderID", req.OrderID),
		zap.String("paymentStatus", paymentStatus),
		zap.Int("totalAmount", totalAmountWithTax))

	return transaction, nil
}

func (s *transactionService) TrashedTransaction(ctx context.Context, transaction_id int) (*db.Transaction, error) {
	const method = "TrashedTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transaction_id))

	defer func() {
		end(status)
	}()

	res, err := s.transactionRepository.TrashTransaction(ctx, transaction_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedTrashedTransaction,
			method,
			span,
			zap.Int("transaction_id", transaction_id),
		)
	}

	s.cache.InvalidateTransactionCache(ctx)

	logSuccess("Successfully trashed transaction", zap.Int("transaction_id", transaction_id))

	return res, nil
}

func (s *transactionService) RestoreTransaction(ctx context.Context, transaction_id int) (*db.Transaction, error) {
	const method = "RestoreTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transaction_id))

	defer func() {
		end(status)
	}()

	res, err := s.transactionRepository.RestoreTransaction(ctx, transaction_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedRestoreTransaction,
			method,
			span,
			zap.Int("transaction_id", transaction_id),
		)
	}

	s.cache.InvalidateTransactionCache(ctx)

	logSuccess("Successfully restored transaction", zap.Int("transaction_id", transaction_id))

	return res, nil
}

func (s *transactionService) DeleteTransactionPermanently(ctx context.Context, transactionID int) (bool, error) {
	const method = "DeleteTransactionPermanently"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transactionID", transactionID))

	defer func() {
		end(status)
	}()

	success, err := s.transactionRepository.DeleteTransactionPermanently(ctx, transactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			transaction_errors.ErrFailedDeleteTransactionPermanently,
			method,
			span,
			zap.Int("transaction_id", transactionID),
		)
	}

	s.cache.InvalidateTransactionCache(ctx)

	logSuccess("Successfully permanently deleted transaction", zap.Int("transactionID", transactionID))

	return success, nil
}

func (s *transactionService) RestoreAllTransactions(ctx context.Context) (bool, error) {
	const method = "RestoreAllTransactions"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.transactionRepository.RestoreAllTransactions(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			transaction_errors.ErrFailedRestoreAllTransactions,
			method,
			span,
		)
	}

	s.cache.InvalidateTransactionCache(ctx)

	logSuccess("Successfully restored all trashed transactions")

	return success, nil
}

func (s *transactionService) DeleteAllTransactionPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllTransactionPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.transactionRepository.DeleteAllTransactionPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			transaction_errors.ErrFailedDeleteAllTransactionPermanent,
			method,
			span,
		)
	}

	s.cache.InvalidateTransactionCache(ctx)

	logSuccess("Successfully permanently deleted all transactions")

	return success, nil
}
