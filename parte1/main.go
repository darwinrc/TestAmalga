package main

import (
	"TestAmalga/parte1/handler"
	"TestAmalga/parte1/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	service := service.NewService()
	handler := handler.NewHandler(service)
	handler.Attach(router)

	if err := http.ListenAndServe(":5000", router); err != nil {
		log.Fatal("ListenAndServe", err)
	}

	log.Println("listening on port 5000")
}
