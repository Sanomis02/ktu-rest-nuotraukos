#!/bin/bash

# Function to display usage
usage() {
  echo "Usage: $0 -t <JWT_token> -i <image_id>"
  exit 1
}

# Parse command-line arguments
while getopts "t:i:" opt; do
  case $opt in
    t) token="$OPTARG" ;;
    i) image_id="$OPTARG" ;;
    *) usage ;;
  esac
done

# Ensure both token and image ID are provided
if [ -z "$token" ] || [ -z "$image_id" ]; then
  usage
fi

# API endpoint
API_URL="https://localhost/api/image/$image_id"

# Perform the DELETE request
response=$(curl -sk -X DELETE "$API_URL" \
  -H "Authorization: Bearer $token")

# Check the response
if [[ "$response" == "Image deleted successfully" ]]; then
  echo "Image deleted successfully."
else
  echo "Error: Failed to delete image. Response: $response"
  exit 1
fi

