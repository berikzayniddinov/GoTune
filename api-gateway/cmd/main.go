package main

import (
	"github.com/gorilla/mux"
	"gotune/api-gateway/internal/client"
	"gotune/api-gateway/internal/handler"
	"log"
	"net/http"
)

func main() {
	userClient := client.NewUserServiceClient("localhost:50051")
	userHandler := handler.NewUserHandler(userClient)

	instrumentClient := client.NewInstrumentServiceClient("localhost:50052")
	instrumentHandler := handler.NewInstrumentHandler(instrumentClient)

	r := mux.NewRouter()

	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")

	r.HandleFunc("/instruments", instrumentHandler.CreateInstrument).Methods("POST")
	r.HandleFunc("/instruments", instrumentHandler.GetAllInstruments).Methods("GET")
	r.HandleFunc("/instruments/{id}", instrumentHandler.DeleteInstrumentByID).Methods("DELETE")
	r.HandleFunc("/instruments/{id}", instrumentHandler.UpdateInstrumentByID).Methods("PATCH")

	log.Println("API Gateway запущен на порту 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Ошибка запуска API Gateway: %v", err)
	}
}
