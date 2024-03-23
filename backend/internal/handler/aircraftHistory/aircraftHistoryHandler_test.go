package aircraftHistory

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/global/geoJSON"
	"adsb-api/internal/utility/testUtility"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	global.InitTestEnv()
	m.Run()
}

func TestInvalidRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := db.NewMockDatabase(ctrl)
	currentEndpoint := httptest.NewServer(HistoryAircraftHandler(mockDB))
	defer currentEndpoint.Close()

	var endpoint = currentEndpoint.URL + global.AircraftHistoryPath

	tests := []struct {
		name, url, httpMethod, errorMsg string
		statusCode, length              int
		setup                           func(mockDB *db.MockDatabase)
	}{
		{
			name:       "Post request",
			httpMethod: http.MethodPost,
			url:        endpoint + "?icao=ABC123",
			statusCode: http.StatusMethodNotAllowed,
			errorMsg:   fmt.Sprintf(global.MethodNotSupported, http.MethodPost),
		},
		{
			name:       "Delete request",
			url:        endpoint + "?icao=ABC123",
			httpMethod: http.MethodDelete,
			statusCode: http.StatusMethodNotAllowed,
			errorMsg:   fmt.Sprintf(global.MethodNotSupported, http.MethodDelete),
		},
		{
			name:       "Database returns nil",
			url:        endpoint + "?icao=ABC123",
			httpMethod: http.MethodGet,
			statusCode: http.StatusInternalServerError,
			setup: func(mockDB *db.MockDatabase) {
				mockDB.EXPECT().GetHistoryByIcao("ABC123").Return([]global.AircraftHistoryModel{}, errors.New("expected error"))
			},
			errorMsg: global.ErrorRetrievingAircraftWithIcao + "ABC123",
		},
		{
			name:       "Get request with too long URL",
			url:        endpoint + "endpoint/",
			httpMethod: http.MethodGet,
			statusCode: http.StatusBadRequest,
			errorMsg:   global.ErrorTongURL,
		},
		{
			name:       "Get request with empty icao",
			url:        endpoint + "?icao=",
			httpMethod: http.MethodGet,
			statusCode: http.StatusBadRequest,
			errorMsg:   global.ErrorInvalidQueryParams + strings.Join(params, ", "),
		},
		{
			name:       "Get request with invalid parameter",
			url:        endpoint + "?param=123",
			httpMethod: http.MethodGet,
			statusCode: http.StatusBadRequest,
			errorMsg:   global.ErrorInvalidQueryParams + strings.Join(params, ", "),
		},
		{
			name:       "Get request with too many parameters",
			url:        endpoint + "?icao=ABC123&param=ABC123",
			httpMethod: http.MethodGet,
			statusCode: http.StatusBadRequest,
			errorMsg:   global.ErrorInvalidQueryParams + strings.Join(params, ", "),
		},
		{
			name:       "Get request without parameters",
			url:        endpoint,
			httpMethod: http.MethodGet,
			statusCode: http.StatusBadRequest,
			errorMsg:   global.ErrorInvalidQueryParams + strings.Join(params, ", "),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(mockDB)
			}

			client := &http.Client{}
			req, err := http.NewRequest(tt.httpMethod, tt.url, nil)
			if err != nil {
				t.Errorf("Test: %s. Error creating request with populateMethod %s and endpoint %s: %s", tt.name, tt.httpMethod, tt.url, err.Error())
			}
			res, err := client.Do(req)
			if err != nil {
				t.Errorf("Test: %s. Error executing %s request: %s", tt.name, tt.httpMethod, err.Error())
			}

			assert.Equal(t, tt.statusCode, res.StatusCode)

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Test: %s. Error reading response body: %s", tt.name, err.Error())
			}
			assert.Equal(t, tt.errorMsg+"\n", string(body))

		})
	}
}

// TODO: Check if result follows GeoJson standard
func TestValidRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := db.NewMockDatabase(ctrl)
	currentEndpoint := httptest.NewServer(HistoryAircraftHandler(mockDB))
	defer currentEndpoint.Close()

	var endpoint = currentEndpoint.URL + global.AircraftHistoryPath

	tests := []struct {
		name, url, httpMethod string
		statusCode            int
		mockData              []global.AircraftHistoryModel
		setup                 func(mockDB *db.MockDatabase, mockData []global.AircraftHistoryModel)
	}{
		{
			name:       "Get request with valid URl",
			url:        endpoint + "?icao=ABC123",
			httpMethod: http.MethodGet,
			statusCode: http.StatusOK,
			mockData:   testUtility.CreateMockHistAircraft(10),
			setup: func(mockDB *db.MockDatabase, mockData []global.AircraftHistoryModel) {
				mockDB.EXPECT().GetHistoryByIcao("ABC123").Return(mockData, nil)
			},
		},
		{
			name:       "Get request with valid URl but unnecessary slashes",
			url:        currentEndpoint.URL + "/../" + global.AircraftHistoryPath + "//../././/?icao=ABC123",
			httpMethod: http.MethodGet,
			statusCode: http.StatusOK,
			mockData:   testUtility.CreateMockHistAircraft(10),
			setup: func(mockDB *db.MockDatabase, mockData []global.AircraftHistoryModel) {
				mockDB.EXPECT().GetHistoryByIcao("ABC123").Return(mockData, nil)
			},
		},
		{
			name:       "Get request with no history",
			url:        endpoint + "?icao=ABC123",
			httpMethod: http.MethodGet,
			statusCode: http.StatusNoContent,
			setup: func(mockDB *db.MockDatabase, mockData []global.AircraftHistoryModel) {
				mockDB.EXPECT().GetHistoryByIcao("ABC123").Return([]global.AircraftHistoryModel{}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(mockDB, tt.mockData)
			}

			client := &http.Client{}
			req, err := http.NewRequest(tt.httpMethod, tt.url, nil)
			if err != nil {
				t.Errorf("Test: %s. Error creating request with populateMethod %s and url %s: %s", tt.name, tt.httpMethod, tt.url, err.Error())
			}
			res, err := client.Do(req)
			if err != nil {
				t.Errorf("Test: %s. Error executing %s request: %s", tt.name, tt.httpMethod, err.Error())
			}

			assert.Equal(t, tt.statusCode, res.StatusCode)

			var actual geoJSON.FeatureCollectionLineString
			_ = json.NewDecoder(res.Body).Decode(&actual)

			mockFeatureCollection, err := geoJSON.ConvertHistoryModelToGeoJson(tt.mockData)

			assert.Equal(t, mockFeatureCollection, actual)
		})
	}
}
