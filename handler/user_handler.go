package handler

import (
	"net/http"
	"os"
	"pojok-baca-api/dto"
	"pojok-baca-api/model"
	"pojok-baca-api/service"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

// CreateUser godoc
// @Summary Register a new user
// @Description Create a new user account (name, email, password required)
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration request"
// @Success 201 {object} dto.RegisterSuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/register [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	var u model.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "Error",
			Code:    http.StatusBadRequest,
			Message: "Invalid JSON",
			Details: err.Error(),
		})
	}

	//name, email, password not fill
	if u.Name == "" || u.Email == "" || u.Password == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "Error",
			Code:    http.StatusBadRequest,
			Message: "name, email, password required",
		})
	}

	//role not fill
	if u.Role == "" {
		u.Role = "user"
	}

	//Duplicate Email
	_, err := h.Service.GetUserByEmail(u.Email)
	if err == nil {
		return c.JSON(http.StatusConflict, dto.ErrorResponse{
			Status:  "Error",
			Code:    http.StatusConflict,
			Message: "Email already exists",
		})
	}

	//hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password : "+err.Error())
	}

	u.Password = string(hashed)
	user, err := h.Service.CreateUser(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  "Error",
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.RegisterSuccessResponse{
		Status:  "success",
		Code:    http.StatusCreated,
		Message: "User created",
		Data: dto.RegisterResponse{
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
	})
}

// Login godoc
// @Summary User login
// @Description Authenticate user and get JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.LoginSuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "Error",
			Code:    http.StatusBadRequest,
			Message: "Invalid JSON",
		})
	}

	user, err := h.Service.GetUserByEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Status:  "Login failed",
			Code:    http.StatusUnauthorized,
			Message: "Email or password is incorrect",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Status:  "Login failed",
			Code:    http.StatusUnauthorized,
			Message: "Email or password is incorrect",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  "Error",
			Code:    http.StatusInternalServerError,
			Message: "Failed to sign token",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": tokenString,
	})
}

// GetDataByID godoc
// @Summary Get current logged-in user data
// @Description Retrieve user profile based on JWT token
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /user/profile [get]
func (h *UserHandler) GetDataByID(c echo.Context) error {

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	user, err := h.Service.GetUserById(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "Error",
			Code:    http.StatusBadRequest,
			Message: "Failed to get your data",
		})
	}

	deposit := 0
	if user.Deposit != nil {
		deposit = *user.Deposit
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Successfully get your data",
		Data: dto.UserDataResponse{
			Name:    user.Name,
			Email:   user.Email,
			Deposit: deposit,
		},
	})
}
