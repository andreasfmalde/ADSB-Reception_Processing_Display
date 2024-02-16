package defaultHandler

import (
	"adsb-api/internal/global"
	"adsb-api/internal/utility/apiUtility"
	"net/http"
)

type DefaultStruct struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	MadeBy    []string          `json:"madeby"`
	Endpoints map[string]string `json:"endpoints"`
}

// DefaultHandler function, which prints info about the service
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		endpoints := make(map[string]string)
		endpoints["current_aircraft"] = global.CurrentAircraftPath

		madeBy := []string{"Andreas Follevaag Malde", "Fredrik Sundt-Hansen"}

		out := DefaultStruct{
			Name:      "ADS-B Reception, Processing, Displaying and Analysis",
			Version:   global.VERSION,
			MadeBy:    madeBy,
			Endpoints: endpoints,
		}

		apiUtility.EncodeData(w, out)
	default:
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
	}
}
