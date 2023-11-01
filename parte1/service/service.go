package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
)

type Service interface {
	CalcularResumen(fecha, dias string) (*Resumen, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

type Resumen struct {
	Total         float64            `json:"total"`
	ComprasPorTDC map[string]float64 `json:"comprasPorTDC"`
	Nocompraron   int                `json:"nocompraron"`
	CompraMasAlta float64            `json:"compraMasAlta"`
}

func (r *Resumen) Normalizar() {
	r.Total = math.Round(r.Total*100) / 100
	r.CompraMasAlta = math.Round(r.CompraMasAlta*100) / 100

	for k, v := range r.ComprasPorTDC {
		r.ComprasPorTDC[k] = math.Round(v*100) / 100
	}
}

type Compra struct {
	ClientID int     `json:"clientId"`
	Phone    string  `json:"phone"`
	Nombre   string  `json:"nombre"`
	Compro   bool    `json:"compro"`
	TDC      string  `json:"tdc"`
	Monto    float64 `json:"monto"`
	Date     string  `json:"date"`
}

const urlStr = "https://apirecruit-gjvkhl2c6a-uc.a.run.app/compras/"

func (s *service) CalcularResumen(fecha, dias string) (*Resumen, error) {
	date, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		fmt.Println("Error parseando fecha:", err)
		return nil, errors.New(fmt.Sprintln("error parseando fecha: ", err))
	}

	numDias, err := strconv.Atoi(dias)
	if err != nil {
		return nil, errors.New(fmt.Sprintln("error convirtiendo dias a int: ", err))
	}

	comprasTotales := []Compra{}

	for i := 0; i < numDias; i++ {
		dia := date.Add(time.Duration(i) * time.Hour * 24).Format("2006-01-02")

		fmt.Println("llamando: ", urlStr+dia)
		resp, err := http.Get(urlStr + dia)
		if err != nil {
			return nil, errors.New(fmt.Sprintln("error obteniendo url: ", err))
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New(fmt.Sprintln("error leyendo body: ", err))
		}

		comprasDia := []Compra{}
		err = json.Unmarshal(body, &comprasDia)
		if err != nil {
			return nil, errors.New(fmt.Sprintln("error unmarshaleando body: ", err))
		}

		comprasTotales = append(comprasTotales, comprasDia...)
	}

	resumen := &Resumen{
		ComprasPorTDC: map[string]float64{},
	}

	for _, compra := range comprasTotales {
		if compra.Compro {
			resumen.Total += compra.Monto
			resumen.ComprasPorTDC[compra.TDC] += compra.Monto

			if compra.Monto > resumen.CompraMasAlta {
				resumen.CompraMasAlta = compra.Monto
			}
		} else {
			resumen.Nocompraron++
		}
	}

	return resumen, nil
}
