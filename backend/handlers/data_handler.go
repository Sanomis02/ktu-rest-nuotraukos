package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// Data represents a row in the database table
type Data struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// DataHandler returns data from the sample_table
func DataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM sample_table")
		if err != nil {
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var data []Data
		for rows.Next() {
			var d Data
			if err := rows.Scan(&d.ID, &d.Name); err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				return
			}
			data = append(data, d)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
