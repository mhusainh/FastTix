package entity

import "time"

type User struct {
	ID        int64  `json:"id"`
	FullName  string `json:"full_name"`
	Gender    string `json:"gender"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}