package repository

import (
	"context"
	"fmt"

	"github.com/mhusainh/FastTix/internal/entity"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]entity.Product, error)
	GetById(ctx context.Context, id int64) (*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, product *entity.Product) error
	FilterProducts(ctx context.Context, filters map[string]interface{}) ([]entity.Product, error)
	SortProducts(ctx context.Context, sortBy string, order string) ([]entity.Product, error)
	SearchProduct(ctx context.Context, keyword string) ([]entity.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) GetAll(ctx context.Context) ([]entity.Product, error) {
	var products []entity.Product
	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
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

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(&product).Error
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Updates(&product).Error
}

func (r *productRepository) Delete(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Delete(&product).Error
}

func (r *productRepository) FilterProducts(ctx context.Context, filters map[string]interface{}) ([]entity.Product, error) {
	var products []entity.Product

	// Build the query dynamically
	query := r.db.WithContext(ctx)
	for key, value := range filters {
		switch key {
		case "min_price":
			query = query.Where("product_price >= ?", value)
		case "max_price":
			query = query.Where("product_price <= ?", value)
		case "category":
			query = query.Where("product_category = ?", value)
		case "location":
			query = query.Where("product_address LIKE ?", fmt.Sprintf("%%%s%%", value))
		case "time":
			query = query.Where("product_time = ?", value)
		case "date":
			query = query.Where("product_date = ?", value)
		case "price":
			query = query.Where("product_price = ?", value)
		}
	}

	// Execute query
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// Generic Sorting
func (r *productRepository) SortProducts(ctx context.Context, sortBy string, order string) ([]entity.Product, error) {
	var products []entity.Product

	// Default sorting if invalid parameters are passed
	if sortBy == "" {
		sortBy = "created_at"
	}
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	// Execute query with sorting
	if err := r.db.WithContext(ctx).Order(fmt.Sprintf("%s %s", sortBy, order)).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// Search Product
func (r *productRepository) SearchProduct(ctx context.Context, keyword string) ([]entity.Product, error) {
	var products []entity.Product

	// Search by name with partial match
	if err := r.db.WithContext(ctx).
		Where("name LIKE ?", fmt.Sprintf("%%%s%%", keyword)).
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
