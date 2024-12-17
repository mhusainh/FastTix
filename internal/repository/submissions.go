package repository

import (
	"context"

	"github.com/mhusainh/FastTix/internal/entity"
	"gorm.io/gorm"
)

type SubmissionRepository interface {
	GetAll(ctx context.Context) ([]entity.Product, error)
	GetById(ctx context.Context, id int64) (*entity.Product, error)
}

type submissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db}
}

func (r *submissionRepository) GetAll(ctx context.Context) ([]entity.Product, error) {
	result := make([]entity.Product, 0)
	if err := r.db.WithContext(ctx).Where("product_status = ?", "pending").Find(&result).Error; err != nil {
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