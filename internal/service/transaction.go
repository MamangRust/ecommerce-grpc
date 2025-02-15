package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"

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

func (s *transactionService) FindAllTransactions(search string, page, pageSize int) ([]*response.TransactionResponse, int, *response.ErrorResponse) {
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

	transactions, totalRecords, err := s.transactionRepository.FindAllTransactions(search, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch transactions",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch transactions"}
	}

	s.logger.Debug("Successfully fetched transactions",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToTransactionsResponse(transactions), totalRecords, nil
}

func (s *transactionService) FindByMerchant(merchant_id int, search string, page, pageSize int) ([]*response.TransactionResponse, int, *response.ErrorResponse) {
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

	transactions, totalRecords, err := s.transactionRepository.FindByMerchant(merchant_id, search, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch transactions",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch transactions"}
	}

	s.logger.Debug("Successfully fetched transactions",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToTransactionsResponse(transactions), totalRecords, nil
}

func (s *transactionService) FindByActive(search string, page, pageSize int) ([]*response.TransactionResponseDeleteAt, int, *response.ErrorResponse) {
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

	transactions, totalRecords, err := s.transactionRepository.FindByActive(search, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch transactions",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch active transactions"}
	}

	s.logger.Debug("Successfully fetched transactions",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToTransactionsResponseDeleteAt(transactions), totalRecords, nil
}

func (s *transactionService) FindByTrashed(search string, page, pageSize int) ([]*response.TransactionResponseDeleteAt, int, *response.ErrorResponse) {
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

	transactions, totalRecords, err := s.transactionRepository.FindByTrashed(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch transactions",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch trashed transactions"}
	}

	s.logger.Debug("Successfully fetched transactions",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToTransactionsResponseDeleteAt(transactions), totalRecords, nil
}

func (s *transactionService) FindById(transactionID int) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching transaction by ID", zap.Int("transactionID", transactionID))

	transaction, err := s.transactionRepository.FindById(transactionID)
	if err != nil {
		s.logger.Error("Failed to fetch transaction", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Transaction not found"}
	}

	return s.mapping.ToTransactionResponse(transaction), nil
}

func (s *transactionService) FindByOrderId(orderID int) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching transaction by Order ID", zap.Int("orderID", orderID))

	transaction, err := s.transactionRepository.FindByOrderId(orderID)
	if err != nil {
		s.logger.Error("Failed to fetch transaction by order ID", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Transaction not found"}
	}

	return s.mapping.ToTransactionResponse(transaction), nil
}

func (s *transactionService) CreateTransaction(req *requests.CreateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new transaction", zap.Int("orderID", req.OrderID))

	_, err := s.merchantRepository.FindById(req.MerchantID)

	if err != nil {
		s.logger.Error("Merchant not found", zap.Int("merchantId", req.MerchantID), zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Merchant not found"}
	}

	_, err = s.orderRepository.FindById(req.OrderID)
	if err != nil {
		s.logger.Error("Order not found", zap.Int("orderID", req.OrderID), zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Order not found"}
	}

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(req.OrderID)
	if err != nil {
		s.logger.Error("Failed to fetch order items", zap.Int("orderID", req.OrderID), zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to fetch order items"}
	}

	var totalAmount int
	for _, item := range orderItems {
		totalAmount += item.Price * item.Quantity
	}

	if req.PaymentStatus != "pending" && req.PaymentStatus != "paid" && req.PaymentStatus != "failed" {
		s.logger.Error("Invalid payment status", zap.String("paymentStatus", req.PaymentStatus))
		return nil, &response.ErrorResponse{Status: "error", Message: "Invalid payment status"}
	}

	if req.PaymentStatus == "paid" {
		if req.Amount < totalAmount {
			s.logger.Error("Insufficient payment amount", zap.Int("amount", req.Amount), zap.Int("totalAmount", totalAmount))
			return nil, &response.ErrorResponse{Status: "error", Message: "Insufficient payment amount"}
		}
		req.ChangeAmount = req.Amount - totalAmount
	} else if req.PaymentStatus == "failed" {
		if req.Amount >= totalAmount {
			s.logger.Error("Invalid amount for failed payment", zap.Int("amount", req.Amount), zap.Int("totalAmount", totalAmount))
			return nil, &response.ErrorResponse{Status: "error", Message: "Invalid amount for failed payment"}
		}
	}

	transaction, err := s.transactionRepository.CreateTransaction(req)
	if err != nil {
		s.logger.Error("Failed to create transaction", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to create transaction"}
	}

	return s.mapping.ToTransactionResponse(transaction), nil
}

func (s *transactionService) UpdateTransaction(req *requests.UpdateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating transaction", zap.Int("transactionID", req.TransactionID))

	_, err := s.transactionRepository.FindById(req.TransactionID)
	if err != nil {
		s.logger.Error("Transaction not found", zap.Int("transactionID", req.TransactionID), zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Transaction not found"}
	}

	_, err = s.merchantRepository.FindById(req.MerchantID)

	if err != nil {
		s.logger.Error("Merchant not found", zap.Int("merchantId", req.MerchantID), zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Merchant not found"}
	}

	_, err = s.orderRepository.FindById(req.OrderID)
	if err != nil {
		s.logger.Error("Order not found", zap.Int("orderID", req.OrderID), zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Order not found"}
	}

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(req.OrderID)
	if err != nil {
		s.logger.Error("Failed to fetch order items", zap.Int("orderID", req.OrderID), zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to fetch order items"}
	}

	var totalAmount int
	for _, item := range orderItems {
		totalAmount += item.Price * item.Quantity
	}

	if req.PaymentStatus != "pending" && req.PaymentStatus != "paid" && req.PaymentStatus != "failed" {
		s.logger.Error("Invalid payment status", zap.String("paymentStatus", req.PaymentStatus))
		return nil, &response.ErrorResponse{Status: "error", Message: "Invalid payment status"}
	}

	if req.PaymentStatus == "paid" {
		if req.Amount < totalAmount {
			s.logger.Error("Insufficient payment amount", zap.Int("amount", req.Amount), zap.Int("totalAmount", totalAmount))
			return nil, &response.ErrorResponse{Status: "error", Message: "Insufficient payment amount"}
		}
		req.ChangeAmount = req.Amount - totalAmount
	} else if req.PaymentStatus == "failed" {
		if req.Amount >= totalAmount {
			s.logger.Error("Invalid amount for failed payment", zap.Int("amount", req.Amount), zap.Int("totalAmount", totalAmount))
			return nil, &response.ErrorResponse{Status: "error", Message: "Invalid amount for failed payment"}
		}
	}

	transaction, err := s.transactionRepository.UpdateTransaction(req)
	if err != nil {
		s.logger.Error("Failed to update transaction", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to update transaction"}
	}

	return s.mapping.ToTransactionResponse(transaction), nil
}

func (s *transactionService) TrashedTransaction(transaction_id int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing transaction", zap.Int("transaction_id", transaction_id))

	res, err := s.transactionRepository.TrashTransaction(transaction_id)

	if err != nil {
		s.logger.Error("Failed to trash transaction", zap.Error(err), zap.Int("transaction_id", transaction_id))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to trash transaction",
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
		s.logger.Error("Failed to restore transaction", zap.Error(err), zap.Int("transaction_id", transaction_id))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore transaction",
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
		s.logger.Error("Failed to permanently delete transaction", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete transaction"}
	}

	return success, nil
}

func (s *transactionService) RestoreAllTransactions() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed transactions")

	success, err := s.transactionRepository.RestoreAllTransactions()
	if err != nil {
		s.logger.Error("Failed to restore all transactions", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to restore all transactions"}
	}

	return success, nil
}

func (s *transactionService) DeleteAllTransactionPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all transactions")

	success, err := s.transactionRepository.DeleteAllTransactionPermanent()
	if err != nil {
		s.logger.Error("Failed to permanently delete all transactions", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete all transactions"}
	}

	return success, nil
}
