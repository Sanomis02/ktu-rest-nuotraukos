package handlers

import (
    "context"
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// Replace with a secure random string from environment variable in production.
var jwtSecret = []byte("MY_SUPER_SECRET_KEY")

// GenerateToken generates a JWT token for a specific user (by username or ID).
func GenerateToken(username string) (string, error) {
    // Set token claims
    claims := jwt.MapClaims{
        "username": username,
        "exp":      time.Now().Add(time.Minute * 5).Unix(), // token expires in 5 minutes
    }

    // Create token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Sign token
    signedToken, err := token.SignedString(jwtSecret)
    if err != nil {
        return "", err
    }
    return signedToken, nil
}

// ValidateToken parses and validates the token string from the Authorization header.
func ValidateToken(tokenString string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Make sure that the token method conform to "SigningMethodHMAC"
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtSecret, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
    return token, nil
}

// AuthenticationMiddleware checks for the Bearer token in the Authorization header.
// If valid, it calls the next handler; otherwise, it returns 401 Unauthorized.
func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
            return
        }

        // Typical Authorization header format is: "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
            return
        }

        tokenString := parts[1]
        token, err := ValidateToken(tokenString)
        if err != nil {
            http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
            return
        }

        // Optionally, you can add claims to the request context
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
           ctx := context.WithValue(r.Context(), "username", claims["username"])
            r = r.WithContext(ctx)
	}

        // Token is valid; proceed with the next handler
        next.ServeHTTP(w, r)
    }
}

