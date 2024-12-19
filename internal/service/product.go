package service

import (
	"context"
	"errors"

	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
)

type ProductService interface {
	GetAll(ctx context.Context) ([]entity.Product, error)
	GetById(ctx context.Context, id int64) (*entity.Product, error)
	Create(ctx context.Context, req dto.CreateProductRequest, t dto.CreateTransactionRequest) error
	Update(ctx context.Context, req dto.UpdateProductRequest) error
	Delete(ctx context.Context, product *entity.Product) error
	Sort(ctx context.Context, sortBy string, order string) ([]entity.Product, error)
	FilterProducts(ctx context.Context, req dto.FilterProductRequest) ([]entity.Product, error)
	SearchProduct(ctx context.Context, keyword string) ([]entity.Product, error)
}

type productService struct {
	productRepository     repository.ProductRepository
	transactionRepository repository.TransactionRepository
}

func NewProductService(
	productRepository repository.ProductRepository,
	transactionRepository repository.TransactionRepository,
) ProductService {
	return &productService{productRepository, transactionRepository}
}

func (s productService) GetAll(ctx context.Context) ([]entity.Product, error) {
	return s.productRepository.GetAll(ctx)
}

func (s productService) GetById(ctx context.Context, id int64) (*entity.Product, error) {
	return s.productRepository.GetById(ctx, id)
}

func (s productService) Create(ctx context.Context, req dto.CreateProductRequest, t dto.CreateTransactionRequest) error {
	userID := req.UserID
	if userID == 0 {
		return errors.New("User ID tidak ditemukan")
	}

	if req.ProductPrice == 0 {
		req.ProductStatus = "pending"
	} else {
		req.ProductStatus = "unpaid"
	}

	product := &entity.Product{
		ProductName:        req.ProductName,
		ProductAddress:     req.ProductAddress,
		ProductTime:        req.ProductTime,
		ProductDate:        req.ProductDate,
		ProductPrice:       req.ProductPrice,
		ProductDescription: req.ProductDescription,
		ProductCategory:    req.ProductCategory,
		ProductQuantity:    req.ProductQuantity,
		ProductType:        "available",
		ProductStatus:      req.ProductStatus,
		UserID:             userID,
	}

	if err := s.productRepository.Create(ctx, product); err != nil {
		return err
	}

	ProductID := product.ID
	TransactionAmount := product.ProductPrice * 0.25

	if TransactionAmount != 0 {
		transaction := &entity.Transaction{
			ProductID:           ProductID,
			UserID:              userID,
			TransactionQuantity: 1,
			TransactionAmount:   TransactionAmount,
			TransactionStatus:   "pending",
		}
		if err := s.transactionRepository.Create(ctx, transaction); err != nil {
			return err
		}
	}

	return nil

}

func (s productService) Update(ctx context.Context, req dto.UpdateProductRequest) error {
	product, err := s.productRepository.GetById(ctx, req.ID)
	if err != nil {
		return err
	}
	if req.ProductName != "" {
		product.ProductName = req.ProductName
	}
	if req.ProductAddress != "" {
		product.ProductAddress = req.ProductAddress
	}
	if req.ProductTime != "" {
		product.ProductTime = req.ProductTime
	}
	if req.ProductDate != "" {
		product.ProductDate = req.ProductDate
	}
	if req.ProductPrice != 0 {
		product.ProductPrice = req.ProductPrice
	}
	if req.ProductDescription != "" {
		product.ProductDescription = req.ProductDescription
	}
	if req.ProductCategory != "" {
		product.ProductCategory = req.ProductCategory
	}
	if req.ProductQuantity != 0 {
		product.ProductQuantity = req.ProductQuantity
	}
	if req.ProductType != "" {
		product.ProductType = req.ProductType
	}
	return s.productRepository.Update(ctx, product)
}

func (s productService) Delete(ctx context.Context, product *entity.Product) error {
	return s.productRepository.Delete(ctx, product)
}

func (s productService) Sort(ctx context.Context, sortBy string, order string) ([]entity.Product, error) {
	return s.productRepository.SortProducts(ctx, sortBy, order)
}

// Filter products
func (s productService) FilterProducts(ctx context.Context, req dto.FilterProductRequest) ([]entity.Product, error) {
	filters := make(map[string]interface{})

	if req.MinPrice != nil {
		filters["min_price"] = *req.MinPrice
	}
	if req.MaxPrice != nil {
		filters["max_price"] = *req.MaxPrice
	}
	if req.Category != nil {
		filters["category"] = *req.Category
	}
	if req.Location != nil {
		filters["location"] = *req.Location
	}
	if req.Price != nil {
		filters["price"] = *req.Price
	}
	if req.Date != nil {
		filters["date"] = *req.Date
	}
	if req.Time != nil {
		filters["time"] = *req.Time
	}

	return s.productRepository.FilterProducts(ctx, filters)
}

// Search products by keyword
func (s productService) SearchProduct(ctx context.Context, keyword string) ([]entity.Product, error) {
	return s.productRepository.SearchProduct(ctx, keyword)
}