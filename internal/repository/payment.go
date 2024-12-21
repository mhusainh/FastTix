package repository

import (
	"context"

	"github.com/mhusainh/FastTix/internal/http/dto"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment dto.CreatePaymentRequest) error
	GetByTokenTransaction(ctx context.Context, tokenTransaction string) (dto.CreatePaymentRequest, error)
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

func (r *paymentRepository) GetByTokenTransaction(ctx context.Context, tokenTransaction string) (dto.CreatePaymentRequest, error) {
	var payment dto.CreatePaymentRequest
	if err := r.db.WithContext(ctx).Where("verification_token = ?", tokenTransaction).First(&payment).Error; err != nil {
		return dto.CreatePaymentRequest{}, err
	}
	return payment, nil
}
