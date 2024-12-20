package dto

type GetProductByIDRequest struct {
	ID int64 `param:"id" validate:"required"`
	UserID int64 `json:"user_id" validate:"required"`
}

type GetProductByUserIDRequest struct {
	UserID int64 `json:"user_id" validate:"required"`
	Order string `query:"order" validate:"required"`
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
	UserID             int64   `json:"user_id" validate:"required"`
}

type DeleteProductRequest struct {
	ID int64 `param:"id" validate:"required"`
}

type FilterProductRequest struct {
	MinPrice *float64 `query:"product_price" validate:"omitempty,gte=0"`
	MaxPrice *float64 `query:"product_price" validate:"omitempty,gte=0"`
	Category *string  `query:"product_category" validate:"omitempty"`
	Location *string  `query:"product_address" validate:"omitempty"`
	Price    *float64 `query:"product_price" validate:"omitempty,gte=0"`
	Date     *string  `query:"product_date" validate:"omitempty,datetime=2006-01-02"` // Validasi format tanggal
	Time     *string  `query:"product_time" validate:"omitempty,datetime=15:04:05"`   // Validasi format waktu
}

type SortProductsRequest struct {
	sortBy string `query:"sort_by" validate:"omitempty,oneof=created_at price"`
	srder  string `query:"order" validate:"omitempty,oneof=ASC DESC"`
}

type SearchProduct struct {
	Keyword string `param:"keyword" validate:"required"`
}

type GetAllProductsRequest struct {
	Page   int  `query:"page" validate:"required"`
	Limit  int  `query:"limit" validate:"required"`
	Search string `query:"search" validate:"required"`
	Sort   string `query:"sort" validate:"required"`
	Order  string `query:"order" validate:"required"`
}