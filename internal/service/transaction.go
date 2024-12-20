package service

import (
	"bytes"
	"context"
	"errors"
	"text/template"

	"github.com/mhusainh/FastTix/config"
	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
	"gopkg.in/gomail.v2"
)

type TransactionService interface {
	GetAll(ctx context.Context) ([]entity.Transaction, error)
	GetById(ctx context.Context, id int64) (*entity.Transaction, error)
	GetByUserId(ctx context.Context, req dto.GetTransactionByUserIDRequest) ([]entity.Transaction, error)
	Create(ctx context.Context, req dto.CreateTransactionRequest) error
	PaymentTicket(ctx context.Context, req dto.UpdateTransactionRequest) error
	PaymentSubmission(ctx context.Context, req dto.UpdateTransactionRequest) error
}

type transactionService struct {
	cfg                   *config.Config
	transactionRepository repository.TransactionRepository
	userRepository        repository.UserRepository
	productRepository     repository.ProductRepository
}

func NewTransactionService(
	cfg *config.Config,
	transactionRepository repository.TransactionRepository,
	userRepository repository.UserRepository,
	productRepository repository.ProductRepository,
) TransactionService {
	return &transactionService{cfg, transactionRepository, userRepository, productRepository}
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

func (s *transactionService) Create(ctx context.Context, req dto.CreateTransactionRequest) error {
	userID := req.UserID
	if userID == 0 {
		return errors.New("user id tidak ditemukan")
	}
	user, err := s.userRepository.GetById(ctx, userID)
	if err != nil {
		return err
	}
	product, err := s.productRepository.GetById(ctx, req.ProductID)
	if err != nil {
		return err
	}

	if product.ProductPrice == 0 {
		req.TransactionStatus = "success"
		templatePath := "./templates/email/ticket.html"
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			return err
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
			return err
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
			return err
		}
	} else {
		req.TransactionStatus = "pending"
	}
	transaction := &entity.Transaction{
		TransactionAmount:   req.TransactionAmount,
		TransactionQuantity: req.TransactionQuantity,
		TransactionStatus:   req.TransactionStatus,
		UserID:              userID,
		ProductID:           product.ID,
	}

	return s.transactionRepository.Create(ctx, transaction)
}

func (s *transactionService) PaymentTicket(ctx context.Context, req dto.UpdateTransactionRequest) error {
	userID := req.UserID
	if userID == 0 {
		return errors.New("user id tidak ditemukan")
	}
	user, err := s.userRepository.GetById(ctx, userID)
	if err != nil {
		return err
	}
	transaction, err := s.transactionRepository.GetById(ctx, req.ID)
	if err != nil {
		return err
	}
	product, err := s.productRepository.GetById(ctx, transaction.ProductID)
	if err != nil {
		return err
	}

	product.ProductQuantity -= transaction.TransactionQuantity
	err = s.productRepository.Update(ctx, product)
	if err != nil {
		return err
	}
	templatePath := "./templates/email/ticket.html"
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			return err
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
			return err
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
	return s.transactionRepository.Update(ctx, transaction)
}

func (s *transactionService) PaymentSubmission(ctx context.Context, req dto.UpdateTransactionRequest) error {
	userID := req.UserID
	if userID == 0 {
		return errors.New("user id tidak ditemukan")
	}
	user, err := s.userRepository.GetById(ctx, userID)
	if err != nil {
		return err
	}
	transaction, err := s.transactionRepository.GetById(ctx, req.ID)
	if err != nil {
		return err
	}
	product, err := s.productRepository.GetById(ctx, transaction.ProductID)
	if err != nil {
		return err
	}
	product.ProductStatus = "pending"
	err = s.productRepository.Update(ctx, product)
	if err != nil {
		return err
	}
	templatePath := "./templates/email/notif-submission.html"
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			return err
		}

		var body bytes.Buffer
		if err := tmpl.Execute(&body, nil); err != nil {
			return err
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
	return s.transactionRepository.Update(ctx, transaction)
}
