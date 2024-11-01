package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/brkcnr/getweatherapi/internal/models"
	"github.com/brkcnr/getweatherapi/internal/services"
	"github.com/brkcnr/getweatherapi/internal/werror"
)

// Helper function to send JSON-formatted errors
func sendJSONError(w http.ResponseWriter, err werror.WError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code())
	json.NewEncoder(w).Encode(models.ErrorMessage{Code: err.Code(), Message: err.Error()})

	// Log if the error is loggable
	if err.(*werror.Error).Loggable {
		log.Println("Error:", err.Error())
	}
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	// Handle only GET requests
	if r.Method != http.MethodGet {
		sendJSONError(w, werror.ErrMethodNotAllowed)
		return
	}

	city := r.URL.Query().Get("city")
	if city == "" {
		sendJSONError(w, werror.ErrCityParameterMissing)
		return
	}

	// Validate city length
	switch {
	case len(city) < 2:
		sendJSONError(w, werror.ErrCharacterLessThan)
		return
	case len(city) > 30:
		sendJSONError(w, werror.ErrCharacterMoreThan)
		return
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		sendJSONError(w, werror.ErrAPIKeyNotFound)
		return
	}

	// Call weather service to get weather data
	weatherData, err := services.GetWeather(apiKey, city)
	if err != nil {
		sendJSONError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weatherData)
}
