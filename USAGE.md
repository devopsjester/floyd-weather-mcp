# Using Floyd in VS Code

Floyd is an MCP (Model Context Protocol) server that integrates with GitHub Copilot in VS Code 1.100 and above. It provides weather information and deployment safety checks for any city in the world.

## Setup

1. The server is already configured in your VS Code settings.json file.
2. Make sure the Floyd executable is built:
   ```bash
   cd /Users/devopsjester/repos/experiment/mcp/floyd-weather-deployer
   go build -o floyd-weather-server
   ```

## Interacting with Floyd via GitHub Copilot

You can interact with Floyd through GitHub Copilot chat. Here are some examples of what you can ask:

1. **Get weather information**:
   - "What's the weather in Paris, France?"
   - "Tell me the current temperature in Tokyo, Japan."

2. **Check deployment safety**:
   - "Is it safe to deploy to London, United Kingdom?"
   - "Can I deploy to New York, USA right now?"

3. **Deploy to a city**:
   - "Deploy to Berlin, Germany."
   - "Deploy my app to Sydney, Australia."

## Understanding the Response

Floyd will provide:
- Current weather conditions
- Temperature (in Fahrenheit for US cities, Celsius for others)
- Whether it's safe to deploy based on:
  - Time of day (must be between 9am-5pm local time)
  - Weather conditions (must be clear or sunny)

## Troubleshooting

If Floyd is not responding or giving errors:
1. Make sure the executable is built and the path in settings.json is correct
2. Check internet connectivity (needed to access the OpenMeteo API)
3. Restart VS Code if necessary
