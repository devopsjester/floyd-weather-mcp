// Main package for the Floyd Weather MCP Server
package main

import (
	"fmt"
	"os"

	"github.com/devopsjester/floyd-weather-deployer/pkg/deployment"
	"github.com/devopsjester/floyd-weather-deployer/pkg/mcp"
	"github.com/devopsjester/floyd-weather-deployer/pkg/weather"
)

func main() {
	// Initialize services
	weatherService := weather.NewAPIService()
	deploymentService := deployment.NewService(weatherService)
	
	// Create handler and server
	mcpHandler := mcp.NewHandler(weatherService, deploymentService)
	mcpServer, err := mcp.NewServer(mcpHandler)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing server: %v\n", err)
		os.Exit(1)
	}
	
	// Start serving requests
	mcpServer.Serve()
}
