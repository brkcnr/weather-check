package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	service "github.com/brkcnr/getweatherapi/internal/services"
)

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
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

	// Call weather.GetWeather instead of just GetWeather
	weatherData, err := service.GetWeather(apiKey, city)
	if err.Code != 0 {
		http.Error(w, err.Message, err.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weatherData)
}
