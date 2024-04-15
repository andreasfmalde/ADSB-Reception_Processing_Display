package apiUtility

import (
	"adsb-api/internal/global/errorMsg"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
)

// EncodeJsonData encodes a struct to json and writes it to the response writer. Returns an error if the encoding fails.
func EncodeJsonData(w http.ResponseWriter, data interface{}) error {
	w.Header().Add("content-type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	return encoder.Encode(data)
}

// ValidateURL checks the validity of an HTTP request URL.
// Cleans the URL.
// Verifies the URL length, and the presence of specified parameters
// if endpoint does not use parameters leaves params nil.
// If the URL length exceeds the maximum length or if any of the specified parameters are missing from the request,
// it returns false.
func ValidateURL(w http.ResponseWriter, r *http.Request, maxLength int, optionalParams []string) error {
	url := strings.Split(path.Clean(r.URL.Path), "/")
	if len(url) > maxLength {
		http.Error(w, errorMsg.ErrorTongURL, http.StatusRequestURITooLong)
		return fmt.Errorf("falied to validate URL")
	}

	if r.URL.Query().Encode() == "" {
		return nil
	}

	query := r.URL.Query()

	if len(query) != len(optionalParams) {
		http.Error(w, fmt.Errorf(errorMsg.ErrorInvalidQueryParams+": %s", strings.Join(optionalParams, ", ")).Error(), http.StatusBadRequest)
		return fmt.Errorf("falied to validate URL")
	}

	for _, param := range optionalParams {
		values, ok := query[param]
		if !ok || len(values) == 0 || values[0] == "" {
			http.Error(w, fmt.Errorf(errorMsg.ErrorInvalidQueryParams+": %s", strings.Join(optionalParams, ", ")).Error(), http.StatusBadRequest)
			return fmt.Errorf("falied to validate URL")
		}
	}

	return nil
}

func NoContent(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusNoContent)
}
