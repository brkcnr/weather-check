package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/brkcnr/getweatherapi/internal/models"
	service "github.com/brkcnr/getweatherapi/internal/services"
)

const (
	errCityParameterMissing = "City parameter is missing"
	errAPIKeyNotFound       = "API key not found"
)

func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorMessage{Code: statusCode, Message: message})
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {

	city := r.URL.Query().Get("city")
	if city == "" {
		sendJSONError(w, errCityParameterMissing, http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Println("WEATHER_API_KEY not found in environment variables")
		sendJSONError(w, errAPIKeyNotFound, http.StatusInternalServerError)
		return
	}

	// Call weather service to get weather data
	weatherData, err := service.GetWeather(apiKey, city)
	if err.Code != 0 {
		sendJSONError(w, err.Message, err.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weatherData)
}
