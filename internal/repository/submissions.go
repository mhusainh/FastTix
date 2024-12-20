package repository

import (
	"context"
	"strings"

	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"gorm.io/gorm"
)

type SubmissionRepository interface {
	GetAll(ctx context.Context, req dto.GetAllProductsRequest) ([]entity.Product, error)
	GetById(ctx context.Context, id int64) (*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, product *entity.Product) error
}

type submissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db}
}

func (r *submissionRepository) GetAll(ctx context.Context, req dto.GetAllProductsRequest) ([]entity.Product, error) {
	result := make([]entity.Product, 0)
	query := r.db.WithContext(ctx)
	if req.Search != "" {
		search := strings.ToLower(req.Search)
		query = query.Where("product_status = ? AND LOWER(product_name) LIKE ?", "pending", "%"+search+"%").
			Or("product_status = ? AND LOWER(product_category) LIKE ?", "pending", "%"+search+"%").
			Or("product_status = ? AND LOWER(product_address) LIKE ?", "pending", "%"+search+"%").
			Or("product_status = ? AND LOWER(product_price) LIKE ?", "pending", "%"+search+"%").
			Or("product_status = ? AND LOWER(product_date) LIKE ?", "pending", "%"+search+"%").
			Or("product_status = ? AND LOWER(product_time) LIKE ?", "pending", "%"+search+"%")
	}
	if req.Sort != "" && req.Order != "" {
		query = query.Order(req.Sort + " " + req.Order)
	}
	if req.Page != 0 && req.Limit != 0 {
		query = query.Offset((req.Page - 1) * req.Limit).Limit(req.Limit)
	}
	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *submissionRepository) GetById(ctx context.Context, id int64) (*entity.Product, error) {
	result := new(entity.Product)
	if err := r.db.WithContext(ctx).Where("id = ? AND product_status = ?", id, "pending").First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *submissionRepository) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(&product).Error
}

func (r *submissionRepository) Update(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Updates(product).Error
}

func (r *submissionRepository) Delete(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Delete(product).Error
}