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
	GetByUserId(ctx context.Context, req dto.GetProductByUserIDRequest) ([]entity.Product, error)
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
	query := r.db.WithContext(ctx).Where("product_status = ?", "pending") // Applies to all queries

	if req.Search != "" {
		search := strings.ToLower(req.Search)
		query = query.Where(
			r.db.Where("LOWER(product_name) LIKE ?", "%"+search+"%").
				Or("LOWER(product_category) LIKE ?", "%"+search+"%").
				Or("LOWER(product_address) LIKE ?", "%"+search+"%").
				Or("LOWER(product_price) LIKE ?", "%"+search+"%").
				Or("LOWER(product_date) LIKE ?", "%"+search+"%").
				Or("LOWER(product_time) LIKE ?", "%"+search+"%"),
		)
	}

	// Additional filtering by price range
	if req.MinPrice > 0 && req.MaxPrice > 0 {
		query = query.Where("product_price BETWEEN ? AND ?", req.MinPrice, req.MaxPrice)
	}

	if req.MinPrice > 0 {
		query = query.Where("product_price >= ?", req.MinPrice)
	}

	if req.MaxPrice > 0 {
		query = query.Where("product_price <= ?", req.MinPrice)
	}

	// Additional filtering by date range
	if req.StartDate != "" && req.EndDate != "" {
		query = query.Where("product_date BETWEEN ? AND ?", req.StartDate, req.EndDate)
	}

	// Sorting
	if req.Sort != "" && req.Order != "" {
		query = query.Order(req.Sort + " " + req.Order)
	}

	// Pagination
	if req.Page != 0 && req.Limit != 0 {
		query = query.Offset((req.Page - 1) * req.Limit).Limit(req.Limit)
	}

	// Execute the query
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

func (r *submissionRepository) GetByUserId(ctx context.Context, req dto.GetProductByUserIDRequest) ([]entity.Product, error) {
	submissions := make([]entity.Product, 0)
	query := r.db.WithContext(ctx).Where("user_id = ?", req.UserID)
	if req.Order != "" {
		query = query.Order("created_at " + req.Order)
	}
	if err := query.Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
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