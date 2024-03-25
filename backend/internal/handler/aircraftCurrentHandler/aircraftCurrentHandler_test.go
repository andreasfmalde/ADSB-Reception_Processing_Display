package aircraftCurrentHandler

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	errors2 "adsb-api/internal/global/errorMsg"
	"adsb-api/internal/global/geoJSON"
	"adsb-api/internal/global/models"
	"adsb-api/internal/utility/testUtility"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
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
	currentEndpoint := httptest.NewServer(CurrentAircraftHandler(mockDB)) // Use mockDB here
	defer currentEndpoint.Close()

	tests := []struct {
		name, url, httpMethod, errorMsg string
		statusCode, length              int
		setup                           func(mockDB *db.MockDatabase)
	}{
		{
			name:       "Post request",
			url:        currentEndpoint.URL + global.AircraftCurrentPath,
			httpMethod: http.MethodPost,
			statusCode: http.StatusMethodNotAllowed,
			errorMsg:   fmt.Sprintf(errors2.MethodNotSupported, http.MethodPost),
		},
		{
			name:       "Delete request",
			url:        currentEndpoint.URL + global.AircraftCurrentPath,
			httpMethod: http.MethodDelete,
			statusCode: http.StatusMethodNotAllowed,
			errorMsg:   fmt.Sprintf(errors2.MethodNotSupported, http.MethodDelete),
		},
		{
			name:       "Database returns nil",
			url:        currentEndpoint.URL + global.AircraftCurrentPath,
			httpMethod: http.MethodGet,
			statusCode: http.StatusInternalServerError,
			setup: func(mockDB *db.MockDatabase) {
				mockDB.EXPECT().GetCurrentAircraft().Return([]models.AircraftCurrentModel{}, errors.New("no new aircraft"))
			},
			errorMsg: errors2.ErrorRetrievingCurrentAircraft,
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
				t.Errorf("Test: %s. Error creating request with populateMethod %s and url %s: %s", tt.name, tt.httpMethod, tt.url, err.Error())
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

func TestValidRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := db.NewMockDatabase(ctrl)
	currentEndpoint := httptest.NewServer(CurrentAircraftHandler(mockDB)) // Use mockDB here
	defer currentEndpoint.Close()

	tests := []struct {
		name, url, httpMethod string
		statusCode            int
		mockData              []models.AircraftCurrentModel
		setup                 func(mockDB *db.MockDatabase, mockData []models.AircraftCurrentModel)
	}{
		{
			name:       "Get request without parameters",
			url:        currentEndpoint.URL + global.AircraftCurrentPath,
			httpMethod: http.MethodGet,
			statusCode: http.StatusOK,
			mockData:   testUtility.CreateMockAircraft(10),
			setup: func(mockDB *db.MockDatabase, mockData []models.AircraftCurrentModel) {
				mockDB.EXPECT().GetCurrentAircraft().Return(mockData, nil)
			},
		},
		{
			name:       "Get request with empty current_time_aircraft table",
			url:        currentEndpoint.URL + global.AircraftCurrentPath,
			httpMethod: http.MethodGet,
			statusCode: http.StatusNoContent,
			setup: func(mockDB *db.MockDatabase, mockData []models.AircraftCurrentModel) {
				mockDB.EXPECT().GetCurrentAircraft().Return([]models.AircraftCurrentModel{}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.setup(mockDB, tt.mockData)

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

			if tt.mockData == nil {
				return
			}

			var actual geoJSON.FeatureCollectionPoint
			err = json.NewDecoder(res.Body).Decode(&actual)
			if err != nil {
				t.Errorf("Test: %s. Error decoing response body: %s", tt.name, err.Error())
			}

			mockFeatureCollection, err := geoJSON.ConvertCurrentModelToGeoJson(tt.mockData)

			assert.Equal(t, mockFeatureCollection, actual)
		})
	}
}
