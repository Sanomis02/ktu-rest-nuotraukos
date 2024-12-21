package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
)

// Credentials is a struct to hold JSON username/password in login request.
type Credentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var creds Credentials
        if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        // For simplicity, let's assume you have a users table with columns 'username' and 'password'
        // Adjust the SQL query and table name to your needs, or handle password hashing, etc.
        // Example only. Definitely secure your logic better in production!
        row := db.QueryRow("SELECT username FROM users WHERE username = ? AND password = ?", creds.Username, creds.Password)
        var dbUsername string
        err := row.Scan(&dbUsername)
        if err != nil {
            // If there is no matching row or any other error
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }

        // Generate JWT token
        token, err := GenerateToken(dbUsername)
        if err != nil {
            http.Error(w, "Could not generate token", http.StatusInternalServerError)
            return
        }

        // Return the token to client
        resp := map[string]string{"token": token}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    }
}
