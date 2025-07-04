package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"pojok-baca-api/handler"
	"pojok-baca-api/model"
	"pojok-baca-api/service"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginSuccess(t *testing.T) {
	// Setup Echo dan Recorder
	e := echo.New()
	reqBody := `{"email": "john@mail.com", "password": "secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Setup Mock Service
	mockService := new(service.UserServiceMock)

	// Password yang sudah di-hash (hasil hash dari "secret123")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	mockUser := model.User{
		Email:    "john@mail.com",
		Password: string(hashedPassword),
		Role:     "user",
	}

	// Ekspektasi mock
	mockService.On("GetUserByEmail", "john@mail.com").Return(mockUser, nil)

	// Set JWT_SECRET untuk test
	os.Setenv("JWT_SECRET", "mysecret")

	// Jalankan handler
	handler := handler.UserHandler{Service: mockService}
	err := handler.Login(c)

	// Validasi hasil
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["token"])

	// Pastikan ekspektasi mock terpenuhi
	mockService.AssertExpectations(t)
}

func TestLogin_EmailNotFound(t *testing.T) {
	e := echo.New()
	reqBody := `{"email": "notfound@mail.com", "password": "secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(service.UserServiceMock)
	mockService.On("GetUserByEmail", "notfound@mail.com").Return(model.User{}, errors.New("user not found"))

	handler := handler.UserHandler{Service: mockService}
	err := handler.Login(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "Login failed", resp["status"])
	mockService.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	e := echo.New()
	reqBody := `{"email": "john@mail.com", "password": "wrongpass"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(service.UserServiceMock)

	// Hash password asli
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
	user := model.User{
		Email:    "john@mail.com",
		Password: string(hashedPassword),
		Role:     "user",
	}

	mockService.On("GetUserByEmail", "john@mail.com").Return(user, nil)

	handler := handler.UserHandler{Service: mockService}
	err := handler.Login(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "Login failed", resp["status"])
	mockService.AssertExpectations(t)
}
