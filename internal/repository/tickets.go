package repository

import (
	"context"

	"github.com/mhusainh/FastTix/internal/entity"
	"gorm.io/gorm"
)

type TicketRepository interface {
	GetAll(ctx context.Context) ([]entity.Product, error)
	GetById(ctx context.Context, id int) (entity.Product, error)
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db}
}

func (r *ticketRepository) GetAll(ctx context.Context) ([]entity.Product, error) {
	result := make([]entity.Product, 0)
	if err := r.db.WithContext(ctx).Where("product_status = ?", "accepted").Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ticketRepository) GetById(ctx context.Context, id int) (entity.Product, error) {
	result := new(entity.Product)
	if err := r.db.WithContext(ctx).Where("id = ? AND product_status = ?", id, "accepted").First(&result).Error; err != nil {
		return entity.Product{}, err
	}
	return *result, nil
}
