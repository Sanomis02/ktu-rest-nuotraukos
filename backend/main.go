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
	baseURL := "https://localhost"

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Public endpoint: /api/login
	http.HandleFunc("/api/login", handlers.LoginHandler(db))

	// Protected endpoints:
	// We wrap them with handlers.AuthenticationMiddleware()
	http.HandleFunc("/api/data", handlers.AuthenticationMiddleware(handlers.DataHandler(db)),)
	http.HandleFunc("/api/upload", handlers.AuthenticationMiddleware(handlers.UploadImageHandler(db, uploadDir)),)
	http.HandleFunc("/api/uploads", handlers.ListImagesHandler(uploadDir, baseURL))
	http.HandleFunc("/api/users", handlers.AuthenticationMiddleware(handlers.UsersHandler(db)),)
	http.HandleFunc("/api/user", handlers.CreateUserHandler(db))

	// Serves the static files in "uploads" directory
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))

	log.Println("Starting server on :8000")
        log.Fatal(http.ListenAndServe(":8000", nil))
}

