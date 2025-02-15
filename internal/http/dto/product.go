package dto

type GetProductByIDRequest struct {
	ID     int64 `param:"id" validate:"required"`
	UserID int64 `json:"user_id" validate:"required"`
}

type GetProductByUserIDRequest struct {
	UserID int64  `json:"user_id" validate:"required"`
	Order  string `query:"order" validate:"required"`
}

type CreateProductRequest struct {
	ProductName        string  `json:"product_name" validate:"required"`
	ProductAddress     string  `json:"product_address" validate:"required"`
	ProductImage       string  `json:"product_image" validate:"required"`
	ProductTime        string  `json:"product_time" validate:"required"`
	ProductDate        string  `json:"product_date" validate:"required"`
	ProductPrice       float64 `json:"product_price" validate:"required"`
	ProductDescription string  `json:"product_description" validate:"required"`
	ProductCategory    string  `json:"product_category" validate:"required"`
	ProductQuantity    int     `json:"product_quantity" validate:"required"`
	ProductType        string  `json:"product_type" validate:"required"`
	ProductStatus      string  `json:"product_status" validate:"required"`
	UserID             int64   `json:"user_id" validate:"required"`
	OrderID            string  `json:"order_id" validate:"required"`
}

type UpdateProductRequest struct {
	ID                 int64   `param:"id" validate:"required"`
	ProductName        string  `json:"product_name" validate:"required"`
	ProductAddress     string  `json:"product_address" validate:"required"`
	ProductImage       string  `json:"product_image" validate:"required"`
	ProductTime        string  `json:"product_time" validate:"required"`
	ProductDate        string  `json:"product_date" validate:"required"`
	ProductPrice       float64 `json:"product_price" validate:"required"`
	ProductSold        int     `json:"product_sold" validate:"required"`
	ProductDescription string  `json:"product_description" validate:"required"`
	ProductCategory    string  `json:"product_category" validate:"required"`
	ProductQuantity    int     `json:"product_quantity" validate:"required"`
	ProductType        string  `json:"product_type" validate:"required"`
	UserID             int64   `json:"user_id" validate:"required"`
}

type UpdateProductStatusRequest struct {
	ID     int64  `param:"id" validate:"required"`
	Status string `param:"status" validate:"required"`
}

type DeleteProductRequest struct {
	ID int64 `param:"id" validate:"required"`
}

type GetAllProductsRequest struct {
	Page      int     `query:"page" validate:"required"`
	Limit     int     `query:"limit" validate:"required"`
	Search    string  `query:"search" validate:"required"`
	Sort      string  `query:"sort" validate:"required"`
	Order     string  `query:"order" validate:"required"`
	MinPrice  float64 `query:"min_price" validate:"required"`
	MaxPrice  float64 `query:"max_price" validate:"required"`
	StartDate string  `query:"start_date" validate:"required"`
	EndDate   string  `query:"end_date" validate:"required"`
}
