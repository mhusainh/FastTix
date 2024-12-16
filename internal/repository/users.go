package repository

import (
	"context"

	"github.com/mhusainh/FastTix/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]entity.User, error)
	GetById(ctx context.Context, id int64) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, user *entity.User) error
	GetByResetPasswordToken(ctx context.Context, token string) (*entity.User, error)
	GetByVerifyEmailToken(ctx context.Context, token string) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	result := make([]entity.User, 0)
	if err := r.db.WithContext(ctx).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *userRepository) GetById(ctx context.Context, id int64) (*entity.User, error) {
	result := new(entity.User)
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Updates(&user).Error
}

func (r *userRepository) Delete(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Delete(&user).Error
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	result := new(entity.User)
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *userRepository) GetByResetPasswordToken(ctx context.Context, token string) (*entity.User, error) {
	result := new(entity.User)
	if err := r.db.WithContext(ctx).Where("reset_password_token = ?", token).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *userRepository) GetByVerifyEmailToken(ctx context.Context, token string) (*entity.User, error) {
	result := new(entity.User)
	if err := r.db.WithContext(ctx).Where("verify_email_token = ?", token).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}