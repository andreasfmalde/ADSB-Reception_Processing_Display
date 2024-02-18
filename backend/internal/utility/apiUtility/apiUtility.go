package apiUtility

import (
	"adsb-api/internal/logger"
	"encoding/json"
	"net/http"
)

// EncodeData encodes a struct to json and writes it to the response writer. It returns an error if the encoding fails.
func EncodeData(w http.ResponseWriter, data interface{}) {
	w.Header().Add("content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(data)
	if err != nil {
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		logger.Error.Println("Failed to encode data: %q", err)
	}
}
