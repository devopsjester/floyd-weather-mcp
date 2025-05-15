# Floyd Weather Deployer MCP Server

[![Go Build](https://github.com/devopsjester/floyd-weather-mcp/actions/workflows/go.yml/badge.svg)](https://github.com/devopsjester/floyd-weather-mcp/actions/workflows/go.yml)

An MCP (Model Context Protocol) server that provides weather information and deployment safety checks, designed to be used with GitHub Copilot in Visual Studio Code 1.100+.

## Features

- Get current weather conditions for any city in the world
- Check if it's safe to deploy to a server in a specific city
- Deploy to a server if conditions are safe
- Temperature display in appropriate units (°F for US, °C elsewhere)

## Deployment Safety Rules

Deployment is considered safe when:
- It's between 9am and 5pm in the target city's local time
- The weather conditions are clear or sunny

## Installation

### Prerequisites

- Go 1.18 or higher
- VS Code 1.100 or higher with GitHub Copilot

### Setup

1. Clone this repository:
   ```bash
   git clone https://github.com/devopsjester/floyd-weather-mcp.git
   cd floyd-weather-mcp
   ```

2. Build the server:
   ```bash
   go build -o floyd-weather-server
   ```

3. Add to your VS Code `settings.json`:
   ```json
   "mcp": {
     "servers": {
       "floyd": {
         "command": "/path/to/floyd-weather-server"
       }
     }
   }
   ```

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
- Uses the OpenMeteo API for geocoding and weather data (free, public, no authentication required)
- Follows MCP protocol for interactions with GitHub Copilot

## API Details

The server uses these OpenMeteo API endpoints:
- Geocoding: `https://geocoding-api.open-meteo.com/v1/search`
- Weather: `https://api.open-meteo.com/v1/forecast`
- Timezone: Part of the forecast API

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
