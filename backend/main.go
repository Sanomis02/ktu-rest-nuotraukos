package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/Sanomis02/ktu-rest-nuotraukos/handlers"
)

func main() {
	// Database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	// Data Source Name (DSN) for MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	uploadDir := "./uploads"
	baseURL := "http://localhost:8080"

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	http.HandleFunc("/api/login", handlers.LoginHandler)                              // Public
	http.HandleFunc("/api/data", handlers.AuthMiddleware(handlers.DataHandler(db))) // Protected
	http.HandleFunc("/api/upload", handlers.AuthMiddleware(handlers.UploadImageHandler(db, uploadDir))) // Protected
	http.HandleFunc("/api/uploads", handlers.ListImagesHandler(uploadDir, baseURL))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir)))) // Serve images statically

	log.Println("Starting server on :8000")
        log.Fatal(http.ListenAndServe(":8000", nil))
}

