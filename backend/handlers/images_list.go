package handlers

import (
	"encoding/json"
	"net/http"
	"os"
)

// Image represents an image file with its name and URL
type Image struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// ListImagesHandler returns a list of all images in the uploads directory
func ListImagesHandler(uploadDir, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var images []Image

		// Read all files in the uploads directory
		files, err := os.ReadDir(uploadDir)
		if err != nil {
			http.Error(w, "Unable to read uploads directory", http.StatusInternalServerError)
			return
		}

		// Generate URLs for each file
		for _, file := range files {
			if !file.IsDir() {
				fileName := file.Name()
				imageURL := baseURL + "/uploads/" + fileName
				images = append(images, Image{Name: fileName, URL: imageURL})
			}
		}

		// Respond with JSON list of images
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(images)
	}
}

