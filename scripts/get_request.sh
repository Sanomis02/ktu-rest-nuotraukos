#!/bin/bash

# Function to display usage
usage() {
	  echo "Usage: $0 -e <endpoint> [-t <token>]"
	    echo "Endpoints:"
	      echo "  uploads - Fetch data from /api/uploads"
	        echo "  users   - Fetch data from /api/users"
		  exit 1
	  }

  # Parse command-line arguments
  while getopts "e:t:" opt; do
	    case $opt in
		        e) endpoint="$OPTARG" ;;
			    t) token="$OPTARG" ;;
			        *) usage ;;
				  esac
			  done

			  # Ensure the endpoint is provided
			  if [ -z "$endpoint" ]; then
				    usage
			  fi

			  # Map endpoint argument to the actual API path
			  if [ "$endpoint" == "uploads" ]; then
				    API_URL="https://localhost/api/uploads"
			    elif [ "$endpoint" == "users" ]; then
				      API_URL="https://localhost/api/users"
			      else
				        echo "Error: Invalid endpoint. Valid options are 'uploads' or 'users'."
					  usage
			  fi

			  # Prepare curl headers
			  headers=()
			  headers+=("-H" "Content-Type: application/json")

			  # Add Authorization header only if token is provided
			  if [ -n "$token" ]; then
				    headers+=("-H" "Authorization: Bearer $token")
			  fi

			  # Perform the GET request
			  response=$(curl -sk -X GET "${headers[@]}" "$API_URL")

			  # Check if the response is successful
			  if echo "$response" | jq -e '.' > /dev/null 2>&1; then
				    echo "Response from $endpoint:"
				      echo "$response" | jq
			      else
				        echo "Error: Failed to fetch data. Response: $response"
					  exit 1
			  fi

