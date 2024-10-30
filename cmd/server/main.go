package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Weather struct {
	City        string  `json:"city"`
	Region      string  `json:"region"`
	Country     string  `json:"country"`
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"feels_like"`
}

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func init() {
	if err := godotenv.Load("cmd/server/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetWeather(apiKey, city string) (Weather, ErrorMessage) {
	baseURL := "http://api.weatherapi.com/v1/current.json"
	fullURL := fmt.Sprintf("%s?key=%s&q=%s&aqi=no", baseURL, apiKey, url.QueryEscape(city))

	response, err := http.Get(fullURL)
	if err != nil {
		return Weather{}, ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to make request"}
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		var apiResponse map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
			return Weather{}, ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to parse response"}
		}

		location, ok := apiResponse["location"].(map[string]interface{})
		if !ok {
			return Weather{}, ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to retrieve location data"}
		}
		cityName, ok := location["name"].(string)
		if !ok {
			return Weather{}, ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to retrieve city name"}
		}
		region, ok := location["region"].(string)
		if !ok {
			return Weather{}, ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to retrieve region"}
		}
		country, ok := location["country"].(string)
		if !ok {
			return Weather{}, ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to retrieve country name"}
		}

		current, ok := apiResponse["current"].(map[string]interface{})
		if !ok {
			return Weather{}, ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to retrieve current weather data"}
		}
		tempC, ok := current["temp_c"].(float64)
		if !ok {
			return Weather{}, ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to retrieve temperature"}
		}
		feelsLike, ok := current["feelslike_c"].(float64)
		if !ok {
			return Weather{}, ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to retrieve feels like temperature"}
		}

		return Weather{
			City:        cityName,
			Region:      region,
			Country:     country,
			Temperature: tempC,
			FeelsLike:   feelsLike,
		}, ErrorMessage{}

	case http.StatusBadRequest:
		return Weather{}, ErrorMessage{Code: http.StatusBadRequest, Message: "Invalid city name. Please try again."}

	case http.StatusForbidden:
		var errorResponse map[string]interface{}
		json.NewDecoder(response.Body).Decode(&errorResponse)
		return Weather{}, ErrorMessage{Code: http.StatusForbidden, Message: "Forbidden access"}

	default:
		return Weather{}, ErrorMessage{Code: response.StatusCode, Message: "Unexpected status code"}
	}
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		http.Error(w, "API key not found", http.StatusInternalServerError)
		return
	}

	weather, err := GetWeather(apiKey, city)
	if err.Code != 0 {
		http.Error(w, err.Message, err.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	htmlPath := filepath.Join("static", "index.html")
	http.ServeFile(w, r, htmlPath)
}

func main() {
	http.HandleFunc("/weather", WeatherHandler)
	http.HandleFunc("/", HomePageHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
