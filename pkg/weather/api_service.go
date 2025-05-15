// Package weather provides functionality for retrieving weather information
package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// APIService implements the Service interface using OpenMeteo API
type APIService struct {
	httpClient *http.Client
}

// NewAPIService creates a new weather service instance
func NewAPIService() *APIService {
	return &APIService{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetCityData retrieves all necessary data about a city
func (s *APIService) GetCityData(city, country string) (CityData, error) {
	// Step 1: Geocoding to get coordinates
	lat, lon, err := s.getGeoLocation(city, country)
	if err != nil {
		return CityData{}, err
	}

	// Step 2: Get weather data
	weather, tempC, err := s.getCurrentWeather(lat, lon)
	if err != nil {
		return CityData{}, err
	}

	// Step 3: Get timezone and current time
	timezone, err := s.getTimeZone(lat, lon)
	if err != nil {
		return CityData{}, err
	}

	localTime := s.getCurrentLocalTime(timezone)
	tempF := s.celsiusToFahrenheit(tempC)

	return CityData{
		Name:      city,
		Country:   country,
		Latitude:  lat,
		Longitude: lon,
		TimeZone:  timezone,
		LocalTime: localTime,
		Weather:   weather,
		TempC:     tempC,
		TempF:     tempF,
	}, nil
}

// IsClearOrSunny checks if the weather description indicates clear or sunny weather
func (s *APIService) IsClearOrSunny(weatherDesc string) bool {
	clearConditions := []string{"Clear", "Mainly clear", "Clear sky", "Sunny"}
	for _, condition := range clearConditions {
		if weatherDesc == condition {
			return true
		}
	}
	return false
}

// IsBusinessHours checks if the given time is within business hours (9am-5pm)
func (s *APIService) IsBusinessHours(t time.Time) bool {
	hour := t.Hour()
	return hour >= 9 && hour < 17
}

// FormatTemperature formats the temperature string according to the country
func (s *APIService) FormatTemperature(city CityData) string {
	if city.Country == "United States" || city.Country == "USA" {
		return fmt.Sprintf("%.1f°F", city.TempF)
	}
	return fmt.Sprintf("%.1f°C", city.TempC)
}

// getGeoLocation gets the latitude and longitude for a city
func (s *APIService) getGeoLocation(city, country string) (float64, float64, error) {
	// URL encode the city name to handle multi-word cities
	encodedCity := url.QueryEscape(city)
	apiURL := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json",
		encodedCity)

	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		return 0, 0, fmt.Errorf("geocoding API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("geocoding API returned non-200 status: %d", resp.StatusCode)
	}

	var response struct {
		Results []struct {
			Name      string  `json:"name"`
			Country   string  `json:"country"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, 0, fmt.Errorf("failed to parse geocoding response: %w", err)
	}

	if len(response.Results) == 0 {
		return 0, 0, fmt.Errorf("city not found: %s, %s", city, country)
	}

	// Find the result that matches the country
	for _, result := range response.Results {
		if result.Country == country {
			return result.Latitude, result.Longitude, nil
		}
	}

	// If no exact country match, use the first result
	return response.Results[0].Latitude, response.Results[0].Longitude, nil
}

// getCurrentWeather gets the current weather conditions for a location
func (s *APIService) getCurrentWeather(lat, lon float64) (string, float64, error) {
	apiURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,weather_code",
		lat, lon)

	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		return "", 0, fmt.Errorf("weather API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("weather API returned non-200 status: %d", resp.StatusCode)
	}

	var response struct {
		Current struct {
			Temperature float64 `json:"temperature_2m"`
			WeatherCode int     `json:"weather_code"`
		} `json:"current"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", 0, fmt.Errorf("failed to parse weather response: %w", err)
	}

	weatherDesc := s.weatherCodeToDescription(response.Current.WeatherCode)
	return weatherDesc, response.Current.Temperature, nil
}

// getTimeZone gets the timezone for a location
func (s *APIService) getTimeZone(lat, lon float64) (string, error) {
	apiURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&timezone=auto",
		lat, lon)

	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("timezone API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("timezone API returned non-200 status: %d", resp.StatusCode)
	}

	var response struct {
		Timezone string `json:"timezone"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to parse timezone response: %w", err)
	}

	return response.Timezone, nil
}

// getCurrentLocalTime gets the current time in a specific timezone
func (s *APIService) getCurrentLocalTime(timezone string) time.Time {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		// Default to UTC if timezone cannot be loaded
		return time.Now().UTC()
	}
	return time.Now().In(loc)
}

// weatherCodeToDescription converts the Open-Meteo weather code to a human-readable description
func (s *APIService) weatherCodeToDescription(code int) string {
	// WMO Weather interpretation codes (WW)
	// https://www.nodc.noaa.gov/archive/arc0021/0002199/1.1/data/0-data/HTML/WMO-CODE/WMO4677.HTM
	switch {
	case code == 0:
		return "Clear sky"
	case code == 1:
		return "Mainly clear"
	case code == 2:
		return "Partly cloudy"
	case code == 3:
		return "Overcast"
	case code >= 45 && code <= 48:
		return "Fog"
	case code >= 51 && code <= 55:
		return "Drizzle"
	case code >= 56 && code <= 57:
		return "Freezing Drizzle"
	case code >= 61 && code <= 65:
		return "Rain"
	case code >= 66 && code <= 67:
		return "Freezing Rain"
	case code >= 71 && code <= 75:
		return "Snow fall"
	case code == 77:
		return "Snow grains"
	case code >= 80 && code <= 82:
		return "Rain showers"
	case code >= 85 && code <= 86:
		return "Snow showers"
	case code == 95:
		return "Thunderstorm"
	case code == 96 || code == 99:
		return "Thunderstorm with hail"
	default:
		return "Unknown weather condition"
	}
}

// celsiusToFahrenheit converts Celsius to Fahrenheit
func (s *APIService) celsiusToFahrenheit(celsius float64) float64 {
	return celsius*9/5 + 32
}
