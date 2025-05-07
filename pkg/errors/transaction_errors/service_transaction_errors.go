package transaction_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrFailedPaymentStatusCannotBeModified = response.NewErrorResponse("Cannot modify payment status", http.StatusBadRequest)
	ErrFailedPaymentStatusInvalid          = response.NewErrorResponse("Invalid payment status", http.StatusBadRequest)
	ErrFailedPaymentInsufficientBalance    = response.NewErrorResponse("Insufficient balance", http.StatusBadRequest)
	ErrFailedOrderItemEmpty                = response.NewErrorResponse("Failed to order item empty", http.StatusInternalServerError)

	ErrFailedFindMonthlyAmountSuccess = response.NewErrorResponse("Failed to find monthly amount success", http.StatusInternalServerError)
	ErrFailedFindYearlyAmountSuccess  = response.NewErrorResponse("Failed to find yearly amount success", http.StatusInternalServerError)
	ErrFailedFindMonthlyAmountFailed  = response.NewErrorResponse("Failed to find monthly amount failed", http.StatusInternalServerError)
	ErrFailedFindYearlyAmountFailed   = response.NewErrorResponse("Failed to find yearly amount failed", http.StatusInternalServerError)

	ErrFailedFindMonthlyAmountSuccessByMerchant = response.NewErrorResponse("Failed to find monthly amount success by merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyAmountSuccessByMerchant  = response.NewErrorResponse("Failed to find yearly amount success by merchant", http.StatusInternalServerError)
	ErrFailedFindMonthlyAmountFailedByMerchant  = response.NewErrorResponse("Failed to find monthly amount failed by merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyAmountFailedByMerchant   = response.NewErrorResponse("Failed to find yearly amount failed by merchant", http.StatusInternalServerError)

	ErrFailedFindMonthlyMethod           = response.NewErrorResponse("Failed to find monthly method", http.StatusInternalServerError)
	ErrFailedFindYearlyMethod            = response.NewErrorResponse("Failed to find yearly method", http.StatusInternalServerError)
	ErrFailedFindMonthlyMethodByMerchant = response.NewErrorResponse("Failed to find monthly method by merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyMethodByMerchant  = response.NewErrorResponse("Failed to find yearly method by merchant", http.StatusInternalServerError)

	ErrFailedFindAllTransactions        = response.NewErrorResponse("Failed to find all transactions", http.StatusInternalServerError)
	ErrFailedFindTransactionsByMerchant = response.NewErrorResponse("Failed to find transactions by merchant", http.StatusInternalServerError)
	ErrFailedFindTransactionsByActive   = response.NewErrorResponse("Failed to find active transactions", http.StatusInternalServerError)
	ErrFailedFindTransactionsByTrashed  = response.NewErrorResponse("Failed to find trashed transactions", http.StatusInternalServerError)
	ErrFailedFindTransactionById        = response.NewErrorResponse("Failed to find transaction by ID", http.StatusInternalServerError)
	ErrFailedFindTransactionByOrderId   = response.NewErrorResponse("Failed to find transaction by order ID", http.StatusInternalServerError)

	ErrFailedCreateTransaction             = response.NewErrorResponse("Failed to create transaction", http.StatusInternalServerError)
	ErrFailedUpdateTransaction             = response.NewErrorResponse("Failed to update transaction", http.StatusInternalServerError)
	ErrFailedTrashedTransaction            = response.NewErrorResponse("Failed to trash transaction", http.StatusInternalServerError)
	ErrFailedRestoreTransaction            = response.NewErrorResponse("Failed to restore transaction", http.StatusInternalServerError)
	ErrFailedDeleteTransactionPermanently  = response.NewErrorResponse("Failed to permanently delete transaction", http.StatusInternalServerError)
	ErrFailedRestoreAllTransactions        = response.NewErrorResponse("Failed to restore all transactions", http.StatusInternalServerError)
	ErrFailedDeleteAllTransactionPermanent = response.NewErrorResponse("Failed to permanently delete all transactions", http.StatusInternalServerError)
)
