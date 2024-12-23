package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"

	"github.com/mhusainh/FastTix/config"
	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
	"github.com/mhusainh/FastTix/utils"
	"gopkg.in/gomail.v2"
)

type SubmissionService interface {
	GetAll(ctx context.Context, req dto.GetAllProductsRequest) ([]entity.Product, error)
	GetById(ctx context.Context, id int64) (*entity.Product, error)
	GetByUserId(ctx context.Context, req dto.GetProductByUserIDRequest, user *entity.User) ([]entity.Product, error)
	Create(ctx context.Context, req dto.CreateProductRequest, t dto.CreateTransactionRequest, user *entity.User) (*entity.Product, error)
	UpdateByUser(ctx context.Context, req dto.UpdateProductRequest, user *entity.User, submission *entity.Product) (*entity.Product, error)
	Approval(ctx context.Context, req dto.UpdateProductStatusRequest, submission *entity.Product, user *entity.User) (*entity.Product, error)
	Cancel(ctx context.Context, submission *entity.Product, req dto.GetProductByIDRequest) error
	HandleMidtransNotification(ctx context.Context, notif map[string]interface{}) error
	UpdatePictureURL(ctx context.Context, req dto.GetProductByIDRequest, pictureURL string) error
}

type submissionService struct {
	cfg                   *config.Config
	submissionRepository  repository.SubmissionRepository
	transactionRepository repository.TransactionRepository
	productRepository     repository.ProductRepository
	userRepository        repository.UserRepository
}

func NewSubmissionService(
	cfg *config.Config,
	submissionRepository repository.SubmissionRepository,
	transactionRepository repository.TransactionRepository,
	productRepository repository.ProductRepository,
	userRepository repository.UserRepository,
) SubmissionService {
	return &submissionService{cfg, submissionRepository, transactionRepository, productRepository, userRepository}
}

func (s *submissionService) GetAll(ctx context.Context, req dto.GetAllProductsRequest) ([]entity.Product, error) {
	return s.submissionRepository.GetAll(ctx, req)
}

func (s *submissionService) GetById(ctx context.Context, id int64) (*entity.Product, error) {
	return s.submissionRepository.GetById(ctx, id)
}

func (s *submissionService) GetByUserId(ctx context.Context, req dto.GetProductByUserIDRequest, user *entity.User) ([]entity.Product, error) {
	req.UserID = user.ID
	return s.submissionRepository.GetByUserId(ctx, req)
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

	req.OrderID = fmt.Sprintf("daftar_id-%s", utils.RandomString(10))

	submission := &entity.Product{
		ProductName:        req.ProductName,
		ProductAddress:     req.ProductAddress,
		ProductTime:        req.ProductTime,
		ProductDate:        req.ProductDate,
		ProductPrice:       req.ProductPrice,
		ProductSold:        0,
		ProductDescription: req.ProductDescription,
		ProductCategory:    req.ProductCategory,
		ProductQuantity:    req.ProductQuantity,
		ProductType:        "available",
		ProductStatus:      req.ProductStatus,
		UserID:             user.ID,
		OrderID:            req.OrderID,
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
			VerificationToken:   t.VerificationToken,
			OrderID:             req.OrderID,
			CheckIn:             1,
		}
		if err := s.transactionRepository.Create(ctx, transaction); err != nil {
			return nil, err
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

func (s *submissionService) Approval(ctx context.Context, req dto.UpdateProductStatusRequest, submission *entity.Product, user *entity.User) (*entity.Product, error) {
	fmt.Println(req.Status)
	fmt.Println("tes")
	if req.Status == "approve" {
		submission.ProductStatus = "accepted"
	} else if req.Status == "reject" {
		submission.ProductStatus = "rejected"
	} else {
		return nil, errors.New("invalid status")
	}

	templatePath := "./templates/email/approval.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	var ReplacerEmail = struct {
		Status string
	}{
		Status: submission.ProductStatus,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, ReplacerEmail); err != nil {
		return nil, err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPConfig.Username)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Fast Tix : Approval !")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(
		s.cfg.SMTPConfig.Host,
		s.cfg.SMTPConfig.Port,
		s.cfg.SMTPConfig.Username,
		s.cfg.SMTPConfig.Password,
	)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	return submission, s.submissionRepository.Update(ctx, submission)
}

func (s *submissionService) Cancel(ctx context.Context, submission *entity.Product, req dto.GetProductByIDRequest) error {
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

func (s *submissionService) UpdatePictureURL(ctx context.Context, req dto.GetProductByIDRequest, pictureURL string) error {
	submisstion, err := s.productRepository.GetById(ctx, req.ID)
	if err != nil {
		return err
	}
	submisstion.ProductImage = pictureURL
	return s.submissionRepository.Update(ctx, submisstion)
}
