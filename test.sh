#!/bin/bash

# Test script for Floyd Weather MCP Server
# This script simulates MCP protocol interactions

echo "Testing Floyd Weather MCP Server..."

COMMAND="/Users/devopsjester/repos/experiment/mcp/floyd-weather-deployer/floyd-weather-server"

# Test weather information
echo "Test 1: Weather information for London, UK"
echo '{"method":"get-weather","parameters":{"city":"London","country":"United Kingdom"}}' | $COMMAND
echo ""

# Test deployment safety check
echo "Test 2: Deployment safety check for San Francisco, USA"
echo '{"method":"check-deployment-safety","parameters":{"city":"San Francisco","country":"USA"}}' | $COMMAND
echo ""

# Test deployment request
echo "Test 3: Attempt to deploy to Tokyo, Japan"
echo '{"method":"deploy-to-city","parameters":{"city":"Tokyo","country":"Japan"}}' | $COMMAND
echo ""

echo "Testing complete."
