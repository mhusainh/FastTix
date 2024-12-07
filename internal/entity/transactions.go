package entity

import "time"

type Transaction struct {
	ID        int64     `json:"id"`
	TransactionStatus string `json:"transaction_status"`
	ProductID int64     `json:"product_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Transaction) TableName() string {
	return "transactions"
}