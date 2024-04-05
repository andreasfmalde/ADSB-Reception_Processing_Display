package aircraftHistory

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
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	global.InitTestEnvironment()
	m.Run()
}

func TestInvalidRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockRestService(ctrl)
	currentEndpoint := httptest.NewServer(HistoryAircraftHandler(mockSvc))
	defer currentEndpoint.Close()

	var endpoint = currentEndpoint.URL + global.AircraftHistoryPath

	tests := []struct {
		name, url, httpMethod, errorMsg string
		statusCode, length              int
		setup                           func(mockDB *mock.MockRestService)
	}{
		{
			name:       "Post request",
			httpMethod: http.MethodPost,
			url:        endpoint + "ABC123",
			statusCode: http.StatusMethodNotAllowed,
			errorMsg:   fmt.Sprintf(errorMsg.MethodNotSupported, http.MethodPost),
		},
		{
			name:       "Delete request",
			url:        endpoint + "ABC123",
			httpMethod: http.MethodDelete,
			statusCode: http.StatusMethodNotAllowed,
			errorMsg:   fmt.Sprintf(errorMsg.MethodNotSupported, http.MethodDelete),
		},
		{
			name:       "Database returns nil",
			url:        endpoint + "ABC123",
			httpMethod: http.MethodGet,
			statusCode: http.StatusInternalServerError,
			setup: func(mockSvc *mock.MockRestService) {
				mockSvc.EXPECT().GetAircraftHistoryByIcao("ABC123").Return([]models.AircraftHistoryModel{}, errors.New("expected error"))
			},
			errorMsg: errorMsg.ErrorRetrievingAircraftWithIcao + "ABC123",
		},
		{
			name:       "Get request with too long URL",
			url:        endpoint + "endpoint/endpoint/",
			httpMethod: http.MethodGet,
			statusCode: http.StatusBadRequest,
			errorMsg:   errorMsg.ErrorTongURL,
		},
		{
			name:       "Get request with empty icao",
			url:        endpoint + "/",
			httpMethod: http.MethodGet,
			statusCode: http.StatusBadRequest,
			errorMsg:   errorMsg.EmptyIcao,
		},
		{
			name:       "Get request with invalid parameter",
			url:        endpoint + "ABC123?param=123",
			httpMethod: http.MethodGet,
			statusCode: http.StatusBadRequest,
			errorMsg:   errorMsg.ErrorInvalidQueryParams + strings.Join(params, ", "),
		},
		{
			name:       "Get request with too many parameters",
			url:        endpoint + "ABC123?url=abc?param=ABC123",
			httpMethod: http.MethodGet,
			statusCode: http.StatusBadRequest,
			errorMsg:   errorMsg.ErrorInvalidQueryParams + strings.Join(params, ", "),
		},
		{
			name:       "Get request without icao or parameter",
			url:        endpoint,
			httpMethod: http.MethodGet,
			statusCode: http.StatusBadRequest,
			errorMsg:   errorMsg.EmptyIcao,
		},
		{
			name:       "Invalid query parameter 'hour'",
			url:        endpoint + "ABC123?hour=ABC123",
			statusCode: http.StatusBadRequest,
			errorMsg:   errorMsg.InvalidQueryParameterHour,
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

func TestValidRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockRestService(ctrl)
	currentEndpoint := httptest.NewServer(HistoryAircraftHandler(mockSvc))
	defer currentEndpoint.Close()

	var endpoint = currentEndpoint.URL + global.AircraftHistoryPath

	tests := []struct {
		name, url, httpMethod string
		statusCode            int
		mockData              []models.AircraftHistoryModel
		setup                 func(mockDB *mock.MockRestService, mockData []models.AircraftHistoryModel)
	}{
		{
			name:       "Get request with valid URl",
			url:        endpoint + "/ABC123",
			httpMethod: http.MethodGet,
			statusCode: http.StatusOK,
			mockData:   testUtility.CreateMockHistAircraft(10),
			setup: func(mockSvc *mock.MockRestService, mockData []models.AircraftHistoryModel) {
				mockSvc.EXPECT().GetAircraftHistoryByIcao("ABC123").Return(mockData, nil)
			},
		},
		{
			name:       "Get request with valid URl but unnecessary slashes",
			url:        currentEndpoint.URL + "/../" + global.AircraftHistoryPath + "//../././/ABC123",
			httpMethod: http.MethodGet,
			statusCode: http.StatusOK,
			mockData:   testUtility.CreateMockHistAircraft(10),
			setup: func(mockSvc *mock.MockRestService, mockData []models.AircraftHistoryModel) {
				mockSvc.EXPECT().GetAircraftHistoryByIcao("ABC123").Return(mockData, nil)
			},
		},
		{
			name:       "Get request with no history",
			url:        endpoint + "ABC123",
			httpMethod: http.MethodGet,
			statusCode: http.StatusNoContent,
			setup: func(mockSvc *mock.MockRestService, mockData []models.AircraftHistoryModel) {
				mockSvc.EXPECT().GetAircraftHistoryByIcao("ABC123").Return([]models.AircraftHistoryModel{}, nil)
			},
		},
		{
			name:       "history with too few coordinates",
			url:        endpoint + "ABC123",
			httpMethod: http.MethodGet,
			statusCode: http.StatusNoContent,
			setup: func(mockDB *mock.MockRestService, mockData []models.AircraftHistoryModel) {
				mockSvc.EXPECT().GetAircraftHistoryByIcao("ABC123").Return([]models.AircraftHistoryModel{}, nil)
			},
		},
		{
			name:       "Get request with valid icao and valid query parameter 'hour'",
			url:        endpoint + "ABC123?hour=2",
			statusCode: http.StatusOK,
			mockData:   testUtility.CreateMockHistAircraft(10),
			setup: func(mockDB *mock.MockRestService, mockData []models.AircraftHistoryModel) {
				mockSvc.EXPECT().GetAircraftHistoryByIcaoFilterByTimestamp("ABC123", 2).Return(mockData, nil)
			},
		},
		{
			name:       "Get request with valid query parameter 'hour' but no history",
			url:        endpoint + "ABC123?hour=1000",
			statusCode: http.StatusNoContent,
			setup: func(mockDB *mock.MockRestService, mockData []models.AircraftHistoryModel) {
				mockSvc.EXPECT().GetAircraftHistoryByIcaoFilterByTimestamp("ABC123", 1000).Return([]models.AircraftHistoryModel{}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(mockSvc, tt.mockData)
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

			if tt.mockData != nil {
				var actual geoJSON.FeatureCollectionLineString
				_ = json.NewDecoder(res.Body).Decode(&actual)

				mockFeatureCollection, err := convert.HistoryModelToGeoJson(tt.mockData)
				if err != nil {
					t.Fatalf("error converting from history model to geo json")
				}

				assert.Equal(t, mockFeatureCollection, actual)
			}
		})
	}
}
