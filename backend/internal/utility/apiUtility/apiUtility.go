package apiUtility

import (
	"net/http"
	"strings"
)

func ValidateURL(w http.ResponseWriter, r *http.Request, maxLength int) bool {
	pathList := strings.Split(r.URL.Path, "/")
	if len(pathList) < maxLength {
		http.Error(w, "URL is too short", http.StatusBadRequest)
		return false
	}
	return true
}
