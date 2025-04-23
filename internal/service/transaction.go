package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"
	"net/http"

	"go.uber.org/zap"
)

type transactionService struct {
	merchantRepository    repository.MerchantRepository
	transactionRepository repository.TransactionRepository
	orderRepository       repository.OrderRepository
	orderItemRepository   repository.OrderItemRepository
	logger                logger.LoggerInterface
	mapping               response_service.TransactionResponseMapper
}

func NewTransactionService(
	merchantRepository repository.MerchantRepository,
	transactionRepository repository.TransactionRepository,
	orderRepository repository.OrderRepository,
	orderItemRepository repository.OrderItemRepository,
	logger logger.LoggerInterface,
	mapping response_service.TransactionResponseMapper,
) *transactionService {
	return &transactionService{
		merchantRepository:    merchantRepository,
		transactionRepository: transactionRepository,
		orderRepository:       orderRepository,
		orderItemRepository:   orderItemRepository,
		mapping:               mapping,
		logger:                logger,
	}
}

func (s *transactionService) FindAllTransactions(req *requests.FindAllTransactions) ([]*response.TransactionResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching transactions",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := s.transactionRepository.FindAllTransactions(req)
	if err != nil {
		s.logger.Error("Failed to retrieve transaction list",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("page_size", pageSize))

		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve transaction list",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched transactions",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToTransactionsResponse(transactions), totalRecords, nil
}

func (s *transactionService) FindByMerchant(req *requests.FindAllTransactionByMerchant) ([]*response.TransactionResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search
	merchant_id := req.MerchantID

	s.logger.Debug("Fetching transactions",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := s.transactionRepository.FindByMerchant(req)

	if err != nil {
		s.logger.Error("Failed to retrieve merchant's transactions from database",
			zap.Int("merchant_id", merchant_id),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Error(err))
		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve merchant's transactions",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched transactions",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToTransactionsResponse(transactions), totalRecords, nil
}

func (s *transactionService) FindByActive(req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching transactions",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := s.transactionRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active transactions from database",
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Error(err))

		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve active transactions",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched transactions",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToTransactionsResponseDeleteAt(transactions), totalRecords, nil
}

func (s *transactionService) FindByTrashed(req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching transactions",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := s.transactionRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed transactions",
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Error(err))

		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve trashed transactions",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched transactions",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToTransactionsResponseDeleteAt(transactions), totalRecords, nil
}

func (s *transactionService) FindMonthlyAmountSuccess(req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	res, err := s.transactionRepository.GetMonthlyAmountSuccess(req)

	if err != nil {
		s.logger.Error("failed to get monthly successful transaction amounts",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve monthly successful transactions data",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToTransactionMonthlyAmountSuccess(res), nil
}

func (s *transactionService) FindYearlyAmountSuccess(year int) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse) {
	res, err := s.transactionRepository.GetYearlyAmountSuccess(year)
	if err != nil {
		s.logger.Error("failed to get yearly successful transaction amounts",
			zap.Int("year", year),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve yearly successful transactions data",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToTransactionYearlyAmountSuccess(res), nil
}

func (s *transactionService) FindMonthlyAmountFailed(req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	res, err := s.transactionRepository.GetMonthlyAmountFailed(req)
	if err != nil {
		s.logger.Error("failed to get monthly failed transaction amounts",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve monthly failed transactions data",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToTransactionMonthlyAmountFailed(res), nil
}

func (s *transactionService) FindYearlyAmountFailed(year int) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse) {
	res, err := s.transactionRepository.GetYearlyAmountFailed(year)

	if err != nil {
		s.logger.Error("failed to get yearly failed transaction amounts",
			zap.Int("year", year),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve yearly failed transactions data",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToTransactionYearlyAmountFailed(res), nil
}

func (s *transactionService) FindMonthlyAmountSuccessByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse) {
	year := req.Year
	month := req.Month
	merchantId := req.MerchantID

	res, err := s.transactionRepository.GetMonthlyAmountSuccessByMerchant(req)

	if err != nil {
		s.logger.Error("failed to get monthly successful transactions by merchant",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Int("merchantID", merchantId),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve merchant's monthly successful transactions",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToTransactionMonthlyAmountSuccess(res), nil
}

func (s *transactionService) FindYearlyAmountSuccessByMerchant(req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse) {
	year := req.Year
	merchantId := req.MerchantID

	res, err := s.transactionRepository.GetYearlyAmountSuccessByMerchant(req)
	if err != nil {
		s.logger.Error("failed to get yearly successful transactions by merchant",
			zap.Int("year", year),
			zap.Int("merchantID", merchantId),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve merchant's yearly successful transactions",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToTransactionYearlyAmountSuccess(res), nil
}

func (s *transactionService) FindMonthlyAmountFailedByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse) {
	year := req.Year
	month := req.Month
	merchantId := req.MerchantID

	res, err := s.transactionRepository.GetMonthlyAmountFailedByMerchant(req)
	if err != nil {
		s.logger.Error("failed to get monthly failed transactions by merchant",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Int("merchantID", merchantId),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve merchant's monthly failed transactions",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToTransactionMonthlyAmountFailed(res), nil
}

func (s *transactionService) FindYearlyAmountFailedByMerchant(req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse) {
	year := req.Year
	merchantId := req.MerchantID

	res, err := s.transactionRepository.GetYearlyAmountFailedByMerchant(req)
	if err != nil {
		s.logger.Error("failed to get yearly failed transactions by merchant",
			zap.Int("year", year),
			zap.Int("merchantID", merchantId),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve merchant's yearly failed transactions",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToTransactionYearlyAmountFailed(res), nil
}

func (s *transactionService) FindMonthlyMethod(year int) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	res, err := s.transactionRepository.GetMonthlyTransactionMethod(year)
	if err != nil {
		s.logger.Error("failed to get monthly transaction methods",
			zap.Int("year", year),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve monthly transaction methods data",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToTransactionMonthlyMethod(res), nil
}

func (s *transactionService) FindYearlyMethod(year int) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	res, err := s.transactionRepository.GetYearlyTransactionMethod(year)

	if err != nil {
		s.logger.Error("failed to get yearly transaction methods",
			zap.Int("year", year),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve yearly transaction methods data",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToTransactionYearlyMethod(res), nil
}

func (s *transactionService) FindMonthlyMethodByMerchant(req *requests.MonthlyYearTransactionMethodMerchant) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	year := req.Year
	merchantId := req.MerchantID

	res, err := s.transactionRepository.GetMonthlyTransactionMethodByMerchant(req)
	if err != nil {
		s.logger.Error("failed to get monthly transaction methods by merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchantId),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve merchant's monthly transaction methods",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToTransactionMonthlyMethod(res), nil
}

func (s *transactionService) FindYearlyMethodByMerchant(req *requests.MonthlyYearTransactionMethodMerchant) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	year := req.Year
	merchantId := req.MerchantID

	res, err := s.transactionRepository.GetYearlyTransactionMethodByMerchant(req)

	if err != nil {
		s.logger.Error("failed to get yearly transaction methods by merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchantId),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve merchant's yearly transaction methods",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToTransactionYearlyMethod(res), nil
}

func (s *transactionService) FindById(transactionID int) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching transaction by ID", zap.Int("transactionID", transactionID))

	transaction, err := s.transactionRepository.FindById(transactionID)

	if err != nil {
		s.logger.Error("Failed to retrieve transaction details",
			zap.Int("transaction_id", transactionID),
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve transaction details",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToTransactionResponse(transaction), nil
}

func (s *transactionService) FindByOrderId(orderID int) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching transaction by Order ID", zap.Int("orderID", orderID))

	transaction, err := s.transactionRepository.FindByOrderId(orderID)

	if err != nil {
		s.logger.Error("Failed to retrieve transaction by order ID",
			zap.Int("order_id", orderID),
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve transaction by order",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToTransactionResponse(transaction), nil
}

func (s *transactionService) CreateTransaction(req *requests.CreateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new transaction", zap.Int("orderID", req.OrderID))

	_, err := s.merchantRepository.FindById(req.MerchantID)

	if err != nil {
		s.logger.Error("Merchant not found", zap.Int("merchantId", req.MerchantID), zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "The requested Merchant was not found",
			Code:    http.StatusNotFound,
		}
	}

	_, err = s.orderRepository.FindById(req.OrderID)

	if err != nil {
		s.logger.Error("Order not found", zap.Int("orderID", req.OrderID), zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "fail",
			Message: "The specified order does not exist in our system",
			Code:    http.StatusNotFound,
		}
	}

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(req.OrderID)

	if err != nil {
		s.logger.Error("Failed to fetch order items", zap.Int("orderID", req.OrderID), zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Unable to retrieve order details at this time",
			Code:    http.StatusInternalServerError,
		}
	}

	if len(orderItems) == 0 {
		return nil, &response.ErrorResponse{
			Status:  "fail",
			Message: "Cannot process transaction for an empty order",
			Code:    http.StatusUnprocessableEntity,
		}
	}

	var totalAmount int
	for _, item := range orderItems {
		if item.Quantity <= 0 {
			return nil, &response.ErrorResponse{
				Status:  "fail",
				Message: "Invalid item quantity in order",
				Code:    http.StatusBadRequest,
			}
		}

		totalAmount += item.Price * item.Quantity
	}

	paymentStatus := "pending"
	if req.Amount > 0 {
		if req.Amount >= totalAmount {
			paymentStatus = "paid"
		} else {
			paymentStatus = "failed"
			return nil, &response.ErrorResponse{
				Status:  "fail",
				Message: "Payment amount is insufficient to complete the transaction",
				Code:    http.StatusPaymentRequired,
			}
		}
	}

	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionRepository.CreateTransaction(req)

	if err != nil {
		s.logger.Error("Failed to create transaction record", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to process transaction due to system error",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToTransactionResponse(transaction), nil
}

func (s *transactionService) UpdateTransaction(req *requests.UpdateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating transaction", zap.Int("transactionID", *req.TransactionID))

	existingTx, err := s.transactionRepository.FindById(*req.TransactionID)

	if err != nil {
		s.logger.Error("Transaction not found", zap.Int("transactionID", *req.TransactionID), zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "fail",
			Message: "The transaction record was not found",
			Code:    http.StatusNotFound,
		}
	}

	if existingTx.PaymentStatus == "paid" || existingTx.PaymentStatus == "refunded" {
		return nil, &response.ErrorResponse{
			Status:  "fail",
			Message: "Completed transactions cannot be modified",
			Code:    http.StatusConflict,
		}
	}

	_, err = s.merchantRepository.FindById(req.MerchantID)

	if err != nil {
		s.logger.Error("Merchant not found", zap.Int("merchantId", req.MerchantID), zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "fail",
			Message: "The merchant account is no longer available",
			Code:    http.StatusNotFound,
		}
	}

	_, err = s.orderRepository.FindById(req.OrderID)

	if err != nil {
		s.logger.Error("Order not found", zap.Int("orderID", req.OrderID), zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "fail",
			Message: "The associated order was not found",
			Code:    http.StatusNotFound,
		}
	}

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(req.OrderID)

	if err != nil {
		s.logger.Error("Failed to retrieve order items", zap.Int("orderID", req.OrderID), zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Unable to retrieve order details",
			Code:    http.StatusInternalServerError,
		}
	}

	var totalAmount int
	for _, item := range orderItems {
		totalAmount += item.Price * item.Quantity
	}

	paymentStatus := "pending"

	if req.Amount > 0 {
		if req.Amount >= totalAmount {
			paymentStatus = "paid"
		} else {
			paymentStatus = "failed"

			return nil, &response.ErrorResponse{
				Status:  "fail",
				Message: "Updated payment amount is insufficient",
				Code:    http.StatusBadRequest,
			}
		}
	}

	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionRepository.UpdateTransaction(req)

	if err != nil {
		s.logger.Error("Failed to update transaction", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update transaction due to system error",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToTransactionResponse(transaction), nil
}

func (s *transactionService) TrashedTransaction(transaction_id int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing transaction", zap.Int("transaction_id", transaction_id))

	res, err := s.transactionRepository.TrashTransaction(transaction_id)

	if err != nil {
		s.logger.Error("Failed to move transaction to trash",
			zap.Int("transaction_id", transaction_id),
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to move transaction to trash",
			Code:    http.StatusInternalServerError,
		}
	}

	so := s.mapping.ToTransactionResponseDeleteAt(res)

	s.logger.Debug("Successfully trashed transaction", zap.Int("transaction_id", transaction_id))

	return so, nil
}

func (s *transactionService) RestoreTransaction(transaction_id int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring transaction", zap.Int("transaction_id", transaction_id))

	res, err := s.transactionRepository.RestoreTransaction(transaction_id)

	if err != nil {
		s.logger.Error("Failed to restore transaction from trash",
			zap.Int("transaction_id", transaction_id),
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore transaction from trash",
			Code:    http.StatusInternalServerError,
		}
	}

	so := s.mapping.ToTransactionResponseDeleteAt(res)

	s.logger.Debug("Successfully restored transaction", zap.Int("transaction_id", transaction_id))

	return so, nil
}

func (s *transactionService) DeleteTransactionPermanently(transactionID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting transaction", zap.Int("transactionID", transactionID))

	success, err := s.transactionRepository.DeleteTransactionPermanently(transactionID)

	if err != nil {
		s.logger.Error("Failed to permanently delete transaction",
			zap.Int("transaction_id", transactionID),
			zap.Error(err))

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete transaction",
			Code:    http.StatusInternalServerError,
		}
	}

	return success, nil
}

func (s *transactionService) RestoreAllTransactions() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed transactions")

	success, err := s.transactionRepository.RestoreAllTransactions()

	if err != nil {
		s.logger.Error("Failed to restore all trashed transactions",
			zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all transactions",
			Code:    http.StatusInternalServerError,
		}
	}

	return success, nil
}

func (s *transactionService) DeleteAllTransactionPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all transactions")

	success, err := s.transactionRepository.DeleteAllTransactionPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed transactions",
			zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all transactions",
			Code:    http.StatusInternalServerError,
		}
	}

	return success, nil
}
