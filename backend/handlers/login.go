package handlers

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// Define the request structure
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// Define the response structure
type LoginResponse struct {
    Message string `json:"message"`
    Token   string `json:"token,omitempty"`
}

var jwtSecret = []byte("your-secure-secret-key")

// GenerateJWT generates a JWT token for the given username
func GenerateJWT(username string) (string, error) {
    claims := jwt.MapClaims{
        "username": username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
        "iat":      time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// LoginHandler handles the login requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var loginReq LoginRequest
    err := json.NewDecoder(r.Body).Decode(&loginReq)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Dummy authentication logic (replace with DB lookup)
    if loginReq.Username != "testuser" || loginReq.Password != "password123" {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(LoginResponse{Message: "Invalid username or password"})
        return
    }

    token, err := GenerateJWT(loginReq.Username)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(LoginResponse{
        Message: "Login successful",
        Token:   token,
    })
}

