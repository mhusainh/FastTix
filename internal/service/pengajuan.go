package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mhusainh/FastTix/config"
	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gopkg.in/gomail.v2"
)

type PengajuanService interface {
	CreatePengajuan(ctx context.Context, req dto.CreateProductRequest) (string, error)
	HandleMidtransNotification(ctx context.Context, notif map[string]interface{}) error
	ApprovePengajuan(ctx context.Context, productID int64) error
	RejectPengajuan(ctx context.Context, productID int64) error
}

type pengajuanService struct {
	cfg                   *config.Config
	productRepository     repository.ProductRepository
	userRepository        repository.UserRepository
	transactionRepository repository.TransactionRepository
}

func NewPengajuanService(cfg *config.Config, productRepo repository.ProductRepository, userRepo repository.UserRepository, transactionRepo repository.TransactionRepository) PengajuanService {
	return &pengajuanService{
		cfg:                   cfg,
		productRepository:     productRepo,
		userRepository:        userRepo,
		transactionRepository: transactionRepo,
	}
}

func (s *pengajuanService) CreatePengajuan(ctx context.Context, req dto.CreateProductRequest) (string, error) {
	status := "pending"
	if req.ProductPrice > 0 {
		status = "unpaid"
	}
	product := &entity.Product{
		ProductName:        req.ProductName,
		ProductAddress:     req.ProductAddress,
		ProductTime:        req.ProductTime,
		ProductDate:        req.ProductDate,
		ProductPrice:       req.ProductPrice,
		ProductDescription: req.ProductDescription,
		ProductCategory:    req.ProductCategory,
		ProductQuantity:    req.ProductQuantity,
		ProductType:        "available",
		ProductStatus:      status,
		UserID:             req.UserID,
	}
	err := s.productRepository.Create(ctx, product)
	if err != nil {
		return "", err
	}

	var redirectURL string
	if product.ProductStatus == "unpaid" {
		user, err := s.userRepository.GetById(ctx, product.UserID)
		if err != nil {
			return "", err
		}
		// Buat order_id
		orderID := fmt.Sprintf("order-%d-%d", product.ID, time.Now().Unix())

		// Buat transaksi pending di DB
		trans := &entity.Transaction{
			OrderID:             orderID,
			TransactionStatus:   "pending",
			ProductID:           product.ID,
			UserID:              product.UserID,
			TransactionQuantity: 1,
			TransactionAmount:   product.ProductPrice,
		}
		if err := s.transactionRepository.Create(ctx, trans); err != nil {
			return "", err
		}

		// Buat transaksi di Midtrans
		redirectURL, err = s.createMidtransTransaction(product, user, orderID)
		if err != nil {
			return "", err
		}

		err = s.sendPaymentConfirmationEmail(user.Email, product, redirectURL)
		if err != nil {
			return "", err
		}
	}

	return redirectURL, nil
}

// Buat request transaksi ke Midtrans
func (s *pengajuanService) createMidtransTransaction(product *entity.Product, user *entity.User, orderID string) (string, error) {
	midclient := snap.Client{}
	env := midtrans.Sandbox

	// Pilih environment berdasarkan kebutuhan
	if s.cfg.MidtransConfig.Environment == "production" {
		env = midtrans.Production
	}

	// Masukkan ServerKey langsung di sini (contoh: "Your-Midtrans-ServerKey")
	serverKey := "SB-Mid-server-HCftsNlI4uRULWqCQkxLGbqJ"

	midclient.New(serverKey, env)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(product.ProductPrice),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			Email: user.Email,
			FName: user.FullName,
		},
		EnabledPayments: []snap.SnapPaymentType{
			snap.PaymentTypeGopay,
			snap.PaymentTypeCreditCard,
			snap.PaymentTypeBankTransfer,
		},
	}

	snapResp, err := midclient.CreateTransaction(req)
	if err != nil {
		return "", err
	}
	return snapResp.RedirectURL, nil
}

// Handle notifikasi dari Midtrans
func (s *pengajuanService) HandleMidtransNotification(ctx context.Context, notif map[string]interface{}) error {
	orderID, ok := notif["order_id"].(string)
	if !ok {
		return errors.New("no order_id in notification")
	}
	transactionStatus, ok := notif["transaction_status"].(string)
	if !ok {
		return errors.New("no transaction_status in notification")
	}

	trans, err := s.transactionRepository.GetByOrderID(ctx, orderID)
	if err != nil {
		return err
	}
	product, err := s.productRepository.GetById(ctx, trans.ProductID)
	if err != nil {
		return err
	}
	user, err := s.userRepository.GetById(ctx, product.UserID)
	if err != nil {
		return err
	}

	if transactionStatus == "capture" || transactionStatus == "settlement" || transactionStatus == "success" {
		// Pembayaran berhasil
		product.ProductStatus = "pending"
		if err := s.productRepository.Update(ctx, product); err != nil {
			return err
		}
		trans.TransactionStatus = "success"
		if err := s.transactionRepository.Update(ctx, trans); err != nil {
			return err
		}

		// Kirim email tiket pengajuan
		if err := s.sendSubmissionTicketEmail(user.Email, product); err != nil {
			return err
		}
	} else if transactionStatus == "deny" || transactionStatus == "cancel" || transactionStatus == "expire" {
		// Pembayaran gagal
		trans.TransactionStatus = "failed"
		if err := s.transactionRepository.Update(ctx, trans); err != nil {
			return err
		}
		// Di sini bisa kirim email bahwa pembayaran gagal, jika diinginkan
	}

	return nil
}

func (s *pengajuanService) ApprovePengajuan(ctx context.Context, productID int64) error {
	product, err := s.productRepository.GetById(ctx, productID)
	if err != nil {
		return errors.New("Pengajuan tidak ditemukan")
	}

	product.ProductStatus = "accepted"
	if err := s.productRepository.Update(ctx, product); err != nil {
		return err
	}

	user, err := s.userRepository.GetById(ctx, product.UserID)
	if err != nil {
		return err
	}

	// Kirim email persetujuan
	return s.sendApprovalEmail(user.Email, product)
}

func (s *pengajuanService) RejectPengajuan(ctx context.Context, productID int64) error {
	product, err := s.productRepository.GetById(ctx, productID)
	if err != nil {
		return errors.New("Pengajuan tidak ditemukan")
	}

	product.ProductStatus = "unaccepted"
	if err := s.productRepository.Update(ctx, product); err != nil {
		return err
	}

	user, err := s.userRepository.GetById(ctx, product.UserID)
	if err != nil {
		return err
	}

	return s.sendRejectionEmail(user.Email, product)
}

// Fungsi bantu kirim email
func (s *pengajuanService) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPConfig.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(s.cfg.SMTPConfig.Host, s.cfg.SMTPConfig.Port, s.cfg.SMTPConfig.Username, s.cfg.SMTPConfig.Password)
	return d.DialAndSend(m)
}

func (s *pengajuanService) sendPaymentConfirmationEmail(email string, product *entity.Product, redirectURL string) error {
	body := fmt.Sprintf("Halo, Anda mengajukan tiket '%s' seharga %.2f.\nSilakan lakukan pembayaran di link berikut:\n%s", product.ProductName, product.ProductPrice, redirectURL)
	return s.sendEmail(email, "Konfirmasi Pembayaran Pengajuan", body)
}

func (s *pengajuanService) sendSubmissionTicketEmail(email string, product *entity.Product) error {
	body := fmt.Sprintf("Halo, pembayaran untuk pengajuan tiket '%s' telah kami terima. Pengajuan Anda sedang direview oleh admin.", product.ProductName)
	return s.sendEmail(email, "Tiket Pengajuan Anda", body)
}

func (s *pengajuanService) sendApprovalEmail(email string, product *entity.Product) error {
	body := fmt.Sprintf("Selamat, pengajuan tiket '%s' Anda disetujui oleh admin. Berikut tiket final Anda!", product.ProductName)
	return s.sendEmail(email, "Pengajuan Disetujui", body)
}

func (s *pengajuanService) sendRejectionEmail(email string, product *entity.Product) error {
	body := fmt.Sprintf("Maaf, pengajuan tiket '%s' Anda ditolak oleh admin. Silakan hubungi kami untuk info lebih lanjut.", product.ProductName)
	return s.sendEmail(email, "Pengajuan Ditolak", body)
}
