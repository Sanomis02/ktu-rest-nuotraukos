package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// DeleteImageHandler handles the deletion of an image by ID
func DeleteImageHandler(db *sql.DB, uploadDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract the image ID from the URL
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid image ID", http.StatusBadRequest)
			return
		}

		// Fetch the image filename from the database
		var filename string
		query := "SELECT filename FROM images WHERE id = ?"
		err = db.QueryRow(query, id).Scan(&filename)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Image not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to query database", http.StatusInternalServerError)
			}
			return
		}

		// Delete the image file from the filesystem
		filePath := fmt.Sprintf("%s/%s", uploadDir, filename)
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			http.Error(w, "Failed to delete image file", http.StatusInternalServerError)
			return
		}

		// Delete the image entry from the database
		deleteQuery := "DELETE FROM images WHERE id = ?"
		_, err = db.Exec(deleteQuery, id)
		if err != nil {
			http.Error(w, "Failed to delete image from database", http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Image deleted successfully"))
	}
}

