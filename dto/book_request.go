package dto

type UpdateBookRequest struct {
	Name       string `json:"name" validate:"required"`
	Stok       int    `json:"stok" validate:"required,gte=0"`
	RentalCost int    `json:"rental_cost" validate:"required,gte=0"`
	Category   string `json:"category" validate:"required"`
}

type CreateBookResponse struct {
	Status  string `json:"status" example:"success"`
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"success create book"`
}

type GetBookData struct {
	ID         uint   `json:"id" example:"1"`
	Name       string `json:"name" example:"Atomic Habits"`
	Stok       int    `json:"stok" example:"5"`
	Category   string `json:"category" example:"Self Development"`
	RentalCost int    `json:"rental_cost" example:"20000"`
}

type GetBookDataResponse struct {
	Status  string `json:"status" example:"success"`
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"success create book"`
	Data GetBookData `json:"data"`
}