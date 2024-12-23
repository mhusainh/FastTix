package entity

import "time"

type User struct {
	ID                 int64     `json:"id"`
	FullName           string    `json:"full_name"`
	Gender             string    `json:"gender"`
	Email              string    `json:"email"`
	Password           string    `json:"password"`
	Role               string    `json:"role"`
	ResetPasswordToken string    `json:"reset_password_token"`
	VerifyEmailToken   string    `json:"verify_email_token"`
	IsVerified         int       `json:"is_verified"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
