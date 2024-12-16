package dto

type GetProductByIDRequest struct {
	ID int64 `param:"id" validate:"required"`
}

type CreateProductRequest struct {
	ProductName        string  `json:"product_name" validate:"required"`
	ProductAddress     string  `json:"product_address" validate:"required"`
	ProductTime        string  `json:"product_time" validate:"required"`
	ProductDate        string  `json:"product_date" validate:"required"`
	ProductPrice       float64 `json:"product_price" validate:"required"`
	ProductDescription string  `json:"product_description" validate:"required"`
	ProductCategory    string  `json:"product_category" validate:"required"`
	ProductQuantity    int     `json:"product_quantity" validate:"required"`
	ProductType        string  `json:"product_type" validate:"required"`
	ProductStatus      string  `json:"product_status" validate:"required"`
	UserID             int64   `json:"user_id" validate:"required"`
}

type UpdateProductRequest struct {
	ID                 int64   `param:"id" validate:"required"`
	ProductName        string  `json:"product_name" validate:"required"`
	ProductAddress     string  `json:"product_address" validate:"required"`
	ProductTime        string  `json:"product_time" validate:"required"`
	ProductDate        string  `json:"product_date" validate:"required"`
	ProductPrice       float64 `json:"product_price" validate:"required"`
	ProductDescription string  `json:"product_description" validate:"required"`
	ProductCategory    string  `json:"product_category" validate:"required"`
	ProductQuantity    int     `json:"product_quantity" validate:"required"`
	ProductType        string  `json:"product_type" validate:"required"`
}

type DeleteProductRequest struct {
	ID int64 `param:"id" validate:"required"`
}