package dto

type LoginSuccessResponse struct {
	Token string `json:"token" example:"your-jwt-token"`
}