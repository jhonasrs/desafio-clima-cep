package main

import (
	"desafio-clima-cep/internal/client"
	"desafio-clima-cep/internal/handler"
	"desafio-clima-cep/internal/service"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	viaCEPClient := client.NewHTTPViaCEPClient("")
	weatherAPIClient := client.NewHTTPWeatherAPIClient("")
	weatherService := service.NewDefaultWeatherService(viaCEPClient, weatherAPIClient)
	weatherHandler := handler.NewWeatherHandler(weatherService)

	mux := http.NewServeMux()
	mux.Handle("/cep/", weatherHandler)

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
