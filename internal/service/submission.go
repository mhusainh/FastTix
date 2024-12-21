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

type SubmissionService interface {
	GetAll(ctx context.Context, req dto.GetAllProductsRequest) ([]entity.Product, error)
	GetById(ctx context.Context, id int64) (*entity.Product, error)
	Create(ctx context.Context, req dto.CreateProductRequest, t dto.CreateTransactionRequest, user *entity.User) (*entity.Product, error)
	UpdateByUser(ctx context.Context, req dto.UpdateProductRequest, user *entity.User, submission *entity.Product) (*entity.Product, error)
	Approve(ctx context.Context, submission *entity.Product) (*entity.Product, error)
	Reject(ctx context.Context, submission *entity.Product) (*entity.Product, error)
	Cancel(ctx context.Context, submission *entity.Product, req dto.GetProductByIDRequest) error
}

type submissionService struct {
	cfg                   *config.Config
	submissionRepository  repository.SubmissionRepository
	transactionRepository repository.TransactionRepository
	productRepository     repository.ProductRepository
}

func NewSubmissionService(
	cfg *config.Config,
	submissionRepository repository.SubmissionRepository,
	transactionRepository repository.TransactionRepository,
	productRepository repository.ProductRepository,
) SubmissionService {
	return &submissionService{cfg, submissionRepository, transactionRepository, productRepository}
}

func (s *submissionService) GetAll(ctx context.Context, req dto.GetAllProductsRequest) ([]entity.Product, error) {
	return s.submissionRepository.GetAll(ctx, req)
}

func (s *submissionService) GetById(ctx context.Context, id int64) (*entity.Product, error) {
	return s.submissionRepository.GetById(ctx, id)
}

func (s *submissionService) Create(ctx context.Context, req dto.CreateProductRequest, t dto.CreateTransactionRequest, user *entity.User) (*entity.Product, error) {
	exist, err := s.productRepository.GetByName(ctx, req.ProductName)
	if err == nil && exist != nil {
		return nil, errors.New("Nama event sudah digunakan")
	}
	if req.ProductPrice == 0 {
		req.ProductStatus = "pending"
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
	} else {
		req.ProductStatus = "unpaid"
	}

	submission := &entity.Product{
		ProductName:        req.ProductName,
		ProductAddress:     req.ProductAddress,
		ProductTime:        req.ProductTime,
		ProductDate:        req.ProductDate,
		ProductPrice:       req.ProductPrice,
		ProductDescription: req.ProductDescription,
		ProductCategory:    req.ProductCategory,
		ProductQuantity:    req.ProductQuantity,
		ProductType:        "available",
		ProductStatus:      req.ProductStatus,
		UserID:             user.ID,
	}

	if err := s.submissionRepository.Create(ctx, submission); err != nil {
		return nil, err
	}

	ProductID := submission.ID
	TransactionAmount := submission.ProductPrice * 0.25

	if TransactionAmount != 0 {
		transaction := &entity.Transaction{
			ProductID:           ProductID,
			UserID:              user.ID,
			TransactionQuantity: 1,
			TransactionAmount:   TransactionAmount,
			TransactionStatus:   "pending",
			TransactionType:     "submission",
		}
		if err := s.transactionRepository.Create(ctx, transaction); err != nil {
			return nil,err
		}
	}

	return submission, nil

}

func (s *submissionService) UpdateByUser(ctx context.Context, req dto.UpdateProductRequest, user *entity.User, submission *entity.Product) (*entity.Product, error) {
	if req.ProductName != "" {
		submission.ProductName = req.ProductName
	}
	if req.ProductAddress != "" {
		submission.ProductAddress = req.ProductAddress
	}
	if req.ProductTime != "" {
		submission.ProductTime = req.ProductTime
	}
	if req.ProductDate != "" {
		submission.ProductDate = req.ProductDate
	}
	if req.ProductDescription != "" {
		submission.ProductDescription = req.ProductDescription
	}
	if req.ProductCategory != "" {
		submission.ProductCategory = req.ProductCategory
	}
	if req.ProductQuantity != 0 {
		submission.ProductQuantity = req.ProductQuantity
	}
	if req.ProductType != "" {
		submission.ProductType = req.ProductType
	}
	return submission, s.submissionRepository.Update(ctx, submission)
}

func (s *submissionService) Approve(ctx context.Context, submission *entity.Product) (*entity.Product, error) {
	submission.ProductStatus = "accepted"
	return submission, s.submissionRepository.Update(ctx, submission)
}

func (s *submissionService) Reject(ctx context.Context, submission *entity.Product) (*entity.Product, error) {
	submission.ProductStatus = "rejected"
	return submission, s.submissionRepository.Update(ctx, submission)
}

func (s *submissionService) Cancel(ctx context.Context, submission *entity.Product, req dto.GetProductByIDRequest) error {
	return s.submissionRepository.Delete(ctx, submission)
}
