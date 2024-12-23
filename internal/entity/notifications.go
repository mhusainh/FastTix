package entity

import "time"

type Notification struct {
	ID        int64     `json:"id"`
	Message   string    `json:"message"`
	IsRead    int       `json:"is_read"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Notification) TableName() string {
	return "notifications"
}
