#!/bin/bash

# Function to display usage
usage() {
  echo "Usage: $0 -u <username> -p <password>"
  exit 1
}

# Parse command-line arguments
while getopts "u:p:" opt; do
  case $opt in
    u) username="$OPTARG" ;;
    p) password="$OPTARG" ;;
    *) usage ;;
  esac
done

# Ensure both username and password are provided
if [ -z "$username" ] || [ -z "$password" ]; then
  usage
fi

# API endpoint
API_URL="https://localhost/api/login"

# Perform the POST request to authenticate the user
response=$(curl -sk -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$username\",\"password\":\"$password\"}")

# Check if the response is valid JSON
if echo "$response" | jq -e . > /dev/null 2>&1; then
  # Extract the token from the JSON response
  token=$(echo "$response" | jq -r '.token')

  # Check if the token was successfully retrieved
  if [ -n "$token" ] && [ "$token" != "null" ]; then
    echo "Login successful. Token: $token"
    exit 0
  else
    echo "Error: Login failed. Response: $response"
    exit 1
  fi
else
  # Handle non-JSON responses
  echo "Error: Login failed. Response: $response"
  exit 1
fi

