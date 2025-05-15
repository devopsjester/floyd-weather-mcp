#!/bin/bash

# run.sh - Helper script for running Floyd Weather MCP Server commands
# Usage: ./run.sh <command> <city> <country>
# Example: ./run.sh check-weather "New York" "USA"

# Set defaults
COMMAND="get-weather"
CITY="London"
COUNTRY="United Kingdom"

# Parse arguments
if [ $# -ge 1 ]; then
  COMMAND=$1
fi

if [ $# -ge 2 ]; then
  CITY=$2
fi

if [ $# -ge 3 ]; then
  COUNTRY=$3
fi

# Normalize command names
case $COMMAND in
  "weather"|"get-weather"|"weather-info")
    COMMAND="get-weather"
    ;;
  "check"|"safety"|"check-deployment"|"check-safety")
    COMMAND="check-deployment-safety"
    ;;
  "deploy"|"deployment"|"deploy-to")
    COMMAND="deploy-to-city"
    ;;
  *)
    if [ "$COMMAND" != "get-weather" ] && [ "$COMMAND" != "check-deployment-safety" ] && [ "$COMMAND" != "deploy-to-city" ]; then
      echo "Unknown command: $COMMAND"
      echo "Usage: ./run.sh <command> <city> <country>"
      echo "Commands:"
      echo "  - get-weather (or weather): Get current weather for a city"
      echo "  - check-deployment-safety (or check): Check if it's safe to deploy"
      echo "  - deploy-to-city (or deploy): Attempt to deploy to a city"
      exit 1
    fi
    ;;
esac

# Execute the command
echo "Executing command: $COMMAND for $CITY, $COUNTRY"

# Check if jq is installed
if command -v jq &> /dev/null; then
  echo "{\"method\":\"$COMMAND\",\"parameters\":{\"city\":\"$CITY\",\"country\":\"$COUNTRY\"}}" | ./floyd-weather-server | jq .
else
  echo "{\"method\":\"$COMMAND\",\"parameters\":{\"city\":\"$CITY\",\"country\":\"$COUNTRY\"}}" | ./floyd-weather-server
fi
