package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"sync"
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

// Normalizar redondea los valores del resumen a dos decimales
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

// CalcularResumen calcula el resumen de las compras a partir de la fecha y los dias indicados
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

	compras, err := s.obtenerCompras(numDias, date)
	if err != nil {
		return nil, errors.New(fmt.Sprintln("error obteniendo compras concurrentemente: ", err))
	}

	resumen := s.obtenerResumen(compras)
	resumen.Normalizar()

	return resumen, nil
}

// obtenerCompras hace las llamadas a la API concurrentemente para obtener las compras del rango de días indicado
func (s *service) obtenerCompras(numDias int, date time.Time) ([]*Compra, error) {
	const urlStr = "https://apirecruit-gjvkhl2c6a-uc.a.run.app/compras/"

	var (
		wg             sync.WaitGroup
		errCh          = make(chan error, numDias)
		comprasTotales = make([][]*Compra, numDias)
	)

	for i := 0; i < numDias; i++ {
		wg.Add(1)

		// Lanzar una goroutine por cada dia con el fin de que se hagan las llamadas concurrentemente
		go func(i int) {
			defer wg.Done()
			dia := date.Add(time.Duration(i) * time.Hour * 24).Format("2006-01-02")

			resp, err := http.Get(urlStr + dia)
			if err != nil {
				errCh <- errors.New(fmt.Sprintln("error obteniendo url: ", err))
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				errCh <- errors.New(fmt.Sprintln("error leyendo body: ", err))
				return
			}

			comprasDia := []*Compra{}
			err = json.Unmarshal(body, &comprasDia)
			if err != nil {
				errCh <- errors.New(fmt.Sprintln("error unmarshaleando body: ", err))
				return
			}

			comprasTotales[i] = comprasDia
		}(i)
	}

	// Manejo de errores en las goroutines
	go func() {
		wg.Wait()
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return nil, err
	default:
	}

	wg.Wait()

	compras := []*Compra{}
	for _, ct := range comprasTotales {
		compras = append(compras, ct...)
	}

	return compras, nil
}

// obtenerResumen calcula el resumen de las compras
func (s *service) obtenerResumen(compras []*Compra) *Resumen {
	resumen := &Resumen{
		ComprasPorTDC: make(map[string]float64),
	}

	for _, compra := range compras {
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

	return resumen
}
