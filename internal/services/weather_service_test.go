package services_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brkcnr/getweatherapi/internal/services"
	"github.com/brkcnr/getweatherapi/internal/werror"
)

func MockServer(responseBody string, statusCode int) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(responseBody))
	})
	return httptest.NewServer(handler)
}

func TestGetWeather_Succes(t *testing.T) {
	mockResponse := `{
		"location": {"name": "Istanbul", "region": "Istanbul", "country": "Turkey", "tz_id": "Europe/Istanbul"},
		"current": {"temp_c": 19.0, "feelslike_c": 17.5, "condition": {"text": "Sunny"}}
	}`
	mockServer := MockServer(mockResponse, http.StatusOK)
	defer mockServer.Close()

	// Override base URL
	services.BaseURL = mockServer.URL

	apiKey := "dummyKey"
	city := "Istanbul"
	weatherData, err := services.GetWeather(apiKey, city)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if weatherData.City != "Istanbul" || weatherData.Temperature != 19.0 {
		t.Errorf("unexpected weather data: %+v", weatherData)
	}
}

func TestGetWeather_InvalidCity(t *testing.T) {
	mockServer := MockServer(`{"error": {"message": "Invalid city name"}}`, http.StatusBadRequest)
	defer mockServer.Close()

	// Override base URL
	services.BaseURL = mockServer.URL

	apiKey := "dummyKey"
	city := "InvalidCity"
	_, err := services.GetWeather(apiKey, city)

	if err == nil || err != werror.ErrInvalidCity {
		t.Fatalf("expected ErrInvalidCity, got %v", err)
	}
}

func TestGetWeather_ForbiddenAccess(t *testing.T) {
	mockServer := MockServer(`{"error": {"message": "Forbidden"}}`, http.StatusForbidden)
	defer mockServer.Close()

	// Override base URL
	services.BaseURL = mockServer.URL

	apiKey := "dummyKey"
	city := "Istanbul"
	_, err := services.GetWeather(apiKey, city)

	if err == nil || err != werror.ErrForbiddenAccess {
		t.Fatalf("expected ErrForbiddenAccess, got %v", err)
	}
}
