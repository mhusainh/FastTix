package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mhusainh/FastTix/config"
	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gopkg.in/gomail.v2"
)

type PurchaseService interface {
	PurchaseTicket(ctx context.Context, userID int64, req dto.PurchaseTicketRequest) (*dto.PurchaseTicketResponse, error)
	CheckPurchaseStatus(ctx context.Context, orderID string) (string, error)
	HandleMidtransNotification(ctx context.Context, notif map[string]interface{}) error
}

type purchaseService struct {
	cfg                   *config.Config
	productRepository     repository.ProductRepository
	userRepository        repository.UserRepository
	transactionRepository repository.TransactionRepository
	tokenService          TokenService
}

func NewPurchaseService(cfg *config.Config, productRepo repository.ProductRepository, userRepo repository.UserRepository, transactionRepo repository.TransactionRepository, tokenService TokenService) PurchaseService {
	return &purchaseService{
		cfg:                   cfg,
		productRepository:     productRepo,
		userRepository:        userRepo,
		transactionRepository: transactionRepo,
		tokenService:          tokenService,
	}
}

func (s *purchaseService) PurchaseTicket(ctx context.Context, userID int64, req dto.PurchaseTicketRequest) (*dto.PurchaseTicketResponse, error) {
	// Fetch the product
	product, err := s.productRepository.GetById(ctx, req.ProductID)
	if err != nil {
		log.Printf("ERROR: Failed to fetch product with ID %d: %v", req.ProductID, err)
		return nil, errors.New("Product not found")
	}

	// Log product and request details
	log.Printf("DEBUG: Product Date: %s, Selected Date: %s", product.ProductDate, req.SelectedDate)
	log.Printf("DEBUG: Product Time: %s, Selected Time: %s", product.ProductTime, req.SelectedTime)

	// Define date and time formats
	const productDateLayout = time.RFC3339  // Expected format for product.ProductDate
	const selectedDateLayout = "2006-01-02" // Format for req.SelectedDate
	const productTimeLayout = "15:04:05"    // Format for product.ProductTime
	const selectedTimeLayout = "15:04:05"   // Format for req.SelectedTime

	// Parse and format product date
	parsedProductDate, err := time.Parse(productDateLayout, product.ProductDate)
	if err != nil {
		log.Printf("ERROR: Invalid product date format for product ID %d: %v", req.ProductID, err)
		return nil, errors.New("Invalid product date format")
	}
	formattedProductDate := parsedProductDate.Format(selectedDateLayout)
	log.Printf("DEBUG: Formatted Product Date: %s, Selected Date: %s", formattedProductDate, req.SelectedDate)

	// Compare formatted date with selected date
	if formattedProductDate != req.SelectedDate {
		log.Printf("ERROR: Date mismatch. Product Date: %s, Selected Date: %s", formattedProductDate, req.SelectedDate)
		return nil, errors.New("Invalid date selection")
	}

	// Parse and format product time
	parsedProductTime, err := time.Parse(productTimeLayout, product.ProductTime)
	if err != nil {
		log.Printf("ERROR: Invalid product time format for product ID %d: %v", req.ProductID, err)
		return nil, errors.New("Invalid product time format")
	}
	formattedProductTime := parsedProductTime.Format(selectedTimeLayout)
	log.Printf("DEBUG: Formatted Product Time: %s, Selected Time: %s", formattedProductTime, req.SelectedTime)

	// Validate time
	if formattedProductTime != req.SelectedTime {
		log.Printf("ERROR: Time mismatch. Product Time: %s, Selected Time: %s", formattedProductTime, req.SelectedTime)
		return nil, errors.New("Invalid time selection")
	}
	if req.Category != product.ProductCategory {
		log.Printf("ERROR: Category mismatch. Product Category: %s, Selected Category: %s", product.ProductCategory, req.Category)
		return nil, errors.New("Invalid category selection")
	}

	// Create a transaction entry
	transaction := &entity.Transaction{
		OrderID:             fmt.Sprintf("order_id-%d-%d", product.ID, time.Now().Unix()), // Updated prefix
		TransactionStatus:   "pending",
		ProductID:           product.ID,
		UserID:              userID,
		TransactionQuantity: 1,
		TransactionAmount:   product.ProductPrice,
	}

	if err := s.transactionRepository.Create(ctx, transaction); err != nil {
		log.Printf("ERROR: Failed to create transaction for product ID %d: %v", req.ProductID, err)
		return nil, errors.New("Failed to create transaction")
	}

	// Handle payment if required
	if product.ProductPrice > 0 {
		// Fetch user details
		user, err := s.userRepository.GetById(ctx, userID)
		if err != nil {
			log.Printf("ERROR: Failed to fetch user with ID %d: %v", userID, err)
			return nil, errors.New("User not found")
		}

		paymentLink, err := s.createMidtransTransaction(ctx, product, transaction.OrderID, user)
		if err != nil {
			log.Printf("ERROR: Failed to create Midtrans transaction for product ID %d: %v", req.ProductID, err)
			return nil, errors.New("Failed to initiate payment")
		}

		if err := s.sendPaymentConfirmationEmail(user.Email, product, paymentLink); err != nil {
			log.Printf("ERROR: Failed to send payment confirmation email: %v", err)
			return nil, errors.New("Failed to send payment confirmation email")
		}

		return &dto.PurchaseTicketResponse{
			Message:     "Payment required. A confirmation email has been sent.",
			PaymentLink: paymentLink,
		}, nil
	}

	// Send ticket directly if payment is not required
	user, err := s.userRepository.GetById(ctx, userID)
	if err != nil {
		log.Printf("ERROR: Failed to fetch user with ID %d: %v", userID, err)
		return nil, errors.New("User not found")
	}

	if err := s.sendTicketEmail(user.Email, product, user.FullName); err != nil {
		log.Printf("ERROR: Failed to send ticket email: %v", err)
		return nil, errors.New("Failed to send ticket email")
	}

	return &dto.PurchaseTicketResponse{
		Message: "Ticket successfully purchased. A ticket has been sent to your email.",
	}, nil
}

// internal/service/purchase.go

func (s *purchaseService) createMidtransTransaction(ctx context.Context, product *entity.Product, orderID string, user *entity.User) (string, error) {
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

func (s *purchaseService) HandleMidtransNotification(ctx context.Context, notif map[string]interface{}) error {
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

func (s *purchaseService) sendPaymentVerificationEmail(email string, product *entity.Product, paymentLink string) error {
	subject := "Payment Verification for Your Ticket Purchase"
	body := fmt.Sprintf(`
        <h1>Payment Required</h1>
        <p>Thank you for purchasing a ticket for <strong>%s</strong>.</p>
        <p>Please complete your payment by clicking the link below:</p>
        <a href="%s">Pay Now</a>
        <p>If you did not initiate this purchase, please ignore this email.</p>
    `, product.ProductName, paymentLink)

	return s.sendEmail(email, subject, body)
}

// internal/service/purchase.go

func (s *purchaseService) sendPaymentConfirmationEmail(email string, product *entity.Product, paymentLink string) error {
	subject := "Payment Confirmation for Your Ticket Purchase"
	body := fmt.Sprintf(`
        <h1>Payment Required</h1>
        <p>Thank you for purchasing a ticket for <strong>%s</strong>.</p>
        <p>Please complete your payment by clicking the link below:</p>
        <a href="%s">Pay Now</a>
        <p>If you did not initiate this purchase, please ignore this email.</p>
    `, product.ProductName, paymentLink)

	return s.sendEmail(email, subject, body)
}

// internal/service/purchase.go

func (s *purchaseService) sendPaymentSuccessEmail(email string, product *entity.Product) error {
	subject := "Payment Successful for Your Ticket Purchase"
	body := fmt.Sprintf(`
        <h1>Payment Successful</h1>
        <p>Thank you for your payment for the ticket to <strong>%s</strong>.</p>
        <p>Your purchase has been confirmed. Enjoy the event!</p>
        <p>If you have any questions, feel free to contact us.</p>
    `, product.ProductName)

	return s.sendEmail(email, subject, body)
}

// internal/service/purchase.go

func (s *purchaseService) sendTicketEmail(email string, product *entity.Product, fullName string) error {
	subject := "Your Ticket for " + product.ProductName
	body := fmt.Sprintf(`
        <h1>Your Ticket</h1>
        <p>Thank you for your purchase, %s!</p>
        <p><strong>Event:</strong> %s</p>
        <p><strong>Date:</strong> %s</p>
        <p><strong>Time:</strong> %s</p>
        <p><strong>Category:</strong> %s</p>
        <p>Please present this ticket at the event entrance.</p>
    `, fullName, product.ProductName, product.ProductDate, product.ProductTime, product.ProductCategory)

	return s.sendEmail(email, subject, body)
}

func (s *purchaseService) CheckPurchaseStatus(ctx context.Context, orderID string) (string, error) {
	trans, err := s.transactionRepository.GetByOrderID(ctx, orderID)
	if err != nil {
		return "", errors.New("Transaction not found")
	}
	return trans.TransactionStatus, nil
}
func (s *purchaseService) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPConfig.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.cfg.SMTPConfig.Host, s.cfg.SMTPConfig.Port, s.cfg.SMTPConfig.Username, s.cfg.SMTPConfig.Password)
	return d.DialAndSend(m)
}
