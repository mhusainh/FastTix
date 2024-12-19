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
	UserID             int64   `json:"user_id" validate:"required"`
}

type DeleteProductRequest struct {
	ID int64 `param:"id" validate:"required"`
}

type SearchProductRequest struct {
	Search string `query:"search" validate:"required"`
}

type FilterProductRequest struct {
	Address  string `query:"address" validate:"required"`   // address
	Category string `query:"category" validate:"required"`  // category
	MinPrice string `query:"min_price" validate:"required"` // min_price
	MaxPrice string `query:"max_price" validate:"required"` // max_price
	Status   string `query:"status" validate:"required"`    // status
	Date     string `query:"date" validate:"required"`      // date
	Time     string `query:"time" validate:"required"`      // time
}

type SortProductRequest struct {
	Sort string `query:"sort" validate:"required"`
}
