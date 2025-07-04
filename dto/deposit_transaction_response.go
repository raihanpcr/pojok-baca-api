package dto

type DepositResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	OrderID string `json:"order_id"`
	PayURL  string `json:"payment_url"`
}
