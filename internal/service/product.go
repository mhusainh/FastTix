package service

import (
	"context"

	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
)

type ProductService interface {
	GetAllPending(ctx context.Context) ([]entity.Product, error)
	GetByIdPending(ctx context.Context, id int64) (*entity.Product, error)
	GetById(ctx context.Context, id int64) (*entity.Product, error)
	Create(ctx context.Context, req dto.CreateProductRequest) error
	Update(ctx context.Context, req dto.UpdateProductRequest) error
	Delete(ctx context.Context, product *entity.Product) error
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{productRepository}
}

func (s productService) GetAllPending(ctx context.Context) ([]entity.Product, error) {
	return s.productRepository.GetAllPending(ctx)
}

func (s productService) GetByIdPending(ctx context.Context, id int64) (*entity.Product, error) {
	return s.productRepository.GetByIdPending(ctx, id)
}

func (s productService) GetById(ctx context.Context, id int64) (*entity.Product, error) {
	return s.productRepository.GetById(ctx, id)
}

func (s productService) Create(ctx context.Context, req dto.CreateProductRequest) error {
	product := &entity.Product{
		ProductName:        req.ProductName,
		ProductAddress:     req.ProductAddress,
		ProductTime:        req.ProductTime,
		ProductDate:        req.ProductDate,
		ProductPrice:       req.ProductPrice,
		ProductDescription: req.ProductDescription,
		ProductStatus:      "pending",
		UserID:             req.UserID,
	}
	return s.productRepository.Create(ctx, product)
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
	if req.ProductTime != nil {
		product.ProductTime = *req.ProductTime
	}
	if req.ProductDate != nil {
		product.ProductDate = *req.ProductDate
	}
	if req.ProductPrice != 0 {
		product.ProductPrice = req.ProductPrice
	}
	if req.ProductDescription != "" {
		product.ProductDescription = req.ProductDescription
	}
	return s.productRepository.Update(ctx, product)
}

func (s productService) Delete(ctx context.Context, product *entity.Product) error {
	return s.productRepository.Delete(ctx, product)
}
