package book_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pojok-baca-api/dto"
	"pojok-baca-api/handler"
	"pojok-baca-api/service"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDeleteBookByID_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/books/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Simulasi token role admin
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "admin",
	})
	c.Set("user", token)

	mockService := new(service.BookServiceMock)
	mockService.On("DeleteBookByID", uint(1)).Return(nil)

	handler := handler.ProductHandler{Service: mockService}
	err := handler.DeleteBookByID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	mockService.AssertExpectations(t)
}

func TestDeleteBookByID_Unauthorized(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/books/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Simulasi token role bukan admin
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "user",
	})
	c.Set("user", token)

	mockService := new(service.BookServiceMock)
	handler := handler.ProductHandler{Service: mockService}

	err := handler.DeleteBookByID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp dto.ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "Unauthorized", resp.Status)
	assert.Equal(t, "You are not allowed to access this resource", resp.Message)
}
