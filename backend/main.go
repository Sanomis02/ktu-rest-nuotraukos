package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Data struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
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
	})

	log.Println("Starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
