package service

import (
	"bytes"
	"context"
	"errors"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mhusainh/FastTix/config"
	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
	"github.com/mhusainh/FastTix/utils"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type UserService interface {
	Login(ctx context.Context, username string, password string) (*entity.JWTCustomClaims, error)
	Register(ctx context.Context, req dto.UserRegisterRequest) error
	VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) error
	GetAll(ctx context.Context) ([]entity.User, error)
	GetById(ctx context.Context, id int64) (*entity.User, error)
	GetByIdByUser(ctx context.Context, id int64) (dto.GetUserByIDByUserRequest, error)
	Update(ctx context.Context, req dto.UpdateUserRequest) error
	Delete(ctx context.Context, user *entity.User) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
	RequestResetPassword(ctx context.Context, username string) error
}

type userService struct {
	tokenService   TokenService
	cfg            *config.Config
	userRepository repository.UserRepository
}

func NewUserService(
	tokenService TokenService,
	cfg *config.Config,
	userRepository repository.UserRepository,
) UserService {
	return &userService{tokenService, cfg, userRepository}
}

func (s *userService) Login(ctx context.Context, email string, password string) (*entity.JWTCustomClaims, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("Email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("Email atau password salah")
	}

	if user.IsVerified == 0 {
		return nil, errors.New("Silahkan verifikasi email terlebih dahulu")
	}

	expiredTime := time.Now().Add(time.Minute * 10)

	claims := &entity.JWTCustomClaims{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Fast Tix",
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}

	return claims, nil
}

func (s *userService) Register(ctx context.Context, req dto.UserRegisterRequest) error {
	user := new(entity.User)
	user.Email = req.Email
	user.FullName = req.FullName
	user.Gender = req.Gender
	user.Role = "User"
	user.VerifyEmailToken = utils.RandomString(16)
	user.IsVerified = 0

	exist, err := s.userRepository.GetByEmail(ctx, req.Email)
	if err == nil && exist != nil {
		return errors.New("Email sudah digunakan")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	templatePath := "./templates/email/verify-email.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var ReplacerEmail = struct {
		Token string
	}{
		Token: user.VerifyEmailToken,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, ReplacerEmail); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPConfig.Username)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Fast Tix : Verifikasi Email!")
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

	err = s.userRepository.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) GetAll(ctx context.Context) ([]entity.User, error) {
	return s.userRepository.GetAll(ctx)
}

func (s *userService) GetById(ctx context.Context, id int64) (*entity.User, error) {
	return s.userRepository.GetById(ctx, id)
}

func (s *userService) GetByIdByUser(ctx context.Context, id int64) (dto.GetUserByIDByUserRequest, error) {
	user, err := s.userRepository.GetById(ctx, id)
	if err != nil {
		return dto.GetUserByIDByUserRequest{}, err
	}
	return dto.GetUserByIDByUserRequest{
		ID:       user.ID,
		FullName: user.FullName,
		Gender:   user.Gender,
		Email:    user.Email,
	}, nil
}
func (s *userService) Update(ctx context.Context, req dto.UpdateUserRequest) error {
	user, err := s.userRepository.GetById(ctx, req.ID)
	if err != nil {
		return err
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	return s.userRepository.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, user *entity.User) error {
	return s.userRepository.Delete(ctx, user)
}

func (s *userService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	user, err := s.userRepository.GetByResetPasswordToken(ctx, req.Token)
	if err != nil {
		return errors.New("Token reset password salah")
	}

	if req.Password == "" {
		return errors.New("Password tidak boleh kosong")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.userRepository.Update(ctx, user)
}

func (s *userService) RequestResetPassword(ctx context.Context, email string) error {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return errors.New("Email tersebut tidak ditemukan")
	}

	expiredTime := time.Now().Add(10 * time.Minute)

	claims := &entity.ResetPasswordClaims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
			Issuer:    "Reset Password",
		},
	}

	token, err := s.tokenService.GenerateResetPasswordToken(ctx, *claims)
	if err != nil {
		return err
	}

	user.ResetPasswordToken = token
	err = s.userRepository.Update(ctx, user)
	if err != nil {
		return err
	}

	templatePath := "./templates/email/reset-password.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var ReplacerEmail = struct {
		Token string
	}{
		Token: user.ResetPasswordToken,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, ReplacerEmail); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPConfig.Username)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Reset Password Request !")
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

func (s *userService) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) error {
	user, err := s.userRepository.GetByVerifyEmailToken(ctx, req.Token)
	if err != nil {
		return errors.New("Token verifikasi email salah")
	}
	user.IsVerified = 1
	return s.userRepository.Update(ctx, user)
}
