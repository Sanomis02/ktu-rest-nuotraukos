package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"

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

	r := mux.NewRouter()

	// Public endpoints
	r.HandleFunc("/api/login", handlers.LoginHandler(db)).Methods(http.MethodPost)
	r.HandleFunc("/api/uploads", handlers.ListImagesHandler(db, uploadDir, baseURL)).Methods(http.MethodGet)

	// Protected endpoints
	r.Handle("/api/upload", handlers.AuthenticationMiddleware(http.HandlerFunc(handlers.UploadImageHandler(db, uploadDir)))).Methods(http.MethodPost)
	r.Handle("/api/image/{id}", handlers.AuthenticationMiddleware(http.HandlerFunc(handlers.DeleteImageHandler(db, uploadDir)))).Methods(http.MethodDelete)
	r.Handle("/api/users", handlers.AuthenticationMiddleware(http.HandlerFunc(handlers.UsersHandler(db)))).Methods(http.MethodGet)
	r.Handle("/api/user", http.HandlerFunc(handlers.CreateUserHandler(db))).Methods(http.MethodPost)

	// Static file serving for uploaded images
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir)))).Methods(http.MethodGet)

	log.Println("Starting server on :8000")
        log.Fatal(http.ListenAndServe(":8000", r))
}

