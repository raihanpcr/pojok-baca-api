package rental

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pojok-baca-api/dto"
	"pojok-baca-api/handler"
	"pojok-baca-api/model"
	"pojok-baca-api/service"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateRental_Success(t *testing.T) {
	e := echo.New()
	body := `{"book_id": 1}`

	req := httptest.NewRequest(http.MethodPost, "/rentals", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock JWT claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(1),
		"email":   "user@mail.com",
	})
	c.Set("user", token)

	// Init mock services
	mockRentalService := new(service.RentalServiceMock)
	mockBookService := new(service.BookServiceMock)
	mockUserService := new(service.UserServiceMock)

	// Mock data
	book := model.Book{
		Model:      gorm.Model{ID: 1},
		Stok:       2,
		RentalCost: 5000,
	}
	deposit := 10000
	user := model.User{
		Model:  gorm.Model{ID: 1},
		Email:  "user@mail.com",
		Deposit: &deposit,
	}
	createdRental := model.Rental{
		BookID:     1,
		UserID:     1,
		RentDate:   time.Now(),
		ReturnDate: ptrToTime(time.Now().AddDate(0, 0, 7)),
		Status:     "Borrowed",
	}

	// Set expectations
	mockBookService.On("GetBookByID", uint(1)).Return(book, nil)
	mockUserService.On("GetUserByEmail", "user@mail.com").Return(user, nil)
	mockUserService.On("UpdateDepositUser", 5000, uint(1)).Return(user, nil)
	mockRentalService.On("CreateRental", mock.AnythingOfType("model.Rental")).Return(createdRental, nil)

	// Call handler
	handler := handler.NewRentalHandler(mockRentalService, mockBookService, mockUserService)
	err := handler.CreateRental(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp dto.RentalResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Success Create Rental", resp.Message)
	assert.Equal(t, uint(1), resp.Data.BookID)

	// Verify expectations
	mockBookService.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
	mockRentalService.AssertExpectations(t)
}

func ptrToTime(t time.Time) *time.Time {
	return &t
}


