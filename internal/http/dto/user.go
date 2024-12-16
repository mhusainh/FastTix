package dto

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserRegisterRequest struct {
	FullName         string `json:"full_name" validate:"required"`
	Gender           string `json:"gender" validate:"required"`
	Email            string `json:"email" validate:"required"`
	Password         string `json:"password" validate:"required"`
	Role             string `json:"role" validate:"required"`
	VerifyEmailToken string `json:"verify_email_token" validate:"required"`
}

type UpdateUserRequest struct {
	ID       int64  `param:"id" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Gender   string `json:"gender" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type DeleteUserRequest struct {
	ID int64 `param:"id" validate:"required"`
}

type GetUserByIDRequest struct {
	ID int64 `param:"id" validate:"required"`
}

type ResetPasswordRequest struct {
	Token              string `param:"token" validate:"required"`
	Password           string `json:"password" validate:"required"`
}

type RequestResetPassword struct {
	Email string `json:"email" validate:"required"`
}

type VerifyEmailRequest struct {
	Token string `param:"token" validate:"required"`
}
