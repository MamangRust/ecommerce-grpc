package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	merchant_errors "ecommerce/pkg/errors/merchant"
	"ecommerce/pkg/errors/order_errors"
	orderitem_errors "ecommerce/pkg/errors/order_item_errors"
	shippingaddress_errors "ecommerce/pkg/errors/shipping_address_errors"
	"ecommerce/pkg/errors/transaction_errors"
	"ecommerce/pkg/logger"

	"go.uber.org/zap"
)

type transactionService struct {
	merchantRepository    repository.MerchantRepository
	transactionRepository repository.TransactionRepository
	orderRepository       repository.OrderRepository
	orderItemRepository   repository.OrderItemRepository
	shippingRepository    repository.ShippingAddressRepository
	logger                logger.LoggerInterface
	mapping               response_service.TransactionResponseMapper
}

func NewTransactionService(
	merchantRepository repository.MerchantRepository,
	transactionRepository repository.TransactionRepository,
	orderRepository repository.OrderRepository,
	orderItemRepository repository.OrderItemRepository,
	shippingRepository repository.ShippingAddressRepository,
	logger logger.LoggerInterface,
	mapping response_service.TransactionResponseMapper,
) *transactionService {
	return &transactionService{
		merchantRepository:    merchantRepository,
		transactionRepository: transactionRepository,
		orderRepository:       orderRepository,
		orderItemRepository:   orderItemRepository,
		shippingRepository:    shippingRepository,
		mapping:               mapping,
		logger:                logger,
	}
}

func (s *transactionService) FindAllTransactions(req *requests.FindAllTransaction) ([]*response.TransactionResponse, *int, *response.ErrorResponse) {
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

		return nil, nil, transaction_errors.ErrFailedFindAllTransactions
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
		return nil, nil, transaction_errors.ErrFailedFindTransactionsByMerchant
	}

	s.logger.Debug("Successfully fetched transactions",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToTransactionsResponse(transactions), totalRecords, nil
}

func (s *transactionService) FindByActive(req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
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

		return nil, nil, transaction_errors.ErrFailedFindTransactionsByActive
	}

	s.logger.Debug("Successfully fetched transactions",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToTransactionsResponseDeleteAt(transactions), totalRecords, nil
}

func (s *transactionService) FindByTrashed(req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
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

		return nil, nil, transaction_errors.ErrFailedFindTransactionsByTrashed
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
		return nil, transaction_errors.ErrFailedFindMonthlyAmountSuccess
	}

	return s.mapping.ToTransactionMonthlyAmountSuccess(res), nil
}

func (s *transactionService) FindYearlyAmountSuccess(year int) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse) {
	res, err := s.transactionRepository.GetYearlyAmountSuccess(year)
	if err != nil {
		s.logger.Error("failed to get yearly successful transaction amounts",
			zap.Int("year", year),
			zap.Error(err))
		return nil, transaction_errors.ErrFailedFindYearlyAmountSuccess
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
		return nil, transaction_errors.ErrFailedFindMonthlyAmountFailed
	}

	return s.mapping.ToTransactionMonthlyAmountFailed(res), nil
}

func (s *transactionService) FindYearlyAmountFailed(year int) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse) {
	res, err := s.transactionRepository.GetYearlyAmountFailed(year)

	if err != nil {
		s.logger.Error("failed to get yearly failed transaction amounts",
			zap.Int("year", year),
			zap.Error(err))
		return nil, transaction_errors.ErrFailedFindYearlyAmountFailed
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
		return nil, transaction_errors.ErrFailedFindMonthlyAmountSuccessByMerchant
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
		return nil, transaction_errors.ErrFailedFindYearlyAmountSuccessByMerchant
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
		return nil, transaction_errors.ErrFailedFindMonthlyAmountFailedByMerchant
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
		return nil, transaction_errors.ErrFailedFindYearlyAmountFailedByMerchant
	}

	return s.mapping.ToTransactionYearlyAmountFailed(res), nil
}

func (s *transactionService) FindMonthlyMethodSuccess(req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	res, err := s.transactionRepository.GetMonthlyTransactionMethodSuccess(req)
	if err != nil {
		s.logger.Error("failed to get monthly transaction methods",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.Error(err))
		return nil, transaction_errors.ErrFailedFindMonthlyMethod
	}

	return s.mapping.ToTransactionMonthlyMethod(res), nil
}

func (s *transactionService) FindYearlyMethodSuccess(year int) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	res, err := s.transactionRepository.GetYearlyTransactionMethodSuccess(year)

	if err != nil {
		s.logger.Error("failed to get yearly transaction methods",
			zap.Int("year", year),
			zap.Error(err))
		return nil, transaction_errors.ErrFailedFindYearlyMethod
	}

	return s.mapping.ToTransactionYearlyMethod(res), nil
}

func (s *transactionService) FindMonthlyMethodByMerchantSuccess(req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	year := req.Year
	merchantId := req.MerchantID

	res, err := s.transactionRepository.GetMonthlyTransactionMethodByMerchantSuccess(req)
	if err != nil {
		s.logger.Error("failed to get monthly transaction methods by merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchantId),
			zap.Error(err))
		return nil, transaction_errors.ErrFailedFindMonthlyMethodByMerchant
	}

	return s.mapping.ToTransactionMonthlyMethod(res), nil
}

func (s *transactionService) FindYearlyMethodByMerchantSuccess(req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	year := req.Year
	merchantId := req.MerchantID

	res, err := s.transactionRepository.GetYearlyTransactionMethodByMerchantSuccess(req)

	if err != nil {
		s.logger.Error("failed to get yearly transaction methods by merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchantId),
			zap.Error(err))
		return nil, transaction_errors.ErrFailedFindYearlyMethodByMerchant
	}

	return s.mapping.ToTransactionYearlyMethod(res), nil
}

func (s *transactionService) FindMonthlyMethodFailed(req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	res, err := s.transactionRepository.GetMonthlyTransactionMethodFailed(req)
	if err != nil {
		s.logger.Error("failed to get monthly transaction methods",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.Error(err))
		return nil, transaction_errors.ErrFailedFindMonthlyMethod
	}

	return s.mapping.ToTransactionMonthlyMethod(res), nil
}

func (s *transactionService) FindYearlyMethodFailed(year int) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	res, err := s.transactionRepository.GetYearlyTransactionMethodFailed(year)

	if err != nil {
		s.logger.Error("failed to get yearly transaction methods",
			zap.Int("year", year),
			zap.Error(err))
		return nil, transaction_errors.ErrFailedFindYearlyMethod
	}

	return s.mapping.ToTransactionYearlyMethod(res), nil
}

func (s *transactionService) FindMonthlyMethodByMerchantFailed(req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	year := req.Year
	merchantId := req.MerchantID

	res, err := s.transactionRepository.GetMonthlyTransactionMethodByMerchantFailed(req)
	if err != nil {
		s.logger.Error("failed to get monthly transaction methods by merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchantId),
			zap.Error(err))
		return nil, transaction_errors.ErrFailedFindMonthlyMethodByMerchant
	}

	return s.mapping.ToTransactionMonthlyMethod(res), nil
}

func (s *transactionService) FindYearlyMethodByMerchantFailed(req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	year := req.Year
	merchantId := req.MerchantID

	res, err := s.transactionRepository.GetYearlyTransactionMethodByMerchantFailed(req)

	if err != nil {
		s.logger.Error("failed to get yearly transaction methods by merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchantId),
			zap.Error(err))
		return nil, transaction_errors.ErrFailedFindYearlyMethodByMerchant
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

		return nil, transaction_errors.ErrFailedFindTransactionById
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

		return nil, transaction_errors.ErrFailedFindTransactionByOrderId
	}

	return s.mapping.ToTransactionResponse(transaction), nil
}

func (s *transactionService) CreateTransaction(req *requests.CreateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new transaction", zap.Int("orderID", req.OrderID))

	_, err := s.merchantRepository.FindById(req.MerchantID)

	if err != nil {
		s.logger.Error("Merchant not found", zap.Int("merchantId", req.MerchantID), zap.Error(err))
		return nil, merchant_errors.ErrFailedFindMerchantById
	}

	_, err = s.orderRepository.FindById(req.OrderID)

	if err != nil {
		s.logger.Error("Order not found", zap.Int("orderID", req.OrderID), zap.Error(err))
		return nil, order_errors.ErrFailedFindOrderById
	}

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(req.OrderID)

	if err != nil {
		s.logger.Error("Failed to fetch order items", zap.Int("orderID", req.OrderID), zap.Error(err))

		return nil, orderitem_errors.ErrFailedFindOrderItemByOrder
	}

	if len(orderItems) == 0 {
		return nil, orderitem_errors.ErrFailedFindOrderItemByOrder
	}

	shipping, err := s.shippingRepository.FindByOrder(req.OrderID)

	if err != nil {
		s.logger.Error("Failed to fetch shipping address", zap.Int("orderID", req.OrderID), zap.Error(err))

		return nil, shippingaddress_errors.ErrFailedFindShippingAddressByOrder
	}

	var totalAmount int
	for _, item := range orderItems {
		if item.Quantity <= 0 {
			return nil, orderitem_errors.ErrFailedFindOrderItemByOrder
		}
		totalAmount += item.Price*item.Quantity + shipping.ShippingCost
	}

	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn + shipping.ShippingCost

	var paymentStatus string
	if req.Amount >= totalAmountWithTax {
		paymentStatus = "success"
	} else {
		paymentStatus = "failed"
		return nil, transaction_errors.ErrFailedPaymentInsufficientBalance
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionRepository.CreateTransaction(req)

	if err != nil {
		s.logger.Error("Failed to create transaction record", zap.Error(err))

		return nil, transaction_errors.ErrFailedCreateTransaction
	}

	s.logger.Debug("Successfully created transaction", zap.Any("transaction", transaction))

	return s.mapping.ToTransactionResponse(transaction), nil
}

func (s *transactionService) UpdateTransaction(req *requests.UpdateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating transaction", zap.Int("transactionID", *req.TransactionID))

	existingTx, err := s.transactionRepository.FindById(*req.TransactionID)

	if err != nil {
		s.logger.Error("Transaction not found", zap.Int("transactionID", *req.TransactionID), zap.Error(err))

		return nil, transaction_errors.ErrFailedFindTransactionById
	}

	if existingTx.PaymentStatus == "success" || existingTx.PaymentStatus == "refunded" {
		return nil, transaction_errors.ErrFailedPaymentStatusCannotBeModified
	}

	_, err = s.merchantRepository.FindById(req.MerchantID)

	if err != nil {
		s.logger.Error("Merchant not found", zap.Int("merchantId", req.MerchantID), zap.Error(err))

		return nil, merchant_errors.ErrFailedFindMerchantById
	}

	_, err = s.orderRepository.FindById(req.OrderID)

	if err != nil {
		s.logger.Error("Order not found", zap.Int("orderID", req.OrderID), zap.Error(err))

		return nil, order_errors.ErrFailedFindOrderById
	}

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(req.OrderID)

	if err != nil {
		s.logger.Error("Failed to retrieve order items", zap.Int("orderID", req.OrderID), zap.Error(err))

		return nil, orderitem_errors.ErrFailedFindOrderItemByOrder
	}

	shipping, err := s.shippingRepository.FindByOrder(req.OrderID)

	if err != nil {
		s.logger.Error("Failed to retrieve shipping address", zap.Int("orderID", req.OrderID), zap.Error(err))

		return nil, shippingaddress_errors.ErrFailedFindShippingAddressByOrder
	}

	var totalAmount int
	for _, item := range orderItems {
		if item.Quantity <= 0 {
			return nil, orderitem_errors.ErrFailedFindOrderItemByOrder
		}
		totalAmount += item.Price*item.Quantity + shipping.ShippingCost
	}

	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn + shipping.ShippingCost

	var paymentStatus string
	if req.Amount >= totalAmountWithTax {
		paymentStatus = "success"
	} else {
		paymentStatus = "failed"
		return nil, transaction_errors.ErrFailedPaymentInsufficientBalance
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionRepository.UpdateTransaction(req)

	if err != nil {
		s.logger.Error("Failed to update transaction", zap.Error(err))

		return nil, transaction_errors.ErrFailedUpdateTransaction
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

		return nil, transaction_errors.ErrFailedTrashedTransaction
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

		return nil, transaction_errors.ErrFailedRestoreTransaction
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

		return false, transaction_errors.ErrFailedDeleteTransactionPermanently
	}

	return success, nil
}

func (s *transactionService) RestoreAllTransactions() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed transactions")

	success, err := s.transactionRepository.RestoreAllTransactions()

	if err != nil {
		s.logger.Error("Failed to restore all trashed transactions",
			zap.Error(err))
		return false, transaction_errors.ErrFailedRestoreAllTransactions
	}

	return success, nil
}

func (s *transactionService) DeleteAllTransactionPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all transactions")

	success, err := s.transactionRepository.DeleteAllTransactionPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed transactions",
			zap.Error(err))
		return false, transaction_errors.ErrFailedDeleteAllTransactionPermanent
	}

	return success, nil
}
