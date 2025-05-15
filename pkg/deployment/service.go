// Package deployment provides functionality for checking deployment safety
package deployment

import (
	"fmt"

	"github.com/devopsjester/floyd-weather-deployer/pkg/weather"
)

// Safety represents the deployment safety status
type Safety struct {
	Safe    bool   `json:"safe"`
	Reason  string `json:"reason"`
	Weather string `json:"weather"`
	Temp    string `json:"temp"`
}

// Service defines the interface for deployment operations
type Service interface {
	// CheckSafety determines if it's safe to deploy based on time and weather
	CheckSafety(city weather.CityData) Safety

	// Deploy attempts to deploy to the specified city
	Deploy(city weather.CityData) (bool, string)
}

// DefaultService implements the Service interface
type DefaultService struct {
	weatherSvc weather.Service
}

// NewService creates a new deployment service
func NewService(weatherSvc weather.Service) *DefaultService {
	return &DefaultService{
		weatherSvc: weatherSvc,
	}
}

// CheckSafety determines if it's safe to deploy based on time and weather
func (s *DefaultService) CheckSafety(city weather.CityData) Safety {
	// Format temperature according to country
	tempStr := s.weatherSvc.FormatTemperature(city)

	// Check if it's business hours (9am-5pm)
	if !s.weatherSvc.IsBusinessHours(city.LocalTime) {
		localTimeStr := city.LocalTime.Format("3:04 PM")
		return Safety{
			Safe:    false,
			Reason:  fmt.Sprintf("Outside of business hours (current time is %s)", localTimeStr),
			Weather: city.Weather,
			Temp:    tempStr,
		}
	}

	// Check if weather is clear or sunny
	if !s.weatherSvc.IsClearOrSunny(city.Weather) {
		return Safety{
			Safe:    false,
			Reason:  fmt.Sprintf("Weather conditions are not clear/sunny (current: %s)", city.Weather),
			Weather: city.Weather,
			Temp:    tempStr,
		}
	}

	// All conditions are met
	return Safety{
		Safe:    true,
		Reason:  "Business hours and clear/sunny weather",
		Weather: city.Weather,
		Temp:    tempStr,
	}
}

// Deploy attempts to deploy to the specified city
func (s *DefaultService) Deploy(city weather.CityData) (bool, string) {
	safety := s.CheckSafety(city)

	if safety.Safe {
		return true, fmt.Sprintf("Successfully deployed to %s, %s. Current weather: %s. Temperature: %s.",
			city.Name, city.Country, safety.Weather, safety.Temp)
	}

	return false, fmt.Sprintf("Could not deploy to %s, %s: %s. Current weather: %s. Temperature: %s.",
		city.Name, city.Country, safety.Reason, safety.Weather, safety.Temp)
}
