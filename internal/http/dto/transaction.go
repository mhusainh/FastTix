package dto

type CreateTransactionRequest struct {
	UserID    int64 `json:"user_id" validate:"required"`
	ProductID int64 `json:"product_id" validate:"required"`
	TransactionQuantity  int   `json:"transaction_quantity" validate:"required"`
	TransactionAmount    float64 `json:"transaction_amount" validate:"required"`
	TransactionStatus    string `json:"transaction_status" validate:"required"`
}