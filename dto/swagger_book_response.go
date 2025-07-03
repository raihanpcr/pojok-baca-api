package dto

type BookSwaggerByIDResponse struct {
	Status  string      `json:"status" example:"success"`
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"success"`
	Data    BookSwagger `json:"data"`
}
type BookSwagger struct {
	ID         uint   `json:"id" example:"1"`
	Name       string `json:"name" example:"Atomic Habits"`
	Stok       int    `json:"stok" example:"5"`
	Category   string `json:"category" example:"Self Development"`
	RentalCost int    `json:"rental_cost" example:"20000"`
}