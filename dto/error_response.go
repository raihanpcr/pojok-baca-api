package dto

type ErrorResponse struct {
	Status  string      `json:"status" example:"error"`
	Code    int         `json:"code" example:"400"`
	Message string      `json:"message" example:"something went wrong"`
	Details interface{} `json:"details,omitempty"`
}
