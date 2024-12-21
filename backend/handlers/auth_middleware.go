package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secure-secret-key")

// AuthMiddleware validates the JWT token
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log request method and path
		log.Printf("Received %s request for %s", r.Method, r.URL.Path)

		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Authorization header missing")
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Check for Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // No "Bearer " prefix
			log.Println("Invalid token format")
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Println("Invalid signing method")
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Printf("Invalid or expired token: %v", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Token is valid; log success
		log.Println("Token validated successfully")

		// Proceed to the next handler
		next(w, r)
	}
}

