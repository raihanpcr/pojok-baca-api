package handler

import (
	"net/http"
	"pojok-baca-api/dto"
	"pojok-baca-api/model"
	"pojok-baca-api/service"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type RentalHandler struct {
	Service     service.RentalService
	bookService service.BookService
	userService service.UserService
}

func NewRentalHandler(rentalService service.RentalService, bookService service.BookService, userService service.UserService) *RentalHandler {
	return &RentalHandler{
		Service:     rentalService,
		bookService: bookService,
		userService: userService,
	}
}

// CreateRental godoc
// @Summary Create a rental
// @Description Create a new book rental by user. Requires login.
// @Tags Rentals
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.RentalRequest true "Rental request payload"
// @Success 201 {object} dto.RentalResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /rentals [post]
func (h *RentalHandler) CreateRental(c echo.Context) error {

	//Get ID user login
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	email := claims["email"].(string)

	//Request data rental
	var req dto.RentalRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "Bad Request",
			Code:    http.StatusBadRequest,
			Message: "Invalid request format",
		})
	}

	//Request required
	if req.BookID == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "Bad Request",
			Code:    http.StatusBadRequest,
			Message: "book_id is required",
		})
	}

	//Check book is available or not
	book, err := h.bookService.GetBookByID(req.BookID)
	if book.Stok == 0 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "Internal Server Error",
			Code:    http.StatusBadRequest,
			Message: "Book not available",
		})
	}

	//Check deposit
	user, err := h.userService.GetUserByEmail(email)
	if user.Deposit == nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "Bad Request",
			Code:    http.StatusBadRequest,
			Message: "Deposit not found for this user",
		})
	}

	userDepo := user.Deposit
	bookCost := book.RentalCost
	if *userDepo < bookCost {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "Bad Request",
			Code:    http.StatusBadRequest,
			Message: "Insufficient deposit",
		})
	}

	//Update Deposit
	depoResult := *userDepo - bookCost
	_, err = h.userService.UpdateDepositUser(depoResult, user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  "Internal Server Error",
			Code:    http.StatusInternalServerError,
			Message: "Update deposit is failed",
		})
	}

	//Set rentDate & return Date
	rentDate := time.Now()
	returnDate := rentDate.AddDate(0, 0, 7)

	//Create New Data
	newRental := model.Rental{
		BookID:     req.BookID,
		UserID:     userID,
		RentDate:   rentDate,
		ReturnDate: &returnDate,
		Status:     "Borrowed",
	}

	//Create Rental
	createdRental, err := h.Service.CreateRental(newRental)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  "Internal Server Error",
			Code:    http.StatusInternalServerError,
			Message: "Create Rental is Failed",
		})
	}

	//Create rental success
	return c.JSON(http.StatusCreated, dto.RentalResponse{
		Status:  "Created",
		Code:    http.StatusCreated,
		Message: "Success Create Rental",
		Data: dto.RentalDataResponse{
			BookID:     createdRental.BookID,
			RentDate:   createdRental.RentDate.Format("2006-01-02"),
			ReturnDate: createdRental.ReturnDate.Format("2006-01-02"),
			Status:     createdRental.Status,
		},
	})
}

// GetRentalByUserID godoc
// @Summary Get rental history by logged-in user
// @Description Returns a list of rental data for the authenticated user
// @Tags Rentals
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.RentalUserResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /rentals [get]
func (h *RentalHandler) GetRentalByUserID(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	rentals, err := h.Service.GetRentalByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  "Internal Server Error",
			Code:    http.StatusInternalServerError,
			Message: "Get Data Rental User is Failed",
		})
	}

	var rentalResponses []dto.RentalUserDataResponse

	for _, rental := range rentals {
		returnDate := ""
		if rental.ReturnDate != nil {
			returnDate = rental.ReturnDate.Format("2006-01-02")
		}

		rentalResponses = append(rentalResponses, dto.RentalUserDataResponse{
			RentalID:   rental.ID,
			BookID:     rental.BookID,
			BookTitle:  rental.Book.Name,
			RentDate:   rental.RentDate.Format("2006-01-02"),
			ReturnDate: returnDate,
			Status:     rental.Status,
		})
	}

	return c.JSON(http.StatusOK, dto.RentalUserResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Success Get Rental By UserID",
		Data:    rentalResponses,
	})
}
