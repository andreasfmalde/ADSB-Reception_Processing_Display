package aircraftCurrentHandler

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/service/restService"
	"adsb-api/internal/utility/apiUtility"
	"adsb-api/internal/utility/convert"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

// CurrentAircraftHandler handles HTTP requests for /aircraft/current/ endpoint.
func CurrentAircraftHandler(svc restService.RestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := apiUtility.ValidateURL(r, global.CurrentPathMaxLength, []string{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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
func handleCurrentAircraftGetRequest(w http.ResponseWriter, r *http.Request, svc restService.RestService) {
	res, err := svc.GetCurrentAircraft()
	if err != nil {
		http.Error(w, errorMsg.ErrorRetrievingCurrentAircraft, http.StatusInternalServerError)
		log.Error().Msgf(errorMsg.ErrorRetrievingCurrentAircraft+": %q Path: %q", err, r.URL)
		return
	}
	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	aircraft, err := convert.CurrentModelToGeoJson(res)
	if err != nil {
		http.Error(w, errorMsg.ErrorConvertingDataToGeoJson, http.StatusInternalServerError)
		log.Error().Msgf(errorMsg.ErrorConvertingDataToGeoJson+": %q", err)
		return
	}

	err = apiUtility.EncodeJsonData(w, aircraft)
	if err != nil {
		http.Error(w, errorMsg.ErrorEncodingJsonData, http.StatusInternalServerError)
		log.Error().Msgf(errorMsg.ErrorEncodingJsonData+": %q", err)
	}
}
