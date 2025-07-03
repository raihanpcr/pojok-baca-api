package dto

type RegisterSuccessResponse struct {
	Status  string           `json:"status" example:"success"`
	Code    int              `json:"code" example:"200"`
	Message string           `json:"message" example:"success"`
	Data    RegisterResponse `json:"data"`
}

type RegisterResponse struct {
	Email string `json:"email" example:"johndoe@gmail.com"`
	Name  string `json:"name" example:"John Doe"`
	Role  string `json:"role" example:"user"`
}
