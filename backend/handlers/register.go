package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strings"

    "golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type CreateUserResponse struct {
    Message string `json:"message"`
}

// HashPassword hashes a plain-text password using bcrypt
func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hashedPassword), err
}

func CreateUserHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var req CreateUserRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        // Validate input
        if req.Username == "" || req.Password == "" {
            http.Error(w, "Username and password are required", http.StatusBadRequest)
            return
        }

        // Hash the password
        hashedPassword, err := HashPassword(req.Password)
        if err != nil {
            http.Error(w, "Failed to hash password", http.StatusInternalServerError)
            return
        }

        // Insert the new user into the database
        query := "INSERT INTO users (username, password) VALUES (?, ?)"
        _, err = db.Exec(query, req.Username, hashedPassword)
	if err != nil {
            // Check if the error is a duplicate entry error
            if strings.Contains(err.Error(), "Duplicate entry") {
                http.Error(w, "Username already exists", http.StatusConflict)
                return
	    }
	}
        // Respond with success message
        response := CreateUserResponse{Message: "User created successfully"}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}

