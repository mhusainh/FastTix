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
	VerificationToken string  `json:"verification_token" form:"verification_token"`
	OrderID           string  `json:"order_id" form:"order_id"`
	Amount            float64 `json:"amount" form:"amount"`
	Email             string  `json:"email" form:"email"`
	UserID            int64   `json:"user_id"`
	NameProduct       string  `json:"name_product" form:"name_product"`
}

type WebhookRequest struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
}
