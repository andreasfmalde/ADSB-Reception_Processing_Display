package aircraftHistory

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/logger"
	"adsb-api/internal/utility/apiUtility"
	"fmt"
	"net/http"
)

var params = []string{"icao"}

// HistoryAircraftHandler handles HTTP requests for /aircraft/history endpoint.
func HistoryAircraftHandler(db db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := apiUtility.ValidateURL(r.URL.Path, r.URL.Query(), len(global.AircraftHistoryPath), params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			handleHistoryAircraftGetRequest(w, r, db)
		default:
			http.Error(w, fmt.Sprintf(global.MethodNotSupported, r.Method), http.StatusMethodNotAllowed)
		}
	}
}

// handleHistoryAircraftGetRequest handles GET requests for the /aircraft/history endpoint.
// Sends history data for aircraft given by the icao query parameter.
// A valid icao: "ABC123"
func handleHistoryAircraftGetRequest(w http.ResponseWriter, r *http.Request, db db.Database) {
	var search = r.URL.Query().Get("icao")
	res, err := db.GetHistoryByIcao(search)
	if err != nil {
		http.Error(w, global.ErrorRetrievingAircraftWithIcao+search, http.StatusInternalServerError)
		logger.Error.Printf(global.ErrorRetrievingAircraftWithIcao+search+": %q URL: %q", err, r.URL)
		return
	}
	if len(res.Features) == 0 {
		http.Error(w, global.NoAircraftFound, http.StatusNoContent)
		return
	}
	apiUtility.EncodeJsonData(w, res)
}
