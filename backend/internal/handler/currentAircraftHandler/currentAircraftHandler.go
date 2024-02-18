package currentAircraftHandler

import (
	"adsb-api/internal/db"
	"adsb-api/internal/logger"
	"adsb-api/internal/utility/apiUtility"
	"database/sql"
	"net/http"
)

// CurrentAircraftHandler handles HTTP requests for /aircraft/current/ endpoint.
func CurrentAircraftHandler(dbConn *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if apiUtility.ValidateURL(w, r, 6) {
				handleCurrentAircraftGetRequest(w, r, dbConn)
			}
		default:
			http.Error(w, "Method "+r.Method+" is not supported", http.StatusMethodNotAllowed)
		}
	}
}

// handleCurrentAircraftGetRequest handles the GET request for the /aircraft/current/ endpoint.
// Sends all current aircraft in the database to the client.
func handleCurrentAircraftGetRequest(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	aircraft, err := db.RetrieveCurrentTimeAircrafts(dbConn)
	if err != nil {
		http.Error(w, "Error during request execution", http.StatusInternalServerError)
		logger.Error.Printf("Error: %q Path: %q", err, r.URL)
		return
	}
	if len(aircraft) == 0 {
		http.Error(w, "No aircraft found.", http.StatusNotFound)
		return
	}
	apiUtility.EncodeData(w, aircraft)
}
