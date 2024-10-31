package models

type Weather struct {
	City             string  `json:"city"`
	Region           string  `json:"region"`
	Country          string  `json:"country"`
	TimeZoneId       string  `json:"tz_id"`
	Temperature      float64 `json:"temperature"`
	FeelsLike        float64 `json:"feels_like"`
	WeatherCondition string  `json:"condition_text"`
}
