package rental

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pojok-baca-api/dto"
	"pojok-baca-api/handler"
	"pojok-baca-api/model"
	"pojok-baca-api/service"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetRentalByUserID_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/rentals/me", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulasi JWT claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(1),
	})
	c.Set("user", token)

	// Mock data
	returnDate := time.Now().AddDate(0, 0, 7)
	mockRentals := []model.Rental{
		{
			Model:      gorm.Model{ID: 1},
			BookID:     1,
			RentDate:   time.Now(),
			ReturnDate: &returnDate,
			Status:     "Borrowed",
			Book:       model.Book{Name: "Clean Architecture"},
		},
	}

	// Mock service
	mockRentalService := new(service.RentalServiceMock)
	mockRentalService.On("GetRentalByUserID", uint(1)).Return(mockRentals, nil)

	// Handler
	handler := handler.NewRentalHandler(mockRentalService, nil, nil)

	// Eksekusi
	err := handler.GetRentalByUserID(c)

	// Validasi
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.RentalUserResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Success Get Rental By UserID", resp.Message)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, uint(1), resp.Data[0].RentalID)
	assert.Equal(t, "Clean Architecture", resp.Data[0].BookTitle)

	mockRentalService.AssertExpectations(t)
}
