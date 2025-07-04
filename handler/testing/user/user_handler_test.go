package handler_test

import (
	"bytes"
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
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser_Success(t *testing.T) {
	e := echo.New()

	mockService := new(service.UserServiceMock)

	userHandler := handler.UserHandler{
		Service: mockService,
	}

	userInput := model.User{
		Name:     "Raihan",
		Email:    "raihan@mail.com",
		Password: "password123",
	}

	mockService.On("GetUserByEmail", userInput.Email).
		Return(model.User{}, errors.New("not found"))

	mockService.On("CreateUser", mock.MatchedBy(func(u model.User) bool {
		err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(userInput.Password))
		return u.Email == userInput.Email &&
			u.Name == userInput.Name &&
			u.Role == "user" &&
			err == nil
	})).Return(model.User{
		Name:  userInput.Name,
		Email: userInput.Email,
		Role:  "user",
	}, nil)

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := userHandler.CreateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response dto.RegisterSuccessResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "User created", response.Message)
	assert.Equal(t, userInput.Email, response.Data.Email)
	assert.Equal(t, userInput.Name, response.Data.Name)
	assert.Equal(t, "user", response.Data.Role)

	mockService.AssertExpectations(t)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	e := echo.New()
	mockService := new(service.UserServiceMock)

	userHandler := handler.UserHandler{
		Service: mockService,
	}

	//simulasi : user daftar
	userInput := model.User{
		Name:     "Raihan",
		Email:    "raihan@mail.com",
		Password: "password123",
	}

	//simulasi : Email yang sudah terdaftar
	mockService.On("GetUserByEmail", userInput.Email).
		Return(model.User{
			Name:  "Raihan",
			Email: "raihan@mail.com",
			Role:  "user",
		}, nil) // nil error = email ditemukan

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := userHandler.CreateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)

	var response dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Error", response.Status)
	assert.Equal(t, http.StatusConflict, response.Code)
	assert.Equal(t, "Email already exists", response.Message)

	mockService.AssertExpectations(t)
}
