// Package mcp provides functionality for the Model Context Protocol
package mcp

import (
	"encoding/json"
	"log"

	"github.com/devopsjester/floyd-weather-deployer/pkg/deployment"
	"github.com/devopsjester/floyd-weather-deployer/pkg/weather"
)

// Request represents the request format from MCP
type Request struct {
	Method     string          `json:"method"`
	Parameters json.RawMessage `json:"parameters"`
}

// Response represents the response format for MCP
type Response struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// CityParams contains parameters with city and country
type CityParams struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

// Handler defines the interface for handling MCP requests
type Handler interface {
	// ProcessRequest handles a single MCP request
	ProcessRequest(request Request) Response
}

// DefaultHandler implements the Handler interface
type DefaultHandler struct {
	weatherSvc    weather.Service
	deploymentSvc deployment.Service
}

// NewHandler creates a new MCP handler
func NewHandler(weatherSvc weather.Service, deploymentSvc deployment.Service) *DefaultHandler {
	return &DefaultHandler{
		weatherSvc:    weatherSvc,
		deploymentSvc: deploymentSvc,
	}
}

// ProcessRequest handles a single MCP request
func (h *DefaultHandler) ProcessRequest(request Request) Response {
	log.Printf("Processing request method: %s", request.Method)

	switch request.Method {
	case "check-deployment-safety":
		log.Printf("Handling check-deployment-safety")
		response := h.handleCheckDeploymentSafety(request.Parameters)
		log.Printf("Completed check-deployment-safety")
		return response

	case "deploy-to-city":
		log.Printf("Handling deploy-to-city")
		response := h.handleDeployToCity(request.Parameters)
		log.Printf("Completed deploy-to-city")
		return response

	case "get-weather":
		log.Printf("Handling get-weather")
		response := h.handleGetWeather(request.Parameters)
		log.Printf("Completed get-weather")
		return response

	default:
		log.Printf("Unknown method: %s", request.Method)
		return h.createErrorResponse("Unknown method: " + request.Method)
	}
}

// handleCheckDeploymentSafety processes the deployment safety check
func (h *DefaultHandler) handleCheckDeploymentSafety(params json.RawMessage) Response {
	var cityParams CityParams
	if err := json.Unmarshal(params, &cityParams); err != nil {
		log.Printf("Error parsing parameters: %v", err)
		return h.createErrorResponse("Error parsing parameters: " + err.Error())
	}

	log.Printf("Getting city data for %s, %s", cityParams.City, cityParams.Country)
	cityData, err := h.weatherSvc.GetCityData(cityParams.City, cityParams.Country)
	if err != nil {
		log.Printf("Error getting city data: %v", err)
		return h.createErrorResponse("Error getting city data: " + err.Error())
	}
	log.Printf("Successfully retrieved data for %s, %s", cityParams.City, cityParams.Country)

	safety := h.deploymentSvc.CheckSafety(cityData)
	return Response{
		Type:    "success",
		Content: safety,
	}
}

// handleDeployToCity processes the deployment request
func (h *DefaultHandler) handleDeployToCity(params json.RawMessage) Response {
	var cityParams CityParams
	if err := json.Unmarshal(params, &cityParams); err != nil {
		return h.createErrorResponse("Error parsing parameters: " + err.Error())
	}

	cityData, err := h.weatherSvc.GetCityData(cityParams.City, cityParams.Country)
	if err != nil {
		return h.createErrorResponse("Error getting city data: " + err.Error())
	}

	deployed, message := h.deploymentSvc.Deploy(cityData)

	result := map[string]interface{}{
		"deployed": deployed,
		"message":  message,
	}

	return Response{
		Type:    "success",
		Content: result,
	}
}

// handleGetWeather processes the weather information request
func (h *DefaultHandler) handleGetWeather(params json.RawMessage) Response {
	var cityParams CityParams
	if err := json.Unmarshal(params, &cityParams); err != nil {
		return h.createErrorResponse("Error parsing parameters: " + err.Error())
	}

	cityData, err := h.weatherSvc.GetCityData(cityParams.City, cityParams.Country)
	if err != nil {
		return h.createErrorResponse("Error getting city data: " + err.Error())
	}

	// Format temperature according to country
	tempStr := h.weatherSvc.FormatTemperature(cityData)

	result := map[string]interface{}{
		"city":    cityParams.City,
		"country": cityParams.Country,
		"weather": cityData.Weather,
		"temp":    tempStr,
	}

	return Response{
		Type:    "success",
		Content: result,
	}
}

// createErrorResponse creates an error response in MCP format
func (h *DefaultHandler) createErrorResponse(message string) Response {
	return Response{
		Type: "error",
		Content: map[string]string{
			"message": message,
		},
	}
}
