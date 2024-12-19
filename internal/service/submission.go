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

type SubmissionService interface {
	CreateSubmission(ctx context.Context, req dto.CreateProductRequest) (string, error)
	HandleMidtransNotification(ctx context.Context, notif map[string]interface{}) error
	ApproveSubmission(ctx context.Context, productID int64) error
	RejectSubmission(ctx context.Context, productID int64) error
}

type submissionService struct {
	cfg                   *config.Config
	productRepository     repository.ProductRepository
	userRepository        repository.UserRepository
	transactionRepository repository.TransactionRepository
}

func NewSubmissionService(cfg *config.Config, productRepo repository.ProductRepository, userRepo repository.UserRepository, transactionRepo repository.TransactionRepository) SubmissionService {
	return &submissionService{
		cfg:                   cfg,
		productRepository:     productRepo,
		userRepository:        userRepo,
		transactionRepository: transactionRepo,
	}
}

func (s *submissionService) CreateSubmission(ctx context.Context, req dto.CreateProductRequest) (string, error) {
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
		// Create order_id with "daftar_id" prefix
		orderID := fmt.Sprintf("daftar_id-%d-%d", product.ID, time.Now().Unix())

		// Create pending transaction in DB
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

		// Create transaction in Midtrans
		redirectURL, err = s.createMidtransTransaction(product, user, orderID)
		if err != nil {
			return "", err
		}

		// Send payment confirmation email
		err = s.sendPaymentConfirmationEmail(user.Email, product, redirectURL)
		if err != nil {
			return "", err
		}
	}

	return redirectURL, nil
}

// CreateMidtransTransaction creates a transaction with Midtrans
func (s *submissionService) createMidtransTransaction(product *entity.Product, user *entity.User, orderID string) (string, error) {
	midclient := snap.Client{}
	env := midtrans.Sandbox

	// Choose environment based on configuration
	if s.cfg.MidtransConfig.Environment == "production" {
		env = midtrans.Production
	}

	midclient.New("SB-Mid-server-HCftsNlI4uRULWqCQkxLGbqJ", env)

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

// HandleMidtransNotification processes Midtrans webhook notifications
func (s *submissionService) HandleMidtransNotification(ctx context.Context, notif map[string]interface{}) error {
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
	user, err := s.userRepository.GetById(ctx, trans.UserID)
	if err != nil {
		return err
	}

	if transactionStatus == "capture" || transactionStatus == "settlement" || transactionStatus == "success" {
		// Payment successful
		product.ProductStatus = "pending"
		if err := s.productRepository.Update(ctx, product); err != nil {
			return err
		}
		trans.TransactionStatus = "success"
		if err := s.transactionRepository.Update(ctx, trans); err != nil {
			return err
		}

		// Send submission ticket email
		if err := s.sendSubmissionTicketEmail(user.Email, product); err != nil {
			return err
		}
	} else if transactionStatus == "deny" || transactionStatus == "cancel" || transactionStatus == "expire" {
		// Payment failed
		trans.TransactionStatus = "failed"
		if err := s.transactionRepository.Update(ctx, trans); err != nil {
			return err
		}
		// Optionally, send a failure notification email
	}

	return nil
}

// ApproveSubmission allows admin to approve a submission
func (s *submissionService) ApproveSubmission(ctx context.Context, productID int64) error {
	product, err := s.productRepository.GetById(ctx, productID)
	if err != nil {
		return errors.New("Submission not found")
	}

	product.ProductStatus = "accepted"
	if err := s.productRepository.Update(ctx, product); err != nil {
		return err
	}

	user, err := s.userRepository.GetById(ctx, product.UserID)
	if err != nil {
		return err
	}

	// Send approval email
	return s.sendApprovalEmail(user.Email, product)
}

// RejectSubmission allows admin to reject a submission
func (s *submissionService) RejectSubmission(ctx context.Context, productID int64) error {
	product, err := s.productRepository.GetById(ctx, productID)
	if err != nil {
		return errors.New("Submission not found")
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

// Helper functions to send emails

func (s *submissionService) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPConfig.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.cfg.SMTPConfig.Host, s.cfg.SMTPConfig.Port, s.cfg.SMTPConfig.Username, s.cfg.SMTPConfig.Password)
	return d.DialAndSend(m)
}

func (s *submissionService) sendPaymentConfirmationEmail(email string, product *entity.Product, redirectURL string) error {
	subject := "Payment Confirmation for Your Submission"
	body := fmt.Sprintf(`
        <h1>Payment Required</h1>
        <p>Thank you for submitting a ticket for <strong>%s</strong>.</p>
        <p>Please complete your payment by clicking the link below:</p>
        <a href="%s">Pay Now</a>
        <p>If you did not initiate this submission, please ignore this email.</p>
    `, product.ProductName, redirectURL)

	return s.sendEmail(email, subject, body)
}

func (s *submissionService) sendSubmissionTicketEmail(email string, product *entity.Product) error {
	subject := "Your Submission Ticket"
	body := fmt.Sprintf(`
        <h1>Submission Successful</h1>
        <p>Your payment for the submission of ticket '<strong>%s</strong>' has been received and is under review by the admin.</p>
    `, product.ProductName)

	return s.sendEmail(email, subject, body)
}

func (s *submissionService) sendApprovalEmail(email string, product *entity.Product) error {
	subject := "Submission Approved"
	body := fmt.Sprintf(`
        <h1>Submission Approved</h1>
        <p>Congratulations! Your submission for ticket '<strong>%s</strong>' has been approved by the admin. Here is your final ticket.</p>
    `, product.ProductName)

	return s.sendEmail(email, subject, body)
}

func (s *submissionService) sendRejectionEmail(email string, product *entity.Product) error {
	subject := "Submission Rejected"
	body := fmt.Sprintf(`
        <h1>Submission Rejected</h1>
        <p>We regret to inform you that your submission for ticket '<strong>%s</strong>' has been rejected by the admin. Please contact us for more details.</p>
    `, product.ProductName)

	return s.sendEmail(email, subject, body)
}
