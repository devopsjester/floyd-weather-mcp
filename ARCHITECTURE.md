# Floyd Weather MCP Server Architecture

This document describes the architecture of the Floyd Weather MCP Server.

## Overview

Floyd is an MCP (Model Context Protocol) server that integrates with GitHub Copilot in VS Code 1.100 and above. It provides weather information and deployment safety checks for any city in the world.

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
