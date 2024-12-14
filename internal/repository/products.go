package repository

import (
	"context"

	"github.com/mhusainh/FastTix/internal/entity"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAllPending(ctx context.Context) ([]entity.Product, error)
	GetByIdPending(ctx context.Context, id int64) (*entity.Product, error)
	GetById(ctx context.Context, id int64) (*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, product *entity.Product) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository{
	return &productRepository{db}
}

func (r *productRepository) GetAllPending(ctx context.Context) ([]entity.Product, error){
	result := make([]entity.Product, 0)
	if err := r.db.WithContext(ctx).Where("product_status = ?", "pending").Find(&result).Error; err != nil{
		return nil, err
	}
	return result, nil
}

func (r *productRepository) GetByIdPending(ctx context.Context, id int64) (*entity.Product, error){
	result := new(entity.Product)
	if err := r.db.WithContext(ctx).Where("id = ? AND product_status = ?", id, "pending").First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *productRepository) GetById(ctx context.Context, id int64) (*entity.Product, error){
	result := new(entity.Product)
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(&product).Error
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Updates(&product).Error
}

func (r *productRepository) Delete(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Delete(&product).Error
}
