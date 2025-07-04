package book_test

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

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateBook_Success(t *testing.T) {
	e := echo.New()
	body := `{"name":"Golang","stok":10,"category":"Programming","rental_cost":5000}`
	req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulasi JWT dengan role admin
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "admin",
	})
	c.Set("user", token)

	mockService := new(service.BookServiceMock)

	mockBook := model.Book{
		Name:       "Golang",
		Stok:       10,
		Category:   "Programming",
		RentalCost: 5000,
	}
	mockService.On("Create", mock.AnythingOfType("model.Book")).Return(mockBook, nil)

	handler := handler.ProductHandler{Service: mockService}
	err := handler.CreateBook(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp dto.GetBookDataResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Success", resp.Status)
	assert.Equal(t, "Success Create Book", resp.Message)
	assert.Equal(t, "Golang", resp.Data.Name)

	mockService.AssertExpectations(t)
}

func TestCreateBook_Unauthorized(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/books", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// role bukan admin
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "user",
	})
	c.Set("user", token)

	mockService := new(service.BookServiceMock)
	handler := handler.ProductHandler{Service: mockService}
	err := handler.CreateBook(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp dto.ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "Unauthorized", resp.Status)
	assert.Equal(t, "You are not allowed to access this resource", resp.Message)
}