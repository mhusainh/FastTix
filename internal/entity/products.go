package entity

import "time"

type Product struct {
	ID                    int64     `json:"id"`
	ProductName           string    `json:"product_name"`
	ProductAddress        string    `json:"product_address"`
	ProductTime           string    `json:"product_time"`
	ProductDate           string    `json:"product_date"`
	ProductPrice          float64   `json:"product_price"`
	ProductDescription    string    `json:"product_description"`
	ProductCategory       string    `json:"product_category"`
	ProductStatus         string    `json:"product_status"`
	VerifySubmissionToken string    `json:"verify_submission_token"`
	IsVerified            int     `json:"is_verified"`
	UserID                int64     `json:"user_id"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}
