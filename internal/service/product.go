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
	"github.com/mhusainh/FastTix/utils"
	"gopkg.in/gomail.v2"
)

type ProductService interface {
	GetAllPending(ctx context.Context) ([]entity.Product, error)
	GetByIdPending(ctx context.Context, id int64) (*entity.Product, error)
	GetById(ctx context.Context, id int64) (*entity.Product, error)
	Create(ctx context.Context, req dto.CreateProductRequest) error
	Update(ctx context.Context, req dto.UpdateProductRequest) error
	Delete(ctx context.Context, product *entity.Product) error
	VerifySubmission(ctx context.Context, req dto.VerifySubmissionRequest) error
}

type productService struct {
	cfg               *config.Config
	productRepository repository.ProductRepository
}

func NewProductService(
	cfg *config.Config,
	productRepository repository.ProductRepository,
) ProductService {
	return &productService{cfg, productRepository}
}

func (s productService) GetAllPending(ctx context.Context) ([]entity.Product, error) {
	return s.productRepository.GetAllPending(ctx)
}

func (s productService) GetByIdPending(ctx context.Context, id int64) (*entity.Product, error) {
	return s.productRepository.GetByIdPending(ctx, id)
}

func (s productService) GetById(ctx context.Context, id int64) (*entity.Product, error) {
	return s.productRepository.GetById(ctx, id)
}

func (s productService) Create(ctx context.Context, req dto.CreateProductRequest) error {
	product := &entity.Product{
		ProductName:           req.ProductName,
		ProductAddress:        req.ProductAddress,
		ProductTime:           req.ProductTime,
		ProductDate:           req.ProductDate,
		ProductPrice:          req.ProductPrice,
		ProductDescription:    req.ProductDescription,
		ProductCategory:       req.ProductCategory,
		ProductStatus:         "pending",
		VerifySubmissionToken: utils.RandomString(16),
		UserID:                req.UserID,
	}
	err := s.productRepository.Create(ctx, product)
	if err != nil {
		return err
	}

	templatePath := "./templates/email/verify-submission.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var ReplacerEmail = struct {
		Token string
	}{
		Token: product.VerifySubmissionToken,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, ReplacerEmail); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPConfig.Username)
	m.SetHeader("To", "jovande19@gmail.com")
	m.SetHeader("Subject", "Fast Tix : Verify your Submission !")
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

	return nil
}

func (s productService) Update(ctx context.Context, req dto.UpdateProductRequest) error {
	product, err := s.productRepository.GetById(ctx, req.ID)
	if err != nil {
		return err
	}
	if req.ProductName != "" {
		product.ProductName = req.ProductName
	}
	if req.ProductAddress != "" {
		product.ProductAddress = req.ProductAddress
	}
	if req.ProductTime != "" {
		product.ProductTime = req.ProductTime
	}
	if req.ProductDate != "" {
		product.ProductDate = req.ProductDate
	}
	if req.ProductPrice != 0 {
		product.ProductPrice = req.ProductPrice
	}
	if req.ProductDescription != "" {
		product.ProductDescription = req.ProductDescription
	}
	if req.ProductCategory != "" {
		product.ProductCategory = req.ProductCategory
	}
	return s.productRepository.Update(ctx, product)
}

func (s productService) Delete(ctx context.Context, product *entity.Product) error {
	return s.productRepository.Delete(ctx, product)
}

func (s productService) VerifySubmission(ctx context.Context, req dto.VerifySubmissionRequest) error {
	product, err := s.productRepository.GetByVerifySubmissionToken(ctx, req.Token)
	if err != nil {
		return errors.New("Token Verification is Wrong")
	}
	product.IsVerified = 1
	return s.productRepository.Update(ctx, product)
}
