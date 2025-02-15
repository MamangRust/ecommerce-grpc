package response

type TransactionResponse struct {
	ID            int    `json:"id"`
	OrderID       int    `json:"order_id"`
	MerchantID    int    `json:"merchant_id"`
	PaymentMethod string `json:"payment_method"`
	Amount        int    `json:"amount"`
	ChangeAmount  int    `json:"change_amount"`
	PaymentStatus string `json:"payment_status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type TransactionResponseDeleteAt struct {
	ID            int    `json:"id"`
	OrderID       int    `json:"order_id"`
	MerchantID    int    `json:"merchant_id"`
	PaymentMethod string `json:"payment_method"`
	Amount        int    `json:"amount"`
	ChangeAmount  int    `json:"change_amount"`
	PaymentStatus string `json:"payment_status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	DeletedAt     string `json:"deleted_at"`
}

type ApiResponseTransaction struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Data    *TransactionResponse `json:"data"`
}

type ApiResponseTransactionDeleteAt struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    *TransactionResponseDeleteAt `json:"data"`
}

type ApiResponsesTransaction struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    []*TransactionResponse `json:"data"`
}

type ApiResponseTransactionDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseTransactionAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationTransactionDeleteAt struct {
	Status     string                         `json:"status"`
	Message    string                         `json:"message"`
	Data       []*TransactionResponseDeleteAt `json:"data"`
	Pagination PaginationMeta                 `json:"pagination"`
}

type ApiResponsePaginationTransaction struct {
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	Data       []*TransactionResponse `json:"data"`
	Pagination PaginationMeta         `json:"pagination"`
}
