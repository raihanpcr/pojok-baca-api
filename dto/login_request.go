package dto

type LoginRequest struct {
	Email    string `json:"email" example:"johndoe@gmail.com" validate:"required,email"`
	Password string `json:"password" example:"123456" validate:"required"`
}
