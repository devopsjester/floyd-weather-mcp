// Package weather provides functionality for retrieving weather information
package weather

import (
	"time"
)

// Condition represents possible weather conditions
type Condition int

const (
	Clear Condition = iota
	Sunny
	Cloudy
	Rainy
	Stormy
	Unknown
)

// CityData contains location data for a city
type CityData struct {
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	TimeZone  string    `json:"timezone"`
	LocalTime time.Time `json:"localTime"`
	Weather   string    `json:"weather"`
	TempC     float64   `json:"tempC"`
	TempF     float64   `json:"tempF"`
}
