#!/bin/bash

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
			API_URL="https://localhost/api/user"

			# Perform the POST request to register a new user
			response=$(curl -sk -X POST "$API_URL" \
				  -H "Content-Type: application/json" \
				    -d "{\"username\":\"$username\",\"password\":\"$password\"}")

			# Check if the response is successful
			if echo "$response" | jq -e '.message' > /dev/null 2>&1; then
				  echo "Registration successful: $(echo "$response" | jq -r '.message')"
			  else
				    echo "Error during registration: $response"
				      exit 1
			fi

