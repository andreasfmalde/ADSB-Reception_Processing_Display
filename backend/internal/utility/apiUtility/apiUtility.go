package apiUtility

import (
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/logger"
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"strings"
)

// EncodeJsonData encodes a struct to json and writes it to the response writer. Returns an error if the encoding fails.
func EncodeJsonData(w http.ResponseWriter, data interface{}) {
	w.Header().Add("content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(data)
	if err != nil {
		http.Error(w, errorMsg.ErrorEncodingJsonData, http.StatusInternalServerError)
		logger.Error.Printf(errorMsg.ErrorEncodingJsonData+": %q", err)
	}
}

// ValidateURL checks the validity of an HTTP request URL.
// It verifies the URL length, and the presence of specified parameters
// if endpoint does not use parameters leaves params nil.
// If the URL length exceeds the maximum length or if any of the specified parameters are missing from the request,
// it returns false.
func ValidateURL(url string, query map[string][]string, maxLength int, params []string) error {
	url = path.Clean(url)
	if len(url) > maxLength {
		return errors.New(errorMsg.ErrorTongURL)
	}

	if params == nil {
		return nil
	}

	if len(query) != len(params) {
		return errors.New(errorMsg.ErrorInvalidQueryParams + strings.Join(params, ", "))
	}

	for _, param := range params {
		values, ok := query[param]
		if !ok || len(values) == 0 || values[0] == "" {
			return errors.New(errorMsg.ErrorInvalidQueryParams + strings.Join(params, ", "))
		}
	}

	return nil
}
