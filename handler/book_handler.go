package handler

import (
	"net/http"
	"pojok-baca-api/dto"
	"pojok-baca-api/model"
	"pojok-baca-api/service"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	Service service.BookService
}

func NewProductHandler(s service.BookService) *ProductHandler {
	return &ProductHandler{Service: s}
}

// GetBooks godoc
// @Summary Get all books
// @Description Retrieve all available books with stock and rental info
// @Tags Books
// @Produce json
// @Success 200 {object} dto.BookResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /products [get]
func (h *ProductHandler) GetBooks(c echo.Context) error {
	products, err := h.Service.GetBooks()
	if err != nil {
		return c.JSON(http.StatusConflict, dto.ErrorResponse{
			Status:  "Duplicate",
			Code:    http.StatusConflict,
			Message: "Email already exists",
		})
	}

	var bookResponse []dto.GetAllBooksResponse

	for _, product := range products {
		bookResponse = append(bookResponse, dto.GetAllBooksResponse{
			ID:         product.ID,
			Name:       product.Name,
			Stok:       product.Stok,
			RentalCost: product.RentalCost,
			Category:   product.Category,
		})
	}

	return c.JSON(http.StatusOK, dto.BookResponse{
		Status:  "Success",
		Code:    http.StatusOK,
		Message: "Success Get All Books",
		Data:    bookResponse,
	})
}

// CreateBook godoc
// @Summary Create a new book
// @Description Only admin users can create books
// @Tags Books
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.GetBookData true "Book creation request"
// @Success 201 {object} dto.CreateBookResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products [post]
func (h *ProductHandler) CreateBook(c echo.Context) error {

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	role := claims["role"].(string)

	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Status:  "Unauthorized",
			Code:    http.StatusUnauthorized,
			Message: "You are not allowed to access this resource",
		})
	}

	var req dto.CreateBookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "BadRequest",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	newBook := model.Book{
		Name:       req.Name,
		Stok:       req.Stok,
		Category:   req.Category,
		RentalCost: req.RentalCost,
	}

	createdBook, err := h.Service.Create(newBook)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  "Error",
			Code:    http.StatusInternalServerError,
			Message: "Failed to create book",
		})
	}

	responseData := dto.GetBookData{
		ID:         createdBook.ID,
		Name:       createdBook.Name,
		Stok:       createdBook.Stok,
		Category:   createdBook.Category,
		RentalCost: createdBook.RentalCost,
	}

	return c.JSON(http.StatusCreated, dto.GetBookDataResponse{
		Status:  "Success",
		Code:    http.StatusCreated,
		Message: "Success Create Book",
		Data: responseData,
	})
}

// GetBookByID godoc
// @Summary Get a book by its ID
// @Description Get a single book using its ID
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} dto.BookSwaggerByIDResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetBookByID(c echo.Context) error {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "BadRequest",
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	book, err := h.Service.GetBookByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Status:  "Not Found",
			Code:    http.StatusNotFound,
			Message: "Book not found",
		})
	}

	return c.JSON(http.StatusOK, dto.BookByIDResponse{
		Status:  "Success",
		Code:    http.StatusOK,
		Message: "Success Get Book",
		Data:    book,
	})
}

// DeleteBookByID godoc
// @Summary Delete a book by its ID
// @Description Only admin can delete a book by ID
// @Tags Books
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 204 {object} nil "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteBookByID(c echo.Context) error {

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	role := claims["role"].(string)

	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Status:  "Unauthorized",
			Code:    http.StatusUnauthorized,
			Message: "You are not allowed to access this resource",
		})
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "BadRequest",
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	err = h.Service.DeleteBookByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  "Error",
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete book",
		})
	}

	return c.JSON(http.StatusNoContent, map[string]string{"message": "Success Delete Book"})
}

// UpdateBookByID godoc
// @Summary Update a book by its ID
// @Description Only admin can update a book's information
// @Tags Books
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param request body dto.UpdateBookRequest true "Update Book Request Body"
// @Success 200 {object} dto.UpdateBookByIDResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateBookByID(c echo.Context) error {
	// Authorization
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	role := claims["role"].(string)

	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Status:  "Unauthorized",
			Code:    http.StatusUnauthorized,
			Message: "You are not allowed to access this resource",
		})
	}

	// Get ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "BadRequest",
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	// Bind JSON to DTO
	var req dto.UpdateBookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "BadRequest",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Update via service
	book, err := h.Service.UpdateBookByID(req, uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  "Error",
			Code:    http.StatusInternalServerError,
			Message: "Failed to update book",
		})
	}

	// Return response
	return c.JSON(http.StatusOK, dto.UpdateBookByIDResponse{
		Status:  "Success",
		Code:    http.StatusOK,
		Message: "Success Update Book",
		Data: dto.UpdateResoponse{
			ID:         uint(book.ID),
			Name:       book.Name,
			Stok:       book.Stok,
			RentalCost: book.RentalCost,
			Category:   book.Category,
		},
	})
}
