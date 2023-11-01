package handler

import (
	"TestAmalga/parte1/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler interface {
	Attach(r *mux.Router)
	HandleGet(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	Service service.Service
}

func NewHandler(s service.Service) Handler {
	return &handler{
		Service: s,
	}
}

func (h *handler) Attach(r *mux.Router) {
	r.HandleFunc("/resumen/{fecha}", h.HandleGet).Methods("GET", "OPTIONS")
}

func (h *handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	fecha := vars["fecha"]

	queryParams := r.URL.Query()
	dias := queryParams.Get("dias")

	if dias == "" {
		dias = "0"
	}

	resumen, err := h.Service.CalcularResumen(fecha, dias)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
		return
	}

	resumen.Normalizar()

	res, err := json.Marshal(resumen)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
