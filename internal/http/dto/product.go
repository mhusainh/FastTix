package dto

import "time"

type GetProductByIDRequest struct {
	ID int64 `param:"id" validate:"required"`
}

type CreateProductRequest struct {
	ProductName        string    `json:"product_name" validate:"required"`
	ProductAddress     string    `json:"product_address" validate:"required"`
	ProductTime        time.Time `json:"product_time" validate:"required"`
	ProductDate        time.Time `json:"product_date" validate:"required"`
	ProductPrice       float64   `json:"product_price" validate:"required"`
	ProductDescription string    `json:"product_description" validate:"required"`
	ProductStatus      string    `json:"product_status" validate:"required"`
	UserID             int64     `json:"user_id" validate:"required"`
}

type UpdateProductRequest struct {
	ID                 int64     `param:"id" validate:"required"`
	ProductName        string    `json:"product_name" validate:"required"`
	ProductAddress     string    `json:"product_address" validate:"required"`
	ProductTime        *time.Time `json:"product_time" validate:"required"`
	ProductDate        *time.Time `json:"product_date" validate:"required"`
	ProductPrice       float64   `json:"product_price" validate:"required"`
	ProductDescription string    `json:"product_description" validate:"required"`
	UserID             int64     `json:"user_id" validate:"required"`
}


type DeleteProductRequest struct{
	ID int64 `param:"id" validate:"required"`
}