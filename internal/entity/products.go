package entity

import "time"

type Product struct {
	ID                    int64     `json:"id"`
	ProductName           string    `json:"product_name"`
	ProductAddress        string    `json:"product_address"`
	ProductImage          string    `json:"product_image"`
	ProductTime           string    `json:"product_time"`
	ProductDate           string    `json:"product_date"`
	ProductPrice          float64   `json:"product_price"`
	ProductSold           int       `json:"product_sold"`
	ProductDescription    string    `json:"product_description"`
	ProductCategory       string    `json:"product_category"`
	ProductQuantity       int       `json:"product_quantity"`
	ProductType           string    `json:"product_type"`
	ProductStatus         string    `json:"product_status"`
	UserID                int64     `json:"user_id"`
	OrderID               string    `json:"order_id"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}
