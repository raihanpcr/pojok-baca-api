package dto

type RegisterRequest struct {
	Name     string `json:"name" example:"John Doe" validate:"required"`
	Email    string `json:"email" example:"johndoe@gmail.com" validate:"required,email"`
	Password string `json:"password" example:"123456" validate:"required"`
}
