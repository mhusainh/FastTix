package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mhusainh/FastTix/internal/entity"
)

type TokenService interface {
	GenerateAccessToken(ctx context.Context, claims entity.JWTCustomClaims) (string, error)
	GenerateResetPasswordToken(ctx context.Context, claims entity.ResetPasswordClaims) (string, error)
	ValidateToken(ctx context.Context, tokenString string) (jwt.MapClaims, error)
	GetUserIDFromToken(ctx echo.Context) (int64, error)
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

func (s *tokenService) ValidateToken(ctx context.Context, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token claims")
	}
}

func (s *tokenService) GetUserIDFromToken(ctx echo.Context) (int64, error) {
	tokenString := ctx.Request().Header.Get("Authorization")
	if tokenString == "" {
		return 0, fmt.Errorf("Authorization token is required")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims, err := s.ValidateToken(ctx.Request().Context(), tokenString)
	if err != nil {
		return 0, fmt.Errorf("invalid token")
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user ID not found in token")
	}
	
	return int64(userID), nil
}
