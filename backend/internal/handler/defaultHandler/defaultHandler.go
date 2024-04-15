package defaultHandler

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/utility/apiUtility"
	"github.com/rs/zerolog/log"
	"net/http"
)

type DefaultStruct struct {
	Name      string   `json:"name"`
	Version   string   `json:"version"`
	MadeBy    []string `json:"madeby"`
	Endpoints []string `json:"endpoints"`
}

// DefaultHandler function, which prints info about the service
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var endpoints []string
		endpoints = append(endpoints, global.AircraftCurrentPath)
		endpoints = append(endpoints, global.AircraftHistoryPath)

		madeBy := []string{"Andreas Follevaag Malde", "Fredrik Sundt-Hansen"}

		out := DefaultStruct{
			Name:      "ADS-B Reception, Processing, Displaying and Analysis",
			Version:   global.VERSION,
			MadeBy:    madeBy,
			Endpoints: endpoints,
		}

		err := apiUtility.EncodeJsonData(w, out)
		if err != nil {
			http.Error(w, errorMsg.ErrorEncodingJsonData, http.StatusInternalServerError)
			log.Error().Msgf(errorMsg.ErrorEncodingJsonData+": %q", err)
		}
	default:
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
	}
}
