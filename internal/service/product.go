package service

import (
	"context"

	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
)

type ProductService interface {
	GetAll(ctx context.Context, req dto.GetAllProductsRequest) ([]entity.Product, error)
	GetById(ctx context.Context, id int64) (*entity.Product, error)
	GetByUserId(ctx context.Context, req dto.GetProductByUserIDRequest) ([]entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, product *entity.Product) error
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{productRepository}
}

func (s *productService) GetAll(ctx context.Context, req dto.GetAllProductsRequest) ([]entity.Product, error) {
	return s.productRepository.GetAll(ctx, req)
}

func (s *productService) GetById(ctx context.Context, id int64) (*entity.Product, error) {
	return s.productRepository.GetById(ctx, id)
}

func (s *productService) Delete(ctx context.Context, product *entity.Product) error {
	return s.productRepository.Delete(ctx, product)
}

func (s *productService) GetByUserId(ctx context.Context, req dto.GetProductByUserIDRequest) ([]entity.Product, error) {
	return s.productRepository.GetByUserId(ctx, req)
}

func (s *productService) Update(ctx context.Context, product *entity.Product) error {
	return s.productRepository.Update(ctx, product)
}
