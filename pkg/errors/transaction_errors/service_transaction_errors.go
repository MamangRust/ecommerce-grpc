package transaction_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrFailedPaymentStatusCannotBeModified = errors.NewErrorResponse("Cannot modify payment status", http.StatusBadRequest)
	ErrFailedPaymentStatusInvalid          = errors.NewErrorResponse("Invalid payment status", http.StatusBadRequest)
	ErrFailedPaymentInsufficientBalance    = errors.NewErrorResponse("Insufficient balance", http.StatusBadRequest)
	ErrFailedOrderItemEmpty                = errors.NewErrorResponse("Failed to order item empty", http.StatusInternalServerError)

	ErrFailedFindMonthlyAmountSuccess = errors.NewErrorResponse("Failed to find monthly amount success", http.StatusInternalServerError)
	ErrFailedFindYearlyAmountSuccess  = errors.NewErrorResponse("Failed to find yearly amount success", http.StatusInternalServerError)
	ErrFailedFindMonthlyAmountFailed  = errors.NewErrorResponse("Failed to find monthly amount failed", http.StatusInternalServerError)
	ErrFailedFindYearlyAmountFailed   = errors.NewErrorResponse("Failed to find yearly amount failed", http.StatusInternalServerError)

	ErrFailedFindMonthlyAmountSuccessByMerchant = errors.NewErrorResponse("Failed to find monthly amount success by merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyAmountSuccessByMerchant  = errors.NewErrorResponse("Failed to find yearly amount success by merchant", http.StatusInternalServerError)
	ErrFailedFindMonthlyAmountFailedByMerchant  = errors.NewErrorResponse("Failed to find monthly amount failed by merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyAmountFailedByMerchant   = errors.NewErrorResponse("Failed to find yearly amount failed by merchant", http.StatusInternalServerError)

	ErrFailedFindMonthlyMethod           = errors.NewErrorResponse("Failed to find monthly method", http.StatusInternalServerError)
	ErrFailedFindYearlyMethod            = errors.NewErrorResponse("Failed to find yearly method", http.StatusInternalServerError)
	ErrFailedFindMonthlyMethodByMerchant = errors.NewErrorResponse("Failed to find monthly method by merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyMethodByMerchant  = errors.NewErrorResponse("Failed to find yearly method by merchant", http.StatusInternalServerError)

	ErrFailedFindAllTransactions        = errors.NewErrorResponse("Failed to find all transactions", http.StatusInternalServerError)
	ErrFailedFindTransactionsByMerchant = errors.NewErrorResponse("Failed to find transactions by merchant", http.StatusInternalServerError)
	ErrFailedFindTransactionsByActive   = errors.NewErrorResponse("Failed to find active transactions", http.StatusInternalServerError)
	ErrFailedFindTransactionsByTrashed  = errors.NewErrorResponse("Failed to find trashed transactions", http.StatusInternalServerError)
	ErrFailedFindTransactionById        = errors.NewErrorResponse("Failed to find transaction by ID", http.StatusInternalServerError)
	ErrFailedFindTransactionByOrderId   = errors.NewErrorResponse("Failed to find transaction by order ID", http.StatusInternalServerError)

	ErrFailedCreateTransaction             = errors.NewErrorResponse("Failed to create transaction", http.StatusInternalServerError)
	ErrFailedUpdateTransaction             = errors.NewErrorResponse("Failed to update transaction", http.StatusInternalServerError)
	ErrFailedTrashedTransaction            = errors.NewErrorResponse("Failed to trash transaction", http.StatusInternalServerError)
	ErrFailedRestoreTransaction            = errors.NewErrorResponse("Failed to restore transaction", http.StatusInternalServerError)
	ErrFailedDeleteTransactionPermanently  = errors.NewErrorResponse("Failed to permanently delete transaction", http.StatusInternalServerError)
	ErrFailedRestoreAllTransactions        = errors.NewErrorResponse("Failed to restore all transactions", http.StatusInternalServerError)
	ErrFailedDeleteAllTransactionPermanent = errors.NewErrorResponse("Failed to permanently delete all transactions", http.StatusInternalServerError)
)
