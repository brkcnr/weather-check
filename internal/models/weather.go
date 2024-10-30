package models

type Weather struct {
	City        string  `json:"city"`
	Region      string  `json:"region"`
	Country     string  `json:"country"`
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"feels_like"`
}
