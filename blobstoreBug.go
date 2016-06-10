package blobstoreBug

import (
	"encoding/json"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine/image"
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
	c := appengine.NewContext(r)
	blobs, _, err := blobstore.ParseUpload(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := []map[string]string{}
	for _, blob := range blobs["image"] {
		blobStatus := map[string]string{
			"BlobKey": string(blob.BlobKey),
		}
		url, err := image.ServingURL(c, blob.BlobKey, &image.ServingURLOptions{Size: 100})
		if err != nil {
			blobStatus["status"] = err.Error()
			continue
		}
		blobStatus["status"] = "success"
		blobStatus["imageURL"] = url.String()
		resp = append(resp, blobStatus)
	}
	json.NewEncoder(w).Encode(resp)
}
