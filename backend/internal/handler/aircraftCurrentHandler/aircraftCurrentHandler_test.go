package aircraftCurrentHandler

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/global/geoJSON"
	"adsb-api/internal/global/models"
	"adsb-api/internal/utility/convert"
	"adsb-api/internal/utility/mock"
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
	global.InitTestEnvironment()
	m.Run()
}

func TestInvalidRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockRestService(ctrl)
	currentEndpoint := httptest.NewServer(CurrentAircraftHandler(mockSvc)) // Use mockSvc here
	defer currentEndpoint.Close()

	var endpoint = currentEndpoint.URL + global.AircraftCurrentPath

	tests := []struct {
		name, url, httpMethod, errorMsg string
		statusCode, length              int
		setup                           func(mockDB *mock.MockRestService)
	}{
		{
			name:       "Post request",
			url:        endpoint,
			httpMethod: http.MethodPost,
			statusCode: http.StatusMethodNotAllowed,
			errorMsg:   fmt.Sprintf(errorMsg.MethodNotSupported, http.MethodPost),
		},
		{
			name:       "Delete request",
			url:        endpoint,
			httpMethod: http.MethodDelete,
			statusCode: http.StatusMethodNotAllowed,
			errorMsg:   fmt.Sprintf(errorMsg.MethodNotSupported, http.MethodDelete),
		},
		{
			name:       "Database returns nil",
			url:        endpoint,
			httpMethod: http.MethodGet,
			statusCode: http.StatusInternalServerError,
			setup: func(mockSvc *mock.MockRestService) {
				mockSvc.EXPECT().GetCurrentAircraft().Return([]models.AircraftCurrentModel{}, errors.New("no new aircraft"))
			},
			errorMsg: errorMsg.ErrorRetrievingCurrentAircraft,
		},
		{
			name:       "Get request with too long URL",
			url:        endpoint + "endpoint/",
			httpMethod: http.MethodGet,
			statusCode: http.StatusBadRequest,
			errorMsg:   errorMsg.ErrorTongURL,
		},
		{
			name:       "Get request with invalid parameter",
			url:        endpoint + "?param=123",
			httpMethod: http.MethodGet,
			statusCode: http.StatusBadRequest,
			errorMsg:   errorMsg.ErrorInvalidQueryParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(mockSvc)
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

	mockSvc := mock.NewMockRestService(ctrl)
	currentEndpoint := httptest.NewServer(CurrentAircraftHandler(mockSvc)) // Use mockSvc here
	defer currentEndpoint.Close()

	var endpoint = currentEndpoint.URL + global.AircraftCurrentPath

	tests := []struct {
		name, url, httpMethod string
		statusCode            int
		mockData              []models.AircraftCurrentModel
		setup                 func(mockSvc *mock.MockRestService, mockData []models.AircraftCurrentModel)
	}{
		{
			name:       "Get request without parameters",
			url:        endpoint,
			httpMethod: http.MethodGet,
			statusCode: http.StatusOK,
			mockData:   testUtility.CreateMockAircraft(10),
			setup: func(mockSvc *mock.MockRestService, mockData []models.AircraftCurrentModel) {
				mockSvc.EXPECT().GetCurrentAircraft().Return(mockData, nil)
			},
		},
		{
			name:       "Get request with empty current_time_aircraft table",
			url:        endpoint,
			httpMethod: http.MethodGet,
			statusCode: http.StatusNoContent,
			setup: func(mockSvc *mock.MockRestService, mockData []models.AircraftCurrentModel) {
				mockSvc.EXPECT().GetCurrentAircraft().Return([]models.AircraftCurrentModel{}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.setup(mockSvc, tt.mockData)

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

			mockFeatureCollection, err := convert.CurrentModelToGeoJson(tt.mockData)

			assert.Equal(t, mockFeatureCollection, actual)
		})
	}
}
