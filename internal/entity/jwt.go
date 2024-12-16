package entity

import "github.com/golang-jwt/jwt/v5"

type JWTCustomClaims struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type ResetPasswordClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}