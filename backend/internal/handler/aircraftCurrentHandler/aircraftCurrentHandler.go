package aircraftCurrentHandler

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/global/geoJSON"
	"adsb-api/internal/logger"
	"adsb-api/internal/utility/apiUtility"
	"fmt"
	"net/http"
)

/*
CurrentAircraftHandler handles HTTP requests for /aircraft/current/ endpoint.
Endpoints:
	- GET /aircraft/current/

Planned options:
	- ?icao=
	- ?callsign=
	- altitude=
		- ?minAlt=
		- ?maxAlt=
	- position=
		- circle
			- ?center= (lat,long)
			- ?radius=
		- polygon
			- ?points= (lat,long),(lat,long),...
*/

// CurrentAircraftHandler handles HTTP requests for /aircraft/current/ endpoint.
func CurrentAircraftHandler(db db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleCurrentAircraftGetRequest(w, r, db)
		default:
			http.Error(w, fmt.Sprintf(global.MethodNotSupported, r.Method), http.StatusMethodNotAllowed)
		}
	}
}

// handleCurrentAircraftGetRequest handles GET requests for the /aircraft/current/ endpoint.
// Sends all current aircraft in the database to the client.
func handleCurrentAircraftGetRequest(w http.ResponseWriter, r *http.Request, db db.Database) {
	res, err := db.GetAllCurrentAircraft()
	if err != nil {
		http.Error(w, global.ErrorRetrievingCurrentAircraft, http.StatusInternalServerError)
		logger.Error.Printf(global.ErrorRetrievingCurrentAircraft+": %q Path: %q", err, r.URL)
		return
	}
	if len(res) == 0 {
		http.Error(w, global.NoAircraftFound, http.StatusNoContent)
		return
	}

	aircraft, err := geoJSON.ConvertCurrentModelToGeoJson(res)
	if err != nil {
		http.Error(w, global.ErrorConvertingDataToGeoJson, http.StatusInternalServerError)
		logger.Error.Printf(global.ErrorConvertingDataToGeoJson+": %q", err)
		return
	}

	apiUtility.EncodeJsonData(w, aircraft)
}
