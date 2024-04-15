package aircraftHistory

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/global/models"
	"adsb-api/internal/service/restService"
	"adsb-api/internal/utility/apiUtility"
	"adsb-api/internal/utility/convert"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"path"
	"strconv"
)

var optionalParams = []string{"hour"}

// HistoryAircraftHandler handles HTTP requests for /aircraft/history/{icao}?hour= endpoint.
func HistoryAircraftHandler(svc restService.RestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := apiUtility.ValidateURL(w, r, global.HistoryPathMaxLength, optionalParams)
		if err != nil {
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

// handleHistoryAircraftGetRequest handles GET requests for the aircraft/history/{icao}?hour= endpoint.
// Sends history data for aircraft given by the icao query parameter.
func handleHistoryAircraftGetRequest(w http.ResponseWriter, r *http.Request, svc restService.RestService) {
	search := path.Base(r.URL.Path)
	if search == "history" {
		http.Error(w, errorMsg.EmptyIcao, http.StatusBadRequest)
		return
	} else if len(search) > 6 {
		http.Error(w, errorMsg.TooLongIcao, http.StatusBadRequest)
		return
	}

	var err error
	var res []models.AircraftHistoryModel

	if r.URL.Query().Has("hour") {
		hour, err := strconv.Atoi(r.URL.Query().Get("hour"))
		if err != nil {
			http.Error(w, errorMsg.InvalidQueryParameterHour, http.StatusBadRequest)
			log.Error().Msgf(errorMsg.InvalidQueryParameterHour+" Error : %q", err)
			return
		}
		res, err = svc.GetAircraftHistoryByIcaoFilterByTimestamp(search, hour)
	} else {
		res, err = svc.GetAircraftHistoryByIcao(search)
	}

	if err != nil {
		http.Error(w, errorMsg.ErrorRetrievingAircraftWithIcao+search, http.StatusInternalServerError)
		log.Error().Msgf(errorMsg.ErrorRetrievingAircraftWithIcao+": %s Error : %q URL: %q", search, err, r.URL)
		return
	}

	if len(res) == 0 || len(res) < 2 {
		apiUtility.NoContent(w)
		return
	}

	aircraft, err := convert.HistoryModelToGeoJson(res)
	if err != nil {
		http.Error(w, errorMsg.ErrorConvertingDataToGeoJson, http.StatusInternalServerError)
		log.Error().Msgf(errorMsg.ErrorConvertingDataToGeoJson+" Error: %q", err)
		return
	}

	err = apiUtility.EncodeJsonData(w, aircraft)
	if err != nil {
		http.Error(w, errorMsg.ErrorEncodingJsonData, http.StatusInternalServerError)
		log.Error().Msgf(errorMsg.ErrorEncodingJsonData+": %q", err)
	}
}
