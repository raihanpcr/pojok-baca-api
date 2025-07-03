package dto

type UserResponse struct {
	Status  string           `json:"status" example:"Success"`
	Code    int              `json:"code" example:"200"`
	Message string           `json:"message" example:"Success Get User"`
	Data    UserDataResponse `json:"data"`
}

type UserDataResponse struct {
	Name    string `json:"name" example:"John doe"`
	Email   string `json:"email" example:"johndoe@example.com"`
	Deposit int    `json:"deposit" example:"4"`
}
