package repository

import (
	"context"

	"github.com/mhusainh/FastTix/internal/http/dto"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment dto.CreatePaymentRequest) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRequestRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) Create(ctx context.Context, payment dto.CreatePaymentRequest) error {
	return r.db.WithContext(ctx).Create(&payment).Error
}
