# Floyd Weather MCP Server Architecture

This document describes the architecture of the Floyd Weather MCP Server.

## Overview

Floyd is an MCP (Model Context Protocol) server that integrates with GitHub Copilot in VS Code 1.100 and above. It provides weather information and deployment safety checks for any city in the world.

## SOLID Principles Implementation

The codebase has been restructured to follow SOLID principles:

### 1. Single Responsibility Principle (SRP)

Each package and class has a single responsibility:

- `weather` package: Responsible only for weather-related functionality
- `deployment` package: Responsible only for deployment-related functionality
- `mcp` package: Responsible only for handling MCP protocol communication

### 2. Open/Closed Principle (OCP)

The code is open for extension but closed for modification:

- New weather data sources can be added by implementing the `weather.Service` interface
- New deployment strategies can be added by implementing the `deployment.Service` interface
- New MCP methods can be added without modifying existing method handlers

### 3. Liskov Substitution Principle (LSP)

Services are defined through interfaces, allowing different implementations to be substituted:

- `weather.Service` can have multiple implementations (e.g., `APIService`, `MockService`)
- `deployment.Service` can be implemented differently for various deployment strategies

### 4. Interface Segregation Principle (ISP)

Interfaces are specific and focused:

- `weather.Service` only includes methods relevant to weather data
- `deployment.Service` only includes methods related to deployment
- `mcp.Handler` only includes methods for handling MCP requests

### 5. Dependency Inversion Principle (DIP)

High-level modules depend on abstractions, not concrete implementations:

- `deployment.Service` depends on the `weather.Service` interface, not a specific implementation
- `mcp.Handler` depends on service interfaces rather than concrete types
- `mcp.Server` depends on the `Handler` interface

## Architecture Diagram

```
┌────────────────┐        ┌─────────────────┐         ┌────────────────┐
│                │        │                 │         │                │
│  VS Code with  │        │  Floyd Weather  │         │  OpenMeteo API │
│  GitHub Copilot│◄──────►│  MCP Server     │◄───────►│                │
│                │        │                 │         │                │
└────────────────┘        └─────────────────┘         └────────────────┘
                                  ▲
                                  │
                                  ▼
                          ┌─────────────────┐

## Package Structure

The application is organized into the following packages:

- **pkg/weather**: Weather data retrieval and processing
  - `models.go`: Defines data structures for weather information
  - `service.go`: Defines the weather service interface
  - `api_service.go`: Implements the weather service using OpenMeteo API

- **pkg/deployment**: Deployment logic and safety checks
  - `service.go`: Defines and implements deployment services

- **pkg/mcp**: Model Context Protocol handling
  - `handler.go`: MCP request processing and routing
  - `server.go`: Server implementation for handling I/O

## Error Handling

The application follows Go's idiomatic error handling:

- Errors are properly propagated up the call stack
- API errors are handled and transformed into user-friendly messages
- Logging is used to track error details for debugging
                          │    Local Time   │
                          │    & Business   │
                          │    Hours Check  │
                          └─────────────────┘
```

## Component Description

### VS Code with GitHub Copilot
- Provides the user interface for interacting with the Floyd server
- Sends MCP requests to the server and displays the responses

### Floyd Weather MCP Server
- Processes natural language requests via MCP protocol
- Handles three main types of requests:
  1. Get weather information
  2. Check deployment safety
  3. Deploy to a city
- Returns appropriate responses in MCP format

### OpenMeteo API
- External service providing:
  - Geocoding (city coordinates)
  - Weather data
  - Timezone information

### Local Time & Business Hours Check
- Internal logic that:
  - Determines the current time in the target city's timezone
  - Checks if it's within business hours (9am-5pm)
  - Evaluates weather conditions for deployment safety

## Data Flow

1. User submits a query through GitHub Copilot in VS Code
2. Query is sent as an MCP request to the Floyd server
3. Floyd server:
   - Parses the request
   - Makes API calls to OpenMeteo for location, weather, and timezone data
   - Processes the data according to business rules
   - Formats the response according to MCP protocol
4. Response is sent back to VS Code
5. GitHub Copilot displays the result to the user
