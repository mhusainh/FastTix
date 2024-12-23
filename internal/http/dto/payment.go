package dto


type WebhookRequest struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
}

type CheckinWebhook struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	CheckIn           int    `json:"checkin"`
}

type PurchaseTicketResponse struct {
	Message     string `json:"message" form:"message"`
	PaymentLink string `json:"payment_link" form:"payment_link"`
}
