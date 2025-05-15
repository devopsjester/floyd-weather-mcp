package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// WeatherCondition represents possible weather conditions
type WeatherCondition int

const (
	Clear WeatherCondition = iota
	Sunny
	Cloudy
	Rainy
	Stormy
	Unknown
)

// DeploymentSafety represents the deployment safety status
type DeploymentSafety struct {
	Safe    bool   `json:"safe"`
	Reason  string `json:"reason"`
	Weather string `json:"weather"`
	Temp    string `json:"temp"`
}

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

// McpRequest represents the request format from MCP (Model Context Protocol)
type McpRequest struct {
	Method     string          `json:"method"`
	Parameters json.RawMessage `json:"parameters"`
}

// CheckDeploymentParams contains parameters for checking deployment safety
type CheckDeploymentParams struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

// McpResponse represents the response format for MCP (Model Context Protocol)
type McpResponse struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

func main() {
	// Read from stdin and write to stdout as per MCP protocol
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var request McpRequest
		if err := decoder.Decode(&request); err != nil {
			if err == io.EOF {
				break
			}
			sendErrorResponse(encoder, "Error decoding request: "+err.Error())
			continue
		}

		switch request.Method {
		case "check-deployment-safety":
			handleCheckDeploymentSafety(encoder, request.Parameters)
		case "deploy-to-city":
			handleDeployToCity(encoder, request.Parameters)
		case "get-weather":
			handleGetWeather(encoder, request.Parameters)
		default:
			sendErrorResponse(encoder, "Unknown method: "+request.Method)
		}
	}
}

// handleCheckDeploymentSafety processes the deployment safety check
func handleCheckDeploymentSafety(encoder *json.Encoder, params json.RawMessage) {
	var checkParams CheckDeploymentParams
	if err := json.Unmarshal(params, &checkParams); err != nil {
		sendErrorResponse(encoder, "Error parsing parameters: "+err.Error())
		return
	}

	cityData, err := getCityData(checkParams.City, checkParams.Country)
	if err != nil {
		sendErrorResponse(encoder, "Error getting city data: "+err.Error())
		return
	}

	safety := checkDeploymentSafety(cityData)
	encoder.Encode(McpResponse{
		Type:    "success",
		Content: safety,
	})
}

// handleDeployToCity processes the deployment request
func handleDeployToCity(encoder *json.Encoder, params json.RawMessage) {
	var deployParams CheckDeploymentParams
	if err := json.Unmarshal(params, &deployParams); err != nil {
		sendErrorResponse(encoder, "Error parsing parameters: "+err.Error())
		return
	}

	cityData, err := getCityData(deployParams.City, deployParams.Country)
	if err != nil {
		sendErrorResponse(encoder, "Error getting city data: "+err.Error())
		return
	}

	safety := checkDeploymentSafety(cityData)

	var result map[string]interface{}
	if safety.Safe {
		result = map[string]interface{}{
			"deployed": true,
			"message": fmt.Sprintf("Successfully deployed to %s, %s. Current weather: %s. Temperature: %s.",
				deployParams.City, deployParams.Country, safety.Weather, safety.Temp),
		}
	} else {
		result = map[string]interface{}{
			"deployed": false,
			"message": fmt.Sprintf("Could not deploy to %s, %s: %s. Current weather: %s. Temperature: %s.",
				deployParams.City, deployParams.Country, safety.Reason, safety.Weather, safety.Temp),
		}
	}

	encoder.Encode(McpResponse{
		Type:    "success",
		Content: result,
	})
}

// handleGetWeather processes the weather information request
func handleGetWeather(encoder *json.Encoder, params json.RawMessage) {
	var weatherParams CheckDeploymentParams
	if err := json.Unmarshal(params, &weatherParams); err != nil {
		sendErrorResponse(encoder, "Error parsing parameters: "+err.Error())
		return
	}

	cityData, err := getCityData(weatherParams.City, weatherParams.Country)
	if err != nil {
		sendErrorResponse(encoder, "Error getting city data: "+err.Error())
		return
	}

	// Format temperature according to country
	tempStr := formatTemperature(cityData)

	result := map[string]interface{}{
		"city":    weatherParams.City,
		"country": weatherParams.Country,
		"weather": cityData.Weather,
		"temp":    tempStr,
	}

	encoder.Encode(McpResponse{
		Type:    "success",
		Content: result,
	})
}

// getCityData retrieves all necessary data about a city
func getCityData(city, country string) (CityData, error) {
	// Step 1: Geocoding to get coordinates
	lat, lon, err := getGeoLocation(city, country)
	if err != nil {
		return CityData{}, err
	}

	// Step 2: Get weather data
	weather, tempC, err := getCurrentWeather(lat, lon)
	if err != nil {
		return CityData{}, err
	}

	// Step 3: Get timezone and current time
	timezone, err := getTimeZone(lat, lon)
	if err != nil {
		return CityData{}, err
	}

	localTime := getCurrentLocalTime(timezone)
	tempF := celsiusToFahrenheit(tempC)

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

// getGeoLocation gets the latitude and longitude for a city
func getGeoLocation(city, country string) (float64, float64, error) {
	// URL encode the city name to handle multi-word cities
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json",
		encodedCity)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
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
func getCurrentWeather(lat, lon float64) (string, float64, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,weather_code",
		lat, lon)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
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

	weatherDesc := weatherCodeToDescription(response.Current.WeatherCode)
	return weatherDesc, response.Current.Temperature, nil
}

// getTimeZone gets the timezone for a location
func getTimeZone(lat, lon float64) (string, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&timezone=auto",
		lat, lon)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
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
func getCurrentLocalTime(timezone string) time.Time {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		// Default to UTC if timezone cannot be loaded
		return time.Now().UTC()
	}
	return time.Now().In(loc)
}

// weatherCodeToDescription converts the Open-Meteo weather code to a human-readable description
func weatherCodeToDescription(code int) string {
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

// isClearOrSunny checks if the weather description indicates clear or sunny weather
func isClearOrSunny(weatherDesc string) bool {
	clearConditions := []string{"Clear", "Mainly clear", "Clear sky", "Sunny"}
	for _, condition := range clearConditions {
		if weatherDesc == condition {
			return true
		}
	}
	return false
}

// isBusinessHours checks if the given time is within business hours (9am-5pm)
func isBusinessHours(t time.Time) bool {
	hour := t.Hour()
	return hour >= 9 && hour < 17
}

// checkDeploymentSafety determines if it's safe to deploy based on time and weather
func checkDeploymentSafety(city CityData) DeploymentSafety {
	// Format temperature according to country
	tempStr := formatTemperature(city)

	// Check if it's business hours (9am-5pm)
	if !isBusinessHours(city.LocalTime) {
		localTimeStr := city.LocalTime.Format("3:04 PM")
		return DeploymentSafety{
			Safe:    false,
			Reason:  fmt.Sprintf("Outside of business hours (current time is %s)", localTimeStr),
			Weather: city.Weather,
			Temp:    tempStr,
		}
	}

	// Check if weather is clear or sunny
	if !isClearOrSunny(city.Weather) {
		return DeploymentSafety{
			Safe:    false,
			Reason:  fmt.Sprintf("Weather conditions are not clear/sunny (current: %s)", city.Weather),
			Weather: city.Weather,
			Temp:    tempStr,
		}
	}

	// All conditions are met
	return DeploymentSafety{
		Safe:    true,
		Reason:  "Business hours and clear/sunny weather",
		Weather: city.Weather,
		Temp:    tempStr,
	}
}

// celsiusToFahrenheit converts Celsius to Fahrenheit
func celsiusToFahrenheit(celsius float64) float64 {
	return celsius*9/5 + 32
}

// formatTemperature formats the temperature string according to the country
func formatTemperature(city CityData) string {
	if city.Country == "United States" || city.Country == "USA" {
		return fmt.Sprintf("%.1f°F", city.TempF)
	}
	return fmt.Sprintf("%.1f°C", city.TempC)
}

// sendErrorResponse sends an error response in MCP format
func sendErrorResponse(encoder *json.Encoder, message string) {
	encoder.Encode(McpResponse{
		Type: "error",
		Content: map[string]string{
			"message": message,
		},
	})
}
