package dto

import "pojok-baca-api/model"

type BookResponse struct {
	Status  string                `json:"status" example:"success"`
	Code    int                   `json:"code" example:"200"`
	Message string                `json:"message" example:"success"`
	Data    []GetAllBooksResponse `json:"data"`
}

type GetAllBooksResponse struct {
	ID         uint   `json:"id" example:"0"`
	Name       string `json:"name" example:"johndoe"`
	Stok       int    `json:"stok" example:"1"`
	RentalCost int    `json:"rental_cost" example:"1"`
	Category   string `json:"category" example:"programming"`
}

type CreateBookRequest struct {
	Name       string `json:"name" example:"johndoe"`
	Stok       int    `json:"stok" example:"1"`
	RentalCost int    `json:"rental_cost" example:"1"`
	Category   string `json:"category" example:"programming"`
}

type BookByIDResponse struct {
	Status  string     `json:"status" example:"success"`
	Code    int        `json:"code" example:"200"`
	Message string     `json:"message" example:"success"`
	Data    model.Book `json:"data"`
}

type UpdateBookByIDResponse struct {
	Status  string          `json:"status" example:"success"`
	Code    int             `json:"code" example:"200"`
	Message string          `json:"message" example:"success"`
	Data    UpdateResoponse `json:"data"`
}

type UpdateResoponse struct {
	ID         uint   `json:"id" example:"2"`
	Name       string `json:"name" example:"johndoe"`
	Stok       int    `json:"stok" example:"1"`
	RentalCost int    `json:"rental_cost" example:"1"`
	Category   string `json:"category" example:"programming"`
}
