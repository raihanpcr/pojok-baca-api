package handler_test

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

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetDataByID_Success(t *testing.T) {
	e := echo.New()

	// Buat token palsu
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(1),
		"email":   "john@mail.com",
		"role":    "user",
	})

	// Setup context Echo dan JWT Token
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", token)

	// Mock service dan data user
	mockService := new(service.UserServiceMock)
	deposit := 10000
	mockUser := model.User{
		Name:    "John Doe",
		Email:   "john@mail.com",
		Deposit: &deposit,
	}

	mockService.On("GetUserById", uint(1)).Return(mockUser, nil)

	// Panggil handler
	handler := handler.NewUserHandler(mockService)
	err := handler.GetDataByID(c)

	// Validasi
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.UserResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Successfully get your data", resp.Message)
	assert.Equal(t, "john@mail.com", resp.Data.Email)
	assert.Equal(t, 10000, resp.Data.Deposit)

	mockService.AssertExpectations(t)
}

func TestGetDataByID_UserNotFound(t *testing.T) {
	e := echo.New()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(99),
		"email":   "notfound@mail.com",
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", token)

	mockService := new(service.UserServiceMock)
	mockService.On("GetUserById", uint(99)).Return(model.User{}, errors.New("user not found"))

	handler := handler.NewUserHandler(mockService)
	err := handler.GetDataByID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dto.ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "Error", resp.Status)
	assert.Equal(t, "Failed to get your data", resp.Message)

	mockService.AssertExpectations(t)
}
