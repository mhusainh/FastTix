package entity

import "github.com/golang-jwt/jwt/v5"

type JWTCustomClaims struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Gender   string `json:"gender"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
