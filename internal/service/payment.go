// internal/service/payment_service.go
package service

import (
	"bytes"
	"context" // Tambahkan ini
	"errors"
	"fmt"
	"image/png"
	"io"
	"log"
	"text/template" // Tambahkan ini // Tambahkan ini

	bc "github.com/boombuler/barcode"
	bcqr "github.com/boombuler/barcode/qr"
	"github.com/mhusainh/FastTix/config"
	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gopkg.in/gomail.v2"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, req dto.CreatePaymentRequest) (map[string]interface{}, error)
	VerifyPayment(ctx context.Context, req dto.VerifyPaymentRequest) error
	HandleMidtransNotification(ctx context.Context, notif map[string]interface{}) error
	CreateTokenPayment(ctx context.Context, req dto.CreatePaymentRequest) (string, error)
	HandleCheckinNotification(ctx context.Context, req dto.CheckinWebhook) error
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

func (s *paymentService) CreatePayment(ctx context.Context, req dto.CreatePaymentRequest) (map[string]interface{}, error) {
	userID := req.UserID
	if userID == 0 {
		return nil, errors.New("user id tidak ditemukan payment.go")
	}
	user, err := s.userRepository.GetById(ctx, userID)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	if err := s.sendPaymentConfirmationEmail(user.Email, req.NameProduct, snapResponse.RedirectURL); err != nil {
		log.Printf("ERROR: Failed to send payment confirmation email: %v", err)
		return nil, errors.New("Failed to send payment confirmation email")
	}

	fmt.Println(snapResponse.RedirectURL)

	// Return JSON response
	return map[string]interface{}{
		"message":     "Payment created successfully",
		"redirectURL": snapResponse.RedirectURL,
	}, nil
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

	templatePath := "./templates/email/verify-payment.html"
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
		if err := s.sendTicketEmail(user.Email, product, user.FullName, orderID); err != nil {
			return err
		}
		// Send success email
		if err := s.sendPaymentSuccessEmail(user.Email, product); err != nil {
			return err
		}

		// Send ticket email

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
        <h1>Payment Successfully Completed</h1>
        <p>Thank you for your payment for the ticket to <strong>%s</strong>.</p>
        <p>Your purchase has been confirmed. Enjoy the event!</p>
        <p>If you have any questions, feel free to contact us.</p>
    `, product.ProductName)

	return s.sendEmail(email, subject, body)
}
func (s *paymentService) sendTicketEmail(email string, product *entity.Product, fullName string, orderID string) error {
	templatePath := "./templates/email/ticket.html"
	log.Println("Parsing template:", templatePath)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("ERROR: Failed to parse template: %v", err)
		return err
	}

	qrData := fmt.Sprintf("localhost:8080/api/v1/qrcode?order_id=%s", orderID)
	qrBytes, err := s.generateQRCode(qrData)
	if err != nil {
		log.Printf("ERROR: Failed to generate QR code: %v", err)
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.SMTPConfig.Username)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Fast Tix : Ticket Purchased")

	// Define a unique Content-ID for the QR code
	contentID := "https://28b2-139-0-237-234.ngrok-free.app/api/v1/checkin/" // You can generate a unique ID if needed

	// Attach the QR code as an inline image
	m.Attach("qrcode.png",
		gomail.SetHeader(map[string][]string{
			"Content-ID":          {fmt.Sprintf("<%s%s>", contentID, orderID)},
			"Content-Disposition": {"inline; filename=\"qrcode.png\""},
		}),
		gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(qrBytes)
			return err
		}),
	)

	// Prepare data for the email template
	var ReplacerEmail = struct {
		ProductName     string
		ProductCategory string
		ProductDate     string
		ProductPrice    float64
		ProductTime     string
		ProductLocation string
		ProductID       int64
		Name            string
		QRCodeCID       string
	}{
		ProductName:     product.ProductName,
		ProductCategory: product.ProductCategory,
		ProductDate:     product.ProductDate,
		ProductPrice:    product.ProductPrice,
		ProductTime:     product.ProductTime,
		ProductLocation: product.ProductAddress,
		ProductID:       product.ID,
		Name:            fullName,
		QRCodeCID:       contentID,
	}

	var bodyBuffer bytes.Buffer
	log.Println("Executing template with data:", ReplacerEmail)
	if err := tmpl.Execute(&bodyBuffer, ReplacerEmail); err != nil {
		log.Printf("ERROR: Failed to execute template: %v", err)
		return err
	}

	emailBody := bodyBuffer.String()
	log.Println("Generated email body:", emailBody)

	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer(
		s.config.SMTPConfig.Host,
		s.config.SMTPConfig.Port,
		s.config.SMTPConfig.Username,
		s.config.SMTPConfig.Password,
	)

	log.Println("Sending email to:", email)
	if err := d.DialAndSend(m); err != nil {
		log.Printf("ERROR: Failed to send email: %v", err)
		return err
	}

	log.Println("Email sent successfully to:", email)
	return nil
}

// Fungsi untuk menghasilkan barcode qr code
// Fungsi untuk menghasilkan QR code
func (s *paymentService) generateQRCode(data string) ([]byte, error) {
	// Encode data into QR code
	qrInstance, err := bcqr.Encode(data, bcqr.M, bcqr.Auto)
	if err != nil {
		return nil, err
	}

	// Scale QR code
	qrInstance, err = bc.Scale(qrInstance, 200, 200)
	if err != nil {
		return nil, err
	}

	// Encode QR code to PNG
	var pngBuffer bytes.Buffer
	if err := png.Encode(&pngBuffer, qrInstance); err != nil {
		return nil, err
	}

	return pngBuffer.Bytes(), nil
}

func (s *paymentService) HandleCheckinNotification(ctx context.Context, req dto.CheckinWebhook) error {
	// Handle check-in notification
	trans, err := s.transactionRepository.GetByOrderID(ctx, req.OrderID)
	if err != nil {
		return err
	}
	if trans.TransactionStatus != "success" {
		return errors.New("transaction not success")
	}
	trans.CheckIn = 1
	if err := s.transactionRepository.Update(ctx, trans); err != nil {
		return err
	}
	return nil
}
