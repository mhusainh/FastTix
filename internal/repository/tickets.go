package repository

import (
	"context"
	"strings"

	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"gorm.io/gorm"
)

type TicketRepository interface {
	GetAll(ctx context.Context, req dto.GetAllProductsRequest) ([]entity.Product, error)
	GetById(ctx context.Context, id int) (entity.Product, error)
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db}
}

func (r *ticketRepository) GetAll(ctx context.Context, req dto.GetAllProductsRequest) ([]entity.Product, error) {
	result := make([]entity.Product, 0)
	query := r.db.WithContext(ctx)
	if req.Search != "" {
		search := strings.ToLower(req.Search)
		query = query.Where("product_status = ? AND LOWER(product_name) LIKE ?", "accepted", "%"+search+"%").
			Or("product_status = ? AND LOWER(product_category) LIKE ?", "accepted", "%"+search+"%").
			Or("product_status = ? AND LOWER(product_address) LIKE ?", "accepted", "%"+search+"%").
			Or("product_status = ? AND LOWER(product_price) LIKE ?", "accepted", "%"+search+"%").
			Or("product_status = ? AND LOWER(product_date) LIKE ?", "accepted", "%"+search+"%").
			Or("product_status = ? AND LOWER(product_time) LIKE ?", "accepted", "%"+search+"%")
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

func (r *ticketRepository) GetById(ctx context.Context, id int) (entity.Product, error) {
	result := new(entity.Product)
	if err := r.db.WithContext(ctx).Where("id = ? AND product_status = ?", id, "accepted").First(&result).Error; err != nil {
		return entity.Product{}, err
	}
	return *result, nil
}