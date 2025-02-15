package service

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/mhusainh/FastTix/config"
	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
	"github.com/mhusainh/FastTix/utils"
	"gopkg.in/gomail.v2"
)

type TransactionService interface {
	GetAll(ctx context.Context) ([]entity.Transaction, error)
	GetById(ctx context.Context, id int64) (*entity.Transaction, error)
	GetByUserId(ctx context.Context, req dto.GetTransactionByUserIDRequest) ([]entity.Transaction, error)
	GetTransactionByToken(ctx context.Context, token string) (*entity.Transaction, error)
	GetByOrderID(ctx context.Context, orderID string) (*entity.Transaction, error)
	Create(ctx context.Context, req dto.CreateTransactionRequest, user *entity.User, product *entity.Product) (*entity.Transaction, error)
	PaymentTicket(ctx context.Context, req dto.UpdateTransactionRequest, user *entity.User, product *entity.Product, transaction *entity.Transaction) (*entity.Transaction, error)
	PaymentSubmission(ctx context.Context, req dto.UpdateTransactionRequest, user *entity.User, product *entity.Product, transaction *entity.Transaction) (*entity.Transaction, error)
}

type transactionService struct {
	cfg                   *config.Config
	transactionRepository repository.TransactionRepository
	productRepository     repository.ProductRepository
}

func NewTransactionService(
	cfg *config.Config,
	transactionRepository repository.TransactionRepository,
	productRepository repository.ProductRepository,
) TransactionService {
	return &transactionService{cfg, transactionRepository, productRepository}
}

func (s *transactionService) GetAll(ctx context.Context) ([]entity.Transaction, error) {
	return s.transactionRepository.GetAll(ctx)
}

func (s *transactionService) GetById(ctx context.Context, id int64) (*entity.Transaction, error) {
	return s.transactionRepository.GetById(ctx, id)
}

func (s *transactionService) GetByUserId(ctx context.Context, req dto.GetTransactionByUserIDRequest) ([]entity.Transaction, error) {
	return s.transactionRepository.GetByUserId(ctx, req)
}

func (s *transactionService) GetByOrderID(ctx context.Context, orderID string) (*entity.Transaction, error) {
	return s.transactionRepository.GetByOrderID(ctx, orderID)
}

func (s *transactionService) Create(ctx context.Context, req dto.CreateTransactionRequest, user *entity.User, product *entity.Product) (*entity.Transaction, error) {
	if product.ProductPrice == 0 {
		req.TransactionStatus = "success"
		templatePath := "./templates/email/ticket.html"
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			return nil, err
		}

		var ReplacerEmail = struct {
			Name    string
			Address string
			Time    string
			Date    string
			Price   float64
		}{
			Name:    product.ProductName,
			Address: product.ProductAddress,
			Time:    product.ProductTime,
			Date:    product.ProductDate,
			Price:   product.ProductPrice,
		}

		var body bytes.Buffer
		if err := tmpl.Execute(&body, ReplacerEmail); err != nil {
			return nil, err
		}

		m := gomail.NewMessage()
		m.SetHeader("From", s.cfg.SMTPConfig.Username)
		m.SetHeader("To", user.Email)
		m.SetHeader("Subject", "Fast Tix : Ticket "+product.ProductName+"!")
		m.SetBody("text/html", body.String())

		d := gomail.NewDialer(
			s.cfg.SMTPConfig.Host,
			s.cfg.SMTPConfig.Port,
			s.cfg.SMTPConfig.Username,
			s.cfg.SMTPConfig.Password,
		)

		// Send the email to Bob, Cora and Dan.
		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
		product.ProductQuantity -= req.TransactionQuantity
		err = s.productRepository.Update(ctx, product)
		if err != nil {
			return nil, err
		}
	} else {
		req.TransactionStatus = "pending"
	}
	req.TransactionAmount = product.ProductPrice * float64(req.TransactionQuantity)
	req.OrderID = fmt.Sprintf("order_id-%s", utils.RandomString(10))
	transaction := &entity.Transaction{
		TransactionAmount:   req.TransactionAmount,
		TransactionQuantity: req.TransactionQuantity,
		TransactionStatus:   req.TransactionStatus,
		TransactionType:     "ticket",
		VerificationToken:   req.VerificationToken,
		OrderID:             req.OrderID,
		UserID:              user.ID,
		ProductID:           product.ID,
	}

	return transaction, s.transactionRepository.Create(ctx, transaction)
}

func (s *transactionService) PaymentTicket(ctx context.Context, req dto.UpdateTransactionRequest, user *entity.User, product *entity.Product, transaction *entity.Transaction) (*entity.Transaction, error) {
	templatePath := "./templates/email/ticket.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	var ReplacerEmail = struct {
		Name    string
		Address string
		Time    string
		Date    string
		Price   float64
	}{
		Name:    product.ProductName,
		Address: product.ProductAddress,
		Time:    product.ProductTime,
		Date:    product.ProductDate,
		Price:   product.ProductPrice,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, ReplacerEmail); err != nil {
		return nil, err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPConfig.Username)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Fast Tix : Ticket "+product.ProductName+"!")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(
		s.cfg.SMTPConfig.Host,
		s.cfg.SMTPConfig.Port,
		s.cfg.SMTPConfig.Username,
		s.cfg.SMTPConfig.Password,
	)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	transaction.TransactionStatus = "success"
	return transaction, s.transactionRepository.Update(ctx, transaction)
}

func (s *transactionService) PaymentSubmission(ctx context.Context, req dto.UpdateTransactionRequest, user *entity.User, product *entity.Product, transaction *entity.Transaction) (*entity.Transaction, error) {
	templatePath := "./templates/email/notif-submission.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, nil); err != nil {
		return nil, err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPConfig.Username)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Fast Tix : Submission Event!")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(
		s.cfg.SMTPConfig.Host,
		s.cfg.SMTPConfig.Port,
		s.cfg.SMTPConfig.Username,
		s.cfg.SMTPConfig.Password,
	)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	transaction.TransactionStatus = "success"
	return transaction, s.transactionRepository.Update(ctx, transaction)
}

func (s *transactionService) GetTransactionByToken(ctx context.Context, token string) (*entity.Transaction, error) {
	return s.transactionRepository.GetTransactionByToken(ctx, token)
}
