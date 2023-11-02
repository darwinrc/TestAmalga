package handler

import (
	mock_service "TestAmalga/parte1/service/mocks"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"TestAmalga/parte1/service"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandler_HandleGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockService(ctrl)

	handler := NewHandler(mockService)

	req, err := http.NewRequest("GET", "/resumen/2019-12-01?dias=5", nil)
	if err != nil {
		t.Fatalf("Error creando request: %v", err)
	}

	rr := httptest.NewRecorder()

	resumen := &service.Resumen{
		Total: 29815625.82,
		ComprasPorTDC: map[string]float64{
			"amex":           2861283.64,
			"amex corp":      2923605.31,
			"maestro":        2901466.24,
			"master classic": 2930035.72,
			"master gold":    2988878.78,
			"master plat":    2395721.04,
			"privada":        2538281.95,
			"visa classic":   2697760.23,
			"visa debit":     2318518.49,
			"visa gold":      2750625.22,
			"visa plat":      2509449.2,
		},
		Nocompraron:   1888,
		CompraMasAlta: 30477.77,
	}

	mockService.EXPECT().CalcularResumen("2019-12-01", "5").Return(resumen, nil).Times(1)

	r := mux.NewRouter()
	handler.Attach(r)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Se espera el status code %d, pero se obtuvo %d", http.StatusOK, rr.Code)
	}

	esperado, err := json.Marshal(resumen)
	if err != nil {
		t.Errorf("Error marshaleando resumen: %s", err)
	}

	assert.Equal(t, rr.Body.Bytes(), esperado)
}
