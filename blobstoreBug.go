package blobstoreBug

import (
	"encoding/json"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
)

func init() {
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/uploadURL", createUploadURL)
}

func createUploadURL(w http.ResponseWriter, r *http.Request) {
	uploadURL, err := blobstore.UploadURL(appengine.NewContext(r), "/upload", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"uploadURL": uploadURL.String(),
	})
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	blobs, _, err := blobstore.ParseUpload(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(blobs)
}
