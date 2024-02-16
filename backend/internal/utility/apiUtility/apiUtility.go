package apiUtility

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ValidateURL checks if the url is longer then the given minimum length.
// If the url is not valid, it returns false.
func ValidateURL(w http.ResponseWriter, r *http.Request, maxLength int) bool {
	pathList := strings.Split(r.URL.Path, "/")
	if len(pathList) < maxLength {
		http.Error(w, "URL is too short", http.StatusBadRequest)
		return false
	}
	return true
}

// EncodeData encodes a struct to json and writes it to the response writer. It returns an error if the encoding fails.
func EncodeData(w http.ResponseWriter, data interface{}) error {
	w.Header().Add("content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	return encoder.Encode(data)
}