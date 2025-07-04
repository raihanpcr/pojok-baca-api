package book_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pojok-baca-api/dto"
	"pojok-baca-api/handler"
	"pojok-baca-api/model"
	"pojok-baca-api/service"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetBookByID_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/books/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockService := new(service.BookServiceMock)
	mockBook := model.Book{
		Name:       "Golang Book",
		Stok:       10,
		Category:   "Programming",
		RentalCost: 5000,
	}
	mockService.On("GetBookByID", uint(1)).Return(mockBook, nil)

	handler := handler.ProductHandler{Service: mockService}
	err := handler.GetBookByID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.BookByIDResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Success", resp.Status)
	assert.Equal(t, mockBook.Name, resp.Data.Name)

	mockService.AssertExpectations(t)
}

func TestGetBookByID_InvalidID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/books/abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:id")
	c.SetParamNames("id")
	c.SetParamValues("abc")

	mockService := new(service.BookServiceMock)
	handler := handler.ProductHandler{Service: mockService}

	err := handler.GetBookByID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dto.ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "BadRequest", resp.Status)
	assert.Equal(t, "Invalid ID", resp.Message)
}

func TestGetBookByID_NotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/books/99", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:id")
	c.SetParamNames("id")
	c.SetParamValues("99")

	mockService := new(service.BookServiceMock)
	mockService.On("GetBookByID", uint(99)).Return(model.Book{}, errors.New("not found"))

	handler := handler.ProductHandler{Service: mockService}
	err := handler.GetBookByID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var resp dto.ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "Not Found", resp.Status)
	assert.Equal(t, "Book not found", resp.Message)

	mockService.AssertExpectations(t)
}