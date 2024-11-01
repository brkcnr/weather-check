package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/brkcnr/getweatherapi/internal/handlers"
)

func TestWeatherHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/weather", nil)
	w := httptest.NewRecorder()

	handlers.WeatherHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, res.StatusCode)
	}
}

func TestWeatherHandler_MissingCityParameter(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/weather", nil)
	w := httptest.NewRecorder()

	handlers.WeatherHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, res.StatusCode)
	}
}

func TestWeatherHandler_CityLengthTooShort(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/weather?city=a", nil)
	w := httptest.NewRecorder()

	handlers.WeatherHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, res.StatusCode)
	}
}

func TestWeatherHandler_CityLengthTooLong(t *testing.T) {
	longCity := strings.Repeat("a", 31)
	req := httptest.NewRequest(http.MethodGet, "/weather?city="+longCity, nil)
	w := httptest.NewRecorder()

	handlers.WeatherHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, res.StatusCode)
	}
}

func TestWeatherHandler_MissingAPIKey(t *testing.T) {
	// Unset API key if exists
	os.Unsetenv("WEATHER_API_KEY")

	req := httptest.NewRequest(http.MethodGet, "/weather?city=Istanbul", nil)
	w := httptest.NewRecorder()

	handlers.WeatherHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, res.StatusCode)
	}
}
