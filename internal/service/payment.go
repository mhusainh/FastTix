// internal/service/payment_service.go
package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"text/template"

	"github.com/mhusainh/FastTix/config"
	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gopkg.in/gomail.v2"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, req dto.CreatePaymentRequest) (string, error)
	VerifyPayment(ctx context.Context, req dto.VerifyPaymentRequest) error
	HandleMidtransNotification(ctx context.Context, notif map[string]interface{}) error
	CreateTokenPayment(ctx context.Context, req dto.CreatePaymentRequest) (string, error)
}

type paymentService struct {
	paymentRepository     repository.PaymentRepository
	config                *config.Config
	userRepository        repository.UserRepository
	transactionRepository repository.TransactionRepository
	productRepository     repository.ProductRepository
}

func NewPaymentService(paymentRepository repository.PaymentRepository, cfg *config.Config, userRepository repository.UserRepository, transactionRepository repository.TransactionRepository, productRepository repository.ProductRepository) PaymentService {
	return &paymentService{paymentRepository, cfg, userRepository, transactionRepository, productRepository}
}

func (s *paymentService) CreatePayment(ctx context.Context, req dto.CreatePaymentRequest) (string, error) {
	userID := req.UserID
	if userID == 0 {
		return "", errors.New("user id tidak ditemukan payment.go")
	}
	user, err := s.userRepository.GetById(ctx, userID)
	if err != nil {
		return "", err
	}
	fmt.Println(req.OrderID)
	fmt.Println(req.Amount)

	midclient := snap.Client{}
	midclient.New("SB-Mid-server-HCftsNlI4uRULWqCQkxLGbqJ", midtrans.Sandbox) // Use config for ServerKey
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
	snapResponse, err := midclient.CreateTransaction(request)
	if err != nil {
		return "", err
	}
	if err := s.sendPaymentConfirmationEmail(user.Email, req.NameProduct, snapResponse.RedirectURL); err != nil {
		log.Printf("ERROR: Failed to send payment confirmation email: %v", err)
		return "", errors.New("Failed to send payment confirmation email")
	}
	fmt.Println(snapResponse.RedirectURL)
	return snapResponse.RedirectURL, nil
}

func (s *paymentService) CreateTokenPayment(ctx context.Context, req dto.CreatePaymentRequest) (string, error) {
	// mau nambahkan verifikasi token
	userID := req.UserID
	if userID == 0 {
		return "", errors.New("user id tidak ditemukan payment.go")
	}
	user, err := s.userRepository.GetById(ctx, userID)
	if err != nil {
		return "", err
	}
	fmt.Println(req.OrderID)
	fmt.Println(req.Amount)

	templatePath := "./templates/email/verify-email.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "err", err
	}

	var ReplacerEmail = struct {
		Token string
	}{
		Token: req.VerificationToken,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, ReplacerEmail); err != nil {
		return "err", err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.SMTPConfig.Username)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Fast Tix : Verifikasi Email!")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(
		s.config.SMTPConfig.Host,
		s.config.SMTPConfig.Port,
		s.config.SMTPConfig.Username,
		s.config.SMTPConfig.Password,
	)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	// midclient := snap.Client{}
	// midclient.New("SB-Mid-server-HCftsNlI4uRULWqCQkxLGbqJ", midtrans.Sandbox) // Use config for ServerKey
	// request := &snap.Request{
	// 	TransactionDetails: midtrans.TransactionDetails{
	// 		OrderID:  req.OrderID,
	// 		GrossAmt: int64(req.Amount),
	// 	},
	// 	CustomerDetail: &midtrans.CustomerDetails{
	// 		FName: user.FullName,
	// 		Email: user.Email,
	// 	},
	// }
	// snapResponse, err := midclient.CreateTransaction(request)
	// if err != nil {
	// 	return "", err
	// }
	// if err := s.sendPaymentConfirmationEmail(user.Email, req.NameProduct, snapResponse.RedirectURL); err != nil {
	// 	log.Printf("ERROR: Failed to send payment confirmation email: %v", err)
	// 	return "", errors.New("Failed to send payment confirmation email")
	// }
	// fmt.Println(snapResponse.RedirectURL)
	// return snapResponse.RedirectURL, nil
	return "", nil
}

func (s *paymentService) VerifyPayment(ctx context.Context, req dto.VerifyPaymentRequest) error {
	// Implement payment verification logic if needed
	return nil
}

func (s *paymentService) sendPaymentConfirmationEmail(email string, NameProduct string, paymentLink string) error {
	subject := "Payment Confirmation for Your Ticket Purchase"
	body := fmt.Sprintf(`
        <h1>Payment Required</h1>
        <p>Thank you for purchasing a ticket for <strong>%s</strong>.</p>
        <p>Please complete your payment by clicking the link below:</p>
        <a href="%s">Pay Now</a>
        <p>If you did not initiate this purchase, please ignore this email.</p>
    `, NameProduct, paymentLink)

	return s.sendEmail(email, subject, body)
}

func (s *paymentService) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.SMTPConfig.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.config.SMTPConfig.Host, s.config.SMTPConfig.Port, s.config.SMTPConfig.Username, s.config.SMTPConfig.Password)
	return d.DialAndSend(m)
}

func (s *paymentService) HandleMidtransNotification(ctx context.Context, notif map[string]interface{}) error {
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
		trans.TransactionStatus = "success"
		if err := s.transactionRepository.Update(ctx, trans); err != nil {
			return err
		}

		// Send success email
		if err := s.sendPaymentSuccessEmail(user.Email, product); err != nil {
			return err
		}

		// Send ticket email
		if err := s.sendTicketEmail(user.Email, product, user.FullName); err != nil {
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

func (s *paymentService) sendPaymentSuccessEmail(email string, product *entity.Product) error {
	subject := "Payment Successful for Your Ticket Purchase"
	body := fmt.Sprintf(`
        <h1>Payment Successful</h1>
        <p>Thank you for your payment for the ticket to <strong>%s</strong>.</p>
        <p>Your purchase has been confirmed. Enjoy the event!</p>
        <p>If you have any questions, feel free to contact us.</p>
    `, product.ProductName)

	return s.sendEmail(email, subject, body)
}
func (s *paymentService) sendTicketEmail(email string, product *entity.Product, fullName string) error {
	subject := "Your Ticket for " + product.ProductName
	body := fmt.Sprintf(`
        <h1>Your Ticket</h1>
        <p>Thank you for your purchase, %s!</p>
        <p><strong>Event:</strong> %s</p>
        <p><strong>Date:</strong> %s</p>
        <p><strong>Time:</strong> %s</p>
        <p>Please present this ticket at the event entrance.</p>
    `, fullName, product.ProductName, product.ProductDate, product.ProductTime) // Use req.Category if available

	return s.sendEmail(email, subject, body)
}
