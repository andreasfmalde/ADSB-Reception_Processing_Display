package aircraftCurrentHandler

import (
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/global/geoJSON"
	"adsb-api/internal/logger"
	"adsb-api/internal/service"
	"adsb-api/internal/utility/apiUtility"
	"fmt"
	"net/http"
)

// CurrentAircraftHandler handles HTTP requests for /aircraft/current/ endpoint.
func CurrentAircraftHandler(svc service.RestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleCurrentAircraftGetRequest(w, r, svc)
		default:
			http.Error(w, fmt.Sprintf(errorMsg.MethodNotSupported, r.Method), http.StatusMethodNotAllowed)
		}
	}
}

// handleCurrentAircraftGetRequest handles GET requests for the /aircraft/current/ endpoint.
// Sends all current aircraft in the database to the client.
func handleCurrentAircraftGetRequest(w http.ResponseWriter, r *http.Request, svc service.RestService) {
	res, err := svc.GetCurrentAircraft()
	if err != nil {
		http.Error(w, errorMsg.ErrorRetrievingCurrentAircraft, http.StatusInternalServerError)
		logger.Error.Printf(errorMsg.ErrorRetrievingCurrentAircraft+": %q Path: %q", err, r.URL)
		return
	}
	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	aircraft, err := geoJSON.ConvertCurrentModelToGeoJson(res)
	if err != nil {
		http.Error(w, errorMsg.ErrorConvertingDataToGeoJson, http.StatusInternalServerError)
		logger.Error.Printf(errorMsg.ErrorConvertingDataToGeoJson+": %q", err)
		return
	}

	apiUtility.EncodeJsonData(w, aircraft)
}
