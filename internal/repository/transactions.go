package repository

import (
	"context"

	"github.com/mhusainh/FastTix/internal/entity"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	GetByOrderID(ctx context.Context, orderID string) (*entity.Transaction, error)
	Update(ctx context.Context, transaction *entity.Transaction) error
	
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	return r.db.WithContext(ctx).Create(&transaction).Error
}

func (r *transactionRepository) GetByOrderID(ctx context.Context, orderID string) (*entity.Transaction, error) {
	var trans entity.Transaction
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&trans).Error; err != nil {
		return nil, err
	}
	return &trans, nil
}

func (r *transactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}


