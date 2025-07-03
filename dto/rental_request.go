package dto

type RentalRequest struct {
	BookID uint `json:"book_id" example:"1" validate:"required"`
}
