package apiUtility

import (
	"adsb-api/internal/global/errorMsg"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"path"
	"strings"
)

// EncodeJsonData encodes a struct to json and writes it to the response writer. Returns an error if the encoding fails.
func EncodeJsonData(w http.ResponseWriter, data interface{}) {
	w.Header().Add("content-type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(data)
	if err != nil {
		http.Error(w, errorMsg.ErrorEncodingJsonData, http.StatusInternalServerError)
		log.Error().Msgf(errorMsg.ErrorEncodingJsonData+": %q", err)
	}
}

// ValidateURL checks the validity of an HTTP request URL.
// Cleans the URL.
// Verifies the URL length, and the presence of specified parameters
// if endpoint does not use parameters leaves params nil.
// If the URL length exceeds the maximum length or if any of the specified parameters are missing from the request,
// it returns false.
func ValidateURL(r *http.Request, maxLength int, optionalParams []string) error {
	url := path.Clean(r.URL.Path)
	if len(strings.SplitAfter(url, "/")) > maxLength {
		return errors.New(errorMsg.ErrorTongURL)
	}

	if r.URL.Query().Encode() == "" {
		return nil
	}

	query := r.URL.Query()

	if len(query) != len(optionalParams) {
		return errors.New(errorMsg.ErrorInvalidQueryParams + strings.Join(optionalParams, ", "))
	}

	for _, param := range optionalParams {
		values, ok := query[param]
		if !ok || len(values) == 0 || values[0] == "" {
			return errors.New(errorMsg.ErrorInvalidQueryParams + strings.Join(optionalParams, ", "))
		}
	}

	return nil
}
