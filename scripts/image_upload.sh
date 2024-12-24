#!/bin/bash

# Function to display usage
usage() {
  echo "Usage: $0 -t <JWT_token> -f <image_path> [-c <comment>]"
  exit 1
}

# Parse command-line arguments
while getopts "t:f:c:" opt; do
  case $opt in
    t) token="$OPTARG" ;;
    f) image_path="$OPTARG" ;;
    c) comment="$OPTARG" ;;
    *) usage ;;
  esac
done

# Ensure both token and image path are provided
if [ -z "$token" ] || [ -z "$image_path" ]; then
  usage
fi

# Check if the file exists
if [ ! -f "$image_path" ]; then
  echo "Error: File '$image_path' does not exist."
  exit 1
fi

# API endpoint
API_URL="https://localhost/api/upload"

# Form the POST request
if [ -n "$comment" ]; then
  # Include comment in the request if provided
  response=$(curl -sk -X POST "$API_URL" \
    -H "Authorization: Bearer $token" \
    -H "Content-Type: multipart/form-data" \
    -F "image=@$image_path" \
    -F "comment=$comment")
else
  # Send the request without comment
  response=$(curl -sk -X POST "$API_URL" \
    -H "Authorization: Bearer $token" \
    -H "Content-Type: multipart/form-data" \
    -F "image=@$image_path")
fi

# Check the response
if echo "$response" | jq -e '.message' > /dev/null 2>&1; then
  echo "Image uploaded successfully: $(echo "$response" | jq -r '.message')"
else
  echo "Error: Failed to upload image. Response: $response"
  exit 1
fi

