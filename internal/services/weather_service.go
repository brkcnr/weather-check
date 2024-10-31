package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/brkcnr/getweatherapi/internal/models"
)

const (
	errRequestFailed    = "Failed to make request"
	errParseResponse    = "Failed to parse response"
	errLocationData     = "Failed to retrieve location data"
	errWeatherData      = "Failed to retrieve current weather data"
	errConditionData    = "Failed to retrieve weather condition data"
	errForbiddenAccess  = "Forbidden access"
	errUnexpectedStatus = "Unexpected status code"
)

func GetWeather(apiKey, city string) (models.Weather, models.ErrorMessage) {
	baseURL := "http://api.weatherapi.com/v1/current.json"
	fullURL := fmt.Sprintf("%s?key=%s&q=%s&aqi=no", baseURL, apiKey, url.QueryEscape(city))
	log.Println("Requesting URL:", fullURL)

	response, err := http.Get(fullURL)
	if err != nil {
		return models.Weather{}, models.ErrorMessage{Code: http.StatusInternalServerError, Message: errRequestFailed}
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		var apiResponse map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
			return models.Weather{}, models.ErrorMessage{Code: http.StatusInternalServerError, Message: errParseResponse}
		}

		location, ok := apiResponse["location"].(map[string]interface{})
		if !ok {
			return models.Weather{}, models.ErrorMessage{Code: http.StatusInternalServerError, Message: errLocationData}
		}
		cityName, _ := location["name"].(string)
		region, _ := location["region"].(string)
		country, _ := location["country"].(string)
		timezoneId, _ := location["tz_id"].(string)

		current, ok := apiResponse["current"].(map[string]interface{})
		if !ok {
			return models.Weather{}, models.ErrorMessage{Code: http.StatusInternalServerError, Message: errWeatherData}
		}
		tempC, _ := current["temp_c"].(float64)
		feelsLike, _ := current["feelslike_c"].(float64)

		weatherCondition, ok := current["condition"].(map[string]interface{})
		if !ok {
			return models.Weather{}, models.ErrorMessage{Code: http.StatusInternalServerError, Message: errConditionData}
		}
		weatherConditionText, _ := weatherCondition["text"].(string)

		return models.Weather{
			City:             cityName,
			Region:           region,
			Country:          country,
			Temperature:      tempC,
			FeelsLike:        feelsLike,
			TimeZoneId:       timezoneId,
			WeatherCondition: weatherConditionText,
		}, models.ErrorMessage{}

	case http.StatusForbidden:
		var errorResponse map[string]interface{}
		json.NewDecoder(response.Body).Decode(&errorResponse)
		return models.Weather{}, models.ErrorMessage{Code: http.StatusForbidden, Message: errForbiddenAccess}

	default:
		return models.Weather{}, models.ErrorMessage{Code: response.StatusCode, Message: errUnexpectedStatus}
	}
}
