package repository

import (
	"context"
	"strings"

	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAll(ctx context.Context, req dto.GetAllProductsRequest) ([]entity.Product, error)
	GetById(ctx context.Context, id int64) (*entity.Product, error)
	GetByName(ctx context.Context, name string) (*entity.Product, error)
	Delete(ctx context.Context, product *entity.Product) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) GetAll(ctx context.Context, req dto.GetAllProductsRequest) ([]entity.Product, error) {
	products := make([]entity.Product, 0)
	query := r.db.WithContext(ctx)
	if req.Search != "" {
		search := strings.ToLower(req.Search)
		query = query.Where("LOWER(product_name) LIKE ?", "%"+search+"%").
			Or("LOWER(product_category) LIKE ?", "%"+search+"%").
			Or("LOWER(product_address) LIKE ?", "%"+search+"%").
			Or("LOWER(product_price) LIKE ?", "%"+search+"%").
			Or("LOWER(product_date) LIKE ?", "%"+search+"%").
			Or("LOWER(product_time) LIKE ?", "%"+search+"%")
	}
	if req.Sort != "" && req.Order != "" {
		query = query.Order(req.Sort + " " + req.Order)
	}
	if req.Page != 0 && req.Limit != 0 {
		query = query.Offset((req.Page - 1) * req.Limit).Limit(req.Limit)
	}
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) GetById(ctx context.Context, id int64) (*entity.Product, error) {
	result := new(entity.Product)
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *productRepository) GetByName(ctx context.Context, name string) (*entity.Product, error) {
	result := new(entity.Product)
	if err := r.db.WithContext(ctx).Where("product_name = ?", name).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(&product).Error
}

func (r *productRepository) Delete(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Delete(&product).Error
}