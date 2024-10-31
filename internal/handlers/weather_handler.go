package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	service "github.com/brkcnr/getweatherapi/internal/services"
)

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	// Handle only GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "This code can only handle GET requests", http.StatusMethodNotAllowed)
		return
	}

	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Println("WEATHER_API_KEY not found in environment variables")
		http.Error(w, "API key not found", http.StatusInternalServerError)
		return
	}

	// Call weather service to get weather data
	weatherData, err := service.GetWeather(apiKey, city)
	if err.Code != 0 {
		http.Error(w, err.Message, err.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weatherData)
}
