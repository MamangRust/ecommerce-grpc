package record

type TransactionRecord struct {
	ID            int     `json:"id"`
	OrderID       int     `json:"order_id"`
	MerchantID    int     `json:"merchant_id"`
	PaymentMethod string  `json:"payment_method"`
	Amount        int     `json:"amount"`
	ChangeAmount  int     `json:"change_amount"`
	PaymentStatus string  `json:"payment_status"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     *string `json:"deleted_at"`
}
