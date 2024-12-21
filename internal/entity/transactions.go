package entity

import "time"

type Transaction struct {
	ID                  int64     `json:"id"`
	TransactionStatus   string    `json:"transaction_status"`
	ProductID           int64     `json:"product_id"`
	UserID              int64     `json:"user_id"`
	TransactionQuantity int       `json:"transaction_quantity"`
	TransactionAmount   float64   `json:"transaction_amount"`
	TransactionType     string    `json:"transaction_type"`
	OrderID             string    `json:"order_id"`
	VerificationToken   string    `json:"verification_token"`
	CheckIn             int       `json:"checkin"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (Transaction) TableName() string {
	return "transactions"
}
