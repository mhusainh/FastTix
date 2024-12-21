package service

import (
	"context"

	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
)

type NotificationService interface {
	SendNotificationSubmission(ctx context.Context, req dto.CreateNotificationRequest, submussion *entity.Product) error
	SendNotificationTransaction(ctx context.Context, req dto.CreateNotificationRequest, product *entity.Product, transaction *entity.Transaction) error
	GetByID(ctx context.Context, id int64) (*entity.Notification, error)
	GetByUserID(ctx context.Context, id int64) ([]entity.Notification, error)
	MarkAsRead(ctx context.Context, notification *entity.Notification) error
}

type notificationService struct {
	notificationRepository repository.NotificationRepository
}

func NewNotificationService(notificationRepository repository.NotificationRepository) NotificationService {
	return &notificationService{notificationRepository}
}

func (s *notificationService) SendNotificationSubmission(ctx context.Context, req dto.CreateNotificationRequest, submussion *entity.Product) error {
	if req.Message == "create" {
		req.Message = "Anda telah membuat pengajuan tiket " + submussion.ProductName
	}
	if req.Message == "update" {
		req.Message = "Anda telah mengubah pengajuan tiket " + submussion.ProductName
	}
	if req.Message == "accept" {
		req.Message = "Pengajuan tiket " + submussion.ProductName + " telah diterima oleh admin"
	}
	if req.Message == "reject" {
		req.Message = "Pengajuan tiket " + submussion.ProductName + " telah ditolak oleh admin"
	}
	if req.Message == "delete" {
		req.Message = "Anda telah membatalkan pengajuan tiket " + submussion.ProductName
	}
	notification := &entity.Notification{
		Message: req.Message,
		IsRead:  0,
		UserID:  submussion.UserID,
	}
	return s.notificationRepository.Create(ctx, notification)
}

func (s *notificationService) SendNotificationTransaction(ctx context.Context, req dto.CreateNotificationRequest, product *entity.Product, transaction *entity.Transaction) error {
	if req.Message == "Checkout Ticket" {
		req.Message = "Anda telah melakukan checkout tiket " + product.ProductName + " dengan order " + transaction.OrderID
	} else {
		req.Message = "Anda telah melakukan pembayaran tagihan dengan order " + transaction.OrderID
	}
	notification := &entity.Notification{
		Message: req.Message,
		IsRead:  0,
		UserID:  transaction.UserID,
	}
	return s.notificationRepository.Create(ctx, notification)
}

func (s *notificationService) GetByID(ctx context.Context, id int64) (*entity.Notification, error) {
	return s.notificationRepository.GetByID(ctx, id)
}

func (s *notificationService) GetByUserID(ctx context.Context, id int64) ([]entity.Notification, error) {
	return s.notificationRepository.GetByUserID(ctx, id)
}

func (s *notificationService) MarkAsRead(ctx context.Context, notification *entity.Notification) error {
	notification.IsRead = 1
	return s.notificationRepository.Update(ctx, notification)
}
