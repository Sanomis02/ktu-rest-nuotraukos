package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// UploadImageHandler handles image uploads and saves metadata to the database
func UploadImageHandler(db *sql.DB, uploadDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse multipart form data
		err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
		if err != nil {
			http.Error(w, "Unable to parse form data", http.StatusBadRequest)
			return
		}

		// Get the uploaded file
		file, header, err := r.FormFile("image") // "image" is the form field name
		if err != nil {
			http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Generate a unique filename with a timestamp
		uniqueFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(header.Filename))
		filePath := filepath.Join(uploadDir, uniqueFilename)

		// Save the file to the file system
		outFile, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Unable to save file to the server", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		// Save metadata to the database
		query := "INSERT INTO images (filename, filepath) VALUES (?, ?)"
		_, err = db.Exec(query, uniqueFilename, filePath)
		if err != nil {
			http.Error(w, "Failed to save metadata to database", http.StatusInternalServerError)
			return
		}

		// Respond to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"message": "File uploaded successfully", "filename": "%s"}`, uniqueFilename)
	}
}

