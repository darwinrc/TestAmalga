package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestCalcularResumen(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCases := []struct {
		fecha, dias   string
		resumen       *Resumen
		errorEsperado bool
		mock          string
	}{
		{
			fecha: "2019-12-01",
			dias:  "5",
			resumen: &Resumen{
				Total: 5.93146013e+06,
				ComprasPorTDC: map[string]float64{
					"amex":           463818.01,
					"amex corp":      493481.79,
					"maestro":        677530.63,
					"master classic": 537349.74,
					"master gold":    704109.56,
					"master plat":    503626.74,
					"privada":        543205.76,
					"visa classic":   537807.08,
					"visa debit":     418471.12,
					"visa gold":      574832.87,
					"visa plat":      477226.83,
				},
				Nocompraron:   407,
				CompraMasAlta: 30408.89,
			},
			errorEsperado: false,
			mock:          leerMock(),
		},
	}

	for _, tc := range testCases {
		httpmock.RegisterResponder("GET", "https://apirecruit-gjvkhl2c6a-uc.a.run.app/compras/"+tc.fecha,
			httpmock.NewStringResponder(200, tc.mock))

		s := NewService()
		res, err := s.CalcularResumen(tc.fecha, tc.dias)

		if tc.errorEsperado {
			if err == nil {
				t.Errorf("se espera un error, pero no se obtuvo ninguno")
			}
			return
		}

		if err != nil {
			t.Errorf("error inesperado: %v", err)
			return
		}

		if res == nil {
			t.Errorf("se espera una respuesta no nula")
			return
		}

		assert.Equal(t, tc.resumen, res)
	}
}

func leerMock() string {
	filePath := "./mocks/mock_api.json"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error abriendo el archivo de mock:", err)
		return ""
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error leyendo el body:", err)
		return ""
	}

	return string(data)
}
