package main

import (
	"log"
	"net/http"

	_ "github.com/brkcnr/getweatherapi/internal/env"
	"github.com/brkcnr/getweatherapi/internal/handlers"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/weather", handlers.WeatherHandler)

	log.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", mux)
}
