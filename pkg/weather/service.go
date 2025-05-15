// Package weather provides functionality for retrieving weather information
package weather

import (
	"time"
)

// Service defines the interface for weather information retrieval
type Service interface {
	// GetCityData retrieves all necessary data about a city
	GetCityData(city, country string) (CityData, error)

	// IsClearOrSunny checks if the weather description indicates clear or sunny weather
	IsClearOrSunny(weatherDesc string) bool

	// IsBusinessHours checks if the given time is within business hours (9am-5pm)
	IsBusinessHours(t time.Time) bool

	// FormatTemperature formats the temperature string according to the country
	FormatTemperature(city CityData) string
}
