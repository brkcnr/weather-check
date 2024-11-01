package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/brkcnr/getweatherapi/internal/models"
	"github.com/brkcnr/getweatherapi/internal/werror"
)

func GetWeather(apiKey, city string) (models.Weather, werror.WError) {
	baseURL := "http://api.weatherapi.com/v1/current.json"
	fullURL := fmt.Sprintf("%s?key=%s&q=%s&aqi=no", baseURL, apiKey, url.QueryEscape(city))
	log.Printf("Requesting URL: %s/%s", baseURL, city)

	response, err := http.Get(fullURL)
	if err != nil {
		return models.Weather{}, werror.ErrRequestFailed.Wrap(err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		var apiResponse map[string]any
		if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
			return models.Weather{}, werror.ErrParseResponse.Wrap(err)
		}

		location, ok := apiResponse["location"].(map[string]any)
		if !ok {
			return models.Weather{}, werror.ErrLocationDataMissing
		}
		cityName, _ := location["name"].(string)
		region, _ := location["region"].(string)
		country, _ := location["country"].(string)
		timezoneId, _ := location["tz_id"].(string)

		current, ok := apiResponse["current"].(map[string]any)
		if !ok {
			return models.Weather{}, werror.ErrWeatherDataMissing
		}
		tempC, _ := current["temp_c"].(float64)
		feelsLike, _ := current["feelslike_c"].(float64)

		weatherCondition, ok := current["condition"].(map[string]any)
		if !ok {
			return models.Weather{}, werror.ErrConditionDataMissing
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
		}, nil

	case http.StatusBadRequest:
		return models.Weather{}, werror.ErrInvalidCity

	case http.StatusForbidden:
		var errorResponse map[string]any
		json.NewDecoder(response.Body).Decode(&errorResponse)
		return models.Weather{}, werror.ErrForbiddenAccess

	default:
		return models.Weather{}, werror.ErrUnexpectedStatus
	}
}
