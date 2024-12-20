package dto

type PaymentRequest struct {
	Amount int `json:"amount" form:"amount"`
}

type PaymentResponse struct {
	Amount int `json:"amount" form:"amount"`
}

type PaymentStatusResponse struct {
	Status string `json:"status" form:"status"`
}

type VerifyPaymentRequest struct {
	OrderID string `json:"order_id" form:"order_id"`
}

type CreatePaymentRequest struct {
	OrderID string `json:"order_id" form:"order_id"`
	Amount  int    `json:"amount" form:"amount"`
	Email   string `json:"email" form:"email"`
	UserID  int64  `json:"user_id"`
}

type WebhookRequest struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
}
