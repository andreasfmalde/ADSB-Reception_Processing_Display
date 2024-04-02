package aircraftHistory

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/service/restService"
	"adsb-api/internal/utility/apiUtility"
	"adsb-api/internal/utility/convert"
	"adsb-api/internal/utility/logger"
	"fmt"
	"net/http"
)

var params = []string{"icao"}

// HistoryAircraftHandler handles HTTP requests for /aircraft/history endpoint.
func HistoryAircraftHandler(svc restService.RestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := apiUtility.ValidateURL(r.URL.Path, r.URL.Query(), len(global.AircraftHistoryPath), params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			handleHistoryAircraftGetRequest(w, r, svc)
		default:
			http.Error(w, fmt.Sprintf(errorMsg.MethodNotSupported, r.Method), http.StatusMethodNotAllowed)
		}
	}
}

// handleHistoryAircraftGetRequest handles GET requests for the /aircraft/history endpoint.
// Sends history data for aircraft given by the icao query parameter.
// A valid icao: "ABC123"
func handleHistoryAircraftGetRequest(w http.ResponseWriter, r *http.Request, svc restService.RestService) {
	var search = r.URL.Query().Get("icao")
	res, err := svc.GetAircraftHistoryByIcao(search)
	if err != nil {
		http.Error(w, errorMsg.ErrorRetrievingAircraftWithIcao+search, http.StatusInternalServerError)
		logger.Error.Printf(errorMsg.ErrorRetrievingAircraftWithIcao+search+": %q URL: %q", err, r.URL)
		return
	}
	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	aircraft, err := convert.HistoryModelToGeoJson(res)
	if err != nil {
		http.Error(w, errorMsg.ErrorConvertingDataToGeoJson, http.StatusInternalServerError)
		logger.Error.Printf(errorMsg.ErrorConvertingDataToGeoJson+": %q", err)
		return
	}

	apiUtility.EncodeJsonData(w, aircraft)
}
