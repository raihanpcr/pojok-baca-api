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

func TestGetBooks_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock service
	mockService := new(service.BookServiceMock)

	mockBooks := []model.Book{
		{Name: "Book A", Stok: 10, RentalCost: 5000, Category: "Fiction"},
		{Name: "Book B", Stok: 5, RentalCost: 3000, Category: "History"},
	}

	mockService.On("GetBooks").Return(mockBooks, nil)

	handler := handler.ProductHandler{Service: mockService}

	err := handler.GetBooks(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.BookResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Success", resp.Status)
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, "Book A", resp.Data[0].Name)

	mockService.AssertExpectations(t)
}

func TestGetBooks_Failure(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/producs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(service.BookServiceMock)

	mockService.On("GetBooks").Return([]model.Book{}, errors.New("failed to fetch"))

	handler := handler.ProductHandler{Service: mockService}

	err := handler.GetBooks(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)

	var resp dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Duplicate", resp.Status)
	assert.Equal(t, "Email already exists", resp.Message)

	mockService.AssertExpectations(t)
}