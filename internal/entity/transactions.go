package entity

import "time"

type Transaction struct {
	ID                  int64     `json:"id"`
	TransactionStatus   string    `json:"transaction_status"`
	ProductID           int64     `json:"product_id"`
	UserID              int64     `json:"user_id"`
	TransactionQuantity int       `json:"transaction_quantity"`
	TransactionAmount   float64   `json:"transaction_amount"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	OrderID             string    `json:"order_id"`
	IsVerified          int       `json:"is_verified"`
	VerificationToken   string    `json:"verification_token"`
	CheckIn             int       `json:"checkin"`
}

func (Transaction) TableName() string {
	return "transactions"
}
