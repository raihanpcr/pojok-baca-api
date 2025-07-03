package dto

type RentalResponse struct {
	Status  string `json:"status" example:"success"`
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"success"`
	Data    RentalDataResponse
}

type RentalDataResponse struct {
	BookID     uint   `json:"book_id" example:"1" validate:"required"`
	RentDate   string `json:"rent_date" example:"2020-01-01"`
	ReturnDate string `json:"return_date" example:"2020-01-01"`
	Status     string `json:"status" example:"borrowed"`
}

type RentalUserResponse struct {
	Status  string                   `json:"status" example:"success"`
	Code    int                      `json:"code" example:"200"`
	Message string                   `json:"message" example:"Rental list"`
	Data    []RentalUserDataResponse `json:"data"`
}

type RentalUserDataResponse struct {
	RentalID   uint   `json:"rental_id" example:"1"`
	BookID     uint   `json:"book_id" example:"2"`
	BookTitle  string `json:"book_title" example:"Atomic Habits"`
	RentDate   string `json:"rent_date" example:"2025-07-03"`
	ReturnDate string `json:"return_date" example:"2025-07-10"`
	Status     string `json:"status" example:"Borrowed"`
}
