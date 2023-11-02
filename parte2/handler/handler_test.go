package handler

import (
	"TestAmalga/parte2/service"
	mock_service "TestAmalga/parte2/service/mocks"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandler_HandleGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockService(ctrl)

	handler := NewHandler(mockService)

	req, err := http.NewRequest("GET", "/organization?file=organization.csv", nil)
	if err != nil {
		t.Fatalf("Error creando request: %v", err)
	}

	rr := httptest.NewRecorder()

	orgs := []service.Organization{
		{
			Organization: "org1",
			Users: []service.User{
				{
					Username: "jperez",
					Roles:    []string{"admin", "superadmin"},
				},
				{
					Username: "asosa",
					Roles:    []string{"writer"},
				},
			},
		},
		{
			Organization: "org2",
			Users: []service.User{
				{
					Username: "jperez",
					Roles:    []string{"admin"},
				},
				{
					Username: "rrodriguez",
					Roles:    []string{"writer", "editor"},
				},
			},
		},
	}
	mockService.EXPECT().GetOrgFromFile("organization.csv").Return(orgs, nil).Times(1)

	r := mux.NewRouter()
	handler.Attach(r)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Se espera el status code %d, pero se obtuvo %d", http.StatusOK, rr.Code)
	}

	esperado, err := json.Marshal(orgs)
	if err != nil {
		t.Errorf("Error marshaleando resumen: %s", err)
	}

	assert.Equal(t, rr.Body.Bytes(), esperado)
}
