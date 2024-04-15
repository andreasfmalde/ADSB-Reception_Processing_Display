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
// 1. Cleans the URL and verifies the URL length
// 2. It checks parameters in the url against parameter optionalParams,
// if endpoint does not use parameters leaves params nil.
//
// If the URL length exceeds the maximum length or if any of the specified parameters are missing from the request, it
// writes to the ResponseWriter with appropriate status codes and returns an error.
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

// NoContent sets the Access-Control-Allow-Origin header to "*"
// and writes a StatusNoContent header to the response writer.
func NoContent(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusNoContent)
}
