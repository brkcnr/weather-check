package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/brkcnr/getweatherapi/internal/models"
)

func GetWeather(apiKey, city string) (models.Weather, models.ErrorMessage) {
	baseURL := "http://api.weatherapi.com/v1/current.json"
	fullURL := fmt.Sprintf("%s?key=%s&q=%s&aqi=no", baseURL, apiKey, url.QueryEscape(city))
	log.Println("Requesting URL:", fullURL)

	response, err := http.Get(fullURL)
	if err != nil {
		return models.Weather{}, models.ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to make request"}
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		var apiResponse map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
			return models.Weather{}, models.ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to parse response"}
		}

		location, ok := apiResponse["location"].(map[string]interface{})
		if !ok {
			return models.Weather{}, models.ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to retrieve location data"}
		}
		cityName, ok := location["name"].(string)
		if !ok {
			return models.Weather{}, models.ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to retrieve city name"}
		}
		region, ok := location["region"].(string)
		if !ok {
			return models.Weather{}, models.ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to retrieve region"}
		}
		country, ok := location["country"].(string)
		if !ok {
			return models.Weather{}, models.ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to retrieve country name"}
		}
		timezoneId, ok := location["tz_id"].(string)
		if !ok {
			return models.Weather{}, models.ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to retrieve timezone"}
		}

		current, ok := apiResponse["current"].(map[string]interface{})
		if !ok {
			return models.Weather{}, models.ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to retrieve current weather data"}
		}
		tempC, ok := current["temp_c"].(float64)
		if !ok {
			return models.Weather{}, models.ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to retrieve temperature"}
		}
		feelsLike, ok := current["feelslike_c"].(float64)
		if !ok {
			return models.Weather{}, models.ErrorMessage{Code: http.StatusInternalServerError, Message: "Failed to retrieve feels like temperature"}
		}
		weatherCondition, ok := current["condition"].(map[string]interface{})
		if !ok {
			return models.Weather{}, models.ErrorMessage{Code: http.StatusInternalServerError, Message: "Weather Condition"}
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

	case http.StatusBadRequest:
		return models.Weather{}, models.ErrorMessage{Code: http.StatusBadRequest, Message: "Invalid city name. Please try again."}

	case http.StatusForbidden:
		var errorResponse map[string]interface{}
		json.NewDecoder(response.Body).Decode(&errorResponse)
		return models.Weather{}, models.ErrorMessage{Code: http.StatusForbidden, Message: "Forbidden access"}

	default:
		return models.Weather{}, models.ErrorMessage{Code: response.StatusCode, Message: "Unexpected status code"}
	}
}
