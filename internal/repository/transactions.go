package repository

import (
	"context"

	"github.com/mhusainh/FastTix/internal/entity"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
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
