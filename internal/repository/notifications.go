package repository

import (
	"context"

	"github.com/mhusainh/FastTix/internal/entity"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) error
	GetByID(ctx context.Context, id int64) (*entity.Notification, error)
	GetByUserID(ctx context.Context, id int64) ([]entity.Notification, error)
	Update(ctx context.Context, notification *entity.Notification) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *entity.Notification) error {
	return r.db.WithContext(ctx).Create(&notification).Error
}

func (r *notificationRepository) GetByID(ctx context.Context, id int64) (*entity.Notification, error) {
	result := new(entity.Notification)
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *notificationRepository) GetByUserID(ctx context.Context, id int64) ([]entity.Notification, error) {
	result := make([]entity.Notification, 0)
	if err := r.db.WithContext(ctx).Where("user_id = ?", id).Order("created_at DESC").Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *notificationRepository) Update(ctx context.Context, notification *entity.Notification) error {
	return r.db.WithContext(ctx).Updates(&notification).Error
}
