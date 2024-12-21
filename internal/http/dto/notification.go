package dto

type GetNotificationByIDRequest struct {
	ID int64 `param:"id" validate:"required"`
}

type CreateNotificationRequest struct {
	Message string `json:"message" validate:"required"`
	IsRead  int    `json:"is_read" validate:"required"`
	UserID  int64  `json:"user_id" validate:"required"`
}
