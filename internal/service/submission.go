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
	Create(ctx context.Context, req dto.CreateProductRequest, t dto.CreateTransactionRequest) error
	UpdateByUser(ctx context.Context, req dto.UpdateProductRequest) error
	Approve(ctx context.Context, id int64) error
	Reject(ctx context.Context, id int64) error
	Cancel(ctx context.Context, submission *entity.Product, req dto.GetProductByIDRequest) error
	HandleMidtransNotification(ctx context.Context, notif map[string]interface{}) error
}

type submissionService struct {
	cfg                   *config.Config
	submissionRepository  repository.SubmissionRepository
	transactionRepository repository.TransactionRepository
	userRepository        repository.UserRepository
	productRepository     repository.ProductRepository
}

func NewSubmissionService(
	cfg *config.Config,
	submissionRepository repository.SubmissionRepository,
	transactionRepository repository.TransactionRepository,
	userRepository repository.UserRepository,
	productRepository repository.ProductRepository,
) SubmissionService {
	return &submissionService{cfg, submissionRepository, transactionRepository, userRepository, productRepository}
}

func (s *submissionService) GetAll(ctx context.Context, req dto.GetAllProductsRequest) ([]entity.Product, error) {
	return s.submissionRepository.GetAll(ctx, req)
}

func (s *submissionService) GetById(ctx context.Context, id int64) (*entity.Product, error) {
	return s.submissionRepository.GetById(ctx, id)
}

func (s *submissionService) Create(ctx context.Context, req dto.CreateProductRequest, t dto.CreateTransactionRequest) error {
	userID := req.UserID
	if userID == 0 {
		return errors.New("User ID tidak ditemukan")
	}
	user, err := s.userRepository.GetById(ctx, userID)
	if err != nil {
		return err
	}
	exist, err := s.productRepository.GetByName(ctx, req.ProductName)
	if err == nil && exist != nil {
		return errors.New("Nama event sudah digunakan")
	}
	if req.ProductPrice == 0 {
		req.ProductStatus = "pending"
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
		OrderID:            t.OrderID,
		UserID:             userID,
	}

	if err := s.submissionRepository.Create(ctx, submission); err != nil {
		return err
	}

	ProductID := submission.ID
	TransactionAmount := submission.ProductPrice * 0.25

	if TransactionAmount != 0 {
		transaction := &entity.Transaction{
			ProductID:           ProductID,
			UserID:              userID,
			TransactionQuantity: 1,
			TransactionAmount:   TransactionAmount,
			TransactionStatus:   "pending",
			OrderID:             t.OrderID,
			VerificationToken:   t.VerificationToken,
		}
		if err := s.transactionRepository.Create(ctx, transaction); err != nil {
			return err
		}
	}

	return nil

}

func (s *submissionService) UpdateByUser(ctx context.Context, req dto.UpdateProductRequest) error {
	userID := req.UserID
	if userID == 0 {
		return errors.New("User ID tidak ditemukan")
	}
	submission, err := s.submissionRepository.GetById(ctx, req.ID)
	if err != nil {
		return err
	}
	if submission.UserID != userID {
		return errors.New("Anda tidak memiliki hak untuk mengupdate pengajuan ini")
	}
	if submission.ProductStatus != "pending" {
		return errors.New("Pengajuan ini sudah tidak dapat diupdate karena sudah diterima atau ditolak oleh admin")
	}
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
	if req.ProductPrice != 0 {
		submission.ProductPrice = req.ProductPrice
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
	return s.submissionRepository.Update(ctx, submission)
}

func (s *submissionService) Approve(ctx context.Context, id int64) error {
	submission, err := s.submissionRepository.GetById(ctx, id)
	if err != nil {
		return err
	}
	if submission.ProductStatus != "pending" {
		return errors.New("Pengajuan ini sudah diterima atau ditolak oleh admin")
	}
	submission.ProductStatus = "accepted"
	return s.submissionRepository.Update(ctx, submission)
}

func (s *submissionService) Reject(ctx context.Context, id int64) error {
	submission, err := s.submissionRepository.GetById(ctx, id)
	if err != nil {
		return err
	}
	if submission.ProductStatus != "pending" {
		return errors.New("Pengajuan ini sudah diterima atau ditolak oleh admin")
	}
	submission.ProductStatus = "rejected"
	return s.submissionRepository.Update(ctx, submission)
}

func (s *submissionService) Cancel(ctx context.Context, submission *entity.Product, req dto.GetProductByIDRequest) error {
	userID := req.UserID
	if userID == 0 {
		return errors.New("User ID tidak ditemukan")
	}
	submission, err := s.submissionRepository.GetById(ctx, req.ID)
	if err != nil {
		return err
	}
	if submission.UserID != userID {
		return errors.New("Anda tidak memiliki hak untuk mengupdate pengajuan ini")
	}
	if submission.ProductStatus != "pending" {
		return errors.New("Pengajuan ini sudah tidak dapat dicancel karena sudah diterima atau ditolak oleh admin")
	}
	return s.submissionRepository.Delete(ctx, submission)
}

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
		if err := s.productRepository.Update(ctx, product); err != nil {
			return err
		}
		trans.TransactionStatus = "success"

		if err := s.transactionRepository.Update(ctx, trans); err != nil {
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
