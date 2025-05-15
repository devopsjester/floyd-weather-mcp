# Floyd Weather Deployer MCP Server

An MCP (Model Context Protocol) server that provides weather information and deployment safety checks.

## Features

- Get current weather conditions for any city in the world
- Check if it's safe to deploy to a server in a specific city
- Deploy to a server if conditions are safe

## Deployment Safety Rules

Deployment is considered safe when:
- It's between 9am and 5pm in the target city's local time
- The weather conditions are clear or sunny

## Usage

This server is designed to be used with GitHub Copilot through VS Code MCP integration.

### Available Commands

1. **Check deployment safety**:
   Ask about deployment safety for a city.

2. **Get weather information**:
   Ask about weather conditions in a city.

3. **Deploy to a city**:
   Request deployment to a city.

### Examples

- "Is it safe to deploy to London, United Kingdom?"
- "What's the weather in New York, USA?"
- "Deploy to Tokyo, Japan"

## Technical Details

- Built in Go
- Uses the OpenMeteo API for geocoding and weather data
- Follows MCP protocol for interactions with GitHub Copilot
