package service

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mhusainh/FastTix/internal/entity"
)

type TokenService interface {
	GenerateAccessToken(ctx context.Context, claims entity.JWTCustomClaims) (string, error)
	GenerateResetPasswordToken(ctx context.Context, claims entity.ResetPasswordClaims) (string, error)
}

type tokenService struct {
	secretKey string
}

func NewTokenService(secretKey string) TokenService {
	return &tokenService{secretKey}
}

func (s *tokenService) GenerateAccessToken(ctx context.Context, claims entity.JWTCustomClaims) (string, error) {
	plainToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	encodedToken, err := plainToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}
	return encodedToken, nil
}

func (s *tokenService) GenerateResetPasswordToken(ctx context.Context, claims entity.ResetPasswordClaims) (string, error) {
	plainToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	encodedToken, err := plainToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}
	return encodedToken, nil
}
