package repository

import (
	"context"

	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	GetAll(ctx context.Context) ([]entity.Transaction, error)
	GetById(ctx context.Context, id int64) (*entity.Transaction, error)
	GetByUserId(ctx context.Context, req dto.GetTransactionByUserIDRequest) ([]entity.Transaction, error)
	Create(ctx context.Context, transaction *entity.Transaction) error
	Update(ctx context.Context, transaction *entity.Transaction) error
	GetByOrderID(ctx context.Context, orderID string) (*entity.Transaction, error)
	GetOrderIdByToken(ctx context.Context, token string) (string, error)
	GetTransactionByToken(ctx context.Context, token string) (*entity.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) GetAll(ctx context.Context) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.db.WithContext(ctx).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) GetById(ctx context.Context, id int64) (*entity.Transaction, error) {
	result := new(entity.Transaction)
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *transactionRepository) GetByUserId(ctx context.Context, req dto.GetTransactionByUserIDRequest) ([]entity.Transaction, error) {
	transactions := make([]entity.Transaction, 0)
	query := r.db.WithContext(ctx).Where("user_id = ?", req.UserID)
	if req.Order != "" {
		query = query.Order("created_at " + req.Order)
	}
	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	return r.db.WithContext(ctx).Create(&transaction).Error
}

func (r *transactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	return r.db.WithContext(ctx).Updates(&transaction).Error
}

func (r *transactionRepository) GetByOrderID(ctx context.Context, orderID string) (*entity.Transaction, error) {
	result := new(entity.Transaction)
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *transactionRepository) GetOrderIdByToken(ctx context.Context, token string) (string, error) {
	result := new(entity.Transaction)
	if err := r.db.WithContext(ctx).Where("verification_token = ?", token).First(&result).Error; err != nil {
		return "", err
	}
	return result.OrderID, nil
}

func (r *transactionRepository) GetTransactionByToken(ctx context.Context, token string) (*entity.Transaction, error) {
	result := new(entity.Transaction)
	if err := r.db.WithContext(ctx).Where("verification_token = ?", token).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
