package currentAircraftHandler

import (
	"adsb-api/internal/utility/apiUtility"
	"net/http"
)

// CurrentAircraftHandler handles HTTP requests for /aircraft/current/ endpoint.
func CurrentAircraftHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if apiUtility.ValidateURL(w, r, 6) {
			handleCurrentAircraftGetRequest(w, r)
		}
	default:
		http.Error(w, "Method "+r.Method+" is not supported", http.StatusMethodNotAllowed)
	}
}

func handleCurrentAircraftGetRequest(w http.ResponseWriter, r *http.Request) {

}
