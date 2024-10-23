package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		Temperature float64 `json:"temp_c"`
		FeelsLike   float64 `json:"feelslike_c"`
	} `json:"current"`
}

func init() {
	// Load the .env file at the beginning of the application
	err := godotenv.Load("cmd/server/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the "city" query parameter from the request URL
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
		return
	}

	// Fetch API key from the environment
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		http.Error(w, "API key not found", http.StatusInternalServerError)
		return
	}

	// Create the API request URL using the API key and city
	apiURL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, city)

	// Make the HTTP request to the weather API
	res, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the API response has an error status code
	if res.StatusCode > 299 {
		http.Error(w, fmt.Sprintf("Response failed with status code: %d and body: %s", res.StatusCode, body), http.StatusInternalServerError)
		return
	}

	// Unmarshal the response body into the Weather struct
	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		http.Error(w, "Failed to unmarshal response", http.StatusInternalServerError)
		return
	}

	// Prepare the response map
	responseMap := map[string]any{
		"city":        weather.Location.Name,
		"country":     weather.Location.Country,
		"temperature": weather.Current.Temperature,
		"feels_like":  weather.Current.FeelsLike,
	}

	// Set the content type to JSON and encode the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(responseMap)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
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
