package handler

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"pojok-baca-api/dto"
	"pojok-baca-api/service"
)

type DepositTransactionHandler struct {
	Service service.DepositTransactionService
}

func NewDepositTransactionHandler(service service.DepositTransactionService) *DepositTransactionHandler {
	return &DepositTransactionHandler{Service: service}
}

func (h *DepositTransactionHandler) Create(c echo.Context) error {
	var req dto.CreateDepositoryRequest
	if err := c.Bind(&req); err != nil || req.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "error",
			Code:    400,
			Message: "Invalid amount",
		})
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	res, err := h.Service.CreateTransaction(userID, req.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  "error",
			Code:    500,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.DepositResponse{
		Status:  "success",
		Code:    200,
		Message: "Success Create Transaction",
		OrderID: res.Token, // or res.OrderID
		PayURL:  res.RedirectURL,
	})
}

func (h *DepositTransactionHandler) Webhook(c echo.Context) error {
	var payload struct {
		OrderID           string `json:"order_id"`
		TransactionStatus string `json:"transaction_status"`
	}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid payload")
	}

	err := h.Service.HandleWebhook(payload.OrderID, payload.TransactionStatus)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to process webhook")
	}

	return c.JSON(http.StatusOK, "OK")
}
