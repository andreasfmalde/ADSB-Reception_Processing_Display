package apiUtility

import (
	"adsb-api/internal/logger"
	"encoding/json"
	"net/http"
)

// EncodeJsonData encodes a struct to json and writes it to the response writer. Returns an error if the encoding fails.
func EncodeJsonData(w http.ResponseWriter, data interface{}) {
	w.Header().Add("content-type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(data)
	if err != nil {
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		logger.Error.Printf("Failed to encode data: %q", err)
	}
}
