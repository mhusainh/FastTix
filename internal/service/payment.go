package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/mhusainh/FastTix/config"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, req dto.CreatePaymentRequest) (string, error)
	VerifyPayment(ctx context.Context, req dto.VerifyPaymentRequest) error
}

type paymentService struct {
	paymentRepository repository.PaymentRepository
	config            *config.Config
	userRepository    repository.UserRepository
}

func NewPaymentService(paymentRepository repository.PaymentRepository, cfg *config.Config, userRepository repository.UserRepository) PaymentService {
	return &paymentService{paymentRepository, cfg, userRepository}
}

func (s *paymentService) CreatePayment(ctx context.Context, req dto.CreatePaymentRequest) (string, error) {
	// details := s.paymentRepository.Create(ctx, req)
	// return details, nil
	userID := req.UserID
	if userID == 0 {
		return "", errors.New("user id tidak ditemukan")
	}
	user, err := s.userRepository.GetById(ctx, userID)
	if err != nil {
		return "", err
	}
	var snaps = snap.Client{}
	snaps.New("SB-Mid-server-HCftsNlI4uRULWqCQkxLGbqJ", midtrans.Sandbox)
	request := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.OrderID,
			GrossAmt: int64(req.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.FullName,
			Email: user.Email,
		},
	}
	snapResponse, err := snap.CreateTransaction(request)
	if err != nil {
		return "", err
	}
	fmt.Println(snapResponse.RedirectURL)
	return snapResponse.RedirectURL, nil
}

func (s *paymentService) VerifyPayment(ctx context.Context, req dto.VerifyPaymentRequest) error {
	return nil
}
